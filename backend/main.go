package main

import (
	"github.com/your-username/blog-backend/config"
	"github.com/your-username/blog-backend/database"
	"github.com/your-username/blog-backend/routes"
)

func main() {
	// Load Config
	config.LoadConfig()

	// Connect Database
	database.ConnectDB()

	// Setup Router
	r := routes.SetupRouter()

	// Run Server
	r.Run(":" + config.AppConfig.ServerPort)
}
