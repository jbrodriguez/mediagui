package storage

import (
	"log"
	"sort"
	"strings"
	"time"

	"mediagui/domain"
	"mediagui/lib"
	"mediagui/logger"
)

func (s *Storage) CopyMovie(movie *domain.Movie) *domain.Movie {
	s.mu.Lock()
	defer s.mu.Unlock()

	logger.Blue("STARTED COPYING MOVIE [%d] %s", movie.ID, movie.Title)

	now := time.Now().UTC().Format(time.RFC3339)

	stmt, err := s.db.Prepare("select all_watched, tmdb_id, score from movie where rowid = ?")
	if err != nil {
		log.Fatalf("at copy movie prepare 1: %s", err)
	}

	// get all watched times for the movie I'm copying from
	var when string
	var tmdb uint64
	var score uint64
	err = stmt.QueryRow(movie.Tmdb_Id).Scan(&when, &tmdb, &score)
	if err != nil {
		log.Fatalf("at copy movie queryrow: %s", err)
	}

	// create an array with all watched times
	watchedTimes := make([]string, 0)
	if when != "" {
		watchedTimes = strings.Split(when, "|")
	}

	// this are the watched times for the current movie
	currentWatched := make([]string, 0)
	if movie.All_Watched != "" {
		currentWatched = strings.Split(movie.All_Watched, "|")
	}

	for _, watched := range currentWatched {
		if !strings.Contains(when, watched) {
			watchedTimes = append(watchedTimes, watched)
		}
	}

	// this sorts the dates in ascending order by default
	sort.Strings(watchedTimes)

	// set final variables
	lastWatched := watchedTimes[len(watchedTimes)-1]
	countWatched := uint64(len(watchedTimes))
	allWatched := strings.Join(watchedTimes, "|")

	tx, err := s.db.Begin()
	if err != nil {
		log.Fatalf("at copy movie begin: %s", err)
	}

	stmt, err = tx.Prepare(`update movie set
								last_watched = ?,
								all_watched = ?,
								count_watched = ?,
								score = ?
								modified = ?
								where rowid = ?`)
	if err != nil {
		rollback(tx)
		log.Fatalf("at copy movie prepare 2: %s", err)
	}
	defer lib.Close(stmt)

	_, err = stmt.Exec(lastWatched, allWatched, countWatched, score, now, movie.ID)
	if err != nil {
		rollback(tx)
		log.Fatalf("at copy movie exec: %s", err)
	}

	commit(tx)
	logger.Blue("FINISHED COPYING MOVIE [%d] %s", movie.ID, movie.Title)

	movie.All_Watched = allWatched
	movie.Count_Watched = countWatched
	movie.Last_Watched = lastWatched
	movie.Modified = now
	movie.Tmdb_Id = tmdb

	return movie
}
