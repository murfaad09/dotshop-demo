package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/harishash/dotshop-be/integration/hooks/convictional/config"
	"github.com/harishash/dotshop-be/integration/hooks/convictional/handlers"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	db := config.SetupDatabase()
	app := fiber.New()

	app.Post("/myWebhookReceiver", handlers.WebhookReceiver(db))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3333"
	}

	go func() {
		log.Printf("Starting server on port %s\n", port)
		if err := app.Listen(":" + port); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Wait for the server to be up
	time.Sleep(2 * time.Second) // Adjust the sleep duration if necessary

	// Register the webhook using the updated config.RegisterWebhook function
	// if err := config.RegisterWebhook(); err != nil {
	// 	//log.Fatalf("Error registering webhook: %v", err)
	// }

	// Block the main function to keep the server running
	select {}
}
