package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccountType string

const (
	AccountTypeChecking   AccountType = "checking"
	AccountTypeSavings    AccountType = "savings"
	AccountTypeCredit     AccountType = "credit"
	AccountTypeCash       AccountType = "cash"
	AccountTypeInvestment AccountType = "investment"
)

type Account struct {
	ID             uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID         uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;index"`
	Name           string         `json:"name" gorm:"not null" validate:"required,min=2,max=100"`
	Type           AccountType    `json:"type" gorm:"not null" validate:"required,oneof=checking savings credit cash investment"`
	Balance        float64        `json:"balance" gorm:"default:0" validate:"numeric"`
	Currency       string         `json:"currency" gorm:"default:'USD'" validate:"len=3"`
	Color          string         `json:"color" gorm:"default:'#6366f1'" validate:"hexcolor"`
	Icon           string         `json:"icon" gorm:"default:'🏦'"`
	Description    string         `json:"description" gorm:"type:text"`
	IsActive       bool           `json:"is_active" gorm:"default:true"`
	IncludeInTotal bool           `json:"include_in_total" gorm:"default:true"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User         User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Transactions []Transaction `json:"transactions,omitempty" gorm:"foreignKey:AccountID"`
}

// BeforeCreate hook
func (a *Account) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}

// UpdateBalance updates the account balance based on transactions
func (a *Account) UpdateBalance(db *gorm.DB) error {
	var totalIncome, totalExpense float64

	// Calculate total income
	db.Model(&Transaction{}).
		Where("account_id = ? AND type = ?", a.ID, TransactionTypeIncome).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalIncome)

	// Calculate total expense
	db.Model(&Transaction{}).
		Where("account_id = ? AND type = ?", a.ID, TransactionTypeExpense).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalExpense)

	// Update balance
	newBalance := totalIncome - totalExpense
	return db.Model(a).Update("balance", newBalance).Error
}

// AccountSummary represents account summary with additional calculated fields
type AccountSummary struct {
	Account
	TotalIncome       float64    `json:"total_income"`
	TotalExpense      float64    `json:"total_expense"`
	TransactionCount  int64      `json:"transaction_count"`
	LastTransactionAt *time.Time `json:"last_transaction_at"`
}
