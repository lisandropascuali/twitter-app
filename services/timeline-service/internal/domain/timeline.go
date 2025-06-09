package domain

import "time"

type Timeline struct {
	Tweets []Tweet `json:"tweets"`
}

type Tweet struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// FollowingUser represents a user that someone follows - matches the User struct from user-service
type FollowingUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
} 