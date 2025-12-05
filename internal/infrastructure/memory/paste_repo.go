package memory

import (
	"errors"
	"sync"
	"trace/internal/domain"
)

// PasteRepository — реализация хранилища в памяти.
// Используем RWMutex для безопасного конкурентного доступа.
type PasteRepository struct {
	mu   sync.RWMutex
	data map[string]*domain.Paste
}

// NewPasteRepository создает новый экземпляр.
func NewPasteRepository() *PasteRepository {
	return &PasteRepository{
		data: make(map[string]*domain.Paste),
	}
}

// Store сохраняет пасту.
func (r *PasteRepository) Store(p *domain.Paste) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// В реальной БД здесь была бы проверка на коллизии,
	// но для in-memory считаем, что генератор ID справляется.
	r.data[p.ID] = p
	return nil
}

// GetByID возвращает пасту по ID.
func (r *PasteRepository) GetByID(id string) (*domain.Paste, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	paste, exists := r.data[id]
	if !exists {
		return nil, errors.New("paste not found")
	}

	return paste, nil
}

// Delete удаляет пасту.
func (r *PasteRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.data, id)
	return nil
}
