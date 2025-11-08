package repository

import (
	"context"
	"encoding/json"
	"time"
)

const (
	PageEntity   = "page"
	SharedEntity = "shared"

	MainContentKey = "content"
)

type SchemaEntity struct {
	ID        int64           `json:"id"`
	Title     string          `json:"title"`
	Type      string          `json:"type"`
	Slug      string          `json:"slug"`
	Content   json.RawMessage `json:"content"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type SchemaEntityRepository interface {
	GetBySlug(ctx context.Context, slug string) (*SchemaEntity, error)
	GetAll(ctx context.Context) ([]*SchemaEntity, error)
	Upsert(ctx context.Context, page *SchemaEntity) error
	DeleteBySlug(ctx context.Context, slug string) error
}
