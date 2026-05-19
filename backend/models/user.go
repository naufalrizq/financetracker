package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID               uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Email            string         `json:"email" gorm:"uniqueIndex;not null" validate:"required,email"`
	Password         string         `json:"-" gorm:"not null" validate:"required,min=8"`
	FirstName        string         `json:"first_name" gorm:"not null" validate:"required,min=2,max=50"`
	LastName         string         `json:"last_name" gorm:"not null" validate:"required,min=2,max=50"`
	Currency         string         `json:"currency" gorm:"default:'USD'" validate:"len=3"`
	Timezone         string         `json:"timezone" gorm:"default:'UTC'"`
	IsActive         bool           `json:"is_active" gorm:"default:true"`
	EmailVerified    bool           `json:"email_verified" gorm:"default:false"`
	RefreshTokenHash string         `json:"-" gorm:"type:text"`
	LastLoginAt      *time.Time     `json:"last_login_at"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Accounts     []Account     `json:"accounts,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Categories   []Category    `json:"categories,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Transactions []Transaction `json:"transactions,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Budgets      []Budget      `json:"budgets,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Goals        []Goal        `json:"goals,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// BeforeCreate hook to hash password
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return u.HashPassword()
}

// BeforeUpdate hook to hash password if changed
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	if tx.Statement.Changed("Password") {
		return u.HashPassword()
	}
	return nil
}

// HashPassword hashes the user password
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword verifies the password
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// GetFullName returns the user's full name
func (u *User) GetFullName() string {
	return u.FirstName + " " + u.LastName
}

// UserResponse represents the user data returned to client (without sensitive fields)
type UserResponse struct {
	ID            uuid.UUID  `json:"id"`
	Email         string     `json:"email"`
	FirstName     string     `json:"first_name"`
	LastName      string     `json:"last_name"`
	Currency      string     `json:"currency"`
	Timezone      string     `json:"timezone"`
	IsActive      bool       `json:"is_active"`
	EmailVerified bool       `json:"email_verified"`
	LastLoginAt   *time.Time `json:"last_login_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:            u.ID,
		Email:         u.Email,
		FirstName:     u.FirstName,
		LastName:      u.LastName,
		Currency:      u.Currency,
		Timezone:      u.Timezone,
		IsActive:      u.IsActive,
		EmailVerified: u.EmailVerified,
		LastLoginAt:   u.LastLoginAt,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}
