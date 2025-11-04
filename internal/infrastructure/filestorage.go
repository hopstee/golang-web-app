package infrastructure

import (
	"fmt"
	"mobile-backend-boilerplate/internal/filestorage/localstorage"
)

func (d *Dependencies) InitFileStorage() error {
	var localStorage *localstorage.LocalStorage

	switch d.Config.FileStorage.Type {
	case "local":
		localStorage = localstorage.NewLocalStorage(
			d.Config.FileStorage.Local.BasePath,
			d.Logger,
		)
	default:
		return fmt.Errorf("unsupported file storage type: %s", d.Config.FileStorage.Type)
	}

	d.FileStorage = localStorage
	return nil
}
