#include .env

LOCAL_BIN:=$(CURDIR)/bin

LOCAL_MIGRATION_DIR=$(MIGRATION_DIR)
LOCAL_MIGRATION_DSN="host=localhost port=$(PG_PORT) user=$(PG_USER) password=$(PG_PASSWORD) dbname=$(PG_DATABASE_NAME) sslmode=disable"

install-deps:
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@latest

local-migration-status:
	${LOCAL_BIN}/goose -dir $(LOCAL_MIGRATION_DIR) postgres  $(LOCAL_MIGRATION_DSN) status -v

local-migration-up:
	${LOCAL_BIN}/goose -dir $(LOCAL_MIGRATION_DIR) postgres  $(LOCAL_MIGRATION_DSN) up -v

local-migration-down:
	${LOCAL_BIN}/goose -dir $(LOCAL_MIGRATION_DIR) postgres  $(LOCAL_MIGRATION_DSN) down -v


BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) -config ./configs/config.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/... ./pkg/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.57.2

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img version test lint


