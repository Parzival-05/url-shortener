# Simple Makefile for a Go project
ENTRYPOINT = cmd/api/main.go
GO_FILES := $(shell find cmd internal -type f -name '*.go')
SWAGGER_FILE := docs/docs.go

# Build the application
all: generate-proto generate-swagger build test 

build:
	@echo "Building..."
	
	@go build -o main $(ENTRYPOINT)

# Run the application
run: generate-swagger
	@echo "Running..."

	@go run $(ENTRYPOINT) $(ARGS)

# Generate the swagger file
generate-swagger: $(SWAGGER_FILE)

$(SWAGGER_FILE): $(GO_FILES)
	@echo "Generating swagger..."
	@swag init -g $(ENTRYPOINT)

# Create DB container
docker-run:
	@if docker compose up --build 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker compose up --build; \
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
	@rm -r main

lint:
	@echo "Linting..."
	@golangci-lint run

generate-proto:
	@protoc --proto_path=api --go_out=api/gen/ --go_opt=paths=source_relative --go-grpc_out=api/gen/ --go-grpc_opt=paths=source_relative \
	proto/url_shortener/v1/url_shortener.proto

download-proto-deps: download-validate download-google-api download-grpc-gateway download-proto-go

download-proto-go:
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

download-validate:
	git clone -b main --single-branch --depth=2 --filter=tree:0 \
		https://github.com/bufbuild/protoc-gen-validate vendor.protogen/tmp && \
		cd vendor.protogen/tmp && \
		git sparse-checkout set --no-cone validate && \
		git checkout
	mkdir -p vendor.protogen/validate
	mv vendor.protogen/tmp/validate vendor.protogen/
	rm -rf vendor.protogen/tmp

download-google-api:
	git clone -b master --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/googleapis/googleapis vendor.protogen/googleapis && \
		cd vendor.protogen/googleapis && \
		git sparse-checkout set --no-cone google/api && \
		git checkout
	mkdir -p vendor.protogen/google
	mv vendor.protogen/googleapis/google/api vendor.protogen/google
	rm -rf vendor.protogen/googleapis

download-grpc-gateway:
	git clone -b main --single-branch -n --depth=1 --filter=tree:0 \
		https://github.com/grpc-ecosystem/grpc-gateway vendor.protogen/grpc-ecosystem && \
		cd vendor.protogen/grpc-ecosystem && \
		git sparse-checkout set --no-cone protoc-gen-openapiv2/options && \
		git checkout
	mkdir -p vendor.protogen/protoc-gen-openapiv2
	mv vendor.protogen/grpc-ecosystem/protoc-gen-openapiv2/options vendor.protogen/protoc-gen-openapiv2
	rm -rf vendor.protogen/grpc-ecosystem


.PHONY: all build test clean watch docker-run docker-down download-proto-deps download-validate download-google-api download-grpc-gateway download-proto-go
