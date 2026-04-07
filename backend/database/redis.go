package database

import (
	"context"
	"fmt"
	"log"

	"teras-vps/backend/config"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

// InitRedis initializes Redis connection
func InitRedis(cfg *config.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPass,
		DB:       0,
	})

	// Test connection
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Printf("⚠️  Warning: Failed to connect to Redis: %v", err)
		return rdb
	}

	log.Println("✅ Redis connected successfully")
	return rdb
}
