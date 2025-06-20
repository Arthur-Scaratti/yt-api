package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    ytapi "github.com/Arthur-Scaratti/yt-api"
)

func createDownloadsDir() error {
    downloadDir := "downloads"
    
    // Verifica se a pasta já existe
    if _, err := os.Stat(downloadDir); os.IsNotExist(err) {
        // Cria a pasta com permissões 0755 (rwxr-xr-x)
        if err := os.Mkdir(downloadDir, 0755); err != nil {
            return err
        }
        log.Println("Pasta 'downloads' criada com sucesso")
    } else {
        log.Println("Pasta 'downloads' já existe")
    }
    
    return nil
}

func main() {
    
    if err := createDownloadsDir(); err != nil {
        log.Fatal("Erro ao criar pasta downloads:", err)
    }

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