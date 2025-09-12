package storage

import (
	"database/sql"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "modernc.org/sqlite"

	"mediagui/domain"
	"mediagui/logger"
)

//go:embed all:migrations
var migrationsFS embed.FS

type Storage struct {
	ctx *domain.Context
	db  *sql.DB
	mu  sync.Mutex
}

func Create(ctx *domain.Context) *Storage {
	return &Storage{
		ctx: ctx,
	}
}

func (s *Storage) Start() error {
	dbDir := filepath.Join(s.ctx.DataDir, "db")
	dbPath := filepath.Join(dbDir, "mediagui.db")

	// Ensure the db directory exists
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return err
	}

	// Check if database exists
	dbExists := true
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		dbExists = false
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}
	s.db = db

	// Run migrations
	if err := s.runMigrations(dbPath); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	if !dbExists {
		logger.Blue("initialized new database with migrations (%s)", dbPath)
	} else {
		logger.Blue("applied migrations to existing database (%s)", dbPath)
	}

	logger.Blue("started service storage (%s)", dbPath)

	return nil
}

func (s *Storage) runMigrations(dbPath string) error {
	// Create source from embedded filesystem
	sourceDriver, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return fmt.Errorf("could not create migration source: %w", err)
	}

	// Create migrate instance using database URL
	databaseURL := fmt.Sprintf("sqlite://%s", dbPath)
	m, err := migrate.NewWithSourceInstance("iofs", sourceDriver, databaseURL)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}
	defer m.Close()

	// Run migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("could not run migrations: %w", err)
	}

	return nil
}
