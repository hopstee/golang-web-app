package sqlite

import (
	"log/slog"
	"mobile-backend-boilerplate/pkg/logger"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSQLiteRepository_Migrate(t *testing.T) {
	var err error
	dsn := ":inmemory:"
	log := logger.New(slog.LevelDebug)

	repo, err := NewSQLiteRepository(dsn, log)
	require.NoError(t, err, "failed to create SQLiteRepository")

	defer repo.DB.Close()

	require.NotNil(t, repo.DB, "DB should not be nil")

	tables := []string{"refresh_tokens", "users", "posts", "requests", "admins"}
	for _, table := range tables {
		var exists int
		err = repo.DB.QueryRow(
			"SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?",
			table,
		).Scan(&exists)
		require.NoError(t, err, "failed to query sqlite_master")
		require.Equal(t, 1, exists, "there should be at three tables after migration")
	}

	err = repo.Migrate(Down, false)
	require.NoError(t, err)

	for _, table := range tables {
		var count int
		err := repo.DB.QueryRow(
			"SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?",
			table,
		).Scan(&count)
		require.NoError(t, err)
		require.Equal(t, 0, count, "table %s should be dropped after Down migration", table)
	}
}
