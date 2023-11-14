package storage

import (
	"log"
	"time"

	"mediagui/domain"
	"mediagui/lib"
	"mediagui/logger"
)

func (s *Storage) SetDuplicate(movie *domain.Movie) *domain.Movie {
	logger.Blue("STARTED UPDATING MOVIE DUPLICATE STATUS [%d] %s (%d)", movie.ID, movie.Title, movie.ShowIfDuplicate)

	tx, err := s.db.Begin()
	if err != nil {
		log.Fatalf("at begin: %s", err)
	}

	stmt, err := tx.Prepare(`update movie set
								show_if_duplicate = ?,
								modified = ?
								where rowid = ?`)
	if err != nil {
		rollback(tx)
		log.Fatalf("at prepare: %s", err)
	}
	defer lib.Close(stmt)

	now := time.Now().UTC().Format(time.RFC3339)

	_, err = stmt.Exec(movie.ShowIfDuplicate, now, movie.ID)
	if err != nil {
		rollback(tx)
		log.Fatalf("at exec: %s", err)
	}

	commit(tx)
	logger.Blue("FINISHED UPDATING MOVIE DUPLICATE STATUS [%d] %s", movie.ID, movie.Title)

	return movie
}
