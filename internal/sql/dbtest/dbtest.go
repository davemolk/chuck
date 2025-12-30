package dbtest

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/davemolk/chuck/internal/sql"

	"github.com/davemolk/chuck/internal/migrations"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/stretchr/testify/require"
	postgresTest "github.com/testcontainers/testcontainers-go/modules/postgres"
	"go.uber.org/zap"
)

func SetupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	ctx := context.Background()

	pgContainer, err := postgresTest.Run(ctx,
		"postgres:17-alpine",
		postgresTest.WithDatabase("chucktest"),
		postgresTest.WithUsername("test"),
		postgresTest.WithPassword("test"),
		postgresTest.BasicWaitStrategies(),
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		err = pgContainer.Terminate(ctx)
		require.NoError(t, err)
	})

	dbURL, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	database, err := sql.New(zap.Must(zap.NewDevelopment()), dbURL)
	require.NoError(t, err)

	err = database.PingContext(ctx)
	require.NoError(t, err)

	err = runMigrations(database)
	require.NoError(t, err)

	return database
}

func runMigrations(db *sql.DB) error {
	source, err := iofs.New(migrations.MigrationFiles, ".")
	if err != nil {
		return fmt.Errorf("failed to get driver from fs: %w", err)
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to get db driver for migration: %w", err)
	}

	m, err := migrate.NewWithInstance(
		"iofs",
		source,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to get migrator: %w", err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migration failed: %w", err)
	}

	return nil
}
