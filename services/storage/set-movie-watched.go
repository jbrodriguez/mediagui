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

func (s *Storage) SetMovieWatched(movie *domain.Movie) *domain.Movie {
	logger.Blue("STARTED UPDATING MOVIE WATCHED DATE [%d] %s (%s)", movie.ID, movie.Title, movie.Last_Watched)

	now := time.Now().UTC().Format(time.RFC3339)

	when := movie.All_Watched

	// create an array with all watched times
	var watchedTimes []string
	if when != "" {
		watchedTimes = strings.Split(when, "|")
	}

	// convert incoming watched time to sane format
	watched, err := parseToday(movie.Last_Watched)
	if err != nil {
		log.Fatalf("at parseToday: %s", err)
	}
	lastWatched := watched.UTC().Format(time.RFC3339)

	// add last watched to array, only if it doesn't already exist
	if !strings.Contains(when, lastWatched) {
		watchedTimes = append(watchedTimes, lastWatched)
	}

	// this sorts the dates in ascending order by default
	sort.Strings(watchedTimes)

	// set final variables
	lastWatched = watchedTimes[len(watchedTimes)-1]
	countWatched := uint64(len(watchedTimes))
	allWatched := strings.Join(watchedTimes, "|")

	tx, err := s.db.Begin()
	if err != nil {
		log.Fatalf("at begin: %s", err)
	}

	stmt, err := tx.Prepare(`update movie set
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
	logger.Blue("FINISHED UPDATING MOVIE WATCHED DATE [%d] %s", movie.ID, movie.Title)

	movie.All_Watched = allWatched
	movie.Count_Watched = countWatched
	movie.Modified = now

	return movie
}

func parseToday(clientToday string) (today time.Time, err error) {
	client, perr := time.Parse(time.RFC3339, clientToday)
	if perr != nil {
		return today, perr
	}

	today = time.Date(client.Year(), client.Month(), client.Day(), 0, 0, 0, 0, client.Location())

	return today, nil
}
