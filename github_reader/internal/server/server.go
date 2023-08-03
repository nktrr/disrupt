package server

import (
	kafka2 "disrupt/pkg/kafka"
)

type server struct {
}

func NewServer() *server {
	return &server{}
}

func (s *server) Run() {
	println("run")
	consumer := kafka2.NewConsumer(kafka2.ParseGithub)
	consumer.Consume()
}
