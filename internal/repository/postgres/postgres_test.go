package postgres

import (
	"fmt"
	"log/slog"
	"mobile-backend-boilerplate/pkg/logger"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPostgreSQLRepository_Migrate(t *testing.T) {
	dsn := "postgres://user:password@localhost:5432/testdb?sslmode=disable" // замените на вашу тестовую БД
	log := logger.New(slog.LevelDebug)

	repo, err := NewPostgreSQLRepository(dsn, log)
	require.NoError(t, err, "failed to create PostgreSQLRepository")
	defer repo.DB.Close()

	require.NotNil(t, repo.DB, "DB should not be nil")

	tables := []string{"refresh_tokens", "users", "posts", "requests", "admins"}
	for _, table := range tables {
		var exists bool
		query := `
			SELECT EXISTS (
				SELECT 1
				FROM information_schema.tables
				WHERE table_schema = 'public' AND table_name = $1
			)
		`
		err = repo.DB.QueryRow(query, table).Scan(&exists)
		require.NoError(t, err, "failed to query information_schema")
		require.True(t, exists, fmt.Sprintf("table %s should exist after migration", table))
	}

	// Проверка Down миграции
	err = repo.Migrate(Down, false)
	require.NoError(t, err, "failed to run Down migration")

	for _, table := range tables {
		var exists bool
		query := `
			SELECT EXISTS (
				SELECT 1
				FROM information_schema.tables
				WHERE table_schema = 'public' AND table_name = $1
			)
		`
		err = repo.DB.QueryRow(query, table).Scan(&exists)
		require.NoError(t, err)
		require.False(t, exists, fmt.Sprintf("table %s should be dropped after Down migration", table))
	}
}
