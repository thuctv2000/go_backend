package repository

import (
	"context"
	"errors"
	"my_backend/internal/database"
	"my_backend/internal/domain"
)

type postgresUserRepository struct{}

// NewPostgresUserRepository creates a new PostgreSQL user repository
func NewPostgresUserRepository() domain.UserRepository {
	return &postgresUserRepository{}
}

func (r *postgresUserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2)
		RETURNING id
	`

	var id int
	err := database.DB.QueryRow(ctx, query, user.Email, user.Password).Scan(&id)
	if err != nil {
		// Check for unique violation
		if err.Error() == "ERROR: duplicate key value violates unique constraint \"users_email_key\" (SQLSTATE 23505)" {
			return errors.New("user already exists")
		}
		return err
	}

	user.ID = string(rune(id))
	return nil
}

func (r *postgresUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, password_hash
		FROM users
		WHERE email = $1
	`

	var user domain.User
	var id int
	err := database.DB.QueryRow(ctx, query, email).Scan(&id, &user.Email, &user.Password)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	user.ID = string(rune(id))
	return &user, nil
}
