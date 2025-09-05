# Makefile for Palworld Helper - Refactored Version

# Docker Compose files
PROD_COMPOSE = docker-compose.yml
DEV_COMPOSE = docker-compose.dev.yml

# Default environment (can be overridden)
ENV ?= prod
COMPOSE_FILE = $(if $(filter dev,$(ENV)),$(DEV_COMPOSE),$(PROD_COMPOSE))
SERVICE_NAME = palworld-helper
PORT = 8080

.PHONY: help build run stop clean logs shell rebuild watch restart dev

# Default target
help:
	@echo "Available commands:"
	@echo ""
	@echo "Generic commands (use ENV=dev for development):"
	@echo "  build     - Build the Docker container"
	@echo "  run       - Build and run the application"
	@echo "  stop      - Stop the running container"
	@echo "  logs      - Show application logs"
	@echo "  shell     - Open shell in running container"
	@echo "  rebuild   - Clean, build and run"
	@echo "  restart   - Quick restart without rebuild"
	@echo "  clean     - Remove containers and images"
	@echo ""
	@echo "Development shortcuts:"
	@echo "  dev       - Quick dev cycle (ENV=dev stop build run logs)"
	@echo "  watch     - Auto-rebuild on file changes (dev only)"
	@echo ""
	@echo "Usage examples:"
	@echo "  make build           # Build production"
	@echo "  make build ENV=dev   # Build development"
	@echo "  make run ENV=dev     # Run development"
	@echo "  make dev             # Development workflow"
	@echo "  make watch END=dev   # Development workflow "

# Generic commands that work with ENV variable
build:
	@echo "Building Palworld Helper ($(if $(filter dev,$(ENV)),Development,Production))..."
	docker compose -f $(COMPOSE_FILE) build

run:
	@echo "Starting Palworld Helper ($(if $(filter dev,$(ENV)),Development,Production))..."
	docker compose -f $(COMPOSE_FILE) up -d
	@echo "Application is running at http://localhost:$(PORT)"

stop:
	@echo "Stopping Palworld Helper ($(if $(filter dev,$(ENV)),Development,Production))..."
	docker compose -f $(COMPOSE_FILE) down

logs:
	docker compose -f $(COMPOSE_FILE) logs -f

shell:
	docker compose -f $(COMPOSE_FILE) exec $(SERVICE_NAME) /bin/sh

restart:
	@echo "Restarting containers ($(if $(filter dev,$(ENV)),Development,Production))..."
	docker compose -f $(COMPOSE_FILE) down
	docker compose -f $(COMPOSE_FILE) up -d
	@echo "Application restarted at http://localhost:$(PORT)"

rebuild: clean build run

# Development workflow shortcut
dev:
	@$(MAKE) stop ENV=dev
	@$(MAKE) build ENV=dev
	@$(MAKE) run ENV=dev
	@echo "Development environment ready!"
	@$(MAKE) logs ENV=dev

# Auto-rebuild for development only
watch:
	@if [ "$(ENV)" != "dev" ]; then \
		echo "Error: watch command only available in development mode."; \
		echo "Usage: make watch ENV=dev"; \
		exit 1; \
	fi
	@echo "Starting development mode with auto-rebuild..."
	@echo "Requirements: inotify-tools (apt install inotify-tools)"
	@echo ""
	@echo "Initial build and start..."
	@$(MAKE) build ENV=dev
	@$(MAKE) run ENV=dev
	@echo "Now watching for file changes... Press Ctrl+C to stop"
	@while true; do \
		echo "Watching for file changes..."; \
		inotifywait -r -e modify,create,delete --exclude '\.(git|tmp|log)' . 2>/dev/null || { \
			echo "Error: inotifywait not found. Install with: apt install inotify-tools"; \
			exit 1; \
		}; \
		echo "Changes detected, rebuilding..."; \
		docker compose -f $(DEV_COMPOSE) build --quiet && \
		docker compose -f $(DEV_COMPOSE) up -d && \
		echo "Rebuild complete. Watching for more changes..."; \
	done

# Clean up - handles both environments
clean:
	@echo "Cleaning up all containers and images..."
	@docker compose -f $(PROD_COMPOSE) down -v --rmi all --remove-orphans 2>/dev/null || true
	@docker compose -f $(DEV_COMPOSE) down -v --rmi all --remove-orphans 2>/dev/null || true
	@docker system prune -f

# Legacy aliases for backward compatibility (optional)
dev-build:
	@$(MAKE) build ENV=dev

dev-run:
	@$(MAKE) run ENV=dev

dev-stop:
	@$(MAKE) stop ENV=dev

dev-logs:
	@$(MAKE) logs ENV=dev

dev-restart:
	@$(MAKE) restart ENV=dev

# Status and diagnostics
status:
	@echo "=== Docker Containers Status ==="
	@docker ps -a --filter "name=palworld" || echo "No palworld containers found"
	@echo ""
	@echo "=== Port 8080 Usage ==="
	@netstat -tulpn 2>/dev/null | grep :8080 || echo "Port 8080 is free"
	@echo ""
	@echo "=== Docker Compose Status ==="
	@echo "Production:"
	@docker compose -f $(PROD_COMPOSE) ps 2>/dev/null || echo "  No containers running"
	@echo "Development:"
	@docker compose -f $(DEV_COMPOSE) ps 2>/dev/null || echo "  No containers running"