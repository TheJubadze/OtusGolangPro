BIN_SRV := "./bin/server"
BIN_SCHD := "./bin/scheduler"
BIN_SND := "./bin/sender"

DOCKER_IMG="calendar:develop"

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
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

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

PROTO_DIR = ./proto
PROTO_FILES = $(PROTO_DIR)/event.proto
GENERATED_DIR = ./proto

generate:
	protoc --go_out=$(GENERATED_DIR) --go-grpc_out=$(GENERATED_DIR) $(PROTO_FILES)

mockgen:
	mockgen -source=app/lib/storage/storage.go -destination=app/lib/storage/mock_storage.go -package=storage

.PHONY: build run build-img run-img version test lint
