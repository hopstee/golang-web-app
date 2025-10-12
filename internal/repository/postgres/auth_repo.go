package postgres

import (
	"database/sql"
	"log/slog"
	"mobile-backend-boilerplate/internal/repository"
)

type authRepo struct {
	DB     *sql.DB
	Logger *slog.Logger
}

func NewAuthRepo(db *sql.DB, logger *slog.Logger) repository.AuthRepository {
	return &authRepo{
		DB:     db,
		Logger: logger,
	}
}

func (r *authRepo) StoreRefreshToken(rt repository.RefreshToken) error {
	_, err := r.DB.Exec(`
		INSERT INTO refresh_tokens (token, user_id, expires_at, is_revoked, created_at, device_id)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, rt.Token, rt.UserID, rt.ExpiresAt, rt.IsRevoked, rt.CreatedAt, rt.DeviceID)

	if err != nil {
		r.Logger.Error("failed to store token", slog.Any("err", err), slog.Int64("user", rt.UserID))
		return err
	}

	r.Logger.Debug("token stored", slog.Int64("user", rt.UserID), slog.String("device", rt.DeviceID))
	return nil
}

func (r *authRepo) GetRefreshToken(token string) (repository.RefreshToken, error) {
	row := r.DB.QueryRow(`
		SELECT token, user_id, is_revoked, expires_at, device_id
		FROM refresh_tokens
		WHERE token = $1
	`, token)

	var rt repository.RefreshToken
	if err := row.Scan(&rt.Token, &rt.UserID, &rt.IsRevoked, &rt.ExpiresAt, &rt.DeviceID); err != nil {
		if err == sql.ErrNoRows {
			r.Logger.Warn("refresh token not found", slog.String("token", token))
			return repository.RefreshToken{}, repository.ErrNotFound
		}
		r.Logger.Error("failed to get token", slog.String("token", token), slog.Any("err", err))
		return repository.RefreshToken{}, err
	}

	return rt, nil
}

func (r *authRepo) InvalidateRefreshToken(token string) error {
	_, err := r.DB.Exec(`
		UPDATE refresh_tokens
		SET is_revoked = TRUE
		WHERE token = $1
	`, token)
	if err != nil {
		r.Logger.Error("failed to invalidate token", slog.String("token", token), slog.Any("err", err))
		return err
	}

	r.Logger.Debug("token invalidated", slog.String("token", token))
	return nil
}
