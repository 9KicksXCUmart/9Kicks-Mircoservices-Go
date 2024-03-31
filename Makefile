.DEFAULT_GOAL := build

.PHONY:fmt vet build
fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build: vet
	air

setup:
	git config core.hooksPath .githooks
	go mod download
	go install github.com/cosmtrek/air@latest
