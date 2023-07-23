package main

import (
	"disrupt/api_gateway/config"
	"disrupt/api_gateway/internal/server"
	"disrupt/pkg/logger"
	"flag"
	"log"
)

func main() {
	println("start")
	flag.Parse()

	cfg, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	appLogger := logger.NewAppLogger(cfg.Logger)
	appLogger.InitLogger()
	appLogger.WithName("ApiGateway")
	s := server.NewServer(appLogger, cfg)
	appLogger.Fatal(s.Run())
}
