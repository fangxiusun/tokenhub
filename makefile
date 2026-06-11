.PHONY: build build-frontend build-backend run clean test lint

# Build frontend
build-frontend:
	cd web/default && npm install && npm run build

# Build backend (frontend must be built first)
build-backend:
	CGO_ENABLED=0 go build -o bin/server .

# Build everything (frontend + backend into single binary)
build: build-frontend build-backend

# Run in development mode
run:
	go run .

# Run in production mode
run-prod:
	GIN_MODE=release go run .

# Clean build artifacts
clean:
	rm -rf bin/
	rm -rf web/default/dist/

# Run tests
test:
	go test ./...

# Run linter
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./...

# Tidy dependencies
tidy:
	go mod tidy

# Docker build
docker-build:
	docker build -t tokenhub .

# Docker run
docker-run:
	docker-compose up -d

# Docker stop
docker-stop:
	docker-compose down

# Help
help:
	@echo "Available commands:"
	@echo "  build          - Build frontend and backend into single binary"
	@echo "  build-frontend - Build frontend only"
	@echo "  build-backend  - Build backend only"
	@echo "  run            - Run in development mode"
	@echo "  run-prod       - Run in production mode"
	@echo "  clean          - Clean build artifacts"
	@echo "  test           - Run tests"
	@echo "  lint           - Run linter"
	@echo "  fmt            - Format code"
	@echo "  tidy           - Tidy dependencies"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Start Docker containers"
	@echo "  docker-stop    - Stop Docker containers"
	@echo "  help           - Show this help"
