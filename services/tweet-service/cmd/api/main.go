package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/lisandro/challenge/services/tweet-service/config"
	_ "github.com/lisandro/challenge/services/tweet-service/docs" // Import generated docs
	"github.com/lisandro/challenge/services/tweet-service/internal/delivery/http"
	dynamorepo "github.com/lisandro/challenge/services/tweet-service/internal/repository/dynamodb"
	opensearchrepo "github.com/lisandro/challenge/services/tweet-service/internal/repository/opensearch"
	"github.com/lisandro/challenge/services/tweet-service/internal/usecase"
	"github.com/opensearch-project/opensearch-go/v2"
)

// @title Tweet Service API
// @version 1.0
// @description This is the tweet service API documentation.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8081
// @BasePath /api/v1
// @schemes http

func main() {
	config.InitLogger()
	log.Println("Starting tweet service...")

	// Initialize AWS SDK with static credentials for LocalStack
	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(getEnvOrDefault("AWS_REGION", "us-east-1")),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			getEnvOrDefault("AWS_ACCESS_KEY_ID", "test"),
			getEnvOrDefault("AWS_SECRET_ACCESS_KEY", "test"),
			"",
		)),
	)
	if err != nil {
		log.Fatalf("Failed to load AWS config: %v", err)
	}

	// Initialize DynamoDB client with custom endpoint
	dynamoClient := dynamodb.NewFromConfig(awsCfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(getEnvOrDefault("DYNAMODB_ENDPOINT", "http://localhost:4566"))
	})

	// Initialize OpenSearch client
	opensearchClient, err := opensearch.NewClient(opensearch.Config{
		Addresses: []string{getEnvOrDefault("OPENSEARCH_ENDPOINT", "http://localhost:9200")},
	})
	if err != nil {
		log.Fatalf("Failed to create OpenSearch client: %v", err)
	}

	// Initialize repositories
	tweetRepo := dynamorepo.NewTweetRepository(dynamoClient, getEnvOrDefault("DYNAMODB_TABLE", "tweets"))
	searchRepo := opensearchrepo.NewSearchRepository(opensearchClient)

	// Initialize usecase with its dependencies
	tweetUsecase := usecase.NewTweetUseCase(tweetRepo, searchRepo)

	// Initialize HTTP server with its dependencies
	server := http.NewServer(tweetUsecase)

	// Start server in a goroutine
	go func() {
		port := getEnvOrDefault("PORT", "8081")
		if err := server.Start(":" + port); err != nil {
			log.Printf("Server error: %v\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 