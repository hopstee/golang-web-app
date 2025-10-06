package service

import (
	"errors"
	"log/slog"
	"mobile-backend-boilerplate/internal/repository"
	"mobile-backend-boilerplate/pkg/helper/jwt"
	"mobile-backend-boilerplate/pkg/helper/password"
	"time"
)

type MobileAuthService struct {
	authRepo repository.AuthRepository
	userRepo repository.UserRepository
	jwtUtil  *jwt.JWTUtil
	logger   *slog.Logger
}

func NewMobileAuthService(
	ar repository.AuthRepository,
	ur repository.UserRepository,
	jwtSecret []byte,
	logger *slog.Logger,
) *MobileAuthService {
	return &MobileAuthService{
		authRepo: ar,
		userRepo: ur,
		jwtUtil:  jwt.NewJWTUtil(jwtSecret),
		logger:   logger,
	}
}

func (s *MobileAuthService) Login(username, passwordStr, deviceID string) (repository.TokensPair, error) {
	s.logger.Info("loging attempt", slog.String("username", username), slog.String("device_id", deviceID))
	authData, err := s.userRepo.GetAuthData(username)
	if err != nil {
		s.logger.Warn("login failed: user not found", slog.String("username", username), slog.Any("err", err))
		return repository.TokensPair{}, err
	}

	if !password.CheckPasswordHash(passwordStr, authData.Password) {
		s.logger.Warn("login failed: invalid password", slog.String("username", username))
		return repository.TokensPair{}, errors.New("invalid password")
	}

	accessToken, exp, err := s.jwtUtil.CreateAccessToken(authData.ID, authData.TokenVersion, false)
	if err != nil {
		s.logger.Error("login failed: failed to create access token", slog.Int64("user_id", authData.ID), slog.Any("err", err))
		return repository.TokensPair{}, err
	}

	rt := repository.RefreshToken{
		Token:     jwt.GenerateRandomString(),
		UserID:    authData.ID,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		IsRevoked: false,
		CreatedAt: time.Now(),
		DeviceID:  deviceID,
	}

	if err := s.authRepo.StoreRefreshToken(rt); err != nil {
		s.logger.Error("login failed: failed to store refresh token", slog.Int64("user_id", authData.ID), slog.Any("err", err))
		return repository.TokensPair{}, err
	}

	s.logger.Info("loging successfull", slog.Int64("user_id", authData.ID), slog.String("device_id", deviceID))
	return repository.TokensPair{
		AccessToken:  accessToken,
		RefreshToken: rt.Token,
		ExpiresAt:    exp,
	}, nil
}

func (s *MobileAuthService) Register(username, passwordStr, email, deviceID string) (repository.TokensPair, error) {
	s.logger.Info("registration attempt", slog.String("username", username), slog.String("email", email), slog.String("device_id", deviceID))

	var err error

	_, err = s.userRepo.GetByUsername(username)
	if err == nil {
		s.logger.Warn("registration failed: username already exists", slog.String("username", username), slog.Any("err", err))
		return repository.TokensPair{}, errors.New("username already exists")
	}

	_, err = s.userRepo.GetByEmail(username)
	if err == nil {
		s.logger.Warn("registration failed: email already exists", slog.String("email", email), slog.Any("err", err))
		return repository.TokensPair{}, errors.New("email already exists")
	}

	hashedPassword, err := password.HashPassword(passwordStr)
	if err != nil {
		s.logger.Error("registration failed: failed to hash password", slog.String("username", username), slog.Any("err", err))
		return repository.TokensPair{}, err
	}

	tokenVersion := 1
	user := repository.User{
		Username:     username,
		Password:     hashedPassword,
		Email:        email,
		TokenVersion: tokenVersion,
	}

	userID, err := s.userRepo.Create(user)
	if err != nil {
		s.logger.Error("registration failed: failed to create user", slog.String("username", username), slog.Any("err", err))
		return repository.TokensPair{}, err
	}

	accessToken, exp, err := s.jwtUtil.CreateAccessToken(userID, tokenVersion, false)
	if err != nil {
		s.logger.Error("registration failed: failed to create access token", slog.Int64("user_id", userID), slog.Any("err", err))
		return repository.TokensPair{}, err
	}

	rt := repository.RefreshToken{
		Token:     jwt.GenerateRandomString(),
		UserID:    userID,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		IsRevoked: false,
		CreatedAt: time.Now(),
		DeviceID:  deviceID,
	}

	if err := s.authRepo.StoreRefreshToken(rt); err != nil {
		s.logger.Error("registration failed: failed to store refresh token", slog.Int64("user_id", userID), slog.Any("err", err))
		return repository.TokensPair{}, err
	}

	s.logger.Info("registration successfull", slog.Int64("user_id", userID), slog.String("username", username))
	return repository.TokensPair{
		AccessToken:  accessToken,
		RefreshToken: rt.Token,
		ExpiresAt:    exp,
	}, nil
}

