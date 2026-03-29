package config

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	DbUrl                string
	RedisUrl             string
	JWTSecret            string
	StripeSecretKey      string
	StripePublicKey      string
	ConvictionalBaseURL  string
	BuyerAPIKey          string
	SellerAPIKey         string
	PaypalClientID       string
	PaypalClientSecret   string
	PaypalApiURL         string
	KlaviyoKey           string
	AWSREGION            string
	AWSAccessID          string
	AWSSecretAccessKey   string
	DotShopStoreEmail    string
	DotShopStorePassword string
	DotShopAdminEmail    string
	DotShopAdminPassword string
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		config, err := load()
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
		instance = config
	})
	return instance
}

func load() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("POSTGRES_SERVER"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
	)

	redisUrl := fmt.Sprintf("redis://%s:%s@%s:%s",
		os.Getenv("REDIS_USER"),
		os.Getenv("REDIS_PASSWORD"),
		os.Getenv("REDIS_SERVER"),
		os.Getenv("REDIS_PORT"))

	return &Config{
		DbUrl:                dsn,
		RedisUrl:             redisUrl,
		JWTSecret:            os.Getenv("JWT_SECRET"),
		StripeSecretKey:      os.Getenv("STRIPE_SECRET_KEY"),
		StripePublicKey:      os.Getenv("STRIPE_PUBLISHABLE_KEY"),
		ConvictionalBaseURL:  os.Getenv("CONVICTIONAL_BASE_URL"),
		BuyerAPIKey:          os.Getenv("BUYER_API_KEY"),
		SellerAPIKey:         os.Getenv("SELLER_API_KEY"),
		PaypalClientID:       os.Getenv("PAYPAL_CLIENT_ID"),
		PaypalClientSecret:   os.Getenv("PAYPAL_CLIENT_SECRET"),
		PaypalApiURL:         os.Getenv("PAYPAY_API_URL"),
		KlaviyoKey:           os.Getenv("KLAVIYO_API_KEY"),
		AWSREGION:            os.Getenv("AWS_REGION"),
		AWSAccessID:          os.Getenv("AWS_ACCESS_KEY_ID"),
		AWSSecretAccessKey:   os.Getenv("AWS_SECRET_ACCESS_KEY"),
		DotShopStoreEmail:    os.Getenv("DOTSHOP_STORE_EMAIL"),
		DotShopStorePassword: os.Getenv("DOTSHOP_STORE_PASSWORD"),
		DotShopAdminEmail:    os.Getenv("DOTSHOP_ADMIN_EMAIL"),
		DotShopAdminPassword: os.Getenv("DOTSHOP_ADMIN_PASSWORD"),
	}, nil
}
