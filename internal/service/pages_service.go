package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"mobile-backend-boilerplate/internal/kvstore"
	"mobile-backend-boilerplate/internal/repository"
)

type PagesService struct {
	pagesRepo repository.PagesRepository
	kvstore   kvstore.KVStore
	logger    *slog.Logger
}

func NewPagesService(pagesRepo repository.PagesRepository, kvstore kvstore.KVStore, logger *slog.Logger) *PagesService {
	return &PagesService{
		pagesRepo: pagesRepo,
		kvstore:   kvstore,
		logger:    logger,
	}
}

func (s *PagesService) GetAllPagesSchemas(ctx context.Context) ([]*kvstore.EntitySchema, error) {
	schemas, err := s.kvstore.GetPages(ctx)
	if err != nil {
		return nil, err
	}

	return schemas, nil
}

func (s *PagesService) GetPageSchema(ctx context.Context, slug string) (*kvstore.EntitySchema, error) {
	schemas, err := s.GetAllPagesSchemas(ctx)
	if err != nil {
		return nil, err
	}

	for _, schema := range schemas {
		if schema.ID == slug {
			return schema, nil
		}
	}

	s.logger.Info("GetPageSchema: schema not found", slog.String("slug", slug))
	return nil, fmt.Errorf("schema not found")
}

func (s *PagesService) GetPageData(ctx context.Context, slug string) (*kvstore.EntitySchema, map[string]interface{}, error) {
	schema, err := s.GetPageSchema(ctx, slug)
	if err != nil {
		return nil, nil, err
	}

	data, err := s.kvstore.GetPageData(ctx, slug)
	if err == nil {
		return schema, data, nil
	}

	page, err := s.pagesRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, nil, err
	}

	var pageData map[string]interface{}
	if page.Content != nil {
		if err := json.Unmarshal(page.Content, &pageData); err != nil {
			pageData = make(map[string]interface{})
		}
	} else {
		pageData = make(map[string]interface{})
	}

	if err := s.kvstore.SetPageData(ctx, slug, pageData); err != nil {
		s.logger.Warn("GetPageData: failed to set page data to kvstore", slog.String("slug", slug))
	}

	return schema, pageData, nil
}

func (s *PagesService) UpdatePageData(ctx context.Context, slug string, data map[string]interface{}) error {
	content, err := json.Marshal(data)
	if err != nil {
		return err
	}

	page := &repository.Page{
		Slug:    slug,
		Content: content,
	}
	if err := s.pagesRepo.Update(ctx, page); err != nil {
		return err
	}

	if err := s.kvstore.SetPageData(ctx, slug, data); err != nil {
		s.logger.Warn("UpdatePageData: failed to set page data to kvstore", slog.String("slug", slug))
		return err
	}

	return nil
}

func (s *PagesService) DeletePageData(ctx context.Context, slug string) error {
	if err := s.pagesRepo.DeleteBySlug(ctx, slug); err != nil {
		s.logger.Warn("DeletePageData: failed to delete page from repository", slog.String("slug", slug))
		return err
	}

	if err := s.kvstore.DeletePageData(ctx, slug); err != nil {
		s.logger.Warn("DeletePageData: failed to delete page data from kvstore", slog.String("slug", slug))
		return err
	}

	return nil
}
