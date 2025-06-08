# Tweet Service

This is the Tweet Service component of the microservices architecture. It handles tweet creation, retrieval, and deletion operations.

## Features

- Create tweets

## Prerequisites

- Go 1.21 or later
- Docker and Docker Compose
- PostgreSQL
- Redis
- AWS SNS (for production)

## Getting Started

1. Clone the repository
2. Install dependencies:
   ```bash
   make deps
   ```

3. Run the service locally:
   ```bash
   make run
   ```

Or using Docker Compose:
   ```bash
   make docker-compose-up
   ```

## API Endpoints

- `POST /tweets` - Create a new tweet

## Development

- Build the service:
  ```bash
  make build
  ```

- Run tests:
  ```bash
  make test
  ```

- Format code:
  ```bash
  make fmt
  ```

- Lint code:
  ```bash
  make lint
  ```

## Docker

Build and run using Docker:
```bash
make docker-build
make docker-run
```

## Configuration

The service can be configured using environment variables:

- `DB_HOST` - PostgreSQL host (default: localhost)
- `DB_USER` - PostgreSQL user (default: postgres)
- `DB_PASSWORD` - PostgreSQL password (default: postgres)
- `DB_NAME` - PostgreSQL database name (default: tweets)
- `DB_PORT` - PostgreSQL port (default: 5433)

## Architecture

The service follows clean architecture principles with the following layers:

- Domain: Core business logic and entities
- Use Case: Application business rules
- Repository: Data access layer
- Delivery: API handlers and external interfaces

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request 