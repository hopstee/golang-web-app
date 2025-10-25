package kvstore

import "context"

type Field struct {
	Name   string  `json:"name"`
	Type   string  `json:"type"`
	Label  string  `json:"label,omitempty"`
	Schema *Schema `json:"schema,omitempty"`
}

type Schema struct {
	Fields []Field `json:"fields"`
}

type EntitySchema struct {
	ID       string          `json:"id"`
	Title    string          `json:"title"`
	Type     string          `json:"type"` // layout / block / page / module
	Layout   string          `json:"layout,omitempty"`
	Parent   string          `json:"parent,omitempty"`
	Blocks   []string        `json:"blocks,omitempty"`
	SEO      []Field         `json:"seo,omitempty"`
	Content  []Field         `json:"content,omitempty"`
	Children []*EntitySchema `json:"children,omitempty"`
}

type EntityData struct {
}

type KVStore interface {
	SetSchemas(ctx context.Context, pages, modules []*EntitySchema) error
	GetPages(ctx context.Context) ([]*EntitySchema, error)
	GetModules(ctx context.Context) ([]*EntitySchema, error)

	SetPageData(ctx context.Context, slug string, data map[string]interface{}) error
	GetPageData(ctx context.Context, slug string) (map[string]interface{}, error)
	DeletePageData(ctx context.Context, slug string) error
}
