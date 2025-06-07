package postgres

import (
	"github.com/lisandro/challenge/services/user-service/internal/domain"
	"gorm.io/gorm"
)

// PostgresRepository handles user data persistence in PostgreSQL
type PostgresRepository struct {
	db *gorm.DB
}

// NewPostgresRepository creates a new PostgreSQL repository
func NewPostgresRepository(db *gorm.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) Follow(followerID, followedID string) error {
	follow := UserFollow{
		FollowerID: followerID,
		FollowedID: followedID,
	}
	return r.db.Create(&follow).Error
}

func (r *PostgresRepository) Unfollow(followerID, followedID string) error {
	return r.db.Where("follower_id = ? AND followed_id = ?", followerID, followedID).Delete(&UserFollow{}).Error
}

func (r *PostgresRepository) GetFollowing(userID string) ([]domain.User, error) {
	var follows []UserFollow
	if err := r.db.Where("follower_id = ?", userID).Find(&follows).Error; err != nil {
		return nil, err
	}

	var users []domain.User
	for _, follow := range follows {
		var user domain.User
		if err := r.db.First(&user, "id = ?", follow.FollowedID).Error; err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *PostgresRepository) GetUser(id string) (*domain.User, error) {
	var user domain.User
	if err := r.db.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetFollowers returns the list of users that follow a user
func (r *PostgresRepository) GetFollowers(userID string) ([]domain.User, error) {
	var follows []UserFollow
	if err := r.db.Where("followed_id = ?", userID).Find(&follows).Error; err != nil {
		return nil, err
	}

	var users []domain.User
	for _, follow := range follows {
		var user domain.User
		if err := r.db.First(&user, "id = ?", follow.FollowerID).Error; err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
} 