package main

import (
	// "fmt"
	"encoding/json"

	"github.com/jbrodriguez/mlog"

	// "github.com/stretchr/testify/assert"
	// "io"

	"jbrodriguez/mediagui/server/src/model"
	"log"
	"os"

	// "os/signal"
	"path/filepath"
	// "strings"
	"testing"
)

// func write(filename, text string) error {
// 	file, err := os.Create(filename)
// 	defer file.Close()

// 	if err != nil {
// 		return err
// 	}

// 	_, err = io.WriteString(file, fmt.Sprintf(text))
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func TestSettingsNotFound(t *testing.T) {
// 	home := os.Getenv("HOME")
// 	// cwd, err := os.Getwd()
// 	// if err != nil {
// 	// 	log.Printf("Unable to get current directory: %s", err.Error)
// 	// 	os.Exit(1)
// 	// }

// 	locations := []string{
// 		filepath.Join(home, "tmp/mgtest/.mediagui/mediagui.conf"),
// 		filepath.Join(home, "tmp/mgtest/usr/local/etc/mediagui.conf"),
// 		filepath.Join(home, "tmp/mgtest/mediagui.conf"),
// 	}

// 	settings, err := lib.NewSettings(Version, home, locations)

// 	if assert.Error(t, err) {
// 		mlog.Info("Ok: error was:\n%s", err.Error())
// 	} else {
// 		mlog.Fatalf("settings: %+v", settings)
// 	}
// }

// func TestSettingsFound(t *testing.T) {
// 	home := os.Getenv("HOME")

// 	// cwd, err := os.Getwd()
// 	// if err != nil {
// 	// 	mlog.Info("Unable to get current directory: %s", err.Error)
// 	// 	os.Exit(1)
// 	// }
// 	path := filepath.Join(home, "tmp/mgtest/.mediagui")
// 	os.MkdirAll(path, 0777)
// 	defer os.RemoveAll(path)

// 	text := "datadir=mg_datadir\nwebdir=mg_webdir\nmediafolders=movies/bluray|tv shows|movies/blurip"
// 	err := write(filepath.Join(path, "mediagui.conf"), text)

// 	assert.NoError(t, err)

// 	b, err := lib.Exists(filepath.Join(path, "mediagui.conf"))

// 	assert.Equal(t, true, b)

// 	locations := []string{
// 		filepath.Join(home, "tmp/mgtest/.mediagui/mediagui.conf"),
// 		filepath.Join(home, "tmp/mgtest/usr/local/etc/mediagui.conf"),
// 		filepath.Join(home, "tmp/mgtest/mediagui.conf"),
// 	}

// 	settings, err := lib.NewSettings(Version, home, locations)

// 	if assert.NoError(t, err) {
// 		assert.Equal(t, "mg_datadir", settings.DataDir)
// 		assert.Equal(t, "mg_webdir", settings.WebDir)
// 		assert.Equal(t, "", settings.LogDir)
// 		assert.Equal(t, strings.Split("movies/bluray|tv shows|movies/blurip", "|"), settings.MediaFolders)
// 	}

// }

// func TestServerService(t *testing.T) {
// 	home := os.Getenv("HOME")

// 	// cwd, err := os.Getwd()
// 	// if err != nil {
// 	// 	mlog.Info("Unable to get current directory: %s", err.Error)
// 	// 	os.Exit(1)
// 	// }
// 	path := filepath.Join(home, "tmp/mgtest/.mediagui")
// 	os.MkdirAll(path, 0777)
// 	defer os.RemoveAll(path)

// 	html := filepath.Join(home, "tmp/mgtest/web/index.html")
// 	text := fmt.Sprintf("datadir=mg_datadir\nwebdir=%s\nmediafolders=movies/bluray|tv shows|movies/blurip", html)
// 	err := write(filepath.Join(path, "mediagui.conf"), text)

// 	assert.NoError(t, err)

// 	b, err := lib.Exists(filepath.Join(path, "mediagui.conf"))

// 	assert.Equal(t, true, b)

// 	locations := []string{
// 		filepath.Join(home, "tmp/mgtest/.mediagui/mediagui.conf"),
// 		filepath.Join(home, "tmp/mgtest/usr/local/etc/mediagui.conf"),
// 		filepath.Join(home, "tmp/mgtest/mediagui.conf"),
// 	}

// 	settings, err := lib.NewSettings(Version, home, locations)

// 	if assert.NoError(t, err) {
// 		assert.Equal(t, "mg_datadir", settings.DataDir)
// 		assert.Equal(t, "mg_webdir", settings.WebDir)
// 		assert.Equal(t, "", settings.LogDir)
// 		assert.Equal(t, strings.Split("movies/bluray|tv shows|movies/blurip", "|"), settings.MediaFolders)
// 	}

// 	text = "<body><h1>we shall overcome</h1></body>"
// 	err = write(html, text)
// 	assert.NoError(t, err)

// 	bus := pubsub.New(623)

// 	server := services.NewServer(bus, settings)
// 	server.Start()
// }

// func TestScraper(t *testing.T) {
// 	// look for mediagui.conf at the following places
// 	// $HOME/.mediagui/mediagui.conf
// 	// /usr/local/etc/mediagui.conf
// 	// <current dir>/mediagui.conf
// 	home := os.Getenv("HOME")

// 	cwd, err := os.Getwd()
// 	if err != nil {
// 		log.Printf("Unable to get current directory: %s", err.Error())
// 		os.Exit(1)
// 	}

// 	locations := []string{
// 		filepath.Join(home, ".mediagui/mediagui.conf"),
// 		"/usr/local/etc/mediagui.conf",
// 		filepath.Join(cwd, "mediagui.conf"),
// 	}

// 	settings, err := lib.NewSettings(Version, home, locations)
// 	if err != nil {
// 		log.Printf("Unable to start the app: %s", err.Error())
// 		os.Exit(2)
// 	}

// 	mlog.Start(mlog.LevelInfo, "")

// 	mlog.Info("mediagui v%s starting ...", Version)

// 	bus := pubsub.New(623)

// 	core := services.NewCore(bus, settings)
// 	scraper := services.NewScraper(bus, settings)

// 	// socket.Start()
// 	scraper.Start()
// 	core.Start()

// 	movie := &model.Movie{
// 		Title: "Pulp Fiction",
// 	}

// 	msg := &pubsub.Message{Payload: movie}
// 	bus.Pub(msg, "/command/movie/scrape")

// 	mlog.Info("Press Ctrl+C to stop ...")

// 	c := make(chan os.Signal, 1)
// 	signal.Notify(c, os.Interrupt)
// 	for _ = range c {
// 		mlog.Info("Received an interrupt, shutting the app down ...")

// 		core.Stop()
// 		scraper.Stop()
// 		// socket.Stop()

// 		break
// 	}
// }

// func TestScraper(t *testing.T) {
// 	// look for mediagui.conf at the following places
// 	// $HOME/.mediagui/mediagui.conf
// 	// /usr/local/etc/mediagui.conf
// 	// <current dir>/mediagui.conf
// 	home := os.Getenv("HOME")

// 	cwd, err := os.Getwd()
// 	if err != nil {
// 		log.Printf("Unable to get current directory: %s", err.Error())
// 		os.Exit(1)
// 	}

// 	locations := []string{
// 		filepath.Join(home, ".mediagui/mediagui.conf"),
// 		"/usr/local/etc/mediagui.conf",
// 		filepath.Join(cwd, "mediagui.conf"),
// 	}

// 	settings, err := lib.NewSettings(Version, home, locations)
// 	if err != nil {
// 		log.Printf("Unable to start the app: %s", err.Error())
// 		os.Exit(2)
// 	}

// 	mlog.Start(mlog.LevelInfo, "")

// 	mlog.Info("mediagui v%s starting ...", Version)

// 	bus := pubsub.New(623)

// 	dal := services.NewDal(bus, settings)

// 	// socket.Start()
// 	dal.Start()

// 	movie := &model.Movie{
// 		Location: "wopr:/mnt/user/films/xvid/Visitor Q (2001)/ils-visitorq.avi",
// 	}

// 	check := &pubsub.Message{Payload: movie, Reply: make(chan interface{}, 3)}
// 	bus.Pub(check, "/command/movie/exists")

// 	reply := <-check.Reply
// 	exists := reply.(bool)

// 	if exists {
// 		mlog.Info("SKIPPED: exists [%s] (%s)", movie.Title, movie.Location)
// 	} else {
// 		mlog.Info("NEW: [%s] (%s)", movie.Title, movie.Location)
// 		// c.bus.Pub(msg, "/command/movie/scrape")
// 	}

// 	// mlog.Info("Press Ctrl+C to stop ...")

// 	// c := make(chan os.Signal, 1)
// 	// signal.Notify(c, os.Interrupt)
// 	// for _ = range c {
// 	// 	mlog.Info("Received an interrupt, shutting the app down ...")

// 	dal.Stop()

// 	// 	break
// 	// }
// }

func TestJsonLd(t *testing.T) {
	jsonString := `
	{
		"@context": "http://schema.org",
		"@type": "Movie",
		"url": "/title/tt5052474/",
		"name": "Sicario: Day of the Soldado",
		"image": "https://m.media-amazon.com/images/M/MV5BMjgyOWRhMDctZTZlNC00M2I1LWI0NDQtYzlmODdmYjY2MThiXkEyXkFqcGdeQXVyMzY0MTE3NzU@._V1_.jpg",
		"genre": [
		"Action",
		"Crime",
		"Drama",
		"Thriller"
		],
		"contentRating": "R",
		"actor": [
		{
			"@type": "Person",
			"url": "/name/nm0001125/",
			"name": "Benicio Del Toro"
		},
		{
			"@type": "Person",
			"url": "/name/nm0000982/",
			"name": "Josh Brolin"
		},
		{
			"@type": "Person",
			"url": "/name/nm5097044/",
			"name": "Isabela Moner"
		},
		{
			"@type": "Person",
			"url": "/name/nm0232998/",
			"name": "Jeffrey Donovan"
		}
		],
		"director": {
		"@type": "Person",
		"url": "/name/nm1356588/",
		"name": "Stefano Sollima"
		},
		"creator": [
		{
			"@type": "Person",
			"url": "/name/nm0792263/",
			"name": "Taylor Sheridan"
		},
		{
			"@type": "Organization",
			"url": "/company/co0410488/"
		},
		{
			"@type": "Organization",
			"url": "/company/co0172670/"
		},
		{
			"@type": "Organization",
			"url": "/company/co0310184/"
		}
		],
		"description": "Sicario: Day of the Soldado is a movie starring Benicio Del Toro, Josh Brolin, and Isabela Moner. The drug war on the U.S.-Mexico border has escalated as the cartels have begun trafficking terrorists across the US border. To fight...",
		"datePublished": "2018-06-27",
		"keywords": "somalia,desert,second part,cartel,gunfight",
		"aggregateRating": {
		"@type": "AggregateRating",
		"ratingCount": 48506,
		"bestRating": "10.0",
		"worstRating": "1.0",
		"ratingValue": "7.2"
		},
		"review": {
		"@type": "Review",
		"itemReviewed": {
			"@type": "CreativeWork",
			"url": "/title/tt5052474/"
		},
		"author": {
			"@type": "Person",
			"name": "rene-806-6219"
		},
		"dateCreated": "2018-07-18",
		"inLanguage": "English",
		"name": "The first 21 reviews: 10/10. You know something is fup.",
		"reviewBody": "My first review: Loved the first Sicario. Raw and different. This time I almost walked out of the cinema. Really bad. Could not believe its the same screenwriter. The story is really thin (and short). The overacting by Isabela Moner was too much on many moments. May not be her fault, but the director. Benico and Josh did an excellent job, but the movie is just bad. Most annoying is the rating her. The first 21 (!!) gave it a 10/10. The Godfather Part II came to seven (7). So something is really wrong. Really annoying because I really trust IMDB.",
		"reviewRating": {
			"@type": "Rating",
			"worstRating": "1",
			"bestRating": "10",
			"ratingValue": "3"
		}
		},
		"duration": "PT2H2M",
		"trailer": {
		"@type": "VideoObject",
		"name": "Trailer [OV]",
		"embedUrl": "/video/imdb/vi42580761",
		"thumbnail": {
			"@type": "ImageObject",
			"contentUrl": "https://m.media-amazon.com/images/M/MV5BMDE5ZjQ1OTgtN2JmYi00MTAwLWE1OTctYTMwNTUzZTBmMzNiXkEyXkFqcGdeQXRodW1ibmFpbC1pbml0aWFsaXplcg@@._V1_.jpg"
		},
		"thumbnailUrl": "https://m.media-amazon.com/images/M/MV5BMDE5ZjQ1OTgtN2JmYi00MTAwLWE1OTctYTMwNTUzZTBmMzNiXkEyXkFqcGdeQXRodW1ibmFpbC1pbml0aWFsaXplcg@@._V1_.jpg",
		"uploadDate": "2018-06-12T16:00:20Z"
		}
	}`

	var imdb model.Imdb

	if err := json.Unmarshal([]byte(jsonString), &imdb); err != nil {
		log.Fatal(err)
	}

	t.Logf("imdb: %+v", &imdb)
}

func TestMain(m *testing.M) {
	mlog.Start(mlog.LevelInfo, "")

	home := os.Getenv("HOME")
	path := filepath.Join(home, "tmp/mgtest")
	os.RemoveAll(path)

	ret := m.Run()

	os.RemoveAll(path)

	// mlog.Stop()

	os.Exit(ret)
}
