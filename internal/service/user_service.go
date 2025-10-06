package service

import (
	"log/slog"
	"mobile-backend-boilerplate/internal/repository"
)

type UserService struct {
	userRepo repository.UserRepository
	logger   *slog.Logger
}

func NewUserService(userRepo repository.UserRepository, logger *slog.Logger) *UserService {
	return &UserService{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (s *UserService) GetUserById(id int64) (repository.User, error) {
	s.logger.Debug("GetUserById attempt")

	user, err := s.userRepo.GetByID(id)
	if err != nil {
		s.logger.Error("GetUserById failed: failed to get user by id", slog.Int64("user_id", id), slog.Any("err", err))
		return repository.User{}, err
	}

	s.logger.Info("GetUserById successful", slog.Int64("user_id", id))
	return user, nil
}
