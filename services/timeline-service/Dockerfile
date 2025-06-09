FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/timeline-service ./cmd/api

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/timeline-service .
COPY --from=builder /app/config ./config

EXPOSE 8082

CMD ["./timeline-service"] 