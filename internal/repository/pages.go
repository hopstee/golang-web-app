package repository

import (
	"context"
	"encoding/json"
	"time"
)

type Page struct {
	ID        int64           `json:"id"`
	Title     string          `json:"title"`
	Slug      string          `json:"slug"`
	Content   json.RawMessage `json:"content"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type PagesRepository interface {
	GetBySlug(ctx context.Context, slug string) (*Page, error)
	GetAll(ctx context.Context) ([]*Page, error)
	Upsert(ctx context.Context, page *Page) error
	DeleteBySlug(ctx context.Context, slug string) error
}
