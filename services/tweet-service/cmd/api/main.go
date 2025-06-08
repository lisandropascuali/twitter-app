package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/lisandro/challenge/services/tweet-service/config"
	_ "github.com/lisandro/challenge/services/tweet-service/docs"
	"github.com/lisandro/challenge/services/tweet-service/internal/delivery/http"
	pgRepo "github.com/lisandro/challenge/services/tweet-service/internal/repository/postgres"
	"github.com/lisandro/challenge/services/tweet-service/internal/usecase"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	log.Println("Starting user service...")

	// Initialize PostgreSQL connection
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		getEnvOrDefault("DB_HOST", "localhost"),
		getEnvOrDefault("DB_USER", "user_service"),
		getEnvOrDefault("DB_PASSWORD", "user_service_pass"),
		getEnvOrDefault("DB_NAME", "user_service_db"),
		getEnvOrDefault("DB_PORT", "5432"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run migrations
	if err := pgRepo.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize repositories
	tweetRepo := pgRepo.NewTweetRepository(db)
	if err != nil {
		log.Fatalf("Failed to create SNS client: %v", err)
	}

	// Initialize usecase with its dependencies
	tweetUsecase := usecase.NewTweetUseCase(tweetRepo)

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