package domain

import (
	"time"
)

// Paste - основная сущность приложения.
type Paste struct {
	ID               string     `json:"id"`
	Content          string     `json:"content"`
	Language         string     `json:"language"` // например: "go", "javascript", "plaintext"
	Hash             string     `json:"hash"`     // Хэш содержимого для дедупликации (опционально)
	Views            int64      `json:"views"`
	ExpiresAt        *time.Time `json:"expires_at,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	PasswordHash     string     `json:"-"`
	BurnAfterReading bool       `json:"burn_after_reading"`
	IsProtected      bool       `json:"is_protected"`
}

// PasteRepository определяет интерфейс для хранения данных.
// Следуем принципу DIP: бизнес-логика зависит от абстракции.
type PasteRepository interface {
	Store(paste *Paste) error
	GetByID(id string) (*Paste, error)
	Delete(id string) error
	Cleanup() error
}

// PasteUseCase определяет бизнес-операции.
type PasteUseCase interface {
	// Обновили сигнатуру Create
	Create(content, language, password string, ttl time.Duration, burn bool) (*Paste, error)
	// Обновили сигнатуру Get (добавили password для проверки)
	Get(id, password string) (*Paste, error)
}
