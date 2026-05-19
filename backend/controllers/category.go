package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"financetracker/config"
	"financetracker/middleware"
	"financetracker/models"
)

// CreateCategoryRequest represents the category creation request
type CreateCategoryRequest struct {
	Name        string              `json:"name" validate:"required,min=2,max=100"`
	Type        models.CategoryType `json:"type" validate:"required,oneof=income expense"`
	Color       string              `json:"color" validate:"hexcolor"`
	Icon        string              `json:"icon"`
	Description string              `json:"description"`
}

// UpdateCategoryRequest represents the category update request
type UpdateCategoryRequest struct {
	Name        string               `json:"name" validate:"omitempty,min=2,max=100"`
	Type        *models.CategoryType `json:"type" validate:"omitempty,oneof=income expense"`
	Color       string               `json:"color" validate:"omitempty,hexcolor"`
	Icon        string               `json:"icon"`
	Description string               `json:"description"`
	IsActive    *bool                `json:"is_active"`
}

// GetCategories returns all categories for the current user
func GetCategories(c *gin.Context) {
	userID, err := middleware.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var categories []models.Category
	query := config.DB.Where("user_id = ?", userID)

	// Filter by type if provided
	if categoryType := c.Query("type"); categoryType != "" {
		query = query.Where("type = ?", categoryType)
	}

	// Filter by active status
	if isActive := c.Query("is_active"); isActive == "true" {
		query = query.Where("is_active = ?", true)
	} else if isActive == "false" {
		query = query.Where("is_active = ?", false)
	}

	if err := query.Order("name ASC").Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"categories": categories})
}

// CreateCategory creates a new category
func CreateCategory(c *gin.Context) {
	userID, err := middleware.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": err.Error(),
		})
		return
	}

	// Check if category name already exists for this user
	var existingCategory models.Category
	if err := config.DB.Where("user_id = ? AND name = ?", userID, req.Name).First(&existingCategory).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Category with this name already exists"})
		return
	}

	// Set defaults
	if req.Color == "" {
		req.Color = "#6366f1"
	}
	if req.Icon == "" {
		req.Icon = "💰"
	}

	category := models.Category{
		UserID:      userID,
		Name:        req.Name,
		Type:        req.Type,
		Color:       req.Color,
		Icon:        req.Icon,
		Description: req.Description,
		IsDefault:   false,
		IsActive:    true,
	}

	if err := config.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Category created successfully",
		"category": category,
	})
}

// GetCategory returns a specific category
func GetCategory(c *gin.Context) {
	userID, err := middleware.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var category models.Category
	if err := config.DB.Where("id = ? AND user_id = ?", categoryID, userID).First(&category).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch category"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"category": category})
}

// UpdateCategory updates an existing category
func UpdateCategory(c *gin.Context) {
	userID, err := middleware.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": err.Error(),
		})
		return
	}

	// Find category
	var category models.Category
	if err := config.DB.Where("id = ? AND user_id = ?", categoryID, userID).First(&category).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch category"})
		}
		return
	}

	// Prevent updating default categories
	if category.IsDefault {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot update default categories"})
		return
	}

	// Build updates map
	updates := make(map[string]interface{})

	if req.Name != "" {
		// Check if new name already exists
		var existingCategory models.Category
		if err := config.DB.Where("user_id = ? AND name = ? AND id != ?", userID, req.Name, categoryID).First(&existingCategory).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Category with this name already exists"})
			return
		}
		updates["name"] = req.Name
	}

	if req.Type != nil {
		updates["type"] = *req.Type
	}

	if req.Color != "" {
		updates["color"] = req.Color
	}

	if req.Icon != "" {
		updates["icon"] = req.Icon
	}

	if req.Description != "" {
		updates["description"] = req.Description
	}

	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	// Update category
	if err := config.DB.Model(&category).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	// Reload category
	config.DB.First(&category, category.ID)

	c.JSON(http.StatusOK, gin.H{
		"message":  "Category updated successfully",
		"category": category,
	})
}

// DeleteCategory deletes a category
func DeleteCategory(c *gin.Context) {
	userID, err := middleware.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	categoryID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	// Find category
	var category models.Category
	if err := config.DB.Where("id = ? AND user_id = ?", categoryID, userID).First(&category).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch category"})
		}
		return
	}

	// Prevent deleting default categories
	if category.IsDefault {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot delete default categories"})
		return
	}

	// Check if category is being used by transactions
	var transactionCount int64
	config.DB.Model(&models.Transaction{}).Where("category_id = ?", categoryID).Count(&transactionCount)

	if transactionCount > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"error":             "Cannot delete category that is being used by transactions",
			"transaction_count": transactionCount,
		})
		return
	}

	// Check if category is being used by budgets
	var budgetCount int64
	config.DB.Model(&models.Budget{}).Where("category_id = ?", categoryID).Count(&budgetCount)

	if budgetCount > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"error":        "Cannot delete category that is being used by budgets",
			"budget_count": budgetCount,
		})
		return
	}

	// Delete category
	if err := config.DB.Delete(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}
