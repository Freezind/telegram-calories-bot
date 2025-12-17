.PHONY: help build run stop logs clean docker-build docker-run docker-stop test-llm

# Default target
help:
	@echo "Available commands:"
	@echo "  make build         - Build the unified backend binary"
	@echo "  make run           - Run the unified backend locally"
	@echo "  make docker-build  - Build Docker image"
	@echo "  make docker-run    - Run with Docker Compose"
	@echo "  make docker-stop   - Stop Docker containers"
	@echo "  make logs          - View Docker logs"
	@echo "  make test-llm      - Run LLM-based tests"
	@echo "  make clean         - Clean build artifacts"

# Build binary
build:
	@echo "Building unified backend..."
	go build -o unified cmd/unified/main.go
	@echo "✓ Build complete"

# Run locally
run: build
	@echo "Starting unified backend..."
	./unified

# Docker commands
docker-build:
	@echo "Building Docker image..."
	docker build -t telegram-calories-bot .
	@echo "✓ Docker image built"

docker-run:
	@echo "Starting with Docker Compose..."
	docker-compose up -d
	@echo "✓ Containers started"
	@echo "View logs: make logs"

docker-stop:
	@echo "Stopping containers..."
	docker-compose down
	@echo "✓ Containers stopped"

logs:
	docker-compose logs -f

# Test
test-llm:
	@echo "Running LLM-based tests..."
	./test-llm.sh

# Clean
clean:
	@echo "Cleaning build artifacts..."
	rm -f unified
	docker-compose down 2>/dev/null || true
	@echo "✓ Cleanup complete"
