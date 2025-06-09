package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/lisandro/timeline-service/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserClient is a mock implementation of client.UserClient
type MockUserClient struct {
	mock.Mock
}

func (m *MockUserClient) GetFollowingUsers(ctx context.Context, userID string) ([]domain.FollowingUser, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]domain.FollowingUser), args.Error(1)
}

// MockTweetClient is a mock implementation of client.TweetClient
type MockTweetClient struct {
	mock.Mock
}

func (m *MockTweetClient) GetUserTweets(ctx context.Context, userIDs []string) ([]domain.Tweet, error) {
	args := m.Called(ctx, userIDs)
	return args.Get(0).([]domain.Tweet), args.Error(1)
}

func TestTimelineUseCase_GetTimeline(t *testing.T) {
	mockFollowingUsers := []domain.FollowingUser{
		{ID: "user2", Username: "alice"},
		{ID: "user3", Username: "bob"},
	}

	mockTweets := []domain.Tweet{
		{
			ID:        "tweet1",
			UserID:    "user2",
			Content:   "Hello from Alice!",
			CreatedAt: time.Now(),
		},
		{
			ID:        "tweet2",
			UserID:    "user2",
			Content:   "Another tweet from Alice",
			CreatedAt: time.Now(),
		},
		{
			ID:        "tweet3",
			UserID:    "user3",
			Content:   "Hello from Bob!",
			CreatedAt: time.Now(),
		},
	}

	tests := []struct {
		name                 string
		userID               string
		mockFollowingUsers   []domain.FollowingUser
		mockFollowingError   error
		mockTweets           []domain.Tweet
		mockTweetsError      error
		expectedTweetCount   int
		expectedError        bool
		expectedErrorMessage string
	}{
		{
			name:               "successful timeline generation",
			userID:             "user1",
			mockFollowingUsers: mockFollowingUsers,
			mockFollowingError: nil,
			mockTweets:         mockTweets,
			mockTweetsError:    nil,
			expectedTweetCount: 3,
			expectedError:      false,
			expectedErrorMessage: "",
		},
		{
			name:                 "error getting following users",
			userID:               "user1",
			mockFollowingUsers:   nil,
			mockFollowingError:   errors.New("user service unavailable"),
			mockTweets:           nil,
			mockTweetsError:      nil,
			expectedTweetCount:   0,
			expectedError:        true,
			expectedErrorMessage: "user service unavailable",
		},
		{
			name:               "error getting tweets",
			userID:             "user1",
			mockFollowingUsers: mockFollowingUsers,
			mockFollowingError: nil,
			mockTweets:         nil,
			mockTweetsError:    errors.New("tweet service unavailable"),
			expectedTweetCount: 0,
			expectedError:      true,
			expectedErrorMessage: "tweet service unavailable",
		},
		{
			name:               "empty following list",
			userID:             "user1",
			mockFollowingUsers: []domain.FollowingUser{},
			mockFollowingError: nil,
			mockTweets:         []domain.Tweet{},
			mockTweetsError:    nil,
			expectedTweetCount: 0,
			expectedError:      false,
			expectedErrorMessage: "",
		},
		{
			name:               "following users with no tweets",
			userID:             "user1",
			mockFollowingUsers: mockFollowingUsers,
			mockFollowingError: nil,
			mockTweets:         []domain.Tweet{},
			mockTweetsError:    nil,
			expectedTweetCount: 0,
			expectedError:      false,
			expectedErrorMessage: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserClient := new(MockUserClient)
			mockTweetClient := new(MockTweetClient)

			// Set up user client mock
			mockUserClient.On("GetFollowingUsers", mock.Anything, tt.userID).
				Return(tt.mockFollowingUsers, tt.mockFollowingError)

			// Set up tweet client mock only if we have following users and no following error
			if tt.mockFollowingError == nil && len(tt.mockFollowingUsers) > 0 {
				mockTweetClient.On("GetUserTweets", mock.Anything, mock.Anything).
					Return(tt.mockTweets, tt.mockTweetsError)
			}

			useCase := NewTimelineUseCase(mockUserClient, mockTweetClient)
			timeline, err := useCase.GetTimeline(context.Background(), tt.userID)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, timeline)
				if tt.expectedErrorMessage != "" {
					assert.Contains(t, err.Error(), tt.expectedErrorMessage)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, timeline)
				assert.Len(t, timeline.Tweets, tt.expectedTweetCount)

				// Verify that all tweets are from following users
				if tt.expectedTweetCount > 0 {
					followingUserIDs := make(map[string]bool)
					for _, user := range tt.mockFollowingUsers {
						followingUserIDs[user.ID] = true
					}

					for _, tweet := range timeline.Tweets {
						assert.True(t, followingUserIDs[tweet.UserID], 
							"Tweet from user %s should be from a following user", tweet.UserID)
					}
				}
			}

			mockUserClient.AssertExpectations(t)
			mockTweetClient.AssertExpectations(t)
		})
	}
}

