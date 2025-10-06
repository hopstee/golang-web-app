package service

import (
	"errors"
	"fmt"
	"log/slog"
	"mobile-backend-boilerplate/internal/repository"
	"mobile-backend-boilerplate/pkg/helper/jwt"
	"mobile-backend-boilerplate/pkg/helper/password"
)

type WebAuthService struct {
	adminRepo repository.AdminRepository
	jwtUtil   *jwt.JWTUtil
	logger    *slog.Logger
}

func NewWebAuthService(
	adminRepo repository.AdminRepository,
	jwtSecret []byte,
	logger *slog.Logger,
) *WebAuthService {
	return &WebAuthService{
		adminRepo: adminRepo,
		jwtUtil:   jwt.NewJWTUtil(jwtSecret),
		logger:    logger,
	}
}

func (s *WebAuthService) Login(username, passwordStr string) (string, error) {
	s.logger.Info("admin loging attempt", slog.String("username", username))
	authData, err := s.adminRepo.GetAuthData(username)
	if err != nil {
		s.logger.Warn("admin login failed: admin not found", slog.String("username", username), slog.Any("err", err))
		return "", err
	}

	fmt.Printf("is pass correct: %v", password.CheckPasswordHash(passwordStr, authData.Password))
	if !password.CheckPasswordHash(passwordStr, authData.Password) {
		s.logger.Warn("admin login failed: invalid password", slog.String("username", username))
		return "", errors.New("invalid password")
	}

	accessToken, _, err := s.jwtUtil.CreateAccessToken(authData.ID, authData.TokenVersion, true)
	if err != nil {
		s.logger.Error("admin login failed: failed to create access token", slog.Int64("user_id", authData.ID), slog.Any("err", err))
		return "", err
	}

	s.logger.Info("admin loging successfull", slog.Int64("user_id", authData.ID))
	return accessToken, nil
}

func (s *WebAuthService) Me(accessToken string) (repository.Admin, error) {
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
