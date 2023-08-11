package service

import (
	"context"
	"log"
	"simple-ddd-cqrs/config"
	"strconv"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

// GetRedisClient
func GetRedisClient() *redis.Client {
	if redisClient != nil {
		return redisClient
	}

	return initRedisClient()
}

// initRedisClient
func initRedisClient() *redis.Client {
	redisDB, _ := strconv.Atoi(config.Get("REDIS_DB"))

	// Inisialisasi Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.Get("REDIS_ADDR"), // localhost:6379
		Password: config.Get("REDIS_PASSWORD"),
		DB:       redisDB,
	})

	ctx := context.Background()
	if err := redisClient.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	return redisClient
}
