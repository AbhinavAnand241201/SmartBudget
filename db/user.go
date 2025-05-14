package db

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// NewUser creates a new user with default values
func NewUser(email string) *User {
	now := time.Now()
	return &User{
		ID:        uuid.New(),
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
