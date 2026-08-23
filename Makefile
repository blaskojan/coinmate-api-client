# Makefile for Coinmate API Go Client
.PHONY: help build test run dev clean docker-build docker-test docker-run docker-dev

# Default target
help:
	@echo "🚀 Coinmate API Go Client - Available Commands:"
	@echo ""
	@echo "📦 Docker Commands:"
	@echo "  make docker-build    - Build Docker image"
	@echo "  make docker-test     - Run tests in Docker"
	@echo "  make docker-run      - Run application in Docker"
	@echo "  make docker-dev      - Run development environment"
	@echo ""
	@echo "🧪 Test Commands:"
	@echo "  make test            - Run tests locally"
	@echo "  make test-coverage   - Run tests with coverage"
	@echo ""
	@echo "🔧 Development Commands:"
	@echo "  make build           - Build application locally"
	@echo "  make run             - Run application locally"
	@echo "  make clean           - Clean build artifacts"
	@echo ""

# Docker commands
docker-build:
	@echo "🔨 Building Docker image..."
	docker build -t coinmate-api-client .

docker-test:
	@echo "🧪 Running tests in Docker..."
	docker run --rm -v $(PWD):/app -w /app golang:1.27-alpine sh -c "go mod download && go test -v -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o coverage.html"

docker-run:
	@echo "🚀 Running application in Docker..."
	docker run --rm -v $(PWD):/app -w /app golang:1.27-alpine sh -c "go mod download && go run ."

docker-dev:
	@echo "🛠️  Starting development environment..."
	docker run --rm -it -v $(PWD):/app -w /app -p 8080:8080 golang:1.27-alpine sh -c "go mod download && go run ."

docker-test-watch:
	@echo "👀 Starting test watcher..."
	docker run --rm -v $(PWD):/app -w /app golang:1.27-alpine sh -c "go mod download && while true; do go test -v ./...; sleep 2; done"

# Local commands (requires Go installed)
build:
	@echo "🔨 Building application..."
	go build -o main .

test:
	@echo "🧪 Running tests..."
	go test -v ./...

test-coverage:
	@echo "🧪 Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report generated: coverage.html"

run:
	@echo "🚀 Running application..."
	go run .

clean:
	@echo "🧹 Cleaning build artifacts..."
	rm -f main
	rm -f test-runner
	rm -rf coverage/
	rm -f *.out
	rm -f *.html

# Development helpers
install-deps:
	@echo "📦 Installing dependencies..."
	go mod download

fmt:
	@echo "🎨 Formatting code..."
	go fmt ./...

lint:
	@echo "🔍 Linting code..."
	golangci-lint run

# Docker development with hot reload
dev-hot-reload:
	@echo "🔥 Starting development with hot reload..."
	docker-compose -f docker-compose.dev.yml up app-dev

# Quick test run
quick-test:
	@echo "⚡ Running quick tests..."
	docker run --rm -v $(PWD):/app -w /app golang:1.27-alpine sh -c "go mod download && go test -v ./..."

# Alternative test run (if docker compose is not available)
test-docker:
	@echo "🧪 Running tests with Docker..."
	docker run --rm -v $(PWD):/app -w /app golang:1.27-alpine sh -c "go mod download && go test -v -coverprofile=coverage.out ./... && go tool cover -html=coverage.out -o coverage.html"

# Production build
prod-build:
	@echo "🏭 Building production image..."
	docker build -t coinmate-api-client:prod .

# Show logs
logs:
	@echo "📋 Showing logs..."
	@echo "No logs available with direct Docker commands"

# Stop all containers
stop:
	@echo "🛑 Stopping all containers..."
	docker stop $(docker ps -q) 2>/dev/null || true
