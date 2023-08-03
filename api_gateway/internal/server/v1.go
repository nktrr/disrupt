package server

import (
	"context"
	"disrupt/api_gateway/config"
	kafka2 "disrupt/pkg/kafka"
	"disrupt/pkg/logger"
	"github.com/labstack/echo"
	"go.opentelemetry.io/otel"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"

	//"github.com/segmentio/kafka-go"
	"log"
	"net/http"
	"time"
)

type handlers struct {
	group          *echo.Group
	log            logger.Logger
	cfg            *config.Config
	producer       *kafka2.Producer
	tracerProvider *tracesdk.TracerProvider
	//tracer         trace.Tracer
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
		println("parse ")
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		defer func(ctx context.Context) {
			// Do not make the application hang when it is shutdown.
			ctx, cancel = context.WithTimeout(ctx, time.Second*5)
			defer cancel()
			if err := h.tracerProvider.Shutdown(ctx); err != nil {
				log.Fatal(err)
			}
		}(ctx)
		println("defer")
		tracer := h.tracerProvider.Tracer("ParseGithub")
		ctx, span := tracer.Start(ctx, "foo")
		println("ctx", ctx)
		println("tracer span")
		defer span.End()
		profile := c.QueryParam("profile")
		repo := c.QueryParam("repo")

		if profile == "" || repo == "" {
			return c.String(400, "no profile/repo")
		}
		Test(ctx)
		err := h.producer.WriteMessage(profile + "/" + repo)
		if err != nil {
			return c.String(400, err.Error())
		}
		return c.String(http.StatusOK, "PARSE OK")
	}
}

func Test(ctx context.Context) {
	tr := otel.Tracer("component-test")
	_, span := tr.Start(ctx, "test")
	defer span.End()
}
