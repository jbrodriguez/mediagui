package services

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	"jbrodriguez/mediagui/server/lib"
	"jbrodriguez/mediagui/server/model"
	"net/http"
	"path/filepath"
)

const (
	apiVersion = "api/v1"
	capacity   = 3
	bufferSize = 8192
)

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
	s.router.GET("/ws", s.handleSocket)
	s.router.Static("/app", filepath.Join(s.settings.WebDir, "app"))
	s.router.Static("/img", filepath.Join(s.settings.WebDir, "img"))

	api := s.router.Group(apiVersion)
	{
		api.GET("/config", s.getConfig)
		api.GET("/movies/cover", s.getMoviesCover)
		api.GET("/movies", s.getMovies)
		api.POST("/import", s.importMovies)
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

func (s *Server) handleSocket(c *gin.Context) {
	ws, err := websocket.Upgrade(c.Writer, c.Request, nil, bufferSize, bufferSize)
	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(c.Writer, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		mlog.Error(err)
		return
	}

	msg := &pubsub.Message{Payload: ws}
	s.bus.Pub(msg, "socket:connections:new")
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

	// mlog.Info("moviesDTO: %+v", dto)
	c.JSON(200, dto)
}

func (s *Server) getMovies(c *gin.Context) {
	var options lib.Options
	c.Bind(&options) // You can also specify which binder to use. We support binding.Form, binding.JSON and binding.XML.

	mlog.Info("server.getMovies.options: %+v", options)

	msg := &pubsub.Message{Payload: options, Reply: make(chan interface{}, capacity)}
	s.bus.Pub(msg, "/get/movies")

	reply := <-msg.Reply
	dto := reply.(*model.MoviesDTO)

	// // mlog.Info("moviesDTO: %+v", dto)
	// c.JSON(200, {dto})
	c.JSON(200, dto)
}

func (s *Server) importMovies(c *gin.Context) {
	s.bus.Pub(nil, "/post/import")
}
