package usecase

import (
	"context"

	"github.com/lisandro/timeline-service/internal/domain"
)

type Timeline interface {
	GetFollowingUsers(ctx context.Context, userID string) ([]domain.FollowingUser, error)
	GetTweetsByUserIDs(ctx context.Context, userIDs []string) ([]domain.Tweet, error)
} 