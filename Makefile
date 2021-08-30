.PHONY: install test binary race cover 

CURRENT_DIR := $(shell pwd)
VERSION := $(shell git describe --always)
GOLIST := $(shell go list -m)
BINARY := fetcher

binary:
	@go build -ldflags "-X ${GOLIST}/internal/version.version=$(VERSION)" -o ${BINARY} ${CURRENT_DIR}/cmd/srv

test:
	@go test ./... -cover 

test-unit:
	@go test ./... -run Unit -v -race 

test-integration:
	@go test ./... -run Integration -v -race

race:
	@go test -race ./...