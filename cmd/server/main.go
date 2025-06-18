package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    ytapi "github.com/Arthur-Scaratti/yt-api"
)

func main() {
    config := ytapi.DefaultServerConfig()
    
    if port := os.Getenv("PORT"); port != "" {
        config.Port = ":" + port
    }
    
    if os.Getenv("DEBUG") == "true" {
        config.Debug = true
    }

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        <-sigChan
        log.Println("Sinal recebido, iniciando shutdown...")
        cancel()
    }()

    ytapi.QuickStartWithContext(ctx)
    log.Println("Servidor finalizado")
}