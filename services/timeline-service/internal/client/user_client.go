package client

import (
	"context"
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
	"github.com/lisandro/timeline-service/internal/domain"
)

type UserClient interface {
	GetFollowingUsers(ctx context.Context, userID string) ([]domain.FollowingUser, error)
}

type userClient struct {
	baseURL string
	client  *resty.Client
}

// FollowingResponse represents the response structure from the user service
type FollowingResponse struct {
	Following []domain.FollowingUser `json:"following"`
}

func NewUserClient(baseURL string) UserClient {
	return &userClient{
		baseURL: baseURL,
		client:  resty.New(),
	}
}

func (c *userClient) GetFollowingUsers(ctx context.Context, userID string) ([]domain.FollowingUser, error) {
	log.Printf("Requesting following users from user service for user %s", userID)
	log.Printf("Making request to: %s/following", c.baseURL)
	
	var response FollowingResponse
	resp, err := c.client.R().
		SetContext(ctx).
		SetHeader("X-User-ID", userID).
		SetResult(&response).
		Get(fmt.Sprintf("%s/following", c.baseURL))

	if err != nil {
		log.Printf("Failed to get following users from user service for user %s: %v", userID, err)
		return nil, fmt.Errorf("failed to get following users: %w", err)
	}

	if resp.StatusCode() != 200 {
		log.Printf("User service returned non-200 status for user %s: %d", userID, resp.StatusCode())
		return nil, fmt.Errorf("failed to get following users: status code %d", resp.StatusCode())
	}

	log.Printf("Successfully retrieved %d following users from user service for user %s", len(response.Following), userID)
	return response.Following, nil
} 