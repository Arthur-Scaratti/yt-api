package utils

import (
    "log"
    "os"
    "os/exec"
    "path/filepath"
)

func ProcessDownload(url, format, id string) {
    dir := filepath.Join("downloads", id)
    
    if err := os.MkdirAll(dir, 0755); err != nil {
        log.Printf("Erro ao criar diretório %s: %v", dir, err)
        RedisClient.Set(Ctx, "media:"+id, "error", 0)
        return
    }

    cmd := buildDownloadCommand(url, format, dir)
    
    log.Printf("Iniciando download para ID %s: %s", id, cmd.String())
    
    if err := cmd.Run(); err != nil {
        log.Printf("Erro ao executar yt-dlp para ID %s: %v", id, err)
        RedisClient.Set(Ctx, "media:"+id, "error", 0)
        return
    }

    log.Printf("Download concluído para ID %s. Definindo status no Redis.", id)
    RedisClient.Set(Ctx, "media:"+id, "completed", 0)
    CleanupDownloads(2048) // 2GB
}

func buildDownloadCommand(url, format, dir string) *exec.Cmd {
    output := "%(title)s.%(ext)s"
    
    if format == "mp3" {
        return exec.Command("yt-dlp", "-x", "--audio-format", "mp3", "-o", output, "-P", dir, url)
    }
    return exec.Command("yt-dlp", "-o", output, "-P", dir, url)
}