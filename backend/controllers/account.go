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

// CreateAccountRequest represents the account creation request
type CreateAccountRequest struct {
	Name           string             `json:"name" validate:"required,min=2,max=100"`
	Type           models.AccountType `json:"type" validate:"required,oneof=checking savings credit cash investment"`
	Balance        float64            `json:"balance" validate:"numeric"`
	Currency       string             `json:"currency" validate:"len=3"`
	Color          string             `json:"color" validate:"hexcolor"`
	Icon           string             `json:"icon"`
	Description    string             `json:"description"`
	IncludeInTotal bool               `json:"include_in_total"`
}

// UpdateAccountRequest represents the account update request
type UpdateAccountRequest struct {
	Name           string              `json:"name" validate:"omitempty,min=2,max=100"`
	Type           *models.AccountType `json:"type" validate:"omitempty,oneof=checking savings credit cash investment"`
	Currency       string              `json:"currency" validate:"omitempty,len=3"`
	Color          string              `json:"color" validate:"omitempty,hexcolor"`
	Icon           string              `json:"icon"`
	Description    string              `json:"description"`
	IsActive       *bool               `json:"is_active"`
	IncludeInTotal *bool               `json:"include_in_total"`
}

// GetAccounts returns all accounts for the current user
func GetAccounts(c *gin.Context) {
	userID, err := middleware.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var accounts []models.Account
	query := config.DB.Where("user_id = ?", userID)

	// Filter by type if provided
	if accountType := c.Query("type"); accountType != "" {
		query = query.Where("type = ?", accountType)
	}

	// Filter by active status
	if isActive := c.Query("is_active"); isActive == "true" {
		query = query.Where("is_active = ?", true)
	} else if isActive == "false" {
		query = query.Where("is_active = ?", false)
	}

	// Filter by include_in_total
	if includeInTotal := c.Query("include_in_total"); includeInTotal == "true" {
		query = query.Where("include_in_total = ?", true)
	} else if includeInTotal == "false" {
		query = query.Where("include_in_total = ?", false)
	}

	if err := query.Order("name ASC").Find(&accounts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch accounts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"accounts": accounts})
}

// CreateAccount creates a new account
func CreateAccount(c *gin.Context) {
	userID, err := middleware.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req CreateAccountRequest
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

	// Check if account name already exists for this user
	var existingAccount models.Account
	if err := config.DB.Where("user_id = ? AND name = ?", userID, req.Name).First(&existingAccount).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Account with this name already exists"})
		return
	}

	// Get user's default currency if not provided
	if req.Currency == "" {
		var user models.User
		if err := config.DB.First(&user, userID).Error; err == nil {
			req.Currency = user.Currency
		} else {
			req.Currency = "USD"
		}
	}

	// Set defaults
	if req.Color == "" {
		req.Color = "#6366f1"
	}
	if req.Icon == "" {
		req.Icon = "🏦"
	}

	account := models.Account{
		UserID:         userID,
		Name:           req.Name,
		Type:           req.Type,
		Balance:        req.Balance,
		Currency:       req.Currency,
		Color:          req.Color,
		Icon:           req.Icon,
		Description:    req.Description,
		IsActive:       true,
		IncludeInTotal: req.IncludeInTotal,
	}

	if err := config.DB.Create(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create account"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Account created successfully",
		"account": account,
	})
}

// GetAccount returns a specific account
func GetAccount(c *gin.Context) {
	userID, err := middleware.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	accountID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	var account models.Account
	if err := config.DB.Where("id = ? AND user_id = ?", accountID, userID).First(&account).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch account"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"account": account})
}

// UpdateAccount updates an existing account
func UpdateAccount(c *gin.Context) {
	userID, err := middleware.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	accountID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	var req UpdateAccountRequest
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

	// Find account
	var account models.Account
	if err := config.DB.Where("id = ? AND user_id = ?", accountID, userID).First(&account).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch account"})
		}
		return
	}

	// Build updates map
	updates := make(map[string]interface{})

	if req.Name != "" {
		// Check if new name already exists
		var existingAccount models.Account
		if err := config.DB.Where("user_id = ? AND name = ? AND id != ?", userID, req.Name, accountID).First(&existingAccount).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Account with this name already exists"})
			return
		}
		updates["name"] = req.Name
	}

	if req.Type != nil {
		updates["type"] = *req.Type
	}

	if req.Currency != "" {
		updates["currency"] = req.Currency
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

	if req.IncludeInTotal != nil {
		updates["include_in_total"] = *req.IncludeInTotal
	}

	// Update account
	if err := config.DB.Model(&account).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update account"})
		return
	}

	// Reload account
	config.DB.First(&account, account.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Account updated successfully",
		"account": account,
	})
}

// DeleteAccount deletes an account
func DeleteAccount(c *gin.Context) {
	userID, err := middleware.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	accountID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	// Find account
	var account models.Account
	if err := config.DB.Where("id = ? AND user_id = ?", accountID, userID).First(&account).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch account"})
		}
		return
	}

	// Check if account is being used by transactions
	var transactionCount int64
	config.DB.Model(&models.Transaction{}).Where("account_id = ? OR to_account_id = ?", accountID, accountID).Count(&transactionCount)

	if transactionCount > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"error":             "Cannot delete account that has transactions",
			"transaction_count": transactionCount,
		})
		return
	}

	// Delete account
	if err := config.DB.Delete(&account).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete account"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}

// GetAccountSummary returns account summary with additional statistics
func GetAccountSummary(c *gin.Context) {
	userID, err := middleware.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	accountID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account ID"})
		return
	}

	// Find account
	var account models.Account
	if err := config.DB.Where("id = ? AND user_id = ?", accountID, userID).First(&account).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch account"})
		}
		return
	}

	// Calculate summary statistics
	var totalIncome, totalExpense float64
	var transactionCount int64
	var lastTransactionDate *string

	// Total income
	config.DB.Model(&models.Transaction{}).
		Where("account_id = ? AND type = ?", accountID, models.TransactionTypeIncome).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalIncome)

	// Total expense
	config.DB.Model(&models.Transaction{}).
		Where("account_id = ? AND type = ?", accountID, models.TransactionTypeExpense).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalExpense)

	// Transaction count
	config.DB.Model(&models.Transaction{}).
		Where("account_id = ?", accountID).
		Count(&transactionCount)

	// Last transaction date
	var lastTransaction models.Transaction
	if err := config.DB.Where("account_id = ?", accountID).
		Order("date DESC").
		First(&lastTransaction).Error; err == nil {
		dateStr := lastTransaction.Date.Format("2006-01-02T15:04:05Z")
		lastTransactionDate = &dateStr
	}

	summary := models.AccountSummary{
		Account:           account,
		TotalIncome:       totalIncome,
		TotalExpense:      totalExpense,
		TransactionCount:  transactionCount,
		LastTransactionAt: nil,
	}

	if lastTransactionDate != nil {
		summary.LastTransactionAt = lastTransactionDate
	}

	c.JSON(http.StatusOK, gin.H{"account_summary": summary})
}
