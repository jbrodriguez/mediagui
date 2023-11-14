package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) getConfig(c echo.Context) error {
	return c.JSON(http.StatusOK, s.core.GetConfig())
}
