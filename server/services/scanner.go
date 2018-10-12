package services

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/jbrodriguez/actor"
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	"github.com/micro/go-micro/client"

	"mediagui/lib"
	"mediagui/model"
	"mediagui/proto"
)

// Scanner -
type Scanner struct {
	bus      *pubsub.PubSub
	settings *lib.Settings

	actor *actor.Actor

	re           []*lib.Rexp
	includedMask string
}

// NewScanner -
func NewScanner(bus *pubsub.PubSub, settings *lib.Settings) *Scanner {
	scanner := &Scanner{
		bus:      bus,
		settings: settings,
		actor:    actor.NewActor(bus),
	}
	return scanner
}

// Start -
func (s *Scanner) Start() {
	mlog.Info("Starting service Scanner ...")

	s.actor.Register("/command/movie/scan", s.scanMovies)
	s.actor.Register("/event/config/changed", s.configChanged)

	re := []string{
		`(?i)(.*)/(?P<Resolution>.*?)/(?P<Name>.*?)\s\((?P<Year>\d\d\d\d)\)/(?:.*/)*bdmv/index.(?P<FileType>bdmv)$`,
		`(?i)(.*)/(?P<Resolution>.*?)/(?P<Name>.*?)\s\((?P<Year>\d\d\d\d)\)/(?:.*/)*.*\.(?P<FileType>iso|img|nrg|mkv|avi|xvid|ts|mpg|dvr-ms|mdf|wmv|mp4)$`,
		`(?i)(.*)/(?P<Resolution>.*?)/(?P<Name>.*?)\s\((?P<Year>\d\d\d\d)\)/(?:.*/)*(?:video_ts|hv000i01)\.(?P<FileType>ifo)$`,
		// `(?i)(?P<Name>.*?)\s\((?P<Year>\d\d\d\d)\)/(?:.*/)*bdmv/index.(?P<FileType>bdmv)$`,
		// `(?i)(?P<Name>.*?)\s\((?P<Year>\d\d\d\d)\)/(?:.*/)*.*\.(?P<FileType>iso|img|nrg|mkv|avi|xvid|ts|mpg|dvr-ms|mdf|wmv)$`,
		// `(?i)(?P<Name>.*?)\s\((?P<Year>\d\d\d\d)\)/(?:.*/)*(?:video_ts|hv000i01)\.(?P<FileType>ifo)$`,
		// `(?i)(?P<Name>.*?)\s\((?P<Year>\d\d\d\d)\)\.(?P<FileType>iso|img|nrg|mkv|avi|xvid|ts|mpg|dvr-ms|mdf|wmv)$`,
	}

	s.re = make([]*lib.Rexp, 0)
	for _, regex := range re {
		s.re = append(s.re, &lib.Rexp{Exp: regexp.MustCompile(regex)})
	}

	s.includedMask = ".bdmv|.iso|.img|.nrg|.mkv|.avi|.xvid|.ts|.mpg|.dvr-ms|.mdf|.wmv|.ifo|.mp4"

	// cmd.Init()

	go s.actor.React()
}

// Stop -
func (s *Scanner) Stop() {
	mlog.Info("Stopped service Scanner ...")
}

