package usecase

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/lisandro/challenge/services/tweet-service/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTweetRepository is a mock implementation of domain.TweetRepository
type MockTweetRepository struct {
	mock.Mock
}

func (m *MockTweetRepository) Create(tweet *domain.Tweet) error {
	args := m.Called(tweet)
	return args.Error(0)
}

// MockSearchRepository is a mock implementation of domain.SearchRepository
type MockSearchRepository struct {
	mock.Mock
}

func (m *MockSearchRepository) GetTweetsByUsersID(userIDs []uuid.UUID, page, pageSize int) ([]domain.Tweet, error) {
	args := m.Called(userIDs, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Tweet), args.Error(1)
}

func (m *MockSearchRepository) IndexTweet(tweet *domain.Tweet) error {
	args := m.Called(tweet)
	return args.Error(0)
}

func TestCreateTweet(t *testing.T) {
	// Setup
	mockRepo := new(MockTweetRepository)
	mockSearchRepo := new(MockSearchRepository)
	usecase := NewTweetUseCase(mockRepo, mockSearchRepo)

	userID := uuid.New()
	content := "Test tweet content"

	// Expectations
	mockRepo.On("Create", mock.AnythingOfType("*domain.Tweet")).Return(nil)
	mockSearchRepo.On("IndexTweet", mock.AnythingOfType("*domain.Tweet")).Return(nil)

	// Execute
	tweet, err := usecase.CreateTweet(userID, content)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, tweet)
	assert.Equal(t, userID, tweet.UserID)
	assert.Equal(t, content, tweet.Content)
	assert.NotEmpty(t, tweet.ID)
	assert.WithinDuration(t, time.Now(), tweet.CreatedAt, time.Second)
	assert.WithinDuration(t, time.Now(), tweet.UpdatedAt, time.Second)

	mockRepo.AssertExpectations(t)
	mockSearchRepo.AssertExpectations(t)
}

func TestCreateTweet_RepositoryError(t *testing.T) {
	// Setup
	mockRepo := new(MockTweetRepository)
	mockSearchRepo := new(MockSearchRepository)
	usecase := NewTweetUseCase(mockRepo, mockSearchRepo)

	userID := uuid.New()
	content := "Test tweet content"

	// Expectations
	mockRepo.On("Create", mock.AnythingOfType("*domain.Tweet")).Return(assert.AnError)

	// Execute
	tweet, err := usecase.CreateTweet(userID, content)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, tweet)
	assert.Equal(t, assert.AnError, err)

	mockRepo.AssertExpectations(t)
	mockSearchRepo.AssertNotCalled(t, "IndexTweet")
}

func TestGetTweetsByUsersID(t *testing.T) {
	// Setup
	mockRepo := new(MockTweetRepository)
	mockSearchRepo := new(MockSearchRepository)
	usecase := NewTweetUseCase(mockRepo, mockSearchRepo)

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
	mockSearchRepo.On("GetTweetsByUsersID", userIDs, page, pageSize).Return(expectedTweets, nil)

	// Execute
	tweets, err := usecase.GetTweetsByUsersID(userIDs, page, pageSize)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedTweets, tweets)

	mockSearchRepo.AssertExpectations(t)
}

func TestGetTweetsByUsersID_InvalidPage(t *testing.T) {
	// Setup
	mockRepo := new(MockTweetRepository)
	mockSearchRepo := new(MockSearchRepository)
	usecase := NewTweetUseCase(mockRepo, mockSearchRepo)

	userIDs := []uuid.UUID{uuid.New()}
	page := 0
	pageSize := 10

	expectedTweets := []domain.Tweet{}

	// Expectations
	mockSearchRepo.On("GetTweetsByUsersID", userIDs, 1, pageSize).Return(expectedTweets, nil)

	// Execute
	tweets, err := usecase.GetTweetsByUsersID(userIDs, page, pageSize)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedTweets, tweets)

	mockSearchRepo.AssertExpectations(t)
}

func TestGetTweetsByUsersID_InvalidPageSize(t *testing.T) {
	// Setup
	mockRepo := new(MockTweetRepository)
	mockSearchRepo := new(MockSearchRepository)
	usecase := NewTweetUseCase(mockRepo, mockSearchRepo)

	userIDs := []uuid.UUID{uuid.New()}
	page := 1
	pageSize := 0

	expectedTweets := []domain.Tweet{}

	// Expectations
	mockSearchRepo.On("GetTweetsByUsersID", userIDs, page, 10).Return(expectedTweets, nil)

	// Execute
	tweets, err := usecase.GetTweetsByUsersID(userIDs, page, pageSize)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedTweets, tweets)

	mockSearchRepo.AssertExpectations(t)
}

func TestGetTweetsByUsersID_SearchError(t *testing.T) {
	// Setup
	mockRepo := new(MockTweetRepository)
	mockSearchRepo := new(MockSearchRepository)
	usecase := NewTweetUseCase(mockRepo, mockSearchRepo)

	userIDs := []uuid.UUID{uuid.New()}
	page := 1
	pageSize := 10

	// Expectations
	mockSearchRepo.On("GetTweetsByUsersID", userIDs, page, pageSize).Return(nil, assert.AnError)

	// Execute
	tweets, err := usecase.GetTweetsByUsersID(userIDs, page, pageSize)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, tweets)
	assert.Equal(t, assert.AnError, err)

	mockSearchRepo.AssertExpectations(t)
} 