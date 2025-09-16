# Docker Setup for Coinmate API Go Client

This document explains how to use Docker to run the Coinmate API Go Client and its tests.

## 🐳 Quick Start

### Prerequisites
- Docker
- Docker Compose

### Available Commands

```bash
# Show all available commands
make help

# Build Docker image
make docker-build

# Run tests in Docker
make docker-test

# Run application in Docker
make docker-run

# Start development environment
make docker-dev
```

## 🚀 Running the Application

### 1. Build the Docker Image
```bash
make docker-build
```

### 2. Run the Application
```bash
# Production mode
make docker-run

# Development mode with hot reload
make docker-dev
```

### 3. Set Environment Variables (Optional)
For secure endpoints, create a `.env` file:
```bash
COINMATE_CLIENT_ID=your_client_id
COINMATE_API_KEY=your_api_key
COINMATE_PRIVATE_KEY=your_private_key
```

## 🧪 Running Tests

### Run All Tests
```bash
make docker-test
```

### Run Tests with Coverage
```bash
docker-compose -f docker-compose.yml run --rm app-test
```

### Watch Tests (Development)
```bash
make docker-test-watch
```

## 🛠️ Development

### Development Environment
```bash
# Start development environment with hot reload
make docker-dev

# Or use the development compose file directly
docker-compose -f docker-compose.dev.yml up app-dev
```

### Quick Test Run
```bash
# Run tests without building full image
make quick-test
```

## 📁 Docker Files Structure

```
├── Dockerfile              # Production multi-stage build
├── Dockerfile.dev          # Development build with hot reload
├── docker-compose.yml      # Production services
├── docker-compose.dev.yml  # Development services
├── .dockerignore          # Files to exclude from build
└── Makefile               # Convenient commands
```

## 🔧 Docker Services

### Production Services (`docker-compose.yml`)
- `app-prod`: Production application
- `app-test`: Test runner
- `app-dev`: Development environment

### Development Services (`docker-compose.dev.yml`)
- `app-dev`: Development with hot reload
- `app-test`: Test runner
- `app-test-watch`: Test watcher

## 📊 Coverage Reports

After running tests, coverage reports are generated:
- `coverage.html`: HTML coverage report
- `coverage.out`: Raw coverage data

## 🐛 Troubleshooting

### Common Issues

1. **Port already in use**
   ```bash
   # Stop all containers
   make stop
   ```

2. **Build cache issues**
   ```bash
   # Clean and rebuild
   docker system prune -f
   make docker-build
   ```

3. **Permission issues**
   ```bash
   # Fix file permissions
   sudo chown -R $USER:$USER .
   ```

### Useful Commands

```bash
# View logs
make logs

# Stop all containers
make stop

# Clean build artifacts
make clean

# Rebuild without cache
docker build --no-cache -t coinmate-api-client .
```

## 🔒 Security

- The production image runs as a non-root user
- Sensitive data should be passed via environment variables
- The `.dockerignore` file excludes sensitive files from the build context

## 📈 Performance

### Development
- Uses volume mounts for fast file changes
- Includes hot reloading with Air
- Shared Go module cache

### Production
- Multi-stage build for smaller image size
- Alpine Linux base for security and size
- Optimized for container deployment

## 🎯 Next Steps

1. Set up your Coinmate API credentials
2. Run the application: `make docker-run`
3. Run tests: `make docker-test`
4. Start development: `make docker-dev`

For more information, see the main [README.md](README.md).
