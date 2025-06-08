package http

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/lisandro/challenge/services/tweet-service/internal/domain"
)

type Server struct {
	app         *fiber.App
	tweetHandler *Handler
}

func NewServer(tu domain.TweetUseCase) *Server {
	// Create Fiber app with custom config
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	
	// Add logger middleware
	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${latency} ${method} ${path}\n",
		TimeFormat: "2006-01-02 15:04:05",
		TimeZone:   "Local",
	}))

	// Add Swagger UI
	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL:         "/swagger/doc.json",
		DeepLinking: true,
	}))
	
	// Create handlers with dependencies
	handler := NewHandler(tu)
	
	// Register routes
	RegisterRoutes(app, handler)
	
	return &Server{
		app:         app,
		tweetHandler: handler,
	}
}

func (s *Server) Start(address string) error {
	log.Printf("Starting server on %s", address)
	return s.app.Listen(address)
} 