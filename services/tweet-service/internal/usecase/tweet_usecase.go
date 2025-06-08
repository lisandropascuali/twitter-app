package usecase

import (
	"time"

	"github.com/google/uuid"
	"github.com/lisandro/challenge/services/tweet-service/internal/domain"
)

// tweetUsecase implements domain.TweetUseCase
type tweetUsecase struct {
	repo      domain.TweetRepository
}

// NewTweetUseCase creates a new tweet usecase instance
func NewTweetUseCase(repo domain.TweetRepository) domain.TweetUseCase {
	return &tweetUsecase{
		repo:      repo,
	}
}

// CreateTweet creates a new tweet for a user
func (u *tweetUsecase) CreateTweet(userID uuid.UUID, content string) (*domain.Tweet, error) {
	tweet := &domain.Tweet{
		ID:        uuid.New(),
		UserID:    userID,
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := u.repo.Create(tweet); err != nil {
		return nil, err
	}

	// TODO: Publish tweet event to SNS
	// This will be implemented when we set up the AWS SNS client

	return tweet, nil
}
