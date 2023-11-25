package storage

import (
	"fmt"
	"log"

	"mediagui/domain"
	"mediagui/lib"
	"mediagui/logger"
)

const listMovies = `
select rowid, title, original_title, file_title, year, runtime, tmdb_id, imdb_id,
				overview, tagline, resolution, filetype, location, cover, backdrop, genres, vote_average,
				vote_count, countries, added, modified, last_watched, all_watched, count_watched, score,
				director, writer, actors, awards, imdb_rating, imdb_votes, show_if_duplicate, stub
				from movie order by %s %s, rowid desc limit ? offset ?
`

func (s *Storage) listMovies(options *domain.Options) (total uint64, movies []*domain.Movie) {
	s.mu.Lock()
	defer s.mu.Unlock()

	tx, err := s.db.Begin()
	if err != nil {
		log.Fatalf("at list movies begin: %s", err)
	}

	statement := fmt.Sprintf(listMovies, options.SortBy, options.SortOrder)

	stmt, err := tx.Prepare(statement)
	if err != nil {
		log.Fatalf("at list movies prepare 1: %s", err)
	}
	defer lib.Close(stmt)

	rows, err := stmt.Query(options.Limit, options.Offset)
	if err != nil {
		log.Fatalf("at list movies query: %s", err)
	}

	items := make([]*domain.Movie, 0)

	// if options.Offset == 0 {
	// if d.count == 0 {
	var count uint64
	stmt, err = tx.Prepare("select count(*) from movie;")
	if err != nil {
		log.Fatalf("at list movies prepare 2: %s", err)
	}
	defer lib.Close(stmt)

	err = stmt.QueryRow().Scan(&count)
	if err != nil {
		log.Fatalf("at list movies queryrow: %s", err)
	}

	for rows.Next() {
		movie := domain.Movie{}
		rows.Scan(
			&movie.ID, &movie.Title, &movie.Original_Title, &movie.FileTitle, &movie.Year,
			&movie.Runtime, &movie.Tmdb_Id, &movie.Imdb_Id, &movie.Overview, &movie.Tagline,
			&movie.Resolution, &movie.FileType, &movie.Location, &movie.Cover, &movie.Backdrop,
			&movie.Genres, &movie.Vote_Average, &movie.Vote_Count, &movie.Production_Countries,
			&movie.Added, &movie.Modified, &movie.Last_Watched, &movie.All_Watched,
			&movie.Count_Watched, &movie.Score, &movie.Director, &movie.Writer, &movie.Actors,
			&movie.Awards, &movie.Imdb_Rating, &movie.Imdb_Votes, &movie.ShowIfDuplicate, &movie.Stub,
		)
		items = append(items, &movie)
	}
	lib.Close(rows)

	commit(tx)

	logger.Blue("Listed %d movies (total %d)", len(items), count)

	return count, items
}
