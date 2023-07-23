package v1

import (
	"disrupt/api_gateway/config"
	"disrupt/pkg/logger"
	"github.com/labstack/echo"
	"net/http"
)

type handlers struct {
	group *echo.Group
	log   logger.Logger
	cfg   *config.Config
}

func NewHandlers(group *echo.Group, log logger.Logger, cfg *config.Config) *handlers {
	return &handlers{
		group: group,
		log:   log,
		cfg:   cfg,
	}
}

func (h *handlers) ParseGithub() echo.HandlerFunc {
	return func(c echo.Context) error {

		// remove
		return c.String(http.StatusOK, "PARSE OK")
	}
}
