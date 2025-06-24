package handlers

import (
    "crypto/sha256"
    "fmt"
    "net/http"
    "os"
    "os/exec"
    "path/filepath"
    "strconv"
    "strings"
    "github.com/Arthur-Scaratti/yt-api/utils"
    "github.com/Arthur-Scaratti/yt-api/config"
    "github.com/gin-gonic/gin"
)

var cfg *config.Config

func init() {
    cfg = config.Load()
}

func createDownloadDir() error {
    if _, err := os.Stat(cfg.DownloadDir); os.IsNotExist(err) {
        return os.MkdirAll(cfg.DownloadDir, cfg.FilePermissions)
    }
    return nil
}

func DownloadHandler(c *gin.Context) {
    createDownloadDir()

    videoURL := c.Query("url")
    format := c.DefaultQuery("format", cfg.DefaultFormat)
    quality := c.DefaultQuery("quality", fmt.Sprintf(cfg.DefaultQuality, "p"))
    playlist := c.DefaultQuery("playlist", cfg.DefaultPlaylist)
    index := c.DefaultQuery("index", cfg.DefaultIndex)

    
    if videoURL == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Missing URL"})
        return
    }
    
    isPlaylist := strings.ToLower(playlist) == "true"
    isIndexSet := index != ""

    hasher := sha256.New()
    inputString := fmt.Sprintf("%s|%s|%s|%s|%s", videoURL, format, quality, playlist, index)
    hasher.Write([]byte(inputString))
    id := fmt.Sprintf("dl_%x", hasher.Sum(nil))

	    // VERIFICAÇÃO SE ID JÁ EXISTE
		if utils.CheckExistingID(id) {
			// ID já existe, retornar arquivo/informações
			if isPlaylist && !isIndexSet {
				// Playlist completa - retornar lista organizada
				fileList, err := utils.GetPlaylistFiles(id)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao ler arquivos da playlist"})
					return
				}
				
				c.JSON(http.StatusOK, gin.H{
					"status":   "Ready",
					"id":       id,
					"count":    len(fileList),
					"files":    fileList,
					"download": fmt.Sprintf("%s?id=%s&index=N", cfg.PlaylistHandler,id),
				})
				return
			} else {
				// Download único ou item específico da playlist - retornar arquivo
				filePath, err := utils.GetSingleFile(id)
				if err != nil {
					c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
					return
				}
				
				fileName := filepath.Base(filePath)
				safeName := utils.SanitizeFilename(fileName)
				c.FileAttachment(filePath, safeName)
				return
			}
		}

	dir := filepath.Join(cfg.DownloadDir, id)
    os.MkdirAll(dir, os.ModePerm)

    ///////////////Retorna imediatamente e roda em background///////////////
    if isPlaylist && !isIndexSet {
        progressURL := fmt.Sprintf("%s?id=%s", cfg.WebSocketHandler,id)
        c.JSON(http.StatusAccepted, gin.H{
            "id":          id,
            "progressUrl": progressURL,
        })
        go RunPlaylistDownload(videoURL, format, quality, id, dir)
        return
    }

////////////// Execução normal (index ou não-playlist)////////////////////////////////////
    outputname := cfg.OutputTemplateSingle

    cmdArgs := []string{
        "--concurrent-fragments", strconv.Itoa(cfg.ConcurrentFragments),
        "--fragment-retries", strconv.Itoa(cfg.FragmentRetries),
        "--retries", strconv.Itoa(cfg.Retries),
        "--extractor-retries", strconv.Itoa(cfg.ExtractorRetries),
        "-o", outputname,
        "-P", dir,
    }

    formatSelector := BuildFormatSelector(format, quality)
    cmdArgs = append(cmdArgs, "-f", formatSelector)

    switch format {
case "mp3":
        cmdArgs = append(cmdArgs, "--extract-audio", "--audio-format", format)
    case "mp4", "mkv", "webm":
        cmdArgs = append(cmdArgs, "--merge-output-format", format)
    }

    if isPlaylist {
        if isIndexSet {
            cmdArgs = append(cmdArgs, "--playlist-items", index)
        }
    } else {
        cmdArgs = append(cmdArgs, "--no-playlist")
    }

    cmdArgs = append(cmdArgs, videoURL)

    cmd := exec.Command("yt-dlp", cmdArgs...)
    output, err := cmd.CombinedOutput()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Download failed", "details": string(output)})
        return
    }

    files, _ := os.ReadDir(dir)
    if isPlaylist && !isIndexSet {
        var names []string
        for _, f := range files {
            names = append(names, f.Name())
        }
        c.JSON(http.StatusOK, gin.H{
            "id":       id,
            "count":    len(names),
            "files":    names,
            "download": fmt.Sprintf("/result?id=%s&index=N", id),
        })
        return
    }
    
    for _, f := range files {
        safeName := utils.SanitizeFilename(f.Name())
        c.FileAttachment(filepath.Join(dir, f.Name()), safeName)
		utils.UpdateLastAccess(id)
        return
    }

    c.JSON(http.StatusInternalServerError, gin.H{"error": "No file found"})
}

func BuildFormatSelector(format string, quality string) string {
    switch format {
    case "mp3":
        return "bestaudio[ext=mp3]/bestaudio"
    case "mp4", "mkv", "webm":
        height := ParseQuality(quality)
        selector := fmt.Sprintf("bestvideo[height<=%d]+bestaudio/best", height)
        selector += fmt.Sprintf("[ext=%s]", format)
        return selector
    default:
        return "best"
    }
}

func ParseQuality(quality string) int {
    if strings.HasSuffix(quality, "p") {
        q := strings.TrimSuffix(quality, "p")
        if h, err := strconv.Atoi(q); err == nil {
            return h
        }
    }
    return cfg.DefaultQualityYTDLP
}