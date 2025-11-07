package redis

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"

	"mobile-backend-boilerplate/internal/kvstore"
)

func getTestRedisAddr() string {
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("REDIS_PORT")
	if port == "" {
		port = "6379"
	}

	return host + ":" + port
}

func getTestRedisPassword() string {
	fmt.Println("REDIS_PASSWORD:", os.Getenv("REDIS_PASSWORD"))
	return os.Getenv("REDIS_PASSWORD")
}

func TestRedisKVStore_SetGetSchemas(t *testing.T) {
	ctx := context.Background()
	kv := NewRedisKVStore(getTestRedisAddr(), getTestRedisPassword(), 0)

	pages := []*kvstore.EntitySchema{
		{
			ID:    "page1",
			Title: "Page One",
			Type:  "page",
		},
	}

	layouts := []*kvstore.EntitySchema{
		{
			ID:    "layout1",
			Title: "Layout One",
			Type:  "layout",
		},
	}

	modules := []*kvstore.EntitySchema{
		{
			ID:    "module1",
			Title: "Module One",
			Type:  "module",
		},
	}

	schemas := &kvstore.SchemasList{
		Pages:   pages,
		Layouts: layouts,
		Modules: modules,
	}

	// Set schemas
	if err := kv.SetSchemas(ctx, schemas); err != nil {
		t.Fatalf("SetSchemas failed: %v", err)
	}

	// Get pages
	gotPages, err := kv.GetPages(ctx)
	if err != nil {
		t.Fatalf("GetPages failed: %v", err)
	}

	if !reflect.DeepEqual(pages, gotPages) {
		t.Errorf("pages mismatch:\nwant: %+v\ngot: %+v", pages, gotPages)
	}
}

func TestRedisKVStore_GetEmpty(t *testing.T) {
	ctx := context.Background()
	kv := NewRedisKVStore(getTestRedisAddr(), getTestRedisPassword(), 0)

	// Очистим ключи
	kv.SetSchemas(ctx, &kvstore.SchemasList{})

	pages, err := kv.GetPages(ctx)
	if err != nil {
		t.Fatalf("GetPages failed: %v", err)
	}
	if len(pages) != 0 {
		t.Errorf("expected no pages, got %d", len(pages))
	}
}
