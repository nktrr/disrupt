package kafka

import (
	"context"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
)

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(topic string) *Consumer {
	partition := 0
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{"kafka:29092"},
		StartOffset: kafka.LastOffset,
		Topic:       topic,
		Partition:   partition,
		MaxBytes:    10e6, //10MB
	})
	consumer := &Consumer{
		reader: reader,
	}
	return consumer
}

func (consumer *Consumer) Consume() {
	for {
		println("start read")
		m, err := consumer.reader.ReadMessage(context.Background())
		if err != nil {
			println(err.Error())
			break
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}
	if err := consumer.reader.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}
