package controllers

import (
	"net/http"
	"time"

	"financetracker/config"
	"financetracker/middleware"
	"financetracker/models"

	"github.com/gin-gonic/gin"
)

// FinancialSummary represents overall financial summary
type FinancialSummary struct {
	TotalIncome       float64 `json:"total_income"`
	TotalExpense      float64 `json:"total_expense"`
	NetIncome         float64 `json:"net_income"`
	TotalAccounts     int64   `json:"total_accounts"`
	TotalTransactions int64   `json:"total_transactions"`
	PeriodStart       string  `json:"period_start"`
	PeriodEnd         string  `json:"period_end"`
	Currency          string  `json:"currency"`
}

// CategoryAnalysis represents spending/income analysis by category
type CategoryAnalysis struct {
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Amount       float64 `json:"amount"`
	Percentage   float64 `json:"percentage"`
	Color        string  `json:"color"`
	Icon         string  `json:"icon"`
}

// ExpenseAnalysis represents expense breakdown
type ExpenseAnalysis struct {
	TotalExpense float64            `json:"total_expense"`
	Categories   []CategoryAnalysis `json:"categories"`
	PeriodStart  string             `json:"period_start"`
	PeriodEnd    string             `json:"period_end"`
}

// IncomeAnalysis represents income breakdown
type IncomeAnalysis struct {
	TotalIncome float64            `json:"total_income"`
	Categories  []CategoryAnalysis `json:"categories"`
	PeriodStart string             `json:"period_start"`
	PeriodEnd   string             `json:"period_end"`
}

// SpendingTrend represents spending trend data point
type SpendingTrend struct {
	Date    string  `json:"date"`
	Income  float64 `json:"income"`
	Expense float64 `json:"expense"`
	Net     float64 `json:"net"`
}

// SpendingTrends represents spending trends over time
type SpendingTrends struct {
	Trends      []SpendingTrend `json:"trends"`
	Period      string          `json:"period"`
	PeriodStart string          `json:"period_start"`
	PeriodEnd   string          `json:"period_end"`
}

// GetFinancialSummary returns overall financial summary
func GetFinancialSummary(c *gin.Context) {
	userID, err := middleware.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse date range
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	// Default to current month if no dates provided
	now := time.Now()
	if dateFrom == "" {
		dateFrom = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")
	}
	if dateTo == "" {
		dateTo = now.Format("2006-01-02")
	}

	// Get user currency
	var user models.User
	config.DB.First(&user, userID)

	// Calculate totals
	var totalIncome, totalExpense float64
	var totalAccounts, totalTransactions int64

	// Total income
	config.DB.Model(&models.Transaction{}).
		Where("user_id = ? AND type = ? AND date >= ? AND date <= ?",
			userID, models.TransactionTypeIncome, dateFrom, dateTo).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalIncome)

	// Total expense
	config.DB.Model(&models.Transaction{}).
		Where("user_id = ? AND type = ? AND date >= ? AND date <= ?",
			userID, models.TransactionTypeExpense, dateFrom, dateTo).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalExpense)

	// Total accounts
	config.DB.Model(&models.Account{}).
		Where("user_id = ? AND is_active = ?", userID, true).
		Count(&totalAccounts)

	// Total transactions in period
	config.DB.Model(&models.Transaction{}).
		Where("user_id = ? AND date >= ? AND date <= ?", userID, dateFrom, dateTo).
		Count(&totalTransactions)

	summary := FinancialSummary{
		TotalIncome:       totalIncome,
		TotalExpense:      totalExpense,
		NetIncome:         totalIncome - totalExpense,
		TotalAccounts:     totalAccounts,
		TotalTransactions: totalTransactions,
		PeriodStart:       dateFrom,
		PeriodEnd:         dateTo,
		Currency:          user.Currency,
	}

	c.JSON(http.StatusOK, gin.H{"summary": summary})
}

