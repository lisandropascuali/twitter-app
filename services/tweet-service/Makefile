.PHONY: up down restart build run test clean seed migrate wait-for-postgres

# Docker compose commands
up:
	docker-compose up -d

down:
	docker-compose down

restart: down up

# Build and run commands
build:
	go build -o bin/tweet-service cmd/api/main.go

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
	@until docker exec tweet-service-postgres-1 pg_isready -h localhost -p 5432 -U tweet_service; do \
		echo "PostgreSQL is unavailable - sleeping"; \
		sleep 1; \
	done
	@echo "PostgreSQL is ready"
	@sleep 2

# Run database migrations
migrate: wait-for-postgres
	@echo "Running database migrations..."
	@docker exec -i tweet-service-postgres-1 psql -U tweet_service -d tweets < scripts/migrations.sql

# Seed database with mock data
seed: migrate
	@echo "Populating database with sample tweets..."
	@docker exec -i tweet-service-postgres-1 psql -U tweet_service -d tweets < scripts/seed.sql

# Build Docker image
docker-build:
	docker build -t tweet-service .

# Run Docker container
docker-run:
	docker run -p 8081:8081 tweet-service

# Run with docker-compose
docker-compose-up:
	docker-compose up --build

# Stop docker-compose
docker-compose-down:
	docker-compose down

# Generate mocks
generate-mocks:
	mockgen -source=internal/domain/tweet.go -destination=internal/mocks/tweet_mock.go -package=mocks

# Update dependencies
update-deps:
	go get -u ./...
	go mod tidy

# Lint code
lint:
	golangci-lint run

# Format code
fmt:
	go fmt ./... 