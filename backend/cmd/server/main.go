package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/anujgupta/level-up-backend/internal/auth"
	"github.com/anujgupta/level-up-backend/internal/config"
	appdb "github.com/anujgupta/level-up-backend/internal/db"
	dbgen "github.com/anujgupta/level-up-backend/generated/db"
	"github.com/anujgupta/level-up-backend/internal/mailer"
	"github.com/anujgupta/level-up-backend/internal/server"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	if err := run(logger); err != nil {
		logger.Error("server error", "err", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	// 1. Load config
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	// 2. Connect to database
	pool, err := appdb.Connect(cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()
	logger.Info("database connected")

	// 3. Run migrations
	if err := appdb.RunMigrations(cfg.DatabaseURL); err != nil {
		return err
	}
	logger.Info("migrations applied")

	// 4. Build SQLC queries
	queries := dbgen.New(pool)

	// 5. Build auth service
	authSvc := auth.NewService(cfg)

	// 6. Build and start mailer
	mailerSvc := mailer.New(cfg, logger)
	mailerSvc.Start(3)
	logger.Info("mailer started", "workers", 3)

	// 7. Build server (wires all handlers + middleware + router)
	srv := server.New(cfg, queries, authSvc, mailerSvc, logger)

	// 8. Graceful shutdown on SIGINT / SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		logger.Info("shutdown signal received")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		// Stop accepting new requests; wait for in-flight to finish
		if err := srv.HTTPServer().Shutdown(shutdownCtx); err != nil {
			logger.Error("http server shutdown error", "err", err)
		}

		// Drain email worker queue
		mailerSvc.Close()

		// Pool closed via defer above
	}()

	// 9. Serve — blocks until shutdown
	if err := srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	logger.Info("server stopped cleanly")
	return nil
}

// Ensure pgxpool is used (satisfies SQLC's pgxpool adapter).
var _ *pgxpool.Pool
