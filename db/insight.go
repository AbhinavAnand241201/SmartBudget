package db

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// InsightType represents the type of insight
type InsightType string

const (
	Alert      InsightType = "alert"
	Tip        InsightType = "tip"
	Prediction InsightType = "prediction"
)

// Insight represents a financial insight for a user
type Insight struct {
	ID        uuid.UUID   `json:"id" db:"id"`
	UserID    uuid.UUID   `json:"user_id" db:"user_id"`
	Type      InsightType `json:"type" db:"type"`
	Content   string      `json:"content" db:"content"`
	CreatedAt time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt time.Time   `json:"updated_at" db:"updated_at"`
}

// NewInsight creates a new insight with default values
func NewInsight(userID uuid.UUID, insightType InsightType, content string) *Insight {
	now := time.Now()
	return &Insight{
		ID:        uuid.New(),
		UserID:    userID,
		Type:      insightType,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// CreateInsight creates a new insight in the database
func CreateInsight(ctx context.Context, insight *Insight) error {
	query := `
		INSERT INTO insights (id, user_id, type, content, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, type, content, created_at, updated_at
	`
	return GetPool().QueryRow(ctx, query,
		insight.ID, insight.UserID, insight.Type, insight.Content,
		insight.CreatedAt, insight.UpdatedAt,
	).Scan(
		&insight.ID, &insight.UserID, &insight.Type, &insight.Content,
		&insight.CreatedAt, &insight.UpdatedAt,
	)
}

// GetInsight retrieves an insight by ID
func GetInsight(ctx context.Context, id uuid.UUID) (*Insight, error) {
	query := `
		SELECT id, user_id, type, content, created_at, updated_at
		FROM insights
		WHERE id = $1
	`
	insight := &Insight{}
	err := GetPool().QueryRow(ctx, query, id).Scan(
		&insight.ID, &insight.UserID, &insight.Type, &insight.Content,
		&insight.CreatedAt, &insight.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	return insight, err
}

// GetUserInsights retrieves all insights for a user
func GetUserInsights(ctx context.Context, userID uuid.UUID) ([]*Insight, error) {
	query := `
		SELECT id, user_id, type, content, created_at, updated_at
		FROM insights
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := GetPool().Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var insights []*Insight
	for rows.Next() {
		insight := &Insight{}
		err := rows.Scan(
			&insight.ID, &insight.UserID, &insight.Type, &insight.Content,
			&insight.CreatedAt, &insight.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		insights = append(insights, insight)
	}
	return insights, rows.Err()
}

// UpdateInsight updates an existing insight
func UpdateInsight(ctx context.Context, insight *Insight) error {
	query := `
		UPDATE insights
		SET type = $1, content = $2, updated_at = $3
		WHERE id = $4
		RETURNING id, user_id, type, content, created_at, updated_at
	`
	return GetPool().QueryRow(ctx, query,
		insight.Type, insight.Content, time.Now(), insight.ID,
	).Scan(
		&insight.ID, &insight.UserID, &insight.Type, &insight.Content,
		&insight.CreatedAt, &insight.UpdatedAt,
	)
}

// DeleteInsight deletes an insight by ID
func DeleteInsight(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM insights WHERE id = $1`
	_, err := GetPool().Exec(ctx, query, id)
	return err
}
