package sqlite

import (
	"database/sql"
	"fmt"
	"log/slog"
	"mobile-backend-boilerplate/internal/repository"

	_ "github.com/mattn/go-sqlite3"
)

type SQLiteRepository struct {
	DB     *sql.DB
	Logger *slog.Logger

	auth    repository.AuthRepository
	user    repository.UserRepository
	admin   repository.AdminRepository
	post    repository.PostRepository
	request repository.RequestRepository
}

func NewSQLiteRepository(dsn string, logger *slog.Logger) (*SQLiteRepository, error) {
	var err error

	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite db: %w", err)
	}

	repository := &SQLiteRepository{
		DB:     db,
		Logger: logger,
	}

	repository.user = NewUserRepo(db, logger)
	repository.auth = NewAuthRepo(db, logger)
	repository.admin = NewAdminRepo(db, logger)
	repository.post = NewPostRepo(db, logger)
	repository.request = NewRequestRepo(db, logger)

	if err = repository.Migrate(Up, false); err != nil {
		logger.Error("migration failed", slog.Any("err", err))
		return nil, fmt.Errorf("migration failed: %w", err)
	}

	return repository, nil
}

func (r *SQLiteRepository) Database() *sql.DB {
	return r.DB
}

func (r *SQLiteRepository) Auth() repository.AuthRepository {
	return r.auth
}

func (r *SQLiteRepository) User() repository.UserRepository {
	return r.user
}

func (r *SQLiteRepository) Admin() repository.AdminRepository {
	return r.admin
}

func (r *SQLiteRepository) Post() repository.PostRepository {
	return r.post
}

func (r *SQLiteRepository) Request() repository.RequestRepository {
	return r.request
}
