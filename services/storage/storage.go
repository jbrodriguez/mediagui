package storage

import (
	"database/sql"
	"path/filepath"

	_ "modernc.org/sqlite"

	"mediagui/domain"
	"mediagui/logger"
)

type Storage struct {
	ctx *domain.Context
	db  *sql.DB
}

func Create(ctx *domain.Context) *Storage {
	return &Storage{
		ctx: ctx,
	}
}

func (s *Storage) Start() error {
	dbPath := filepath.Join(s.ctx.DataDir, "mediagui.db")

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}
	s.db = db

	logger.Blue("started service storage (%s)", dbPath)

	return nil
}
