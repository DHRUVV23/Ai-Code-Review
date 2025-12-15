package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Global variable for the connection pool
var Pool *pgxpool.Pool

func InitDB() error {
	// Build connection string from .env variables
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("config error: %w", err)
	}

	config.MaxConns = 10
	config.MinConns = 2

	// Connect to database
	Pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return fmt.Errorf("connection error: %w", err)
	}

	// Verify connection with a Ping
	if err := Pool.Ping(context.Background()); err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}

	fmt.Println("✅ DB Connected Successfully!")
	return nil
}

func CloseDB() {
	if Pool != nil {
		Pool.Close()
	}
}
