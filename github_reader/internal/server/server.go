package server

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
)

type server struct {
}

func NewServer() *server {
	return &server{}
}

func (s *server) Run() {
	println("run")
	topic := "parse-github"
	partition := 0

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"kafka:29092"},
		Topic:     topic,
		Partition: partition,
		MaxBytes:  10e6, // 10MB
	})
	//r.SetOffset(42)

	for {
		println("start read")

		m, err := r.ReadMessage(context.Background())
		println("read")
		if err != nil {
			println(err.Error())
			break
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}
