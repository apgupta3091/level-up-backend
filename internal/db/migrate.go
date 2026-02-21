package db

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations applies all pending up migrations from db/migrations/.
// The binary must be run from the project root so the relative path resolves correctly.
func RunMigrations(databaseURL string) error {
	// golang-migrate pgx/v5 driver expects "pgx5://" scheme
	pgxURL := "pgx5://" + databaseURL[len("postgres://"):]

	m, err := migrate.New("file://db/migrations", pgxURL)
	if err != nil {
		return fmt.Errorf("create migrator: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}
