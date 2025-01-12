APP_NAME := tempest
VERSION := $(shell git describe --tags --always)
BUILD_TIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
COMMIT_HASH := $(shell git rev-parse HEAD)

GO_FILES := $(shell find . -name '*.go' -type f)
LDFLAGS := -X 'main.Version=$(VERSION)' -X 'main.BuildTimestamp=$(BUILD_TIME)' -X 'main.CommitHash=$(COMMIT_HASH)'

.PHONY: all build clean

all: build

build: $(GO_FILES)
	@echo "Building $(APP_NAME)..."
	go build -ldflags "$(LDFLAGS)" -o dist/$(APP_NAME) pkg/cmd/main.go

test:
	@echo "Running tests..."
	go test -v ./...

clean:
	@echo "Cleaning up..."
	rm -f $(APP_NAME)
