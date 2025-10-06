package service

import (
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"mobile-backend-boilerplate/internal/repository"
	"os"
	"path/filepath"
	"time"
)

type PostService struct {
	postRepo repository.PostRepository
	logger   *slog.Logger
}

func NewPostService(postRepo repository.PostRepository, logger *slog.Logger) *PostService {
	return &PostService{
		postRepo: postRepo,
		logger:   logger,
	}
}

func (s *PostService) GetAllPublicPosts() ([]repository.Post, error) {
	s.logger.Debug("GetAllPublicPosts attempt")

	posts, err := s.postRepo.GetAllPublic()
	if err != nil {
		s.logger.Error("GetAllPublicPosts failed: failed to get public posts", slog.Any("err", err))
		return nil, err
	}

	s.logger.Info("GetAllPublicPosts successful")
	return posts, nil
}

func (s *PostService) GetPublicPost(id int64) (repository.Post, error) {
	s.logger.Debug("GetPublicPost attempt")

	post, err := s.postRepo.GetPublicByID(id)
	if err != nil {
		s.logger.Error("GetPublicPost failed: failed to get public post", slog.Int64("post", id), slog.Any("err", err))
		return repository.Post{}, err
	}

	s.logger.Info("GetPublicPost successful")
	return post, nil
}

func (s *PostService) GetAllPosts() ([]repository.Post, error) {
	s.logger.Debug("GetAllPosts attempt")

	posts, err := s.postRepo.GetAll()
	if err != nil {
		s.logger.Error("GetAllPosts failed: failed to get all posts", slog.Any("err", err))
		return nil, err
	}

	s.logger.Info("GetAllPosts successful")
	return posts, nil
}

func (s *PostService) GetPost(id int64) (repository.Post, error) {
	s.logger.Debug("GetPost attempt")

	post, err := s.postRepo.GetByID(id)
	if err != nil {
		s.logger.Error("GetPost failed: failed to get post by id", slog.Int64("post", id), slog.Any("err", err))
		return repository.Post{}, err
	}

	s.logger.Info("GetPost successful")
	return post, nil
}

func (s *PostService) CreatePost(post repository.Post, file *multipart.FileHeader) (int64, error) {
	s.logger.Debug("CreatePost attempt")

	if file != nil {
		url, err := s.saveHeroImage(file)
		if err != nil {
			s.logger.Error("failed to save hero image", slog.Any("err", err))
			return 0, err
		}
		post.HeroImgURL = url
	}

	id, err := s.postRepo.Create(post)
	if err != nil {
		s.logger.Error("CreatePost failed: failed to create post", slog.Int64("post", id), slog.Any("err", err))
		return 0, err
	}

	s.logger.Info("CreatePost successful")
	return id, nil
}

func (s *PostService) UpdatePost(post repository.Post, file *multipart.FileHeader) error {
	s.logger.Debug("UpdatePost attempt")

	if file != nil {
		url, err := s.saveHeroImage(file)
		if err != nil {
			s.logger.Error("failed to save hero image", slog.Any("err", err))
			return err
		}
		post.HeroImgURL = url
	}

	err := s.postRepo.Update(post)
	if err != nil {
		s.logger.Error("UpdatePost failed: failed to create post", slog.Int64("post", post.ID), slog.Any("err", err))
		return err
	}

	s.logger.Info("UpdatePost successful")
	return nil
}

func (s *PostService) DeletePost(id int64) error {
	s.logger.Debug("DeletePost attempt")

	err := s.postRepo.Delete(id)
	if err != nil {
		s.logger.Error("DeletePost failed: failed to create post", slog.Int64("post", id), slog.Any("err", err))
		return err
	}

	s.logger.Info("DeletePost successful")
	return nil
}

func (s *PostService) saveHeroImage(file *multipart.FileHeader) (string, error) {
	if err := os.MkdirAll("static/uploads", os.ModePerm); err != nil {
		return "", err
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
	filepath := filepath.Join("static/uploads", filename)

	dst, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}

	return "/" + filepath, nil
}
