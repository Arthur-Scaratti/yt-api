package ytapi

import (
    "context"
    "log"
    "net/http"
    "time"
    "github.com/Arthur-Scaratti/yt-api/handlers"
    "github.com/Arthur-Scaratti/yt-api/utils"

    "github.com/gin-gonic/gin"
)


type ServerConfig struct {
    Port         string
    Debug        bool
    ReadTimeout  time.Duration
    WriteTimeout time.Duration
}


func DefaultServerConfig() *ServerConfig {
    return &ServerConfig{
        Port:         ":8080",
        Debug:        false,
        ReadTimeout:  30 * time.Second,
        WriteTimeout: 30 * time.Second,
    }
}


func InitializeServer(config *ServerConfig) *gin.Engine {
    utils.InitRedis()
    if !config.Debug {
        gin.SetMode(gin.ReleaseMode)
    }
    r := gin.Default()
    SetupRoutes(r)
    return r
}


func SetupRoutes(r *gin.Engine) {
    r.GET("/download", handlers.DownloadHandler)
    r.GET("/result", handlers.ResultHandler)
    r.GET("/health", healthCheckHandler)
}


func StartServer(config *ServerConfig) error {
    r := InitializeServer(config)

    srv := &http.Server{
        Addr:         config.Port,
        Handler:      r,
        ReadTimeout:  config.ReadTimeout,
        WriteTimeout: config.WriteTimeout,
    }
    log.Printf("Servidor iniciando na porta %s", config.Port)
    return srv.ListenAndServe()
}


func StartServerWithContext(ctx context.Context, config *ServerConfig) error {
    r := InitializeServer(config)

    srv := &http.Server{
        Addr:         config.Port,
        Handler:      r,
        ReadTimeout:  config.ReadTimeout,
        WriteTimeout: config.WriteTimeout,
    }

    serverErr := make(chan error, 1)

    go func() {
        log.Printf("Servidor iniciando na porta %s", config.Port)
        serverErr <- srv.ListenAndServe()
    }()

    select {
    case err := <-serverErr:
        return err
    case <-ctx.Done():
        log.Println("Iniciando graceful shutdown...")

        shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()
        
        return srv.Shutdown(shutdownCtx)
    }
}

func healthCheckHandler(c *gin.Context) {
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

func GetServerInfo() map[string]interface{} {
    return map[string]interface{}{
        "name":    "YT MP3 Downloader API",
        "version": "1.0.0",
        "status":  "running",
    }
}