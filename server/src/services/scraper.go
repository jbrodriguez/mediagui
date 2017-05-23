package services

import (
	"github.com/jbrodriguez/actor"
	"github.com/jbrodriguez/go-tmdb"
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	// "io/ioutil"
	"fmt"
	"jbrodriguez/mediagui/server/src/dto"
	"jbrodriguez/mediagui/server/src/lib"
	"jbrodriguez/mediagui/server/src/model"

	"regexp"
	"strconv"
	"strings"
	"time"
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

	// reRating = `<span[^>]*itemprop="ratingValue">([^<]*)</span>`
	// reVotes = `<span[^>]*itemprop="ratingCount">([^<]*)</span>`
	// reDirector = `<a href="/name[^"]*dr"[^>]*><span[^>]*>([^<]*)?</span>`
	// reWriter = `<a href="/name[^"]*wr"[^>]*><span[^>]*>([^<]*)?</span>`
	// reActor = `<a href="/name[^"]*st_sm"[^>]*><span[^>]*>([^<]*)?</span>`

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

	_scrape(wid, s.tmdb, id, movie)

	s.dto.BaseUrl = s.tmdb.BaseUrl
	s.dto.SecureBaseUrl = s.tmdb.SecureBaseUrl

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

	// lib.Notify(s.bus, "import:progress", fmt.Sprintf("RESCRAPE REQUESTED (%d) [%d] %s", wid, movie.Id, movie.Title))
	mlog.Info("RESCRAPE REQUESTED (%d) [%d] %s", wid, movie.Id, movie.Title)

	now := time.Now().UTC().Format(time.RFC3339)
	movie.Modified = now

	id := movie.Tmdb_Id

	err := _scrape(wid, s.tmdb, id, movie)
	if err != nil {
		// lib.Notify(s.bus, "import:progress", fmt.Sprintf("RESCRAPE FAILED (%d) [%d] %s: %s", wid, movie.Id, movie.Title, err))
		mlog.Warning("RESCRAPE FAILED (%d) [%d] %s: %s", wid, movie.Id, movie.Title, err)
		s.bus.Pub(nil, "/event/workunit/done")

		return
	}

	s.dto.BaseUrl = s.tmdb.BaseUrl
	s.dto.SecureBaseUrl = s.tmdb.SecureBaseUrl

	// lib.Notify(s.bus, "import:progress", fmt.Sprintf("RESCRAPE COMPLETED (%d) [%d] %s", wid, movie.Id, movie.Title))
	mlog.Info("RESCRAPE COMPLETED (%d) [%d] %s", wid, movie.Id, movie.Title)

	msg := &pubsub.Message{Payload: s.dto}
	s.bus.Pub(msg, "/event/movie/rescraped")
}

func _scrape(wid int, tmdb *tmdb.Tmdb, id uint64, movie *model.Movie) error {
	// log.Printf("before getmovie [%d] %s", id, media.Movie.Title)
	// mlog.Info("[%s] before getmovie [%s]", movie.Title)
	gmr, err := tmdb.GetMovie(id)
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

	movie.Year = gmr.Release_Date[0:4]

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

	// var omdb model.Omdb

	// lib.Notify(s.bus, "import:progress", fmt.Sprintf("STARTED OMDB [%s]", movie.Title))
	data, err := lib.RestGet(fmt.Sprintf("http://www.im%s.com/title/%s", "db", movie.Imdb_Id))
	if err != nil {
		return fmt.Errorf("OMDB Error: %s", err)
	}

	omdb := getOmdb(data)

	// lib.Notify(s.bus, "import:progress", fmt.Sprintf("omdb: %+v", omdb))

	vote := strings.Replace(omdb.Imdb_Vote, ",", "", -1)
	imdbVote, _ := strconv.ParseUint(vote, 0, 64)
	imdbRating, _ := strconv.ParseFloat(omdb.Imdb_Rating, 64)

	movie.Director = omdb.Director
	movie.Writer = omdb.Writer
	movie.Actors = omdb.Actors
	movie.Awards = omdb.Awards
	movie.Imdb_Rating = imdbRating
	movie.Imdb_Votes = imdbVote

	return nil
}

func (s *Scraper) configChanged(msg *pubsub.Message) {
	s.settings = msg.Payload.(*lib.Settings)
}

var reRating *regexp.Regexp = regexp.MustCompile(`<span[^>]*itemprop="ratingValue">([^<]*)</span>`)
var reVotes *regexp.Regexp = regexp.MustCompile(`<span[^>]*itemprop="ratingCount">([^<]*)</span>`)
var reDirector *regexp.Regexp = regexp.MustCompile(`<a href="/name[^"]*dr"[^>]*><span[^>]*>([^<]*)?</span>`)
var reWriter *regexp.Regexp = regexp.MustCompile(`<a href="/name[^"]*wr"[^>]*><span[^>]*>([^<]*)?</span>`)
var reActor *regexp.Regexp = regexp.MustCompile(`<a href="/name[^"]*st_sm"[^>]*><span[^>]*>([^<]*)?</span>`)

func getOmdb(data string) *model.Omdb {
	omdb := &model.Omdb{}

	rating := reRating.FindStringSubmatch(data)
	if len(rating) > 0 {
		omdb.Imdb_Rating = rating[1]
	}

	votes := reVotes.FindStringSubmatch(data)
	if len(votes) > 0 {
		omdb.Imdb_Vote = votes[1]
	}

	directors := reDirector.FindAllStringSubmatch(data, -1)
	if len(directors) > 0 {
		for _, director := range directors {
			if omdb.Director == "" {
				omdb.Director = director[1]
			} else {
				omdb.Director += ", " + director[1]
			}
		}
	}

	writers := reWriter.FindAllStringSubmatch(data, -1)
	if len(writers) > 0 {
		for _, writer := range writers {
			if omdb.Writer == "" {
				omdb.Writer = writer[1]
			} else {
				omdb.Writer += ", " + writer[1]
			}
		}
	}

	actors := reActor.FindAllStringSubmatch(data, -1)
	if len(actors) > 0 {
		for _, actor := range actors {
			if omdb.Actors == "" {
				omdb.Actors = actor[1]
			} else {
				omdb.Actors += ", " + actor[1]
			}
		}
	}

	return omdb
}
