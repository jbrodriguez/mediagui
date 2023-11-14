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
	logger.Blue("STARTED COPYING MOVIE WATCHED TIMES [%d] %s (%s)", movie.ID, movie.Title, movie.Last_Watched)

	now := time.Now().UTC().Format(time.RFC3339)

	stmt, err := s.db.Prepare("select all_watched, tmdb_id from movie where rowid = ?")
	if err != nil {
		log.Fatalf("at prepare: %s", err)
	}

	// get all watched times for the movie I'm copying from
	var when string
	var tmdb uint64
	err = stmt.QueryRow(movie.Tmdb_Id).Scan(&when, &tmdb)
	if err != nil {
		log.Fatalf("at queryrow: %s", err)
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
		log.Fatalf("at begin: %s", err)
	}

	stmt, err = tx.Prepare(`update movie set
								last_watched = ?,
								all_watched = ?,
								count_watched = ?,
								modified = ?
								where rowid = ?`)
	if err != nil {
		rollback(tx)
		log.Fatalf("at prepare: %s", err)
	}
	defer lib.Close(stmt)

	_, err = stmt.Exec(lastWatched, allWatched, countWatched, now, movie.ID)
	if err != nil {
		rollback(tx)
		log.Fatalf("at exec: %s", err)
	}

	commit(tx)
	logger.Blue("FINISHED COPYING MOVIE WATCHED TIMES [%d] %s", movie.ID, movie.Title)

	movie.All_Watched = allWatched
	movie.Count_Watched = countWatched
	movie.Last_Watched = lastWatched
	movie.Modified = now
	movie.Tmdb_Id = tmdb

	return movie
}
