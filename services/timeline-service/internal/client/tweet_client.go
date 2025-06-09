package client

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/lisandro/timeline-service/internal/domain"
)

type TweetClient interface {
	GetUserTweets(ctx context.Context, userIDs []string) ([]domain.Tweet, error)
}

type tweetClient struct {
	baseURL string
	client  *resty.Client
}

func NewTweetClient(baseURL string) TweetClient {
	return &tweetClient{
		baseURL: baseURL,
		client:  resty.New(),
	}
}

func (c *tweetClient) GetUserTweets(ctx context.Context, userIDs []string) ([]domain.Tweet, error) {
	log.Printf("Requesting tweets from tweet service")
	log.Printf("Making request to: %s/tweets/following", c.baseURL)
	
	var tweets []domain.Tweet
	resp, err := c.client.R().
		SetContext(ctx).
		SetResult(&tweets).
		Get(fmt.Sprintf("%s/tweets/following?user_ids=%s", c.baseURL, strings.Join(userIDs, ",")))

	if err != nil {
		log.Printf("Failed to get tweets from tweet service: %v", err)
		return nil, fmt.Errorf("failed to get user tweets: %w", err)
	}

	if resp.StatusCode() != 200 {
		log.Printf("Tweet service returned non-200 status: %d", resp.StatusCode())
		return nil, fmt.Errorf("failed to get user tweets: status code %d", resp.StatusCode())
	}

	log.Printf("Successfully retrieved %d tweets from tweet service", len(tweets))
	return tweets, nil
} 