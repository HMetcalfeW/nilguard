MODULE := github.com/your-org/nilguard

BIN_DIR := bin
CLI_BIN := $(BIN_DIR)/nilguard
VET_BIN := $(BIN_DIR)/nilguard-vet
PLUGIN  := $(BIN_DIR)/nilguard.so

.PHONY: all build cli vettool plugin test lint tidy ci clean

all: build

build: cli vettool

cli:
	go build -o $(CLI_BIN) ./pkg/nilguard

vettool:
	go build -o $(VET_BIN) ./pkg/vettool

plugin:
	GOFLAGS=-buildmode=plugin go build -o $(PLUGIN) ./plugin

test:
	go test ./... -v

lint:
	golangci-lint run
	staticcheck ./...

tidy:
	go mod tidy

ci: tidy build test lint

clean:
	rm -rf $(BIN_DIR)
