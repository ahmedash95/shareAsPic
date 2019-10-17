package main

import (
	"fmt"
	"log"
	"github.com/go-redis/redis/v7"
)

var redisClient *redis.Client

func initRedisClient() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := redisClient.Ping().Result()
	if err != nil {
		Logger.Error(fmt.Sprintf("Faild to connect to redis, %s", err.Error()))
		log.Fatal(err)
	}
}