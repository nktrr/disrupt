FROM golang:1.21rc3-alpine3.18

WORKDIR /app

ENV CONFIG=docker

COPY ./ ./


RUN go install -mod=mod github.com/githubnemo/CompileDaemon
RUN go mod download


ENTRYPOINT CompileDaemon --build="go build github_reader/cmd/main.go" --command=./main

