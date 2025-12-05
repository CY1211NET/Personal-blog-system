package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/your-username/blog-backend/controllers"
	"github.com/your-username/blog-backend/middlewares"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// CORS Middleware
	r.Use(middlewares.CORSMiddleware())

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

			user := v1.Group("/user")
			user.Use(middlewares.JwtAuthMiddleware())
			{
				user.GET("/profile", controllers.GetProfile)
				user.PUT("/profile", controllers.UpdateProfile)
			}

			// Public User Routes
			v1.GET("/author-profile", controllers.GetAuthorProfile)
			v1.GET("/registration-status", controllers.CheckRegistrationStatus)

			// Public Article Routes
			v1.GET("/articles", controllers.GetArticles)
			v1.GET("/articles/:id", controllers.GetArticle)
			v1.POST("/articles/:id/view", controllers.ViewArticle)

			// Protected Article Routes
			articles := v1.Group("/articles")
			articles.Use(middlewares.JwtAuthMiddleware())
			{
				articles.POST("/", controllers.CreateArticle)
				articles.PUT("/:id", controllers.UpdateArticle)
				articles.DELETE("/:id", controllers.DeleteArticle)
				articles.POST("/batch-delete", controllers.BatchDeleteArticles)
				articles.POST("/:id/comments", controllers.CreateComment)
				articles.POST("/:id/like", controllers.LikeArticle)
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

			// Category Routes
			v1.GET("/categories", controllers.GetCategories)
			v1.POST("/categories", middlewares.JwtAuthMiddleware(), controllers.CreateCategory)
			v1.PUT("/categories/:id", middlewares.JwtAuthMiddleware(), controllers.UpdateCategory)
			v1.DELETE("/categories/:id", middlewares.JwtAuthMiddleware(), controllers.DeleteCategory)

			// Tag Routes
			v1.GET("/tags", controllers.GetTags)
			v1.PUT("/tags/:id", middlewares.JwtAuthMiddleware(), controllers.UpdateTag)
			v1.DELETE("/tags/:id", middlewares.JwtAuthMiddleware(), controllers.DeleteTag)
		}
	}

	// Serve static files
	r.Static("/uploads", "./uploads")

	return r
}
