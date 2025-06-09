package domain

import (
	"time"

	"github.com/google/uuid"
)

// Tweet represents a tweet in the system
type Tweet struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	UserID    uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	CreatedAt time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt time.Time `json:"updated_at" gorm:"not null"`
}

// TweetRepository defines the interface for tweet data operations
type TweetRepository interface {
	Create(tweet *Tweet) error
}

type SearchRepository interface {
	GetTweetsByUsersID(userIDs []uuid.UUID, page, pageSize int) ([]Tweet, error)
	IndexTweet(tweet *Tweet) error
}

// TweetUseCase defines the interface for tweet business logic
type TweetUseCase interface {
	CreateTweet(userID uuid.UUID, content string) (*Tweet, error)
	GetTweetsByUsersID(userIDs []uuid.UUID, page, pageSize int) ([]Tweet, error)
} 