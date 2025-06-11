package utils

import (
    "context"
    "github.com/redis/go-redis/v9"
)

var Ctx = context.Background()
var RedisClient *redis.Client

func InitRedis() {
    RedisClient = redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
        DB:   0,
    })
}
