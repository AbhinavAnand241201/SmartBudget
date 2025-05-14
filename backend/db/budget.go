package db

import (
	"time"

	"github.com/google/uuid"
)

// BudgetPeriod represents the time period for a budget
type BudgetPeriod string

const (
	Daily   BudgetPeriod = "daily"
	Weekly  BudgetPeriod = "weekly"
	Monthly BudgetPeriod = "monthly"
)

// Budget represents a user's budget for a category
type Budget struct {
	ID        uuid.UUID    `json:"id" db:"id"`
	UserID    uuid.UUID    `json:"user_id" db:"user_id"`
	Category  string       `json:"category" db:"category"`
	Amount    float64      `json:"amount" db:"amount"`
	Period    BudgetPeriod `json:"period" db:"period"`
	StartDate time.Time    `json:"start_date" db:"start_date"`
	EndDate   *time.Time   `json:"end_date" db:"end_date"`
	CreatedAt time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt time.Time    `json:"updated_at" db:"updated_at"`
}

// NewBudget creates a new budget with default values
func NewBudget(userID uuid.UUID, category string, amount float64, period BudgetPeriod) *Budget {
	now := time.Now()
	return &Budget{
		ID:        uuid.New(),
		UserID:    userID,
		Category:  category,
		Amount:    amount,
		Period:    period,
		StartDate: now,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