// GetExpenseAnalysis returns expense breakdown by category
func GetExpenseAnalysis(c *gin.Context) {
	userID, err := middleware.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse date range
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	now := time.Now()
	if dateFrom == "" {
		dateFrom = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")
	}
	if dateTo == "" {
		dateTo = now.Format("2006-01-02")
	}

	// Get total expense
	var totalExpense float64
	config.DB.Model(&models.Transaction{}).
		Where("user_id = ? AND type = ? AND date >= ? AND date <= ?",
			userID, models.TransactionTypeExpense, dateFrom, dateTo).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalExpense)

	// Get expense by category
	var results []struct {
		CategoryID   string  `json:"category_id"`
		CategoryName string  `json:"category_name"`
		Amount       float64 `json:"amount"`
		Color        string  `json:"color"`
		Icon         string  `json:"icon"`
	}

	config.DB.Table("transactions t").
		Select("COALESCE(c.id::text, '') as category_id, COALESCE(c.name, 'Uncategorized') as category_name, SUM(t.amount) as amount, COALESCE(c.color, '#6b7280') as color, COALESCE(c.icon, '💸') as icon").
		Joins("LEFT JOIN categories c ON t.category_id = c.id").
		Where("t.user_id = ? AND t.type = ? AND t.date >= ? AND t.date <= ?",
			userID, models.TransactionTypeExpense, dateFrom, dateTo).
		Group("c.id, c.name, c.color, c.icon").
		Order("amount DESC").
		Scan(&results)

	// Calculate percentages
	var categories []CategoryAnalysis
	for _, result := range results {
		percentage := float64(0)
		if totalExpense > 0 {
			percentage = (result.Amount / totalExpense) * 100
		}

		categories = append(categories, CategoryAnalysis{
			CategoryID:   result.CategoryID,
			CategoryName: result.CategoryName,
			Amount:       result.Amount,
			Percentage:   percentage,
			Color:        result.Color,
			Icon:         result.Icon,
		})
	}

	analysis := ExpenseAnalysis{
		TotalExpense: totalExpense,
		Categories:   categories,
		PeriodStart:  dateFrom,
		PeriodEnd:    dateTo,
	}

	c.JSON(http.StatusOK, gin.H{"expense_analysis": analysis})
}

// GetIncomeAnalysis returns income breakdown by category
func GetIncomeAnalysis(c *gin.Context) {
	userID, err := middleware.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse date range
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	now := time.Now()
	if dateFrom == "" {
		dateFrom = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")
	}
	if dateTo == "" {
		dateTo = now.Format("2006-01-02")
	}

	// Get total income
	var totalIncome float64
	config.DB.Model(&models.Transaction{}).
		Where("user_id = ? AND type = ? AND date >= ? AND date <= ?",
			userID, models.TransactionTypeIncome, dateFrom, dateTo).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalIncome)

	// Get income by category
	var results []struct {
		CategoryID   string  `json:"category_id"`
		CategoryName string  `json:"category_name"`
		Amount       float64 `json:"amount"`
		Color        string  `json:"color"`
		Icon         string  `json:"icon"`
	}

	config.DB.Table("transactions t").
		Select("COALESCE(c.id::text, '') as category_id, COALESCE(c.name, 'Uncategorized') as category_name, SUM(t.amount) as amount, COALESCE(c.color, '#10b981') as color, COALESCE(c.icon, '💰') as icon").
		Joins("LEFT JOIN categories c ON t.category_id = c.id").
		Where("t.user_id = ? AND t.type = ? AND t.date >= ? AND t.date <= ?",
			userID, models.TransactionTypeIncome, dateFrom, dateTo).
		Group("c.id, c.name, c.color, c.icon").
		Order("amount DESC").
		Scan(&results)

	// Calculate percentages
	var categories []CategoryAnalysis
	for _, result := range results {
		percentage := float64(0)
		if totalIncome > 0 {
			percentage = (result.Amount / totalIncome) * 100
		}

		categories = append(categories, CategoryAnalysis{
			CategoryID:   result.CategoryID,
			CategoryName: result.CategoryName,
			Amount:       result.Amount,
			Percentage:   percentage,
			Color:        result.Color,
			Icon:         result.Icon,
		})
	}

	analysis := IncomeAnalysis{
		TotalIncome: totalIncome,
		Categories:  categories,
		PeriodStart: dateFrom,
		PeriodEnd:   dateTo,
	}

	c.JSON(http.StatusOK, gin.H{"income_analysis": analysis})
}

