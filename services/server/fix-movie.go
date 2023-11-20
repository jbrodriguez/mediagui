package server

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"mediagui/domain"
	"mediagui/logger"
)

func (s *Server) fixMovie(c echo.Context) error {
	var dto domain.FixMovieDTO
	if err := c.Bind(&dto); err != nil {
		logger.Yellow("Unable to bind fixMovie body: %s", err.Error())
		return c.JSON(http.StatusBadRequest, nil)
	}

	r := s.core.FixMovie(&dto)

	return c.JSON(http.StatusOK, &r)
}
