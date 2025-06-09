.PHONY: up down restart build run test clean create-table check-aws-cli create-opensearch-index

# Docker compose commands
up:
	docker-compose up -d

down:
	docker-compose down

restart: down up

# Build and run commands
build:
	go build -o bin/tweet-service cmd/api/main.go

run: up create-table create-opensearch-index
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

# Check if AWS CLI is installed
check-aws-cli:
	@if ! command -v aws &> /dev/null; then \
		echo "AWS CLI is not installed. Please install it first:"; \
		echo "  - For macOS: brew install awscli"; \
		echo "  - For Linux: pip install awscli"; \
		exit 1; \
	fi

# Create DynamoDB table
create-table: check-aws-cli
	@echo "Creating DynamoDB table..."
	@chmod +x scripts/create-table.sh
	@./scripts/create-table.sh

# Create OpenSearch index
create-opensearch-index:
	@echo "Creating OpenSearch index..."
	@chmod +x scripts/create-opensearch-index.sh
	@./scripts/create-opensearch-index.sh

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