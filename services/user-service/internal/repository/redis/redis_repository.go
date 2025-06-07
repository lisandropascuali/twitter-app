package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/lisandro/challenge/services/user-service/internal/domain"
)

// RedisRepository handles user data caching in Redis
type RedisRepository struct {
	client *redis.Client
	ctx    context.Context
}

// NewRedisRepository creates a new Redis repository
func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{
		client: client,
		ctx:    context.Background(),
	}
}

func (r *RedisRepository) CacheUser(user *domain.User) error {
	userJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return r.client.Set(r.ctx, fmt.Sprintf("user:%s", user.ID), userJSON, 0).Err()
}

func (r *RedisRepository) GetCachedUser(id string) (*domain.User, error) {
	val, err := r.client.Get(r.ctx, fmt.Sprintf("user:%s", id)).Result()
	if err != nil {
		return nil, err
	}

	var user domain.User
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *RedisRepository) CacheFollowing(userID string, following []domain.User) error {
	log.Printf("Caching following list for user %s with %d users", userID, len(following))
	key := fmt.Sprintf("following:%s", userID)
	pipe := r.client.Pipeline()
	
	for _, user := range following {
		userJSON, err := json.Marshal(user)
		if err != nil {
			log.Printf("Error marshaling user %s for cache: %v", user.ID, err)
			return err
		}
		pipe.SAdd(r.ctx, key, userJSON)
	}
	
	_, err := pipe.Exec(r.ctx)
	if err != nil {
		log.Printf("Error executing pipeline for user %s: %v", userID, err)
	}
	return err
}

func (r *RedisRepository) GetCachedFollowing(userID string) ([]domain.User, error) {
	log.Printf("Attempting to get cached following list for user %s", userID)
	key := fmt.Sprintf("following:%s", userID)
	
	// Check if the key exists first
	exists, err := r.client.Exists(r.ctx, key).Result()
	if err != nil {
		log.Printf("Error checking if following key exists for user %s: %v", userID, err)
		return nil, err
	}
	if exists == 0 {
		log.Printf("Cache MISS: Key does not exist for user %s", userID)
		return nil, fmt.Errorf("cache miss")
	}

	vals, err := r.client.SMembers(r.ctx, key).Result()
	if err != nil {
		log.Printf("Error getting cached following for user %s: %v", userID, err)
		return nil, err
	}

	var following []domain.User
	for _, val := range vals {
		var user domain.User
		if err := json.Unmarshal([]byte(val), &user); err != nil {
			log.Printf("Error unmarshaling cached user for user %s: %v", userID, err)
			return nil, err
		}
		following = append(following, user)
	}
	log.Printf("Successfully retrieved %d users from cache for user %s", len(following), userID)
	return following, nil
}

func (r *RedisRepository) InvalidateUserCache(userID string) error {
	return r.client.Del(r.ctx, fmt.Sprintf("user:%s", userID)).Err()
}

func (r *RedisRepository) InvalidateFollowingCache(userID string) error {
	return r.client.Del(r.ctx, fmt.Sprintf("following:%s", userID)).Err()
}

func (r *RedisRepository) CacheFollowers(userID string, followers []domain.User) error {
	log.Printf("Caching followers list for user %s with %d users", userID, len(followers))
	key := fmt.Sprintf("followers:%s", userID)
	pipe := r.client.Pipeline()
	
	for _, user := range followers {
		userJSON, err := json.Marshal(user)
		if err != nil {
			log.Printf("Error marshaling user %s for cache: %v", user.ID, err)
			return err
		}
		pipe.SAdd(r.ctx, key, userJSON)
	}
	
	_, err := pipe.Exec(r.ctx)
	if err != nil {
		log.Printf("Error executing pipeline for user %s: %v", userID, err)
	}
	return err
}

func (r *RedisRepository) GetCachedFollowers(userID string) ([]domain.User, error) {
	log.Printf("Attempting to get cached followers list for user %s", userID)
	key := fmt.Sprintf("followers:%s", userID)
	
	// Check if the key exists first
	exists, err := r.client.Exists(r.ctx, key).Result()
	if err != nil {
		log.Printf("Error checking if followers key exists for user %s: %v", userID, err)
		return nil, err
	}
	if exists == 0 {
		log.Printf("Cache MISS: Key does not exist for user %s", userID)
		return nil, fmt.Errorf("cache miss")
	}

	vals, err := r.client.SMembers(r.ctx, key).Result()
	if err != nil {
		log.Printf("Error getting cached followers for user %s: %v", userID, err)
		return nil, err
	}

	var followers []domain.User
	for _, val := range vals {
		var user domain.User
		if err := json.Unmarshal([]byte(val), &user); err != nil {
			log.Printf("Error unmarshaling cached user for user %s: %v", userID, err)
			return nil, err
		}
		followers = append(followers, user)
	}
	log.Printf("Successfully retrieved %d users from cache for user %s", len(followers), userID)
	return followers, nil
}

func (r *RedisRepository) InvalidateFollowersCache(userID string) error {
	return r.client.Del(r.ctx, fmt.Sprintf("followers:%s", userID)).Err()
} 