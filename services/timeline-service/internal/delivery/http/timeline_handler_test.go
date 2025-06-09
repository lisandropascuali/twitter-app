package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lisandro/timeline-service/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTimelineUseCase is a mock implementation of usecase.TimelineUseCase
type MockTimelineUseCase struct {
	mock.Mock
}

func (m *MockTimelineUseCase) GetTimeline(ctx context.Context, userID string) (*domain.Timeline, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(*domain.Timeline), args.Error(1)
}

func setupTest() (*gin.Engine, *MockTimelineUseCase, *TimelineHandler) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockUseCase := new(MockTimelineUseCase)
	handler := NewTimelineHandler(mockUseCase)
	return router, mockUseCase, handler
}

func TestTimelineHandler_GetTimeline(t *testing.T) {
	mockTweets := []domain.Tweet{
		{
			ID:        "tweet1",
			UserID:    "user2",
			Content:   "Hello world!",
			CreatedAt: time.Now(),
		},
		{
			ID:        "tweet2",
			UserID:    "user3",
			Content:   "Another tweet",
			CreatedAt: time.Now(),
		},
	}

	mockTimeline := &domain.Timeline{
		Tweets: mockTweets,
	}

	tests := []struct {
		name               string
		userID             string
		mockTimeline       *domain.Timeline
		mockError          error
		expectedStatus     int
		expectedTweetCount int
		expectedError      string
	}{
		{
			name:               "successful get timeline",
			userID:             "user1",
			mockTimeline:       mockTimeline,
			mockError:          nil,
			expectedStatus:     http.StatusOK,
			expectedTweetCount: 2,
			expectedError:      "",
		},
		{
			name:               "missing user ID header",
			userID:             "",
			mockTimeline:       nil,
			mockError:          nil,
			expectedStatus:     http.StatusBadRequest,
			expectedTweetCount: 0,
			expectedError:      "X-User-ID header is required",
		},
		{
			name:               "usecase error",
			userID:             "user1",
			mockTimeline:       nil,
			mockError:          errors.New("failed to get following users"),
			expectedStatus:     http.StatusInternalServerError,
			expectedTweetCount: 0,
			expectedError:      "Failed to get timeline",
		},
		{
			name:               "empty timeline",
			userID:             "user1",
			mockTimeline:       &domain.Timeline{Tweets: []domain.Tweet{}},
			mockError:          nil,
			expectedStatus:     http.StatusOK,
			expectedTweetCount: 0,
			expectedError:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router, mockUseCase, handler := setupTest()
			router.GET("/timeline", handler.GetTimeline)

			// Set up mock expectations only if userID is provided
			if tt.userID != "" {
				mockUseCase.On("GetTimeline", mock.Anything, tt.userID).Return(tt.mockTimeline, tt.mockError)
			}

			// Create request
			req := httptest.NewRequest("GET", "/timeline", nil)
			if tt.userID != "" {
				req.Header.Set("X-User-ID", tt.userID)
			}

			// Perform request
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assert status code
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Parse response body
			var response map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			assert.NoError(t, err)

			if tt.expectedStatus == http.StatusOK {
				// Check timeline response
				timeline, exists := response["tweets"]
				assert.True(t, exists)
				
				tweets, ok := timeline.([]interface{})
				assert.True(t, ok)
				assert.Len(t, tweets, tt.expectedTweetCount)

				// Verify tweet structure for non-empty timelines
				if tt.expectedTweetCount > 0 {
					firstTweet, ok := tweets[0].(map[string]interface{})
					assert.True(t, ok)
					assert.Contains(t, firstTweet, "id")
					assert.Contains(t, firstTweet, "user_id")
					assert.Contains(t, firstTweet, "content")
					assert.Contains(t, firstTweet, "created_at")
				}
			} else {
				// Check error response
				errorMsg, exists := response["error"]
				assert.True(t, exists)
				assert.Equal(t, tt.expectedError, errorMsg)
			}

			// Assert mock expectations
			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestTimelineHandler_GetTimeline_ContextPropagation(t *testing.T) {
	router, mockUseCase, handler := setupTest()
	router.GET("/timeline", handler.GetTimeline)

	userID := "user1"
	mockTimeline := &domain.Timeline{Tweets: []domain.Tweet{}}

	// Set up mock to capture the context
	mockUseCase.On("GetTimeline", mock.Anything, userID).
		Run(func(args mock.Arguments) {
			ctx := args.Get(0).(context.Context)
			// Verify that the context is properly passed
			assert.NotNil(t, ctx)
		}).
		Return(mockTimeline, nil)

	req := httptest.NewRequest("GET", "/timeline", nil)
	req.Header.Set("X-User-ID", userID)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUseCase.AssertExpectations(t)
}

func TestTimelineHandler_GetTimeline_ErrorResponse(t *testing.T) {
	router, mockUseCase, handler := setupTest()
	router.GET("/timeline", handler.GetTimeline)

	userID := "user1"
	expectedError := errors.New("external service unavailable")

	mockUseCase.On("GetTimeline", mock.Anything, userID).
		Return((*domain.Timeline)(nil), expectedError)

	req := httptest.NewRequest("GET", "/timeline", nil)
	req.Header.Set("X-User-ID", userID)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response ErrorResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Failed to get timeline", response.Error)

	mockUseCase.AssertExpectations(t)
} 