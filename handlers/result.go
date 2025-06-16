package handlers

import (
	"archive/zip"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Arthur-Scaratti/yt-api/utils"

	"github.com/gin-gonic/gin"
)

func ResultHandler(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID requerido"})
		return
	}

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

func zipFiles(sourceDirectoryPath, destinationZipPath string) error {
	outputZipFile, err := os.Create(destinationZipPath)
	if err != nil {
		return err
	}
	defer outputZipFile.Close()

	zipArchiveWriter := zip.NewWriter(outputZipFile)
	defer zipArchiveWriter.Close()

	return filepath.Walk(sourceDirectoryPath, func(filePath string, fileInfo os.FileInfo, errWalk error) error {
		if errWalk != nil {
			return errWalk
		}
		if fileInfo.IsDir() {
			return nil
		}

		relativePathInZip := strings.TrimPrefix(filePath, sourceDirectoryPath+string(os.PathSeparator))
		fileToZip, errOpen := os.Open(filePath)
		if errOpen != nil {
			return errOpen
		}
		defer fileToZip.Close()

		zipFileEntry, errCreate := zipArchiveWriter.Create(relativePathInZip)
		if errCreate != nil {
			return errCreate
		}

		_, errCopy := io.Copy(zipFileEntry, fileToZip)
		return errCopy
	})
}
