package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"financetracker/config"
	"financetracker/middleware"
	"financetracker/models"
)

// CreateTransactionRequest represents the transaction creation request
type CreateTransactionRequest struct {
	AccountID   uuid.UUID              `json:"account_id" validate:"required"`
	CategoryID  *uuid.UUID             `json:"category_id"`
	Type        models.TransactionType `json:"type" validate:"required,oneof=income expense transfer"`
	Amount      float64                `json:"amount" validate:"required,gt=0"`
	Currency    string                 `json:"currency" validate:"len=3"`
	Description string                 `json:"description" validate:"required,min=2,max=255"`
	Notes       string                 `json:"notes"`
	Date        time.Time              `json:"date" validate:"required"`
	ToAccountID *uuid.UUID             `json:"to_account_id"`
	Tags        []string               `json:"tags"`
}

// UpdateTransactionRequest represents the transaction update request
type UpdateTransactionRequest struct {
	CategoryID  *uuid.UUID              `json:"category_id"`
	Type        *models.TransactionType `json:"type" validate:"omitempty,oneof=income expense transfer"`
	Amount      *float64                `json:"amount" validate:"omitempty,gt=0"`
	Currency    string                  `json:"currency" validate:"omitempty,len=3"`
	Description string                  `json:"description" validate:"omitempty,min=2,max=255"`
	Notes       string                  `json:"notes"`
	Date        *time.Time              `json:"date"`
	ToAccountID *uuid.UUID              `json:"to_account_id"`
	Tags        []string                `json:"tags"`
}

// GetTransactions returns paginated transactions with filters
func GetTransactions(c *gin.Context) {
	userID, err := middleware.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse query parameters
	filter := models.GetDefaultFilter()
	filter.UserID = userID

	if accountID := c.Query("account_id"); accountID != "" {
		if id, err := uuid.Parse(accountID); err == nil {
			filter.AccountID = &id
		}
	}

	if categoryID := c.Query("category_id"); categoryID != "" {
		if id, err := uuid.Parse(categoryID); err == nil {
			filter.CategoryID = &id
		}
	}

	if transactionType := c.Query("type"); transactionType != "" {
		t := models.TransactionType(transactionType)
		filter.Type = &t
	}

	if dateFrom := c.Query("date_from"); dateFrom != "" {
		if date, err := time.Parse("2006-01-02", dateFrom); err == nil {
			filter.DateFrom = &date
		}
	}

	if dateTo := c.Query("date_to"); dateTo != "" {
		if date, err := time.Parse("2006-01-02", dateTo); err == nil {
			filter.DateTo = &date
		}
	}

	if amountMin := c.Query("amount_min"); amountMin != "" {
		if amount, err := strconv.ParseFloat(amountMin, 64); err == nil {
			filter.AmountMin = &amount
		}
	}

	if amountMax := c.Query("amount_max"); amountMax != "" {
		if amount, err := strconv.ParseFloat(amountMax, 64); err == nil {
			filter.AmountMax = &amount
		}
	}

	if search := c.Query("search"); search != "" {
		filter.Search = search
	}

	if page := c.Query("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			filter.Page = p
		}
	}

	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 && l <= 100 {
			filter.Limit = l
		}
	}

	if sortBy := c.Query("sort_by"); sortBy != "" {
		filter.SortBy = sortBy
	}

	if sortOrder := c.Query("sort_order"); sortOrder == "asc" || sortOrder == "desc" {
		filter.SortOrder = sortOrder
	}

	// Build query
	query := config.DB.Model(&models.Transaction{}).
		Preload("Category").
		Preload("Account").
		Preload("ToAccount").
		Where("user_id = ?", filter.UserID)

	// Apply filters
	if filter.AccountID != nil {
		query = query.Where("account_id = ?", *filter.AccountID)
	}

	if filter.CategoryID != nil {
		query = query.Where("category_id = ?", *filter.CategoryID)
	}

	if filter.Type != nil {
		query = query.Where("type = ?", *filter.Type)
	}

	if filter.DateFrom != nil {
		query = query.Where("date >= ?", *filter.DateFrom)
	}

	if filter.DateTo != nil {
		query = query.Where("date <= ?", *filter.DateTo)
	}

	if filter.AmountMin != nil {
		query = query.Where("amount >= ?", *filter.AmountMin)
	}

	if filter.AmountMax != nil {
		query = query.Where("amount <= ?", *filter.AmountMax)
	}

	if filter.Search != "" {
		query = query.Where("description ILIKE ? OR notes ILIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
	}

	// Count total records
	var total int64
	query.Count(&total)

	// Apply pagination and sorting
	offset := (filter.Page - 1) * filter.Limit
	orderClause := filter.SortBy + " " + filter.SortOrder

	var transactions []models.Transaction
	if err := query.Order(orderClause).Offset(offset).Limit(filter.Limit).Find(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transactions"})
		return
	}

	// Calculate pagination info
	totalPages := (total + int64(filter.Limit) - 1) / int64(filter.Limit)

	c.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
		"pagination": gin.H{
			"page":        filter.Page,
			"limit":       filter.Limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// CreateTransaction creates a new transaction
func CreateTransaction(c *gin.Context) {
	userID, err := middleware.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req CreateTransactionRequest
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

	// Verify account ownership
	var account models.Account
	if err := config.DB.Where("id = ? AND user_id = ?", req.AccountID, userID).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
		return
	}

	// Verify category ownership if provided
	if req.CategoryID != nil {
		var category models.Category
		if err := config.DB.Where("id = ? AND user_id = ?", *req.CategoryID, userID).First(&category).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}
	}

	// Verify destination account for transfers
	if req.Type == models.TransactionTypeTransfer {
		if req.ToAccountID == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Destination account is required for transfers"})
			return
		}

		var toAccount models.Account
		if err := config.DB.Where("id = ? AND user_id = ?", *req.ToAccountID, userID).First(&toAccount).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Destination account not found"})
			return
		}
	}

	// Set currency from account if not provided
	if req.Currency == "" {
		req.Currency = account.Currency
	}

	// Create transaction
	transaction := models.Transaction{
		UserID:      userID,
		AccountID:   req.AccountID,
		CategoryID:  req.CategoryID,
		Type:        req.Type,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Description: req.Description,
		Notes:       req.Notes,
		Date:        req.Date,
		ToAccountID: req.ToAccountID,
	}

	if err := config.DB.Create(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	// Load relationships
	config.DB.Preload("Category").Preload("Account").Preload("ToAccount").First(&transaction, transaction.ID)

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Transaction created successfully",
		"transaction": transaction,
	})
}

