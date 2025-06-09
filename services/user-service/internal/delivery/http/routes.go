package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

// RegisterRoutes registers all the user routes
func RegisterRoutes(app *fiber.App, handler *UserHandler) {
	// Swagger documentation
	app.Get("/swagger/*", swagger.HandlerDefault)

	api := app.Group("/api/v1")
	users := api.Group("/users")

	// User CRUD operations
	// @Summary Get all users
	// @Description Get a list of all users in the system
	// @Tags users
	// @Accept json
	// @Produce json
	// @Success 200 {object} map[string][]domain.User
	// @Failure 500 {object} map[string]string
	// @Router /users [get]
	users.Get("/", handler.GetAllUsers)

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
	users.Post("/", handler.CreateUser)

	// Following functionality
	// @Summary Follow a user
	// @Description Follow another user by their ID
	// @Tags users
	// @Accept json
	// @Produce json
	// @Param followedID path string true "ID of the user to follow"
	// @Header 200 {string} X-User-ID "ID of the current user"
	// @Success 200 {object} map[string]interface{}
	// @Failure 401 {object} map[string]string
	// @Failure 500 {object} map[string]string
	// @Router /users/{followedID}/follow [post]
	users.Post("/:followedID/follow", handler.Follow)

	// @Summary Unfollow a user
	// @Description Unfollow another user by their ID
	// @Tags users
	// @Accept json
	// @Produce json
	// @Param followedID path string true "ID of the user to unfollow"
	// @Header 200 {string} X-User-ID "ID of the current user"
	// @Success 200 {object} map[string]interface{}
	// @Failure 401 {object} map[string]string
	// @Failure 500 {object} map[string]string
	// @Router /users/{followedID}/follow [delete]
	users.Delete("/:followedID/follow", handler.Unfollow)

	// @Summary Get following list
	// @Description Get the list of users that the current user follows
	// @Tags users
	// @Accept json
	// @Produce json
	// @Header 200 {string} X-User-ID "ID of the current user"
	// @Success 200 {object} map[string][]domain.User
	// @Failure 401 {object} map[string]string
	// @Failure 500 {object} map[string]string
	// @Router /users/following [get]
	users.Get("/following", handler.GetFollowing)

	// @Summary Get followers list
	// @Description Get the list of users that follow the current user
	// @Tags users
	// @Accept json
	// @Produce json
	// @Header 200 {string} X-User-ID "ID of the current user"
	// @Success 200 {object} map[string][]domain.User
	// @Failure 401 {object} map[string]string
	// @Failure 500 {object} map[string]string
	// @Router /users/followers [get]
	users.Get("/followers", handler.GetFollowers)
} 