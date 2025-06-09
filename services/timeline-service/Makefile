.PHONY: build run test clean docker-build docker-run debug

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Binary name
BINARY_NAME=timeline-service

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) ./cmd/api

# Run the application
run:
	$(GOBUILD) -o $(BINARY_NAME) ./cmd/api
	./$(BINARY_NAME)

# Test the application
test:
	$(GOTEST) -v ./...

# Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Build docker image
docker-build:
	docker build -t timeline-service .

# Run docker container
docker-run:
	docker-compose up --build

# Stop docker container
docker-stop:
	docker-compose down

# Install dependencies
deps:
	$(GOGET) -v ./...

# Format code
fmt:
	$(GOCMD) fmt ./...

# Run linter
lint:
	golangci-lint run
