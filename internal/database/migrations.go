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
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
	`

	_, err := DB.Exec(ctx, createUsersTable)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	fmt.Println("✅ Database migrations completed successfully!")
	return nil
}
