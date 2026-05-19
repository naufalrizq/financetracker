package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionType string

const (
	TransactionTypeIncome   TransactionType = "income"
	TransactionTypeExpense  TransactionType = "expense"
	TransactionTypeTransfer TransactionType = "transfer"
)

type Transaction struct {
	ID          uuid.UUID       `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID      uuid.UUID       `json:"user_id" gorm:"type:uuid;not null;index"`
	AccountID   uuid.UUID       `json:"account_id" gorm:"type:uuid;not null;index"`
	CategoryID  *uuid.UUID      `json:"category_id" gorm:"type:uuid;index"`
	Type        TransactionType `json:"type" gorm:"not null" validate:"required,oneof=income expense transfer"`
	Amount      float64         `json:"amount" gorm:"not null" validate:"required,gt=0"`
	Currency    string          `json:"currency" gorm:"default:'USD'" validate:"len=3"`
	Description string          `json:"description" gorm:"not null" validate:"required,min=2,max=255"`
	Notes       string          `json:"notes" gorm:"type:text"`
	Date        time.Time       `json:"date" gorm:"not null;index" validate:"required"`

	// Transfer specific fields
	ToAccountID *uuid.UUID `json:"to_account_id" gorm:"type:uuid;index"`

	// Recurring transaction fields
	IsRecurring    bool       `json:"is_recurring" gorm:"default:false"`
	RecurringType  string     `json:"recurring_type" gorm:"default:''" validate:"omitempty,oneof=daily weekly monthly yearly"`
	RecurringUntil *time.Time `json:"recurring_until"`

	// Metadata
	Tags      string         `json:"tags" gorm:"type:text"` // JSON array of tags
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User      User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Account   Account   `json:"account,omitempty" gorm:"foreignKey:AccountID"`
	Category  *Category `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	ToAccount *Account  `json:"to_account,omitempty" gorm:"foreignKey:ToAccountID"`
}

// BeforeCreate hook
func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}

// AfterCreate hook to update account balance
func (t *Transaction) AfterCreate(tx *gorm.DB) error {
	return t.updateAccountBalances(tx)
}

// AfterUpdate hook to update account balance
func (t *Transaction) AfterUpdate(tx *gorm.DB) error {
	return t.updateAccountBalances(tx)
}

// AfterDelete hook to update account balance
func (t *Transaction) AfterDelete(tx *gorm.DB) error {
	return t.updateAccountBalances(tx)
}

// updateAccountBalances updates the balance of affected accounts
func (t *Transaction) updateAccountBalances(tx *gorm.DB) error {
	// Update primary account balance
	var account Account
	if err := tx.First(&account, t.AccountID).Error; err != nil {
		return err
	}
	if err := account.UpdateBalance(tx); err != nil {
		return err
	}

	// Update destination account balance for transfers
	if t.Type == TransactionTypeTransfer && t.ToAccountID != nil {
		var toAccount Account
		if err := tx.First(&toAccount, *t.ToAccountID).Error; err != nil {
			return err
		}
		if err := toAccount.UpdateBalance(tx); err != nil {
			return err
		}
	}

	return nil
}

// TransactionSummary represents transaction data with additional calculated fields
type TransactionSummary struct {
	Transaction
	CategoryName  string `json:"category_name"`
	AccountName   string `json:"account_name"`
	ToAccountName string `json:"to_account_name,omitempty"`
}

// TransactionFilter represents filters for transaction queries
type TransactionFilter struct {
	UserID     uuid.UUID        `json:"user_id"`
	AccountID  *uuid.UUID       `json:"account_id"`
	CategoryID *uuid.UUID       `json:"category_id"`
	Type       *TransactionType `json:"type"`
	DateFrom   *time.Time       `json:"date_from"`
	DateTo     *time.Time       `json:"date_to"`
	AmountMin  *float64         `json:"amount_min"`
	AmountMax  *float64         `json:"amount_max"`
	Search     string           `json:"search"`
	Tags       []string         `json:"tags"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	SortBy     string           `json:"sort_by"`
	SortOrder  string           `json:"sort_order"`
}

// GetDefaultFilter returns default filter values
func GetDefaultFilter() TransactionFilter {
	return TransactionFilter{
		Page:      1,
		Limit:     20,
		SortBy:    "date",
		SortOrder: "desc",
	}
}
