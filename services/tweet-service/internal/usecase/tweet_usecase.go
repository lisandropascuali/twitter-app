package usecase

import (
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/lisandro/challenge/services/tweet-service/internal/domain"
)

// tweetUsecase implements domain.TweetUseCase
type tweetUsecase struct {
	repo      domain.TweetRepository
	searchRepo domain.SearchRepository
}

// NewTweetUseCase creates a new tweet usecase instance
func NewTweetUseCase(repo domain.TweetRepository, searchRepo domain.SearchRepository) domain.TweetUseCase {
	return &tweetUsecase{
		repo:      repo,
		searchRepo: searchRepo,
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

	// Index the tweet in OpenSearch
	if err := u.searchRepo.IndexTweet(tweet); err != nil {
		// Log the error but don't fail the request
		// In a production environment, you might want to handle this differently
		// For example, you could use a message queue to retry the indexing
		log.Printf("Failed to index tweet in OpenSearch: %v", err)
	}

	// TODO: Publish tweet event to SNS
	// This will be implemented when we set up the AWS SNS client

	return tweet, nil
}

// GetTweetsByUsersID retrieves tweets from a list of user IDs
func (u *tweetUsecase) GetTweetsByUsersID(userIDs []uuid.UUID, page, pageSize int) ([]domain.Tweet, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	
	return u.searchRepo.GetTweetsByUsersID(userIDs, page, pageSize)
}
