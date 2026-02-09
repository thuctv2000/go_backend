package service

import (
	"context"
	"errors"

	"my_backend/internal/domain"
)

type lixiService struct {
	lixiRepo domain.LixiRepository
}

func NewLixiService(lixiRepo domain.LixiRepository) domain.LixiService {
	return &lixiService{
		lixiRepo: lixiRepo,
	}
}

func (s *lixiService) CreateConfig(ctx context.Context, name string, envelopes []domain.LixiEnvelope) (*domain.LixiConfig, error) {
	if name == "" {
		return nil, errors.New("name is required")
	}

	if len(envelopes) != 12 {
		return nil, errors.New("exactly 12 envelopes are required")
	}

	// Validate each envelope
	for i, env := range envelopes {
		if env.Amount == "" {
			return nil, errors.New("amount is required for all envelopes")
		}
		if env.Message == "" {
			return nil, errors.New("message is required for all envelopes")
		}
		if env.Rate <= 0 {
			return nil, errors.New("rate must be greater than 0 for all envelopes")
		}
		// Set ID based on position (1-12)
		envelopes[i].ID = i + 1
	}

	config := &domain.LixiConfig{
		Name:      name,
		Envelopes: envelopes,
		IsActive:  false,
	}

	if err := s.lixiRepo.Create(ctx, config); err != nil {
		return nil, err
	}

	return config, nil
}

func (s *lixiService) GetActiveConfig(ctx context.Context) (*domain.LixiConfig, error) {
	return s.lixiRepo.GetActive(ctx)
}

func (s *lixiService) GetAllConfigs(ctx context.Context) ([]*domain.LixiConfig, error) {
	return s.lixiRepo.GetAll(ctx)
}

func (s *lixiService) UpdateConfig(ctx context.Context, id string, name string, envelopes []domain.LixiEnvelope) (*domain.LixiConfig, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	// Get existing config
	config, err := s.lixiRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if name != "" {
		config.Name = name
	}

	if len(envelopes) > 0 {
		if len(envelopes) != 12 {
			return nil, errors.New("exactly 12 envelopes are required")
		}

		// Validate each envelope
		for i, env := range envelopes {
			if env.Amount == "" {
				return nil, errors.New("amount is required for all envelopes")
			}
			if env.Message == "" {
				return nil, errors.New("message is required for all envelopes")
			}
			if env.Rate <= 0 {
				return nil, errors.New("rate must be greater than 0 for all envelopes")
			}
			// Set ID based on position (1-12)
			envelopes[i].ID = i + 1
		}
		config.Envelopes = envelopes
	}

	if err := s.lixiRepo.Update(ctx, config); err != nil {
		return nil, err
	}

	return config, nil
}

func (s *lixiService) DeleteConfig(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is required")
	}

	// Check if config exists and is not active
	config, err := s.lixiRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if config.IsActive {
		return errors.New("cannot delete active config")
	}

	return s.lixiRepo.Delete(ctx, id)
}

func (s *lixiService) SetActiveConfig(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id is required")
	}

	// Verify config exists
	_, err := s.lixiRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.lixiRepo.SetActive(ctx, id)
}
