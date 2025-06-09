package repository

import "github.com/lisandro/challenge/services/user-service/internal/domain"

// UserRepository represents the user's repository contract
type UserRepository interface {
    GetAllUsers() ([]domain.User, error)
    CreateUser(req domain.CreateUserRequest) (*domain.User, error)
    Follow(followerID, followedID string) error
    Unfollow(followerID, followedID string) error
    GetFollowing(userID string) ([]domain.User, error)
    GetFollowers(userID string) ([]domain.User, error)
    GetUser(id string) (*domain.User, error)
} 