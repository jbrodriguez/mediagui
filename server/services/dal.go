package services

import (
	"database/sql"
	"fmt"
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	_ "github.com/mattn/go-sqlite3"
	"jbrodriguez/mediagui/server/lib"
	"jbrodriguez/mediagui/server/model"
	"path/filepath"
	"strings"
	// "time"
)

type Dal struct {
	Service

	bus      *pubsub.PubSub
	settings *lib.Settings

	mailbox chan *pubsub.Mailbox

	db *sql.DB
	// dbase       string
	// err         error
	count       uint64
	searchCount uint64
	searchArgs  string

	countRows  *sql.Stmt
	storeMovie *sql.Stmt

	// listMovies      *sql.Stmt
	// listByRuntime   *sql.Stmt
	// listMoviesToFix *sql.Stmt
}

func NewDal(bus *pubsub.PubSub, settings *lib.Settings) *Dal {
	dal := &Dal{bus: bus, settings: settings}
	dal.init()
	return dal
}

func (d *Dal) Start() {
	mlog.Info("Starting service Dal ...")

	dbPath := filepath.Join(d.settings.DataDir, "mediagui.db")

	var err error
	d.db, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		mlog.Fatalf("Unable to open database (%s): %s", dbPath, err)
	}

	d.mailbox = d.register(d.bus, "/get/movies/cover", d.getCover)
	d.registerAdditional(d.bus, "/get/movies", d.getMovies, d.mailbox)
	d.registerAdditional(d.bus, "/command/movie/exists", d.checkExists, d.mailbox)
	d.registerAdditional(d.bus, "/put/movies/score", d.setScore, d.mailbox)

	d.countRows = d.prepare("select count(*) from movie;")

	mlog.Info("Connected to database %s", dbPath)

	go d.react()
}

func (d *Dal) Stop() {
	mlog.Info("Stopped service Dal ...")
}

func (d *Dal) react() {
	for mbox := range d.mailbox {
		// mlog.Info("DAL:Topic: %s", mbox.Topic)
		d.dispatch(mbox.Topic, mbox.Content)
	}
}

func (d *Dal) getCover(msg *pubsub.Message) {
	options := lib.Options{
		Limit:     60,
		SortBy:    "added",
		SortOrder: "desc",
	}

	total, items := d.listMovies(options)

	msg.Reply <- &model.MoviesDTO{Total: total, Items: items}
}

func (d *Dal) getMovies(msg *pubsub.Message) {
	options := msg.Payload.(lib.Options)

	var total uint64
	var items []*model.Movie

	if options.Query == "" {
		total, items = d.listMovies(options)
	} else {
		total, items = d.searchMovies(options)
	}

	msg.Reply <- &model.MoviesDTO{Total: total, Items: items}
}

