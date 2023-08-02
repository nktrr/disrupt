package server

import (
	"context"
	"disrupt/api_gateway/config"
	kafka2 "disrupt/pkg/kafka"
	"disrupt/pkg/logger"
	"github.com/labstack/echo"
	"github.com/segmentio/kafka-go"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"

	//"github.com/segmentio/kafka-go"
	"log"
	"net/http"
	"time"
)

type handlers struct {
	group    *echo.Group
	log      logger.Logger
	cfg      *config.Config
	producer kafka2.Producer
	tracer   *tracesdk.TracerProvider
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
		profile := c.QueryParam("profile")
		repo := c.QueryParam("repo")

		if profile == "" || repo == "" {
			return c.String(400, "no profile/repo")
		}

		h.producer.PublishMessage(c)

		topic := "parse-github"
		partition := 0

		conn, err := kafka.DialLeader(context.Background(), "tcp", "kafka:29092", topic, partition)
		if err != nil {
			log.Fatal("failed to dial leader:", err)
		}

		conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		_, err = conn.WriteMessages(
			kafka.Message{Value: []byte(profile + "/" + repo)},
		)
		if err != nil {
			log.Fatal("failed to write messages:", err)
		}

		if err := conn.Close(); err != nil {
			log.Fatal("failed to close writer:", err)
		}
		// remove
		return c.String(http.StatusOK, "PARSE OK")
	}
}