func (s *MobileAuthService) Refresh(refreshToken, deviceID string) (repository.TokensPair, error) {
	s.logger.Info("refresh attempt", slog.String("device_id", deviceID))

	rt, err := s.authRepo.GetRefreshToken(refreshToken)
	if err != nil {
		s.logger.Warn("refresh failed: token not found", slog.String("device_id", deviceID))
		return repository.TokensPair{}, errors.New("refresh token not found")
	}

	if rt.IsRevoked || rt.ExpiresAt.Before(time.Now()) {
		s.logger.Warn("refresh failed: token invalid", slog.String("device_id", deviceID), slog.Int64("user_id", rt.UserID))
		return repository.TokensPair{}, errors.New("refresh token invalid")
	}

	if err := s.authRepo.InvalidateRefreshToken(refreshToken); err != nil {
		s.logger.Error("refresh failed: failed to invalidate old refresh token", slog.Int64("user_id", rt.UserID), slog.Any("err", err))
		return repository.TokensPair{}, err
	}

	if err := s.userRepo.IncrementTokenVersion(rt.UserID); err != nil {
		s.logger.Error("refresh failed: failed to increment access token version", slog.Int64("user_id", rt.UserID), slog.Any("err", err))
		return repository.TokensPair{}, err
	}

	tokenVersion, err := s.userRepo.GetTokenVersion(rt.UserID)
	if err != nil {
		s.logger.Error("refresh failed: failed to get access token version", slog.Int64("user_id", rt.UserID), slog.Any("err", err))
		return repository.TokensPair{}, err
	}

	accessToken, exp, err := s.jwtUtil.CreateAccessToken(rt.UserID, tokenVersion, false)
	if err != nil {
		s.logger.Error("refresh failed: failed to create access token", slog.Int64("user_id", rt.UserID), slog.Any("err", err))
		return repository.TokensPair{}, err
	}

	newRt := repository.RefreshToken{
		Token:     jwt.GenerateRandomString(),
		UserID:    rt.UserID,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		IsRevoked: false,
		CreatedAt: time.Now(),
		DeviceID:  deviceID,
	}

	if err := s.authRepo.StoreRefreshToken(newRt); err != nil {
		s.logger.Error("refresh failed: failed to get store new refresh token", slog.Int64("user_id", rt.UserID), slog.Any("err", err))
		return repository.TokensPair{}, err
	}

	s.logger.Info("refresh successfull", slog.Int64("user_id", rt.UserID), slog.String("device_id", deviceID))
	return repository.TokensPair{
		AccessToken:  accessToken,
		RefreshToken: newRt.Token,
		ExpiresAt:    exp,
	}, nil
}

func (s *MobileAuthService) Logout(refreshToken string) error {
	s.logger.Info("logout attempt")

	rt, err := s.authRepo.GetRefreshToken(refreshToken)
	if err != nil {
		s.logger.Warn("logout failed: invalid refresh token")
		return errors.New("invalid refresh token")
	}

	if err := s.authRepo.InvalidateRefreshToken(refreshToken); err != nil {
		s.logger.Error("logout failed: failed to invalidate refresh token", slog.Int64("user_id", rt.UserID), slog.Any("err", err))
		return err
	}

	if err := s.userRepo.IncrementTokenVersion(rt.UserID); err != nil {
		s.logger.Error("logout failed: failed to increment access token version", slog.Int64("user_id", rt.UserID), slog.Any("err", err))
		return err
	}

	s.logger.Info("logout successful", slog.Int64("user_id", rt.UserID))
	return nil
}

func (s *MobileAuthService) Me(accessToken string) (repository.User, error) {
	s.logger.Debug("me request")

	claims, err := s.jwtUtil.ParseToken(accessToken)
	if err != nil {
		s.logger.Warn("me request failed: invalid access token")
		return repository.User{}, errors.New("invalid access token")
	}

	userID, ok := claims["sub"].(float64)
	if !ok {
		s.logger.Error("me request failed: invalid claims in access token")
		return repository.User{}, errors.New("invalid claims")
	}

	user, err := s.userRepo.GetByID(int64(userID))
	if err != nil {
		s.logger.Error("me request failed: failed to get user by id", slog.Int64("user_id", int64(userID)), slog.Any("err", err))
		return repository.User{}, err
	}

	s.logger.Info("me request successful", slog.Int64("user_id", int64(userID)))
	return user, nil
}
