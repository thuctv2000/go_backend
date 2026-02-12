package repository

import (
	"context"
	"fmt"

	"my_backend/internal/database"
	"my_backend/internal/domain"
)

type postgresLixiGreetingRepository struct{}

func NewPostgresLixiGreetingRepository() domain.LixiGreetingRepository {
	return &postgresLixiGreetingRepository{}
}

func (r *postgresLixiGreetingRepository) Create(ctx context.Context, greeting *domain.LixiGreeting) error {
	query := `
		INSERT INTO lixi_greetings (name, amount, message, image)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	var id int64
	err := database.DB.QueryRow(ctx, query, greeting.Name, greeting.Amount, greeting.Message, greeting.Image).Scan(&id, &greeting.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create lixi greeting: %w", err)
	}

	greeting.ID = fmt.Sprintf("%d", id)
	return nil
}

func (r *postgresLixiGreetingRepository) GetAll(ctx context.Context) ([]*domain.LixiGreeting, error) {
	query := `
		SELECT id, name, amount, message, image, created_at
		FROM lixi_greetings
		ORDER BY created_at DESC
	`

	rows, err := database.DB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all lixi greetings: %w", err)
	}
	defer rows.Close()

	var greetings []*domain.LixiGreeting
	for rows.Next() {
		var greeting domain.LixiGreeting
		var id int64

		if err := rows.Scan(&id, &greeting.Name, &greeting.Amount, &greeting.Message, &greeting.Image, &greeting.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan lixi greeting: %w", err)
		}

		greeting.ID = fmt.Sprintf("%d", id)
		greetings = append(greetings, &greeting)
	}

	return greetings, nil
}
