package server

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"mediagui/domain"
	"mediagui/logger"
)

func (s *Server) setMovieWatched(c echo.Context) error {
	var dto domain.MovieDTO
	if err := c.Bind(&dto); err != nil {
		logger.Yellow("Unable to bind watchedMovie body: %s", err.Error())
		return c.JSON(http.StatusBadRequest, nil)
	}

	r := s.core.SetMovieWatched(&dto)

	return c.JSON(http.StatusOK, &r)
}
