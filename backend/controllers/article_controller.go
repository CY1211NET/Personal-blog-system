package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/your-username/blog-backend/database"
	"github.com/your-username/blog-backend/models"
)

type CreateArticleInput struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type UpdateArticleInput struct {
	Title   string `json:"title"`
	Content string `json:"content"`
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

	article := models.Article{
		Title:    input.Title,
		Content:  input.Content,
		AuthorID: userID.(uint),
	}

	if err := database.DB.Create(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, article)
}

func GetArticles(c *gin.Context) {
	var articles []models.Article
	query := database.DB.Preload("Author")

	search := c.Query("search")
	if search != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Find(&articles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, articles)
}

func GetArticle(c *gin.Context) {
	var article models.Article
	if err := database.DB.Preload("Author").Where("id = ?", c.Param("id")).First(&article).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Article not found"})
		return
	}

	c.JSON(http.StatusOK, article)
}

func UpdateArticle(c *gin.Context) {
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

	var input UpdateArticleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Model(&article).Updates(models.Article{Title: input.Title, Content: input.Content}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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
