package postgres

import (
	"log"

	"github.com/lisandro/challenge/services/user-service/internal/domain"
	"gorm.io/gorm"
)

// UserFollow represents the database model for user follows
type UserFollow struct {
	FollowerID string `gorm:"type:uuid;primaryKey"`
	FollowedID string `gorm:"type:uuid;primaryKey"`
}

// RunMigrations performs database migrations using GORM
func RunMigrations(db *gorm.DB) error {
	log.Println("Running database migrations...")
	
	// AutoMigrate will create tables and add missing columns/indexes
	err := db.AutoMigrate(&domain.User{}, &UserFollow{})
	if err != nil {
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
} 