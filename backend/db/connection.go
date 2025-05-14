package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool *pgxpool.Pool
)

// InitDB initializes the database connection pool
func InitDB() error {
	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		return fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return fmt.Errorf("unable to parse database URL: %v", err)
	}

	pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %v", err)
	}

	// Test the connection
	if err := pool.Ping(context.Background()); err != nil {
		return fmt.Errorf("unable to ping database: %v", err)
	}

	log.Println("Successfully connected to database")
	return nil
}

// GetPool returns the database connection pool
func GetPool() *pgxpool.Pool {
	return pool
}

// CloseDB closes the database connection pool
func CloseDB() {
	if pool != nil {
		pool.Close()
	}
}
