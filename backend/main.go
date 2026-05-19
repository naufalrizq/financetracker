package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"financetracker/config"
	"financetracker/controllers"
	"financetracker/middleware"
	"financetracker/models"
	"financetracker/routes"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Initialize database
	config.ConnectDatabase()

	// Auto-migrate database schemas
	if err := config.DB.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Account{},
		&models.Transaction{},
		&models.Budget{},
		&models.Goal{},
	); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Seed default categories
	models.SeedDefaultCategories(config.DB)

	// Initialize Gin router
	router := gin.Default()

	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "FinanceTracker API is running",
		})
	})

	// API routes
	api := router.Group("/api")
	{
		// Public routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
			auth.POST("/refresh", controllers.RefreshToken)
		}

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			routes.UserRoutes(protected)
			routes.CategoryRoutes(protected)
			routes.AccountRoutes(protected)
			routes.TransactionRoutes(protected)
			routes.BudgetRoutes(protected)
			routes.GoalRoutes(protected)
			routes.ReportRoutes(protected)
		}
	}

	// API documentation
	router.GET("/docs", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "FinanceTracker API Documentation",
			"version": "1.0.0",
			"endpoints": gin.H{
				"auth": gin.H{
					"POST /api/auth/register": "Register new user",
					"POST /api/auth/login":    "Login user",
					"POST /api/auth/refresh":  "Refresh JWT token",
				},
				"users": gin.H{
					"GET /api/users/profile":    "Get user profile",
					"PUT /api/users/profile":    "Update user profile",
					"DELETE /api/users/account": "Delete user account",
				},
				"categories": gin.H{
					"GET /api/categories":        "Get all categories",
					"POST /api/categories":       "Create category",
					"PUT /api/categories/:id":    "Update category",
					"DELETE /api/categories/:id": "Delete category",
				},
				"accounts": gin.H{
					"GET /api/accounts":        "Get user accounts",
					"POST /api/accounts":       "Create account",
					"PUT /api/accounts/:id":    "Update account",
					"DELETE /api/accounts/:id": "Delete account",
				},
				"transactions": gin.H{
					"GET /api/transactions":        "Get transactions with filters",
					"POST /api/transactions":       "Create transaction",
					"GET /api/transactions/:id":    "Get transaction by ID",
					"PUT /api/transactions/:id":    "Update transaction",
					"DELETE /api/transactions/:id": "Delete transaction",
				},
				"budgets": gin.H{
					"GET /api/budgets":        "Get user budgets",
					"POST /api/budgets":       "Create budget",
					"PUT /api/budgets/:id":    "Update budget",
					"DELETE /api/budgets/:id": "Delete budget",
				},
				"goals": gin.H{
					"GET /api/goals":        "Get financial goals",
					"POST /api/goals":       "Create goal",
					"PUT /api/goals/:id":    "Update goal",
					"DELETE /api/goals/:id": "Delete goal",
				},
				"reports": gin.H{
					"GET /api/reports/summary":  "Get financial summary",
					"GET /api/reports/expenses": "Get expense analysis",
					"GET /api/reports/income":   "Get income analysis",
					"GET /api/reports/trends":   "Get spending trends",
					"GET /api/reports/export":   "Export data (CSV/PDF)",
				},
			},
		})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
