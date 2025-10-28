package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"mobile-backend-boilerplate/internal/repository"
	"time"
)

type pagesRepo struct {
	DB     *sql.DB
	Logger *slog.Logger
}

func NewPagesRepo(db *sql.DB, logger *slog.Logger) *pagesRepo {
	return &pagesRepo{
		DB:     db,
		Logger: logger,
	}
}

func (r *pagesRepo) GetBySlug(ctx context.Context, slug string) (*repository.Page, error) {
	query := `
		SELECT id, title, slug, content, created_at, updated_at
		FROM pages
		WHERE slug = $1
		LIMIT 1;
	`

	var p repository.Page
	err := r.DB.QueryRowContext(ctx, query, slug).Scan(
		&p.ID, &p.Title, &p.Slug, &p.Content, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.Logger.Error("page not found", slog.String("slug", slug))
			return nil, repository.ErrNotFound
		}
		r.Logger.Error("failed to query page", slog.String("slug", slug))
		return nil, err
	}

	return &p, nil
}

func (r *pagesRepo) GetAll(ctx context.Context) ([]*repository.Page, error) {
	query := `
		SELECT id, title, slug, content, created_at, updated_at
		FROM pages
		ORDER BY created_at DESC;
	`

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		r.Logger.Error("failed to get all pages", slog.Any("err", err))
		return nil, err
	}
	defer rows.Close()

	var pages []*repository.Page
	for rows.Next() {
		var p repository.Page
		if err := rows.Scan(&p.ID, &p.Title, &p.Slug, &p.Content, &p.CreatedAt, &p.UpdatedAt); err != nil {
			r.Logger.Error("failed to scan page", slog.Any("err", err))
			return nil, err
		}
		pages = append(pages, &p)
	}

	return pages, nil
}

func (r *pagesRepo) Upsert(ctx context.Context, page *repository.Page) error {
	var id int64
	var updatedAt time.Time

	query := `
		INSERT INTO pages (title, slug, content, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (slug)
		DO UPDATE SET
			title = EXCLUDED.title,
			content = EXCLUDED.content,
			updated_at = EXCLUDED.updated_at
		RETURNING id, updated_at
	`

	err := r.DB.QueryRowContext(
		ctx, query, page.Title, page.Slug, page.Content, time.Now(), time.Now(),
	).Scan(&id, &updatedAt)
	if err != nil {
		r.Logger.Error("failed to upsert page", slog.Any("err", err))
		return err
	}

	page.ID = id
	page.UpdatedAt = updatedAt
	return nil
}

func (r *pagesRepo) DeleteBySlug(ctx context.Context, slug string) error {
	query := `
		DELETE FROM pages
		WHERE slug = $1
	`
	_, err := r.DB.ExecContext(ctx, query, slug)
	if err != nil {
		r.Logger.Error("failed to delete page", slog.Any("err", err))
		return err
	}

	return nil
}
