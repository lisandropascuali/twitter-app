package http

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/lisandro/challenge/services/user-service/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserUsecase is a mock implementation of domain.UserUsecase
type MockUserUsecase struct {
	mock.Mock
}

func (m *MockUserUsecase) Follow(followerID, followedID string) error {
	args := m.Called(followerID, followedID)
	return args.Error(0)
}

func (m *MockUserUsecase) Unfollow(followerID, followedID string) error {
	args := m.Called(followerID, followedID)
	return args.Error(0)
}

func (m *MockUserUsecase) GetFollowing(userID string) ([]domain.User, error) {
	args := m.Called(userID)
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockUserUsecase) GetFollowers(userID string) ([]domain.User, error) {
	args := m.Called(userID)
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockUserUsecase) CreateUser(req domain.CreateUserRequest) (*domain.User, error) {
	args := m.Called(req)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserUsecase) GetAllUsers() ([]domain.User, error) {
	args := m.Called()
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockUserUsecase) GetUser(id string) (*domain.User, error) {
	args := m.Called(id)
	return args.Get(0).(*domain.User), args.Error(1)
}


func setupTest() (*fiber.App, *MockUserUsecase, *UserHandler) {
	app := fiber.New()
	mockUsecase := new(MockUserUsecase)
	handler := NewUserHandler(mockUsecase)
	return app, mockUsecase, handler
}

