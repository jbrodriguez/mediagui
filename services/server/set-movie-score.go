package server

import (
	"mediagui/domain"
	"mediagui/logger"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (s *Server) setMovieScore(c echo.Context) error {
	var movie domain.Movie
	if err := c.Bind(&movie); err != nil {
		logger.Yellow("Unable to obtain setMovieScore: %s", err.Error())
		return c.JSON(http.StatusBadRequest, nil)
	}

	movie.ID, _ = strconv.ParseUint(c.Param("id"), 0, 64)

	r := s.core.SetMovieScore(&movie)

	return c.JSON(http.StatusOK, &r)
}
