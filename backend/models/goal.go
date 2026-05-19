package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GoalType string

const (
	GoalTypeSaving     GoalType = "saving"
	GoalTypeDebtPayoff GoalType = "debt_payoff"
	GoalTypeInvestment GoalType = "investment"
	GoalTypeOther      GoalType = "other"
)

type Goal struct {
	ID            uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID        uuid.UUID  `json:"user_id" gorm:"type:uuid;not null;index"`
	Name          string     `json:"name" gorm:"not null" validate:"required,min=2,max=100"`
	Description   string     `json:"description" gorm:"type:text"`
	Type          GoalType   `json:"type" gorm:"not null" validate:"required,oneof=saving debt_payoff investment other"`
	TargetAmount  float64    `json:"target_amount" gorm:"not null" validate:"required,gt=0"`
	CurrentAmount float64    `json:"current_amount" gorm:"default:0" validate:"gte=0"`
	Currency      string     `json:"currency" gorm:"default:'USD'" validate:"len=3"`
	TargetDate    *time.Time `json:"target_date"`
	IsCompleted   bool       `json:"is_completed" gorm:"default:false"`
	CompletedAt   *time.Time `json:"completed_at"`
	Color         string     `json:"color" gorm:"default:'#6366f1'" validate:"hexcolor"`
	Icon          string     `json:"icon" gorm:"default:'🎯'"`
	IsActive      bool       `json:"is_active" gorm:"default:true"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	User User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// BeforeCreate hook
func (g *Goal) BeforeCreate(tx *gorm.DB) error {
	if g.ID == uuid.Nil {
		g.ID = uuid.New()
	}
	return nil
}

// BeforeUpdate hook to check if goal is completed
func (g *Goal) BeforeUpdate(tx *gorm.DB) error {
	if g.CurrentAmount >= g.TargetAmount && !g.IsCompleted {
		g.IsCompleted = true
		now := time.Now()
		g.CompletedAt = &now
	} else if g.CurrentAmount < g.TargetAmount && g.IsCompleted {
		g.IsCompleted = false
		g.CompletedAt = nil
	}
	return nil
}

// UpdateProgress updates the current amount towards the goal
func (g *Goal) UpdateProgress(amount float64) {
	g.CurrentAmount += amount
	if g.CurrentAmount < 0 {
		g.CurrentAmount = 0
	}
}

// GetProgress returns the progress percentage
func (g *Goal) GetProgress() float64 {
	if g.TargetAmount == 0 {
		return 0
	}
	progress := (g.CurrentAmount / g.TargetAmount) * 100
	if progress > 100 {
		return 100
	}
	return progress
}

// GetRemainingAmount returns the remaining amount to reach the goal
func (g *Goal) GetRemainingAmount() float64 {
	remaining := g.TargetAmount - g.CurrentAmount
	if remaining < 0 {
		return 0
	}
	return remaining
}

// GetDaysRemaining returns the number of days remaining to reach the target date
func (g *Goal) GetDaysRemaining() *int {
	if g.TargetDate == nil {
		return nil
	}

	days := int(g.TargetDate.Sub(time.Now()).Hours() / 24)
	return &days
}

// GetRequiredMonthlyAmount calculates the monthly amount needed to reach the goal
func (g *Goal) GetRequiredMonthlyAmount() *float64 {
	if g.TargetDate == nil || g.IsCompleted {
		return nil
	}

	remaining := g.GetRemainingAmount()
	if remaining <= 0 {
		return nil
	}

	monthsRemaining := g.TargetDate.Sub(time.Now()).Hours() / (24 * 30)
	if monthsRemaining <= 0 {
		return nil
	}

	monthlyAmount := remaining / monthsRemaining
	return &monthlyAmount
}

// GoalStatus represents the current status of a goal with calculated fields
type GoalStatus struct {
	Goal
	Progress              float64  `json:"progress"`
	RemainingAmount       float64  `json:"remaining_amount"`
	DaysRemaining         *int     `json:"days_remaining"`
	RequiredMonthlyAmount *float64 `json:"required_monthly_amount"`
	IsOnTrack             bool     `json:"is_on_track"`
	Status                string   `json:"status"` // "completed", "on_track", "behind", "overdue"
}

// CalculateStatus calculates the current status of the goal
func (g *Goal) CalculateStatus() GoalStatus {
	progress := g.GetProgress()
	remainingAmount := g.GetRemainingAmount()
	daysRemaining := g.GetDaysRemaining()
	requiredMonthlyAmount := g.GetRequiredMonthlyAmount()

	// Determine status
	var status string
	var isOnTrack bool

	if g.IsCompleted {
		status = "completed"
		isOnTrack = true
	} else if g.TargetDate == nil {
		status = "in_progress"
		isOnTrack = true
	} else if daysRemaining != nil && *daysRemaining < 0 {
		status = "overdue"
		isOnTrack = false
	} else {
		// Calculate if on track based on time elapsed vs progress made
		if g.TargetDate != nil {
			totalDays := g.TargetDate.Sub(g.CreatedAt).Hours() / 24
			elapsedDays := time.Now().Sub(g.CreatedAt).Hours() / 24
			expectedProgress := (elapsedDays / totalDays) * 100

			if progress >= expectedProgress*0.9 { // 90% of expected progress
				status = "on_track"
				isOnTrack = true
			} else {
				status = "behind"
				isOnTrack = false
			}
		} else {
			status = "in_progress"
			isOnTrack = true
		}
	}

	return GoalStatus{
		Goal:                  *g,
		Progress:              progress,
		RemainingAmount:       remainingAmount,
		DaysRemaining:         daysRemaining,
		RequiredMonthlyAmount: requiredMonthlyAmount,
		IsOnTrack:             isOnTrack,
		Status:                status,
	}
}
