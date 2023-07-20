package cmd

import (
	"disrupt/api_gateway/config"
	"disrupt/pkg/logger"
	"flag"
	"log"
)

func main() {
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
