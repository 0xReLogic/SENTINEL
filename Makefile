.PHONY: build clean test lint run help

# Binary name
BINARY_NAME=sentinel

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOLINT=golangci-lint

# Build flags
LDFLAGS=-ldflags "-s -w"

all: test build

build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) -v

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME).exe

test:
	$(GOTEST) -v ./...

lint:
	$(GOLINT) run

run:
	$(GOBUILD) -o $(BINARY_NAME) -v
	./$(BINARY_NAME)

validate:
	./$(BINARY_NAME) validate

once:
	./$(BINARY_NAME) once

deps:
	$(GOMOD) tidy
	$(GOGET) -u

help:
	@echo "make - Compile the binary and run tests"
	@echo "make build - Compile the binary"
	@echo "make clean - Remove binary and object files"
	@echo "make test - Run tests"
	@echo "make lint - Run linter"
	@echo "make run - Build and run the application"
	@echo "make validate - Build and validate the configuration"
	@echo "make once - Build and run a single check"
	@echo "make deps - Update dependencies"