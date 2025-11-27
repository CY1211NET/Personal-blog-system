package models

import (
	"time"

	"gorm.io/gorm"
)

type Article struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Title      string         `gorm:"type:varchar(255);not null" json:"title"`
	Content    string         `gorm:"type:text;not null" json:"content"`
	AuthorID   uint           `gorm:"not null" json:"author_id"`
	Author     User           `gorm:"foreignKey:AuthorID" json:"author"`
	CategoryID *uint          `json:"category_id"`
	Category   Category       `gorm:"foreignKey:CategoryID" json:"category"`
	Tags       []Tag          `gorm:"many2many:article_tags;" json:"tags"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
