package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDatabase() *gorm.DB {
	// dsn := os.Getenv("DATABASE_URL")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("POSTGRES_SERVER"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// Auto-migrate your models here
	// db.AutoMigrate(&models.Product{}, &models.Order{})

	return db
}

func RegisterWebhook() error {
	url := "https://api.convictional.com/webhooks"
	payload := strings.NewReader("{\"limiterBurst\":0,\"limiterRate\":0,\"topics\":[\"order.cancelled\"],\"secrets\":[{\"expiresDate\":\"2025-06-25T19:00:00.000+00:00\",\"secret\":\"" + os.Getenv("SECRET") + "\"}],\"targetUrl\":\"" + os.Getenv("WEBHOOK_URL") + "\"}")

	//payload := strings.NewReader("{\"limiterBurst\":0,\"limiterRate\":0,\"topics\":[\"order.cancelled\"],\"secrets\":[{\"expiresDate\":\"2025-06-25T19:00:00.000+00:00\",\"secret\":\"wCOYcm6OleDDDW5cOUzZAeBWx29BLoRf\"}],\"targetUrl\":\"https://e9f1-110-93-204-190.ngrok-free.app/myWebhookReceiver\"}")
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", os.Getenv("CONVICTIONAL_API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Response Body: %s", string(body))
		return fmt.Errorf("error registering webhook: received status code %d", resp.StatusCode)
	}

	log.Println("Webhook registered successfully")
	fmt.Println(string(body))
	return nil
}
