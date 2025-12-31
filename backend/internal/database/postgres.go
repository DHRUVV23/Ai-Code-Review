package database

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func InitDB() error {
	// 1. Connect
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

	Pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return fmt.Errorf("connection error: %w", err)
	}

	if err := Pool.Ping(context.Background()); err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}

	// 2. Create Tables (Safe "IF NOT EXISTS")

	// A. Users
	if _, err := Pool.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            github_id BIGINT UNIQUE,
            username TEXT,
            email TEXT UNIQUE NOT NULL,
            access_token TEXT,  -- 👈 ADD THIS LINE!
            created_at TIMESTAMP DEFAULT NOW(),
            updated_at TIMESTAMP DEFAULT NOW()
        );`); err != nil {
        return err
    }
	// B. Repositories
	if _, err := Pool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS repositories (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL,
			name TEXT NOT NULL,
			owner TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW(),
			CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id)
		);`); err != nil {
		return err
	}

	// C. Reviews (Base Table)
	if _, err := Pool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS reviews (
			id SERIAL PRIMARY KEY,
			repository_id INT NOT NULL,
			pr_number INT NOT NULL,
			status TEXT DEFAULT 'pending',
			created_at TIMESTAMP DEFAULT NOW(),
			CONSTRAINT fk_repo FOREIGN KEY(repository_id) REFERENCES repositories(id)
		);`); err != nil {
		return err
	}

	// D. Configurations (Added this back! ✅)
	if _, err := Pool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS configurations (
			id SERIAL PRIMARY KEY,
			repository_id INT NOT NULL UNIQUE,
			review_style TEXT DEFAULT 'concise',
			ignore_patterns TEXT,
			created_at TIMESTAMP DEFAULT NOW(),
			updated_at TIMESTAMP DEFAULT NOW(),
			CONSTRAINT fk_repo_config FOREIGN KEY(repository_id) REFERENCES repositories(id)
		);`); err != nil {
		return err
	}

	// 3. SMART MIGRATION: Add columns individually if they are missing
	migrations := []string{
		"ALTER TABLE reviews ADD COLUMN IF NOT EXISTS content TEXT;",
		"ALTER TABLE reviews ADD COLUMN IF NOT EXISTS commit_sha TEXT;",
		"ALTER TABLE reviews ADD COLUMN IF NOT EXISTS status TEXT;", 
	}

	for _, query := range migrations {
		if _, err := Pool.Exec(context.Background(), query); err != nil {
			fmt.Printf("⚠️ Migration Warning (Safe to ignore): %v\n", err)
		}
	}

	fmt.Println("✅ DB Connected & Schema Updated!")
	return nil
}

func CloseDB() {
	if Pool != nil {
		Pool.Close()
	}
}