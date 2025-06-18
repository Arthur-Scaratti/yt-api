package utils

import (
    "context"
    "log"
    "github.com/redis/go-redis/v9"
)

var Ctx = context.Background()
var RedisClient *redis.Client

func InitRedis(address, password string, db int) {
    log.Printf("🔌 Conectando ao Redis em %s (DB: %d)", address, db)
    
    RedisClient = redis.NewClient(&redis.Options{
        Addr:     address,
        Password: password,
        DB:       db,
    })

    // Testa a conexão
    if err := RedisClient.Ping(Ctx).Err(); err != nil {
        log.Printf("❌ Falha ao conectar com Redis: %v", err)
    } else {
        log.Printf("✅ Redis conectado com sucesso")
    }
}

func GetRedisClient() *redis.Client {
    return RedisClient
}