# Simple Makefile for a Go project
ENTRYPOINT = cmd/api/main.go
GO_FILES := $(shell find cmd internal -type f -name '*.go')
SWAGGER_FILE := docs/docs.go

# Build the application
all: generate build test

build:
	@echo "Building..."
	
	@go build -o main $(ENTRYPOINT)

# Run the application
run: generate
	@echo "Running..."

	@go run $(ENTRYPOINT) $(ARGS)

# Generate the swagger file
generate: $(SWAGGER_FILE)

$(SWAGGER_FILE): $(GO_FILES)
	@echo "Generating swagger..."
	@swag init -g $(ENTRYPOINT)

# Create DB container
docker-run:
	@if docker compose up --build 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up --build; \
	fi

# Shutdown DB container
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -r bin

lint:
	@echo "Linting..."
	@golangci-lint run

.PHONY: all build test clean watch docker-run docker-down
