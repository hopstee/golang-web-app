package postgres

import (
	"database/sql"
	"fmt"
	"log/slog"
	"mobile-backend-boilerplate/internal/repository"

	_ "github.com/lib/pq"
)

type PostgreSQLRepository struct {
	DB     *sql.DB
	Logger *slog.Logger

	auth    repository.AuthRepository
	user    repository.UserRepository
	admin   repository.AdminRepository
	post    repository.PostRepository
	request repository.RequestRepository
}

func NewPostgreSQLRepository(dsn string, logger *slog.Logger) (*PostgreSQLRepository, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres db: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	repository := &PostgreSQLRepository{
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

	logger.Info("connected to PostgreSQL database")

	return repository, nil
}

func (r *PostgreSQLRepository) Database() *sql.DB {
	return r.DB
}

func (r *PostgreSQLRepository) Auth() repository.AuthRepository {
	return r.auth
}

func (r *PostgreSQLRepository) User() repository.UserRepository {
	return r.user
}

func (r *PostgreSQLRepository) Admin() repository.AdminRepository {
	return r.admin
}

func (r *PostgreSQLRepository) Post() repository.PostRepository {
	return r.post
}

func (r *PostgreSQLRepository) Request() repository.RequestRepository {
	return r.request
}
