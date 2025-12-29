package dbtest

import (
	"context"
	"fmt"
	"testing"

	"github.com/davemolk/chuck/internal/sql"

	"github.com/davemolk/chuck/internal/migrations"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
)

func SetupTestDB(t *testing.T) *sql.DB {
	t.Helper()
	ctx := context.Background()

	postgres, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:15-alpine",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_DB":       "chucktest",
				"POSTGRES_USER":     "test",
				"POSTGRES_PASSWORD": "test",
			},
			WaitingFor: wait.ForListeningPort("5432/tcp"),
		},
		Started: true,
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		err = postgres.Terminate(ctx)
		require.NoError(t, err)
	})

	host, err := postgres.Host(ctx)
	require.NoError(t, err)

	port, err := postgres.MappedPort(ctx, "5432")
	require.NoError(t, err)

	dbURL := fmt.Sprintf("postgres://test:test@%s:%s/chucktest?sslmode=disable", host, port.Port())

	database, err := sql.New(zap.Must(zap.NewDevelopment()), dbURL)
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

	return m.Up()
}
