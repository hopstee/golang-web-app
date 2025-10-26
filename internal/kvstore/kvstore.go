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
	ID           string          `json:"id"`
	Title        string          `json:"title"`
	Type         string          `json:"type"` // layout / block / page / module / shared
	Layout       string          `json:"layout,omitempty"`
	Parent       string          `json:"parent,omitempty"`
	Mode         string          `json:"mode,omitempty"`
	Blocks       []string        `json:"blocks,omitempty"`
	Refs         []string        `json:"refs,omitempty"`
	LayoutFields []Field         `json:"layout_fields,omitempty"`
	Content      []Field         `json:"content,omitempty"`
	Children     []*EntitySchema `json:"children,omitempty"`
}

type SchemasList struct {
	Schema  []*EntitySchema
	Pages   []*EntitySchema
	Layouts []*EntitySchema
	Modules []*EntitySchema
	Blocks  []*EntitySchema
	Shared  []*EntitySchema
}

type KVStore interface {
	SetSchemas(ctx context.Context, schemas *SchemasList) error
	GetPages(ctx context.Context) ([]*EntitySchema, error)
	GetModules(ctx context.Context) ([]*EntitySchema, error)

	SetPageData(ctx context.Context, slug string, data map[string]interface{}) error
	GetPageData(ctx context.Context, slug string) (map[string]interface{}, error)
	DeletePageData(ctx context.Context, slug string) error
}
