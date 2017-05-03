package services

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"jbrodriguez/mediagui/server/src/lib"
	"jbrodriguez/mediagui/server/src/model"

	"github.com/jbrodriguez/actor"
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
)

// Importer -
type Importer struct {
	bus      *pubsub.PubSub
	settings *lib.Settings

	actor *actor.Actor
}

// NewImporter -
func NewImporter(bus *pubsub.PubSub, settings *lib.Settings) *Importer {
	importer := &Importer{
		bus:      bus,
		settings: settings,
		actor:    actor.NewActor(bus),
	}
	return importer
}

// Start -
func (i *Importer) Start() {
	mlog.Info("Starting service Importer ...")

	i.actor.Register("/cli/import", i.runImport)

	go i.actor.React()
}

// Stop -
func (i *Importer) Stop() {
	mlog.Info("Stopped service Importer")
}

func (i *Importer) runImport(msg *pubsub.Message) {
	file := filepath.Join(i.settings.WorkDir, "moviesi.txt")
	f, err := os.Open(file)
	if err != nil {
		mlog.Fatal(err)
	}
	defer f.Close() // this needs to be after the err check

	psv := csv.NewReader(f)
	psv.Comma = '|'

	lines, err := psv.ReadAll()
	if err != nil {
		mlog.Fatal(err)
	}

	mlog.Info("Is check only (%s)", i.settings.CheckOnly)

	if i.settings.CheckOnly {
		ofile := filepath.Join(i.settings.WorkDir, "movieo.txt")

		for _, line := range lines {
			name := strings.TrimSpace(line[0])

			options := lib.Options{
				Query:     name,
				Limit:     60,
				SortBy:    "added",
				SortOrder: "desc",
				FilterBy:  "title",
			}

			search := &pubsub.Message{Payload: &options, Reply: make(chan interface{}, capacity)}
			i.bus.Pub(search, "/get/movies")

			reply := <-search.Reply
			dto := reply.(*model.MoviesDTO)

			// sign := "-"
			if dto.Total == 0 {
				// sign = "+"

				res := fmt.Sprintf("%s|%s|%s|%s", line[0], line[1], line[2], line[3])
				WriteLine(ofile, res)
			}

		}
	} else {
		for index, line := range lines {
			movie := model.Movie{
				Title: line[0],
			}

			msg := &pubsub.Message{Payload: &movie, Reply: make(chan interface{}, capacity)}
			i.bus.Pub(msg, "/post/add")

			reply := <-msg.Reply
			stub := reply.(*model.Movie)

			stub.Score, _ = strconv.ParseUint(strings.TrimSpace(line[1]), 0, 64)

			watched := strings.TrimSpace(line[2])
			stub.Last_Watched = fmt.Sprintf("%s-%s-%sT00:00:00-05:00", watched[:4], watched[4:6], watched[6:])

			setMsg := &pubsub.Message{Payload: stub, Reply: make(chan interface{}, capacity)}
			i.bus.Pub(setMsg, "/put/movies/score")
			i.bus.Pub(setMsg, "/put/movies/watched")

			if index%30 == 0 {
				time.Sleep(10 * time.Second)
			}
		}
	}

	// movie := model.Movie{}
	// movie.Title = ""
	// add := &pubsub.Message{Payload: &movie, Reply: make(chan interface{}, capacity)}
	// i.bus.Pub(add, "/post/add")

}

// WriteLine -
func WriteLine(fullpath, line string) error {
	f, err := os.OpenFile(fullpath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(line + "\n")
	if err != nil {
		return err
	}

	return nil
}
