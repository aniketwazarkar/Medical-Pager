package redis

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
	"medical-pager/utils"
)

var Client *redis.Client
var Ctx = context.Background()

// Connect establishes connection to Redis
func Connect() {
	url := utils.GetEnv("REDIS_URL", "redis://localhost:6379/0")
	opt, err := redis.ParseURL(url)
	if err != nil {
		log.Fatalf("Failed to parse Redis URL: %v", err)
	}

	client := redis.NewClient(opt)

	_, err = client.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Successfully connected to Redis!")
	Client = client
}
