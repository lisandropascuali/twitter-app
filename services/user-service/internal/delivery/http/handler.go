package http

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/lisandro/challenge/services/user-service/internal/domain"
)

type UserHandler struct {
	userUsecase domain.UserUsecase
}

// NewUserHandler creates a new user handler with its dependencies
func NewUserHandler(uu domain.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: uu,
	}
}

// GetAllUsers godoc
// @Summary Get all users
// @Description Get a list of all users in the system
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} map[string][]domain.User
// @Failure 500 {object} map[string]string
// @Router /users [get]
func (h *UserHandler) GetAllUsers(c *fiber.Ctx) error {
	log.Println("Getting all users")
	users, err := h.userUsecase.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get users",
		})
	}

	return c.JSON(fiber.Map{
		"users": users,
	})
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with the provided information
// @Tags users
// @Accept json
// @Produce json
// @Param user body domain.CreateUserRequest true "User information"
// @Success 201 {object} map[string]domain.User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users [post]
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req domain.CreateUserRequest
	
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	log.Printf("Creating new user: %+v", req)
	user, err := h.userUsecase.CreateUser(req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user": user,
	})
}

// Follow godoc
// @Summary Follow a user
// @Description Follow another user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param followedID path string true "ID of the user to follow"
// @Param X-User-ID header string true "ID of the current user"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{followedID}/follow [post]
func (h *UserHandler) Follow(c *fiber.Ctx) error {
	followerID := c.Get("X-User-ID")
	followedID := c.Params("followedID")

	if followerID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	// Check if the user is already following the target user
	following, err := h.userUsecase.GetFollowing(followerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to check following status",
		})
	}

	for _, user := range following {
		if user.ID == followedID {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "User is already following this user",
			})
		}
	}

	if err := h.userUsecase.Follow(followerID, followedID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to follow user",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

// Unfollow godoc
// @Summary Unfollow a user
// @Description Unfollow another user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Param followedID path string true "ID of the user to unfollow"
// @Param X-User-ID header string true "ID of the current user"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/{followedID}/follow [delete]
func (h *UserHandler) Unfollow(c *fiber.Ctx) error {
	followerID := c.Get("X-User-ID")
	followedID := c.Params("followedID")

	if followerID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}

	if err := h.userUsecase.Unfollow(followerID, followedID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to unfollow user",
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

// GetFollowing godoc
// @Summary Get following list
// @Description Get the list of users that the current user follows
// @Tags users
// @Accept json
// @Produce json
// @Param X-User-ID header string true "ID of the current user"
// @Success 200 {object} map[string][]domain.User
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/following [get]
func (h *UserHandler) GetFollowing(c *fiber.Ctx) error {
	userID := c.Get("X-User-ID")
	
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}
	log.Println("Getting following list for user:", userID)
	following, err := h.userUsecase.GetFollowing(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get following list",
		})
	}

	return c.JSON(fiber.Map{
		"following": following,
	})
}

// GetFollowers godoc
// @Summary Get followers list
// @Description Get the list of users that follow the current user
// @Tags users
// @Accept json
// @Produce json
// @Param X-User-ID header string true "ID of the current user"
// @Success 200 {object} map[string][]domain.User
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/followers [get]
func (h *UserHandler) GetFollowers(c *fiber.Ctx) error {
	userID := c.Get("X-User-ID")
	
	if userID == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User ID is required",
		})
	}
	log.Println("Getting followers list for user:", userID)
	followers, err := h.userUsecase.GetFollowers(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get followers list",
		})
	}

	return c.JSON(fiber.Map{
		"followers": followers,
	})
} 