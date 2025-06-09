package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/lisandro/challenge/services/tweet-service/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTweetUseCase is a mock implementation of domain.TweetUseCase
type MockTweetUseCase struct {
	mock.Mock
}

func (m *MockTweetUseCase) CreateTweet(userID uuid.UUID, content string) (*domain.Tweet, error) {
	args := m.Called(userID, content)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Tweet), args.Error(1)
}

func (m *MockTweetUseCase) GetTweetsByUsersID(userIDs []uuid.UUID, page, pageSize int) ([]domain.Tweet, error) {
	args := m.Called(userIDs, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Tweet), args.Error(1)
}

func setupTest() (*fiber.App, *MockTweetUseCase) {
	app := fiber.New()
	mockUseCase := new(MockTweetUseCase)
	handler := NewHandler(mockUseCase)

	// Register routes
	RegisterRoutes(app, handler)

	return app, mockUseCase
}

func TestCreateTweet(t *testing.T) {
	// Setup
	app, mockUseCase := setupTest()

	userID := uuid.New()
	content := "Test tweet content"
	now := time.Now()

	expectedTweet := &domain.Tweet{
		ID:        uuid.New(),
		UserID:    userID,
		Content:   content,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Expectations
	mockUseCase.On("CreateTweet", userID, content).Return(expectedTweet, nil)

	// Create request
	reqBody := CreateTweetRequest{
		Content: content,
	}
	reqBodyJSON, _ := json.Marshal(reqBody)

	// Execute
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tweets", bytes.NewReader(reqBodyJSON))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-User-ID", userID.String())
	resp, err := app.Test(req)
	assert.NoError(t, err)

	// Assert
	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var response Tweet
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, expectedTweet.ID, response.ID)
	assert.Equal(t, expectedTweet.UserID, response.UserID)
	assert.Equal(t, expectedTweet.Content, response.Content)

	mockUseCase.AssertExpectations(t)
}

func TestCreateTweet_InvalidRequest(t *testing.T) {
	// Setup
	app, mockUseCase := setupTest()

	// Create invalid request
	reqBody := map[string]interface{}{
		"content": "Test tweet content",
	}
	reqBodyJSON, _ := json.Marshal(reqBody)

	// Execute
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tweets", bytes.NewReader(reqBodyJSON))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	assert.NoError(t, err)

	// Assert
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response ErrorResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.Error)

	// Verify that CreateTweet was not called
	mockUseCase.AssertNotCalled(t, "CreateTweet")
}

func TestGetTweetsByUsersID(t *testing.T) {
	// Setup
	app, mockUseCase := setupTest()

	userIDs := []uuid.UUID{uuid.New(), uuid.New()}
	page := 1
	pageSize := 10

	expectedTweets := []domain.Tweet{
		{
			ID:        uuid.New(),
			UserID:    userIDs[0],
			Content:   "Tweet 1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ID:        uuid.New(),
			UserID:    userIDs[1],
			Content:   "Tweet 2",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Expectations
	mockUseCase.On("GetTweetsByUsersID", userIDs, page, pageSize).Return(expectedTweets, nil)

	// Execute
	req := httptest.NewRequest(http.MethodGet, "/api/v1/tweets/following?user_ids="+userIDs[0].String()+","+userIDs[1].String()+"&page=1&page_size=10", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)

	// Assert
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var response []Tweet
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, expectedTweets[0].ID, response[0].ID)
	assert.Equal(t, expectedTweets[1].ID, response[1].ID)

	mockUseCase.AssertExpectations(t)
}

func TestGetTweetsByUsersID_MissingUserIDs(t *testing.T) {
	// Setup
	app, mockUseCase := setupTest()

	// Execute
	req := httptest.NewRequest(http.MethodGet, "/api/v1/tweets/following", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)

	// Assert
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response ErrorResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.Error)

	mockUseCase.AssertNotCalled(t, "GetTweetsByUsersID")
}

func TestGetTweetsByUsersID_InvalidUserID(t *testing.T) {
	// Setup
	app, mockUseCase := setupTest()

	// Execute
	req := httptest.NewRequest(http.MethodGet, "/api/v1/tweets/following?user_ids=invalid-uuid", nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)

	// Assert
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response ErrorResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.Error)

	mockUseCase.AssertNotCalled(t, "GetTweetsByUsersID")
}

func TestGetTweetsByUsersID_UseCaseError(t *testing.T) {
	// Setup
	app, mockUseCase := setupTest()

	userID := uuid.New()

	// Expectations
	mockUseCase.On("GetTweetsByUsersID", []uuid.UUID{userID}, 1, 10).Return(nil, assert.AnError)

	// Execute
	req := httptest.NewRequest(http.MethodGet, "/api/v1/tweets/following?user_ids="+userID.String(), nil)

	resp, err := app.Test(req)
	assert.NoError(t, err)

	// Assert
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response ErrorResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response.Error)

	mockUseCase.AssertExpectations(t)
} 