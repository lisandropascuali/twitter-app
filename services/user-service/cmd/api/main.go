package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-redis/redis/v8"
	"github.com/lisandro/challenge/services/user-service/config"
	_ "github.com/lisandro/challenge/services/user-service/docs" // This is important!
	"github.com/lisandro/challenge/services/user-service/internal/delivery/http"
	"github.com/lisandro/challenge/services/user-service/internal/repository"
	pgRepo "github.com/lisandro/challenge/services/user-service/internal/repository/postgres"
	redisRepo "github.com/lisandro/challenge/services/user-service/internal/repository/redis"
	"github.com/lisandro/challenge/services/user-service/internal/usecase"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// @title User Service API
// @version 1.0
// @description This is the user service API documentation.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
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

	// Initialize Redis connection
	rdb := redis.NewClient(&redis.Options{
		Addr:     getEnvOrDefault("REDIS_ADDR", "localhost:6379"),
		Password: getEnvOrDefault("REDIS_PASSWORD", ""),
		DB:       0,
	})

	// Initialize repositories
	pgRepository := pgRepo.NewPostgresRepository(db)
	redisRepository := redisRepo.NewRedisRepository(rdb)
	
	// Initialize composite repository
	userRepo := repository.NewCompositeRepository(pgRepository, redisRepository)

	// Initialize usecase with its dependencies
	userUsecase := usecase.NewUserUsecase(userRepo)

	// Initialize HTTP server with its dependencies
	server := http.NewServer(userUsecase)

	// Start server in a goroutine
	go func() {
		port := getEnvOrDefault("PORT", "8080")
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