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

	modules := []*kvstore.EntitySchema{
		{
			ID:    "module1",
			Title: "Module One",
			Type:  "module",
		},
	}

	// Set schemas
	if err := kv.SetSchemas(ctx, pages, modules); err != nil {
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

	// Get modules
	gotModules, err := kv.GetModules(ctx)
	if err != nil {
		t.Fatalf("GetModules failed: %v", err)
	}

	if !reflect.DeepEqual(modules, gotModules) {
		t.Errorf("modules mismatch:\nwant: %+v\ngot: %+v", modules, gotModules)
	}
}

func TestRedisKVStore_GetEmpty(t *testing.T) {
	ctx := context.Background()
	kv := NewRedisKVStore(getTestRedisAddr(), getTestRedisPassword(), 0)

	// Очистим ключи
	kv.SetSchemas(ctx, []*kvstore.EntitySchema{}, []*kvstore.EntitySchema{})

	pages, err := kv.GetPages(ctx)
	if err != nil {
		t.Fatalf("GetPages failed: %v", err)
	}
	if len(pages) != 0 {
		t.Errorf("expected no pages, got %d", len(pages))
	}

	modules, err := kv.GetModules(ctx)
	if err != nil {
		t.Fatalf("GetModules failed: %v", err)
	}
	if len(modules) != 0 {
		t.Errorf("expected no modules, got %d", len(modules))
	}
}
