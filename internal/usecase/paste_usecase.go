package usecase

import (
	"errors"
	"time"
	"trace/internal/domain"
)

// IDGenerator - абстракция для генерации уникальных ID.
type IDGenerator interface {
	Generate() string
}

type pasteUseCase struct {
	repo    domain.PasteRepository
	idGen   IDGenerator
	timeout time.Duration
}

// NewPasteUseCase - конструктор.
func NewPasteUseCase(repo domain.PasteRepository, idGen IDGenerator, timeout time.Duration) domain.PasteUseCase {
	return &pasteUseCase{
		repo:    repo,
		idGen:   idGen,
		timeout: timeout,
	}
}

func (uc *pasteUseCase) Create(content, language string, ttl time.Duration) (*domain.Paste, error) {
	if content == "" {
		return nil, errors.New("content cannot be empty")
	}

	id := uc.idGen.Generate()
	now := time.Now()

	paste := &domain.Paste{
		ID:        id,
		Content:   content,
		Language:  language,
		CreatedAt: now,
		Views:     0,
	}

	if ttl > 0 {
		exp := now.Add(ttl)
		paste.ExpiresAt = &exp
	}

	if err := uc.repo.Store(paste); err != nil {
		return nil, err
	}

	return paste, nil
}

func (uc *pasteUseCase) Get(id string) (*domain.Paste, error) {
	paste, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Здесь можно добавить инкремент просмотров асинхронно

	return paste, nil
}
