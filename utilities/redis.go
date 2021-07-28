package utilities

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var ctx = context.Background()

var Nil = redis.Nil

func ConnectRedis() error {
	redisHost := os.Getenv("REDIS_HOST")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisPort := os.Getenv("REDIS_PORT")

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: redisPassword,
		DB:       0,
	})

	_, pingError := RedisClient.Ping(ctx).Result()
	if pingError != nil {
		return pingError
	}

	fmt.Println("-- redis: connected")

	return nil
}
