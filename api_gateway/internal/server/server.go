package server

import (
	"context"
	"disrupt/api_gateway/config"
	"disrupt/pkg/kafka"
	"disrupt/pkg/logger"
	"github.com/labstack/echo"
	"os"
	"os/signal"
	"syscall"
)

type server struct {
	log      logger.Logger
	cfg      *config.Config
	echo     *echo.Echo
	handlers *handlers
}

func NewServer(log logger.Logger, cfg *config.Config) *server {
	return &server{
		log:  log,
		cfg:  cfg,
		echo: echo.New(),
	}
}

func (s *server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()
	kafkaProducer := kafka.NewProducer(s.log, s.cfg.Kafka.Brokers)
	defer kafkaProducer.Close()
	s.handlers = NewHandlers(s.echo.Group(s.cfg.Http.BasePath), s.log, s.cfg)
	s.handlers.MapRoutes()
	go func() {
		if err := s.echo.Start(s.cfg.Http.Port); err != nil {
			s.log.Errorf(" s.runHttpServer: %v", err)
			cancel()
		}
	}()
	s.log.Infof("API is listening on: %s", s.cfg.Http.Port)
	<-ctx.Done()
	if err := s.echo.Server.Shutdown(ctx); err != nil {
		s.log.WarnMsg("echo.Server.Shutdown", err)
	}
	return nil
}