// GetTransaction returns a specific transaction
func GetTransaction(c *gin.Context) {
	userID, err := middleware.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	transactionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	var transaction models.Transaction
	if err := config.DB.Preload("Category").Preload("Account").Preload("ToAccount").
		Where("id = ? AND user_id = ?", transactionID, userID).First(&transaction).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transaction"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"transaction": transaction})
}

// UpdateTransaction updates an existing transaction
func UpdateTransaction(c *gin.Context) {
	userID, err := middleware.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	transactionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	var req UpdateTransactionRequest
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

	// Find transaction
	var transaction models.Transaction
	if err := config.DB.Where("id = ? AND user_id = ?", transactionID, userID).First(&transaction).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transaction"})
		}
		return
	}

	// Build updates map
	updates := make(map[string]interface{})

	if req.CategoryID != nil {
		// Verify category ownership
		var category models.Category
		if err := config.DB.Where("id = ? AND user_id = ?", *req.CategoryID, userID).First(&category).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
			return
		}
		updates["category_id"] = *req.CategoryID
	}

	if req.Type != nil {
		updates["type"] = *req.Type
	}

	if req.Amount != nil {
		updates["amount"] = *req.Amount
	}

	if req.Currency != "" {
		updates["currency"] = req.Currency
	}

	if req.Description != "" {
		updates["description"] = req.Description
	}

	if req.Notes != "" {
		updates["notes"] = req.Notes
	}

	if req.Date != nil {
		updates["date"] = *req.Date
	}

	if req.ToAccountID != nil {
		// Verify destination account ownership
		var toAccount models.Account
		if err := config.DB.Where("id = ? AND user_id = ?", *req.ToAccountID, userID).First(&toAccount).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Destination account not found"})
			return
		}
		updates["to_account_id"] = *req.ToAccountID
	}

	// Update transaction
	if err := config.DB.Model(&transaction).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaction"})
		return
	}

	// Reload with relationships
	config.DB.Preload("Category").Preload("Account").Preload("ToAccount").First(&transaction, transaction.ID)

	c.JSON(http.StatusOK, gin.H{
		"message":     "Transaction updated successfully",
		"transaction": transaction,
	})
}

// DeleteTransaction deletes a transaction
func DeleteTransaction(c *gin.Context) {
	userID, err := middleware.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	transactionID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	// Find and delete transaction
	var transaction models.Transaction
	if err := config.DB.Where("id = ? AND user_id = ?", transactionID, userID).First(&transaction).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transaction"})
		}
		return
	}

	if err := config.DB.Delete(&transaction).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted successfully"})
}

// CreateBulkTransactions creates multiple transactions at once
func CreateBulkTransactions(c *gin.Context) {
	userID, err := middleware.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		Transactions []CreateTransactionRequest `json:"transactions" validate:"required,min=1,max=100"`
	}

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

	var transactions []models.Transaction
	var errors []string

	// Process each transaction
	for i, transactionReq := range req.Transactions {
		// Validate individual transaction
		if err := validate.Struct(transactionReq); err != nil {
			errors = append(errors, "Transaction "+strconv.Itoa(i+1)+": "+err.Error())
			continue
		}

		// Verify account ownership
		var account models.Account
		if err := config.DB.Where("id = ? AND user_id = ?", transactionReq.AccountID, userID).First(&account).Error; err != nil {
			errors = append(errors, "Transaction "+strconv.Itoa(i+1)+": Account not found")
			continue
		}

		// Set currency from account if not provided
		if transactionReq.Currency == "" {
			transactionReq.Currency = account.Currency
		}

		transaction := models.Transaction{
			UserID:      userID,
			AccountID:   transactionReq.AccountID,
			CategoryID:  transactionReq.CategoryID,
			Type:        transactionReq.Type,
			Amount:      transactionReq.Amount,
			Currency:    transactionReq.Currency,
			Description: transactionReq.Description,
			Notes:       transactionReq.Notes,
			Date:        transactionReq.Date,
			ToAccountID: transactionReq.ToAccountID,
		}

		transactions = append(transactions, transaction)
	}

	if len(errors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation errors",
			"details": errors,
		})
		return
	}

	// Create all transactions in a transaction
	if err := config.DB.Create(&transactions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transactions"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":      "Transactions created successfully",
		"count":        len(transactions),
		"transactions": transactions,
	})
}

// ExportTransactions exports transactions to CSV
func ExportTransactions(c *gin.Context) {
	// This would implement CSV export functionality
	// For now, return a placeholder response
	c.JSON(http.StatusOK, gin.H{
		"message": "Export functionality will be implemented",
		"format":  "CSV",
	})
}
