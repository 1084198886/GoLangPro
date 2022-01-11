GO ?=go
GOMOBILE ?=gomobile
OUTPUT ?=./build
DIST ?= ./dist
RM ?=rm
VERSION ?=$(shell git describe --tags --abbrev=6 --always --dirty)
BUILD_DATE ?=$(shell date '+%Y-%m-%d %H:%M:%S %Z')

.PHONY: clean build generate client dist

build:
	$(GO) mod vendor
	$(GO) build -v -o $(OUTPUT)/ cmd/weight_agent.go

generate:
	$(GO) generate -v -x ./agent

bootstrap: generate
	$(GO) mod vendor

server: generate
	GOOS=linux GOARCH=amd64 $(GO) build -v \
		-ldflags 'main.Version="$(VERSION)" main.BuildDate="$(BUILD_DATE)"' \
		-o $(DIST)/go_demo_$(VERSION) ./agent

client:
	$(GO) generate ./client
	$(GO) build -v -o build/ ./client/*.go

test: bootstrap
	$(GO) test ./agent

mobile: test
	$(RM) -rf vendor
	mkdir -p dist
	$(GOMOBILE) bind -v -target=android/arm64 -classpath com.supwisdom.diancan \
		-o $(DIST)/godemo_$(VERSION).aar ./agent

clean:
	$(RM) -rf $(OUTPUT)/*
