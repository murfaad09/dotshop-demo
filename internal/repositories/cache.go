package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/harishash/dotshop-be/internal/config"
	"github.com/redis/go-redis/v9"
)

type IUserCache interface {
	Store(key string, value interface{}) error
	Retrieve(key string) (interface{}, error)
}

type Cache struct {
	Client *redis.Client
}

var instance *Cache

// var once *sync.Once

var (
	ErrNil = errors.New("no matching record found in redis cache")
	Ctx    = context.TODO()
)

func NewCacheClient() *Cache {
	instance = &Cache{
		Client: connect(),
	}
	return instance
}

func connect() *redis.Client {
	c := config.GetConfig()

	fmt.Println("Redis URL: ", c.RedisUrl)

	opt, err := redis.ParseURL(c.RedisUrl)

	fmt.Println("Opt: ", opt)

	if err != nil {
		log.Fatalf("Error getting redis config: %v", err)
	}

	client := redis.NewClient(opt)
	if err := client.Ping(Ctx).Err(); err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	return client
}

func (c *Cache) Store(key string, value interface{}) error {
	storeValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = c.Client.Set(Ctx, key, storeValue, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (c *Cache) Retrieve(key string) (interface{}, error) {
	value, err := c.Client.Get(Ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return value, nil
}
