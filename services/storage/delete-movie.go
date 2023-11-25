package storage

import (
	"log"

	"mediagui/lib"
)

func (s *Storage) DeleteMovie(id uint64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// lib.Notify(d.bus, "prune:delete", fmt.Sprintf("STARTED DELETING [%d] %s", movie.Id, movie.Title))

	tx, err := s.db.Begin()
	if err != nil {
		log.Fatalf("at delete movie begin: %s", err)
	}

	stmt, err := tx.Prepare("delete from movie where rowid = ?")
	if err != nil {
		rollback(tx)
		log.Fatalf("at delete movie prepare: %s", err)
	}
	defer lib.Close(stmt)

	_, err = stmt.Exec(id)
	if err != nil {
		rollback(tx)
		log.Fatalf("at delete movie exec: %s", err)
	}

	commit(tx)

	// lib.Notify(d.bus, "prune:delete", fmt.Sprintf("FINISHED DELETING [%d] %s", movie.Id, movie.Title))
}
