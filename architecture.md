# Architecture Diagrams

## System Architecture Overview

```mermaid
graph TB
    Client[Client Applications]
    AG[API Gateway]
    
    subgraph AWS Cloud
        subgraph EKS Cluster
            subgraph Services
                subgraph User Service
                    USALB[Service ALB]
                    US[User Service Pods]
                end
                
                subgraph Tweet Service
                    TSALB[Service ALB]
                    TS[Tweet Service Pods]
                end
                
                subgraph Timeline Service
                    TLSALB[Service ALB]
                    TLS[Timeline Service Pods]
                end
            end
        end
        
        subgraph Data Storage
            RDS[(RDS PostgreSQL)]
            REDIS[(ElastiCache Redis)]
            CASS[(Cassandra)]
        end
        
        subgraph Message Bus
            SNS[SNS Topics]
            SQS[SQS Queues]
        end
        
        subgraph Infrastructure
            S3[S3 Storage]
            CF[CloudFront CDN]
        end
    end
    
    Client --> AG
    AG --> USALB & TSALB & TLSALB
    
    USALB --> US
    TSALB --> TS
    TLSALB --> TLS
    
    US --> RDS
    US --> REDIS
    
    TS --> RDS
    TS --> REDIS
    TS --> SNS
    
    TLS --> REDIS
    TLS --> CASS
    TLS --> SQS
    
    SNS --> SQS
    
    S3 --> CF
```

## Service Interaction - Tweet Creation Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant AG as API Gateway
    participant TS as Tweet Service
    participant SNS as SNS Topic
    participant SQS as SQS Queue
    participant TLS as Timeline Service
    participant RDS as PostgreSQL
    participant R as Redis Cache

    C->>AG: POST /tweets
    AG->>TS: Create Tweet
    TS->>RDS: Store Tweet
    TS->>SNS: Publish Tweet Event
    SNS->>SQS: Forward to Queue
    TLS->>SQS: Consume Tweet Event
    TLS->>R: Update Timeline Cache
    TLS->>AG: Tweet Created
    AG->>C: Success Response
```

## Timeline Retrieval Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant AG as API Gateway
    participant TLS as Timeline Service
    participant R as Redis Cache
    participant CASS as Cassandra

    C->>AG: GET /timeline
    AG->>TLS: Request Timeline
    TLS->>R: Check Cache
    alt Cache Hit
        R->>TLS: Return Cached Timeline
    else Cache Miss
        TLS->>CASS: Fetch Timeline
        CASS->>TLS: Return Timeline Data
        TLS->>R: Update Cache
    end
    TLS->>AG: Timeline Data
    AG->>C: Timeline Response
```

## User Service Flow

```mermaid
sequenceDiagram
    participant C as Client
    participant AG as API Gateway
    participant US as User Service
    participant RDS as PostgreSQL
    participant R as Redis Cache
    participant SNS as SNS Topic

    C->>AG: POST /follow
    AG->>US: Follow Request
    US->>RDS: Update Follow Relationship
    US->>R: Update Cache
    US->>SNS: Publish Follow Event
    SNS->>US: Success
    US->>AG: Follow Success
    AG->>C: Success Response
```

## Data Storage Architecture

```mermaid
graph LR
    subgraph Primary Storage
        RDS[(RDS PostgreSQL)]
        CASS[(Cassandra)]
    end
    
    subgraph Cache Layer
        RC[(Redis Cache)]
    end
    
    subgraph Services
        US[User Service]
        TS[Tweet Service]
        TLS[Timeline Service]
    end
    
    US --> RDS
    US --> RC
    
    TS --> RDS
    TS --> RC
    
    TLS --> CASS
    TLS --> RC
```

## Message Flow Architecture

```mermaid
graph LR
    subgraph Publishers
        TS[Tweet Service]
        US[User Service]
    end
    
    subgraph Message Bus
        SNS[SNS Topics]
        SQS[SQS Queues]
    end
    
    subgraph Consumers
        TLS[Timeline Service]
    end
    
    TS --> SNS
    US --> SNS
    SNS --> SQS
    SQS --> TLS
``` 