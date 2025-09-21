# ==========================================================================
# Build Variables
# ==========================================================================
VERSION := $(shell git describe --tags --always)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
COMMIT_SHA := $(shell git rev-parse --short HEAD)
BUILD_TIME := $(shell date +%Y-%m-%d_%H:%M:%S)

# Go build flags
LDFLAGS := -X main.Version=$(VERSION) \
           -X main.Branch=$(BRANCH) \
           -X main.CommitSHA=$(COMMIT_SHA) \
           -X main.BuildTime=$(BUILD_TIME)

# System information
GOHOSTOS := $(shell go env GOHOSTOS)
GOPATH := $(shell go env GOPATH)
GOARCH := $(shell go env GOARCH)
GO := go

# Project information
PROJECT_NAME := kratos-kit
MAIN_PACKAGES := ./cmd/api/... ./cmd/admin/...

# Proto files
ifeq ($(GOHOSTOS), windows)
	Git_Bash := $(subst \,/,$(subst cmd\,bin\bash.exe,$(dir $(shell where git))))
	INTERNAL_PROTO_FILES := $(shell $(Git_Bash) -c "find internal -name *.proto")
	API_PROTO_FILES := $(shell $(Git_Bash) -c "find api -name *.proto")
else
	INTERNAL_PROTO_FILES := $(shell find internal -name *.proto)
	API_PROTO_FILES := $(shell find api -name *.proto)
endif

# Docker variables
DOCKER_REGISTRY ?= your-registry
DOCKER_IMAGE := $(DOCKER_REGISTRY)/$(PROJECT_NAME)
DOCKER_TAG ?= $(VERSION)

# ==========================================================================
# Development Environment Setup
# ==========================================================================
.PHONY: init
# Install development tools
init:
	@echo "Installing development tools..."
	$(GO) install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
	$(GO) install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
	$(GO) install github.com/google/wire/cmd/wire@latest
	$(GO) install github.com/bufbuild/buf/cmd/buf@latest
	$(GO) install github.com/air-verse/air@latest
	$(GO) install github.com/golangci/golint/cmd/golint@latest
	$(GO) install github.com/swaggo/swag/cmd/swag@latest
	$(GO) install github.com/vektra/mockery/v2@latest
	$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: mod
# Download and tidy dependencies
mod:
	@echo "Downloading and tidying dependencies..."
	$(GO) mod download
	$(GO) mod tidy
	$(GO) mod verify

# ==========================================================================
# Code Generation
# ==========================================================================
.PHONY: gen
# Generate proto files
gen:
	@echo "Generating proto files..."
	buf format -w
	buf generate
	bin/codegen.exe

.PHONY: wire
# Generate wire dependency injection for API service
wire:
	@echo "Generating wire dependency injection..."
	cd cmd/api && wire

.PHONY: wire-admin
# Generate wire for admin service
wire-admin:
	@echo "Generating wire for admin service..."
	cd cmd/admin && wire

.PHONY: gorm-gen
# Generate GORM models
gorm-gen:
	@echo "Generating GORM models..."
	go run tools/gorm-gen/main.go

.PHONY: swagger
# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	swag init -g cmd/api/main.go -o api/swagger

.PHONY: mock
# Generate mocks for interfaces
mock:
	@echo "Generating mocks..."
	mockery --all --keeptree

# ==========================================================================
# Build and Run
# ==========================================================================
.PHONY: build
# Build all applications
build:
	@echo "Building applications..."
	mkdir -p bin/ && $(GO) build -ldflags "$(LDFLAGS)" -o ./bin/ $(MAIN_PACKAGES)

.PHONY: build-linux
# Build for Linux
build-linux:
	@echo "Building for Linux..."
	GOOS=linux GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" -o ./bin/ $(MAIN_PACKAGES)

.PHONY: build-windows
# Build for Windows
build-windows:
	@echo "Building for Windows..."
	GOOS=windows GOARCH=amd64 $(GO) build -ldflags "$(LDFLAGS)" -o ./bin/ $(MAIN_PACKAGES)

.PHONY: run
# Run application with air for hot reload
run:
	@echo "Running application with hot reload..."
	air -- -f configs/config.local.yaml

.PHONY: run-prod
# Run application in production mode
run-prod:
	@echo "Running application in production mode..."
	./bin/api -f configs/config.prod.yaml

# ==========================================================================
# Testing and Quality
# ==========================================================================
.PHONY: test
# Run tests
test:
	@echo "Running tests..."
	go test -v -race -cover ./...

.PHONY: test-coverage
# Run tests with coverage report
test-coverage:
	@echo "Running tests with coverage report..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: lint
# Run golangci-lint
lint:
	@echo "Running linter..."
	golangci-lint run ./...

.PHONY: pprof
# Start pprof profiling
pprof:
	@echo "Starting pprof profiling..."
	go tool pprof -http=:8080 http://localhost:9003/debug/pprof/profile

# ==========================================================================
# Docker Operations
# ==========================================================================
.PHONY: up
# Start docker containers
up:
	@echo "Starting docker containers..."
	docker compose up -d

.PHONY: down
# Stop docker containers
down:
	@echo "Stopping docker containers..."
	docker compose down

.PHONY: docker-build
# Build docker image
docker-build:
	@echo "Building docker image..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	docker tag $(DOCKER_IMAGE):$(DOCKER_TAG) $(DOCKER_IMAGE):latest

.PHONY: docker-push
# Push docker image
docker-push:
	@echo "Pushing docker image..."
	docker push $(DOCKER_IMAGE):$(DOCKER_TAG)
	docker push $(DOCKER_IMAGE):latest

.PHONY: compose-up
# Start docker compose with specific environment
compose-up:
	@echo "Starting docker compose..."
	docker compose -f docker-compose.yml -f docker-compose.$(ENV).yml up -d

.PHONY: compose-down
# Stop docker compose
compose-down:
	@echo "Stopping docker compose..."
	docker compose -f docker-compose.yml -f docker-compose.$(ENV).yml down

# ==========================================================================
# Deployment
# ==========================================================================
.PHONY: deploy
# Deploy application
deploy:
	@echo "Deploying application..."
	chmod +x scripts/deploy.sh
	./scripts/deploy.sh

# ==========================================================================
# Maintenance
# ==========================================================================
.PHONY: clean
# Clean build artifacts and temporary files
clean:
	@echo "Cleaning build artifacts and temporary files..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	rm -f api/swagger/*
	find . -type f -name '*.test' -delete
	find . -type f -name '*.out' -delete
	find . -type d -name 'vendor' -exec rm -rf {} +
	$(GO) clean -cache -testcache -modcache

.PHONY: update-deps
# Update all dependencies
update-deps:
	@echo "Updating dependencies..."
	$(GO) get -u ./...
	$(GO) mod tidy
	$(GO) mod verify

.PHONY: all
# Generate and build everything
all:
	@echo "Running complete build process..."
	make mod
	make gen
	make wire
	make generate
	make build

# ==========================================================================
# Help
# ==========================================================================
.PHONY: help
# Show help information
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[36m%-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help
