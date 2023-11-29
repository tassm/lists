# Define the shell to use for running make commands
SHELL := /bin/bash

# Define the default goal
.DEFAULT_GOAL := build

all: clean test build

# serve http content
serve-web:
	http-server web/

# Build the project
build:
	go build -o ./bin/main cmd/listr/main.go

# Run the project
run: build
	./bin/main

# Test the project
test:
	go test ./...

# Clean the project
clean:
	go fmt ./...
	go mod tidy
	rm -rf ./bin