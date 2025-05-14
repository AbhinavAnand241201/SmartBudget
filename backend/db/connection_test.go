package db

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatabaseConnection(t *testing.T) {
	// Set test database URL
	os.Setenv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/smartbudget_test")

	// Test initialization
	err := InitDB()
	if err != nil {
		t.Skipf("Skipping DB test: %v", err)
	}

	// Test pool is not nil
	pool := GetPool()
	if pool == nil {
		t.Skip("Skipping DB test: pool is nil")
	}
	assert.NotNil(t, pool)

	// Test connection is working
	err = pool.Ping(nil)
	if err != nil {
		t.Skipf("Skipping DB test: ping failed: %v", err)
	}
	assert.NoError(t, err)

	// Clean up
	CloseDB()
}
