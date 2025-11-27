package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/your-username/blog-backend/controllers"
	"github.com/your-username/blog-backend/middlewares"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			auth := v1.Group("/auth")
			{
				auth.POST("/register", controllers.Register)
				auth.POST("/login", controllers.Login)
			}

			protected := v1.Group("/user")
			protected.Use(middlewares.JwtAuthMiddleware())
			{
				protected.GET("/profile", func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "Welcome to protected route"})
				})
			}

			// Public Article Routes
			v1.GET("/articles", controllers.GetArticles)
			v1.GET("/articles/:id", controllers.GetArticle)

			// Protected Article Routes
			articles := v1.Group("/articles")
			articles.Use(middlewares.JwtAuthMiddleware())
			{
				articles.POST("/", controllers.CreateArticle)
				articles.PUT("/:id", controllers.UpdateArticle)
				articles.DELETE("/:id", controllers.DeleteArticle)
				articles.POST("/:id/comments", controllers.CreateComment)
			}

			// Public Comment Routes
			v1.GET("/articles/:id/comments", controllers.GetComments)

			// Protected Comment Routes
			comments := v1.Group("/comments")
			comments.Use(middlewares.JwtAuthMiddleware())
			{
				comments.DELETE("/:id", controllers.DeleteComment)
			}

			// Upload Route
			v1.POST("/upload", middlewares.JwtAuthMiddleware(), controllers.UploadFile)
		}
	}

	// Serve static files
	r.Static("/uploads", "./uploads")

	return r
}
