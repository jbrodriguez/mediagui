package services

import (
	"github.com/gin-gonic/gin"
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	"jbrodriguez/mediagui/server/lib"
	"jbrodriguez/mediagui/server/model"
	"path/filepath"
)

const apiVersion = "api/v1"
const capacity = 3

type Server struct {
	bus      *pubsub.PubSub
	settings *lib.Settings
	router   *gin.Engine
	// socket   *Socket
}

func NewServer(bus *pubsub.PubSub, settings *lib.Settings) *Server {
	server := &Server{bus: bus, settings: settings}
	return server
}

func (s *Server) Start() {
	mlog.Info("Starting service Server ...")

	html := filepath.Join(s.settings.WebDir, "index.html")
	mlog.Info("html is %s", html)
	if b, _ := lib.Exists(html); !b {
		mlog.Fatalf("Looked for index.html in %s, but didn't find it", s.settings.WebDir)
	}

	mlog.Info("Serving files from %s", s.settings.WebDir)

	s.router = gin.Default()

	s.router.GET("/", s.index)
	s.router.Static("/app", filepath.Join(s.settings.WebDir, "app"))
	s.router.Static("/img", filepath.Join(s.settings.WebDir, "img"))

	api := s.router.Group(apiVersion)
	{
		api.GET("/config", s.getConfig)
		api.GET("/movies/cover", s.getMoviesCover)
	}

	port := ":7623"
	go s.router.Run(port)
	mlog.Info("Listening on %s", port)
}

func (s *Server) Stop() {
	mlog.Info("Stopped service Server ...")
}

func (s *Server) index(c *gin.Context) {
	c.File(filepath.Join(s.settings.WebDir, "index.html"))
}

func (s *Server) getConfig(c *gin.Context) {
	msg := &pubsub.Message{Reply: make(chan interface{}, capacity)}
	s.bus.Pub(msg, "/get/config")

	reply := <-msg.Reply
	resp := reply.(*lib.Config)
	mlog.Info("config: %+v", resp)
	c.JSON(200, &resp)
}

func (s *Server) getMoviesCover(c *gin.Context) {
	msg := &pubsub.Message{Reply: make(chan interface{}, capacity)}
	s.bus.Pub(msg, "/get/movies/cover")

	reply := <-msg.Reply
	dto := reply.(*model.MoviesDTO)

	// movies := make([]*model.Movie, 0)
	// movies = append(
	// 	movies,
	// 	&model.Movie{Id: 1, Title: "The Godfather", Year: "1971"},
	// 	&model.Movie{Id: 2, Title: "Pulp Fiction", Year: "1990"},
	// )

	// dto := &model.MoviesDTO{
	// 	Total: 2,
	// 	Items: movies,
	// }

	mlog.Info("moviesDTO: %+v", dto)
	c.JSON(200, dto)
}
