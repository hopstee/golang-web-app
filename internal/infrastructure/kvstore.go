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

type schemaFile struct {
	Path string
	Dst  *[]*kvstore.EntitySchema
}

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
	var (
		schema, pages, layouts, blocks, modules, shared []*kvstore.EntitySchema
		files                                           = []schemaFile{
			{config.Schemas.Schema, &schema},
			{config.Schemas.Pages, &pages},
			{config.Schemas.Layouts, &layouts},
			{config.Schemas.Blocks, &blocks},
			{config.Schemas.Modules, &modules},
			{config.Schemas.Shared, &shared},
		}
	)

	for _, f := range files {
		if err := readAndUnmarshal(f); err != nil {
			return err
		}
	}

	schemas := &kvstore.SchemasList{
		Schema:  schema,
		Pages:   pages,
		Layouts: layouts,
		Blocks:  blocks,
		Modules: modules,
		Shared:  shared,
	}

	if err := kvs.SetSchemas(ctx, schemas); err != nil {
		return err
	}

	fmt.Println("schema loaded")
	return nil
}

func readAndUnmarshal(file schemaFile) error {
	data, err := os.ReadFile(file.Path)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", file.Path, err)
	}

	var schemas []*kvstore.EntitySchema
	if err := json.Unmarshal(data, &schemas); err != nil {
		return fmt.Errorf("unmarshal file %s: %w", file.Path, err)
	}

	*file.Dst = schemas
	return nil
}
