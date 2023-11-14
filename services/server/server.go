package server

import (
	"embed"
	"io/fs"
	"mediagui/domain"
	"mediagui/logger"
	"mediagui/services/core"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	web "mediagui/ui"
)

const (
	apiVersion = "api/v1"
	capacity   = 3
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Server struct {
	ctx    *domain.Context
	engine *echo.Echo
	core   *core.Core
	ws     *websocket.Conn
}

func Create(ctx *domain.Context, core *core.Core) *Server {
	return &Server{
		ctx:  ctx,
		core: core,
	}
}

func (s *Server) Start() error {
	// s.ctx.Hub.Sub("socket:broadcast", s.onBroadcast)

	s.engine = echo.New()

	s.engine.HideBanner = true

	s.engine.Use(middleware.Recover())
	s.engine.Use(middleware.CORS())
	// s.engine.Use(middleware.Logger())
	s.engine.Use(middleware.Gzip())

	// Define a "/" endpoint to serve index.html from the embed FS
	s.engine.GET("/*", indexHandler)

	s.engine.GET("/assets/*", echo.WrapHandler(assetsHandler(web.Dist)))
	s.engine.Static("/img/*", filepath.Join(s.ctx.DataDir, "img"))

	s.engine.GET("/ws", s.wsHandler)

	api := s.engine.Group(apiVersion)
	api.GET("/config", s.getConfig)
	api.GET("/movies/:id", s.getMovie) // prev api.GET("/movies/single/:id", s.getMovie)

	api.GET("/covers", s.getCovers) // prev api.GET("/movies/cover", s.getMoviesCover)
	api.GET("/movies", s.getMovies)
	api.GET("/duplicates", s.getDuplicates) // prev api.GET("/movies/duplicates", s.getDuplicates)

	api.POST("/import", s.importMovies)
	// api.POST("/add", s.addMovie)
	api.POST("/prune", s.pruneMovies)

	// api.PUT("/config/folder", s.addMediaFolder)
	api.PUT("/movies/:id/score", s.setMovieScore)
	api.PUT("/movies/:id/watched", s.setMovieWatched)
	api.PUT("/movies/:id/fix", s.fixMovie)
	api.PUT("/movies/:id/copy", s.copyMovie)
	api.PUT("/movies/:id/duplicate", s.setDuplicate)

	// Always listen on http port, but based on above setting, we could be redirecting to https
	go func() {
		err := s.engine.Start(":7623")
		if err != nil {
			logger.Red("unable to start http server: %s", err)
			os.Exit(1)
		}
	}()

	logger.Blue("started service server (listening http on :7623) ...")

	return nil
}

func indexHandler(c echo.Context) error {
	data, err := web.Dist.ReadFile("dist/index.html")
	if err != nil {
		return err
	}
	return c.Blob(http.StatusOK, "text/html", data)
}

func assetsHandler(content embed.FS) http.Handler {
	fsys, err := fs.Sub(content, "dist")
	if err != nil {
		panic(err)
	}
	return http.FileServer(http.FS(fsys))
}

func (s *Server) wsHandler(c echo.Context) error {
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		logger.Red("unable to upgrade websocket: %s", err)
		return err
	}
	defer conn.Close()

	s.ws = conn

	return s.wsRead()
}

func (s *Server) wsRead() (err error) {
	for {
		var packet domain.Packet
		err = s.ws.ReadJSON(&packet)
		if err != nil {
			logger.Red("unable to read websocket message: %s", err)
			return err
		}

		s.ctx.Hub.Pub(packet.Payload, packet.Topic)
	}
}

func (s *Server) wsWrite(packet *domain.Packet) (err error) {
	err = s.ws.WriteJSON(packet)
	return
}

func (s *Server) onBroadcast(msg any) {
	message := msg.(*domain.Packet)

	packet := &domain.Packet{
		Topic:   message.Topic,
		Payload: message.Payload,
	}

	err := s.wsWrite(packet)
	if err != nil {
		logger.Red("unable to write websocket message: %s", err)
	}
}
