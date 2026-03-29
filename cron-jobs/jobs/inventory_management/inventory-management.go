package inventory_management

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/harishash/dotshop-be/integration/vndr/convictional/client"
	"github.com/harishash/dotshop-be/internal/dto"
	"gorm.io/gorm"
)

const (
	batchSize    = 50 // Number of products to process in one batch
	apiRateLimit = 4  // Max API calls per second to Convictional
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

	productIds := getProductIds(db)

	if len(productIds) == 0 {
		fmt.Println("No products found")
		return
	}

	concurrentBatches := len(productIds) / batchSize
	if len(productIds)%batchSize != 0 {
		concurrentBatches++
	}

	if concurrentBatches > apiRateLimit {
		concurrentBatches = apiRateLimit
	}

	// Create channels for batching and result processing
	idChan := make(chan string, len(productIds))
	resultChan := make(chan error, len(productIds))

	// Start workers for processing product batches
	var wg sync.WaitGroup
	for i := 0; i < concurrentBatches; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			processBatch(db, convClient, idChan, resultChan)
		}()
	}

	// Distribute product IDs into the channel
	go func() {
		for _, id := range productIds {
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

func getProductIds(db *gorm.DB) []string {
	var productIds []string
	db.Table("products").Pluck("product_id", &productIds)
	return productIds
}

func processBatch(db *gorm.DB, convClient *client.APIClient, idChan <-chan string, resultChan chan<- error) {
	ticker := time.NewTicker(time.Second / apiRateLimit)
	defer ticker.Stop()

	for id := range idChan {
		<-ticker.C
		url := fmt.Sprintf("https://api.convictional.com/buyer/products/%s/variants", id)
		variants, err := getVariants(convClient, url)
		if err != nil {
			resultChan <- fmt.Errorf("error getting variants for product %s: %w", id, err)
			continue
		}

		for _, variant := range variants.Variants {
			fmt.Printf("Variant: %v\n", variant)
			if err := addVariantsData(db, &variant); err != nil {
				resultChan <- fmt.Errorf("error adding variants data for product %s: %w", id, err)
				continue
			}
		}
		resultChan <- nil
	}
}

func getVariants(c *client.APIClient, url string) (*dto.Convictional_Variants_Data, error) {
	data, err := client.GetVariants(c, url)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func addVariantsData(db *gorm.DB, variant *dto.Variants) error {
	var existingVariant dto.Variants
	if err := db.Table("variants").Where("id = ?", variant.ID).First(&existingVariant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			fmt.Println("Variant does not exist in the database")
			return nil
		}
		return err
	}

	updateFields := map[string]interface{}{
		"inventory_amount": variant.InventoryAmount,
		"retail_price":     variant.RetailPrice,
		"retail_currency":  variant.RetailCurrency,
		"base_price":       variant.BasePrice,
		"base_currency":    variant.BaseCurrency,
	}

	if err := db.Table("variants").Where("id = ?", variant.ID).Updates(updateFields).Error; err != nil {
		return err
	}

	return nil
}
