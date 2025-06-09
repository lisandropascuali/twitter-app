package usecase

import (
	"context"
	"log"

	"github.com/lisandro/timeline-service/internal/client"
	"github.com/lisandro/timeline-service/internal/domain"
)

type TimelineUseCase interface {
	GetTimeline(ctx context.Context, userID string) (*domain.Timeline, error)
}

type timelineUseCase struct {
	userClient  client.UserClient
	tweetClient client.TweetClient
}

func NewTimelineUseCase(userClient client.UserClient, tweetClient client.TweetClient) TimelineUseCase {
	return &timelineUseCase{
		userClient:  userClient,
		tweetClient: tweetClient,
	}
}

func (uc *timelineUseCase) GetTimeline(ctx context.Context, userID string) (*domain.Timeline, error) {
	log.Printf("Starting timeline generation for user %s", userID)
	
	// Get following users
	log.Printf("Fetching following users for user %s", userID)
	followingUsers, err := uc.userClient.GetFollowingUsers(ctx, userID)
	if err != nil {
		log.Printf("Error fetching following users for user %s: %v", userID, err)
		return nil, err
	}
	log.Printf("Found %d following users for user %s", len(followingUsers), userID)
	var followingUserIDs []string
	for _, followingUser := range followingUsers {
		followingUserIDs = append(followingUserIDs, followingUser.ID)
	}
	// Get tweets for each following user
	log.Printf("Fetching tweets for following user %s", followingUserIDs)
	tweets, err := uc.tweetClient.GetUserTweets(ctx, followingUserIDs)
	if err != nil {
		log.Printf("Error fetching tweets for user %s: %v", followingUserIDs, err)
		return nil, err
	}
	log.Printf("Tweets: %v", tweets)
	log.Printf("Found %d tweets for user %s", len(tweets), followingUserIDs)

	log.Printf("Timeline generation completed for user %s with %d total tweets", userID, len(tweets))
	return &domain.Timeline{
		Tweets: tweets,
	}, nil
} 