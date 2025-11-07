package kvstore

import "context"

const (
	SchemaKeySchema  = "schema:schema"
	SchemaKeyPages   = "schema:pages"
	SchemaKeyLayouts = "schema:layouts"
	SchemaKeyModules = "schema:modules"
	SchemaKeyBlocks  = "schema:blocks"
	SchemaKeyShared  = "schema:shared"

	EntityDataPrefix = "entity:"

	PagesListCacheKey  = "pages_list"
	SharedListCacheKey = "shared_list"
)

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

	Get(ctx context.Context, key string) (any, error)
	Set(ctx context.Context, key string, data any) error
	Delete(ctx context.Context, key string) error

	GetEntityNamesByType(ctx context.Context, entityType string) ([]string, error)
	GetEntitySchema(ctx context.Context, prefix string, slug string) (*EntitySchema, error)
	GetEntityData(ctx context.Context, key string) (map[string]interface{}, error)
	SetEntityData(ctx context.Context, key string, data map[string]interface{}) error
	DeleteEntityData(ctx context.Context, key string) error
}
