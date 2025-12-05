package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"trace/internal/domain"

	_ "github.com/jackc/pgx/v5/stdlib" // Регистрация драйвера pgx
)

type PasteRepository struct {
	db *sql.DB
}

func NewPasteRepository(db *sql.DB) *PasteRepository {
	return &PasteRepository{db: db}
}

func (r *PasteRepository) Store(p *domain.Paste) error {
	query := `
		INSERT INTO pastes (id, content, language, hash, views, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	// Используем context с таймаутом для безопасности
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(ctx, query,
		p.ID,
		p.Content,
		p.Language,
		p.Hash,
		p.Views,
		p.ExpiresAt,
		p.CreatedAt,
	)
	return err
}

func (r *PasteRepository) GetByID(id string) (*domain.Paste, error) {
	query := `
		SELECT id, content, language, hash, views, expires_at, created_at
		FROM pastes WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := r.db.QueryRowContext(ctx, query, id)

	var p domain.Paste
	// Сканируем данные. expires_at может быть NULL, поэтому используем sql.NullTime или указатель
	// Но pgx/stdlib умеет мапить NULL в *time.Time, попробуем напрямую.
	// Если возникнут проблемы, можно использовать sql.NullTime для expires_at.

	err := row.Scan(
		&p.ID,
		&p.Content,
		&p.Language,
		&p.Hash,
		&p.Views,
		&p.ExpiresAt,
		&p.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("paste not found") // Можно создать кастомную ошибку domain.ErrNotFound
		}
		return nil, err
	}

	return &p, nil
}

func (r *PasteRepository) Delete(id string) error {
	query := `DELETE FROM pastes WHERE id = $1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
