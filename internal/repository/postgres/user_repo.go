package postgres

import (
	"database/sql"
	"log/slog"
	"mobile-backend-boilerplate/internal/repository"
	"time"
)

type userRepo struct {
	DB     *sql.DB
	Logger *slog.Logger
}

func NewUserRepo(db *sql.DB, logger *slog.Logger) repository.UserRepository {
	return &userRepo{
		DB:     db,
		Logger: logger,
	}
}

func (r *userRepo) GetByID(id int64) (repository.User, error) {
	return r.getUser(`
		SELECT id, username, email, created_at, updated_at
		FROM users
		WHERE id = $1
	`, id)
}

func (r *userRepo) GetByUsername(username string) (repository.User, error) {
	return r.getUser(`
		SELECT id, username, email, created_at, updated_at
		FROM users
		WHERE username = $1
	`, username)
}

func (r *userRepo) GetByEmail(email string) (repository.User, error) {
	return r.getUser(`
		SELECT id, username, email, created_at, updated_at
		FROM users
		WHERE email = $1
	`, email)
}

func (r *userRepo) getUser(query string, args ...interface{}) (repository.User, error) {
	var u repository.User
	row := r.DB.QueryRow(query, args...)
	err := row.Scan(&u.ID, &u.Username, &u.Email, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			r.Logger.Warn("user not found", slog.Any("args", args))
			return u, repository.ErrNotFound
		}
		r.Logger.Error("failed to query user", slog.Any("args", args), slog.Any("err", err))
		return u, err
	}
	return u, nil
}

func (r *userRepo) GetAuthData(username string) (repository.AuthData, error) {
	row := r.DB.QueryRow(`
		SELECT id, password, token_version
		FROM users
		WHERE username = $1
	`, username)

	var authData repository.AuthData

	err := row.Scan(&authData.ID, &authData.Password, &authData.TokenVersion)
	if err != nil {
		if err == sql.ErrNoRows {
			r.Logger.Warn("auth data not found", slog.String("username", username))
			return authData, repository.ErrNotFound
		}
		r.Logger.Error("failed to get auth data", slog.String("username", username), slog.Any("err", err))
		return authData, err
	}

	r.Logger.Debug("retrieved auth data", slog.Int64("id", authData.ID), slog.String("username", username))
	return authData, nil
}

func (r *userRepo) GetTokenVersion(id int64) (int, error) {
	row := r.DB.QueryRow(`
		SELECT token_version
		FROM users
		WHERE id = $1
	`, id)

	var tokenVersion int

	if err := row.Scan(&tokenVersion); err != nil {
		if err == sql.ErrNoRows {
			r.Logger.Warn("token version not found", slog.Int64("id", id))
			return tokenVersion, repository.ErrNotFound
		}
		r.Logger.Error("failed to get token version", slog.Int64("id", id), slog.Any("err", err))
		return tokenVersion, err
	}

	r.Logger.Debug("retrieved token version", slog.Int64("id", id), slog.Int("version", tokenVersion))
	return tokenVersion, nil
}

func (r *userRepo) Create(user repository.User) (int64, error) {
	var id int64
	err := r.DB.QueryRow(`
		INSERT INTO users(username, password, email, token_version, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, user.Username, user.Password, user.Email, 1, time.Now(), time.Now()).Scan(&id)
	if err != nil {
		r.Logger.Error("failed to create user", slog.String("username", user.Username), slog.Any("err", err))
		return 0, err
	}

	r.Logger.Debug("user created", slog.Int64("id", id), slog.String("username", user.Username))
	return id, nil
}

func (r *userRepo) Update(user repository.User) error {
	_, err := r.DB.Exec(`
		UPDATE users
		SET username = $1, password = $2, email = $3, updated_at = $4
		WHERE id = $5
	`, user.Username, user.Password, user.Email, time.Now(), user.ID)
	if err != nil {
		r.Logger.Error("failed to update user", slog.String("username", user.Username), slog.Any("err", err))
		return err
	}

	r.Logger.Debug("user updated", slog.Int64("id", user.ID), slog.String("username", user.Username))
	return nil
}

func (r *userRepo) Delete(id int64) error {
	_, err := r.DB.Exec(`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		r.Logger.Error("failed to delete user", slog.Int64("id", id), slog.Any("err", err))
		return err
	}

	r.Logger.Debug("user deleted", slog.Int64("id", id))
	return nil
}

func (r *userRepo) IncrementTokenVersion(id int64) error {
	_, err := r.DB.Exec(`
		UPDATE users
		SET token_version = token_version + 1
		WHERE id = $1
	`, id)
	if err != nil {
		r.Logger.Error("failed to increment token version", slog.Int64("id", id), slog.Any("err", err))
		return err
	}

	r.Logger.Debug("token version incremented", slog.Int64("id", id))
	return nil
}
