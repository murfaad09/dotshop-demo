package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func WebhookReceiver(db *gorm.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// secret := os.Getenv("SECRET")
		// if secret == "" {
		// 	return c.Status(http.StatusInternalServerError).SendString("Secret not set")
		// }

		// if !validateSignature(c, secret) {
		// 	return c.Status(http.StatusForbidden).SendString("INVALID SIGNATURE")
		// }

		var data map[string]interface{}
		if err := json.Unmarshal(c.Body(), &data); err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Error unmarshaling request body")
		}
		fmt.Println("data =================================", data)

		logReceivedData(data)

		if err := processAndInsertData(db, data); err != nil {
			return c.Status(http.StatusInternalServerError).SendString("Error processing data")
		}

		return c.SendString("OK")
	}
}

func validateSignature(c *fiber.Ctx, secret string) bool {
	convictionalSignature := c.Get("Convictional-Signature")
	if convictionalSignature == "" {
		return false
	}

	timestampAndSignatures := strings.Split(convictionalSignature, ",")
	if len(timestampAndSignatures) < 2 {
		return false
	}
	timestamp := timestampAndSignatures[0]
	signatures := timestampAndSignatures[1:]

	body := c.Body()
	for _, receivedSignature := range signatures {
		h := hmac.New(sha256.New, []byte(secret))
		h.Write([]byte(fmt.Sprintf("%s.%s", timestamp, string(body))))
		generatedSignature := hex.EncodeToString(h.Sum(nil))
		if receivedSignature == generatedSignature {
			return true
		}
	}
	return false
}

func logReceivedData(data map[string]interface{}) {
	fmt.Println("Load Date +++++++++++++++++++")

	log.Println("Received data:")
	for key, value := range data {
		log.Printf("%s: %v\n", key, value)
	}
}

func processAndInsertData(db *gorm.DB, data map[string]interface{}) error {
	fmt.Println("processAndInsertData =================================")
	eventType, ok := data["event_type"].(string)
	if !ok {
		return fmt.Errorf("event_type not found or not a string")
	}

	switch eventType {
	case "order.cancelled":
		return handleProductUpdated(db, data)
	case "order.fulfilled":
		return handleOrderCreated(db, data)
	default:
		log.Printf("Unhandled event type: %s\n", eventType)
		return nil
	}
}

func handleProductUpdated(db *gorm.DB, data map[string]interface{}) error {
	// if product, ok := data["product"].(map[string]interface{}); ok {
	// 	id := product["id"].(string)
	// 	quantity := product["quantity"].(float64) // adjust as needed

	// 	result := db.Exec(`INSERT INTO products (id, quantity) VALUES ($1, $2)
	//                        ON CONFLICT (id) DO UPDATE SET quantity = EXCLUDED.quantity`,
	// 		id, quantity)
	// 	if result.Error != nil {
	// 		return fmt.Errorf("error inserting data: %v", result.Error)
	// 	}
	// 	log.Printf("Inserted/Updated product with ID %s and quantity %f\n", id, quantity)
	// }
	// return nil
	log.Printf("============================= Order cancelled ===============================: %v\n", data)
	return nil

}

func handleOrderCreated(db *gorm.DB, data map[string]interface{}) error {
	// if order, ok := data["order"].(map[string]interface{}); ok {
	// 	id := order["id"].(string)
	// 	status := order["status"].(string)

	// 	result := db.Exec(`INSERT INTO orders (id, status) VALUES ($1, $2)
	//                        ON CONFLICT (id) DO UPDATE SET status = EXCLUDED.status`,
	// 		id, status)
	// 	if result.Error != nil {
	// 		return fmt.Errorf("error inserting data: %v", result.Error)
	// 	}
	// 	log.Printf("Inserted/Updated order with ID %s and status %s\n", id, status)
	// }
	log.Printf("============================= Order fulfilled ===============================: %v\n", data)
	return nil

}
