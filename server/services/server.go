package services

import (
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	"github.com/julienschmidt/httprouter"
	"jbrodriguez/mediagui/server/lib"
	"net/http"
	"path/filepath"
)

type Server struct {
	bus      *pubsub.PubSub
	settings *lib.Settings
	router   *httprouter.Router
	// socket   *Socket
}

func NewServer(bus *pubsub.PubSub, settings *lib.Settings) *Server {
	server := &Server{bus: bus, settings: settings}
	return server
}

func (s *Server) Start() {
	mlog.Info("Starting service Server ...")

	html := filepath.Join(s.settings.WebDir, "index.html")
	if b, _ := lib.Exists(html); !b {
		mlog.Fatalf("Looked for index.html in %s, but didn't find it", s.settings.WebDir)
	}

	mlog.Info("Serving files from %s", s.settings.WebDir)

	s.router = httprouter.New()

	s.router.HandlerFunc("GET", "/", s.index)
	s.router.ServeFiles("/app/*filepath", http.Dir(filepath.Join(s.settings.WebDir, "app/")))

	port := ":7623"
	mlog.Info("Listening on %s", port)

	go http.ListenAndServe(port, s.router)
}

func (s *Server) Stop() {
	mlog.Info("Stopped service Server ...")
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join(s.settings.WebDir, "index.html"))
}
