GIT_VERSION = $(shell git rev-parse --short HEAD)

.PHONY: build
build:
	go build -o mcp-k8s-eye main.go