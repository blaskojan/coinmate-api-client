#!/bin/bash

# Test runner script for Coinmate API Go Client
# This script runs all tests and generates coverage reports

echo "🧪 Running Coinmate API Go Client Tests"
echo "========================================"

# Create coverage directory
mkdir -p coverage

# Run tests for main client
echo "📦 Testing main client..."
go test -v -coverprofile=coverage/client.out ./coinmate/

# Run tests for public endpoints
echo "🌐 Testing public endpoints..."
go test -v -coverprofile=coverage/public.out ./coinmate/public/

# Run tests for secure endpoints
echo "🔒 Testing secure endpoints..."
go test -v -coverprofile=coverage/secure.out ./coinmate/secure/

# Generate combined coverage report
echo "📊 Generating coverage report..."
go tool cover -html=coverage/client.out -o coverage/client.html
go tool cover -html=coverage/public.out -o coverage/public.html
go tool cover -html=coverage/secure.out -o coverage/secure.html

# Show coverage summary
echo ""
echo "📈 Coverage Summary:"
echo "==================="
echo "Main Client:"
go tool cover -func=coverage/client.out | tail -1
echo ""
echo "Public Endpoints:"
go tool cover -func=coverage/public.out | tail -1
echo ""
echo "Secure Endpoints:"
go tool cover -func=coverage/secure.out | tail -1

# Run all tests together
echo ""
echo "🚀 Running all tests..."
go test -v ./...

# Check if all tests passed
if [ $? -eq 0 ]; then
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

