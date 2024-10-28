package main

import (
	"demo-go/api/handlers"
	"demo-go/internal/config"
	"demo-go/internal/database"

	"github.com/gin-gonic/gin"
)

func main() {
    // Load configuration
    cfg := config.LoadConfig()

    // Initialize database
    db := database.InitDB(cfg.Database.DSN)

    // Set up the Gin router
    r := gin.Default()

    // Register routes
    handlers.RegisterRoutes(r, db)

    // Start the server
    r.Run(cfg.Server.Address)
}
