package scraper

import (
	"errors"

	"github.com/jbrodriguez/go-tmdb"

	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"mediagui/domain"
	"mediagui/lib"
	"mediagui/logger"
)

// Scraper -
type Scraper struct {
	ctx *domain.Context

	tmdb *tmdb.Tmdb
}

func Create(ctx *domain.Context) *Scraper {
	scraper := &Scraper{
		ctx: ctx,
	}
	return scraper
}

// Start -
func (s *Scraper) Start() error {
	var err error
	s.tmdb, err = tmdb.NewClient(s.ctx.TmdbKey, false)
	if err != nil {
		return fmt.Errorf("unable to create tmdb client: %s", err)
	}

	// s.actor.Register("/command/movie/scrape", s.scrapeMovie)
	// s.actor.Register("/command/movie/rescrape", s.reScrapeMovie)
	// s.actor.Register("/event/config/changed", s.configChanged)

	logger.Blue("started service scraper ...")

	return nil
}

// Stop -
func (s *Scraper) Stop() {
	logger.Blue("stopped service scraper ...")
}

func (s *Scraper) ScrapeMovie(movie *domain.Movie) (code int, e error) {
	lib.Notify(s.ctx.Hub, "import:progress", fmt.Sprintf("SCRAPE REQUESTED [%s]", movie.Title))

	now := time.Now().UTC().Format(time.RFC3339)
	movie.Added = now
	movie.Modified = now

	movie.Score = 0
	movie.ShowIfDuplicate = 1

	lib.Notify(s.ctx.Hub, "import:progress", fmt.Sprintf("STARTED TMDB [%s]", movie.Title))
	movies, err := s.tmdb.SearchMovie(movie.Title)
	if err != nil {
		return 1, err
	}

	if movies.Total_Results == 0 {
		lib.Notify(s.ctx.Hub, "import:progress", fmt.Sprintf("TMDB: NO MATCH FOUND [%s]", movie.Title))
		return 2, errors.New("no match found")
	} else if movies.Total_Results > 1 {
		lib.Notify(s.ctx.Hub, "import:progress", fmt.Sprintf("TMDB: MORE THAN ONE [%s]", movie.Title))
	}

	id := movies.Results[0].Id

	if err := scrape(s.tmdb, id, movie, s.ctx.UserAgent); err != nil {
		logger.Yellow("SCRAPE FAILED [%d] %s: %s", movie.ID, movie.Title, err)
	}

	movie.BaseURL = s.tmdb.BaseUrl
	movie.SecureBaseURL = s.tmdb.SecureBaseUrl

	lib.Notify(s.ctx.Hub, "import:progress", fmt.Sprintf("SCRAPE COMPLETED [%s]", movie.Title))

	return 0, nil
}

func (s *Scraper) ReScrape(movie *domain.Movie) (*domain.Movie, error) {
	logger.Blue("RESCRAPE REQUESTED TMDB [%d]", movie.Tmdb_Id)

	now := time.Now().UTC().Format(time.RFC3339)
	movie.Modified = now

	id := movie.Tmdb_Id

	err := scrape(s.tmdb, id, movie, s.ctx.UserAgent)
	if err != nil {
		// logger.Yellow("RESCRAPE FAILED [%d] %s: %s", movie.ID, movie.Title, err)
		return nil, err
	}

	movie.BaseURL = s.tmdb.BaseUrl
	movie.SecureBaseURL = s.tmdb.SecureBaseUrl

	logger.Blue("RESCRAPE COMPLETED TMDB [%d]", movie.Tmdb_Id)

	// msg := &pubsub.Message{Payload: s.dto}
	// s.bus.Pub(msg, "/event/movie/rescraped")

	return movie, nil
}

func scrape(client *tmdb.Tmdb, id uint64, movie *domain.Movie, agent string) error {
	gmr, err := client.GetMovie(id)
	if err != nil {
		return fmt.Errorf("FAILED GETTING MOVIE [%d]", movie.Tmdb_Id)
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
		logger.Yellow("release date is short: (%s)", gmr.Release_Date)
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
		data, err := lib.RestGet(url, agent)
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

// func (s *Scraper) configChanged(msg *pubsub.Message) {
// 	s.settings = msg.Payload.(*lib.Settings)
// }

var reJSONLD = regexp.MustCompile(`<script type="application/ld\+json">([^\<]*)`)

func getImdb(data string) *domain.Imdb {
	imdb := &domain.Imdb{}

	ld := reJSONLD.FindStringSubmatch(data)

	if err := json.Unmarshal([]byte(ld[1]), &imdb); err != nil {
		imdbAlt := &domain.ImdbAlt{}
		if err := json.Unmarshal([]byte(ld[1]), &imdbAlt); err != nil {
			logger.Yellow("Unable to unmarshal imdb data: %s\n%s\n", err, ld[1])
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
