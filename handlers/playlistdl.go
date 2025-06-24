package handlers

import (
    "bufio"
    "fmt"
    "io"
    "os/exec"
    "path/filepath"
    "strconv"
    "strings"
    "github.com/Arthur-Scaratti/yt-api/utils"
)

func RunPlaylistDownload(videoURL, format, quality, id, dir string) {
    outputname := cfg.OutputTemplatePlaylist

    cmdArgs := []string{
        "--concurrent-fragments", strconv.Itoa(cfg.ConcurrentFragments),
        "--fragment-retries", strconv.Itoa(cfg.FragmentRetries),
        "--retries", strconv.Itoa(cfg.Retries),
        "--extractor-retries", strconv.Itoa(cfg.ExtractorRetries),
        "--progress-template", cfg.ProgressTemplate,
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
    cmdArgs = append(cmdArgs, videoURL)

    cmd := exec.Command("yt-dlp", cmdArgs...)
    stdout, _ := cmd.StdoutPipe()
    stderr, _ := cmd.StderrPipe()
    _ = cmd.Start()
    go io.Copy(io.Discard, stderr)

    scanner := bufio.NewScanner(stdout)
    var CURRENT_NAME string // Variável para armazenar o nome atual

    for scanner.Scan() {
        line := scanner.Text()
        
        // Atualizar CURRENT_NAME baseado no formato
        if format == "mp3" && strings.Contains(line, "[ExtractAudio] Destination:") {
            // Para MP3: pegar o nome após "[ExtractAudio] Destination:"
            parts := strings.SplitN(line, "[ExtractAudio] Destination: ", 2)
            if len(parts) == 2 {
                fileName := strings.TrimSpace(parts[1])
                if strings.HasSuffix(fileName, ".mp3") {
                    CURRENT_NAME = filepath.Base(fileName)
                }
            }
        } else if format != "mp3" && strings.Contains(line, "[Merger] Merging formats into") {
            // Para vídeos: pegar o nome após "[Merger] Merging formats into"
            parts := strings.SplitN(line, "[Merger] Merging formats into", 2)
            if len(parts) == 2 {
                fileName := strings.TrimSpace(parts[1])
                fileName = strings.Trim(fileName, "\"")
                CURRENT_NAME = filepath.Base(fileName)
            }
        }
        
        // Fazer broadcast quando detectar "[download] Downloading item"
        if strings.Contains(line, "[download] Downloading item") && CURRENT_NAME != "" {
            title := strings.TrimSuffix(CURRENT_NAME, filepath.Ext(CURRENT_NAME))
            title = utils.SanitizeFilename(title)
            broadcastItem(id, title)
            fmt.Printf("✔ Broadcast enviado: %s\n", title)
            // Reset CURRENT_NAME após broadcast
            CURRENT_NAME = ""
        }
   
    }
    
    cmd.Wait()
    broadcastItem(id, "completed")
    closeWebSocketConnections(id)
}