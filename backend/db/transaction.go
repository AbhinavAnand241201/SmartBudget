package db

import (
	"time"

	"github.com/google/uuid"
)

// Transaction represents a financial transaction
type Transaction struct {
	ID          uuid.UUID `json:"id" db:"id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	Amount      float64   `json:"amount" db:"amount"`
	Category    string    `json:"category" db:"category"`
	Description string    `json:"description" db:"description"`
	Date        time.Time `json:"date" db:"date"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// NewTransaction creates a new transaction with default values
func NewTransaction(userID uuid.UUID, amount float64, category, description string) *Transaction {
	now := time.Now()
	return &Transaction{
		ID:          uuid.New(),
		UserID:      userID,
		Amount:      amount,
		Category:    category,
		Description: description,
		Date:        now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
