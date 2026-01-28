.PHONY: all build clean run test vendor wire docker docker-run docker-up docker-down help

# Variables
APP_NAME := key-spy
BINARY_NAME := scanner
ROOT_DIR := ../..
DOCKER_IMAGE := $(APP_NAME):latest

# Build docker image
docker:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) -f Dockerfile $(ROOT_DIR)
	@echo "Docker image built: $(DOCKER_IMAGE)"

# Run docker container
docker-run:
	@echo "Running Docker container..."
	docker run --rm \
		-v $(ROOT_DIR)/output:/data/key-spy/output \
		-v $(ROOT_DIR)/logs:/data/key-spy/logs \
		$(DOCKER_IMAGE)

# Run with docker-compose
docker-up:
	@echo "Starting with docker-compose..."
	docker compose up -d

# Stop docker-compose
docker-down:
	@echo "Stopping docker-compose..."
	docker compose down

# Run in development mode (single scan)
dev:
	@echo "Running in development mode..."
	cd $(ROOT_DIR) && go build -o bin/$(BINARY_NAME) ./cmd/scanner
	cd $(ROOT_DIR) && bin/$(BINARY_NAME) -cfg_path config/scanner/localcfg -cfg_name test.yaml

# Show help
help:
	@echo "Available targets:"
	@echo "  docker     - Build Docker image"
	@echo "  docker-run - Run Docker container"
	@echo "  docker-up  - Start with docker-compose"
	@echo "  docker-down- Stop docker-compose"
	@echo "  dev        - Run in development mode"
	@echo "  help       - Show this help"
