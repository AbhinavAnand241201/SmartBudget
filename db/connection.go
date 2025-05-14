package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

// DB represents a database connection
type DB struct {
	*sql.DB
}

// NewDB creates a new database connection
func NewDB(databaseURL string) (*DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	return &DB{db}, nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
}

// CreateUser creates a new user
func (db *DB) CreateUser(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (id, email, name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	_, err := db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.Name,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	return nil
}

// GetUser retrieves a user by ID
func (db *DB) GetUser(ctx context.Context, id uuid.UUID) (*User, error) {
	query := `
		SELECT id, email, name, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	user := &User{}
	err := db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting user: %w", err)
	}

	return user, nil
}

// CreateInsight creates a new insight
func (db *DB) CreateInsight(ctx context.Context, insight *Insight) error {
	query := `
		INSERT INTO insights (id, user_id, type, title, description, data, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	now := time.Now()
	insight.CreatedAt = now
	insight.UpdatedAt = now

	_, err := db.ExecContext(ctx, query,
		insight.ID,
		insight.UserID,
		insight.Type,
		insight.Title,
		insight.Description,
		insight.Data,
		insight.CreatedAt,
		insight.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("error creating insight: %w", err)
	}

	return nil
}

// GetInsight retrieves an insight by ID
func (db *DB) GetInsight(ctx context.Context, id uuid.UUID) (*Insight, error) {
	query := `
		SELECT id, user_id, type, title, description, data, created_at, updated_at
		FROM insights
		WHERE id = $1
	`
	insight := &Insight{}
	err := db.QueryRowContext(ctx, query, id).Scan(
		&insight.ID,
		&insight.UserID,
		&insight.Type,
		&insight.Title,
		&insight.Description,
		&insight.Data,
		&insight.CreatedAt,
		&insight.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("error getting insight: %w", err)
	}

	return insight, nil
}

// GetUserInsights retrieves all insights for a user
func (db *DB) GetUserInsights(ctx context.Context, userID uuid.UUID) ([]*Insight, error) {
	query := `
		SELECT id, user_id, type, title, description, data, created_at, updated_at
		FROM insights
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user insights: %w", err)
	}
	defer rows.Close()

	var insights []*Insight
	for rows.Next() {
		insight := &Insight{}
		err := rows.Scan(
			&insight.ID,
			&insight.UserID,
			&insight.Type,
			&insight.Title,
			&insight.Description,
			&insight.Data,
			&insight.CreatedAt,
			&insight.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning insight: %w", err)
		}
		insights = append(insights, insight)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating insights: %w", err)
	}

	return insights, nil
}

// UpdateInsight updates an existing insight
func (db *DB) UpdateInsight(ctx context.Context, insight *Insight) error {
	query := `
		UPDATE insights
		SET type = $1, title = $2, description = $3, data = $4, updated_at = $5
		WHERE id = $6
	`
	insight.UpdatedAt = time.Now()

	_, err := db.ExecContext(ctx, query,
		insight.Type,
		insight.Title,
		insight.Description,
		insight.Data,
		insight.UpdatedAt,
		insight.ID,
	)
	if err != nil {
		return fmt.Errorf("error updating insight: %w", err)
	}

	return nil
}

// DeleteInsight deletes an insight
func (db *DB) DeleteInsight(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM insights
		WHERE id = $1
	`
	_, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting insight: %w", err)
	}

	return nil
}
