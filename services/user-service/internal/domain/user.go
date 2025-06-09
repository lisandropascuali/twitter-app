package domain

// User represents a user in the system
type User struct {
    ID       string `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
    Username string `json:"username" gorm:"type:varchar(255);uniqueIndex:idx_users_username;not null"`
}

// CreateUserRequest represents the request to create a new user
type CreateUserRequest struct {
    Username string `json:"username" validate:"required,min=3,max=50"`
}

// UserRepository represents the user's repository contract
type UserRepository interface {
    GetAllUsers() ([]User, error)
    CreateUser(req CreateUserRequest) (*User, error)
    Follow(followerID, followedID string) error
    Unfollow(followerID, followedID string) error
    GetFollowing(userID string) ([]User, error)
    GetFollowers(userID string) ([]User, error)
	GetUser(id string) (*User, error)
}

// UserUsecase represents the user's business logic contract
type UserUsecase interface {
    GetAllUsers() ([]User, error)
    CreateUser(req CreateUserRequest) (*User, error)
    Follow(followerID, followedID string) error
    Unfollow(followerID, followedID string) error
    GetFollowing(userID string) ([]User, error)
    GetFollowers(userID string) ([]User, error)
} 