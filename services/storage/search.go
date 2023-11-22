package storage

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"strings"

	"mediagui/domain"
	"mediagui/lib"
	"mediagui/logger"
)

func (s *Storage) searchMovies(options *domain.Options) (total uint64, movies []*domain.Movie) {
	if options.FilterBy == "year" {
		total, movies = s.searchByYear(options)
	} else {
		total, movies = s.regularSearch(options)
	}

	return total, movies
}

func (s *Storage) searchByYear(options *domain.Options) (total uint64, movies []*domain.Movie) {
	var start, end, year uint64
	decade := false

	if strings.Contains(options.Query, "-") {
		decade = true
		parts := strings.Split(options.Query, "-")

		start, _ = strconv.ParseUint(parts[0], 0, 64)
		end, _ = strconv.ParseUint(parts[1], 0, 64)
	} else {
		year, _ = strconv.ParseUint(options.Query, 0, 64)
	}

	var countQuery string
	if decade {
		countQuery = `select count(*) from movie where year between ? and ?;`
	} else {
		countQuery = `select count(*) from movie where year = ?;`
	}

	tx, err := s.db.Begin()
	if err != nil {
		log.Fatalf("searchByYear:Unable to begin transaction: %s", err)
	}

	stmt, err := tx.Prepare(countQuery)
	if err != nil {
		log.Fatalf("searchByYear:Unable to prepare transaction: %s", err)
	}
	defer lib.Close(stmt)

	var count uint64
	stmt, err = tx.Prepare("select count(*) from movie;")
	if err != nil {
		log.Fatalf("Unable to prepare count rows transaction: %s", err)
	}
	defer lib.Close(stmt)

	err = stmt.QueryRow().Scan(&count)
	if err != nil {
		log.Fatalf("Unable to count rows: %s", err)
	}

	var rows *sql.Rows
	if decade {
		err = stmt.QueryRow(start, end).Scan(&count)
		if err != nil {
			log.Fatalf("searchByYear:Unable to count rows: %s", err)
		}

		query := fmt.Sprintf(`select rowid, title, original_title, file_title, year, runtime, tmdb_id, imdb_id,
					overview, tagline, resolution, filetype, location, cover, backdrop, genres, vote_average,
					vote_count, countries, added, modified, last_watched, all_watched, count_watched, score,
					director, writer, actors, awards, imdb_rating, imdb_votes, show_if_duplicate, stub
					from movie where year between ? and ? order by %s %s, rowid desc limit ? offset ?`, options.SortBy, options.SortOrder)

		stmt, err = tx.Prepare(query)
		if err != nil {
			log.Fatalf("Unable to prepare transaction: %s", err)
		}
		defer lib.Close(stmt)

		rows, err = stmt.Query(start, end, options.Limit, options.Offset)
		if err != nil {
			log.Fatalf("Unable to prepare transaction: %s", err)
		}
	} else {
		err = stmt.QueryRow(year).Scan(&count)
		if err != nil {
			log.Fatalf("searchByYear:Unable to count rows: %s", err)
		}

		query := fmt.Sprintf(`select rowid, title, original_title, file_title, year, runtime, tmdb_id, imdb_id,
					overview, tagline, resolution, filetype, location, cover, backdrop, genres, vote_average,
					vote_count, countries, added, modified, last_watched, all_watched, count_watched, score,
					director, writer, actors, awards, imdb_rating, imdb_votes, show_if_duplicate, stub
					from movie where year = ? order by %s %s, rowid desc limit ? offset ?`, options.SortBy, options.SortOrder)

		logger.Blue("sql.: %s", query)

		stmt, err := tx.Prepare(query)
		if err != nil {
			log.Fatalf("Unable to prepare transaction: %s", err)
		}
		defer lib.Close(stmt)

		rows, err = stmt.Query(year, options.Limit, options.Offset)
		if err != nil {
			log.Fatalf("Unable to prepare transaction: %s", err)
		}
	}

	items := make([]*domain.Movie, 0)

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

	logger.Blue("searchByYear:Listed %d movies (total %d)", len(items), count)

	return count, items
}

func (s *Storage) regularSearch(options *domain.Options) (total uint64, movies []*domain.Movie) {
	logger.Blue("searchMovies.options: %+v", options)

	tx, err := s.db.Begin()
	if err != nil {
		log.Fatalf("searchMovies:Unable to begin transaction: %s", err)
	}

	term := options.Query + "*"
	args := options.FilterBy

	logger.Blue("this is: %s %s", term, args)

	// if self.searchArgs != args {
	// self.searchArgs = args

	countQuery := fmt.Sprintf(`select count(*) from movie dt, %s vt
		where vt.%s match ? and dt.rowid = vt.rowid;`,
		"movie"+options.FilterBy, "movie"+options.FilterBy)

	stmt, err := tx.Prepare(countQuery)
	if err != nil {
		log.Fatalf("searchMovies:Unable to prepare transaction: %s", err)
	}
	defer lib.Close(stmt)

	logger.Blue("searchMovies:CountQuery:%s", countQuery)

	var count uint64
	err = stmt.QueryRow(term).Scan(&count)
	if err != nil {
		log.Fatalf("searchMovies:Unable to count rows: %s", err)
	}

	listQuery := fmt.Sprintf(`select dt.rowid, dt.title, dt.original_title, dt.year, dt.runtime,
			dt.tmdb_id, dt.imdb_id, dt.overview, dt.tagline, dt.resolution,
			dt.filetype, dt.location, dt.cover, dt.backdrop, dt.genres, dt.vote_average,
			dt.vote_count, dt.countries, dt.added, dt.modified, dt.last_watched,
			dt.all_watched, dt.count_watched, dt.score, dt.director, dt.writer, dt.actors,
			dt.awards, dt.imdb_rating, dt.imdb_votes, dt.show_if_duplicate, dt.stub
			from movie dt, %s vt
			where vt.%s match ? and dt.rowid = vt.rowid order by dt.%s %s limit ? offset ?`,
		"movie"+options.FilterBy, "movie"+options.FilterBy, options.SortBy, options.SortOrder)

	stmt, err = tx.Prepare(listQuery)
	if err != nil {
		log.Fatalf("searchMovies:listQuery:Unable to prepare transaction: %s", err)
	}
	defer lib.Close(stmt)

	rows, err := stmt.Query(term, options.Limit, options.Offset)
	if err != nil {
		log.Fatalf("searchMovies:listQuery:Unable to begin transaction: %s", err)
	}

	items := make([]*domain.Movie, 0)

	for rows.Next() {
		movie := domain.Movie{}
		rows.Scan(&movie.ID, &movie.Title, &movie.Original_Title, &movie.Year, &movie.Runtime, &movie.Tmdb_Id, &movie.Imdb_Id, &movie.Overview, &movie.Tagline, &movie.Resolution, &movie.FileType, &movie.Location, &movie.Cover, &movie.Backdrop, &movie.Genres, &movie.Vote_Average, &movie.Vote_Count, &movie.Production_Countries, &movie.Added, &movie.Modified, &movie.Last_Watched, &movie.All_Watched, &movie.Count_Watched, &movie.Score, &movie.Director, &movie.Writer, &movie.Actors, &movie.Awards, &movie.Imdb_Rating, &movie.Imdb_Votes, &movie.ShowIfDuplicate, &movie.Stub)
		items = append(items, &movie)
	}
	lib.Close(rows)

	commit(tx)

	logger.Blue("searchMovies:Listed %d movies (total %d)", len(items), count)

	return count, items
}
