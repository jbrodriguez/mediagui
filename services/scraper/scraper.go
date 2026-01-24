package scraper

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jbrodriguez/go-tmdb"

	"mediagui/domain"
	"mediagui/lib"
	"mediagui/logger"
)

// Scraper -
type Scraper struct {
	ctx *domain.Context

	tmdb            *tmdb.Tmdb
	tmdbRateLimiter *lib.RateLimiter
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

	// Initialize rate limiter
	s.tmdbRateLimiter = lib.NewRateLimiter(s.ctx.TmdbRateLimit, s.ctx.TmdbRateInterval)

	// s.actor.Register("/command/movie/scrape", s.scrapeMovie)
	// s.actor.Register("/command/movie/rescrape", s.reScrapeMovie)
	// s.actor.Register("/event/config/changed", s.configChanged)

	logger.Blue("started service scraper ...")

	return nil
}

// Stop -
func (s *Scraper) Stop() {
	// Clean up rate limiter
	if s.tmdbRateLimiter != nil {
		s.tmdbRateLimiter.Stop()
	}

	logger.Blue("stopped service scraper ...")
}

func (s *Scraper) ScrapeMovie(movie *domain.Movie) (code int, e error) {
	lib.Notify(s.ctx.Hub, "import:progress", fmt.Sprintf("SCRAPE REQUESTED [%s]", movie.Title))

	now := time.Now().UTC().Format(time.RFC3339)
	movie.Added = now
	movie.Modified = now

	movie.Score = 0
	movie.ShowIfDuplicate = 1

	// // Add delay with random jitter to avoid rate limiting
	// delay := 5*time.Second + time.Duration(rand.Intn(3000))*time.Millisecond
	// logger.Blue("TMDB rate limit delay: %v for [%s]", delay, movie.Title)
	// time.Sleep(delay)

	lib.Notify(s.ctx.Hub, "import:progress", fmt.Sprintf("STARTED TMDB [%s]", movie.Title))

	// Rate limit TMDB search API call
	logger.Blue("TMDB: Waiting for rate limit token for SearchMovie [%s]", movie.Title)
	s.tmdbRateLimiter.Wait()

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

	if err := scrape(s.tmdb, id, movie, s.tmdbRateLimiter, s.ctx.OmdbKey); err != nil {
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

	err := scrape(s.tmdb, id, movie, s.tmdbRateLimiter, s.ctx.OmdbKey)
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

func scrape(client *tmdb.Tmdb, id uint64, movie *domain.Movie, tmdbRateLimiter *lib.RateLimiter, omdbKey string) error {
	// Rate limit TMDB GetMovie API call
	logger.Blue("TMDB: Waiting for rate limit token for GetMovie [%d]", id)
	tmdbRateLimiter.Wait()

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

	// Get data from OMDb (director, actors, writer, rating)
	if movie.Imdb_Id != "" && omdbKey != "" {
		logger.Blue("OMDb: Fetching data for [%s]", movie.Imdb_Id)
		omdb, err := lib.FetchOmdb(movie.Imdb_Id, omdbKey)
		if err != nil {
			logger.Yellow("OMDb: Failed to fetch data for [%s]: %s", movie.Imdb_Id, err)
		} else {
			movie.Director = omdb.Director
			movie.Writer = omdb.Writer
			movie.Actors = omdb.Actors
			movie.Awards = omdb.Awards
			movie.Imdb_Rating = omdb.GetImdbRating()
			movie.Imdb_Votes = omdb.GetImdbVotes()

			// Use OMDb genre if TMDB didn't provide one
			if movie.Genres == "" && omdb.Genre != "" {
				movie.Genres = strings.ReplaceAll(omdb.Genre, ", ", "|")
			}
		}
	}

	return nil
}

