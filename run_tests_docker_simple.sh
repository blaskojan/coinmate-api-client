#!/bin/bash

echo "🧪 Running Coinmate API Go Client Tests with Docker (Simple Mode)"
echo "================================================================"

# Check if Docker is available
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed or not in PATH"
    echo ""
    echo "📦 Installation Instructions:"
    echo "============================="
    echo ""
    echo "For macOS:"
    echo "1. Install Docker Desktop: https://www.docker.com/products/docker-desktop"
    echo "2. Or use Homebrew: brew install --cask docker"
    echo ""
    echo "For Ubuntu/Debian:"
    echo "1. sudo apt-get update"
    echo "2. sudo apt-get install docker.io"
    echo "3. sudo usermod -aG docker \$USER"
    echo ""
    echo "After installation, restart your terminal and run this script again."
    exit 1
fi

echo "✅ Docker is available: $(docker --version)"
echo ""

# Check if Docker daemon is running
if ! docker info &> /dev/null; then
    echo "❌ Docker daemon is not running"
    echo ""
    echo "🚀 Start Docker:"
    echo "==============="
    echo ""
    echo "For macOS:"
    echo "1. Open Docker Desktop application"
    echo "2. Wait for Docker to start"
    echo ""
    echo "For Linux:"
    echo "1. sudo systemctl start docker"
    echo "2. sudo systemctl enable docker"
    echo ""
    echo "After starting Docker, run this script again."
    exit 1
fi

echo "✅ Docker daemon is running"
echo ""

# Create coverage directory
mkdir -p coverage

echo "🔨 Building test environment..."
echo ""

# Run tests for main client
echo "📦 Testing main client..."
if docker run --rm -v $(PWD):/app -w /app golang:1.25-alpine sh -c "go mod download && go test -v -coverprofile=coverage/client.out ./coinmate/"; then
    echo "✅ Main client tests passed"
else
    echo "❌ Main client tests failed"
    exit 1
fi

# Run tests for public endpoints
echo ""
echo "🌐 Testing public endpoints..."
if docker run --rm -v $(PWD):/app -w /app golang:1.25-alpine sh -c "go mod download && go test -v -coverprofile=coverage/public.out ./coinmate/public/"; then
    echo "✅ Public endpoints tests passed"
else
    echo "❌ Public endpoints tests failed"
    exit 1
fi

# Run tests for secure endpoints
echo ""
echo "🔒 Testing secure endpoints..."
if docker run --rm -v $(PWD):/app -w /app golang:1.25-alpine sh -c "go mod download && go test -v -coverprofile=coverage/secure.out ./coinmate/secure/"; then
    echo "✅ Secure endpoints tests passed"
else
    echo "❌ Secure endpoints tests failed"
    exit 1
fi

# Generate coverage reports
echo ""
echo "📊 Generating coverage reports..."
if [ -f coverage/client.out ]; then
    docker run --rm -v $(PWD):/app -w /app golang:1.25-alpine sh -c "go tool cover -html=coverage/client.out -o coverage/client.html"
    echo "✅ Client coverage report: coverage/client.html"
fi

if [ -f coverage/public.out ]; then
    docker run --rm -v $(PWD):/app -w /app golang:1.25-alpine sh -c "go tool cover -html=coverage/public.out -o coverage/public.html"
    echo "✅ Public coverage report: coverage/public.html"
fi

if [ -f coverage/secure.out ]; then
    docker run --rm -v $(PWD):/app -w /app golang:1.25-alpine sh -c "go tool cover -html=coverage/secure.out -o coverage/secure.html"
    echo "✅ Secure coverage report: coverage/secure.html"
fi

# Show coverage summary
echo ""
echo "📈 Coverage Summary:"
echo "==================="
if [ -f coverage/client.out ]; then
    echo "Main Client:"
    docker run --rm -v $(PWD):/app -w /app golang:1.25-alpine sh -c "go tool cover -func=coverage/client.out | tail -1"
    echo ""
fi

if [ -f coverage/public.out ]; then
    echo "Public Endpoints:"
    docker run --rm -v $(PWD):/app -w /app golang:1.25-alpine sh -c "go tool cover -func=coverage/public.out | tail -1"
    echo ""
fi

if [ -f coverage/secure.out ]; then
    echo "Secure Endpoints:"
    docker run --rm -v $(PWD):/app -w /app golang:1.25-alpine sh -c "go tool cover -func=coverage/secure.out | tail -1"
    echo ""
fi

# Run all tests together
echo ""
echo "🚀 Running all tests together..."
if docker run --rm -v $(PWD):/app -w /app golang:1.25-alpine sh -c "go mod download && go test -v ./..."; then
    echo ""
    echo "✅ All tests passed!"
    echo ""
    echo "📁 Coverage reports generated in:"
    echo "   - coverage/client.html"
    echo "   - coverage/public.html"
    echo "   - coverage/secure.html"
else
    echo ""
    echo "❌ Some tests failed!"
    exit 1
fi

