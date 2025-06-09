package usecase

import (
	"github.com/lisandro/challenge/services/user-service/internal/domain"
)

// userUsecase implements domain.UserUsecase
type userUsecase struct {
	repo domain.UserRepository
}

// NewUserUsecase creates a new user usecase instance
func NewUserUsecase(repo domain.UserRepository) domain.UserUsecase {
	return &userUsecase{
		repo: repo,
	}
}

// Follow makes a user follow another user
func (u *userUsecase) Follow(followerID, followedID string) error {
	return u.repo.Follow(followerID, followedID)
}

// Unfollow makes a user unfollow another user
func (u *userUsecase) Unfollow(followerID, followedID string) error {
	return u.repo.Unfollow(followerID, followedID)
}

// GetFollowing returns the list of users that a user follows
func (u *userUsecase) GetFollowing(userID string) ([]domain.User, error) {
	return u.repo.GetFollowing(userID)
}

// GetFollowers returns the list of users that follow a user
func (u *userUsecase) GetFollowers(userID string) ([]domain.User, error) {
	return u.repo.GetFollowers(userID)
}

func (u *userUsecase) CreateUser(req domain.CreateUserRequest) (*domain.User, error) {
	return u.repo.CreateUser(req)
}

// GetAllUsers returns all users in the system
func (u *userUsecase) GetAllUsers() ([]domain.User, error) {
	return u.repo.GetAllUsers()
}

// GetUser returns a user by ID
func (u *userUsecase) GetUser(id string) (*domain.User, error) {
	return u.repo.GetUser(id)
}