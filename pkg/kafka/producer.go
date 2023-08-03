package kafka

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

type Producer struct {
	conn *kafka.Conn
}

func NewProducer(topic string) (*Producer, error) {
	partition := 0
	conn, err := kafka.DialLeader(context.Background(), "tcp", "kafka:29092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
		return nil, err
	}
	err = conn.SetWriteDeadline(time.Time{})
	if err != nil {
		return nil, err
	}
	producer := &Producer{conn: conn}
	return producer, err
}

func (producer Producer) WriteMessage(message string) error {
	_, err := producer.conn.Write([]byte(message))
	if err != nil {
		return err
	}
	return nil
}

func (producer Producer) Close() {
	err := producer.conn.Close()
	if err != nil {
		return
	}
}

const (
	ParseGithub = "parse-github"
)
