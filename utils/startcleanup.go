package utils

import (
    "fmt"
    "time"
)

func StartAutoCleanup() {
    fmt.Println("🚀 Iniciando sistema de cleanup automático (12h)")
    ticker := time.NewTicker(500 * time.Minute)
    go func() {
        for range ticker.C {
            RunSimpleCleanup()
        }
    }()
    
    fmt.Println("⏰ Cleanup agendado para rodar a cada 12 horas")
}