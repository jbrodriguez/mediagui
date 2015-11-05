package services

import (
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	"github.com/myodc/go-micro/client"
	"github.com/myodc/go-micro/cmd"
	"jbrodriguez/mediagui/proto"
	"jbrodriguez/mediagui/server/lib"
	"jbrodriguez/mediagui/server/model"
	// "os"
	// "path/filepath"
	"regexp"
	// "strings"

	"golang.org/x/net/context"
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
	s.registerAdditional(s.bus, "/event/config/changed", s.configChanged, s.mailbox)

	re := []string{
		`(?i)(.*?):/mnt/user/films/(?P<Resolution>.*?)/(?P<Name>.*?)\s\((?P<Year>\d\d\d\d)\)/(?:.*/)*bdmv/index.(?P<FileType>bdmv)$`,
		`(?i)(.*?):/mnt/user/films/(?P<Resolution>.*?)/(?P<Name>.*?)\s\((?P<Year>\d\d\d\d)\)/(?:.*/)*.*\.(?P<FileType>iso|img|nrg|mkv|avi|xvid|ts|mpg|dvr-ms|mdf|wmv)$`,
		`(?i)(.*?):/mnt/user/films/(?P<Resolution>.*?)/(?P<Name>.*?)\s\((?P<Year>\d\d\d\d)\)/(?:.*/)*(?:video_ts|hv000i01)\.(?P<FileType>ifo)$`,
		// `(?i)(?P<Name>.*?)\s\((?P<Year>\d\d\d\d)\)/(?:.*/)*bdmv/index.(?P<FileType>bdmv)$`,
		// `(?i)(?P<Name>.*?)\s\((?P<Year>\d\d\d\d)\)/(?:.*/)*.*\.(?P<FileType>iso|img|nrg|mkv|avi|xvid|ts|mpg|dvr-ms|mdf|wmv)$`,
		// `(?i)(?P<Name>.*?)\s\((?P<Year>\d\d\d\d)\)/(?:.*/)*(?:video_ts|hv000i01)\.(?P<FileType>ifo)$`,
		// `(?i)(?P<Name>.*?)\s\((?P<Year>\d\d\d\d)\)\.(?P<FileType>iso|img|nrg|mkv|avi|xvid|ts|mpg|dvr-ms|mdf|wmv)$`,
	}

	s.re = make([]*lib.Rexp, 0)
	for _, regex := range re {
		s.re = append(s.re, &lib.Rexp{Exp: regexp.MustCompile(regex)})
	}

	s.includedMask = ".bdmv|.iso|.img|.nrg|.mkv|.avi|.xvid|.ts|.mpg|.dvr-ms|.mdf|.wmv|.ifo"

	cmd.Init()

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
	// folders := []string{
	// 	"/Volumes/hal-films",
	// 	"/Volumes/wopr-films",
	// }

	// ping := "ping -c1 %s > /dev/null && echo \"YES\" || echo \"NO\""

	lib.Notify(s.bus, "import:begin", "Import process started")

	// for _, folder := range s.settings.MediaFolders {
	// 	err := s.walk(folder)
	// 	if err != nil {
	// 		mlog.Warning("Unable to scan folder (%s): %s", folder, err)
	// 	}
	// }

	folders := []string{
		`/mnt/user/films`,
	}

	// Create new request to service go.micro.srv.example, method Example.Call
	req := client.NewRequest("io.jbrodriguez.mediagui.scanner.wopr", "Scanner.Scan", &scan.Request{
		// Folders: s.settings.MediaFolders,
		Folders: folders,
	})

	rsp := &scan.Response{}

	// Call service
	if err := client.Call(context.Background(), req, rsp); err != nil {
		lib.Notify(s.bus, "import:progress", "Unable to connect to scanning service")
		lib.Notify(s.bus, "import:end", "Import process finished")
		return
	}

	s.analyze(rsp.Filenames)

	lib.Notify(s.bus, "import:end", "Import process finished")
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

			movie := &model.Movie{Title: rmap["Name"], File_Title: rmap["Name"], Year: rmap["Year"], Resolution: resolution, FileType: rmap["FileType"], Location: file}
			// mlog.Info("FOUND [%s] (%s)", movie.Title, movie.Location)

			// mlog.Info("found %s", movie.Location)
			msg := &pubsub.Message{Payload: movie}
			s.bus.Pub(msg, "/event/movie/found")
			// self.Bus.MovieFound <- movie
		}
	}
}

// func (s *Scanner) walk(folder string) error {
// 	if folder[len(folder)-1] != '/' {
// 		folder = folder + "/"
// 	}

// 	err := walk.Walk(folder, func(path string, f os.FileInfo, err error) error {
// 		if err != nil {
// 			mlog.Info("scanner.walk: %s (%s) - [%+v]", err, path, f)
// 		}

// 		if !strings.Contains(s.includedMask, strings.ToLower(filepath.Ext(path))) {
// 			// mlog.Info("[%s] excluding %s", filepath.Ext(path), path)
// 			return nil
// 		}

// 		comparePath := strings.TrimPrefix(path, folder)
// 		// mlog.Info("folder: %s, comparePath: %s", folder, comparePath)
// 		// TODO: remove folder from path to match against rexp

// 		for i := 0; i < 3; i++ {
// 			var rmap = s.re[i].Match(comparePath)
// 			if rmap == nil {
// 				continue
// 			}

// 			var resolution string
// 			var ok bool
// 			if resolution, ok = rmap["Resolution"]; !ok {
// 				resolution = ""
// 			}

// 			movie := &model.Movie{Title: rmap["Name"], File_Title: rmap["Name"], Year: rmap["Year"], Resolution: resolution, FileType: rmap["FileType"], Location: path}
// 			// mlog.Info("FOUND [%s] (%s)", movie.Title, movie.Location)

// 			msg := &pubsub.Message{Payload: movie}
// 			s.bus.Pub(msg, "/event/movie/found")
// 			// self.Bus.MovieFound <- movie

// 			return nil
// 		}

// 		return nil
// 	})

// 	return err
// }

func (s *Scanner) configChanged(msg *pubsub.Message) {
	s.settings = msg.Payload.(*lib.Settings)
}
