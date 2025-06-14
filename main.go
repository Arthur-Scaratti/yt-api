package main

import (
	"log"
	"net/http" 
	"yt-mp3-api/handlers"
	"yt-mp3-api/utils" 

	"github.com/gin-gonic/gin"
)

func main() {
    utils.InitRedis()

    r := gin.Default()

    r.POST("/download", handlers.DownloadHandler)
    r.GET("/result", handlers.ResultHandler) 

    r.GET("/health", func(c *gin.Context) {
        if err := utils.RedisClient.Ping(utils.Ctx).Err(); err != nil {
            c.JSON(http.StatusServiceUnavailable, gin.H{
                "status": "error",
                "reason": "Redis connection failed",
                "error":  err.Error(),
            })
            return
        }
        c.JSON(http.StatusOK, gin.H{"status": "ok", "redis": "connected"})
    })

    log.Println("Servidor iniciando na porta :8080")
    if err := r.Run(":8080"); err != nil {
        log.Fatalf("Falha ao iniciar o servidor: %v", err)
    }
}