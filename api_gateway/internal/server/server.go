package server

import (
	"disrupt/api_gateway/config"
	"log"
)

type server struct {
	log  log.Logger
	cfg  *config.Config
	echo *echo.Echo
}
