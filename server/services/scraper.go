package services

import (
	"github.com/jbrodriguez/actor"
	"github.com/jbrodriguez/go-tmdb"
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"

	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"mediagui/dto"
	"mediagui/lib"
	"mediagui/model"
)

// Scraper -
type Scraper struct {
	bus      *pubsub.PubSub
	settings *lib.Settings
	pool     *lib.Pool
	tmdb     *tmdb.Tmdb

	actor *actor.Actor

	// reRating   *regexp.Regexp
	// reVotes    *regexp.Regexp
	// reDirector *regexp.Regexp
	// reWriter   *regexp.Regexp
	// reActor    *regexp.Regexp
}

// NewScraper -
func NewScraper(bus *pubsub.PubSub, settings *lib.Settings) *Scraper {
	scraper := &Scraper{
		bus:      bus,
		settings: settings,
		actor:    actor.NewActor(bus),
	}

	return scraper
}

// Start -
func (s *Scraper) Start() {
	mlog.Info("Starting service Scraper ...")

	var err error
	s.tmdb, err = tmdb.NewClient("e610ded10c3f47d05fe797961d90fea6", false)
	if err != nil {
		mlog.Fatalf("Unable to create tmdb client: %s", err)
	}

	s.actor.Register("/command/movie/scrape", s.scrapeMovie)
	s.actor.Register("/command/movie/rescrape", s.reScrapeMovie)
	s.actor.Register("/event/config/changed", s.configChanged)

	s.pool = lib.NewPool(12, 4000)

	go s.actor.React()
}

// Stop -
func (s *Scraper) Stop() {
	mlog.Info("Stopped service Scraper ...")
}

func (s *Scraper) scrapeMovie(msg *pubsub.Message) {
	movie := msg.Payload.(*model.Movie)

	scrape := &Scrape{
		s.bus,
		s.tmdb,
		&dto.Scrape{
			// BasePath: s.settings.WebDir,
			Movie:  movie,
			Forced: false,
		},
	}

	s.pool.Exec(scrape)
}

// Scrape -
type Scrape struct {
	bus  *pubsub.PubSub
	tmdb *tmdb.Tmdb
	dto  *dto.Scrape
}

// Execute -
func (s *Scrape) Execute(wid int) {
	movie := s.dto.Movie.(*model.Movie)

	lib.Notify(s.bus, "import:progress", fmt.Sprintf("SCRAPE REQUESTED (%d) [%s]", wid, movie.Title))

	now := time.Now().UTC().Format(time.RFC3339)
	movie.Added = now
	movie.Modified = now

	movie.Score = 0
	movie.ShowIfDuplicate = 1

	lib.Notify(s.bus, "import:progress", fmt.Sprintf("STARTED TMDB (%d) [%s]", wid, movie.Title))
	movies, err := s.tmdb.SearchMovie(movie.Title)
	if err != nil {
		s.bus.Pub(nil, "/event/workunit/done")

		mlog.Error(err)
		return
	}

	if movies.Total_Results == 0 {
		lib.Notify(s.bus, "import:progress", fmt.Sprintf("TMDB: NO MATCH FOUND (%d) [%s]", wid, movie.Title))

		msg := &pubsub.Message{Payload: s.dto}
		s.bus.Pub(msg, "/event/movie/tmdbnotfound")

		return
	} else if movies.Total_Results > 1 {
		lib.Notify(s.bus, "import:progress", fmt.Sprintf("TMDB: MORE THAN ONE (%d) [%s]", wid, movie.Title))
	}

	id := movies.Results[0].Id

	if err := _scrape(wid, s.tmdb, id, movie); err != nil {
		mlog.Warning("unable to scrape: %s", err)
	}

	s.dto.BaseURL = s.tmdb.BaseUrl
	s.dto.SecureBaseURL = s.tmdb.SecureBaseUrl

	lib.Notify(s.bus, "import:progress", fmt.Sprintf("SCRAPE COMPLETED (%d) [%s]", wid, movie.Title))

	msg := &pubsub.Message{Payload: s.dto}
	s.bus.Pub(msg, "/event/movie/scraped")
}

func (s *Scraper) reScrapeMovie(msg *pubsub.Message) {
	movie := msg.Payload.(*model.Movie)

	reScrape := &ReScrape{
		s.bus,
		s.tmdb,
		&dto.Scrape{
			// BasePath: s.settings.WebDir,
			Movie:  movie,
			Forced: true,
		},
	}

	s.pool.Exec(reScrape)
}

