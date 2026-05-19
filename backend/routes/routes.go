package routes

import (
	"financetracker/controllers"

	"github.com/gin-gonic/gin"
)

// UserRoutes defines user-related routes
func UserRoutes(router *gin.RouterGroup) {
	users := router.Group("/users")
	{
		users.GET("/profile", controllers.GetProfile)
		users.PUT("/profile", controllers.UpdateProfile)
		users.DELETE("/account", controllers.DeleteAccount)
	}
}

// CategoryRoutes defines category-related routes
func CategoryRoutes(router *gin.RouterGroup) {
	categories := router.Group("/categories")
	{
		categories.GET("", controllers.GetCategories)
		categories.POST("", controllers.CreateCategory)
		categories.GET("/:id", controllers.GetCategory)
		categories.PUT("/:id", controllers.UpdateCategory)
		categories.DELETE("/:id", controllers.DeleteCategory)
	}
}

// AccountRoutes defines account-related routes
func AccountRoutes(router *gin.RouterGroup) {
	accounts := router.Group("/accounts")
	{
		accounts.GET("", controllers.GetAccounts)
		accounts.POST("", controllers.CreateAccount)
		accounts.GET("/:id", controllers.GetAccount)
		accounts.PUT("/:id", controllers.UpdateAccount)
		accounts.DELETE("/:id", controllers.DeleteAccount)
		accounts.GET("/:id/summary", controllers.GetAccountSummary)
	}
}

// TransactionRoutes defines transaction-related routes
func TransactionRoutes(router *gin.RouterGroup) {
	transactions := router.Group("/transactions")
	{
		transactions.GET("", controllers.GetTransactions)
		transactions.POST("", controllers.CreateTransaction)
		transactions.GET("/:id", controllers.GetTransaction)
		transactions.PUT("/:id", controllers.UpdateTransaction)
		transactions.DELETE("/:id", controllers.DeleteTransaction)
		transactions.POST("/bulk", controllers.CreateBulkTransactions)
		transactions.GET("/export", controllers.ExportTransactions)
	}
}

// BudgetRoutes defines budget-related routes
func BudgetRoutes(router *gin.RouterGroup) {
	budgets := router.Group("/budgets")
	{
		budgets.GET("", controllers.GetBudgets)
		budgets.POST("", controllers.CreateBudget)
		budgets.GET("/:id", controllers.GetBudget)
		budgets.PUT("/:id", controllers.UpdateBudget)
		budgets.DELETE("/:id", controllers.DeleteBudget)
		budgets.GET("/:id/status", controllers.GetBudgetStatus)
	}
}

// GoalRoutes defines goal-related routes
func GoalRoutes(router *gin.RouterGroup) {
	goals := router.Group("/goals")
	{
		goals.GET("", controllers.GetGoals)
		goals.POST("", controllers.CreateGoal)
		goals.GET("/:id", controllers.GetGoal)
		goals.PUT("/:id", controllers.UpdateGoal)
		goals.DELETE("/:id", controllers.DeleteGoal)
		goals.POST("/:id/progress", controllers.UpdateGoalProgress)
		goals.GET("/:id/status", controllers.GetGoalStatus)
	}
}

// ReportRoutes defines report-related routes
func ReportRoutes(router *gin.RouterGroup) {
	reports := router.Group("/reports")
	{
		reports.GET("/summary", controllers.GetFinancialSummary)
		reports.GET("/expenses", controllers.GetExpenseAnalysis)
		reports.GET("/income", controllers.GetIncomeAnalysis)
		reports.GET("/trends", controllers.GetSpendingTrends)
		reports.GET("/categories", controllers.GetCategoryAnalysis)
		reports.GET("/monthly", controllers.GetMonthlyReport)
		reports.GET("/export", controllers.ExportReport)
	}
}
