package server

import (
	"context"
	"disrupt/api_gateway/config"
	"disrupt/pkg/kafka"
	"disrupt/pkg/logger"
	"disrupt/pkg/tracing"
	"github.com/labstack/echo"
	"go.opentelemetry.io/otel"
	"log"
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
	producer, err := kafka.NewProducer(kafka.ParseGithub)
	if err != nil {
		log.Fatal("failed to create kafka producer:", err)
	}
	s.handlers = NewHandlers(s.echo.Group(s.cfg.Http.BasePath), s.log, s.cfg)
	s.handlers.MapRoutes()
	s.handlers.producer = producer
	defer s.handlers.producer.Close()

	tp, err := tracing.TracerProvider("http://jaeger:14268/api/traces", s.cfg.ServiceName)
	println("start tracing service:", s.cfg.ServiceName)
	if err != nil {
		s.log.Errorf("run tracing: %v", err)
	}
	s.handlers.tracerProvider = tp
	println(s.handlers.tracerProvider)
	otel.SetTracerProvider(s.handlers.tracerProvider)

	// ADD TRACING LOGIC IN HANDLERS
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
