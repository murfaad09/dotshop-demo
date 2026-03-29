package main

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/harishash/dotshop-be/cron-jobs/jobs/inventory_management"
	"github.com/harishash/dotshop-be/cron-jobs/jobs/order_status_management"
	"github.com/harishash/dotshop-be/integration/vndr/convictional/client"
	"github.com/robfig/cron/v3"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading .env file: %v", err)
	}
	viper.AutomaticEnv()
}

func main() {

	convClient := setUpConvictionalClient()

	inventoryManagementJobSchedule := viper.GetString("INVENTORY_MANAGEMENT_CRON_SCHEDULE")
	orderConfirmationJobSchedule := viper.GetString("ORDER_STATUS_MANAGEMENT_CRON_SCHEDULE")

	port := viper.GetString("PORT")

	// Create database connection
	db, err := setupDatabaseConnection()
	if err != nil {
		log.Fatalf("Error setting up database connection: %v", err)
	}
	defer dbSQLClose(db)

	c := cron.New()
	// Schedule Inventory Management Job with DB injection
	_, err = c.AddFunc(inventoryManagementJobSchedule, func() {
		inventory_management.Run(db, convClient)
	})
	if err != nil {
		log.Fatalf("Error scheduling inventory management job: %v", err)
	}

	// Schedule Order Confirmation Management Job with DB injection
	_, err = c.AddFunc(orderConfirmationJobSchedule, func() {
		order_status_management.Run(db, convClient)
	})
	if err != nil {
		log.Fatalf("Error scheduling order confirmation job: %v", err)
	}

	c.Start()

	// Create a new Fiber instance
	app := fiber.New()

	// Define a simple route to verify the server is running
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Server is running")
	})

	log.Printf("Starting server on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

func setupDatabaseConnection() (*gorm.DB, error) {
	user := viper.GetString("POSTGRES_USER")
	password := viper.GetString("POSTGRES_PASSWORD")
	dbName := viper.GetString("POSTGRES_DB")
	port := viper.GetString("POSTGRES_PORT")
	host := viper.GetString("POSTGRES_SERVER")
	dsn := "host=" + host +
		" user=" + user +
		" password=" + password +
		" dbname=" + dbName +
		" port=" + port +
		" sslmode=disable TimeZone=Asia/Karachi"

	// Open the database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Check if the database is alive
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = sqlDB.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func dbSQLClose(db *gorm.DB) {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Error getting sqlDB: %v", err)
	}
	sqlDB.Close()
}

func setUpConvictionalClient() *client.APIClient {
	baseUrl := viper.GetString("CONVICTIONAL_BASE_URL")
	apiKey := viper.GetString("BUYER_API_KEY")
	return client.NewAPIClient(baseUrl, apiKey)
}
