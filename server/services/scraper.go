package services

import (
	// sc "bitbucket.org/jdpalmer/statecraft"
	"github.com/jbrodriguez/go-tmdb"
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	// "io/ioutil"
	"errors"
	"fmt"
	"jbrodriguez/mediagui/server/lib"
	"jbrodriguez/mediagui/server/model"
	"strconv"
	"strings"
	"time"
)

type Scraper struct {
	Service

	bus      *pubsub.PubSub
	settings *lib.Settings
	pool     *lib.Pool
	tmdb     *tmdb.Tmdb
	// socket   *Socket

	mailbox chan *pubsub.Mailbox

	// m *sc.Machine
}

func NewScraper(bus *pubsub.PubSub, settings *lib.Settings) *Scraper {
	scraper := &Scraper{bus: bus, settings: settings}
	scraper.init()
	return scraper
}

func (s *Scraper) Start() {
	mlog.Info("Starting service Scraper ...")

	var err error
	s.tmdb, err = tmdb.NewClient("e610ded10c3f47d05fe797961d90fea6", false)
	if err != nil {
		mlog.Fatalf("Unable to create tmdb client: %s", err)
	}

	s.mailbox = s.register(s.bus, "/command/movie/scrape", s.scrapeMovie)
	s.registerAdditional(s.bus, "/command/movie/rescrape", s.reScrapeMovie, s.mailbox)
	s.registerAdditional(s.bus, "/event/config/changed", s.configChanged, s.mailbox)

	s.pool = lib.NewPool(12, 4000)

	go s.react()
}

func (s *Scraper) Stop() {
	mlog.Info("Stopped service Scraper ...")
}

func (s *Scraper) react() {
	for mbox := range s.mailbox {
		// mlog.Info("Scraper:Topic: %s", mbox.Topic)
		s.dispatch(mbox.Topic, mbox.Content)
	}
}

func (s *Scraper) scrapeMovie(msg *pubsub.Message) {
	movie := msg.Payload.(*model.Movie)

	scrape := &Scrape{
		s.bus,
		s.tmdb,
		&model.ScrapeDTO{
			BasePath: s.settings.WebDir,
			Movie:    movie,
			Forced:   false,
		},
	}

	s.pool.Exec(scrape)
}

type Scrape struct {
	bus  *pubsub.PubSub
	tmdb *tmdb.Tmdb
	dto  *model.ScrapeDTO
}

func (s *Scrape) Execute() {
	lib.Notify(s.bus, "import:progress", fmt.Sprintf("SCRAPE REQUESTED [%s]", s.dto.Movie.Title))

	now := time.Now().UTC().Format(time.RFC3339)
	s.dto.Movie.Added = now
	s.dto.Movie.Modified = now

	s.dto.Movie.Score = 0

	lib.Notify(s.bus, "import:progress", fmt.Sprintf("STARTED TMDB [%s]", s.dto.Movie.Title))
	movies, err := s.tmdb.SearchMovie(s.dto.Movie.Title)
	if err != nil {
		mlog.Error(err)
		return
	}

	if movies.Total_Results == 0 {
		lib.Notify(s.bus, "import:progress", fmt.Sprintf("TMDB: NO MATCH FOUND [%s]", s.dto.Movie.Title))
		return
	} else if movies.Total_Results > 1 {
		lib.Notify(s.bus, "import:progress", fmt.Sprintf("TMDB: MORE THAN ONE [%s]", s.dto.Movie.Title))
	}

	id := movies.Results[0].Id

	_scrape(s.tmdb, id, s.dto.Movie)

	s.dto.BaseUrl = s.tmdb.BaseUrl
	s.dto.SecureBaseUrl = s.tmdb.SecureBaseUrl

	lib.Notify(s.bus, "import:progress", fmt.Sprintf("SCRAPE COMPLETED [%s]", s.dto.Movie.Title))

	msg := &pubsub.Message{Payload: s.dto}
	s.bus.Pub(msg, "/event/movie/scraped")
}

func (s *Scraper) reScrapeMovie(msg *pubsub.Message) {
	movie := msg.Payload.(*model.Movie)

	reScrape := &ReScrape{
		s.bus,
		s.tmdb,
		&model.ScrapeDTO{
			BasePath: s.settings.WebDir,
			Movie:    movie,
			Forced:   true,
		},
	}

	s.pool.Exec(reScrape)
}

type ReScrape struct {
	bus  *pubsub.PubSub
	tmdb *tmdb.Tmdb
	dto  *model.ScrapeDTO
}

func (s *ReScrape) Execute() {
	lib.Notify(s.bus, "import:progress", fmt.Sprintf("RESCRAPE REQUESTED [%s]", s.dto.Movie.Title))

	now := time.Now().UTC().Format(time.RFC3339)
	s.dto.Movie.Modified = now

	id := s.dto.Movie.Tmdb_Id

	err := _scrape(s.tmdb, id, s.dto.Movie)
	if err != nil {
		lib.Notify(s.bus, "import:progress", err.Error())
		return
	}

	s.dto.BaseUrl = s.tmdb.BaseUrl
	s.dto.SecureBaseUrl = s.tmdb.SecureBaseUrl

	lib.Notify(s.bus, "import:progress", fmt.Sprintf("RESCRAPE COMPLETED [%s]", s.dto.Movie.Title))

	msg := &pubsub.Message{Payload: s.dto}
	s.bus.Pub(msg, "/event/movie/rescraped")
}

func _scrape(tmdb *tmdb.Tmdb, id uint64, movie *model.Movie) error {
	// log.Printf("before getmovie [%d] %s", id, media.Movie.Title)
	// mlog.Info("[%s] before getmovie [%s]", s.dto.Movie.Title)
	gmr, err := tmdb.GetMovie(id)
	if err != nil {
		return errors.New(fmt.Sprintf("FAILED GETTING MOVIE [%s]", movie.Title))
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

	for i := 0; i < len(gmr.Genres); i++ {
		attr := &gmr.Genres[i]
		if movie.Genres == "" {
			movie.Genres = attr.Name
		} else {
			movie.Genres += " " + attr.Name
		}
	}

	movie.Vote_Average = gmr.Vote_Average
	movie.Vote_Count = gmr.Vote_Count

	for i := 0; i < len(gmr.Production_Countries); i++ {
		attr := &gmr.Production_Countries[i]
		if movie.Production_Countries == "" {
			movie.Production_Countries = attr.Name
		} else {
			movie.Production_Countries += "|" + attr.Name
		}
	}

	var omdb model.Omdb

	// lib.Notify(s.bus, "import:progress", fmt.Sprintf("STARTED OMDB [%s]", movie.Title))
	err = lib.RestGet(fmt.Sprintf("http://www.omdbapi.com/?i=%s", movie.Imdb_Id), &omdb)
	if err != nil {
		return errors.New(fmt.Sprintf("OMDB Error: %s", err))
	}

	// lib.Notify(s.bus, "import:progress", fmt.Sprintf("omdb: %+v", omdb))

	vote := strings.Replace(omdb.Imdb_Vote, ",", "", -1)
	imdb_vote, _ := strconv.ParseUint(vote, 0, 64)
	imdb_rating, _ := strconv.ParseFloat(omdb.Imdb_Rating, 64)

	movie.Director = omdb.Director
	movie.Writer = omdb.Writer
	movie.Actors = omdb.Actors
	movie.Awards = omdb.Awards
	movie.Imdb_Rating = imdb_rating
	movie.Imdb_Votes = imdb_vote

	return nil
}

func (s *Scraper) configChanged(msg *pubsub.Message) {
	s.settings = msg.Payload.(*lib.Settings)
}
