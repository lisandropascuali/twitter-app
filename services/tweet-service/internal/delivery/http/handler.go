package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lisandro/challenge/services/tweet-service/internal/domain"
)

// Handler handles HTTP requests
type Handler struct {
	tweetUseCase domain.TweetUseCase
}

// Tweet represents a tweet in the API
// @Description Tweet information
type Tweet struct {
	ID        uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174000"`
	UserID    uuid.UUID `json:"user_id" example:"123e4567-e89b-12d3-a456-426614174000"`
	Content   string    `json:"content" example:"Hello, this is my first tweet!"`
	CreatedAt string    `json:"created_at" example:"2024-06-07T22:04:25Z"`
	UpdatedAt string    `json:"updated_at" example:"2024-06-07T22:04:25Z"`
}

// CreateTweetRequest represents the request body for creating a tweet
// @Description Request body for creating a tweet
type CreateTweetRequest struct {
	UserID  uuid.UUID `json:"user_id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
	Content string    `json:"content" binding:"required" example:"Hello, this is my first tweet!"`
}

// UpdateTweetRequest represents the request body for updating a tweet
// @Description Request body for updating a tweet
type UpdateTweetRequest struct {
	Content string `json:"content" binding:"required" example:"Updated tweet content"`
}

// ErrorResponse represents an error response
// @Description Error response
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid request"`
}

// SuccessResponse represents a success response
// @Description Success response
type SuccessResponse struct {
	Message string `json:"message" example:"Tweet created successfully"`
}

func NewHandler(tweetUseCase domain.TweetUseCase) *Handler {
	return &Handler{
		tweetUseCase: tweetUseCase,
	}
}

// CreateTweet godoc
// @Summary Create a new tweet
// @Description Create a new tweet for a user
// @Tags tweets
// @Accept json
// @Produce json
// @Param tweet body CreateTweetRequest true "Tweet object"
// @Success 201 {object} Tweet
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tweets [post]
func (h *Handler) CreateTweet(c *fiber.Ctx) error {
	var req CreateTweetRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{Error: err.Error()})
	}

	tweet, err := h.tweetUseCase.CreateTweet(req.UserID, req.Content)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{Error: err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(tweet)
}
