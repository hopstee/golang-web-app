package postgres

import (
	"database/sql"
	"log/slog"
	"mobile-backend-boilerplate/internal/repository"
	"time"
)

type adminRepo struct {
	DB     *sql.DB
	Logger *slog.Logger
}

func NewAdminRepo(db *sql.DB, logger *slog.Logger) repository.AdminRepository {
	return &adminRepo{
		DB:     db,
		Logger: logger,
	}
}

func (r *adminRepo) GetByID(id int64) (repository.Admin, error) {
	return r.getAdmin(`
		SELECT id, username, role, created_at, updated_at
		FROM admins
		WHERE id = $1
	`, id)
}

func (r *adminRepo) GetByUsername(username string) (repository.Admin, error) {
	return r.getAdmin(`
		SELECT id, username, role, created_at, updated_at
		FROM admins
		WHERE username = $1
	`, username)
}

func (r *adminRepo) GetByEmail(email string) (repository.Admin, error) {
	return r.getAdmin(`
		SELECT id, username, role, created_at, updated_at
		FROM admins
		WHERE email = $1
	`, email)
}

func (r *adminRepo) getAdmin(query string, args ...interface{}) (repository.Admin, error) {
	var a repository.Admin
	row := r.DB.QueryRow(query, args...)
	err := row.Scan(&a.ID, &a.Username, &a.Role, &a.CreatedAt, &a.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			r.Logger.Warn("admin not found", slog.Any("args", args))
			return a, repository.ErrNotFound
		}
		r.Logger.Error("failed to query admin", slog.Any("args", args), slog.Any("err", err))
		return a, err
	}
	return a, nil
}

func (r *adminRepo) GetAuthData(username string) (repository.AuthData, error) {
	row := r.DB.QueryRow(`
		SELECT id, password
		FROM admins
		WHERE username = $1
	`, username)

	var authData repository.AuthData

	err := row.Scan(&authData.ID, &authData.Password)
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

func (r *adminRepo) Create(admin repository.Admin) (int64, error) {
	var id int64
	err := r.DB.QueryRow(`
		INSERT INTO admins (username, password, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, admin.Username, admin.Password, admin.Role, time.Now(), time.Now()).Scan(&id)

	if err != nil {
		r.Logger.Error("failed to create admin", slog.String("username", admin.Username), slog.Any("err", err))
		return 0, err
	}

	r.Logger.Debug("admin created", slog.Int64("id", id), slog.String("username", admin.Username))
	return id, nil
}

func (r *adminRepo) Update(admin repository.Admin) error {
	_, err := r.DB.Exec(`
		UPDATE admins
		SET username = $1, password = $2, role = $3, updated_at = $4
		WHERE id = $5
	`, admin.Username, admin.Password, admin.Role, time.Now(), admin.ID)
	if err != nil {
		r.Logger.Error("failed to update admin", slog.String("username", admin.Username), slog.Any("err", err))
		return err
	}

	r.Logger.Debug("admin updated", slog.Int64("id", admin.ID), slog.String("username", admin.Username))
	return nil
}

func (r *adminRepo) Delete(id int64) error {
	_, err := r.DB.Exec(`DELETE FROM admins WHERE id = $1`, id)
	if err != nil {
		r.Logger.Error("failed to delete admin", slog.Int64("id", id), slog.Any("err", err))
		return err
	}

	r.Logger.Debug("admin deleted", slog.Int64("id", id))
	return nil
}
