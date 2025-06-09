package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lisandro/timeline-service/config"
	_ "github.com/lisandro/timeline-service/docs" // This is important!
	"github.com/lisandro/timeline-service/internal/client"
	"github.com/lisandro/timeline-service/internal/delivery/http"
	"github.com/lisandro/timeline-service/internal/usecase"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Timeline Service API
// @version 1.0
// @description Service that provides timeline functionality for the Twitter clone
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8082
// @BasePath /api/v1
// @schemes http

func main() {
	config.InitLogger()
	log.Println("Starting timeline service...")

	// Initialize clients
	userServiceURL := getEnvOrDefault("USER_SERVICE_URL", "http://localhost:8080/api/v1/users")
	tweetServiceURL := getEnvOrDefault("TWEET_SERVICE_URL", "http://localhost:8081/api/v1")
	
	log.Printf("User service URL: %s", userServiceURL)
	log.Printf("Tweet service URL: %s", tweetServiceURL)
	
	userClient := client.NewUserClient(userServiceURL)
	tweetClient := client.NewTweetClient(tweetServiceURL)

	// Initialize usecase
	timelineUseCase := usecase.NewTimelineUseCase(userClient, tweetClient)

	// Initialize handler
	timelineHandler := http.NewTimelineHandler(timelineUseCase)

	// Initialize router
	router := gin.New()
	
	// Add custom logger middleware
	router.Use(customLogger())
	router.Use(gin.Recovery())

	// Register routes
	v1 := router.Group("/api/v1")
	{
		v1.GET("/timeline", timelineHandler.GetTimeline)
	}

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	port := getEnvOrDefault("PORT", "8082")
	log.Printf("Starting server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// customLogger creates a custom logger middleware that matches the user-service format
func customLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate request duration
		latency := time.Since(start)

		// Get status code
		statusCode := c.Writer.Status()

		// Format path with query string if present
		if raw != "" {
			path = path + "?" + raw
		}

		// Log in format similar to user-service: [timestamp] status - latency method path
		log.Printf("%d - %v %s %s", statusCode, latency, c.Request.Method, path)
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 