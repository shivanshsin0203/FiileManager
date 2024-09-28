package redis

import (
	"context"
	"fmt"
	"log"
	"os"
    	"github.com/joho/godotenv"
	"github.com/go-redis/redis/v8"
)

var (
	rdb *redis.Client
	ctx = context.Background()
)

// Init initializes the Redis client
func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	redisURL := os.Getenv("UPSTASH_REDIS_URL")
	if redisURL == "" {
		log.Fatal("UPSTASH_REDIS_URL environment variable is not set")
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatal("Failed to parse Redis URL:", err)
	}

	rdb = redis.NewClient(opt)

	// Test the connection
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}

	fmt.Println("Successfully connected to Redis")
}

// Enqueue adds an item to the queue
func Enqueue(queueName, item string) error {
	if rdb == nil {
		return fmt.Errorf("redis client not initialized. Call Init() first")
	}

	err := rdb.RPush(ctx, queueName, item).Err()
	if err != nil {
		return fmt.Errorf("failed to enqueue item: %v", err)
	}

	fmt.Printf("Successfully enqueued item '%s' to queue '%s'\n", item, queueName)
	return nil
}

// Close closes the Redis connection
func Close() {
	if rdb != nil {
		err := rdb.Close()
		if err != nil {
			log.Printf("Error closing Redis connection: %v", err)
		} else {
			fmt.Println("Redis connection closed")
		}
	}
}