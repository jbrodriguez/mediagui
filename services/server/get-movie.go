package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) getMovie(c echo.Context) error {
	id := c.Param("id")
	return c.JSON(http.StatusOK, s.core.GetMovie(id))
}
