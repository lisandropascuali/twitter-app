.PHONY: up down restart build run test clean seed migrate wait-for-postgres

# Docker compose commands
up:
	docker-compose up -d

down:
	docker-compose down

restart: down up

# Build and run commands
build:
	go build -o bin/user-service cmd/api/main.go

run: up wait-for-postgres
	go run cmd/api/main.go

# Testing commands
test:
	go test -v ./...

# Clean up
clean:
	docker-compose down -v
	rm -rf bin/

# Show logs
logs:
	docker-compose logs -f

# Install dependencies
deps:
	go mod tidy
	go mod download

# Wait for PostgreSQL to be ready
wait-for-postgres:
	@echo "Waiting for PostgreSQL to be ready..."
	@until docker exec user-service-postgres-1 pg_isready -h localhost -p 5432 -U user_service; do \
		echo "PostgreSQL is unavailable - sleeping"; \
		sleep 1; \
	done
	@echo "PostgreSQL is ready"
	@sleep 2

# Run database migrations
migrate: wait-for-postgres
	@echo "Running database migrations..."
	@docker exec -i user-service-postgres-1 psql -U user_service -d user_service_db < scripts/migrations.sql

# Seed database with mock data
seed: migrate
	@echo "Populating database with sample users..."
	@docker exec -i user-service-postgres-1 psql -U user_service -d user_service_db < scripts/seed.sql