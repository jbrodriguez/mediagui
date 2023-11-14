package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) getCovers(c echo.Context) error {
	return c.JSON(http.StatusOK, s.core.GetCovers())
}
