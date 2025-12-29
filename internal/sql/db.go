package sql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type DB struct {
	logger *zap.Logger
	*sql.DB
}

func New(logger *zap.Logger, dbURL string) (*DB, error) {
	sqlDB, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("faield to open db: %w", err)
	}

	// Note, with more time, would expose these to caller for more control
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	return &DB{
		logger: logger,
		DB:     sqlDB,
	}, nil
}

func (db *DB) Ping(ctx context.Context) error {
	if err := db.PingContext(ctx); err != nil {
		return err
	}

	db.logger.Info("ping successful, database is connected")
	return nil
}

func (db *DB) Close() error {
	db.logger.Info("closing database")
	return db.DB.Close()
}

// note: while it's true we only have one insert statement for mvp on this
// project, this is helpful to have.
func (db *DB) RunInTx(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			// rollback and panic again
			_ = tx.Rollback()
			panic(r)
		} else if err != nil {
			// just rollback
			_ = tx.Rollback()
		}
	}()

	if err = fn(tx); err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit tx: %w", err)
	}

	return nil
}
