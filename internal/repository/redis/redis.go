package redis

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func InitRedis() (*redis.Client) {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Couldn't load .env variables %v\n", err)
	}
	redisAddr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	var redisClient *redis.Client
	redisClient = redis.NewClient(&redis.Options{
		Addr: redisAddr,
		Password: "",
		DB: 0,
	})
	return redisClient
}