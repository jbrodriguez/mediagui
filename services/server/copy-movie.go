package server

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"mediagui/domain"
	"mediagui/logger"
)

func (s *Server) copyMovie(c echo.Context) error {
	var dto domain.MovieDTO
	if err := c.Bind(&dto); err != nil {
		logger.Yellow("Unable to bind fixMovie body: %s", err.Error())
		return c.JSON(http.StatusBadRequest, nil)
	}

	m := s.core.CopyMovie(&dto)

	return c.JSON(http.StatusOK, &m)
}