func TestTimelineUseCase_GetTimeline_ContextPropagation(t *testing.T) {
	mockUserClient := new(MockUserClient)
	mockTweetClient := new(MockTweetClient)

	userID := "user1"
	followingUsers := []domain.FollowingUser{
		{ID: "user2", Username: "alice"},
	}

	tweets := []domain.Tweet{
		{ID: "tweet1", UserID: "user2", Content: "Test tweet", CreatedAt: time.Now()},
	}

	ctx := context.WithValue(context.Background(), "test_key", "test_value")

	// Verify context is passed through to clients
	mockUserClient.On("GetFollowingUsers", ctx, userID).
		Run(func(args mock.Arguments) {
			receivedCtx := args.Get(0).(context.Context)
			assert.Equal(t, "test_value", receivedCtx.Value("test_key"))
		}).
		Return(followingUsers, nil)

	mockTweetClient.On("GetUserTweets", ctx, []string{"user2"}).
		Run(func(args mock.Arguments) {
			receivedCtx := args.Get(0).(context.Context)
			assert.Equal(t, "test_value", receivedCtx.Value("test_key"))
		}).
		Return(tweets, nil)

	useCase := NewTimelineUseCase(mockUserClient, mockTweetClient)
	timeline, err := useCase.GetTimeline(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, timeline)
	assert.Len(t, timeline.Tweets, 1)

	mockUserClient.AssertExpectations(t)
	mockTweetClient.AssertExpectations(t)
}

func TestTimelineUseCase_GetTimeline_UserIDsCollection(t *testing.T) {
	mockUserClient := new(MockUserClient)
	mockTweetClient := new(MockTweetClient)

	userID := "user1"
	followingUsers := []domain.FollowingUser{
		{ID: "user2", Username: "alice"},
		{ID: "user3", Username: "bob"},
		{ID: "user4", Username: "charlie"},
	}

	expectedUserIDs := []string{"user2", "user3", "user4"}
	tweets := []domain.Tweet{
		{ID: "tweet1", UserID: "user2", Content: "Tweet from Alice", CreatedAt: time.Now()},
		{ID: "tweet2", UserID: "user3", Content: "Tweet from Bob", CreatedAt: time.Now()},
		{ID: "tweet3", UserID: "user4", Content: "Tweet from Charlie", CreatedAt: time.Now()},
	}

	mockUserClient.On("GetFollowingUsers", mock.Anything, userID).
		Return(followingUsers, nil)

	// Verify that the exact user IDs are passed to the tweet client
	mockTweetClient.On("GetUserTweets", mock.Anything, expectedUserIDs).
		Run(func(args mock.Arguments) {
			receivedUserIDs := args.Get(1).([]string)
			assert.ElementsMatch(t, expectedUserIDs, receivedUserIDs)
		}).
		Return(tweets, nil)

	useCase := NewTimelineUseCase(mockUserClient, mockTweetClient)
	timeline, err := useCase.GetTimeline(context.Background(), userID)

	assert.NoError(t, err)
	assert.NotNil(t, timeline)
	assert.Len(t, timeline.Tweets, 3)

	mockUserClient.AssertExpectations(t)
	mockTweetClient.AssertExpectations(t)
} 