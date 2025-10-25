package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"mobile-backend-boilerplate/internal/config"
	"mobile-backend-boilerplate/internal/kvstore"
	"mobile-backend-boilerplate/internal/kvstore/redis"
	"os"
)

func (d *Dependencies) InitKVStore(ctx context.Context) error {
	var kvStore kvstore.KVStore

	switch d.Config.KVStore.Type {
	case "redis":
		kvStore = redis.NewRedisKVStore(
			d.Config.KVStore.Redis.Addr,
			d.Config.KVStore.Redis.Password,
			d.Config.KVStore.Redis.DB,
		)
	default:
		return fmt.Errorf("unsupported kvstore type: %s", d.Config.KVStore.Type)
	}

	if err := loadSchemas(ctx, kvStore, d.Config); err != nil {
		return err
	}

	d.KVStore = kvStore
	return nil
}

func loadSchemas(ctx context.Context, kvs kvstore.KVStore, config *config.Config) error {
	pagesData, err := os.ReadFile(config.Schemas.Pages)
	if err != nil {
		return fmt.Errorf("failed to read pages schema: %w", err)
	}

	var pages []*kvstore.EntitySchema
	if err := json.Unmarshal(pagesData, &pages); err != nil {
		return fmt.Errorf("unmarshal pages schema: %w", err)
	}

	modulesData, err := os.ReadFile(config.Schemas.Modules)
	if err != nil {
		return fmt.Errorf("failed to read modules schema: %w", err)
	}

	var modules []*kvstore.EntitySchema
	if err := json.Unmarshal(modulesData, &modules); err != nil {
		return fmt.Errorf("unmarshal modules schema: %w", err)
	}

	if err := kvs.SetSchemas(ctx, pages, modules); err != nil {
		return err
	}

	return nil
}