// ReScrape -
type ReScrape struct {
	bus  *pubsub.PubSub
	tmdb *tmdb.Tmdb
	dto  *dto.Scrape
}

// Execute -
func (s *ReScrape) Execute(wid int) {
	movie := s.dto.Movie.(*model.Movie)

	mlog.Info("RESCRAPE REQUESTED (%d) [%d] %s", wid, movie.ID, movie.Title)

	now := time.Now().UTC().Format(time.RFC3339)
	movie.Modified = now

	id := movie.Tmdb_Id

	err := _scrape(wid, s.tmdb, id, movie)
	if err != nil {
		mlog.Warning("RESCRAPE FAILED (%d) [%d] %s: %s", wid, movie.ID, movie.Title, err)
		s.bus.Pub(nil, "/event/workunit/done")

		return
	}

	s.dto.BaseURL = s.tmdb.BaseUrl
	s.dto.SecureBaseURL = s.tmdb.SecureBaseUrl

	mlog.Info("RESCRAPE COMPLETED (%d) [%d] %s", wid, movie.ID, movie.Title)

	msg := &pubsub.Message{Payload: s.dto}
	s.bus.Pub(msg, "/event/movie/rescraped")
}

func _scrape(wid int, client *tmdb.Tmdb, id uint64, movie *model.Movie) error {
	gmr, err := client.GetMovie(id)
	if err != nil {
		return fmt.Errorf("FAILED GETTING MOVIE (%d) [%s]", wid, movie.Title)
	}

	movie.Title = gmr.Title
	movie.Original_Title = gmr.Original_Title
	movie.Runtime = gmr.Runtime
	movie.Tmdb_Id = gmr.Id
	movie.Imdb_Id = gmr.Imdb_Id
	movie.Overview = gmr.Overview
	movie.Tagline = gmr.Tagline
	movie.Cover = gmr.Poster_Path
	movie.Backdrop = gmr.Backdrop_Path

	if len(gmr.Release_Date) > 0 {
		movie.Year = gmr.Release_Date[0:4]
	} else {
		mlog.Warning("release date is short: (%s)", gmr.Release_Date)
	}

	movie.Genres = ""
	for i := 0; i < len(gmr.Genres); i++ {
		attr := &gmr.Genres[i]
		if movie.Genres == "" {
			movie.Genres = attr.Name
		} else {
			if strings.Contains(movie.Genres, attr.Name) {
				continue
			}

			movie.Genres += "|" + attr.Name
		}
	}

	movie.Vote_Average = gmr.Vote_Average
	movie.Vote_Count = gmr.Vote_Count

	movie.Production_Countries = ""
	for i := 0; i < len(gmr.Production_Countries); i++ {
		attr := &gmr.Production_Countries[i]
		if movie.Production_Countries == "" {
			movie.Production_Countries = attr.Name
		} else {
			if strings.Contains(movie.Production_Countries, attr.Name) {
				continue
			}

			movie.Production_Countries += "|" + attr.Name
		}
	}

	if movie.Imdb_Id != "" {
		url := fmt.Sprintf("https://www.imdb.com/title/%s", movie.Imdb_Id)
		data, err := lib.RestGet(url)
		if err != nil {
			return fmt.Errorf("IMDB Error: %s", err)
		}

		imdb := getImdb(data)

		movie.Director = imdb.Director
		movie.Writer = imdb.Writers
		movie.Actors = imdb.Actors
		movie.Awards = imdb.Awards
		movie.Imdb_Rating = imdb.Rating
		movie.Imdb_Votes = imdb.Votes
	}

	return nil
}

func (s *Scraper) configChanged(msg *pubsub.Message) {
	s.settings = msg.Payload.(*lib.Settings)
}

var reJSONLD = regexp.MustCompile(`<script type="application/ld\+json">([^\<]*)`)

func getImdb(data string) *model.Imdb {
	imdb := &model.Imdb{}

	ld := reJSONLD.FindStringSubmatch(data)

	if err := json.Unmarshal([]byte(ld[1]), &imdb); err != nil {
		imdbAlt := &model.ImdbAlt{}
		if err := json.Unmarshal([]byte(ld[1]), &imdbAlt); err != nil {
			mlog.Warning("Unable to unmarshal imdb data: %s\n%s\n", err, ld[1])
			return imdb
		}

		imdb.Actors = imdbAlt.Actors
		imdb.Awards = imdbAlt.Awards
		imdb.Director = imdbAlt.Director
		imdb.Rating = imdbAlt.Rating
		imdb.Votes = imdbAlt.Votes
		imdb.Writers = imdbAlt.Writers

		return imdb
	}

	return imdb
}
