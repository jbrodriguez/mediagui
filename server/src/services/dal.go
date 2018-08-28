package services

import (
	"jbrodriguez/mediagui/server/src/lib"
	"jbrodriguez/mediagui/server/src/model"
	"jbrodriguez/mediagui/server/src/proto"

	"database/sql"
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jbrodriguez/actor"
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	_ "github.com/mattn/go-sqlite3" // sqlite3 doesn't need to be named
)

// const DATETIME_LAYOUT = "2006-01-02T15:04:05-07:00"

// Dal -
type Dal struct {
	bus      *pubsub.PubSub
	settings *lib.Settings

	actor *actor.Actor

	db *sql.DB
	// dbase       string
	// err         error
	// count       uint64
	// searchCount uint64
	// searchArgs  string

	countRows *sql.Stmt
	// storeMovie *sql.Stmt

	// listMovies      *sql.Stmt
	// listByRuntime   *sql.Stmt
	// listMoviesToFix *sql.Stmt
}

// NewDal -
func NewDal(bus *pubsub.PubSub, settings *lib.Settings) *Dal {
	dal := &Dal{
		bus:      bus,
		settings: settings,
		actor:    actor.NewActor(bus),
	}
	return dal
}

// Start -
func (d *Dal) Start() {
	mlog.Info("Starting service Dal ...")

	dbPath := filepath.Join(d.settings.DataDir, "mediagui.db")

	var err error
	d.db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		mlog.Fatalf("Unable to open database (%s): %s", dbPath, err)
	}

	d.actor.Register("/get/movies/cover", d.getCover)
	d.actor.Register("/get/movies", d.getMovies)
	d.actor.Register("/get/movies/duplicates", d.getDuplicates)
	d.actor.Register("/get/movie", d.getMovie)
	d.actor.Register("/command/movie/exists", d.checkExists)
	d.actor.Register("/command/movie/store", d.storeMovie)
	d.actor.Register("/command/movie/partialstore", d.partialStoreMovie)
	d.actor.Register("/command/movie/update", d.updateMovie)
	d.actor.Register("/command/movie/delete", d.deleteMovie)
	d.actor.Register("/put/movies/score", d.setScore)
	d.actor.Register("/put/movies/watched", d.setWatched)
	d.actor.Register("/put/movies/duplicate", d.setDuplicate)

	d.countRows = d.prepare("select count(*) from movie;")

	mlog.Info("Connected to database %s", dbPath)

	go d.actor.React()
}

// Stop -
func (d *Dal) Stop() {
	mlog.Info("Stopped service Dal ...")
}

func (d *Dal) getCover(msg *pubsub.Message) {
	options := lib.Options{
		Limit:     60,
		SortBy:    "added",
		SortOrder: "desc",
	}

	total, items := d.listMovies(&options)

	msg.Reply <- &model.MoviesDTO{Total: total, Items: items}
}

func (d *Dal) getMovies(msg *pubsub.Message) {
	options := msg.Payload.(*lib.Options)

	mlog.Info("options: %+v", options)

	var total uint64
	var items []*model.Movie

	if options.Query == "" {
		total, items = d.listMovies(options)
	} else {
		total, items = d.searchMovies(options)
	}

	msg.Reply <- &model.MoviesDTO{Total: total, Items: items}
}

