package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	vendor_routes "github.com/harishash/dotshop-be/integration/vndr/convictional/routes"
	"github.com/harishash/dotshop-be/internal/routes"
	"github.com/harishash/dotshop-be/internal/utils/logger"
)

func main() {

	// Create Fiber app
	app := fiber.New()

	// Initialize routes
	routes.InitRoutes(app)
	vendor_routes.InitConvictionalRoutes(app)

	app.Use(logger.Fiber())
	// Initialize cache
	// client := cache.NewCacheClient()

	// Start server
	err := app.Listen(":8081")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
