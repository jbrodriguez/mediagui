package storage

import (
	"log"

	"mediagui/domain"
	"mediagui/lib"
	"mediagui/logger"
)

func (s *Storage) UpdateMovie(movie *domain.Movie) {
	s.mu.Lock()
	defer s.mu.Unlock()

	logger.Blue("STARTED UPDATING [%d] %s", movie.ID, movie.Title)

	tx, err := s.db.Begin()
	if err != nil {
		log.Fatalf("at update movie begin: %s", err)
	}

	stmt, err := tx.Prepare(`update movie set title = ?,
								original_title = ?,
								year = ?,
								runtime = ?,
								tmdb_id = ?,
								imdb_id = ?,
								overview = ?,
								tagline = ?,
								cover = ?,
								backdrop = ?,
								genres = ?,
								vote_average = ?,
								vote_count = ?,
								countries = ?,
								modified = ?,
								director = ?,
								writer = ?,
								actors = ?,
								awards = ?,
								imdb_rating = ?,
								imdb_votes = ?,
								show_if_duplicate = ?
								where rowid = ?`)
	if err != nil {
		rollback(tx)
		log.Fatalf("at update movie prepare: %s", err)
	}
	defer lib.Close(stmt)

	_, err = stmt.Exec(movie.Title, movie.Original_Title, movie.Year, movie.Runtime, movie.Tmdb_Id, movie.Imdb_Id, movie.Overview, movie.Tagline, movie.Cover, movie.Backdrop, movie.Genres, movie.Vote_Average, movie.Vote_Count, movie.Production_Countries, movie.Modified, movie.Director, movie.Writer, movie.Actors, movie.Awards, movie.Imdb_Rating, movie.Imdb_Votes, movie.ShowIfDuplicate, movie.ID)
	if err != nil {
		rollback(tx)
		log.Fatalf("at update movie exec: %s", err)
	}

	commit(tx)
	logger.Blue("FINISHED UPDATING [%d] %s", movie.ID, movie.Title)
}
