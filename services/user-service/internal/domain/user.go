package domain

// User represents a user in the system
type User struct {
    ID       string `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    Username string `json:"username" gorm:"type:varchar(255);uniqueIndex:idx_users_username;not null"`
}

// UserRepository represents the user's repository contract
type UserRepository interface {
    Follow(followerID, followedID string) error
    Unfollow(followerID, followedID string) error
    GetFollowing(userID string) ([]User, error)
    GetFollowers(userID string) ([]User, error)
}

// UserUsecase represents the user's business logic contract
type UserUsecase interface {
    Follow(followerID, followedID string) error
    Unfollow(followerID, followedID string) error
    GetFollowing(userID string) ([]User, error)
    GetFollowers(userID string) ([]User, error)
} 