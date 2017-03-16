package services

import (
	"jbrodriguez/mediagui/server/src/dto"
	"jbrodriguez/mediagui/server/src/lib"
	"jbrodriguez/mediagui/server/src/model"
	"jbrodriguez/mediagui/server/src/net"

	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"golang.org/x/net/websocket"

	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
)

const (
	apiVersion = "api/v1"
	capacity   = 3
	bufferSize = 8192
)

// Server -
type Server struct {
	Service

	bus      *pubsub.PubSub
	settings *lib.Settings
	router   *echo.Echo
	mailbox  chan *pubsub.Mailbox

	pool map[*net.Connection]bool
}

// NewServer -
func NewServer(bus *pubsub.PubSub, settings *lib.Settings) *Server {
	server := &Server{
		bus:      bus,
		settings: settings,
		pool:     make(map[*net.Connection]bool),
	}
	server.init()
	return server
}

// Start -
func (s *Server) Start() {
	mlog.Info("Starting service Server ...")

	// html := filepath.Join(s.settings.WebDir, "index.html")
	// mlog.Info("html is %s", html)
	// if b, _ := lib.Exists(html); !b {
	// 	mlog.Fatalf("Looked for index.html in %s, but didn't find it", s.settings.WebDir)
	// }

	// mlog.Info("Serving files from %s", s.settings.WebDir)

	// gin.SetMode(s.settings.GinMode)

	locations := []string{
		s.settings.WebDir,
		"/usr/local/share/mediagui",
		".",
	}

	location := lib.SearchFile("index.html", locations)
	if location == "" {
		msg := ""
		for _, loc := range locations {
			msg += fmt.Sprintf("%s, ", loc)
		}
		mlog.Fatalf("Unable to find index.html. Exiting now. (searched in %s)", msg)
	}

	mlog.Info("Serving files from %s", location)

	s.router = echo.New()
	s.router.Use(mw.Logger())
	s.router.Use(mw.Recover())

	s.router.Static("/", filepath.Join(location, "index.html"))
	s.router.Static("/img", filepath.Join(location, "img"))
	s.router.Static("/js", filepath.Join(location, "js"))
	s.router.Static("/css", filepath.Join(location, "css"))
	s.router.Static("/fonts", filepath.Join(location, "fonts"))

	s.router.GET("/ws", echo.WrapHandler(websocket.Handler(s.handleWs)))

	// s.router.GET("/", s.index)
	// s.router.GET("/ws", s.handleSocket)

	// s.router.Static("/app", filepath.Join(s.settings.WebDir, "app"))
	// s.router.Static("/img", filepath.Join(s.settings.WebDir, "img"))
	// s.router.Static("/js", filepath.Join(s.settings.WebDir, "js"))
	// s.router.Static("/css", filepath.Join(s.settings.WebDir, "css"))

	api := s.router.Group(apiVersion)
	api.GET("/config", s.getConfig)
	api.GET("/movies/single/:id", s.getMovie)
	api.GET("/movies/cover", s.getMoviesCover)
	api.GET("/movies", s.getMovies)
	api.GET("/movies/duplicates", s.getDuplicates)

	api.POST("/import", s.importMovies)
	api.POST("/prune", s.pruneMovies)

	api.PUT("/config/folder", s.addMediaFolder)
	api.PUT("/movies/:id/score", s.setMovieScore)
	api.PUT("/movies/:id/watched", s.setMovieWatched)
	api.PUT("/movies/:id/fix", s.fixMovie)

	// api := s.router.Group(apiVersion)
	// {
	// 	api.GET("/config", s.getConfig)
	// 	api.GET("/movies/single/:id", s.getMovie)
	// 	api.GET("/movies/cover", s.getMoviesCover)
	// 	api.GET("/movies", s.getMovies)
	// 	api.GET("/movies/duplicates", s.getDuplicates)

	// 	api.POST("/import", s.importMovies)
	// 	api.POST("/prune", s.pruneMovies)

	// 	api.PUT("/config/folder", s.addMediaFolder)
	// 	api.PUT("/movies/:id/score", s.setMovieScore)
	// 	api.PUT("/movies/:id/watched", s.setMovieWatched)
	// 	api.PUT("/movies/:id/fix", s.fixMovie)
	// }

	port := ":7623"
	go s.router.Start(port)

	s.mailbox = s.register(s.bus, "socket:broadcast", s.broadcast)
	go s.react()

	mlog.Info("Listening on %s", port)
}

// Stop -
func (s *Server) Stop() {
	mlog.Info("Stopped service Server ...")
}

func (s *Server) react() {
	for mbox := range s.mailbox {
		// mlog.Info("Core:Topic: %s", mbox.Topic)
		s.dispatch(mbox.Topic, mbox.Content)
	}
}

// // Closer -
// func Closer() gin.HandlerFunc {
// 	return func(c echo.Context) {
// 		c.Header("Connection", "close")
// 	}
// }

func (s *Server) getConfig(c echo.Context) error {
	msg := &pubsub.Message{Reply: make(chan interface{}, capacity)}
	s.bus.Pub(msg, "/get/config")

	reply := <-msg.Reply
	resp := reply.(*lib.Config)
	mlog.Info("config: %+v", resp)

	return c.JSON(http.StatusOK, &resp)
}

func (s *Server) getMoviesCover(c echo.Context) error {
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
	return c.JSON(http.StatusOK, dto)
}

