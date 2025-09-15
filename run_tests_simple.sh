#!/bin/bash

echo "🧪 Running Coinmate API Go Client Tests (Simple Mode)"
echo "===================================================="

# Check if Go is available
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed or not in PATH"
    echo ""
    echo "📦 Installation Instructions:"
    echo "============================="
    echo ""
    echo "For macOS:"
    echo "1. Install Go: brew install go"
    echo "2. Or download from: https://golang.org/dl/"
    echo ""
    echo "For Ubuntu/Debian:"
    echo "1. sudo apt-get update"
    echo "2. sudo apt-get install golang-go"
    echo ""
    echo "For CentOS/RHEL:"
    echo "1. sudo yum install -y golang"
    echo ""
    echo "After installation, restart your terminal and run this script again."
    exit 1
fi

echo "✅ Go is available: $(go version)"
echo ""

# Create coverage directory
mkdir -p coverage

# Run tests for main client
echo "📦 Testing main client..."
if go test -v -coverprofile=coverage/client.out ./coinmate/; then
    echo "✅ Main client tests passed"
else
    echo "❌ Main client tests failed"
    exit 1
fi

# Run tests for public endpoints
echo ""
echo "🌐 Testing public endpoints..."
if go test -v -coverprofile=coverage/public.out ./coinmate/public/; then
    echo "✅ Public endpoints tests passed"
else
    echo "❌ Public endpoints tests failed"
    exit 1
fi

# Run tests for secure endpoints
echo ""
echo "🔒 Testing secure endpoints..."
if go test -v -coverprofile=coverage/secure.out ./coinmate/secure/; then
    echo "✅ Secure endpoints tests passed"
else
    echo "❌ Secure endpoints tests failed"
    exit 1
fi

# Generate coverage reports
echo ""
echo "📊 Generating coverage reports..."
if [ -f coverage/client.out ]; then
    go tool cover -html=coverage/client.out -o coverage/client.html
    echo "✅ Client coverage report: coverage/client.html"
fi

if [ -f coverage/public.out ]; then
    go tool cover -html=coverage/public.out -o coverage/public.html
    echo "✅ Public coverage report: coverage/public.html"
fi

if [ -f coverage/secure.out ]; then
    go tool cover -html=coverage/secure.out -o coverage/secure.html
    echo "✅ Secure coverage report: coverage/secure.html"
fi

# Show coverage summary
echo ""
echo "📈 Coverage Summary:"
echo "==================="
if [ -f coverage/client.out ]; then
    echo "Main Client:"
    go tool cover -func=coverage/client.out | tail -1
    echo ""
fi

if [ -f coverage/public.out ]; then
    echo "Public Endpoints:"
    go tool cover -func=coverage/public.out | tail -1
    echo ""
fi

if [ -f coverage/secure.out ]; then
    echo "Secure Endpoints:"
    go tool cover -func=coverage/secure.out | tail -1
    echo ""
fi

# Run all tests together
echo ""
echo "🚀 Running all tests together..."
if go test -v ./...; then
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


