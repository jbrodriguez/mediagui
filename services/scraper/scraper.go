package scraper

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
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
	imdbRateLimiter *lib.RateLimiter
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

	// Initialize rate limiters
	s.tmdbRateLimiter = lib.NewRateLimiter(s.ctx.TmdbRateLimit, s.ctx.TmdbRateInterval)
	s.imdbRateLimiter = lib.NewRateLimiter(s.ctx.ImdbRateLimit, s.ctx.ImdbRateInterval)

	// s.actor.Register("/command/movie/scrape", s.scrapeMovie)
	// s.actor.Register("/command/movie/rescrape", s.reScrapeMovie)
	// s.actor.Register("/event/config/changed", s.configChanged)

	logger.Blue("started service scraper ...")

	return nil
}

// Stop -
func (s *Scraper) Stop() {
	// Clean up rate limiters
	if s.tmdbRateLimiter != nil {
		s.tmdbRateLimiter.Stop()
	}
	if s.imdbRateLimiter != nil {
		s.imdbRateLimiter.Stop()
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

	if err := scrape(s.tmdb, id, movie, s.ctx.UserAgent, s.tmdbRateLimiter, s.imdbRateLimiter); err != nil {
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

	err := scrape(s.tmdb, id, movie, s.ctx.UserAgent, s.tmdbRateLimiter, s.imdbRateLimiter)
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

func scrape(client *tmdb.Tmdb, id uint64, movie *domain.Movie, agent string, tmdbRateLimiter *lib.RateLimiter, imdbRateLimiter *lib.RateLimiter) error {
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

	if movie.Imdb_Id != "" {
		url := fmt.Sprintf("https://www.imdb.com/title/%s", movie.Imdb_Id)

		// Rate limit IMDB scraping call
		logger.Blue("IMDB: Waiting for rate limit token for scraping [%s]", movie.Imdb_Id)
		imdbRateLimiter.Wait()

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
var reNextData = regexp.MustCompile(`<script id="__NEXT_DATA__" type="application/json">([^<]+)</script>`)

func getImdb(data string) *domain.Imdb {
	imdb := &domain.Imdb{}

	// Try JSON-LD first (legacy)
	ld := reJSONLD.FindStringSubmatch(data)
	if len(ld) > 1 {
		if err := json.Unmarshal([]byte(ld[1]), &imdb); err != nil {
			imdbAlt := &domain.ImdbAlt{}
			if err := json.Unmarshal([]byte(ld[1]), &imdbAlt); err != nil {
				logger.Yellow("Unable to unmarshal imdb JSON-LD data: %s", err)
			} else {
				imdb.Actors = imdbAlt.Actors
				imdb.Awards = imdbAlt.Awards
				imdb.Director = imdbAlt.Director
				imdb.Rating = imdbAlt.Rating
				imdb.Votes = imdbAlt.Votes
				imdb.Writers = imdbAlt.Writers
			}
		}
	}

	// If director/actors are missing, try __NEXT_DATA__ fallback
	if imdb.Director == "" || imdb.Actors == "" {
		logger.Blue("Director or actors missing from JSON-LD, trying __NEXT_DATA__ fallback")
		getImdbFromNextData(data, imdb)
	}

	return imdb
}

func getImdbFromNextData(data string, imdb *domain.Imdb) {
	nextData := reNextData.FindStringSubmatch(data)
	if len(nextData) < 2 {
		logger.Yellow("No __NEXT_DATA__ found in IMDB page")
		return
	}

	var nextJSON map[string]interface{}
	if err := json.Unmarshal([]byte(nextData[1]), &nextJSON); err != nil {
		logger.Yellow("Unable to unmarshal __NEXT_DATA__: %s", err)
		return
	}

	// Navigate: props -> pageProps -> aboveTheFoldData -> principalCreditsV2
	props, ok := nextJSON["props"].(map[string]interface{})
	if !ok {
		return
	}

	pageProps, ok := props["pageProps"].(map[string]interface{})
	if !ok {
		return
	}

	atfData, ok := pageProps["aboveTheFoldData"].(map[string]interface{})
	if !ok {
		return
	}

	// Extract rating if missing
	if imdb.Rating == 0 {
		if ratingsSummary, ok := atfData["ratingsSummary"].(map[string]interface{}); ok {
			if aggRating, ok := ratingsSummary["aggregateRating"].(float64); ok {
				imdb.Rating = aggRating
			}
			if voteCount, ok := ratingsSummary["voteCount"].(float64); ok {
				imdb.Votes = uint64(voteCount)
			}
		}
	}

	// Extract credits
	principalCredits, ok := atfData["principalCreditsV2"].([]interface{})
	if !ok {
		logger.Yellow("No principalCreditsV2 found in __NEXT_DATA__")
		return
	}

	for _, pc := range principalCredits {
		pcMap, ok := pc.(map[string]interface{})
		if !ok {
			continue
		}

		// Get grouping (category)
		grouping := ""
		if group, ok := pcMap["grouping"].(map[string]interface{}); ok {
			if text, ok := group["text"].(string); ok {
				grouping = text
			}
		}

		// Get credits (people)
		credits, ok := pcMap["credits"].([]interface{})
		if !ok {
			continue
		}

		var names []string
		for _, credit := range credits {
			creditMap, ok := credit.(map[string]interface{})
			if !ok {
				continue
			}

			if name, ok := creditMap["name"].(map[string]interface{}); ok {
				if nameText, ok := name["nameText"].(map[string]interface{}); ok {
					if text, ok := nameText["text"].(string); ok {
						names = append(names, text)
					}
				}
			}
		}

		// Assign to appropriate field based on grouping
		nameStr := ""
		for i, name := range names {
			if i > 0 {
				nameStr += ", "
			}
			nameStr += name
		}

		switch grouping {
		case "Director", "Directors":
			if imdb.Director == "" {
				imdb.Director = nameStr
			}
		case "Writer", "Writers":
			if imdb.Writers == "" {
				imdb.Writers = nameStr
			}
		case "Star", "Stars":
			if imdb.Actors == "" {
				imdb.Actors = nameStr
			}
		}
	}

	logger.Blue("Extracted from __NEXT_DATA__: Director=%s, Writers=%s, Actors=%s", imdb.Director, imdb.Writers, imdb.Actors)
}
