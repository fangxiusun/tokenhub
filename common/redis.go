package common

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	RedisClient *redis.Client
	ctx         = context.Background()
)

// InitRedis initializes the Redis client
func InitRedis() {
	redisConn := GetEnvOrDefault("REDIS_CONN_STRING", "")
	if redisConn == "" {
		log.Println("Redis connection string not provided, using in-memory cache")
		return
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisConn,
		Password: GetEnvOrDefault("REDIS_PASSWORD", ""),
		DB:       GetEnvOrDefaultInt("REDIS_DB", 0),
	})

	// Test connection
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Printf("Warning: Failed to connect to Redis: %v", err)
		RedisClient = nil
		return
	}

	log.Println("Redis connected successfully")
}

// IsRedisAvailable returns whether Redis is available
func IsRedisAvailable() bool {
	return RedisClient != nil
}

// SetRedisKey sets a key with expiration
func SetRedisKey(key string, value interface{}, expiration time.Duration) error {
	if !IsRedisAvailable() {
		return nil
	}
	return RedisClient.Set(ctx, key, value, expiration).Err()
}

// GetRedisKey gets a key value
func GetRedisKey(key string) (string, error) {
	if !IsRedisAvailable() {
		return "", redis.Nil
	}
	return RedisClient.Get(ctx, key).Result()
}

// DeleteRedisKey deletes a key
func DeleteRedisKey(key string) error {
	if !IsRedisAvailable() {
		return nil
	}
	return RedisClient.Del(ctx, key).Err()
}
