package storage

import (
	"database/sql"
	"log"

	"mediagui/domain"
	"mediagui/lib"
)

func (s *Storage) CheckExists(movie *domain.Movie) bool {
	tx, err := s.db.Begin()
	if err != nil {
		log.Fatalf("at begin: %s", err)
	}

	stmt, err := tx.Prepare("select rowid from movie where location = ?")
	if err != nil {
		rollback(tx)
		log.Fatalf("at prepare: %s", err)
	}
	defer lib.Close(stmt)

	var id int
	err = stmt.QueryRow(movie.Location).Scan(&id)
	if err != sql.ErrNoRows && err != nil {
		rollback(tx)
		log.Fatalf("at queryrow: %s", err)
	}

	commit(tx)

	return (id != 0)
}
