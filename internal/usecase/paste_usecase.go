package usecase

import (
	"errors"
	"time"
	"trace/internal/domain"

	"golang.org/x/crypto/bcrypt"
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

const MaxDuration = 14 * 24 * time.Hour // 14 дней

var ErrPasswordRequired = errors.New("password required")
var ErrInvalidPassword = errors.New("invalid password")

func (uc *pasteUseCase) Create(content, language, password string, ttl time.Duration, burn bool) (*domain.Paste, error) {
	if content == "" {
		return nil, errors.New("content cannot be empty")
	}

	// Ограничиваем TTL
	if ttl > MaxDuration {
		ttl = MaxDuration
	}
	// Если TTL не задан (0), ставим дефолт (например 14 дней), чтобы не хранить вечно
	if ttl == 0 {
		ttl = MaxDuration
	}

	id := uc.idGen.Generate()
	now := time.Now()

	paste := &domain.Paste{
		ID:               id,
		Content:          content,
		Language:         language,
		CreatedAt:        now,
		Views:            0,
		BurnAfterReading: burn,
	}

	// Хеширование пароля
	if password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		paste.PasswordHash = string(hash)
		paste.IsProtected = true
	}

	exp := now.Add(ttl)
	paste.ExpiresAt = &exp

	if err := uc.repo.Store(paste); err != nil {
		return nil, err
	}

	return paste, nil
}

func (uc *pasteUseCase) Get(id, password string) (*domain.Paste, error) {
	paste, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Проверка пароля
	if paste.IsProtected {
		if password == "" {
			// Возвращаем пасту БЕЗ контента, но с флагом, что нужен пароль
			return &domain.Paste{
				ID:          paste.ID,
				Language:    paste.Language,
				IsProtected: true,
				CreatedAt:   paste.CreatedAt,
				// Content пустой!
			}, ErrPasswordRequired
		}

		// Сверяем хеш
		if err := bcrypt.CompareHashAndPassword([]byte(paste.PasswordHash), []byte(password)); err != nil {
			return nil, ErrInvalidPassword
		}
	}

	// Логика сжигания
	if paste.BurnAfterReading {
		// Удаляем асинхронно, чтобы не тормозить отдачу (или синхронно для надежности)
		_ = uc.repo.Delete(id)
	}

	return paste, nil
}
