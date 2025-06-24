package main

import (
   	"fmt"
    "github.com/Arthur-Scaratti/yt-api/config"
	"github.com/Arthur-Scaratti/yt-api/handlers"
    utils "github.com/Arthur-Scaratti/yt-api/utils"
	"github.com/gin-gonic/gin"
)

func main() {
    cfg := config.Load()
    
    if cfg.GinMode == "debug" {
        gin.SetMode(gin.DebugMode)
    } else {
        gin.SetMode(gin.ReleaseMode)
    }
    
    r := gin.Default()
    
    // Configurar rotas usando variáveis de ambiente
    r.GET(cfg.DownloadHandler, handlers.DownloadHandler)
    r.GET(cfg.PlaylistHandler, handlers.PlaylistHandler)
    r.GET(cfg.WebSocketHandler, handlers.WebSocketHandler)
    
    utils.StartAutoCleanup()
    utils.RunSimpleCleanup()
    // Usar host e porta das variáveis de ambiente
    address := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
    r.Run(address)
}
