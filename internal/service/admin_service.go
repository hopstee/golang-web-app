package service

import (
	"log/slog"
	"mobile-backend-boilerplate/internal/repository"
)

type AdminService struct {
	adminRepo repository.AdminRepository
	logger    *slog.Logger
}

func NewAdminService(adminRepo repository.AdminRepository, logger *slog.Logger) *AdminService {
	return &AdminService{
		adminRepo: adminRepo,
		logger:    logger,
	}
}

func (s *AdminService) GetAdminById(id int64) (repository.Admin, error) {
	s.logger.Debug("GetUserById attempt")

	user, err := s.adminRepo.GetByID(id)
	if err != nil {
		s.logger.Error("GetUserById failed: failed to get user by id", slog.Int64("user_id", id), slog.Any("err", err))
		return repository.Admin{}, err
	}

	s.logger.Info("GetUserById successful", slog.Int64("user_id", id))
	return user, nil
}
