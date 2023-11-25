package storage

import (
	"log"

	"mediagui/domain"
	"mediagui/lib"
)

const getMovie = `
select rowid, title, original_title, file_title, year, runtime, tmdb_id, imdb_id,
				overview, tagline, resolution, filetype, location, cover, backdrop, genres, vote_average,
				vote_count, countries, added, modified, last_watched, all_watched, count_watched, score,
				director, writer, actors, awards, imdb_rating, imdb_votes, show_if_duplicate, stub
				from movie where rowid = ?
`

func (s *Storage) GetMovie(id uint64) *domain.Movie {
	s.mu.Lock()
	defer s.mu.Unlock()

	stmt, err := s.db.Prepare(getMovie)
	if err != nil {
		log.Fatalf("at get movie prepare: %s", err)
	}
	defer lib.Close(stmt)

	movie := domain.Movie{}

	err = stmt.
		QueryRow(id).
		Scan(&movie.ID, &movie.Title, &movie.Original_Title, &movie.FileTitle, &movie.Year, &movie.Runtime, &movie.Tmdb_Id, &movie.Imdb_Id, &movie.Overview, &movie.Tagline, &movie.Resolution, &movie.FileType, &movie.Location, &movie.Cover, &movie.Backdrop, &movie.Genres, &movie.Vote_Average, &movie.Vote_Count, &movie.Production_Countries, &movie.Added, &movie.Modified, &movie.Last_Watched, &movie.All_Watched, &movie.Count_Watched, &movie.Score, &movie.Director, &movie.Writer, &movie.Actors, &movie.Awards, &movie.Imdb_Rating, &movie.Imdb_Votes, &movie.ShowIfDuplicate, &movie.Stub)
	if err != nil {
		log.Fatalf("at get movie queryrow: %s", err)
	}

	return &movie
}
