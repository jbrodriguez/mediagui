package server

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"mediagui/domain"
	"mediagui/logger"
)

func (s *Server) copyMovie(c echo.Context) error {
	var movie domain.Movie
	if err := c.Bind(&movie); err != nil {
		logger.Yellow("Unable to bind copyMovie body: %s", err.Error())
		return c.JSON(http.StatusBadRequest, nil)
	}

	movie.ID, _ = strconv.ParseUint(c.Param("id"), 0, 64)

	m := s.core.CopyMovie(&movie)

	return c.JSON(http.StatusOK, &m)
}