func (s *Scanner) scanMovies(_ *pubsub.Message) {
	defer s.bus.Pub(nil, "/event/workunit/done")

	// folders := []string{
	// 	"/Volumes/hal-films",
	// 	"/Volumes/wopr-films",
	// }

	// ping := "ping -c1 %s > /dev/null && echo \"YES\" || echo \"NO\""
	// lib.Notify(s.bus, "import:begin", "Import process started")

	if s.settings.UnraidMode {
		// folders := []string{
		// 	`/mnt/user/films`,
		// }
		// filenames := []string{
		// 	"wopr:/mnt/user/films/bluray/ There Be Dragons (2011)/BDMV/BACKUP/MovieObject.bdmv",
		// 	"wopr:/mnt/user/films/bluray/ There Be Dragons (2011)/BDMV/BACKUP/index.bdmv",
		// 	"wopr:/mnt/user/films/bluray/ There Be Dragons (2011)/BDMV/MovieObject.bdmv",
		// 	"wopr:/mnt/user/films/bluray/ There Be Dragons (2011)/BDMV/index.bdmv",
		// 	"wopr:/mnt/user/films/bluray/ There Be Dragons (2011)/CERTIFICATE/BACKUP/id.bdmv",
		// 	"wopr:/mnt/user/films/bluray/ There Be Dragons (2011)/CERTIFICATE/id.bdmv",
		// 	"wopr:/mnt/user/films/bluray/'71 (2014)/BDMV/BACKUP/MovieObject.bdmv",
		// 	"wopr:/mnt/user/films/bluray/'71 (2014)/BDMV/BACKUP/index.bdmv",
		// 	"wopr:/mnt/user/films/bluray/'71 (2014)/BDMV/MovieObject.bdmv",
		// 	"wopr:/mnt/user/films/bluray/'71 (2014)/BDMV/index.bdmv",
		// 	"wopr:/mnt/user/films/bluray/10 Things I Hate About You (1999)/BDMV/BACKUP/MovieObject.bdmv",
		// 	"wopr:/mnt/user/films/bluray/10 Things I Hate About You (1999)/BDMV/BACKUP/index.bdmv",
		// 	"wopr:/mnt/user/films/bluray/10 Things I Hate About You (1999)/BDMV/MovieObject.bdmv",
		// 	"wopr:/mnt/user/films/bluray/10 Things I Hate About You (1999)/BDMV/index.bdmv",
		// 	"wopr:/mnt/user/films/bluray/12 Years A Slave (2013)/BDMV/BACKUP/MovieObject.bdmv",
		// 	"wopr:/mnt/user/films/bluray/12 Years A Slave (2013)/BDMV/BACKUP/index.bdmv",
		// 	"wopr:/mnt/user/films/bluray/12 Years A Slave (2013)/BDMV/MovieObject.bdmv",
		// 	"wopr:/mnt/user/films/bluray/12 Years A Slave (2013)/BDMV/index.bdmv",
		// 	"wopr:/mnt/user/films/bluray/12 Years A Slave (2013)/CERTIFICATE/BACKUP/id.bdmv",
		// 	"wopr:/mnt/user/films/bluray/12 Years A Slave (2013)/CERTIFICATE/id.bdmv",
		// 	"wopr:/mnt/user/films/bluray/13 (2010)/BDMV/BACKUP/MovieObject.bdmv",
		// 	"wopr:/mnt/user/films/bluray/13 (2010)/BDMV/BACKUP/index.bdmv",
		// 	"wopr:/mnt/user/films/bluray/13 (2010)/BDMV/MovieObject.bdmv",
		// 	"wopr:/mnt/user/films/bluray/13 (2010)/BDMV/index.bdmv",
		// 	"wopr:/mnt/user/films/bluray/13 (2010)/CERTIFICATE/BACKUP/id.bdmv",
		// 	"wopr:/mnt/user/films/bluray/13 (2010)/CERTIFICATE/id.bdmv",
		// 	"wopr:/mnt/user/films/blurip/10 Things I Hate About You (1999)/movie.mkv",
		// }

		// mlog.Info("started analysis")
		// s.analyze(filenames)
		// mlog.Info("finished analysis")

		for _, host := range s.settings.UnraidHosts {
			// Create new request to service go.micro.srv.example, method Example.Call
			req := client.NewRequest("io.jbrodriguez.mediagui.agent."+host, "Agent.Scan", &agent.ScanReq{
				// Folders: s.settings.MediaFolders,
				Folders: s.settings.MediaFolders,
				Mask:    s.includedMask,
			})

			rsp := &agent.ScanRsp{}

			// Call service
			if err := client.Call(context.Background(), req, rsp, client.WithRequestTimeout(time.Duration(5)*time.Minute)); err != nil {
				mlog.Warning("Unable to connect to service (%s): %s", "io.jbrodriguez.mediagui.agent."+host, err)
				lib.Notify(s.bus, "import:progress", "Unable to connect to host "+host)
				// lib.Notify(s.bus, "import:end", "Import process finished")
				// return
				continue
			}

			s.analyze(rsp.Filenames)
		}
	} else {
		for _, folder := range s.settings.MediaFolders {
			err := s.walk(folder)
			if err != nil {
				mlog.Warning("Unable to scan folder (%s): %s", folder, err)
			}
		}
	}

	// lib.Notify(s.bus, "import:end", "Import process finished")
}

func (s *Scanner) analyze(files []string) {
	mlog.Info("Found %d files", len(files))

	for _, file := range files {
		// comparePath := strings.TrimPrefix(path, file)
		// mlog.Info("folder: %s, comparePath: %s", folder, comparePath)
		// TODO: remove folder from path to match against rexp

		for i := 0; i < 3; i++ {
			rmap := s.re[i].Match(file)
			if rmap == nil {
				continue
			}

			resolution, ok := rmap["Resolution"]
			if !ok {
				resolution = ""
			}

			movie := &model.Movie{Title: rmap["Name"], File_Title: rmap["Name"], Year: rmap["Year"], Resolution: resolution, FileType: rmap["FileType"], Location: file, Stub: 0}
			// mlog.Info("FOUND [%s] (%s)", movie.Title, movie.Location)

			// mlog.Info("found %s", movie.Location)
			msg := &pubsub.Message{Payload: movie}
			s.bus.Pub(msg, "/event/movie/found")
			// self.Bus.MovieFound <- movie
		}
	}
}

func (s *Scanner) walk(folder string) error {
	if folder[len(folder)-1] != '/' {
		folder = folder + "/"
	}

	err := filepath.Walk(folder, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			mlog.Info("scanner.walk: %s (%s) - [%+v]", err, path, f)
		}

		if !strings.Contains(s.includedMask, strings.ToLower(filepath.Ext(path))) {
			// mlog.Info("[%s] excluding %s", filepath.Ext(path), path)
			return nil
		}

		// comparePath := strings.TrimPrefix(path, folder)
		// mlog.Info("folder: %s, comparePath: %s", folder, comparePath)
		// TODO: remove folder from path to match against rexp

		for i := 0; i < 3; i++ {
			var rmap = s.re[i].Match(path)
			if rmap == nil {
				continue
			}

			var resolution string
			var ok bool
			if resolution, ok = rmap["Resolution"]; !ok {
				resolution = ""
			}

			movie := &model.Movie{Title: rmap["Name"], File_Title: rmap["Name"], Year: rmap["Year"], Resolution: resolution, FileType: rmap["FileType"], Location: path}
			// mlog.Info("FOUND [%s] (%s)", movie.Title, movie.Location)

			msg := &pubsub.Message{Payload: movie}
			s.bus.Pub(msg, "/event/movie/found")
			// self.Bus.MovieFound <- movie

			return nil
		}

		return nil
	})

	return err
}

func (s *Scanner) configChanged(msg *pubsub.Message) {
	s.settings = msg.Payload.(*lib.Settings)
}
