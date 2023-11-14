package storage

import (
	"log"
	"mediagui/domain"
	"mediagui/lib"
	"mediagui/logger"
)

func (s *Storage) GetDuplicates() (total uint64, items []*domain.Movie) {
	logger.Blue("getDuplicates.starting")

	tx, err := s.db.Begin()
	if err != nil {
		log.Fatalf("Unable to begin transaction: %s", err)
	}

	rows, err := s.db.Query(`select a.rowid, a.title, a.original_title, a.file_title,
				a.year, a.runtime, a.tmdb_id, a.imdb_id, a.overview, a.tagline, a.resolution,
				a.filetype, a.location, a.cover, a.backdrop, a.genres, a.vote_average,
				a.vote_count, a.countries, a.added, a.modified, a.last_watched, a.all_watched,
				a.count_watched, a.score, a.director, a.writer, a.actors, a.awards, a.imdb_rating,
				a.imdb_votes, a.show_if_duplicate, a.stub
				from
				movie a
				join
				(select title, show_if_duplicate from movie where show_if_duplicate = 1 group by title having count(*) > 1) b
				on a.title = b.title;`)
	if err != nil {
		log.Fatalf("Unable to prepare transaction: %s", err)
	}

	items = make([]*domain.Movie, 0)

	for rows.Next() {
		movie := domain.Movie{}
		rows.Scan(&movie.ID, &movie.Title, &movie.Original_Title, &movie.FileTitle, &movie.Year, &movie.Runtime, &movie.Tmdb_Id, &movie.Imdb_Id, &movie.Overview, &movie.Tagline, &movie.Resolution, &movie.FileType, &movie.Location, &movie.Cover, &movie.Backdrop, &movie.Genres, &movie.Vote_Average, &movie.Vote_Count, &movie.Production_Countries, &movie.Added, &movie.Modified, &movie.Last_Watched, &movie.All_Watched, &movie.Count_Watched, &movie.Score, &movie.Director, &movie.Writer, &movie.Actors, &movie.Awards, &movie.Imdb_Rating, &movie.Imdb_Votes, &movie.ShowIfDuplicate, &movie.Stub)
		items = append(items, &movie)
	}
	lib.Close(rows)

	commit(tx)

	logger.Blue("Found %d duplicate movies", len(items))

	return uint64(len(items)), items
}
