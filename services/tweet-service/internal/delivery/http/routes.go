package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// RegisterRoutes registers all the tweet routes
func RegisterRoutes(app *fiber.App, handler *Handler) {
	// Swagger documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

	api := app.Group("/api/v1")
	tweets := api.Group("/tweets")

	// @Summary Create a new tweet
	// @Description Create a new tweet
	// @Tags tweets
	// @Accept json
	// @Produce json
	// @Param tweet body CreateTweetRequest true "Tweet object"
	// @Success 201 {object} map[string]interface{}
	// @Failure 400 {object} map[string]string
	// @Failure 500 {object} map[string]string
	// @Router /tweets [post]
	tweets.Post("", handler.CreateTweet)

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})
} 