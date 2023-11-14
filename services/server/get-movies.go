package server

import (
	"mediagui/domain"
	"mediagui/logger"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) getMovies(c echo.Context) error {
	var options domain.Options
	c.Bind(&options) // You can also specify which binder to use. We support binding.Form, binding.JSON and binding.XML.

	logger.Blue("server.getMovies.options: %+v", options)

	return c.JSON(http.StatusOK, s.core.GetMovies(&options))
}