// GetSpendingTrends returns spending trends over time
func GetSpendingTrends(c *gin.Context) {
	userID, err := middleware.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse parameters
	period := c.DefaultQuery("period", "daily") // daily, weekly, monthly
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	now := time.Now()
	if dateFrom == "" {
		// Default to last 30 days for daily, last 12 weeks for weekly, last 12 months for monthly
		switch period {
		case "weekly":
			dateFrom = now.AddDate(0, 0, -84).Format("2006-01-02") // 12 weeks
		case "monthly":
			dateFrom = now.AddDate(-1, 0, 0).Format("2006-01-02") // 12 months
		default:
			dateFrom = now.AddDate(0, 0, -30).Format("2006-01-02") // 30 days
		}
	}
	if dateTo == "" {
		dateTo = now.Format("2006-01-02")
	}

	// Build date grouping based on period
	var dateFormat string
	switch period {
	case "weekly":
		dateFormat = "DATE_TRUNC('week', date)"
	case "monthly":
		dateFormat = "DATE_TRUNC('month', date)"
	default:
		dateFormat = "DATE_TRUNC('day', date)"
	}

	// Get trends data
	var results []struct {
		Date    time.Time `json:"date"`
		Income  float64   `json:"income"`
		Expense float64   `json:"expense"`
	}

	config.DB.Table("transactions").
		Select(dateFormat+" as date, "+
			"COALESCE(SUM(CASE WHEN type = 'income' THEN amount ELSE 0 END), 0) as income, "+
			"COALESCE(SUM(CASE WHEN type = 'expense' THEN amount ELSE 0 END), 0) as expense").
		Where("user_id = ? AND date >= ? AND date <= ?", userID, dateFrom, dateTo).
		Group("DATE_TRUNC('" + period + "', date)").
		Order("date ASC").
		Scan(&results)

	// Convert to response format
	var trends []SpendingTrend
	for _, result := range results {
		var dateStr string
		switch period {
		case "weekly":
			dateStr = result.Date.Format("2006-01-02") // Start of week
		case "monthly":
			dateStr = result.Date.Format("2006-01")
		default:
			dateStr = result.Date.Format("2006-01-02")
		}

		trends = append(trends, SpendingTrend{
			Date:    dateStr,
			Income:  result.Income,
			Expense: result.Expense,
			Net:     result.Income - result.Expense,
		})
	}

	response := SpendingTrends{
		Trends:      trends,
		Period:      period,
		PeriodStart: dateFrom,
		PeriodEnd:   dateTo,
	}

	c.JSON(http.StatusOK, gin.H{"spending_trends": response})
}

// GetCategoryAnalysis returns detailed category analysis
func GetCategoryAnalysis(c *gin.Context) {
	userID, err := middleware.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse date range
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	now := time.Now()
	if dateFrom == "" {
		dateFrom = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")
	}
	if dateTo == "" {
		dateTo = now.Format("2006-01-02")
	}

	// Get category analysis with transaction counts
	var results []struct {
		CategoryID       string  `json:"category_id"`
		CategoryName     string  `json:"category_name"`
		CategoryType     string  `json:"category_type"`
		TotalAmount      float64 `json:"total_amount"`
		TransactionCount int64   `json:"transaction_count"`
		AverageAmount    float64 `json:"average_amount"`
		Color            string  `json:"color"`
		Icon             string  `json:"icon"`
	}

	config.DB.Table("transactions t").
		Select("COALESCE(c.id::text, '') as category_id, "+
			"COALESCE(c.name, 'Uncategorized') as category_name, "+
			"COALESCE(c.type, 'expense') as category_type, "+
			"SUM(t.amount) as total_amount, "+
			"COUNT(t.id) as transaction_count, "+
			"AVG(t.amount) as average_amount, "+
			"COALESCE(c.color, '#6b7280') as color, "+
			"COALESCE(c.icon, '💸') as icon").
		Joins("LEFT JOIN categories c ON t.category_id = c.id").
		Where("t.user_id = ? AND t.date >= ? AND t.date <= ?", userID, dateFrom, dateTo).
		Group("c.id, c.name, c.type, c.color, c.icon").
		Order("total_amount DESC").
		Scan(&results)

	c.JSON(http.StatusOK, gin.H{
		"category_analysis": results,
		"period_start":      dateFrom,
		"period_end":        dateTo,
	})
}

