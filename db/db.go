package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Insight represents a financial insight
type Insight struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Type        string    `json:"type"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Data        string    `json:"data"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// DB defines the interface for database operations
type DBInterface interface {
	// User operations
	CreateUser(ctx context.Context, user *User) error
	GetUser(ctx context.Context, id uuid.UUID) (*User, error)

	// Insight operations
	CreateInsight(ctx context.Context, insight *Insight) error
	GetInsight(ctx context.Context, id uuid.UUID) (*Insight, error)
	GetUserInsights(ctx context.Context, userID uuid.UUID) ([]*Insight, error)
	UpdateInsight(ctx context.Context, insight *Insight) error
	DeleteInsight(ctx context.Context, id uuid.UUID) error
}

// Ensure DB implements DBInterface
var _ DBInterface = (*DB)(nil)
