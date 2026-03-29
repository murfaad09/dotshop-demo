package order_status_management

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/harishash/dotshop-be/integration/vndr/convictional/client"
	"github.com/harishash/dotshop-be/internal/dto"
	"gorm.io/gorm"
)

const (
	batchSize    = 50 // Number of orders to process in one batch
	apiRateLimit = 4  // Max API calls per second to Convictional
	StatusPending   = "Pending"
	StatusConfirmed = "Confirmed"
	StatusDelivered = "Delivered"
	StatusCancelled = "Cancelled"
)

func Run(db *gorm.DB, convClient *client.APIClient) {
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Printf("Error getting sqlDB: %v\n", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = sqlDB.PingContext(ctx)
	if err != nil {
		fmt.Printf("Error pinging database: %v\n", err)
		return
	}

	orderIDs := getPendingOrderIds(db)

	if len(orderIDs) == 0 {
		fmt.Println("No orders found")
		return
	}

	concurrentBatches := len(orderIDs) / batchSize
	if len(orderIDs)%batchSize != 0 {
		concurrentBatches++
	}

	if concurrentBatches > apiRateLimit {
		concurrentBatches = apiRateLimit
	}

	// Create channels for batching and result processing
	idChan := make(chan string, len(orderIDs))
	resultChan := make(chan error, len(orderIDs))

	// Start workers for processing order batches
	var wg sync.WaitGroup
	for i := 0; i < concurrentBatches; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			processBatch(db, convClient, idChan, resultChan)
		}()
	}

	// Distribute order IDs into the channel
	go func() {
		for _, id := range orderIDs {
			idChan <- id
		}
		close(idChan)
	}()

	// Wait for all batches to complete
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Handle results
	for err := range resultChan {
		if err != nil {
			fmt.Printf("Error processing product: %v\n", err)
		}
	}
}

func getPendingOrderIds(db *gorm.DB) []string {
	var orderIDs []string
	db.Table("orders").Where("status ILIKE ?", "pending").Pluck("id", &orderIDs)
	return orderIDs
}

func processBatch(db *gorm.DB, convClient *client.APIClient, idChan <-chan string, resultChan chan<- error) {
	ticker := time.NewTicker(time.Second / apiRateLimit)
	defer ticker.Stop()

	for id := range idChan {
		<-ticker.C
		url := fmt.Sprintf("https://api.convictional.com/buyer/orders/%s", id)
		orders, err := getOrders(convClient, url)
		if err != nil {
			resultChan <- fmt.Errorf("error getting variants for product %s: %w", id, err)
			continue
		}
		status := aggregateOrderStatus(*orders)
		fmt.Printf("Orders: %v\n", orders)
		if status != StatusPending {
			if err := updateOrderStatusInDB(db, status, orders); err != nil {
				resultChan <- fmt.Errorf("error adding variants data for product %s: %w", id, err)
			}
		}

		resultChan <- nil
	}
}

func aggregateOrderStatus(order dto.Convictional_Order_Status_Data) string {
	allConfirmed := true
	for _, sellerOrder := range order.Data.SellerOrders {
		if !sellerOrder.Fulfilled {
			allConfirmed = false
			break
		}
	}
	if allConfirmed {
		return StatusConfirmed
	}

	return StatusPending
}

func getOrders(c *client.APIClient, url string) (*dto.Convictional_Order_Status_Data, error) {
	data, err := client.GetPendingOrders(c, url)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func updateOrderStatusInDB(db *gorm.DB, status string, order *dto.Convictional_Order_Status_Data) error {

	updateFields := map[string]interface{}{
		"status":            status,
		"status_updated_at": order.Data.SellerOrders[0].FulfilledDate,
	}

	if err := db.Table("orders").Where("id = ?", order.Data.ID).Updates(updateFields).Error; err != nil {
		return err
	}

	return nil
}
