package storage

import (
	"database/sql"
	_ "embed"
	"os"
	"path/filepath"
	"sync"

	_ "modernc.org/sqlite"

	"mediagui/domain"
	"mediagui/logger"
)

//go:embed ../../metadata/schema.sql
var schemaSQL string

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

	// Initialize database schema if it doesn't exist
	if !dbExists {
		if err := s.initializeSchema(); err != nil {
			return err
		}
		logger.Blue("initialized new database with schema (%s)", dbPath)
	}

	logger.Blue("started service storage (%s)", dbPath)

	return nil
}

func (s *Storage) initializeSchema() error {
	_, err := s.db.Exec(schemaSQL)
	return err
}
