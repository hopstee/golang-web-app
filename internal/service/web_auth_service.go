package service

import (
	"errors"
	"log/slog"
	"mobile-backend-boilerplate/internal/repository"
	"mobile-backend-boilerplate/pkg/helper/jwt"
	"mobile-backend-boilerplate/pkg/helper/password"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
)

type AdminAuthService struct {
	adminRepo repository.AdminRepository
	jwtUtil   *jwt.JWTUtil
	logger    *slog.Logger
}

func NewAdminAuthService(
	adminRepo repository.AdminRepository,
	jwtSecret []byte,
	logger *slog.Logger,
) *AdminAuthService {
	return &AdminAuthService{
		adminRepo: adminRepo,
		jwtUtil:   jwt.NewJWTUtil(jwtSecret),
		logger:    logger,
	}
}

func (s *AdminAuthService) Login(username, passwordStr string) (string, error) {
	s.logger.Info("admin loging attempt", slog.String("username", username))
	authData, err := s.adminRepo.GetAuthData(username)
	if err != nil {
		s.logger.Warn("admin login failed: admin not found", slog.String("username", username), slog.Any("err", err))
		return "", ErrUserNotFound
	}

	if !password.CheckPasswordHash(passwordStr, authData.Password) {
		s.logger.Warn("admin login failed: invalid password", slog.String("username", username))
		return "", ErrInvalidPassword
	}

	accessToken, _, err := s.jwtUtil.CreateAccessToken(authData.ID, authData.TokenVersion, true)
	if err != nil {
		s.logger.Error("admin login failed: failed to create access token", slog.Int64("user_id", authData.ID), slog.Any("err", err))
		return "", err
	}

	s.logger.Info("admin loging successfull", slog.Int64("user_id", authData.ID))
	return accessToken, nil
}

func (s *AdminAuthService) Me(accessToken string) (repository.Admin, error) {
	s.logger.Debug("admin me request")

	claims, err := s.jwtUtil.ParseToken(accessToken)
	if err != nil {
		s.logger.Warn("admin me request failed: invalid access token")
		return repository.Admin{}, errors.New("invalid access token")
	}

	userID, ok := claims["sub"].(float64)
	if !ok {
		s.logger.Error("admin me request failed: invalid claims in access token")
		return repository.Admin{}, errors.New("invalid claims")
	}

	admin, err := s.adminRepo.GetByID(int64(userID))
	if err != nil {
		s.logger.Error("admin me request failed: failed to get user by id", slog.Int64("user_id", int64(userID)), slog.Any("err", err))
		return repository.Admin{}, err
	}

	s.logger.Info("admin me request successful", slog.Int64("user_id", int64(userID)))
	return admin, nil
}
