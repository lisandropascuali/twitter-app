package usecase

import "github.com/lisandro/challenge/services/user-service/internal/domain"

// UserUsecase represents the user's business logic contract
type UserUsecase interface {
    Follow(followerID, followedID string) error
    Unfollow(followerID, followedID string) error
    GetFollowing(userID string) ([]string, error)
    GetUser(id string) (*domain.User, error)
} 