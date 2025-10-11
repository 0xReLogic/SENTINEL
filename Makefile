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
	@echo ""
	@echo "Docker targets:"
	@echo "make docker-build - Build Docker image"
	@echo "make docker-run - Run Docker container"
	@echo "make docker-compose-up - Start services with Docker Compose"
	@echo "make docker-compose-down - Stop services with Docker Compose"
	@echo "make docker-compose-logs - View Docker Compose logs"
	@echo "make docker-clean - Clean Docker artifacts"
	@echo "make quick-start - Quick start with Docker Compose"

# Docker targets
DOCKER_IMAGE=sentinel:latest

.PHONY: docker-build
docker-build:
	docker build -t $(DOCKER_IMAGE) .

.PHONY: docker-run
docker-run:
	docker run --rm --env-file .env -v $(PWD)/sentinel.yaml:/app/sentinel.yaml $(DOCKER_IMAGE)

.PHONY: docker-compose-up
docker-compose-up:
	docker compose up -d

.PHONY: docker-compose-down
docker-compose-down: ## Stop services with Docker Compose V2
	docker compose down

.PHONY: docker-compose-logs
docker-compose-logs: ## View Docker Compose logs
	docker compose logs -f

.PHONY: docker-clean
docker-clean: ## Clean Docker artifacts
	docker compose down -v
	docker image rm $(DOCKER_IMAGE) 2>/dev/null || true

.PHONY: quick-start
quick-start: ## Quick start with Docker Compose
	@echo "Setting up SENTINEL with Docker..."
	@test -f .env || cp .env.example .env
	@echo "Please edit .env file with your notification tokens (optional)"
	@echo "Starting SENTINEL..."
	@$(MAKE) docker-compose-up
	@echo "SENTINEL is running! Check status with: make docker-compose-logs"

.PHONY: docker-validate
docker-validate: ## Validate Docker setup and configuration
	@echo "🔍 Validating Docker setup..."
	@test -f Dockerfile || (echo "Dockerfile not found" && exit 1)
	@test -f docker-compose.yml || (echo "docker-compose.yml not found" && exit 1)
	@test -f .dockerignore || (echo ".dockerignore not found" && exit 1)
	@test -f .env.example || (echo ".env.example not found" && exit 1)
	@echo "All Docker files present"

.PHONY: docker-test-build
docker-test-build: ## Test Docker image build and validate size
	@echo "🔨 Testing Docker image build..."
	@docker build -t sentinel:test .
	@echo "📏 Checking image size..."
	@SIZE_BYTES=$$(docker inspect -f '{{.Size}}' sentinel:test) && \
	SIZE_MB=$$((SIZE_BYTES / 1024 / 1024)) && \
	echo "Image size: ~$$SIZE_MB MB ($$SIZE_BYTES bytes)" && \
	if [ $$SIZE_BYTES -gt 31457280 ]; then \
		echo "Image size exceeds 30MB limit!"; exit 1; \
	else \
		echo "Image size is within limits"; \
	fi

.PHONY: docker-test-run
docker-test-run:
	@echo "Testing container functionality..."
	@echo "Testing sentinel once command..."
	@docker run --rm sentinel:test ./sentinel once > /dev/null && echo "sentinel once works" || (echo "sentinel once failed" && exit 1)
	@echo "Testing sentinel validate command..."
	@docker run --rm -v $(PWD)/sentinel.yaml:/app/sentinel.yaml sentinel:test ./sentinel validate > /dev/null && echo "sentinel validate works" || (echo "sentinel validate failed" && exit 1)

.PHONY: docker-security-scan
docker-security-scan:
	@echo "Running security scan..."
	@which trivy > /dev/null || (echo "Install trivy for security scanning: https://trivy.dev" && exit 0)
	@echo "Scanning for HIGH and CRITICAL vulnerabilities..."
	@trivy image sentinel:test --exit-code 1 --severity HIGH,CRITICAL --format json --output security-report.json

.PHONY: docker-full-test
docker-full-test: ## Run complete Docker test suite
	@echo "Running complete Docker test suite..."
	@$(MAKE) docker-validate
	@$(MAKE) docker-test-build
	@$(MAKE) docker-test-run
	@echo "All Docker tests passed!"
