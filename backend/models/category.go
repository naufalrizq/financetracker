package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryType string

const (
	CategoryTypeIncome  CategoryType = "income"
	CategoryTypeExpense CategoryType = "expense"
)

type Category struct {
	ID          uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;index"`
	Name        string         `json:"name" gorm:"not null" validate:"required,min=2,max=100"`
	Type        CategoryType   `json:"type" gorm:"not null" validate:"required,oneof=income expense"`
	Color       string         `json:"color" gorm:"default:'#6366f1'" validate:"hexcolor"`
	Icon        string         `json:"icon" gorm:"default:'💰'"`
	Description string         `json:"description" gorm:"type:text"`
	IsDefault   bool           `json:"is_default" gorm:"default:false"`
	IsActive    bool           `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User         User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Transactions []Transaction `json:"transactions,omitempty" gorm:"foreignKey:CategoryID"`
	Budgets      []Budget      `json:"budgets,omitempty" gorm:"foreignKey:CategoryID"`
}

// BeforeCreate hook
func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

// DefaultCategories represents the default categories to seed
var DefaultCategories = []Category{
	// Income Categories
	{Name: "Salary", Type: CategoryTypeIncome, Color: "#10b981", Icon: "💼", Description: "Monthly salary and wages", IsDefault: true},
	{Name: "Freelance", Type: CategoryTypeIncome, Color: "#059669", Icon: "💻", Description: "Freelance work income", IsDefault: true},
	{Name: "Investment", Type: CategoryTypeIncome, Color: "#0d9488", Icon: "📈", Description: "Investment returns and dividends", IsDefault: true},
	{Name: "Business", Type: CategoryTypeIncome, Color: "#0891b2", Icon: "🏢", Description: "Business income", IsDefault: true},
	{Name: "Other Income", Type: CategoryTypeIncome, Color: "#0284c7", Icon: "💰", Description: "Other sources of income", IsDefault: true},

	// Expense Categories
	{Name: "Food & Dining", Type: CategoryTypeExpense, Color: "#dc2626", Icon: "🍽️", Description: "Restaurants, groceries, and food", IsDefault: true},
	{Name: "Transportation", Type: CategoryTypeExpense, Color: "#ea580c", Icon: "🚗", Description: "Gas, public transport, car maintenance", IsDefault: true},
	{Name: "Shopping", Type: CategoryTypeExpense, Color: "#d97706", Icon: "🛍️", Description: "Clothing, electronics, and general shopping", IsDefault: true},
	{Name: "Entertainment", Type: CategoryTypeExpense, Color: "#ca8a04", Icon: "🎬", Description: "Movies, games, hobbies", IsDefault: true},
	{Name: "Bills & Utilities", Type: CategoryTypeExpense, Color: "#65a30d", Icon: "📄", Description: "Electricity, water, internet, phone", IsDefault: true},
	{Name: "Healthcare", Type: CategoryTypeExpense, Color: "#16a34a", Icon: "🏥", Description: "Medical expenses, insurance", IsDefault: true},
	{Name: "Education", Type: CategoryTypeExpense, Color: "#059669", Icon: "📚", Description: "Books, courses, tuition", IsDefault: true},
	{Name: "Travel", Type: CategoryTypeExpense, Color: "#0891b2", Icon: "✈️", Description: "Vacation, business trips", IsDefault: true},
	{Name: "Home & Garden", Type: CategoryTypeExpense, Color: "#0284c7", Icon: "🏠", Description: "Rent, mortgage, home improvement", IsDefault: true},
	{Name: "Personal Care", Type: CategoryTypeExpense, Color: "#7c3aed", Icon: "💄", Description: "Haircut, cosmetics, personal items", IsDefault: true},
	{Name: "Gifts & Donations", Type: CategoryTypeExpense, Color: "#c026d3", Icon: "🎁", Description: "Gifts, charity, donations", IsDefault: true},
	{Name: "Other Expenses", Type: CategoryTypeExpense, Color: "#dc2626", Icon: "💸", Description: "Miscellaneous expenses", IsDefault: true},
}

// SeedDefaultCategories creates default categories for new users
func SeedDefaultCategories(db *gorm.DB) error {
	// This function will be called when a new user registers
	// to create default categories for them
	return nil
}

// CreateDefaultCategoriesForUser creates default categories for a specific user
func CreateDefaultCategoriesForUser(db *gorm.DB, userID uuid.UUID) error {
	var categories []Category

	for _, defaultCat := range DefaultCategories {
		category := Category{
			ID:          uuid.New(),
			UserID:      userID,
			Name:        defaultCat.Name,
			Type:        defaultCat.Type,
			Color:       defaultCat.Color,
			Icon:        defaultCat.Icon,
			Description: defaultCat.Description,
			IsDefault:   true,
			IsActive:    true,
		}
		categories = append(categories, category)
	}

	return db.Create(&categories).Error
}
