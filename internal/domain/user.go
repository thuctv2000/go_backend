package domain

import "context"

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"` // omitempty allows hiding password in response
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
}

type AuthService interface {
	Register(ctx context.Context, email, password string) (*User, error)
	Login(ctx context.Context, email, password string) (*User, string, error) // Returns User and JWT token
}
