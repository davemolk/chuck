package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	apihttp "github.com/davemolk/chuck/internal/api/http"
	"github.com/davemolk/chuck/internal/clients/chuck"
	"github.com/davemolk/chuck/internal/service/joke"
	"github.com/davemolk/chuck/internal/service/user"
	"github.com/davemolk/chuck/internal/sql"
	"go.uber.org/zap"
)

type config struct {
	Production bool
	Port       int
	DbURL      string
}

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	cfg, err := fromEnv()
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	logger, err := rootLogger(cfg.Production)
	if err != nil {
		return fmt.Errorf("failed to get logger: %w", err)
	}

	defer func() { _ = logger.Sync() }()

	db, err := sql.New(logger, cfg.DbURL)
	if err != nil {
		return fmt.Errorf("failed to create db: %w", err)
	}

	if err = db.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping db: %w", err)
	}

	chuckClient := chuck.NewClient(logger)
	jokeService := joke.NewService(logger, db, chuckClient)
	userService := user.NewService(logger, db, nil)

	router := apihttp.NewRoutes(logger, &apihttp.Services{
		JokeService: jokeService,
		UserService: userService,
	})

	srv := apihttp.NewServer(logger, cfg.Port, router)

	shutdownErr := make(chan error)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sigRecv := <-quit

		logger.Info("starting graceful shutdown", zap.Stringer("signal", sigRecv))

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err := srv.Shutdown(shutdownCtx)
		if err != nil {
			shutdownErr <- err
		}

		logger.Info("chuck norris delivered a roundhouse to the server")

		err = db.Close()
		if err != nil {
			shutdownErr <- err
		}

		logger.Info("chuck norris told the db to get chucked")

		shutdownErr <- nil
	}()

	err = srv.Run()
	if !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server failed: %w", err)
	}

	err = <-shutdownErr
	if err != nil {
		return fmt.Errorf("error shutting down: %w", err)
	}

	logger.Info("exiting app...")

	return nil
}

func rootLogger(prod bool) (*zap.Logger, error) {
	var l *zap.Logger
	var err error

	if prod {
		l, err = zap.NewProduction()
	} else {
		l, err = zap.NewDevelopment()
	}
	if err != nil {
		return nil, fmt.Errorf("failed to init logger: %w", err)
	}

	return l, nil
}

func fromEnv() (*config, error) {
	prodStr := os.Getenv("PRODUCTION")
	portStr := os.Getenv("PORT")
	dbURL := os.Getenv("DATABASE_URL")

	prod, err := strconv.ParseBool(prodStr)
	if err != nil {
		return nil, fmt.Errorf("parsing production bool: %w", err)
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("parsing port: %w", err)
	}
	if port == 0 {
		return nil, errors.New("port must be set")
	}

	if dbURL == "" {
		return nil, errors.New("database url must be set")
	}

	return &config{
		Production: prod,
		Port:       port,
		DbURL:      dbURL,
	}, nil
}
