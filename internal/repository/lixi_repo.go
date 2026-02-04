package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"my_backend/internal/database"
	"my_backend/internal/domain"
)

type postgresLixiRepository struct{}

func NewPostgresLixiRepository() domain.LixiRepository {
	return &postgresLixiRepository{}
}

func (r *postgresLixiRepository) Create(ctx context.Context, config *domain.LixiConfig) error {
	envelopesJSON, err := json.Marshal(config.Envelopes)
	if err != nil {
		return fmt.Errorf("failed to marshal envelopes: %w", err)
	}

	query := `
		INSERT INTO lixi_configs (name, envelopes, is_active)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	var id int64
	err = database.DB.QueryRow(ctx, query, config.Name, envelopesJSON, config.IsActive).Scan(&id, &config.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create lixi config: %w", err)
	}

	config.ID = fmt.Sprintf("%d", id)
	return nil
}

func (r *postgresLixiRepository) GetActive(ctx context.Context) (*domain.LixiConfig, error) {
	query := `
		SELECT id, name, envelopes, is_active, created_at
		FROM lixi_configs
		WHERE is_active = TRUE
		LIMIT 1
	`

	var config domain.LixiConfig
	var id int64
	var envelopesJSON []byte

	err := database.DB.QueryRow(ctx, query).Scan(&id, &config.Name, &envelopesJSON, &config.IsActive, &config.CreatedAt)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, errors.New("no active lixi config found")
		}
		return nil, fmt.Errorf("failed to get active lixi config: %w", err)
	}

	if err := json.Unmarshal(envelopesJSON, &config.Envelopes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal envelopes: %w", err)
	}

	config.ID = fmt.Sprintf("%d", id)
	return &config, nil
}

func (r *postgresLixiRepository) GetByID(ctx context.Context, id string) (*domain.LixiConfig, error) {
	query := `
		SELECT id, name, envelopes, is_active, created_at
		FROM lixi_configs
		WHERE id = $1
	`

	var config domain.LixiConfig
	var dbID int64
	var envelopesJSON []byte

	err := database.DB.QueryRow(ctx, query, id).Scan(&dbID, &config.Name, &envelopesJSON, &config.IsActive, &config.CreatedAt)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, errors.New("lixi config not found")
		}
		return nil, fmt.Errorf("failed to get lixi config: %w", err)
	}

	if err := json.Unmarshal(envelopesJSON, &config.Envelopes); err != nil {
		return nil, fmt.Errorf("failed to unmarshal envelopes: %w", err)
	}

	config.ID = fmt.Sprintf("%d", dbID)
	return &config, nil
}

func (r *postgresLixiRepository) GetAll(ctx context.Context) ([]*domain.LixiConfig, error) {
	query := `
		SELECT id, name, envelopes, is_active, created_at
		FROM lixi_configs
		ORDER BY created_at DESC
	`

	rows, err := database.DB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all lixi configs: %w", err)
	}
	defer rows.Close()

	var configs []*domain.LixiConfig
	for rows.Next() {
		var config domain.LixiConfig
		var id int64
		var envelopesJSON []byte

		if err := rows.Scan(&id, &config.Name, &envelopesJSON, &config.IsActive, &config.CreatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan lixi config: %w", err)
		}

		if err := json.Unmarshal(envelopesJSON, &config.Envelopes); err != nil {
			return nil, fmt.Errorf("failed to unmarshal envelopes: %w", err)
		}

		config.ID = fmt.Sprintf("%d", id)
		configs = append(configs, &config)
	}

	return configs, nil
}

func (r *postgresLixiRepository) Update(ctx context.Context, config *domain.LixiConfig) error {
	envelopesJSON, err := json.Marshal(config.Envelopes)
	if err != nil {
		return fmt.Errorf("failed to marshal envelopes: %w", err)
	}

	query := `
		UPDATE lixi_configs
		SET name = $1, envelopes = $2
		WHERE id = $3
	`

	result, err := database.DB.Exec(ctx, query, config.Name, envelopesJSON, config.ID)
	if err != nil {
		return fmt.Errorf("failed to update lixi config: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("lixi config not found")
	}

	return nil
}

func (r *postgresLixiRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM lixi_configs WHERE id = $1`

	result, err := database.DB.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete lixi config: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("lixi config not found")
	}

	return nil
}

func (r *postgresLixiRepository) SetActive(ctx context.Context, id string) error {
	// Use transaction to ensure atomicity
	tx, err := database.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Deactivate all configs
	_, err = tx.Exec(ctx, `UPDATE lixi_configs SET is_active = FALSE WHERE is_active = TRUE`)
	if err != nil {
		return fmt.Errorf("failed to deactivate configs: %w", err)
	}

	// Activate the specified config
	result, err := tx.Exec(ctx, `UPDATE lixi_configs SET is_active = TRUE WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to activate config: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.New("lixi config not found")
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
