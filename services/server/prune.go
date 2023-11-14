package server

import (
	"github.com/labstack/echo/v4"
)

func (s *Server) pruneMovies(c echo.Context) error {
	s.core.PruneMovies()
	return nil
}
