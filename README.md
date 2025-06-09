# Twitter-like Platform Architecture

## Quick Start 🚀

For evaluators and developers, you can get the entire platform running with just a few commands:

### Prerequisites
- Docker & Docker Compose
- Go (1.19+)
- AWS CLI (for local DynamoDB and OpenSearch)
- Make

### One-Command Setup
```bash
# Complete setup and run all services
make setup && make run-all
```

### Step-by-Step Setup
```bash
# 1. Check dependencies
make check-deps

# 2. Install missing dependencies (macOS only)
make install-deps

# 3. Setup infrastructure (databases, tables, indexes)
make setup-infra

# 4. Build all services
make build

# 5. Run all services
make run-all
```

### Available Commands
- `make help` - Show all available commands
- `make setup` - Complete setup (deps + infra + build)
- `make run-all` - Run all services
- `make stop-all` - Stop all services
- `make logs` - Show logs from all services
- `make test` - Run all tests
- `make clean` - Clean up everything
- `make health` - Check service health

### Service URLs
After running `make run-all`, the services will be available at:
- **User Service**: http://localhost:8080
- **Tweet Service**: http://localhost:8081  
- **Timeline Service**: http://localhost:8082

## Overview
This project implements a scalable, resilient, and high-performance Twitter-like platform using microservices architecture. The system is designed to handle millions of users while maintaining optimal read performance and eventual consistency.

## Tech Stack
- **Programming Language**: Go (Golang)
- **Databases**: 
  - DynamoDB (Tweet storage)
  - OpenSearch (Tweet search and queries)
  - PostgreSQL (User data)
- **Cloud Provider**: AWS
- **Architecture Pattern**: Clean Architecture
- **Infrastructure**: Docker + Kubernetes
- **Caching**: Redis

## System Architecture

### Microservices Breakdown

1. **User Service**
   - User profile management
   - User search and recommendations
   - Follow/Unfollow functionality
   - Follower/Following relationships
   - Technologies:
     - PostgreSQL for user data
     - Redis for caching user profiles

2. **Tweet Service**
   - Tweet creation and management
   - Tweet validation
   - Tweet storage and retrieval
   - Tweet search and queries
   - Technologies:
     - DynamoDB for tweet storage
     - OpenSearch for tweet search and queries
     - Redis for hot tweets caching

3. **Timeline Service**
   - Timeline generation and management
   - Feed aggregation
   - Technologies:
     - Redis for timeline caching
     - HTTP calls to User Service for following relationships
     - HTTP calls to Tweet Service for tweet data
     - Acts as an aggregator service

4. **API Gateway**
   - Request routing
   - Rate limiting
   - Authentication & Authorization
   - Technologies:
     - AWS API Gateway
     - Custom Go implementation

### Data Flow

1. **Tweet Creation Flow**
   ```
   Client -> API Gateway -> Tweet Service -> DynamoDB
                                        -> OpenSearch
                                        -> Redis Cache
   ```

2. **Timeline Retrieval Flow**
   ```
   Client -> API Gateway -> Timeline Service -> Redis Cache
                                           -> User Service (get following users)
                                           -> Tweet Service (get tweets)
   ```

### Data Storage Strategy

1. **Tweet Storage**
   - **DynamoDB**: Primary storage for tweets
     - Stores complete tweet data
     - Optimized for tweet retrieval by ID
     - Enables fast access to individual tweets
   
   - **OpenSearch**: Tweet search and queries
     - Stores tweet metadata and content for search
     - Optimized for tweet queries and full-text search
     - Enables efficient tweet search capabilities

2. **Data Access Pattern**
   - Fanout on read strategy
   - Timeline Service aggregates data from User and Tweet services
   - Tweet Service uses OpenSearch for efficient tweet queries
   - Redis cache provides fast access to recent timelines
   - Eventual consistency model for timeline updates

### Scalability Considerations

1. **Read Optimization**
   - Aggressive caching with Redis
   - OpenSearch for efficient timeline queries
   - Read replicas for PostgreSQL
   - Content Delivery Network (CDN) for static content

2. **Write Scalability**
   - DynamoDB for high-throughput tweet storage
   - Asynchronous processing for non-critical operations

3. **High Availability**
   - Multi-AZ deployment
   - Kubernetes for container orchestration
   - Circuit breakers for service resilience
   - Service mesh for improved reliability

## Clean Architecture Implementation

Each microservice follows Clean Architecture principles with the following layers:

1. **Domain Layer**
   - Business entities
   - Domain rules
   - Interface definitions

2. **Use Case Layer**
   - Application business rules
   - Use case implementations
   - Port interfaces

3. **Interface Layer**
   - Controllers
   - Presenters
   - Gateway implementations

4. **Infrastructure Layer**
   - Database implementations
   - External service implementations
   - Framework-specific code

## AWS Infrastructure

- **Compute**: EKS (Elastic Kubernetes Service)
- **Database**: Amazon RDS for PostgreSQL
- **Caching**: Amazon ElastiCache (Redis)
- **Message Queue**: Amazon SNS / SQS
- **Storage**: S3 for media storage
- **CDN**: CloudFront
- **Load Balancing**: ALB (Application Load Balancer)

## Development Setup

[To be added: Local development setup instructions]

## Deployment

[To be added: Deployment instructions and pipeline details]

## Performance Considerations

1. **Caching Strategy**
   - Multi-level caching
   - Cache warming for popular content
   - Cache invalidation patterns

2. **Database Optimization**
   - Proper indexing
   - Query optimization
   - Connection pooling
   - Read replicas

3. **Monitoring and Alerting**
   - Service metrics
   - Business metrics
   - SLO/SLA monitoring
   - Error tracking and alerting

## Security Measures

1. **API Security**
   - Rate limiting
   - JWT authentication
   - API key management
   - Request validation

2. **Data Security**
   - Encryption at rest
   - Encryption in transit
   - Secure communication between services

## Future Improvements

1. **Feature Enhancements**
   - Media handling
   - Hashtag support
   - User mentions
   - Direct messaging

2. **Technical Improvements**
   - Real-time notifications
   - Advanced analytics
   - Machine learning integration

## Contributing

[To be added: Contribution guidelines]

## License

[To be added: License information] 