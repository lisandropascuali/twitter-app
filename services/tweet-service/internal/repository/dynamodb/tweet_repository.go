package dynamodb

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/lisandro/challenge/services/tweet-service/internal/domain"
)

type tweetRepository struct {
	client    *dynamodb.Client
	tableName string
}

// NewTweetRepository creates a new instance of tweet repository
func NewTweetRepository(client *dynamodb.Client, tableName string) domain.TweetRepository {
	return &tweetRepository{
		client:    client,
		tableName: tableName,
	}
}

func (r *tweetRepository) Create(tweet *domain.Tweet) error {
	now := time.Now()
	tweet.CreatedAt = now
	tweet.UpdatedAt = now

	item := map[string]types.AttributeValue{
		"id": &types.AttributeValueMemberS{
			Value: tweet.ID.String(),
		},
		"user_id": &types.AttributeValueMemberS{
			Value: tweet.UserID.String(),
		},
		"content": &types.AttributeValueMemberS{
			Value: tweet.Content,
		},
		"created_at": &types.AttributeValueMemberS{
			Value: tweet.CreatedAt.Format(time.RFC3339),
		},
		"updated_at": &types.AttributeValueMemberS{
			Value: tweet.UpdatedAt.Format(time.RFC3339),
		},
	}

	_, err := r.client.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String(r.tableName),
		Item:      item,
	})

	return err
} 