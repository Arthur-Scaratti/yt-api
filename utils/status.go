package utils

import (
    "log"
    "os"
    "path/filepath"
    "time"
)

type StatusInfo struct {
    Status                string
    ShouldStartProcessing bool
    CurrentStatusForClient string
}

func GetDownloadStatus(id string) StatusInfo {
    redisStatus, redisErr := RedisClient.Get(Ctx, "media:"+id).Result()
    
    info := StatusInfo{}

    if redisErr == nil {
        log.Printf("Status encontrado no Redis para ID %s: %s", id, redisStatus)
        info.Status = redisStatus
        info.CurrentStatusForClient = redisStatus

        switch redisStatus {
        case "completed":
            updateDownloadTimestamp(id)
            info.ShouldStartProcessing = false
        case "processing":
            info.ShouldStartProcessing = false
        default:
            log.Printf("Status no Redis para ID %s é '%s'. Iniciando novo processamento.", id, redisStatus)
            info.ShouldStartProcessing = true
            info.CurrentStatusForClient = "processing"
        }
    } else {
        log.Printf("Nenhum status no Redis para ID %s (Erro: %v). Iniciando novo processamento.", id, redisErr)
        info.ShouldStartProcessing = true
        info.CurrentStatusForClient = "processing"
    }

    return info
}

func SetProcessingStatus(id string) {
    RedisClient.Set(Ctx, "media:"+id, "processing", 0)
}

func updateDownloadTimestamp(id string) {
    downloadDir := filepath.Join("downloads", id)
    if _, statErr := os.Stat(downloadDir); !os.IsNotExist(statErr) {
        os.Chtimes(downloadDir, time.Now(), time.Now())
    }
}