func TestUserHandler_Follow(t *testing.T) {
	tests := []struct {
		name           string
		followerID     string
		followedID     string
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "successful follow",
			followerID:     "user1",
			followedID:     "user2",
			mockError:      nil,
			expectedStatus: fiber.StatusOK,
			expectedBody:   nil,
		},
		{
			name:           "missing user ID",
			followerID:     "",
			followedID:     "user2",
			mockError:      nil,
			expectedStatus: fiber.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "User ID is required",
			},
		},
		{
			name:           "usecase error",
			followerID:     "user1",
			followedID:     "user2",
			mockError:      errors.New("database error"),
			expectedStatus: fiber.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "Failed to follow user",
			},
		},
		{
			name:           "already following",
			followerID:     "user1",
			followedID:     "user2",
			mockError:      nil,
			expectedStatus: fiber.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "User is already following this user",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockUsecase, handler := setupTest()
			app.Post("/:followedID/follow", handler.Follow)

			if tt.name == "already following" {
				mockUsecase.On("GetFollowing", tt.followerID).Return([]domain.User{{ID: tt.followedID, Username: "user2"}}, nil)
			} else if tt.followerID != "" {
				mockUsecase.On("Follow", tt.followerID, tt.followedID).Return(tt.mockError)
			}

			req := httptest.NewRequest("POST", "/"+tt.followedID+"/follow", nil)
			if tt.followerID != "" {
				req.Header.Set("X-User-ID", tt.followerID)
			}

			resp, _ := app.Test(req)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectedBody != nil {
				var body map[string]interface{}
				json.NewDecoder(resp.Body).Decode(&body)
				assert.Equal(t, tt.expectedBody, body)
			}

			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestUserHandler_Unfollow(t *testing.T) {
	tests := []struct {
		name           string
		followerID     string
		followedID     string
		mockError      error
		expectedStatus int
		expectedBody   map[string]interface{}
	}{
		{
			name:           "successful unfollow",
			followerID:     "user1",
			followedID:     "user2",
			mockError:      nil,
			expectedStatus: fiber.StatusOK,
			expectedBody:   nil,
		},
		{
			name:           "missing user ID",
			followerID:     "",
			followedID:     "user2",
			mockError:      nil,
			expectedStatus: fiber.StatusUnauthorized,
			expectedBody: map[string]interface{}{
				"error": "User ID is required",
			},
		},
		{
			name:           "usecase error",
			followerID:     "user1",
			followedID:     "user2",
			mockError:      errors.New("database error"),
			expectedStatus: fiber.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "Failed to unfollow user",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockUsecase, handler := setupTest()
			app.Delete("/:followedID/follow", handler.Unfollow)

			if tt.followerID != "" {
				mockUsecase.On("Unfollow", tt.followerID, tt.followedID).Return(tt.mockError)
			}

			req := httptest.NewRequest("DELETE", "/"+tt.followedID+"/follow", nil)
			if tt.followerID != "" {
				req.Header.Set("X-User-ID", tt.followerID)
			}

			resp, _ := app.Test(req)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			if tt.expectedBody != nil {
				var body map[string]interface{}
				json.NewDecoder(resp.Body).Decode(&body)
				assert.Equal(t, tt.expectedBody, body)
			}

			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestUserHandler_GetFollowing(t *testing.T) {
	mockUsers := []domain.User{
		{ID: "user2", Username: "user2"},
		{ID: "user3", Username: "user3"},
	}

	tests := []struct {
		name           string
		userID         string
		mockUsers      []domain.User
		mockError      error
		expectedStatus int
		expectedBody   []map[string]interface{}
		expectedError  string
	}{
		{
			name:           "successful get following",
			userID:         "user1",
			mockUsers:      mockUsers,
			mockError:      nil,
			expectedStatus: fiber.StatusOK,
			expectedBody: []map[string]interface{}{
				{"id": "user2", "username": "user2"},
				{"id": "user3", "username": "user3"},
			},
		},
		{
			name:           "missing user ID",
			userID:         "",
			mockUsers:      nil,
			mockError:      nil,
			expectedStatus: fiber.StatusUnauthorized,
			expectedError:  "User ID is required",
		},
		{
			name:           "usecase error",
			userID:         "user1",
			mockUsers:      nil,
			mockError:      errors.New("database error"),
			expectedStatus: fiber.StatusInternalServerError,
			expectedError:  "Failed to get following list",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockUsecase, handler := setupTest()
			app.Get("/following", handler.GetFollowing)

			if tt.userID != "" {
				mockUsecase.On("GetFollowing", tt.userID).Return(tt.mockUsers, tt.mockError)
			}

			req := httptest.NewRequest("GET", "/following", nil)
			if tt.userID != "" {
				req.Header.Set("X-User-ID", tt.userID)
			}

			resp, _ := app.Test(req)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var body map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&body)

			if tt.expectedStatus == fiber.StatusOK {
				following, ok := body["following"].([]interface{})
				assert.True(t, ok)
				assert.Len(t, following, len(tt.expectedBody))
				for i, u := range following {
					userMap, ok := u.(map[string]interface{})
					assert.True(t, ok)
					assert.Equal(t, tt.expectedBody[i]["id"], userMap["id"])
					assert.Equal(t, tt.expectedBody[i]["username"], userMap["username"])
				}
			} else if tt.expectedError != "" {
				errMsg, ok := body["error"].(string)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedError, errMsg)
			}

			mockUsecase.AssertExpectations(t)
		})
	}
}

func TestUserHandler_GetFollowers(t *testing.T) {
	mockUsers := []domain.User{
		{ID: "user2", Username: "user2"},
		{ID: "user3", Username: "user3"},
	}

	tests := []struct {
		name           string
		userID         string
		mockUsers      []domain.User
		mockError      error
		expectedStatus int
		expectedBody   []map[string]interface{}
		expectedError  string
	}{
		{
			name:           "successful get followers",
			userID:         "user1",
			mockUsers:      mockUsers,
			mockError:      nil,
			expectedStatus: fiber.StatusOK,
			expectedBody: []map[string]interface{}{
				{"id": "user2", "username": "user2"},
				{"id": "user3", "username": "user3"},
			},
		},
		{
			name:           "missing user ID",
			userID:         "",
			mockUsers:      nil,
			mockError:      nil,
			expectedStatus: fiber.StatusUnauthorized,
			expectedError:  "User ID is required",
		},
		{
			name:           "usecase error",
			userID:         "user1",
			mockUsers:      nil,
			mockError:      errors.New("database error"),
			expectedStatus: fiber.StatusInternalServerError,
			expectedError:  "Failed to get followers list",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, mockUsecase, handler := setupTest()
			app.Get("/followers", handler.GetFollowers)

			if tt.userID != "" {
				mockUsecase.On("GetFollowers", tt.userID).Return(tt.mockUsers, tt.mockError)
			}

			req := httptest.NewRequest("GET", "/followers", nil)
			if tt.userID != "" {
				req.Header.Set("X-User-ID", tt.userID)
			}

			resp, _ := app.Test(req)

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			var body map[string]interface{}
			json.NewDecoder(resp.Body).Decode(&body)

			if tt.expectedStatus == fiber.StatusOK {
				followers, ok := body["followers"].([]interface{})
				assert.True(t, ok)
				assert.Len(t, followers, len(tt.expectedBody))
				for i, u := range followers {
					userMap, ok := u.(map[string]interface{})
					assert.True(t, ok)
					assert.Equal(t, tt.expectedBody[i]["id"], userMap["id"])
					assert.Equal(t, tt.expectedBody[i]["username"], userMap["username"])
				}
			} else if tt.expectedError != "" {
				errMsg, ok := body["error"].(string)
				assert.True(t, ok)
				assert.Equal(t, tt.expectedError, errMsg)
			}

			mockUsecase.AssertExpectations(t)
		})
	}
} 