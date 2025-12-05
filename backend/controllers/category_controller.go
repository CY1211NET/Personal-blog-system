package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/your-username/blog-backend/database"
	"github.com/your-username/blog-backend/models"
)

type CreateCategoryInput struct {
	Name string `json:"name" binding:"required"`
}

func CreateCategory(c *gin.Context) {
	var input CreateCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category := models.Category{Name: input.Name}
	if err := database.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

func GetCategories(c *gin.Context) {
	var categories []models.Category
	if err := database.DB.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category

	if err := database.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	// Find all articles in this category
	var articles []models.Article
	if err := database.DB.Where("category_id = ?", id).Find(&articles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find associated articles"})
		return
	}

	var articleIDs []uint
	for _, a := range articles {
		articleIDs = append(articleIDs, a.ID)
	}

	if len(articleIDs) > 0 {
		// Delete mappings in article_tags
		if err := database.DB.Exec("DELETE FROM article_tags WHERE article_id IN ?", articleIDs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete article tags"})
			return
		}

		// Delete comments
		if err := database.DB.Exec("DELETE FROM comments WHERE article_id IN ?", articleIDs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete article comments"})
			return
		}

		// Hard delete articles
		if err := database.DB.Unscoped().Delete(&models.Article{}, articleIDs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete associated articles"})
			return
		}
	}

	if err := database.DB.Delete(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category and associated articles deleted successfully"})
}

type UpdateCategoryInput struct {
	Name string `json:"name" binding:"required"`
}

func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category

	if err := database.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	var input UpdateCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category.Name = input.Name
	if err := database.DB.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}
