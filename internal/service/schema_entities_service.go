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

func (s *SchemaEntityService) CollectFullEntityData(ctx context.Context, slug string) (pageData map[string]interface{}, layoutName string, err error) {
	schema, err := s.kvstore.GetEntitySchema(ctx, kvstore.SchemaKeyPages, slug)
	if err != nil {
		s.logger.Warn("CollectEntityData: failed to get entity schema", slog.String("slug", slug))
		return nil, layoutName, err
	}

	layoutName = schema.Layout
	key := kvstore.EntityDataPrefix + slug

	data, err := s.kvstore.GetEntityData(ctx, key)
	if err != nil {
		schemaEntity, repoErr := s.schemaEntitiesRepo.GetBySlug(ctx, slug)
		if repoErr != nil {
			if errors.Is(repoErr, repository.ErrNotFound) {
				emptyData := make(map[string]interface{})
				return emptyData, layoutName, nil
			}
			return nil, layoutName, repoErr
		}

		if schemaEntity.Content != nil {
			if err := json.Unmarshal(schemaEntity.Content, &data); err != nil {
				data = make(map[string]interface{})
			}
		} else {
			data = make(map[string]interface{})
		}

		if cacheErr := s.kvstore.SetEntityData(ctx, key, data); cacheErr != nil {
			s.logger.Warn("CollectFullEntityData: failed to set schemaEntity data to kvstore", slog.String("slug", slug), slog.Any("err", cacheErr))
		}
	}

	contentRaw, ok := data[repository.MainContentKey]
	if !ok {
		s.logger.Warn("CollectFullEntityData: entity has no content", slog.String("slug", slug))
		data[repository.MainContentKey] = map[string]interface{}{}
		contentRaw = data[repository.MainContentKey]
	}

	content, ok := contentRaw.(map[string]interface{})
	if !ok {
		content = map[string]interface{}{}
		data[repository.MainContentKey] = content
	}

	if err := s.collectSharedData(ctx, schema, content); err != nil {
		s.logger.Warn("CollectFullEntityData: failed to collect shared data", slog.String("slug", slug), slog.Any("err", err))
		return nil, layoutName, err
	}

	return data, layoutName, nil
}

func (s *SchemaEntityService) collectSharedData(ctx context.Context, node *kvstore.EntitySchema, data map[string]interface{}) error {
	if node == nil || data == nil {
		return nil
	}

	for _, shared := range node.Shared {
		key := kvstore.EntityDataPrefix + shared
		var sharedMap map[string]interface{}

		if sharedData, err := s.kvstore.GetEntityData(ctx, key); err == nil {
			if sharedDataContent, ok := sharedData[repository.MainContentKey]; ok {
				data[shared] = sharedDataContent
				continue
			}
			data[shared] = sharedData
			continue
		}

		schemaEntity, err := s.schemaEntitiesRepo.GetBySlug(ctx, shared)
		if err != nil {
			if errors.Is(err, repository.ErrNotFound) {
				data[shared] = make(map[string]interface{})
				continue
			}
			s.logger.Warn("collectSharedData: failed to get schema entity", slog.String("slug", shared), slog.Any("err", err))
			continue
		}

		if schemaEntity.Content != nil {
			if err := json.Unmarshal(schemaEntity.Content, &sharedMap); err != nil {
				s.logger.Warn("collectSharedData: failed to unmarshal content", slog.String("slug", shared), slog.Any("err", err))
				sharedMap = map[string]interface{}{}
			}
		} else {
			sharedMap = map[string]interface{}{}
		}

		if sharedContent, ok := sharedMap[repository.MainContentKey]; ok {
			data[shared] = sharedContent
		} else {
			data[shared] = sharedMap
		}
	}

	for _, child := range node.Children {
		childDataRaw, ok := data[child.ID]
		if !ok {
			continue
		}

		childData, ok := childDataRaw.(map[string]interface{})
		if !ok {
			continue
		}

		if err := s.collectSharedData(ctx, child, childData); err != nil {
			s.logger.Warn("collectSharedData: failed for child", slog.String("child", child.ID), slog.Any("err", err))
		}
	}

	return nil
}