func (s *Server) getMovies(c echo.Context) error {
	var options lib.Options
	c.Bind(&options) // You can also specify which binder to use. We support binding.Form, binding.JSON and binding.XML.

	mlog.Info("server.getMovies.options: %+v", options)
	// mlog.Info("request: ", c.Request)

	msg := &pubsub.Message{Payload: &options, Reply: make(chan interface{}, capacity)}
	s.bus.Pub(msg, "/get/movies")

	reply := <-msg.Reply
	dto := reply.(*model.MoviesDTO)

	// // mlog.Info("moviesDTO: %+v", dto)
	// return c.JSON(http.StatusOK, {dto})
	return c.JSON(http.StatusOK, dto)
}

func (s *Server) getDuplicates(c echo.Context) error {
	msg := &pubsub.Message{Reply: make(chan interface{}, capacity)}
	s.bus.Pub(msg, "/get/movies/duplicates")

	reply := <-msg.Reply
	dto := reply.(*model.MoviesDTO)

	return c.JSON(http.StatusOK, dto)
}

func (s *Server) getMovie(c echo.Context) error {
	id := c.Param("id")
	// mlog.Info("server.getMovies.options: %+v", options)
	// mlog.Info("request: ", c.Request)

	msg := &pubsub.Message{Payload: id, Reply: make(chan interface{}, capacity)}
	s.bus.Pub(msg, "/get/movie")

	reply := <-msg.Reply
	dto := reply.(*model.Movie)

	// // mlog.Info("moviesDTO: %+v", dto)
	// return c.JSON(http.StatusOK, {dto})
	return c.JSON(http.StatusOK, dto)
}

func (s *Server) importMovies(c echo.Context) error {
	s.bus.Pub(nil, "/post/import")

	return nil
}

func (s *Server) pruneMovies(c echo.Context) error {
	s.bus.Pub(nil, "/post/prune")

	return nil
}

func (s *Server) addMediaFolder(c echo.Context) error {
	var pkt dto.Packet
	if err := c.Bind(&pkt); err != nil {
		mlog.Warning("Unable to obtain folder: %s", err.Error())
	}

	msg := &pubsub.Message{Payload: pkt.Payload, Reply: make(chan interface{}, capacity)}
	s.bus.Pub(msg, "/put/config/folder")

	reply := <-msg.Reply
	resp := reply.(*lib.Config)

	return c.JSON(http.StatusOK, &resp)
}

func (s *Server) setMovieScore(c echo.Context) error {
	var movie model.Movie
	if err := c.Bind(&movie); err != nil {
		mlog.Warning("Unable to obtain setMovieScore: %s", err.Error())
		return c.JSON(http.StatusBadRequest, nil)
	}

	movie.Id, _ = strconv.ParseUint(c.Param("id"), 0, 64)

	msg := &pubsub.Message{Payload: &movie, Reply: make(chan interface{}, capacity)}
	s.bus.Pub(msg, "/put/movies/score")

	reply := <-msg.Reply
	resp := reply.(*model.Movie)

	return c.JSON(http.StatusOK, &resp)
}

func (s *Server) setMovieWatched(c echo.Context) error {
	var movie model.Movie
	if err := c.Bind(&movie); err != nil {
		mlog.Warning("Unable to obtain setMovieWatched: %s", err.Error())
		return c.JSON(http.StatusBadRequest, nil)
	}

	movie.Id, _ = strconv.ParseUint(c.Param("id"), 0, 64)

	msg := &pubsub.Message{Payload: &movie, Reply: make(chan interface{}, capacity)}
	s.bus.Pub(msg, "/put/movies/watched")

	reply := <-msg.Reply
	resp := reply.(*model.Movie)

	return c.JSON(http.StatusOK, &resp)
}

func (s *Server) fixMovie(c echo.Context) error {
	var movie model.Movie
	if err := c.Bind(&movie); err != nil {
		mlog.Warning("Unable to bind fixMovie body: %s", err.Error())
		return c.JSON(http.StatusBadRequest, nil)
	}

	movie.Id, _ = strconv.ParseUint(c.Param("id"), 0, 64)

	msg := &pubsub.Message{Payload: &movie, Reply: make(chan interface{}, capacity)}
	s.bus.Pub(msg, "/put/movies/fix")

	reply := <-msg.Reply
	resp := reply.(*model.Movie)

	return c.JSON(http.StatusOK, &resp)
}

func (s *Server) handleWs(ws *websocket.Conn) {
	conn := net.NewConnection(ws, s.onMessage, s.onClose)
	s.pool[conn] = true
	conn.Read()
}

func (s *Server) onMessage(packet *dto.Packet) {
	// mlog.Info("topic(%s)-payload(%+v)", packet.Topic, packet.Payload)
	s.bus.Pub(&pubsub.Message{Payload: packet.Payload}, packet.Topic)
}

func (s *Server) onClose(c *net.Connection, err error) {
	mlog.Warning("closing socket (%+v): %s", c, err)
	if _, ok := s.pool[c]; ok {
		delete(s.pool, c)
	}
}

func (s *Server) broadcast(msg *pubsub.Message) {
	packet := msg.Payload.(*dto.Packet)
	// mlog.Info("paylod-%+v", packet)
	for conn := range s.pool {
		conn.Write(packet)
	}
}
