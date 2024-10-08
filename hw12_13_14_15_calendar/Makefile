BIN_SRV := "./bin/server"
BIN_SCHD := "./bin/scheduler"
BIN_SND := "./bin/sender"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN_SRV) -ldflags "$(LDFLAGS)" ./app/service/cmd/
	go build -v -o $(BIN_SCHD) -ldflags "$(LDFLAGS)" ./app/service/scheduler/
	go build -v -o $(BIN_SND) -ldflags "$(LDFLAGS)" ./app/service/sender/

run: build
	$(BIN_SRV) -config config.yaml &
	$(BIN_SCHD) -config config.yaml &
	$(BIN_SND) -config config.yaml

build-img:
	docker compose build

up: build-img
	docker compose up -d

down:
	docker compose down

version: build
	$(BIN_SRV) version

test:
	go test -race ./app/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.59.1

lint: install-lint-deps
	golangci-lint run ./...

lint-fix:
	go mod tidy
	gofumpt -l -w .
	golangci-lint run --fix

PROTO_DIR = ./app/service/proto
PROTO_FILES = $(PROTO_DIR)/event.proto
GENERATED_DIR = ./app/service/proto

generate:
	protoc --go_out=$(GENERATED_DIR) --go-grpc_out=$(GENERATED_DIR) $(PROTO_FILES)

mockgen:
	mockgen -source=app/lib/storage/storage.go -destination=app/lib/storage/mock_storage.go -package=storage

integration-tests:
	# Run the integration tests inside a Docker container
	docker compose -f docker-compose.integration.yml up --build --abort-on-container-exit \

	# Clean up the environment
	docker compose -f docker-compose.integration.yml down

DOCKER_REGISTRY := docker.io
DOCKER_USER := thejubadze
DOCKER_PATH := $(DOCKER_REGISTRY)/$(DOCKER_USER)
push-img: build-img
	docker push $(DOCKER_PATH)/calendar
	docker push $(DOCKER_PATH)/calendar-scheduler
	docker push $(DOCKER_PATH)/calendar-sender

.PHONY: build run build-img up down version test lint lint-fix generate mockgen integration-tests push-img