// type Options struct {
// 	SearchTerm string `json:"searchTerm"`
// 	Current    uint64 `json:"current"`
// 	Limit      uint64 `json:"limit"`
// 	SortBy     string `json:"sortBy"`
// 	SortOrder  string `json:"sortOrder"`
// 	FilterBy   string `json:"filterBy"`
// }
func (d *Dal) listMovies(options lib.Options) (total uint64, movies []*model.Movie) {
	mlog.Info("listMovies.options: %+v", options)

	tx, err := d.db.Begin()
	if err != nil {
		mlog.Fatalf("listMovies:Unable to begin transaction: %s", err)
	}

	sql := fmt.Sprintf(`select rowid, title, original_title, file_title, year, runtime, tmdb_id, imdb_id, 
				overview, tagline, resolution, filetype, location, cover, backdrop, genres, vote_average, 
				vote_count, countries, added, modified, last_watched, all_watched, count_watched, score, 
				director, writer, actors, awards, imdb_rating, imdb_votes 
				from movie order by %s %s, rowid desc limit ? offset ?`, options.SortBy, options.SortOrder)

	stmt, err := tx.Prepare(sql)
	if err != nil {
		mlog.Fatalf("Unable to prepare transaction: %s", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(options.Limit, options.Offset)
	if err != nil {
		mlog.Fatalf("Unable to prepare transaction: %s", err)
	}

	items := make([]*model.Movie, 0)

	// if options.Offset == 0 {
	if d.count == 0 {
		err = d.countRows.QueryRow().Scan(&d.count)
		if err != nil {
			mlog.Fatalf("Unable to count rows: %s", err)
		}
	}

	for rows.Next() {
		movie := model.Movie{}
		rows.Scan(
			&movie.Id, &movie.Title, &movie.Original_Title, &movie.File_Title, &movie.Year,
			&movie.Runtime, &movie.Tmdb_Id, &movie.Imdb_Id, &movie.Overview, &movie.Tagline,
			&movie.Resolution, &movie.FileType, &movie.Location, &movie.Cover, &movie.Backdrop,
			&movie.Genres, &movie.Vote_Average, &movie.Vote_Count, &movie.Production_Countries,
			&movie.Added, &movie.Modified, &movie.Last_Watched, &movie.All_Watched,
			&movie.Count_Watched, &movie.Score, &movie.Director, &movie.Writer, &movie.Actors,
			&movie.Awards, &movie.Imdb_Rating, &movie.Imdb_Votes,
		)
		items = append(items, &movie)
	}
	rows.Close()

	tx.Commit()

	mlog.Info("Listed %d movies (total %d)", len(items), d.count)

	return d.count, items
}

func (d *Dal) searchMovies(options lib.Options) (total uint64, movies []*model.Movie) {
	mlog.Info("searchMovies.options: %+v", options)

	tx, err := d.db.Begin()
	if err != nil {
		mlog.Fatalf("searchMovies:Unable to begin transaction: %s", err)
	}

	term := options.Query + "*"
	args := options.FilterBy

	mlog.Info("this is: %s %s %s", term, d.searchArgs, args)

	// if self.searchArgs != args {
	// self.searchArgs = args

	countQuery := fmt.Sprintf(`select count(*) from movie dt, %s vt 
		where vt.%s match ? and dt.rowid = vt.docid;`,
		"movie"+options.FilterBy, "movie"+options.FilterBy)

	stmt, err := tx.Prepare(countQuery)
	if err != nil {
		mlog.Fatalf("searchMovies:Unable to prepare transaction: %s", err)
	}
	defer stmt.Close()

	mlog.Info("searchMovies:CountQuery:%s", countQuery)

	err = stmt.QueryRow(term).Scan(&d.searchCount)
	if err != nil {
		mlog.Fatalf("searchMovies:Unable to count rows: %s", err)
	}

	listQuery := fmt.Sprintf(`select dt.rowid, dt.title, dt.original_title, dt.year, dt.runtime, 
			dt.tmdb_id, dt.imdb_id, dt.overview, dt.tagline, dt.resolution, 
			dt.filetype, dt.location, dt.cover, dt.backdrop, dt.genres, dt.vote_average, 
			dt.vote_count, dt.countries, dt.added, dt.modified, dt.last_watched, 
			dt.all_watched, dt.count_watched, dt.score, dt.director, dt.writer, dt.actors, 
			dt.awards, dt.imdb_rating, dt.imdb_votes
			from movie dt, %s vt 
			where vt.%s match ? and dt.rowid = vt.docid order by dt.%s %s limit ? offset ?`,
		"movie"+options.FilterBy, "movie"+options.FilterBy, options.SortBy, options.SortOrder)

	// mlog.Info("my main man: %s", sql)

	stmt, err = tx.Prepare(listQuery)
	if err != nil {
		mlog.Fatalf("searchMovies:listQuery:Unable to prepare transaction: %s", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(term, options.Limit, options.Offset)
	if err != nil {
		mlog.Fatalf("searchMovies:listQuery:Unable to begin transaction: %s", err)
	}

	items := make([]*model.Movie, 0)

	for rows.Next() {
		movie := model.Movie{}
		rows.Scan(&movie.Id, &movie.Title, &movie.Original_Title, &movie.Year, &movie.Runtime, &movie.Tmdb_Id, &movie.Imdb_Id, &movie.Overview, &movie.Tagline, &movie.Resolution, &movie.FileType, &movie.Location, &movie.Cover, &movie.Backdrop, &movie.Genres, &movie.Vote_Average, &movie.Vote_Count, &movie.Production_Countries, &movie.Added, &movie.Modified, &movie.Last_Watched, &movie.All_Watched, &movie.Count_Watched, &movie.Score, &movie.Director, &movie.Writer, &movie.Actors, &movie.Awards, &movie.Imdb_Rating, &movie.Imdb_Votes)
		// movie := &model.Movie{}
		// rows.Scan(movie.Id, movie.Title, movie.Original_Title, movie.Year, movie.Runtime, movie.Tmdb_Id, movie.Imdb_Id, movie.Overview, movie.Tagline, movie.Resolution, movie.FileType, movie.Location, movie.Cover, movie.Backdrop)
		// mlog.Info("title: (%s)", movie.Title)
		items = append(items, &movie)
	}
	rows.Close()

	tx.Commit()

	mlog.Info("searchMovies:Listed %d movies (total %d)", len(items), d.searchCount)

	return d.searchCount, items
}

func (d *Dal) checkExists(msg *pubsub.Message) {
	tx, err := d.db.Begin()
	if err != nil {
		mlog.Fatalf("at begin: %s", err)
	}

	stmt, err := tx.Prepare("select rowid from movie where upper(location) = ?")
	if err != nil {
		tx.Rollback()
		mlog.Fatalf("at prepare: %s", err)
	}
	defer stmt.Close()

	movie := msg.Payload.(*model.Movie)

	var id int
	err = stmt.QueryRow(strings.ToUpper(movie.Location)).Scan(&id)

	// if err == sql.ErrNoRows {
	// 	mlog.Fatalf("id = %d, err = %d", id, err)
	// }

	// mlog.Fatalf("gone and done")
	if err != sql.ErrNoRows && err != nil {
		tx.Rollback()
		mlog.Fatalf("at queryrow: %s", err)
	}

	tx.Commit()

	msg.Reply <- (id != 0)
}

func (d *Dal) setScore(msg *pubsub.Message) {
	dto := msg.Payload.(*model.Movie)

	mlog.Info("STARTED UPDATING MOVIE SCORE [%d] %s (%d)", dto.Id, dto.Title, dto.Score)

	tx, err := d.db.Begin()
	if err != nil {
		mlog.Fatalf("at begin: %s", err)
	}

	stmt, err := tx.Prepare("update movie set score = ? where rowid = ?")
	if err != nil {
		tx.Rollback()
		mlog.Fatalf("at prepare: %s", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(dto.Score, dto.Id)
	if err != nil {
		tx.Rollback()
		mlog.Fatalf("at exec: %s", err)
	}

	tx.Commit()
	mlog.Info("FINISHED UPDATING MOVIE SCORE [%d] %s", dto.Id, dto.Title)

	msg.Reply <- dto
}

func (d *Dal) prepare(sql string) *sql.Stmt {
	stmt, err := d.db.Prepare(sql)
	if err != nil {
		mlog.Fatalf("prepare sql: %s (%s)", err, sql)
	}
	return stmt
}
