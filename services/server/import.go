package server

import "github.com/labstack/echo/v4"

func (s *Server) importMovies(c echo.Context) error {
	s.core.ImportMovies()
	return nil
}
