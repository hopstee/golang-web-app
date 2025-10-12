package postgres

import (
	"context"
	"fmt"
	"path"

	"github.com/mattermost/morph"
	"github.com/mattermost/morph/drivers/sqlite"
	mbindata "github.com/mattermost/morph/sources/embedded"
)

const (
	migrationsTableName = "db_migrations"
)

type Direction int

const (
	Up Direction = iota
	Down
)

func (s *PostgreSQLRepository) initMorph(dryRun bool, timeoutSeconds int) (*morph.Morph, error) {
	driver, err := sqlite.WithInstance(s.DB)
	if err != nil {
		return nil, fmt.Errorf("failed to init morph sqlite driver: %w", err)
	}

	assetsList, err := assets.ReadDir("migrations")
	if err != nil {
		return nil, err
	}

	assetNamesForDriver := make([]string, len(assetsList))
	for i, entry := range assetsList {
		assetNamesForDriver[i] = entry.Name()
	}

	src, err := mbindata.WithInstance(&mbindata.AssetSource{
		Names: assetNamesForDriver,
		AssetFunc: func(name string) ([]byte, error) {
			return assets.ReadFile(path.Join("migrations", name))
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to generate source assets: %w", err)
	}

	opts := []morph.EngineOption{
		// PostgreSQL does not support locking
		// morph.WithLock("migrations-lock-key"),
		morph.SetMigrationTableName(migrationsTableName),
		morph.SetStatementTimeoutInSeconds(timeoutSeconds),
		morph.SetDryRun(dryRun),
	}

	engine, err := morph.New(context.Background(), driver, src, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to creare morph engine: %w", err)
	}

	return engine, nil
}

func (s *PostgreSQLRepository) Migrate(direction Direction, dryRun bool) error {
	engine, err := s.initMorph(dryRun, 30)
	if err != nil {
		return fmt.Errorf("failed to initialize morph: %w", err)
	}
	defer engine.Close()

	switch direction {
	case Down:
		_, err = engine.ApplyDown(-1)
		return err
	default:
		return engine.ApplyAll()
	}
}