func (d *Dal) listMovies(options *lib.Options) (total uint64, movies []*model.Movie) {
	// mlog.Info("listMovies.options: %+v", options)

	tx, err := d.db.Begin()
	if err != nil {
		mlog.Fatalf("listMovies:Unable to begin transaction: %s", err)
	}

	sql := fmt.Sprintf(`select rowid, title, original_title, file_title, year, runtime, tmdb_id, imdb_id,
				overview, tagline, resolution, filetype, location, cover, backdrop, genres, vote_average,
				vote_count, countries, added, modified, last_watched, all_watched, count_watched, score,
				director, writer, actors, awards, imdb_rating, imdb_votes, show_if_duplicate, stub
				from movie order by %s %s, rowid desc limit ? offset ?`, options.SortBy, options.SortOrder)

	stmt, err := tx.Prepare(sql)
	if err != nil {
		mlog.Fatalf("Unable to prepare transaction: %s", err)
	}
	defer lib.Close(stmt)

	rows, err := stmt.Query(options.Limit, options.Offset)
	if err != nil {
		mlog.Fatalf("Unable to prepare transaction: %s", err)
	}

	items := make([]*model.Movie, 0)

	// if options.Offset == 0 {
	// if d.count == 0 {
	var count uint64

	err = d.countRows.QueryRow().Scan(&count)
	if err != nil {
		mlog.Fatalf("Unable to count rows: %s", err)
	}
	// }

	for rows.Next() {
		movie := model.Movie{}
		rows.Scan(
			&movie.ID, &movie.Title, &movie.Original_Title, &movie.File_Title, &movie.Year,
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

	mlog.Info("Listed %d movies (total %d)", len(items), count)

	return count, items
}

func (d *Dal) searchMovies(options *lib.Options) (total uint64, movies []*model.Movie) {
	if options.FilterBy == "year" {
		total, movies = d.searchByYear(options)
	} else {
		total, movies = d.regularSearch(options)
	}

	return total, movies
}

func (d *Dal) searchByYear(options *lib.Options) (total uint64, movies []*model.Movie) {
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

	// mlog.Info("this is: %s %s", term, args)

	tx, err := d.db.Begin()
	if err != nil {
		mlog.Fatalf("searchByYear:Unable to begin transaction: %s", err)
	}

	stmt, err := tx.Prepare(countQuery)
	if err != nil {
		mlog.Fatalf("searchByYear:Unable to prepare transaction: %s", err)
	}
	defer lib.Close(stmt)

	var count uint64

	var rows *sql.Rows
	if decade {
		err = stmt.QueryRow(start, end).Scan(&count)
		if err != nil {
			mlog.Fatalf("searchByYear:Unable to count rows: %s", err)
		}

		sql := fmt.Sprintf(`select rowid, title, original_title, file_title, year, runtime, tmdb_id, imdb_id,
					overview, tagline, resolution, filetype, location, cover, backdrop, genres, vote_average,
					vote_count, countries, added, modified, last_watched, all_watched, count_watched, score,
					director, writer, actors, awards, imdb_rating, imdb_votes, show_if_duplicate, stub
					from movie where year between ? and ? order by %s %s, rowid desc limit ? offset ?`, options.SortBy, options.SortOrder)

		stmt, err = tx.Prepare(sql)
		if err != nil {
			mlog.Fatalf("Unable to prepare transaction: %s", err)
		}
		defer lib.Close(stmt)

		rows, err = stmt.Query(start, end, options.Limit, options.Offset)
		if err != nil {
			mlog.Fatalf("Unable to prepare transaction: %s", err)
		}
	} else {
		err = stmt.QueryRow(year).Scan(&count)
		if err != nil {
			mlog.Fatalf("searchByYear:Unable to count rows: %s", err)
		}

		sql := fmt.Sprintf(`select rowid, title, original_title, file_title, year, runtime, tmdb_id, imdb_id,
					overview, tagline, resolution, filetype, location, cover, backdrop, genres, vote_average,
					vote_count, countries, added, modified, last_watched, all_watched, count_watched, score,
					director, writer, actors, awards, imdb_rating, imdb_votes, show_if_duplicate, stub
					from movie where year = ? order by %s %s, rowid desc limit ? offset ?`, options.SortBy, options.SortOrder)

		mlog.Info("sql.: %s", sql)

		stmt, err := tx.Prepare(sql)
		if err != nil {
			mlog.Fatalf("Unable to prepare transaction: %s", err)
		}
		defer lib.Close(stmt)

		rows, err = stmt.Query(year, options.Limit, options.Offset)
		if err != nil {
			mlog.Fatalf("Unable to prepare transaction: %s", err)
		}
	}

	items := make([]*model.Movie, 0)

	for rows.Next() {
		movie := model.Movie{}
		rows.Scan(
			&movie.ID, &movie.Title, &movie.Original_Title, &movie.File_Title, &movie.Year,
			&movie.Runtime, &movie.Tmdb_Id, &movie.Imdb_Id, &movie.Overview, &movie.Tagline,
			&movie.Resolution, &movie.FileType, &movie.Location, &movie.Cover, &movie.Backdrop,
			&movie.Genres, &movie.Vote_Average, &movie.Vote_Count, &movie.Production_Countries,
			&movie.Added, &movie.Modified, &movie.Last_Watched, &movie.All_Watched,
			&movie.Count_Watched, &movie.Score, &movie.Director, &movie.Writer, &movie.Actors,
			&movie.Awards, &movie.Imdb_Rating, &movie.Imdb_Votes, &movie.ShowIfDuplicate, &movie.Stub,
		)
		// movie := &model.Movie{}
		// rows.Scan(movie.ID, movie.Title, movie.Original_Title, movie.Year, movie.Runtime, movie.Tmdb_Id, movie.Imdb_Id, movie.Overview, movie.Tagline, movie.Resolution, movie.FileType, movie.Location, movie.Cover, movie.Backdrop)
		// mlog.Info("title: (%s)", movie.Title)
		items = append(items, &movie)
	}
	lib.Close(rows)

	commit(tx)

	mlog.Info("searchByYear:Listed %d movies (total %d)", len(items), count)

	return count, items
}

func (d *Dal) regularSearch(options *lib.Options) (total uint64, movies []*model.Movie) {
	mlog.Info("searchMovies.options: %+v", options)

	tx, err := d.db.Begin()
	if err != nil {
		mlog.Fatalf("searchMovies:Unable to begin transaction: %s", err)
	}

	term := options.Query + "*"
	args := options.FilterBy

	mlog.Info("this is: %s %s", term, args)

	// if self.searchArgs != args {
	// self.searchArgs = args

	countQuery := fmt.Sprintf(`select count(*) from movie dt, %s vt
		where vt.%s match ? and dt.rowid = vt.docid;`,
		"movie"+options.FilterBy, "movie"+options.FilterBy)

	stmt, err := tx.Prepare(countQuery)
	if err != nil {
		mlog.Fatalf("searchMovies:Unable to prepare transaction: %s", err)
	}
	defer lib.Close(stmt)

	mlog.Info("searchMovies:CountQuery:%s", countQuery)

	var count uint64
	err = stmt.QueryRow(term).Scan(&count)
	if err != nil {
		mlog.Fatalf("searchMovies:Unable to count rows: %s", err)
	}

	listQuery := fmt.Sprintf(`select dt.rowid, dt.title, dt.original_title, dt.year, dt.runtime,
			dt.tmdb_id, dt.imdb_id, dt.overview, dt.tagline, dt.resolution,
			dt.filetype, dt.location, dt.cover, dt.backdrop, dt.genres, dt.vote_average,
			dt.vote_count, dt.countries, dt.added, dt.modified, dt.last_watched,
			dt.all_watched, dt.count_watched, dt.score, dt.director, dt.writer, dt.actors,
			dt.awards, dt.imdb_rating, dt.imdb_votes, dt.show_if_duplicate, dt.stub
			from movie dt, %s vt
			where vt.%s match ? and dt.rowid = vt.docid order by dt.%s %s limit ? offset ?`,
		"movie"+options.FilterBy, "movie"+options.FilterBy, options.SortBy, options.SortOrder)

	// mlog.Info("my main man: %s", sql)

	stmt, err = tx.Prepare(listQuery)
	if err != nil {
		mlog.Fatalf("searchMovies:listQuery:Unable to prepare transaction: %s", err)
	}
	defer lib.Close(stmt)

	rows, err := stmt.Query(term, options.Limit, options.Offset)
	if err != nil {
		mlog.Fatalf("searchMovies:listQuery:Unable to begin transaction: %s", err)
	}

	items := make([]*model.Movie, 0)

	for rows.Next() {
		movie := model.Movie{}
		rows.Scan(&movie.ID, &movie.Title, &movie.Original_Title, &movie.Year, &movie.Runtime, &movie.Tmdb_Id, &movie.Imdb_Id, &movie.Overview, &movie.Tagline, &movie.Resolution, &movie.FileType, &movie.Location, &movie.Cover, &movie.Backdrop, &movie.Genres, &movie.Vote_Average, &movie.Vote_Count, &movie.Production_Countries, &movie.Added, &movie.Modified, &movie.Last_Watched, &movie.All_Watched, &movie.Count_Watched, &movie.Score, &movie.Director, &movie.Writer, &movie.Actors, &movie.Awards, &movie.Imdb_Rating, &movie.Imdb_Votes, &movie.ShowIfDuplicate, &movie.Stub)
		// movie := &model.Movie{}
		// rows.Scan(movie.ID, movie.Title, movie.Original_Title, movie.Year, movie.Runtime, movie.Tmdb_Id, movie.Imdb_Id, movie.Overview, movie.Tagline, movie.Resolution, movie.FileType, movie.Location, movie.Cover, movie.Backdrop)
		// mlog.Info("title: (%s)", movie.Title)
		items = append(items, &movie)
	}
	lib.Close(rows)

	commit(tx)

	mlog.Info("searchMovies:Listed %d movies (total %d)", len(items), count)

	return count, items
}

func (d *Dal) getDuplicates(msg *pubsub.Message) {
	mlog.Info("getDuplicates.starting")

	tx, err := d.db.Begin()
	if err != nil {
		mlog.Fatalf("Unable to begin transaction: %s", err)
	}

	// rows, err := self.listMovies.Query()
	// if err != nil {
	// 	mlog.Fatalf("unable to prepare transaction: %s", err)
	// }

	// rows, err := self.db.Query("select rowid, title, original_title, file_title, year, runtime, tmdb_id, imdb_id, overview, tagline, resolution, filetype, location, cover, backdrop, genres, vote_average, vote_count, countries, added, modified, last_watched, all_watched, count_watched from movie where title in (select title from movie group by title having count(*) > 1);")
	rows, err := d.db.Query(`select a.rowid, a.title, a.original_title, a.file_title,
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
		mlog.Fatalf("Unable to prepare transaction: %s", err)
	}

	items := make([]*model.Movie, 0)

	for rows.Next() {
		movie := model.Movie{}
		rows.Scan(&movie.ID, &movie.Title, &movie.Original_Title, &movie.File_Title, &movie.Year, &movie.Runtime, &movie.Tmdb_Id, &movie.Imdb_Id, &movie.Overview, &movie.Tagline, &movie.Resolution, &movie.FileType, &movie.Location, &movie.Cover, &movie.Backdrop, &movie.Genres, &movie.Vote_Average, &movie.Vote_Count, &movie.Production_Countries, &movie.Added, &movie.Modified, &movie.Last_Watched, &movie.All_Watched, &movie.Count_Watched, &movie.Score, &movie.Director, &movie.Writer, &movie.Actors, &movie.Awards, &movie.Imdb_Rating, &movie.Imdb_Votes, &movie.ShowIfDuplicate, &movie.Stub)
		items = append(items, &movie)
	}
	lib.Close(rows)

	commit(tx)

	mlog.Info("Found %d duplicate movies", len(items))

	msg.Reply <- &model.MoviesDTO{Total: uint64(len(items)), Items: items}
}

func (d *Dal) getMovie(msg *pubsub.Message) {
	id := msg.Payload.(string)

	sql := `select rowid, title, original_title, file_title, year, runtime, tmdb_id, imdb_id,
				overview, tagline, resolution, filetype, location, cover, backdrop, genres, vote_average,
				vote_count, countries, added, modified, last_watched, all_watched, count_watched, score,
				director, writer, actors, awards, imdb_rating, imdb_votes, show_if_duplicate, stub
				from movie where rowid = ?`

	stmt, err := d.db.Prepare(sql)
	if err != nil {
		mlog.Fatalf("Unable to prepare transaction: %s", err)
	}
	defer lib.Close(stmt)

	movie := model.Movie{}

	err = stmt.
		QueryRow(id).
		Scan(&movie.ID, &movie.Title, &movie.Original_Title, &movie.File_Title, &movie.Year, &movie.Runtime, &movie.Tmdb_Id, &movie.Imdb_Id, &movie.Overview, &movie.Tagline, &movie.Resolution, &movie.FileType, &movie.Location, &movie.Cover, &movie.Backdrop, &movie.Genres, &movie.Vote_Average, &movie.Vote_Count, &movie.Production_Countries, &movie.Added, &movie.Modified, &movie.Last_Watched, &movie.All_Watched, &movie.Count_Watched, &movie.Score, &movie.Director, &movie.Writer, &movie.Actors, &movie.Awards, &movie.Imdb_Rating, &movie.Imdb_Votes, &movie.ShowIfDuplicate, &movie.Stub)
	if err != nil {
		mlog.Fatalf("Unable to prepare transaction: %s", err)
	}

	msg.Reply <- &movie
}

func (d *Dal) checkExists(msg *pubsub.Message) {
	tx, err := d.db.Begin()
	if err != nil {
		mlog.Fatalf("at begin: %s", err)
	}

	stmt, err := tx.Prepare("select rowid from movie where location = ?")
	if err != nil {
		rollback(tx)
		mlog.Fatalf("at prepare: %s", err)
	}
	defer lib.Close(stmt)

	movie := msg.Payload.(*model.Movie)

	var id int
	err = stmt.QueryRow(movie.Location).Scan(&id)

	// if err == sql.ErrNoRows {
	// 	mlog.Fatalf("id = %d, err = %d", id, err)
	// }

	// mlog.Fatalf("gone and done")
	if err != sql.ErrNoRows && err != nil {
		rollback(tx)
		mlog.Fatalf("at queryrow: %s", err)
	}

	commit(tx)

	// mlog.Info("Check exists: [%d] (%s)", id, movie.Location)

	msg.Reply <- (id != 0)
}

func (d *Dal) storeMovie(msg *pubsub.Message) {
	defer d.bus.Pub(nil, "/event/workunit/done")

	movie := msg.Payload.(*model.Movie)

	// d.count = 0

	mlog.Info("STARTED SAVING %s", movie.Title)

	tx, err := d.db.Begin()
	if err != nil {
		mlog.Fatalf("at begin: %s", err)
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
		mlog.Fatalf("at prepare: %s", err)
	}
	defer lib.Close(stmt)

	res, e := stmt.Exec(movie.Title, movie.Original_Title, movie.File_Title, movie.Year,
		movie.Runtime, movie.Tmdb_Id, movie.Imdb_Id, movie.Overview, movie.Tagline,
		movie.Resolution, movie.FileType, movie.Location, movie.Cover, movie.Backdrop,
		movie.Genres, movie.Vote_Average, movie.Vote_Count, movie.Production_Countries,
		movie.Added, movie.Modified, movie.Last_Watched, movie.All_Watched, movie.Count_Watched,
		movie.Score, movie.Director, movie.Writer, movie.Actors, movie.Awards, movie.Imdb_Rating,
		movie.Imdb_Votes, movie.ShowIfDuplicate, movie.Stub)
	if e != nil {
		rollback(tx)
		mlog.Fatalf("at exec: %s", e)
	}

	id, _ := res.LastInsertId()
	movie.ID = uint64(id)

	commit(tx)
	mlog.Info("FINISHED SAVING %s [%d]", movie.Title, movie.ID)
}

func (d *Dal) partialStoreMovie(msg *pubsub.Message) {
	defer d.bus.Pub(nil, "/event/workunit/done")

	movie := msg.Payload.(*model.Movie)

	// d.count = 0
	now := time.Now().UTC().Format(time.RFC3339)
	movie.Added = now
	movie.Modified = now
	movie.ShowIfDuplicate = 1

	mlog.Info("STARTED PARTIAL SAVING %s", movie.Title)

	tx, err := d.db.Begin()
	if err != nil {
		mlog.Fatalf("at begin: %s", err)
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
		mlog.Fatalf("at prepare: %s", err)
	}
	defer lib.Close(stmt)

	res, e := stmt.Exec(movie.Title, movie.Original_Title, movie.File_Title, movie.Year,
		movie.Runtime, movie.Tmdb_Id, movie.Imdb_Id, movie.Overview, movie.Tagline,
		movie.Resolution, movie.FileType, movie.Location, movie.Cover, movie.Backdrop,
		movie.Genres, movie.Vote_Average, movie.Vote_Count, movie.Production_Countries,
		movie.Added, movie.Modified, movie.Last_Watched, movie.All_Watched, movie.Count_Watched,
		movie.Score, movie.Director, movie.Writer, movie.Actors, movie.Awards, movie.Imdb_Rating,
		movie.Imdb_Votes, movie.ShowIfDuplicate, movie.Stub)
	if e != nil {
		rollback(tx)
		mlog.Fatalf("at exec: %s", e)
	}

	id, _ := res.LastInsertId()
	movie.ID = uint64(id)

	commit(tx)
	mlog.Info("FINISHED PARTIAL SAVING %s [%d]", movie.Title, movie.ID)
}

func (d *Dal) updateMovie(msg *pubsub.Message) {
	defer d.bus.Pub(nil, "/event/workunit/done")

	movie := msg.Payload.(*model.Movie)

	mlog.Info("STARTED UPDATING [%d] %s", movie.ID, movie.Title)

	tx, err := d.db.Begin()
	if err != nil {
		mlog.Fatalf("at begin: %s", err)
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
		mlog.Fatalf("at prepare: %s", err)
	}
	defer lib.Close(stmt)

	_, err = stmt.Exec(movie.Title, movie.Original_Title, movie.Year, movie.Runtime, movie.Tmdb_Id, movie.Imdb_Id, movie.Overview, movie.Tagline, movie.Cover, movie.Backdrop, movie.Genres, movie.Vote_Average, movie.Vote_Count, movie.Production_Countries, movie.Modified, movie.Director, movie.Writer, movie.Actors, movie.Awards, movie.Imdb_Rating, movie.Imdb_Votes, movie.ShowIfDuplicate, movie.ID)
	if err != nil {
		rollback(tx)
		mlog.Fatalf("at exec: %s", err)
	}

	commit(tx)
	mlog.Info("FINISHED UPDATING [%d] %s", movie.ID, movie.Title)

	// updated := &pubsub.Message{}
	// d.bus.Pub(updated, "/event/movie/updated")
}

func (d *Dal) deleteMovie(msg *pubsub.Message) {
	movie := msg.Payload.(*agent.Item)

	// d.count = 0

	lib.Notify(d.bus, "prune:delete", fmt.Sprintf("STARTED DELETING [%d] %s", movie.Id, movie.Title))

	tx, err := d.db.Begin()
	if err != nil {
		mlog.Fatalf("at begin: %s", err)
	}

	stmt, err := tx.Prepare("delete from movie where rowid = ?")
	if err != nil {
		rollback(tx)
		mlog.Fatalf("at prepare: %s", err)
	}
	defer lib.Close(stmt)

	_, err = stmt.Exec(movie.Id)
	if err != nil {
		rollback(tx)
		mlog.Fatalf("at exec: %s", err)
	}

	commit(tx)

	lib.Notify(d.bus, "prune:delete", fmt.Sprintf("FINISHED DELETING [%d] %s", movie.Id, movie.Title))
}

func (d *Dal) setScore(msg *pubsub.Message) {
	movie := msg.Payload.(*model.Movie)

	mlog.Info("STARTED UPDATING MOVIE SCORE [%d] %s (%d)", movie.ID, movie.Title, movie.Score)

	tx, err := d.db.Begin()
	if err != nil {
		mlog.Fatalf("at begin: %s", err)
	}

	stmt, err := tx.Prepare(`update movie set
								score = ?,
								modified = ?
								where rowid = ?`)
	if err != nil {
		rollback(tx)
		mlog.Fatalf("at prepare: %s", err)
	}
	defer lib.Close(stmt)

	now := time.Now().UTC().Format(time.RFC3339)

	_, err = stmt.Exec(movie.Score, now, movie.ID)
	if err != nil {
		rollback(tx)
		mlog.Fatalf("at exec: %s", err)
	}

	commit(tx)
	mlog.Info("FINISHED UPDATING MOVIE SCORE [%d] %s", movie.ID, movie.Title)

	msg.Reply <- movie
}

func (d *Dal) setWatched(msg *pubsub.Message) {
	dto := msg.Payload.(*model.Movie)

	mlog.Info("STARTED UPDATING MOVIE WATCHED DATE [%d] %s (%s)", dto.ID, dto.Title, dto.Last_Watched)

	now := time.Now().UTC().Format(time.RFC3339)

	stmt, err := d.db.Prepare("select all_watched from movie where rowid = ?")
	if err != nil {
		mlog.Fatalf("at prepare: %s", err)
	}

	// get all watched times from the db
	var when string
	err = stmt.QueryRow(dto.ID).Scan(&when)
	if err != nil {
		mlog.Fatalf("at queryrow: %s", err)
	}

	// create an array with all watched times
	var watchedTimes []string
	if when != "" {
		watchedTimes = strings.Split(when, "|")
	}

	// convert incoming watched time to sane format
	watched, err := parseToday(dto.Last_Watched)
	if err != nil {
		mlog.Fatalf("at parseToday: %s", err)
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

	// var allWatched string
	// countWatched := dto.Count_Watched
	// if !strings.Contains(dto.All_Watched, dto.Last_Watched) {
	// 	countWatched++
	// 	if dto.All_Watched == "" {
	// 		allWatched = dto.Last_Watched
	// 	} else {
	// 		allWatched += "|" + dto.Last_Watched
	// 	}
	// }

	tx, err := d.db.Begin()
	if err != nil {
		mlog.Fatalf("at begin: %s", err)
	}

	stmt, err = tx.Prepare(`update movie set
								last_watched = ?,
								all_watched = ?,
								count_watched = ?,
								modified = ?
								where rowid = ?`)
	if err != nil {
		rollback(tx)
		mlog.Fatalf("at prepare: %s", err)
	}
	defer lib.Close(stmt)

	_, err = stmt.Exec(lastWatched, allWatched, countWatched, now, dto.ID)
	if err != nil {
		rollback(tx)
		mlog.Fatalf("at exec: %s", err)
	}

	commit(tx)
	mlog.Info("FINISHED UPDATING MOVIE WATCHED DATE [%d] %s", dto.ID, dto.Title)

	dto.All_Watched = allWatched
	dto.Count_Watched = countWatched
	dto.Modified = now

	msg.Reply <- dto
}

func (d *Dal) setDuplicate(msg *pubsub.Message) {
	movie := msg.Payload.(*model.Movie)

	mlog.Info("STARTED UPDATING MOVIE DUPLICATE STATUS [%d] %s (%d)", movie.ID, movie.Title, movie.ShowIfDuplicate)

	tx, err := d.db.Begin()
	if err != nil {
		mlog.Fatalf("at begin: %s", err)
	}

	stmt, err := tx.Prepare(`update movie set
								show_if_duplicate = ?,
								modified = ?
								where rowid = ?`)
	if err != nil {
		rollback(tx)
		mlog.Fatalf("at prepare: %s", err)
	}
	defer lib.Close(stmt)

	now := time.Now().UTC().Format(time.RFC3339)

	_, err = stmt.Exec(movie.ShowIfDuplicate, now, movie.ID)
	if err != nil {
		rollback(tx)
		mlog.Fatalf("at exec: %s", err)
	}

	commit(tx)
	mlog.Info("FINISHED UPDATING MOVIE DUPLICATE STATUS [%d] %s", movie.ID, movie.Title)

	msg.Reply <- movie
}

func (d *Dal) prepare(sql string) *sql.Stmt {
	stmt, err := d.db.Prepare(sql)
	if err != nil {
		mlog.Fatalf("prepare sql: %s (%s)", err, sql)
	}
	return stmt
}

func parseToday(clientToday string) (today time.Time, err error) {
	client, perr := time.Parse(time.RFC3339, clientToday)
	if perr != nil {
		return today, perr
	}

	today = time.Date(client.Year(), client.Month(), client.Day(), 0, 0, 0, 0, client.Location())

	return today, nil
}

func rollback(tx *sql.Tx) {
	err := tx.Rollback()
	if err != nil {
		mlog.Warning("ROLLBACK: %s", err)
	}
}

func commit(tx *sql.Tx) {
	err := tx.Commit()
	if err != nil {
		mlog.Warning("COMMIT: %s", err)
	}
}
