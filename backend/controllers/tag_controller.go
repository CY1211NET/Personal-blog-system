package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/your-username/blog-backend/database"
	"github.com/your-username/blog-backend/models"
)

func GetTags(c *gin.Context) {
	var tags []models.Tag
	if err := database.DB.Find(&tags).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tags)
}

func DeleteTag(c *gin.Context) {
	id := c.Param("id")
	var tag models.Tag

	if err := database.DB.First(&tag, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		return
	}

	// Delete all articles associated with this tag
	// First find article IDs
	var articleTags []struct {
		ArticleID uint
	}
	// Assuming article_tags table name. GORM default is `article_tags` for many2many
	if err := database.DB.Table("article_tags").Select("article_id").Where("tag_id = ?", id).Scan(&articleTags).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find associated articles"})
		return
	}

	var articleIDs []uint
	for _, at := range articleTags {
		articleIDs = append(articleIDs, at.ArticleID)
	}

	if len(articleIDs) > 0 {
		// Delete ALL tag associations for these articles (not just the current tag)
		if err := database.DB.Exec("DELETE FROM article_tags WHERE article_id IN ?", articleIDs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete article tags"})
			return
		}

		// Delete comments
		if err := database.DB.Exec("DELETE FROM comments WHERE article_id IN ?", articleIDs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete article comments"})
			return
		}

		// Hard delete articles to avoid issues with soft deleted records still referencing things
		if err := database.DB.Unscoped().Delete(&models.Article{}, articleIDs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete associated articles"})
			return
		}
	}

	// Delete the tag itself
	if err := database.DB.Delete(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag and associated articles deleted successfully"})
}

type UpdateTagInput struct {
	Name string `json:"name" binding:"required"`
}

func UpdateTag(c *gin.Context) {
	id := c.Param("id")
	var tag models.Tag

	if err := database.DB.First(&tag, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		return
	}

	var input UpdateTagInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tag.Name = input.Name
	if err := database.DB.Save(&tag).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tag)
}
