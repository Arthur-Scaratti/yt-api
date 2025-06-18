package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/Arthur-Scaratti/yt-api/utils"
)

func HealthCheckHandler(c *gin.Context) {
    if err := utils.RedisClient.Ping(utils.Ctx).Err(); err != nil {
        c.JSON(http.StatusServiceUnavailable, gin.H{
            "status": "error",
            "reason": "Redis connection failed",
            "error":  err.Error(),
        })
        return
    }
    c.JSON(http.StatusOK, gin.H{"status": "ok", "redis": "connected"})
}