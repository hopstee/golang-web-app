package infrastructure

import (
	"errors"
	"log/slog"
	"mobile-backend-boilerplate/internal/repository/sqlite"
)

func (d *Dependencies) InitRepos() error {
	switch d.Config.Database.Driver {
	case "sqlite":
		d.AuthRepo = sqlite.NewAuthRepo(d.DB, d.Logger)
		d.UserRepo = sqlite.NewUserRepo(d.DB, d.Logger)
		d.AdminRepo = sqlite.NewAdminRepo(d.DB, d.Logger)
		d.RequestRepo = sqlite.NewRequestRepo(d.DB, d.Logger)
		d.PostRepo = sqlite.NewPostRepo(d.DB, d.Logger)
	default:
		d.Logger.Error("unsupported database driver", slog.String("driver", d.Config.Database.Driver))
		err := errors.New("unsupported database driver")
		return err
	}

	return nil
}
