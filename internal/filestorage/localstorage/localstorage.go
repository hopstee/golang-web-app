package localstorage

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
)

type LocalStorage struct {
	basePath string
	logger   *slog.Logger
}

func NewLocalStorage(basePath string, logger *slog.Logger) *LocalStorage {
	return &LocalStorage{
		basePath: basePath,
		logger:   logger,
	}
}

func (ls *LocalStorage) SaveFile(name string, data []byte) (path string, err error) {
	path = fmt.Sprintf("%s/%s", ls.basePath, name)

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		ls.logger.Error("Failed to open file", slog.String("path", path), slog.Any("err", err))
		return "", err
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		ls.logger.Error("Failed to save file", slog.String("path", path), slog.Any("err", err))
		return "", err
	}

	return path, nil
}

func (ls *LocalStorage) DeleteFile(path string) error {
	if _, err := os.Stat(path); err == nil {
		err := os.Remove(path)
		if err != nil {
			ls.logger.Error("Failed to delete file", slog.String("path", path), slog.Any("err", err))
			return err
		}
	} else if errors.Is(err, os.ErrNotExist) {
		ls.logger.Info("File does not exist, nothing to delete", slog.String("path", path))
	} else {
		ls.logger.Error("Error checking file existence", slog.String("path", path), slog.Any("err", err))
	}

	return nil
}
