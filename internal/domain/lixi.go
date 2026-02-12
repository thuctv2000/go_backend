package domain

import (
	"context"
	"time"
)

type LixiEnvelope struct {
	ID      int     `json:"id"`
	Amount  string  `json:"amount"`  // "100K VNĐ", "1 Triệu VNĐ"
	Message string  `json:"message"` // "Phát Tài Phát Lộc!"
	Rate    float64 `json:"rate"`    // Probability weight (e.g., 0.5 = 50% chance)
}

type LixiConfig struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`      // "Tết 2025"
	Envelopes []LixiEnvelope  `json:"envelopes"` // 12 envelopes
	IsActive  bool            `json:"is_active"`
	CreatedAt time.Time       `json:"created_at"`
}

type LixiRepository interface {
	Create(ctx context.Context, config *LixiConfig) error
	GetActive(ctx context.Context) (*LixiConfig, error)
	GetByID(ctx context.Context, id string) (*LixiConfig, error)
	GetAll(ctx context.Context) ([]*LixiConfig, error)
	Update(ctx context.Context, config *LixiConfig) error
	Delete(ctx context.Context, id string) error
	SetActive(ctx context.Context, id string) error
}

type LixiService interface {
	CreateConfig(ctx context.Context, name string, envelopes []LixiEnvelope) (*LixiConfig, error)
	GetActiveConfig(ctx context.Context) (*LixiConfig, error)
	GetAllConfigs(ctx context.Context) ([]*LixiConfig, error)
	UpdateConfig(ctx context.Context, id string, name string, envelopes []LixiEnvelope) (*LixiConfig, error)
	DeleteConfig(ctx context.Context, id string) error
	SetActiveConfig(ctx context.Context, id string) error
	SubmitGreeting(ctx context.Context, name, amount, message, image string) (*LixiGreeting, error)
	GetAllGreetings(ctx context.Context) ([]*LixiGreeting, error)
}

type LixiGreeting struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Amount    string    `json:"amount"`
	Message   string    `json:"message"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"created_at"`
}

type LixiGreetingRepository interface {
	Create(ctx context.Context, greeting *LixiGreeting) error
	GetAll(ctx context.Context) ([]*LixiGreeting, error)
}
