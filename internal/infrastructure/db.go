package infrastructure

import (
	"database/sql"
	"errors"
	"log/slog"
	"mobile-backend-boilerplate/internal/config"
	"mobile-backend-boilerplate/internal/repository/sqlite"
)

func InitDB(cfg *config.Config, logger *slog.Logger) (*sql.DB, error) {
	var db *sql.DB

	switch cfg.Database.Driver {
	case "sqlite":
		repo, err := sqlite.NewSQLiteRepository(cfg.Database.DataSource, logger)
		if err != nil {
			logger.Error("failed to init sqlite repository", slog.Any("err", err))
			return nil, err
		}

		db = repo.DB

		logger.Info("SQLite repository initialized", slog.String("DSN", cfg.Database.DataSource))
	default:
		logger.Error("unsupported database driver", slog.String("driver", cfg.Database.Driver))
		err := errors.New("unsupported database driver")
		return nil, err
	}

	return db, nil
}
