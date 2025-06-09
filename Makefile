.PHONY: help setup up down restart build run test clean logs check-deps install-deps setup-infra run-all stop-all

# Colors for output
GREEN := \033[0;32m
YELLOW := \033[0;33m
RED := \033[0;31m
NC := \033[0m # No Color

# Default target
help:
	@echo "$(GREEN)Twitter-like Platform - Available Commands:$(NC)"
	@echo ""
	@echo "$(YELLOW)Quick Start:$(NC)"
	@echo "  make setup    - Complete setup (check deps, setup infra, build services)"
	@echo "  make run-all  - Run all services after setup"
	@echo "  make stop-all - Stop all services"
	@echo ""
	@echo "$(YELLOW)Individual Commands:$(NC)"
	@echo "  make check-deps     - Check if required dependencies are installed"
	@echo "  make install-deps   - Install missing dependencies (macOS only)"
	@echo "  make setup-infra    - Setup infrastructure (databases, tables, indexes)"
	@echo "  make build          - Build all services"
	@echo "  make up             - Start all docker-compose services"
	@echo "  make down           - Stop all docker-compose services"
	@echo "  make test           - Run tests for all services"
	@echo "  make clean          - Clean up all build artifacts and containers"
	@echo "  make logs           - Show logs from all services"
	@echo ""
	@echo "$(YELLOW)Service URLs (after running):$(NC)"
	@echo "  User Service:     http://localhost:8080"
	@echo "  Tweet Service:    http://localhost:8081"
	@echo "  Timeline Service: http://localhost:8082"

# Check if required dependencies are installed
check-deps:
	@echo "$(GREEN)Checking dependencies...$(NC)"
	@command -v docker >/dev/null 2>&1 || { echo "$(RED)Docker is required but not installed.$(NC)" >&2; exit 1; }
	@command -v docker-compose >/dev/null 2>&1 || { echo "$(RED)Docker Compose is required but not installed.$(NC)" >&2; exit 1; }
	@command -v go >/dev/null 2>&1 || { echo "$(RED)Go is required but not installed.$(NC)" >&2; exit 1; }
	@command -v aws >/dev/null 2>&1 || { echo "$(RED)AWS CLI is required but not installed.$(NC)" >&2; exit 1; }
	@echo "$(GREEN)✓ All dependencies are installed$(NC)"

# Install missing dependencies (macOS only)
install-deps:
	@echo "$(GREEN)Installing dependencies (macOS)...$(NC)"
	@if ! command -v brew >/dev/null 2>&1; then \
		echo "$(RED)Homebrew is required but not installed. Please install it first.$(NC)"; \
		exit 1; \
	fi
	@if ! command -v docker >/dev/null 2>&1; then \
		echo "Installing Docker..."; \
		brew install --cask docker; \
	fi
	@if ! command -v go >/dev/null 2>&1; then \
		echo "Installing Go..."; \
		brew install go; \
	fi
	@if ! command -v aws >/dev/null 2>&1; then \
		echo "Installing AWS CLI..."; \
		brew install awscli; \
	fi
	@echo "$(GREEN)✓ Dependencies installed$(NC)"

# Setup infrastructure (databases, tables, indexes)
setup-infra:
	@echo "$(GREEN)Setting up infrastructure...$(NC)"
	@echo "$(YELLOW)Starting PostgreSQL for User Service...$(NC)"
	@cd services/user-service && make up
	@echo "$(YELLOW)Starting DynamoDB and OpenSearch for Tweet Service...$(NC)"
	@cd services/tweet-service && make up
	@echo "$(YELLOW)Waiting for services to be ready...$(NC)"
	@sleep 10
	@echo "$(YELLOW)Setting up User Service database...$(NC)"
	@cd services/user-service && make seed
	@echo "$(YELLOW)Setting up Tweet Service infrastructure...$(NC)"
	@cd services/tweet-service && make create-table create-opensearch-index
	@echo "$(GREEN)✓ Infrastructure setup complete$(NC)"

# Build all services
build:
	@echo "$(GREEN)Building all services...$(NC)"
	@echo "$(YELLOW)Building User Service...$(NC)"
	@cd services/user-service && make build
	@echo "$(YELLOW)Building Tweet Service...$(NC)"
	@cd services/tweet-service && make build
	@echo "$(YELLOW)Building Timeline Service...$(NC)"
	@cd services/timeline-service && make build
	@echo "$(GREEN)✓ All services built$(NC)"

# Start all docker-compose services
up:
	@echo "$(GREEN)Starting all docker-compose services...$(NC)"
	@cd services/user-service && make up
	@cd services/tweet-service && make up
	@echo "$(GREEN)✓ All docker services started$(NC)"

# Stop all docker-compose services
down:
	@echo "$(GREEN)Stopping all docker-compose services...$(NC)"
	@cd services/user-service && make down
	@cd services/tweet-service && make down
	@cd services/timeline-service && make docker-stop
	@echo "$(GREEN)✓ All docker services stopped$(NC)"

# Complete setup - check deps, setup infra, build services
setup: check-deps setup-infra build
	@echo "$(GREEN)✓ Complete setup finished!$(NC)"
	@echo "$(YELLOW)Ready to run services with: make run-all$(NC)"

# Run all services (use in separate terminals or background)
run-all:
	@echo "$(GREEN)Starting all services...$(NC)"
	@echo "$(YELLOW)Note: This will run services in the background$(NC)"
	@echo "$(YELLOW)Use 'make logs' to see output or 'make stop-all' to stop$(NC)"
	@echo "$(YELLOW)Starting User Service on port 8080...$(NC)"
	@cd services/user-service && nohup make run > ../user-service.log 2>&1 &
	@sleep 5
	@echo "$(YELLOW)Starting Tweet Service on port 8081...$(NC)"
	@cd services/tweet-service && nohup make run > ../tweet-service.log 2>&1 &
	@sleep 5
	@echo "$(YELLOW)Starting Timeline Service on port 8082...$(NC)"
	@cd services/timeline-service && nohup make run > ../timeline-service.log 2>&1 &
	@sleep 5
	@echo "$(GREEN)✓ All services started$(NC)"
	@echo "$(YELLOW)Service URLs:$(NC)"
	@echo "  User Service:     http://localhost:8080"
	@echo "  Tweet Service:    http://localhost:8081"
	@echo "  Timeline Service: http://localhost:8082"
	@echo ""
	@echo "$(YELLOW)Use 'make logs' to see service logs$(NC)"

# Stop all running services
stop-all:
	@echo "$(GREEN)Stopping all services...$(NC)"
	@pkill -f "user-service" || true
	@pkill -f "tweet-service" || true
	@pkill -f "timeline-service" || true
	@make down
	@echo "$(GREEN)✓ All services stopped$(NC)"

# Show logs from all services
logs:
	@echo "$(GREEN)Showing logs from all services...$(NC)"
	@echo "$(YELLOW)=== User Service Logs ====$(NC)"
	@tail -n 20 services/user-service.log 2>/dev/null || echo "No user service logs found"
	@echo "$(YELLOW)=== Tweet Service Logs ====$(NC)"
	@tail -n 20 services/tweet-service.log 2>/dev/null || echo "No tweet service logs found"
	@echo "$(YELLOW)=== Timeline Service Logs ====$(NC)"
	@tail -n 20 services/timeline-service.log 2>/dev/null || echo "No timeline service logs found"
	@echo "$(YELLOW)=== Docker Compose Logs ====$(NC)"
	@cd services/user-service && docker-compose logs --tail=10 || true
	@cd services/tweet-service && docker-compose logs --tail=10 || true

# Run tests for all services
test:
	@echo "$(GREEN)Running tests for all services...$(NC)"
	@echo "$(YELLOW)Testing User Service...$(NC)"
	@cd services/user-service && make test
	@echo "$(YELLOW)Testing Tweet Service...$(NC)"
	@cd services/tweet-service && make test
	@echo "$(YELLOW)Testing Timeline Service...$(NC)"
	@cd services/timeline-service && make test
	@echo "$(GREEN)✓ All tests completed$(NC)"

# Clean up everything
clean:
	@echo "$(GREEN)Cleaning up all services and containers...$(NC)"
	@cd services/user-service && make clean
	@cd services/tweet-service && make clean
	@cd services/timeline-service && make clean
	@rm -f services/*.log
	@echo "$(GREEN)✓ Cleanup completed$(NC)"

# Restart everything
restart: stop-all setup run-all

# Check service health
health:
	@echo "$(GREEN)Checking service health...$(NC)"
	@echo "$(YELLOW)User Service:$(NC)"
	@curl -s http://localhost:8080/health || echo "$(RED)User Service not responding$(NC)"
	@echo "$(YELLOW)Tweet Service:$(NC)"
	@curl -s http://localhost:8081/health || echo "$(RED)Tweet Service not responding$(NC)"
	@echo "$(YELLOW)Timeline Service:$(NC)"
	@curl -s http://localhost:8082/health || echo "$(RED)Timeline Service not responding$(NC)" 