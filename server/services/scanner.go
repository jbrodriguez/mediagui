package services

import (
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	"jbrodriguez/mediagui/server/lib"
	"jbrodriguez/mediagui/server/model"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Scanner struct {
	Service

	bus      *pubsub.PubSub
	settings *lib.Settings
	// socket   *Socket

	mailbox chan *pubsub.Mailbox

	re           []*lib.Rexp
	includedMask string
}

func NewScanner(bus *pubsub.PubSub, settings *lib.Settings) *Scanner {
	scanner := &Scanner{bus: bus, settings: settings}
	scanner.init()
	return scanner
}

func (s *Scanner) Start() {
	mlog.Info("Starting service Scanner ...")

	s.mailbox = s.register(s.bus, "/command/movie/scan", s.scanMovies)

	re := []string{
		`(?i)(?P<Resolution>.*?)/(?P<Name>.*?)\s\((?P<Year>\d\d\d\d)\)/(?:.*/)*bdmv/index.(?P<FileType>bdmv)$`,
		`(?i)(?P<Resolution>.*?)/(?P<Name>.*?)\s\((?P<Year>\d\d\d\d)\)/(?:.*/)*.*\.(?P<FileType>iso|img|nrg|mkv|avi|xvid|ts|mpg|dvr-ms|mdf|wmv)$`,
		`(?i)(?P<Resolution>.*?)/(?P<Name>.*?)\s\((?P<Year>\d\d\d\d)\)/(?:.*/)*(?:video_ts|hv000i01)\.(?P<FileType>ifo)$`,
		`(?i)(?P<Name>.*?)\s\((?P<Year>\d\d\d\d)\)/(?:.*/)*bdmv/index.(?P<FileType>bdmv)$`,
		`(?i)(?P<Name>.*?)\s\((?P<Year>\d\d\d\d)\)/(?:.*/)*.*\.(?P<FileType>iso|img|nrg|mkv|avi|xvid|ts|mpg|dvr-ms|mdf|wmv)$`,
		`(?i)(?P<Name>.*?)\s\((?P<Year>\d\d\d\d)\)/(?:.*/)*(?:video_ts|hv000i01)\.(?P<FileType>ifo)$`,
		`(?i)(?P<Name>.*?)\s\((?P<Year>\d\d\d\d)\)\.(?P<FileType>iso|img|nrg|mkv|avi|xvid|ts|mpg|dvr-ms|mdf|wmv)$`,
	}

	s.re = make([]*lib.Rexp, 0)
	for _, regex := range re {
		s.re = append(s.re, &lib.Rexp{Exp: regexp.MustCompile(regex)})
	}

	s.includedMask = ".bdmv|.iso|.img|.nrg|.mkv|.avi|.xvid|.ts|.mpg|.dvr-ms|.mdf|.wmv|.ifo"

	go s.react()
}

func (s *Scanner) Stop() {
	mlog.Info("Stopped service Scanner ...")
}

func (s *Scanner) react() {
	for mbox := range s.mailbox {
		// mlog.Info("Scanner:Topic: %s", mbox.Topic)
		s.dispatch(mbox.Topic, mbox.Content)
	}
}

func (s *Scanner) scanMovies(msg *pubsub.Message) {
	folders := []string{
		"/Volumes/hal-films",
		"/Volumes/wopr-films",
	}

	// ping := "ping -c1 %s > /dev/null && echo \"YES\" || echo \"NO\""

	begin := &pubsub.Message{Payload: "Import process started"}
	s.bus.Pub(begin, "import:begin")

	for _, folder := range folders {
		err := s.walk(folder)
		if err != nil {
			mlog.Warning("Unable to scan folder (%s): %s", folder, err)
		}
	}

	end := &pubsub.Message{Payload: "Import process finished."}
	s.bus.Pub(end, "import:end")
}

func (s *Scanner) walk(folder string) error {
	if folder[len(folder)-1] != '/' {
		folder = folder + "/"
	}

	err := filepath.Walk(folder, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			mlog.Info("scanner.walk: %s", err)
		}

		if !strings.Contains(s.includedMask, strings.ToLower(filepath.Ext(path))) {
			// mlog.Info("[%s] excluding %s", filepath.Ext(path), path)
			return nil
		}

		comparePath := strings.TrimPrefix(path, folder)
		// mlog.Info("folder: %s, comparePath: %s", folder, comparePath)
		// TODO: remove folder from path to match against rexp

		for i := 0; i < 3; i++ {
			var rmap = s.re[i].Match(comparePath)
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
