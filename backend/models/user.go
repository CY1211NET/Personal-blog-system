package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"username"`
	Email        string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password     string         `gorm:"not null" json:"-"` // Don't return password in JSON
	AvatarURL    string         `json:"avatar_url"`
	Bio          string         `gorm:"type:text" json:"bio"`
	SocialLinks  string         `gorm:"type:text" json:"social_links"`  // JSON string
	SponsorLinks string         `gorm:"type:text" json:"sponsor_links"` // JSON string
	FriendLinks  string         `gorm:"type:text" json:"friend_links"`  // JSON string
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
