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

func NewRedisKVStore(addr, password string, db int) *RedisKVStore {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	return &RedisKVStore{client: rdb}
}

func (r *RedisKVStore) SetSchemas(ctx context.Context, pages, modules []*kvstore.EntitySchema) error {
	if err := r.setJSON(ctx, "schema:pages", pages); err != nil {
		return fmt.Errorf("set pages: %w", err)
	}
	if err := r.setJSON(ctx, "schema:modules", modules); err != nil {
		return fmt.Errorf("set modules: %w", err)
	}
	return nil
}

func (r *RedisKVStore) GetPages(ctx context.Context) (pages []*kvstore.EntitySchema, err error) {
	if err := r.getJSON(ctx, "schema:pages", &pages); err != nil {
		return nil, err
	}
	return pages, nil
}

func (r *RedisKVStore) GetModules(ctx context.Context) (modules []*kvstore.EntitySchema, err error) {
	if err := r.getJSON(ctx, "schema:modules", &modules); err != nil {
		return nil, err
	}
	return modules, nil
}

func (r *RedisKVStore) SetPageData(ctx context.Context, slug string, data map[string]interface{}) error {
	key := fmt.Sprintf("page:%s", slug)
	return r.setJSON(ctx, key, data)
}

func (r *RedisKVStore) GetPageData(ctx context.Context, slug string) (map[string]interface{}, error) {
	key := fmt.Sprintf("page:%s", slug)
	var data map[string]interface{}
	if err := r.getJSON(ctx, key, &data); err != nil {
		return nil, err
	}
	return data, nil
}

func (r *RedisKVStore) DeletePageData(ctx context.Context, slug string) error {
	key := fmt.Sprintf("page:%s", slug)
	return r.client.Del(ctx, key).Err()
}

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
