package database

import (
	"context"
	"fmt"
)

// RunMigrations creates the necessary database tables
func RunMigrations() error {
	ctx := context.Background()

	// Create users table
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := DB.Exec(ctx, createUsersTable)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	// Create lixi_configs table
	createLixiConfigsTable := `
	CREATE TABLE IF NOT EXISTS lixi_configs (
		id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
		name TEXT NOT NULL,
		envelopes JSONB NOT NULL,
		is_active BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err = DB.Exec(ctx, createLixiConfigsTable)
	if err != nil {
		return fmt.Errorf("failed to create lixi_configs table: %w", err)
	}

	// Ensure only one active config - create unique partial index
	createLixiActiveIndex := `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_lixi_active
	ON lixi_configs (is_active) WHERE is_active = TRUE;
	`

	_, err = DB.Exec(ctx, createLixiActiveIndex)
	if err != nil {
		return fmt.Errorf("failed to create lixi active index: %w", err)
	}

	fmt.Println("✅ Database migrations completed successfully!")
	return nil
}
