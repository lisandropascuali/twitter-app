package postgres

import (
	"log"
	"time"

	"gorm.io/gorm"
)

// Tweet represents the database model for tweets
type Tweet struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	UserID    string `gorm:"type:uuid"`
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// RunMigrations performs database migrations using GORM
func RunMigrations(db *gorm.DB) error {
	log.Println("Running database migrations...")
	
	// AutoMigrate will create tables and add missing columns/indexes
	err := db.AutoMigrate(&Tweet{})
	if err != nil {
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
} 