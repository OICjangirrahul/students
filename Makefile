.PHONY: all build test clean docker-up docker-down docs

# Default target
all: build

# Build the application
build:
	go build -o bin/app cmd/main.go

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Start Docker containers
docker-up:
	docker-compose up -d

# Stop Docker containers
docker-down:
	docker-compose down

# Run the application
run: docker-up
	go run cmd/main.go

# Run tests with Docker
test-with-docker: docker-up
	go test -v ./...

# Generate Swagger documentation
docs:
	swag init -g cmd/main.go -o docs 