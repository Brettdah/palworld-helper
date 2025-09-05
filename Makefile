# Makefile for Palworld Helper

.PHONY: help build run stop clean logs shell

# Default target
help:
	@echo "Available commands:"
	@echo "  build    - Build the Docker container"
	@echo "  run      - Build and run the application"
	@echo "  stop     - Stop the running container"
	@echo "  clean    - Remove container and image"
	@echo "  logs     - Show application logs"
	@echo "  shell    - Open shell in running container"
	@echo "  rebuild  - Clean, build and run"

# Build the Docker image
build:
	@echo "Building Palworld Helper..."
	docker compose build

# Build and run the application
run:
	@echo "Starting Palworld Helper..."
	docker compose up -d
	@echo "Application is running at http://localhost:8080"

# Stop the running container
stop:
	@echo "Stopping Palworld Helper..."
	docker compose down

# Clean up containers and images
clean:
	@echo "Cleaning up..."
	docker compose down -v --rmi all --remove-orphans
	docker system prune -f

# Show logs
logs:
	docker compose logs -f palworld-helper

# Open shell in the running container
shell:
	docker exec -it palworld-helper-app /bin/sh

# Rebuild everything
rebuild: clean build run

# Quick development cycle
dev: stop build run logs