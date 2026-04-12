APP_NAME := space-launch-server
BINARY := bin/server
MAIN_PACKAGE := ./cmd/server
DOCKER_IMAGE ?= $(APP_NAME):latest
SERVER_HOST ?= 10.0.10.1
SERVER_USER ?= root
SERVER_PORT ?= 22

.PHONY: help tidy fmt test build run dev air-install docker-build docker-export

help:
	@printf "Available targets:\n"
	@printf "  make tidy         - Run go mod tidy\n"
	@printf "  make fmt          - Format Go source files\n"
	@printf "  make test         - Run all tests\n"
	@printf "  make build        - Build binary to $(BINARY)\n"
	@printf "  make run          - Run server directly\n"
	@printf "  make air-install  - Install air hot reload tool\n"
	@printf "  make dev          - Run server with air\n"
	@printf "  make docker-build - Build Docker image ($(DOCKER_IMAGE))\n"
	@printf "  make docker-export - Build and transfer image to $(SERVER_USER)@$(SERVER_HOST)\n"

tidy:
	go mod tidy

fmt:
	go fmt ./...

test:
	go test ./...

build:
	mkdir -p bin
	go build -trimpath -o $(BINARY) $(MAIN_PACKAGE)

run:
	go run $(MAIN_PACKAGE)

air-install:
	go install github.com/air-verse/air@latest

dev:
	air -c .air.toml

docker-build:
	docker build -t $(DOCKER_IMAGE) .

docker-export: docker-build
	docker save $(DOCKER_IMAGE) | gzip | ssh -p $(SERVER_PORT) $(SERVER_USER)@$(SERVER_HOST) 'gunzip | docker load'
