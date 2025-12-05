package domain

import (
	"time"
)

// Paste - основная сущность приложения.
type Paste struct {
	ID        string     `json:"id"`
	Content   string     `json:"content"`
	Language  string     `json:"language"` // например: "go", "javascript", "plaintext"
	Hash      string     `json:"hash"`     // Хэш содержимого для дедупликации (опционально)
	Views     int64      `json:"views"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

// PasteRepository определяет интерфейс для хранения данных.
// Следуем принципу DIP: бизнес-логика зависит от абстракции.
type PasteRepository interface {
	Store(paste *Paste) error
	GetByID(id string) (*Paste, error)
	Delete(id string) error
}

// PasteUseCase определяет бизнес-операции.
type PasteUseCase interface {
	Create(content, language string, ttl time.Duration) (*Paste, error)
	Get(id string) (*Paste, error)
}
