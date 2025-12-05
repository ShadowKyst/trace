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
		INSERT INTO pastes (id, content, language, hash, views, expires_at, created_at, password_hash, burn_after_reading)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Обработка NULL для password_hash
	var pwdHash sql.NullString
	if p.PasswordHash != "" {
		pwdHash = sql.NullString{String: p.PasswordHash, Valid: true}
	}

	_, err := r.db.ExecContext(ctx, query,
		p.ID, p.Content, p.Language, p.Hash, p.Views, p.ExpiresAt, p.CreatedAt,
		pwdHash, p.BurnAfterReading,
	)
	return err
}

func (r *PasteRepository) GetByID(id string) (*domain.Paste, error) {
	// Добавляем проверку на expires_at в сам запрос.
	// Если паста протухла, база её просто не найдет.
	query := `
		SELECT id, content, language, hash, views, expires_at, created_at, password_hash, burn_after_reading
		FROM pastes 
		WHERE id = $1 AND (expires_at IS NULL OR expires_at > NOW())
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := r.db.QueryRowContext(ctx, query, id)

	var p domain.Paste
	var pwdHash sql.NullString

	err := row.Scan(
		&p.ID, &p.Content, &p.Language, &p.Hash, &p.Views, &p.ExpiresAt, &p.CreatedAt,
		&pwdHash, &p.BurnAfterReading,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("paste not found or expired")
		}
		return nil, err
	}

	if pwdHash.Valid {
		p.PasswordHash = pwdHash.String
		p.IsProtected = true
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

func (r *PasteRepository) Cleanup() error {
	query := `DELETE FROM pastes WHERE expires_at < NOW()`
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := r.db.ExecContext(ctx, query)
	return err
}
