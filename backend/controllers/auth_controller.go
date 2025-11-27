package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/your-username/blog-backend/database"
	"github.com/your-username/blog-backend/models"
	"github.com/your-username/blog-backend/utils"
)

type RegisterInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if a user already exists (single user blog)
	var count int64
	database.DB.Model(&models.User{}).Count(&count)
	if count > 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Registration is closed. This is a single-user blog."})
		return
	}

	u := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
		return
	}
	u.Password = hashedPassword

	if err := database.DB.Create(&u).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Registration success"})
}

func Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}

	if err := utils.CheckPassword(input.Password, user.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid password"})
		return
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

type UpdateProfileInput struct {
	AvatarURL    string `json:"avatar_url"`
	Bio          string `json:"bio"`
	SocialLinks  string `json:"social_links"`
	SponsorLinks string `json:"sponsor_links"`
	FriendLinks  string `json:"friend_links"`
}

func UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var input UpdateProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Update all fields
	if input.AvatarURL != "" {
		user.AvatarURL = input.AvatarURL
	}
	user.Bio = input.Bio
	user.SocialLinks = input.SocialLinks
	user.SponsorLinks = input.SponsorLinks
	user.FriendLinks = input.FriendLinks

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func GetAuthorProfile(c *gin.Context) {
	// For single user blog, just get the first user or specific ID.
	// Let's assume ID 1 is the admin/author.
	var user models.User
	if err := database.DB.First(&user, 1).Error; err != nil {
		// Try first user if ID 1 not found (maybe different ID)
		if err := database.DB.First(&user).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Author not found"})
			return
		}
	}
	c.JSON(http.StatusOK, user)
}

func CheckRegistrationStatus(c *gin.Context) {
	var count int64
	database.DB.Model(&models.User{}).Count(&count)
	c.JSON(http.StatusOK, gin.H{
		"registration_allowed": count == 0,
		"user_count":           count,
	})
}
