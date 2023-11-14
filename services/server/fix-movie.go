package server

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"mediagui/domain"
	"mediagui/logger"
)

func (s *Server) fixMovie(c echo.Context) error {
	var movie domain.Movie
	if err := c.Bind(&movie); err != nil {
		logger.Yellow("Unable to bind fixMovie body: %s", err.Error())
		return c.JSON(http.StatusBadRequest, nil)
	}

	movie.ID, _ = strconv.ParseUint(c.Param("id"), 0, 64)

	r := s.core.FixMovie(&movie)

	return c.JSON(http.StatusOK, &r)
}
