package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BudgetPeriod string

const (
	BudgetPeriodWeekly  BudgetPeriod = "weekly"
	BudgetPeriodMonthly BudgetPeriod = "monthly"
	BudgetPeriodYearly  BudgetPeriod = "yearly"
)

type Budget struct {
	ID         uuid.UUID    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID     uuid.UUID    `json:"user_id" gorm:"type:uuid;not null;index"`
	CategoryID *uuid.UUID   `json:"category_id" gorm:"type:uuid;index"`
	Name       string       `json:"name" gorm:"not null" validate:"required,min=2,max=100"`
	Amount     float64      `json:"amount" gorm:"not null" validate:"required,gt=0"`
	Currency   string       `json:"currency" gorm:"default:'USD'" validate:"len=3"`
	Period     BudgetPeriod `json:"period" gorm:"not null" validate:"required,oneof=weekly monthly yearly"`
	StartDate  time.Time    `json:"start_date" gorm:"not null" validate:"required"`
	EndDate    *time.Time   `json:"end_date"`
	IsActive   bool         `json:"is_active" gorm:"default:true"`

	// Alert settings
	AlertThreshold float64    `json:"alert_threshold" gorm:"default:80" validate:"min=0,max=100"` // Percentage
	AlertEnabled   bool       `json:"alert_enabled" gorm:"default:true"`
	LastAlertSent  *time.Time `json:"last_alert_sent"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User     User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Category *Category `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
}

// BeforeCreate hook
func (b *Budget) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

// GetCurrentPeriodStart returns the start date of the current budget period
func (b *Budget) GetCurrentPeriodStart() time.Time {
	now := time.Now()

	switch b.Period {
	case BudgetPeriodWeekly:
		// Start of current week (Monday)
		weekday := int(now.Weekday())
		if weekday == 0 { // Sunday
			weekday = 7
		}
		return now.AddDate(0, 0, -(weekday - 1)).Truncate(24 * time.Hour)
	case BudgetPeriodMonthly:
		// Start of current month
		return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	case BudgetPeriodYearly:
		// Start of current year
		return time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	default:
		return b.StartDate
	}
}

// GetCurrentPeriodEnd returns the end date of the current budget period
func (b *Budget) GetCurrentPeriodEnd() time.Time {
	start := b.GetCurrentPeriodStart()

	switch b.Period {
	case BudgetPeriodWeekly:
		return start.AddDate(0, 0, 7).Add(-time.Nanosecond)
	case BudgetPeriodMonthly:
		return start.AddDate(0, 1, 0).Add(-time.Nanosecond)
	case BudgetPeriodYearly:
		return start.AddDate(1, 0, 0).Add(-time.Nanosecond)
	default:
		if b.EndDate != nil {
			return *b.EndDate
		}
		return start.AddDate(0, 1, 0).Add(-time.Nanosecond)
	}
}

// BudgetStatus represents the current status of a budget
type BudgetStatus struct {
	Budget
	Spent          float64   `json:"spent"`
	Remaining      float64   `json:"remaining"`
	PercentageUsed float64   `json:"percentage_used"`
	IsOverBudget   bool      `json:"is_over_budget"`
	DaysRemaining  int       `json:"days_remaining"`
	CategoryName   string    `json:"category_name,omitempty"`
	PeriodStart    time.Time `json:"period_start"`
	PeriodEnd      time.Time `json:"period_end"`
}

// CalculateStatus calculates the current status of the budget
func (b *Budget) CalculateStatus(db *gorm.DB) (*BudgetStatus, error) {
	periodStart := b.GetCurrentPeriodStart()
	periodEnd := b.GetCurrentPeriodEnd()

	// Calculate spent amount
	var spent float64
	query := db.Model(&Transaction{}).
		Where("user_id = ? AND type = ? AND date >= ? AND date <= ?",
			b.UserID, TransactionTypeExpense, periodStart, periodEnd)

	if b.CategoryID != nil {
		query = query.Where("category_id = ?", *b.CategoryID)
	}

	if err := query.Select("COALESCE(SUM(amount), 0)").Scan(&spent).Error; err != nil {
		return nil, err
	}

	// Calculate remaining and percentage
	remaining := b.Amount - spent
	percentageUsed := (spent / b.Amount) * 100
	isOverBudget := spent > b.Amount

	// Calculate days remaining in period
	daysRemaining := int(periodEnd.Sub(time.Now()).Hours() / 24)
	if daysRemaining < 0 {
		daysRemaining = 0
	}

	// Get category name if applicable
	var categoryName string
	if b.CategoryID != nil && b.Category != nil {
		categoryName = b.Category.Name
	}

	return &BudgetStatus{
		Budget:         *b,
		Spent:          spent,
		Remaining:      remaining,
		PercentageUsed: percentageUsed,
		IsOverBudget:   isOverBudget,
		DaysRemaining:  daysRemaining,
		CategoryName:   categoryName,
		PeriodStart:    periodStart,
		PeriodEnd:      periodEnd,
	}, nil
}
