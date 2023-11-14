package storage

import (
	"log"
	"time"

	"mediagui/domain"
	"mediagui/lib"
	"mediagui/logger"
)

func (s *Storage) SetMovieScore(movie *domain.Movie) {
	logger.Blue("STARTED UPDATING MOVIE SCORE [%d] %s (%d)", movie.ID, movie.Title, movie.Score)

	tx, err := s.db.Begin()
	if err != nil {
		log.Fatalf("at begin: %s", err)
	}

	stmt, err := tx.Prepare(`update movie set
								score = ?,
								modified = ?
								where rowid = ?`)
	if err != nil {
		rollback(tx)
		log.Fatalf("at prepare: %s", err)
	}
	defer lib.Close(stmt)

	now := time.Now().UTC().Format(time.RFC3339)

	_, err = stmt.Exec(movie.Score, now, movie.ID)
	if err != nil {
		rollback(tx)
		log.Fatalf("at exec: %s", err)
	}

	commit(tx)
	logger.Blue("FINISHED UPDATING MOVIE SCORE [%d] %s", movie.ID, movie.Title)
}
