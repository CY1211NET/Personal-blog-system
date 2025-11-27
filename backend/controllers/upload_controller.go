package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file is received"})
		return
	}

	// Validate file extension
	extension := filepath.Ext(file.Filename)
	if extension != ".jpg" && extension != ".jpeg" && extension != ".png" && extension != ".gif" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only jpg, jpeg, png, gif are allowed"})
		return
	}

	// Generate unique filename
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), extension)
	dst := filepath.Join("uploads", filename)

	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save file"})
		return
	}

	// Return the URL to access the file
	// Assuming the server is running on localhost:8080 and serving static files from /uploads
	fileURL := fmt.Sprintf("http://localhost:8080/uploads/%s", filename)

	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
		"url":     fileURL,
	})
}
