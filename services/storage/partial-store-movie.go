package storage

import (
	"log"
	"time"

	"mediagui/domain"
	"mediagui/lib"
	"mediagui/logger"
)

func (s *Storage) PartialStoreMovie(movie *domain.Movie) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UTC().Format(time.RFC3339)
	movie.Added = now
	movie.Modified = now
	movie.ShowIfDuplicate = 1

	logger.Blue("STARTED PARTIAL SAVING %s", movie.Title)

	tx, err := s.db.Begin()
	if err != nil {
		log.Fatalf("at partial store movie begin: %s", err)
	}

	stmt, err := tx.Prepare(`insert into movie(title, original_title, file_title,
								year, runtime, tmdb_id, imdb_id, overview, tagline,
								resolution, filetype, location, cover, backdrop, genres,
								vote_average, vote_count, countries, added, modified,
								last_watched, all_watched, count_watched, score, director,
								writer, actors, awards, imdb_rating, imdb_votes, show_if_duplicate, stub)
								values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
									?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		rollback(tx)
		log.Fatalf("at partial store movie prepare: %s", err)
	}
	defer lib.Close(stmt)

	res, e := stmt.Exec(movie.Title, movie.Original_Title, movie.FileTitle, movie.Year,
		movie.Runtime, movie.Tmdb_Id, movie.Imdb_Id, movie.Overview, movie.Tagline,
		movie.Resolution, movie.FileType, movie.Location, movie.Cover, movie.Backdrop,
		movie.Genres, movie.Vote_Average, movie.Vote_Count, movie.Production_Countries,
		movie.Added, movie.Modified, movie.Last_Watched, movie.All_Watched, movie.Count_Watched,
		movie.Score, movie.Director, movie.Writer, movie.Actors, movie.Awards, movie.Imdb_Rating,
		movie.Imdb_Votes, movie.ShowIfDuplicate, movie.Stub)
	if e != nil {
		rollback(tx)
		log.Fatalf("at partial store movie exec: %s", e)
	}

	id, _ := res.LastInsertId()
	movie.ID = uint64(id)

	commit(tx)
	logger.Blue("FINISHED PARTIAL SAVING %s [%d]", movie.Title, movie.ID)
}
