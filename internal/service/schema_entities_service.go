package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"mobile-backend-boilerplate/internal/kvstore"
	"mobile-backend-boilerplate/internal/repository"
)

type SchemaEntityService struct {
	schemaEntitiesRepo repository.SchemaEntityRepository
	kvstore            kvstore.KVStore
	logger             *slog.Logger
}

func NewSchemaEntityService(schemaEntitiesRepo repository.SchemaEntityRepository, kvstore kvstore.KVStore, logger *slog.Logger) *SchemaEntityService {
	return &SchemaEntityService{
		schemaEntitiesRepo: schemaEntitiesRepo,
		kvstore:            kvstore,
		logger:             logger,
	}
}

func (s *SchemaEntityService) GetEntitiesName(ctx context.Context, entityType string) ([]kvstore.ShortEntityData, error) {
	var key string
	switch entityType {
	case repository.PageEntity:
		key = kvstore.PagesListCacheKey
	case repository.SharedEntity:
		key = kvstore.SharedListCacheKey
	default:
		return nil, fmt.Errorf("unknown entity type: %s", entityType)
	}

	names, err := s.kvstore.GetEntityDataByType(ctx, key)
	if err != nil {
		return nil, err
	}

	return names, nil
}

func (s *SchemaEntityService) GetEntitySchema(ctx context.Context, entityType string, slug string) (*kvstore.EntitySchema, error) {
	var key string
	switch entityType {
	case repository.PageEntity:
		key = kvstore.SchemaKeyPages
	case repository.SharedEntity:
		key = kvstore.SchemaKeyShared
	default:
		return nil, fmt.Errorf("unknown entity type: %s", entityType)
	}
	return s.kvstore.GetEntitySchema(ctx, key, slug)
}

func (s *SchemaEntityService) GetEntityData(ctx context.Context, slug string) (map[string]interface{}, error) {
	key := kvstore.EntityDataPrefix + slug
	data, err := s.kvstore.GetEntityData(ctx, key)
	if err == nil {
		return data, nil
	}

	schemaEntity, err := s.schemaEntitiesRepo.GetBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			emptyData := make(map[string]interface{})
			return emptyData, nil
		}
		return nil, err
	}

	var schemaEntityData map[string]interface{}
	if schemaEntity.Content != nil {
		if err := json.Unmarshal(schemaEntity.Content, &schemaEntityData); err != nil {
			schemaEntityData = make(map[string]interface{})
		}
	} else {
		schemaEntityData = make(map[string]interface{})
	}

	if err := s.kvstore.SetEntityData(ctx, key, schemaEntityData); err != nil {
		s.logger.Warn("GetSchemaEntityData: failed to set schemaEntity data to kvstore", slog.String("slug", slug))
	}

	return schemaEntityData, nil
}

func (s *SchemaEntityService) UpdateEntityData(ctx context.Context, slug string, data map[string]interface{}) error {
	content, err := json.Marshal(data)
	if err != nil {
		return err
	}

	schemaEntity := &repository.SchemaEntity{
		Slug:    slug,
		Content: content,
	}
	if err := s.schemaEntitiesRepo.Upsert(ctx, schemaEntity); err != nil {
		return err
	}

	key := kvstore.EntityDataPrefix + slug
	if err := s.kvstore.SetEntityData(ctx, key, data); err != nil {
		s.logger.Warn("UpdateEntityData: failed to set entity data to kvstore", slog.String("slug", slug))
		return err
	}

	return nil
}

func (s *SchemaEntityService) CollectFullEntityData(ctx context.Context, slug string) (pageData map[string]interface{}, err error) {
	key := kvstore.EntityDataPrefix + slug
	schema, err := s.kvstore.GetEntitySchema(ctx, kvstore.SchemaKeyPages, slug)
	if err != nil {
		s.logger.Warn("CollectEntityData: failed to get entity schema", slog.String("slug", slug))
		return nil, err
	}

	data, err := s.kvstore.GetEntityData(ctx, key)
	if err == nil {
		s.collectSharedData(ctx, schema.Shared, data)
		return data, nil
	}

	schemaEntity, err := s.schemaEntitiesRepo.GetBySlug(ctx, slug)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			emptyData := make(map[string]interface{})
			return emptyData, nil
		}
		return nil, err
	}

	var schemaEntityData map[string]interface{}
	if schemaEntity.Content != nil {
		if err := json.Unmarshal(schemaEntity.Content, &schemaEntityData); err != nil {
			schemaEntityData = make(map[string]interface{})
		}
	} else {
		schemaEntityData = make(map[string]interface{})
	}

	if err := s.kvstore.SetEntityData(ctx, key, schemaEntityData); err != nil {
		s.logger.Warn("GetSchemaEntityData: failed to set schemaEntity data to kvstore", slog.String("slug", slug))
	}

	s.collectSharedData(ctx, schema.Shared, schemaEntityData)

	return schemaEntityData, nil
}

func (s *SchemaEntityService) collectSharedData(ctx context.Context, sharedList []string, data map[string]interface{}) error {
	for _, shared := range sharedList {
		key := kvstore.EntityDataPrefix + shared
		sharedData, err := s.kvstore.GetEntityData(ctx, key)

		if err == nil {
			data[shared] = sharedData
		} else {
			schemaEntity, err := s.schemaEntitiesRepo.GetBySlug(ctx, shared)
			if err != nil {
				if errors.Is(err, repository.ErrNotFound) {
					emptyData := make(map[string]interface{})
					data[shared] = emptyData
				}
			} else {
				data[shared] = schemaEntity
			}
		}
	}
	return nil
}
