package db

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
)

// MockDB is a mock implementation of the database interface
type MockDB struct {
	users    map[uuid.UUID]*User
	insights map[uuid.UUID]*Insight
	mu       sync.RWMutex
}

// NewMockDB creates a new mock database
func NewMockDB() *MockDB {
	return &MockDB{
		users:    make(map[uuid.UUID]*User),
		insights: make(map[uuid.UUID]*Insight),
	}
}

// CreateUser creates a new user
func (m *MockDB) CreateUser(ctx context.Context, user *User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.users[user.ID]; exists {
		return errors.New("user already exists")
	}

	m.users[user.ID] = user
	return nil
}

// GetUser gets a user by ID
func (m *MockDB) GetUser(ctx context.Context, id uuid.UUID) (*User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	user, exists := m.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// CreateInsight creates a new insight
func (m *MockDB) CreateInsight(ctx context.Context, insight *Insight) (*Insight, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.insights[insight.ID]; exists {
		return nil, errors.New("insight already exists")
	}

	m.insights[insight.ID] = insight
	return insight, nil
}

// GetInsight gets an insight by ID
func (m *MockDB) GetInsight(ctx context.Context, id uuid.UUID) (*Insight, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	insight, exists := m.insights[id]
	if !exists {
		return nil, errors.New("insight not found")
	}

	return insight, nil
}

// GetUserInsights gets all insights for a user
func (m *MockDB) GetUserInsights(ctx context.Context, userID uuid.UUID) ([]*Insight, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var insights []*Insight
	for _, insight := range m.insights {
		if insight.UserID == userID {
			insights = append(insights, insight)
		}
	}

	return insights, nil
}

// UpdateInsight updates an insight
func (m *MockDB) UpdateInsight(ctx context.Context, insight *Insight) (*Insight, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.insights[insight.ID]; !exists {
		return nil, errors.New("insight not found")
	}

	m.insights[insight.ID] = insight
	return insight, nil
}

// DeleteInsight deletes an insight
func (m *MockDB) DeleteInsight(ctx context.Context, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.insights[id]; !exists {
		return errors.New("insight not found")
	}

	delete(m.insights, id)
	return nil
}