// GetMonthlyReport returns comprehensive monthly report
func GetMonthlyReport(c *gin.Context) {
	userID, err := middleware.GetCurrentUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Parse month and year
	month := c.DefaultQuery("month", time.Now().Format("01"))
	year := c.DefaultQuery("year", time.Now().Format("2006"))

	// Calculate date range for the month
	monthInt := 1
	yearInt := 2024
	if m, err := time.Parse("01", month); err == nil {
		monthInt = int(m.Month())
	}
	if y, err := time.Parse("2006", year); err == nil {
		yearInt = y.Year()
	}

	startDate := time.Date(yearInt, time.Month(monthInt), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Nanosecond)

	dateFrom := startDate.Format("2006-01-02")
	dateTo := endDate.Format("2006-01-02")

	// Get comprehensive monthly data
	var totalIncome, totalExpense float64
	var transactionCount int64

	// Totals
	config.DB.Model(&models.Transaction{}).
		Where("user_id = ? AND type = ? AND date >= ? AND date <= ?",
			userID, models.TransactionTypeIncome, dateFrom, dateTo).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalIncome)

	config.DB.Model(&models.Transaction{}).
		Where("user_id = ? AND type = ? AND date >= ? AND date <= ?",
			userID, models.TransactionTypeExpense, dateFrom, dateTo).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalExpense)

	config.DB.Model(&models.Transaction{}).
		Where("user_id = ? AND date >= ? AND date <= ?", userID, dateFrom, dateTo).
		Count(&transactionCount)

	// Top categories
	var topExpenseCategories []CategoryAnalysis
	var topIncomeCategories []CategoryAnalysis

	// Top expense categories
	var expenseResults []struct {
		CategoryName string  `json:"category_name"`
		Amount       float64 `json:"amount"`
		Color        string  `json:"color"`
		Icon         string  `json:"icon"`
	}

	config.DB.Table("transactions t").
		Select("COALESCE(c.name, 'Uncategorized') as category_name, SUM(t.amount) as amount, COALESCE(c.color, '#6b7280') as color, COALESCE(c.icon, '💸') as icon").
		Joins("LEFT JOIN categories c ON t.category_id = c.id").
		Where("t.user_id = ? AND t.type = ? AND t.date >= ? AND t.date <= ?",
			userID, models.TransactionTypeExpense, dateFrom, dateTo).
		Group("c.name, c.color, c.icon").
		Order("amount DESC").
		Limit(5).
		Scan(&expenseResults)

	for _, result := range expenseResults {
		percentage := float64(0)
		if totalExpense > 0 {
			percentage = (result.Amount / totalExpense) * 100
		}
		topExpenseCategories = append(topExpenseCategories, CategoryAnalysis{
			CategoryName: result.CategoryName,
			Amount:       result.Amount,
			Percentage:   percentage,
			Color:        result.Color,
			Icon:         result.Icon,
		})
	}

	// Top income categories
	var incomeResults []struct {
		CategoryName string  `json:"category_name"`
		Amount       float64 `json:"amount"`
		Color        string  `json:"color"`
		Icon         string  `json:"icon"`
	}

	config.DB.Table("transactions t").
		Select("COALESCE(c.name, 'Uncategorized') as category_name, SUM(t.amount) as amount, COALESCE(c.color, '#10b981') as color, COALESCE(c.icon, '💰') as icon").
		Joins("LEFT JOIN categories c ON t.category_id = c.id").
		Where("t.user_id = ? AND t.type = ? AND t.date >= ? AND t.date <= ?",
			userID, models.TransactionTypeIncome, dateFrom, dateTo).
		Group("c.name, c.color, c.icon").
		Order("amount DESC").
		Limit(5).
		Scan(&incomeResults)

	for _, result := range incomeResults {
		percentage := float64(0)
		if totalIncome > 0 {
			percentage = (result.Amount / totalIncome) * 100
		}
		topIncomeCategories = append(topIncomeCategories, CategoryAnalysis{
			CategoryName: result.CategoryName,
			Amount:       result.Amount,
			Percentage:   percentage,
			Color:        result.Color,
			Icon:         result.Icon,
		})
	}

	report := gin.H{
		"month":                  month,
		"year":                   year,
		"period_start":           dateFrom,
		"period_end":             dateTo,
		"total_income":           totalIncome,
		"total_expense":          totalExpense,
		"net_income":             totalIncome - totalExpense,
		"transaction_count":      transactionCount,
		"top_expense_categories": topExpenseCategories,
		"top_income_categories":  topIncomeCategories,
		"savings_rate": func() float64 {
			if totalIncome > 0 {
				return ((totalIncome - totalExpense) / totalIncome) * 100
			}
			return 0
		}(),
	}

	c.JSON(http.StatusOK, gin.H{"monthly_report": report})
}

// ExportReport exports financial data
func ExportReport(c *gin.Context) {
	// This would implement CSV/PDF export functionality
	// For now, return a placeholder response
	c.JSON(http.StatusOK, gin.H{
		"message": "Export functionality will be implemented",
		"formats": []string{"CSV", "PDF", "Excel"},
	})
}
