package repository

import (
	"log"

	"github.com/lisandro/challenge/services/user-service/internal/domain"
)

// PersistentRepository defines the interface for persistent storage (e.g., PostgreSQL)
type PersistentRepository interface {
	Follow(followerID, followedID string) error
	Unfollow(followerID, followedID string) error
	GetFollowing(userID string) ([]domain.User, error)
	GetFollowers(userID string) ([]domain.User, error)
	GetUser(id string) (*domain.User, error)
}

// CacheRepository defines the interface for caching storage (e.g., Redis)
type CacheRepository interface {
	GetCachedUser(id string) (*domain.User, error)
	CacheUser(user *domain.User) error
	InvalidateUserCache(userID string) error
	GetCachedFollowing(userID string) ([]domain.User, error)
	CacheFollowing(userID string, following []domain.User) error
	InvalidateFollowingCache(userID string) error
	GetCachedFollowers(userID string) ([]domain.User, error)
	CacheFollowers(userID string, followers []domain.User) error
	InvalidateFollowersCache(userID string) error
}

type compositeRepository struct {
	persistent PersistentRepository
	cache      CacheRepository
}

// NewCompositeRepository creates a new composite repository that combines persistent and cache storage
func NewCompositeRepository(persistent PersistentRepository, cache CacheRepository) domain.UserRepository {
	return &compositeRepository{
		persistent: persistent,
		cache:      cache,
	}
}

func (r *compositeRepository) Follow(followerID, followedID string) error {
	// First update persistent storage
	if err := r.persistent.Follow(followerID, followedID); err != nil {
		return err
	}
	
	// Invalidate cache
	return r.cache.InvalidateFollowingCache(followerID)
}

func (r *compositeRepository) Unfollow(followerID, followedID string) error {
	// First update persistent storage
	if err := r.persistent.Unfollow(followerID, followedID); err != nil {
		return err
	}
	
	// Invalidate cache
	return r.cache.InvalidateFollowingCache(followerID)
}

func (r *compositeRepository) GetFollowing(userID string) ([]domain.User, error) {
	log.Printf("Getting following list for user %s", userID)
	
	// Try cache first
	following, err := r.cache.GetCachedFollowing(userID)
	if err == nil {
		log.Printf("Cache HIT: Found following list for user %s with %d users", userID, len(following))
		return following, nil
	}
	log.Printf("Cache MISS: No following list found in cache for user %s, error: %v", userID, err)

	// On cache miss, get from persistent storage
	following, err = r.persistent.GetFollowing(userID)
	if err != nil {
		log.Printf("Error getting following from persistent storage for user %s: %v", userID, err)
		return nil, err
	}
	log.Printf("Retrieved %d following users from persistent storage for user %s", len(following), userID)

	// Update cache
	if err := r.cache.CacheFollowing(userID, following); err != nil {
		log.Printf("Failed to cache following list for user %s: %v", userID, err)
	} else {
		log.Printf("Successfully cached following list for user %s", userID)
	}

	return following, nil
}

func (r *compositeRepository) GetUser(id string) (*domain.User, error) {
	// Try cache first
	user, err := r.cache.GetCachedUser(id)
	if err == nil {
		return user, nil
	}

	// On cache miss, get from persistent storage
	user, err = r.persistent.GetUser(id)
	if err != nil {
		return nil, err
	}

	// Update cache
	if err := r.cache.CacheUser(user); err != nil {
		// Log error but don't fail the request
		// log.Printf("Failed to cache user: %v", err)
	}

	return user, nil
}

func (r *compositeRepository) GetFollowers(userID string) ([]domain.User, error) {
	log.Printf("Getting followers list for user %s", userID)
	
	// Try cache first
	followers, err := r.cache.GetCachedFollowers(userID)
	if err == nil {
		log.Printf("Cache HIT: Found followers list for user %s with %d users", userID, len(followers))
		return followers, nil
	}
	log.Printf("Cache MISS: No followers list found in cache for user %s, error: %v", userID, err)

	// On cache miss, get from persistent storage
	followers, err = r.persistent.GetFollowers(userID)
	if err != nil {
		log.Printf("Error getting followers from persistent storage for user %s: %v", userID, err)
		return nil, err
	}
	log.Printf("Retrieved %d followers from persistent storage for user %s", len(followers), userID)

	// Update cache
	if err := r.cache.CacheFollowers(userID, followers); err != nil {
		log.Printf("Failed to cache followers list for user %s: %v", userID, err)
	} else {
		log.Printf("Successfully cached followers list for user %s", userID)
	}

	return followers, nil
} 