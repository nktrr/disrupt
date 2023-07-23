package v1

import (
	"github.com/labstack/echo"
	"net/http"
)

func (h *handlers) MapRoutes() {
	h.group.Any("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})
	h.group.GET("/parse", h.ParseGithub())
}
