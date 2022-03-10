BIN_API := "./bin/calendar_api"
BIN_SCHEDULER := "./bin/calendar_scheduler"
BIN_SENDER := "./bin/calendar_scheduler"
DOCKER_API_IMG="calendarapi:develop"
REPO=github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X '${REPO}/cmd/shared/cmd.Release=develop' -X '${REPO}/cmd/shared/cmd.BuildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S)' -X '${REPO}/cmd/shared/cmd.GitHash=$(GIT_HASH)'

APPS:=calendar scheduler sender

$(addprefix build-, $(APPS)): build-%: %
	go build -v -o ./bin/$< -ldflags "$(LDFLAGS)" ./cmd/$<

run-migrate:
	./bin/calendar migrate --config ./configs/sample.toml

$(addprefix run-, $(APPS)): run-%: % build-% run-migrate
	./bin/$< --config ./configs/sample.toml

build: build-calendar build-scheduler build-sender

build-api-img:
	docker build --build-arg=LDFLAGS="$(LDFLAGS)" -t $(DOCKER_API_IMG) -f build/api.Dockerfile .

version: build
	$(BIN) version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

generate-grpc:
	protoc --go_out=internal/server/grpc --go-grpc_out=internal/server/grpc api/EventService.proto

.PHONY: $(APPS) build run-api version test lint
