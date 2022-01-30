BIN := "./bin/calendar"
DOCKER_IMG="calendar:develop"
REPO=github.com/pls87/OtusGolang_homework/hw12_13_14_15_calendar

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X '${REPO}/cmd/calendar/cmd.Release=develop' -X '${REPO}/cmd/calendar/cmd.BuildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S)' -X '${REPO}/cmd/calendar/cmd.GitHash=$(GIT_HASH)'

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/calendar

run: build
	$(BIN) server --config ./configs/templates/config.toml

#build-img:
#	docker build \
#		--build-arg=LDFLAGS="$(LDFLAGS)" \
#		-t $(DOCKER_IMG) \
#		-f build/Dockerfile .

#run-img: build-img
#	docker run $(DOCKER_IMG)

version: build
	$(BIN) version

test:
	go test -race ./internal/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run version test lint
