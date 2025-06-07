package usecase

import (
	"errors"
	"testing"

	"github.com/lisandro/challenge/services/user-service/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of domain.UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Follow(followerID, followedID string) error {
	args := m.Called(followerID, followedID)
	return args.Error(0)
}

func (m *MockUserRepository) Unfollow(followerID, followedID string) error {
	args := m.Called(followerID, followedID)
	return args.Error(0)
}

func (m *MockUserRepository) GetFollowing(userID string) ([]domain.User, error) {
	args := m.Called(userID)
	return args.Get(0).([]domain.User), args.Error(1)
}

func (m *MockUserRepository) GetFollowers(userID string) ([]domain.User, error) {
	args := m.Called(userID)
	return args.Get(0).([]domain.User), args.Error(1)
}

func TestUserUsecase_Follow(t *testing.T) {
	tests := []struct {
		name        string
		followerID  string
		followedID  string
		mockError   error
		expectError bool
	}{
		{
			name:        "successful follow",
			followerID:  "user1",
			followedID:  "user2",
			mockError:   nil,
			expectError: false,
		},
		{
			name:        "repository error",
			followerID:  "user1",
			followedID:  "user2",
			mockError:   errors.New("database error"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			mockRepo.On("Follow", tt.followerID, tt.followedID).Return(tt.mockError)

			usecase := NewUserUsecase(mockRepo)
			err := usecase.Follow(tt.followerID, tt.followedID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserUsecase_Unfollow(t *testing.T) {
	tests := []struct {
		name        string
		followerID  string
		followedID  string
		mockError   error
		expectError bool
	}{
		{
			name:        "successful unfollow",
			followerID:  "user1",
			followedID:  "user2",
			mockError:   nil,
			expectError: false,
		},
		{
			name:        "repository error",
			followerID:  "user1",
			followedID:  "user2",
			mockError:   errors.New("database error"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			mockRepo.On("Unfollow", tt.followerID, tt.followedID).Return(tt.mockError)

			usecase := NewUserUsecase(mockRepo)
			err := usecase.Unfollow(tt.followerID, tt.followedID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserUsecase_GetFollowing(t *testing.T) {
	mockUsers := []domain.User{
		{ID: "user2", Username: "user2"},
		{ID: "user3", Username: "user3"},
	}

	tests := []struct {
		name        string
		userID      string
		mockUsers   []domain.User
		mockError   error
		expectError bool
	}{
		{
			name:        "successful get following",
			userID:      "user1",
			mockUsers:   mockUsers,
			mockError:   nil,
			expectError: false,
		},
		{
			name:        "repository error",
			userID:      "user1",
			mockUsers:   nil,
			mockError:   errors.New("database error"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			mockRepo.On("GetFollowing", tt.userID).Return(tt.mockUsers, tt.mockError)

			usecase := NewUserUsecase(mockRepo)
			users, err := usecase.GetFollowing(tt.userID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, users)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockUsers, users)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserUsecase_GetFollowers(t *testing.T) {
	mockUsers := []domain.User{
		{ID: "user2", Username: "user2"},
		{ID: "user3", Username: "user3"},
	}

	tests := []struct {
		name        string
		userID      string
		mockUsers   []domain.User
		mockError   error
		expectError bool
	}{
		{
			name:        "successful get followers",
			userID:      "user1",
			mockUsers:   mockUsers,
			mockError:   nil,
			expectError: false,
		},
		{
			name:        "repository error",
			userID:      "user1",
			mockUsers:   nil,
			mockError:   errors.New("database error"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			mockRepo.On("GetFollowers", tt.userID).Return(tt.mockUsers, tt.mockError)

			usecase := NewUserUsecase(mockRepo)
			users, err := usecase.GetFollowers(tt.userID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, users)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.mockUsers, users)
			}
			mockRepo.AssertExpectations(t)
		})
	}
} 