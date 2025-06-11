package handlers

import (
    "archive/zip"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "strings"

    "github.com/gin-gonic/gin"
    "yt-mp3-api/utils"
)

func ResultHandler(c *gin.Context) {
    id := c.Query("id")
    if id == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "ID requerido"})
        return
    }

    // Verificar status no Redis
    status, err := utils.RedisClient.Get(utils.Ctx, "media:"+id).Result()
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"status": "not_found"})
        return
    }

    if status != "completed" {
        c.JSON(http.StatusAccepted, gin.H{
            "status":  status,
            "message": "A mídia ainda está sendo processada. Tente novamente em breve.",
        })
        return
    }

    dir := filepath.Join("downloads", id)
    files, err := os.ReadDir(dir)
    if err != nil || len(files) == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "Arquivo não encontrado"})
        return
    }

    if len(files) == 1 {
        path := filepath.Join(dir, files[0].Name())
        c.FileAttachment(path, files[0].Name())
    } else {
        zipPath := filepath.Join("downloads", id+".zip")
        if err := zipFiles(dir, zipPath); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao compactar arquivos"})
            return
        }
        c.FileAttachment(zipPath, id+".zip")
    }
}

func zipFiles(sourceDir, zipPath string) error {
    zipFile, err := os.Create(zipPath)
    if err != nil {
        return err
    }
    defer zipFile.Close()

    zipWriter := zip.NewWriter(zipFile)
    defer zipWriter.Close()

    return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.IsDir() {
            return nil
        }

        relPath := strings.TrimPrefix(path, sourceDir+string(os.PathSeparator))
        file, err := os.Open(path)
        if err != nil {
            return err
        }
        defer file.Close()

        zipEntry, err := zipWriter.Create(relPath)
        if err != nil {
            return err
        }

        _, err = io.Copy(zipEntry, file)
        return err
    })
}
