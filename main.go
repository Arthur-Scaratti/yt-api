package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    "yt-mp3-api/serverlib"
)

func main() {
    config := serverlib.DefaultServerConfig()
    
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

    if err := serverlib.StartServerWithContext(ctx, config); err != nil {
        log.Fatalf("Falha ao iniciar o servidor: %v", err)
    }
    log.Println("Servidor finalizado")
}