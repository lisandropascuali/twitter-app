# Build stage
FROM golang:1.21-alpine

# Install AWS CLI and required dependencies
RUN apk add --no-cache \
    python3 \
    py3-pip \
    curl \
    bash \
    && pip3 install --no-cache-dir awscli

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o main cmd/api/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .

# Expose port
EXPOSE 8081

# Run the application
CMD ["./main"] 