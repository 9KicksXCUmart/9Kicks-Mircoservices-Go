.DEFAULT_GOAL := build

.PHONY:fmt vet build
fmt:
	go fmt ./...

vet: fmt
	go vet ./...

build: vet
	air

clean:
	rm -f hello_world

setup:
	go mod download
	go install github.com/cosmtrek/air@latest
