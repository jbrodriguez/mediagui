package server

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (s *Server) getMovie(c echo.Context) error {
	sid := c.Param("id")

	id, err := strconv.ParseUint(sid, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, nil)
	}

	return c.JSON(http.StatusOK, s.core.GetMovie(id))
}
