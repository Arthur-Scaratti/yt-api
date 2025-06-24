package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Arthur-Scaratti/yt-api/config"
)

var cfg *config.Config

func init() {
    cfg = config.Load()
}
// Adicionar essas funções no arquivo download.go

// Verifica se um ID já existe (pasta existe no diretório de downloads)
func CheckExistingID(id string) bool {
    dir := filepath.Join(cfg.DownloadDir, id)
    if _, err := os.Stat(dir); os.IsNotExist(err) {
        return false
    }
	UpdateLastAccess(id)
    return true
}

// Retorna lista de arquivos organizados para playlist
func GetPlaylistFiles(id string) ([]map[string]string, error) {
    dir := filepath.Join(cfg.DownloadDir, id)
    files, err := os.ReadDir(dir)
    if err != nil {
        return nil, err
    }
    
    var fileList []map[string]string
    for _, file := range files {
        if file.Name() == "playlist.zip" {
            continue // ignora o zip
        }
        
        // Extrai índice e título do nome do arquivo
        fileName := file.Name()
        title := strings.TrimSuffix(fileName, filepath.Ext(fileName))
        
        // Se tem formato "N - Título", separa
        var index, cleanTitle string
        if strings.Contains(title, " - ") {
            parts := strings.SplitN(title, " - ", 2)
            index = parts[0]
            cleanTitle = parts[1]
        } else {
            cleanTitle = title
        }
        
        fileList = append(fileList, map[string]string{
            "index":    index,
            "title":    cleanTitle,
            "filename": fileName,
        })
    }
    
    return fileList, nil
}

// Retorna o primeiro arquivo encontrado para download único
func GetSingleFile(id string) (string, error) {
    dir := filepath.Join(cfg.DownloadDir, id)
    files, err := os.ReadDir(dir)
    if err != nil || len(files) == 0 {
        return "", fmt.Errorf("nenhum arquivo encontrado")
    }
    
    // Retorna o primeiro arquivo (ignora .zip se existir)
    for _, file := range files {
        if file.Name() != "playlist.zip" {
            return filepath.Join(dir, file.Name()), nil
        }
    }
    
    return "", fmt.Errorf("nenhum arquivo válido encontrado")
}