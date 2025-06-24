package handlers

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"github.com/Arthur-Scaratti/yt-api/utils"
	"github.com/gin-gonic/gin"
)

func PlaylistHandler(c *gin.Context) {
    id := c.Query("id")
    index := c.Query("index")
    zipRequested := index == "" || strings.ToUpper(index) == "N"
    dir := filepath.Join(cfg.DownloadDir, id)
    
    files, err := os.ReadDir(dir)
    if err != nil || len(files) == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "ID inválido ou nenhum arquivo encontrado"})
        return
    }
	utils.UpdateLastAccess(id)
    if !zipRequested {
        prefix := fmt.Sprintf("%s -", index)
        var matched os.DirEntry
        for _, file := range files {
            if strings.HasPrefix(file.Name(), prefix) {
                matched = file
                break
            }
        }
        if matched == nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "Arquivo com índice não encontrado"})
            return
        }
        filepath := filepath.Join(dir, matched.Name())
        safeName := utils.SanitizeFilename(matched.Name())
        c.FileAttachment(filepath, safeName)
        return
    }
    
    zipPath := filepath.Join(dir, "playlist.zip")
    if _, err := os.Stat(zipPath); os.IsNotExist(err) {
        err := createZip(zipPath, dir, files)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao criar ZIP"})
            return
        }
    }
    c.FileAttachment(zipPath, "playlist.zip")
}
/////////////////////////////////////////////////////////////

func createZip(zipPath string, dir string, files []os.DirEntry) error {
		zipFile, err := os.Create(zipPath)
		if err != nil {
			return err
		}
		defer zipFile.Close()
	
		zipWriter := zip.NewWriter(zipFile)
		defer zipWriter.Close()
	
		for _, file := range files {
			if file.Name() == "playlist.zip" {
				continue
			}
			filePath := filepath.Join(dir, file.Name())
			f, err := os.Open(filePath)
			if err != nil {
				continue
			}
			defer f.Close()
	
			w, err := zipWriter.Create(file.Name())
			if err != nil {
				continue
			}
			_, err = io.Copy(w, f)
			if err != nil {
				continue
			}
		}
		return nil
	}