package usecase

import "github.com/lisandro/challenge/services/user-service/internal/domain"

// UserUsecase represents the user's business logic contract
type UserUsecase interface {
    GetAllUsers() ([]domain.User, error)
    CreateUser(req domain.CreateUserRequest) (*domain.User, error)
    Follow(followerID, followedID string) error
    Unfollow(followerID, followedID string) error
    GetFollowing(userID string) ([]domain.User, error)
    GetFollowers(userID string) ([]domain.User, error)
    GetUser(id string) (*domain.User, error)
} 