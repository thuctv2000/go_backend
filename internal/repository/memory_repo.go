package repository

import (
	"context"
	"errors"
	"my_backend/internal/domain"
	"sync"
)

type memoryUserRepository struct {
	users map[string]*domain.User
	mu    sync.RWMutex
}

func NewMemoryUserRepository() domain.UserRepository {
	return &memoryUserRepository{
		users: make(map[string]*domain.User),
	}
}

func (r *memoryUserRepository) Create(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.Email]; exists {
		return errors.New("user already exists")
	}

	r.users[user.Email] = user
	return nil
}

func (r *memoryUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[email]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}
