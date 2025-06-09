package opensearch

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lisandro/challenge/services/tweet-service/internal/domain"
	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
)

const (
	tweetsIndex = "tweets"
)

type searchRepository struct {
	client *opensearch.Client
}

func NewSearchRepository(client *opensearch.Client) domain.SearchRepository {
	return &searchRepository{
		client: client,
	}
}

// IndexTweet indexes a tweet in OpenSearch
func (r *searchRepository) IndexTweet(tweet *domain.Tweet) error {
	doc := map[string]interface{}{
		"id":         tweet.ID.String(),
		"user_id":    tweet.UserID.String(),
		"content":    tweet.Content,
		"created_at": tweet.CreatedAt.Format(time.RFC3339),
		"updated_at": tweet.UpdatedAt.Format(time.RFC3339),
	}

	docJSON, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("failed to marshal tweet: %w", err)
	}

	req := opensearchapi.IndexRequest{
		Index:      tweetsIndex,
		DocumentID: tweet.ID.String(),
		Body:       strings.NewReader(string(docJSON)),
	}

	res, err := req.Do(context.Background(), r.client)
	if err != nil {
		return fmt.Errorf("failed to index tweet: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error indexing tweet: %s", res.String())
	}

	return nil
}

func (r *searchRepository) GetTweetsByUsersID(userIDs []uuid.UUID, page, pageSize int) ([]domain.Tweet, error) {
	from := (page - 1) * pageSize

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"terms": map[string]interface{}{
				"user_id": userIDs,
			},
		},
		"sort": []map[string]interface{}{
			{
				"created_at": map[string]interface{}{
					"order": "desc",
				},
			},
		},
		"from": from,
		"size": pageSize,
	}

	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal query: %w", err)
	}

	searchRequest := opensearchapi.SearchRequest{
		Index: []string{tweetsIndex},
		Body:  strings.NewReader(string(queryJSON)),
	}

	response, err := searchRequest.Do(context.Background(), r.client)
	if err != nil {
		return nil, fmt.Errorf("failed to search tweets: %w", err)
	}
	defer response.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	hits := result["hits"].(map[string]interface{})["hits"].([]interface{})
	tweets := make([]domain.Tweet, 0, len(hits))

	for _, hit := range hits {
		hitMap := hit.(map[string]interface{})
		source := hitMap["_source"].(map[string]interface{})

		tweetID, err := uuid.Parse(source["id"].(string))
		if err != nil {
			return nil, fmt.Errorf("failed to parse tweet ID: %w", err)
		}

		userID, err := uuid.Parse(source["user_id"].(string))
		if err != nil {
			return nil, fmt.Errorf("failed to parse user ID: %w", err)
		}

		createdAt, err := time.Parse(time.RFC3339, source["created_at"].(string))
		if err != nil {
			return nil, fmt.Errorf("failed to parse created_at: %w", err)
		}

		updatedAt, err := time.Parse(time.RFC3339, source["updated_at"].(string))
		if err != nil {
			return nil, fmt.Errorf("failed to parse updated_at: %w", err)
		}

		tweets = append(tweets, domain.Tweet{
			ID:        tweetID,
			UserID:    userID,
			Content:   source["content"].(string),
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})
	}

	return tweets, nil
}
