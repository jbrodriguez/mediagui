package server

import (
	"mediagui/domain"
	"mediagui/logger"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) setMovieScore(c echo.Context) error {
	var dto domain.MovieDTO
	if err := c.Bind(&dto); err != nil {
		logger.Yellow("Unable to bind rateMovie body: %s", err.Error())
		return c.JSON(http.StatusBadRequest, nil)
	}

	r := s.core.SetMovieScore(&dto)

	return c.JSON(http.StatusOK, &r)
}
