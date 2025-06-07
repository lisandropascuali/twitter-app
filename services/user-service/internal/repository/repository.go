package repository

import "github.com/lisandro/challenge/services/user-service/internal/domain"

// UserRepository represents the user's repository contract
type UserRepository interface {
    Follow(followerID, followedID string) error
    Unfollow(followerID, followedID string) error
    GetFollowing(userID string) ([]string, error)
    GetUser(id string) (*domain.User, error)
} 