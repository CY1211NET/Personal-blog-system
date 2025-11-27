package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/your-username/blog-backend/database"
	"github.com/your-username/blog-backend/models"
)

type CreateArticleInput struct {
	Title      string   `json:"title" binding:"required"`
	Content    string   `json:"content" binding:"required"`
	CategoryID *uint    `json:"category_id"`
	Tags       []string `json:"tags"` // List of tag names
}

type UpdateArticleInput struct {
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	CategoryID *uint    `json:"category_id"`
	Tags       []string `json:"tags"`
}

func CreateArticle(c *gin.Context) {
	var input CreateArticleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Handle Tags
	var tags []models.Tag
	for _, tagName := range input.Tags {
		var tag models.Tag
		if err := database.DB.FirstOrCreate(&tag, models.Tag{Name: tagName}).Error; err != nil {
			continue
		}
		tags = append(tags, tag)
	}

	article := models.Article{
		Title:      input.Title,
		Content:    input.Content,
		AuthorID:   userID.(uint),
		CategoryID: input.CategoryID,
		Tags:       tags,
	}

	if err := database.DB.Create(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Preload associations for response
	database.DB.Preload("Author").Preload("Category").Preload("Tags").First(&article, article.ID)

	c.JSON(http.StatusOK, article)
}

func GetArticles(c *gin.Context) {
	var articles []models.Article
	query := database.DB.Preload("Author").Preload("Category").Preload("Tags")

	search := c.Query("search")
	if search != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	categoryID := c.Query("category_id")
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	if err := query.Find(&articles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, articles)
}

func GetArticle(c *gin.Context) {
	id := c.Param("id")
	var article models.Article

	if err := database.DB.Preload("Author").Preload("Category").Preload("Tags").First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	c.JSON(http.StatusOK, article)
}

func UpdateArticle(c *gin.Context) {
	id := c.Param("id")
	var article models.Article

	if err := database.DB.First(&article, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if article.AuthorID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not the author of this article"})
		return
	}

	var input UpdateArticleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update basic fields
	if input.Title != "" {
		article.Title = input.Title
	}
	if input.Content != "" {
		article.Content = input.Content
	}
	if input.CategoryID != nil {
		article.CategoryID = input.CategoryID
	}

	// Update Tags if provided
	if input.Tags != nil {
		var tags []models.Tag
		for _, tagName := range input.Tags {
			var tag models.Tag
			if err := database.DB.FirstOrCreate(&tag, models.Tag{Name: tagName}).Error; err != nil {
				continue
			}
			tags = append(tags, tag)
		}
		// Replace tags
		database.DB.Model(&article).Association("Tags").Replace(tags)
	}

	if err := database.DB.Save(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Preload associations for response
	database.DB.Preload("Author").Preload("Category").Preload("Tags").First(&article, article.ID)

	c.JSON(http.StatusOK, article)
}

func DeleteArticle(c *gin.Context) {
	var article models.Article
	if err := database.DB.Where("id = ?", c.Param("id")).First(&article).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Article not found"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if article.AuthorID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not the author of this article"})
		return
	}

	if err := database.DB.Delete(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Article deleted successfully"})
}

func LikeArticle(c *gin.Context) {
	var article models.Article
	if err := database.DB.First(&article, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	// In a real app, you'd track *who* liked it to prevent duplicates
	article.Likes++
	database.DB.Save(&article)
	c.JSON(http.StatusOK, gin.H{"likes": article.Likes})
}

func ViewArticle(c *gin.Context) {
	var article models.Article
	if err := database.DB.First(&article, c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Article not found"})
		return
	}

	article.Views++
	database.DB.Save(&article)
	c.JSON(http.StatusOK, gin.H{"views": article.Views})
}
