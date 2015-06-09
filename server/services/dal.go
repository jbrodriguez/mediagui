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
	// "strings"
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

	d.mailbox = d.register(d.bus, d.getCover, "/get/movies/cover")

	d.countRows = d.prepare("select count(*) from movie;")

	mlog.Info("Connected to database %s", dbPath)

	go d.react()
}

func (d *Dal) Stop() {
	mlog.Info("Stopped service Dal ...")
}

func (d *Dal) react() {
	for mbox := range d.mailbox {
		mlog.Info("DAL:Topic: %s", mbox.Topic)
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
		mlog.Fatalf("Unable to begin transaction: %s", err)
	}

	sql := fmt.Sprintf(`select rowid, title, original_title, file_title, year, runtime, tmdb_id, imdb_id, 
				overview, tagline, resolution, filetype, location, cover, backdrop, genres, vote_average, 
				vote_count, countries, added, modified, last_watched, all_watched, count_watched, score, 
				director, writer, actors, awards, imdb_rating, imdb_votes 
				from movie order by %s %s limit ? offset ?`, options.SortBy, options.SortOrder)

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

	var count = 0
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
		count++
	}
	rows.Close()

	tx.Commit()

	mlog.Info("Listed %d movies (total %d)", count, d.count)

	return d.count, items
}

func (d *Dal) prepare(sql string) *sql.Stmt {
	stmt, err := d.db.Prepare(sql)
	if err != nil {
		mlog.Fatalf("prepare sql: %s (%s)", err, sql)
	}
	return stmt
}
