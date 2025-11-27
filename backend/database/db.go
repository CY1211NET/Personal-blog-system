package database

import (
	"fmt"
	"log"

	"github.com/your-username/blog-backend/config"
	"github.com/your-username/blog-backend/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.AppConfig.DBUsername,
		config.AppConfig.DBPassword,
		config.AppConfig.DBHost,
		config.AppConfig.DBPort,
		config.AppConfig.DBName,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	log.Println("Database connected successfully")

	if err := DB.AutoMigrate(&models.User{}, &models.Article{}, &models.Comment{}, &models.Category{}, &models.Tag{}); err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}
}
