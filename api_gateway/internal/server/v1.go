package server

import (
	"context"
	"disrupt/api_gateway/config"
	"disrupt/pkg/logger"
	"github.com/labstack/echo"
	"github.com/segmentio/kafka-go"
	"log"
	"net/http"
	"strconv"
	"time"
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
		println("here")
		profile := "nktrr"
		repository := "disrupt_old"

		topic := "parse-github"
		partition := 0

		conn, err := kafka.DialLeader(context.Background(), "tcp", "kafka:29092", topic, partition)
		if err != nil {
			log.Fatal("failed to dial leader:", err)
		}

		conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
		i, err := conn.WriteMessages(
			kafka.Message{Value: []byte(profile + "/" + repository)},
		)
		println("write ", strconv.Itoa(i))
		if err != nil {
			println("write err")
			log.Fatal("failed to write messages:", err)
		}

		if err := conn.Close(); err != nil {
			println("close writter")
			log.Fatal("failed to close writer:", err)
		}
		// remove
		return c.String(http.StatusOK, "PARSE OK")
	}
}
