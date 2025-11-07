package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"mobile-backend-boilerplate/internal/kvstore"

	"github.com/redis/go-redis/v9"
)

type RedisKVStore struct {
	client *redis.Client
}

type KVData struct {
	Key  string
	Data []*kvstore.EntitySchema
}

func NewRedisKVStore(addr, password string, db int) *RedisKVStore {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisKVStore{client: rdb}
}

func (r *RedisKVStore) SetSchemas(ctx context.Context, schemas *kvstore.SchemasList) error {
	kvData := []KVData{
		{kvstore.SchemaKeySchema, schemas.Schema},
		{kvstore.SchemaKeyPages, schemas.Pages},
		{kvstore.SchemaKeyLayouts, schemas.Layouts},
		{kvstore.SchemaKeyModules, schemas.Modules},
		{kvstore.SchemaKeyBlocks, schemas.Blocks},
		{kvstore.SchemaKeyShared, schemas.Shared},
	}

	for _, schema := range kvData {
		if err := r.client.Del(ctx, schema.Key).Err(); err != nil {
			return fmt.Errorf("delete %s: %v", schema.Key, err)
		}
		if err := r.setJSON(ctx, schema.Key, schema.Data); err != nil {
			return fmt.Errorf("set %s: %w", schema.Key, err)
		}
	}

	var pagesList []string
	for _, page := range schemas.Pages {
		key := kvstore.EntityDataPrefix + page.ID
		if err := r.client.Del(ctx, key).Err(); err != nil {
			return fmt.Errorf("delete page data %s: %v", key, err)
		}
		if err := r.setJSON(ctx, key, page.Content); err != nil {
			return fmt.Errorf("set %s: %w", key, err)
		}
		pagesList = append(pagesList, page.Title)
	}

	err := r.Set(ctx, kvstore.PagesListCacheKey, pagesList)
	if err != nil {
		return fmt.Errorf("set %s: %w", kvstore.PagesListCacheKey, err)
	}

	var sharedList []string
	for _, shared := range schemas.Shared {
		key := kvstore.EntityDataPrefix + shared.ID
		if err := r.client.Del(ctx, key).Err(); err != nil {
			return fmt.Errorf("delete shared data %s: %v", key, err)
		}
		if err := r.setJSON(ctx, key, shared.Content); err != nil {
			return fmt.Errorf("set %s: %w", key, err)
		}
		sharedList = append(sharedList, shared.Title)
	}

	err = r.Set(ctx, kvstore.SharedListCacheKey, sharedList)
	if err != nil {
		return fmt.Errorf("set %s: %w", kvstore.SharedListCacheKey, err)
	}

	return nil
}

func (r *RedisKVStore) Get(ctx context.Context, key string) (data interface{}, err error) {
	if err = r.getJSON(ctx, key, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (r *RedisKVStore) Set(ctx context.Context, key string, data interface{}) error {
	return r.setJSON(ctx, key, data)
}

func (r *RedisKVStore) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// Helper methods for entites
func (r *RedisKVStore) GetEntityNamesByType(ctx context.Context, entityType string) ([]string, error) {
	var entityNames []string
	if err := r.getJSON(ctx, entityType, &entityNames); err != nil {
		return nil, err
	}
	return entityNames, nil
}

func (r *RedisKVStore) GetEntitySchema(ctx context.Context, key string, slug string) (*kvstore.EntitySchema, error) {
	var schemas []*kvstore.EntitySchema
	if err := r.getJSON(ctx, key, &schemas); err != nil {
		return nil, err
	}

	var schema *kvstore.EntitySchema
	for _, s := range schemas {
		if s.ID == slug {
			schema = s
			break
		}
	}

	if schema == nil {
		return nil, fmt.Errorf("schema not found")
	}

	return schema, nil
}

func (r *RedisKVStore) GetEntityData(ctx context.Context, key string) (map[string]interface{}, error) {
	var data map[string]interface{}
	if err := r.getJSON(ctx, key, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (r *RedisKVStore) SetEntityData(ctx context.Context, key string, data map[string]interface{}) error {
	return r.setJSON(ctx, key, data)
}

func (r *RedisKVStore) DeleteEntityData(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

// Utility methods
func (r *RedisKVStore) setJSON(ctx context.Context, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, data, 0).Err()
}

func (r *RedisKVStore) getJSON(ctx context.Context, key string, out interface{}) error {
	value, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(value, out)
}
