package scanner

import (
	"context"
	"fmt"
	"regexp"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"mediagui/domain"
	"mediagui/lib"
	"mediagui/logger"
	pb "mediagui/mediaagent"
)

// Scanner -
type Scanner struct {
	ctx *domain.Context

	re           []*lib.Rexp
	includedMask string
}

func Create(ctx *domain.Context) *Scanner {
	scanner := &Scanner{
		ctx: ctx,
	}
	return scanner
}

// Start -
func (s *Scanner) Start() {
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

	logger.Blue("started service scanner ...")
}

func (s *Scanner) ScanMovies() {
	defer s.ctx.Hub.Pub(nil, "/event/workunit/done")

	for _, host := range s.ctx.UnraidHosts {
		address := fmt.Sprintf("%s:7624", host)

		conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			// logger.Yellow("Unable to connect to host (%s): %s", address, err)
			lib.Notify(s.ctx.Hub, "import:progress", "Unable to connect to host "+host)
			continue
		}
		defer conn.Close()

		client := pb.NewMediaAgentClient(conn)

		rsp, err := client.Scan(context.Background(), &pb.ScanReq{Folders: s.ctx.MediaFolders, Mask: s.includedMask})
		if err != nil {
			// logger.Yellow("Unable to scan (%s): %s", address, err)
			lib.Notify(s.ctx.Hub, "import:progress", "Unable to scan host "+host)
			continue
		}

		s.analyze(rsp.Filenames)
	}
}

func (s *Scanner) analyze(files []string) {
	logger.Blue("Found %d files", len(files))

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

			movie := &domain.Movie{Title: rmap["Name"], FileTitle: rmap["Name"], Year: rmap["Year"], Resolution: resolution, FileType: rmap["FileType"], Location: file, Stub: 0}

			s.ctx.Hub.Pub(movie, "/event/movie/found")
		}
	}
}

// func (s *Scanner) configChanged(msg *pubsub.Message) {
// 	s.settings = msg.Payload.(*lib.Settings)
// }
