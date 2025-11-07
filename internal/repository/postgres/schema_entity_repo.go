package postgres

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"mobile-backend-boilerplate/internal/repository"
	"time"
)

type schemaEntityRepo struct {
	DB     *sql.DB
	Logger *slog.Logger
}

func NewSchemaEntityRepo(db *sql.DB, logger *slog.Logger) *schemaEntityRepo {
	return &schemaEntityRepo{
		DB:     db,
		Logger: logger,
	}
}

func (r *schemaEntityRepo) GetBySlug(ctx context.Context, slug string) (*repository.SchemaEntity, error) {
	query := `
		SELECT id, title, slug, content, created_at, updated_at
		FROM schemaEntities
		WHERE slug = $1
		LIMIT 1;
	`

	var p repository.SchemaEntity
	err := r.DB.QueryRowContext(ctx, query, slug).Scan(
		&p.ID, &p.Title, &p.Slug, &p.Content, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.Logger.Error("schemaEntity not found", slog.String("slug", slug))
			return nil, repository.ErrNotFound
		}
		r.Logger.Error("failed to query schemaEntity", slog.String("slug", slug))
		return nil, err
	}

	return &p, nil
}

func (r *schemaEntityRepo) GetAll(ctx context.Context) ([]*repository.SchemaEntity, error) {
	query := `
		SELECT id, title, slug, content, created_at, updated_at
		FROM schemaEntities
		ORDER BY created_at DESC;
	`

	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		r.Logger.Error("failed to get all schemaEntities", slog.Any("err", err))
		return nil, err
	}
	defer rows.Close()

	var schemaEntities []*repository.SchemaEntity
	for rows.Next() {
		var p repository.SchemaEntity
		if err := rows.Scan(&p.ID, &p.Title, &p.Slug, &p.Content, &p.CreatedAt, &p.UpdatedAt); err != nil {
			r.Logger.Error("failed to scan schemaEntity", slog.Any("err", err))
			return nil, err
		}
		schemaEntities = append(schemaEntities, &p)
	}

	return schemaEntities, nil
}

func (r *schemaEntityRepo) Upsert(ctx context.Context, schemaEntity *repository.SchemaEntity) error {
	var id int64
	var updatedAt time.Time

	query := `
		INSERT INTO schemaEntities (title, slug, content, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (slug)
		DO UPDATE SET
			title = EXCLUDED.title,
			content = EXCLUDED.content,
			updated_at = EXCLUDED.updated_at
		RETURNING id, updated_at
	`

	err := r.DB.QueryRowContext(
		ctx, query, schemaEntity.Title, schemaEntity.Slug, schemaEntity.Content, time.Now(), time.Now(),
	).Scan(&id, &updatedAt)
	if err != nil {
		r.Logger.Error("failed to upsert schemaEntity", slog.Any("err", err))
		return err
	}

	schemaEntity.ID = id
	schemaEntity.UpdatedAt = updatedAt
	return nil
}

func (r *schemaEntityRepo) DeleteBySlug(ctx context.Context, slug string) error {
	query := `
		DELETE FROM schemaEntities
		WHERE slug = $1
	`
	_, err := r.DB.ExecContext(ctx, query, slug)
	if err != nil {
		r.Logger.Error("failed to delete schemaEntity", slog.Any("err", err))
		return err
	}

	return nil
}
