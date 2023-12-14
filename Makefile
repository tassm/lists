# Define the shell to use for running make commands
SHELL := /bin/bash

# .PHONY serve-web build run test clean package

# Define the default goal
.DEFAULT_GOAL := build

all: clean test build

# serve http content
serve-web:
	http-server web/

# Build the project
build:
	go build -o ./bin/lists-server/main cmd/lists-server/main.go
	go build -o ./bin/lists-lambda/main cmd/lists-lambda/main.go

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
	rm -r bin dist

package-lambda:
	mkdir -p dist
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o ./bin/lists-lambda/bootstrap cmd/lists-lambda/main.go
	zip -j bin/list-service-lambda.zip bin/lists-lambda/bootstrap
	mv bin/list-service-lambda.zip dist