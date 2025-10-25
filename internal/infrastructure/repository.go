package infrastructure

import (
	"errors"
	"log/slog"
	"mobile-backend-boilerplate/internal/repository"
	"mobile-backend-boilerplate/internal/repository/postgres"
)

func (d *Dependencies) InitRepository() error {
	var repo repository.Repository
	var err error

	switch d.Config.Database.Driver {
	case "postgres":
		repo, err = postgres.NewPostgreSQLRepository(d.Config.Database.DataSource, d.Logger)
	default:
		d.Logger.Error("unsupported database driver", slog.String("driver", d.Config.Database.Driver))
		err := errors.New("unsupported database driver")
		return err
	}

	if err != nil {
		d.Logger.Error("failed to init repository", slog.Any("err", err))
		return err
	}

	d.Repository = repo
	return nil
}
