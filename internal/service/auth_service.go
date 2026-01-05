package service

import (
	"context"
	"errors"
	"time"

	"my_backend/internal/domain"

	"github.com/golang-jwt/jwt/v5"
)

type authService struct {
	userRepo  domain.UserRepository
	jwtSecret []byte
}

func NewAuthService(userRepo domain.UserRepository, jwtSecret string) domain.AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: []byte(jwtSecret),
	}
}

func (s *authService) Register(ctx context.Context, email, password string) (*domain.User, error) {
	// In production, you MUST hash the password here (e.g. using bcrypt).
	// For this example, we store it as plain text (Do NOT do this in real apps).
	user := &domain.User{
		ID:       time.Now().String(), // Simple ID generation
		Email:    email,
		Password: password,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (*domain.User, string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}

	if user.Password != password {
		return nil, "", errors.New("invalid credentials")
	}

	// Create JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return nil, "", err
	}

	return user, tokenString, nil
}
