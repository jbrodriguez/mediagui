package services

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"golang.org/x/net/websocket"

	"github.com/jbrodriguez/actor"
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"

	"mediagui/dto"
	"mediagui/lib"
	"mediagui/model"
	"mediagui/net"
)

const (
	apiVersion = "api/v1"
	capacity   = 3
	// bufferSize = 8192
)

// Server -
type Server struct {
	bus      *pubsub.PubSub
	settings *lib.Settings
	router   *echo.Echo
	actor    *actor.Actor

	pool map[*net.Connection]bool
}

// NewServer -
func NewServer(bus *pubsub.PubSub, settings *lib.Settings) *Server {
	server := &Server{
		bus:      bus,
		settings: settings,
		actor:    actor.NewActor(bus),
		pool:     make(map[*net.Connection]bool),
	}
	return server
}

// Start -
func (s *Server) Start() {
	mlog.Info("Starting service Server ...")

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

	s.router.HideBanner = true

	s.router.Use(mw.Logger())
	s.router.Use(mw.Recover())
	s.router.Use(mw.CORS())
	s.router.Use(mw.Gzip())

	s.router.Static("/", filepath.Join(location, "index.html"))
	s.router.Static("/img", filepath.Join(location, "img"))
	s.router.Static("/js", filepath.Join(location, "js"))
	s.router.Static("/css", filepath.Join(location, "css"))
	s.router.Static("/fonts", filepath.Join(location, "fonts"))

	s.router.GET("/ws", echo.WrapHandler(websocket.Handler(s.handleWs)))

	api := s.router.Group(apiVersion)
	api.GET("/config", s.getConfig)
	api.GET("/movies/single/:id", s.getMovie)
	api.GET("/movies/cover", s.getMoviesCover)
	api.GET("/movies", s.getMovies)
	api.GET("/movies/duplicates", s.getDuplicates)

	api.POST("/import", s.importMovies)
	api.POST("/add", s.addMovie)
	api.POST("/prune", s.pruneMovies)

	api.PUT("/config/folder", s.addMediaFolder)
	api.PUT("/movies/:id/score", s.setMovieScore)
	api.PUT("/movies/:id/watched", s.setMovieWatched)
	api.PUT("/movies/:id/fix", s.fixMovie)
	api.PUT("/movies/:id/duplicate", s.setDuplicate)

	port := ":7623"
	go s.router.Start(port)

	s.actor.Register("socket:broadcast", s.broadcast)
	go s.actor.React()

	mlog.Info("Listening on %s", port)
}

// Stop -
func (s *Server) Stop() {
	mlog.Info("Stopped service Server ...")
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
	movie := reply.(*model.MoviesDTO)
	return c.JSON(http.StatusOK, movie)
}

func (s *Server) getMovies(c echo.Context) error {
	var options lib.Options
	c.Bind(&options) // You can also specify which binder to use. We support binding.Form, binding.JSON and binding.XML.

	mlog.Info("server.getMovies.options: %+v", options)

	msg := &pubsub.Message{Payload: &options, Reply: make(chan interface{}, capacity)}
	s.bus.Pub(msg, "/get/movies")

	reply := <-msg.Reply
	movie := reply.(*model.MoviesDTO)

	// // mlog.Info("moviesDTO: %+v", dto)
	// return c.JSON(http.StatusOK, {dto})
	return c.JSON(http.StatusOK, movie)
}

func (s *Server) getDuplicates(c echo.Context) error {
	msg := &pubsub.Message{Reply: make(chan interface{}, capacity)}
	s.bus.Pub(msg, "/get/movies/duplicates")

	reply := <-msg.Reply
	movie := reply.(*model.MoviesDTO)

	return c.JSON(http.StatusOK, movie)
}

func (s *Server) getMovie(c echo.Context) error {
	id := c.Param("id")

	msg := &pubsub.Message{Payload: id, Reply: make(chan interface{}, capacity)}
	s.bus.Pub(msg, "/get/movie")

	reply := <-msg.Reply
	movie := reply.(*model.Movie)
	return c.JSON(http.StatusOK, movie)
}

func (s *Server) importMovies(_ echo.Context) error {
	s.bus.Pub(nil, "/post/import")

	return nil
}

func (s *Server) addMovie(c echo.Context) error {
	var movie model.Movie
	if err := c.Bind(&movie); err != nil {
		mlog.Warning("Unable to bind addMovie body: %s", err.Error())
		return c.JSON(http.StatusBadRequest, nil)
	}

	msg := &pubsub.Message{Payload: &movie, Reply: make(chan interface{}, capacity)}
	s.bus.Pub(msg, "/post/add")

	reply := <-msg.Reply
	resp := reply.(*model.Movie)

	return c.JSON(http.StatusOK, &resp)
}

func (s *Server) pruneMovies(_ echo.Context) error {
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

	movie.ID, _ = strconv.ParseUint(c.Param("id"), 0, 64)

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

	movie.ID, _ = strconv.ParseUint(c.Param("id"), 0, 64)

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

	movie.ID, _ = strconv.ParseUint(c.Param("id"), 0, 64)

	msg := &pubsub.Message{Payload: &movie, Reply: make(chan interface{}, capacity)}
	s.bus.Pub(msg, "/put/movies/fix")

	reply := <-msg.Reply
	resp := reply.(*model.Movie)

	return c.JSON(http.StatusOK, &resp)
}

func (s *Server) setDuplicate(c echo.Context) error {
	var movie model.Movie
	if err := c.Bind(&movie); err != nil {
		mlog.Warning("Unable to bind setDuplicate body: %s", err.Error())
		return c.JSON(http.StatusBadRequest, nil)
	}

	movie.ID, _ = strconv.ParseUint(c.Param("id"), 0, 64)

	msg := &pubsub.Message{Payload: &movie, Reply: make(chan interface{}, capacity)}
	s.bus.Pub(msg, "/put/movies/duplicate")

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
	s.bus.Pub(&pubsub.Message{Payload: packet.Payload}, packet.Topic)
}

func (s *Server) onClose(c *net.Connection, err error) {
	mlog.Warning("closing socket (%+v): %s", c, err)
	delete(s.pool, c)
}

func (s *Server) broadcast(msg *pubsub.Message) {
	packet := msg.Payload.(*dto.Packet)
	for conn := range s.pool {
		conn.Write(packet)
	}
}
