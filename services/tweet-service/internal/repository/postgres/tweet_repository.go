package postgres

import (
	"github.com/lisandro/challenge/services/tweet-service/internal/domain"
	"gorm.io/gorm"
)

type tweetRepository struct {
	db *gorm.DB
}

// NewTweetRepository creates a new instance of tweet repository
func NewTweetRepository(db *gorm.DB) domain.TweetRepository {
	return &tweetRepository{db: db}
}

func (r *tweetRepository) Create(tweet *domain.Tweet) error {
	return r.db.Create(tweet).Error
}
