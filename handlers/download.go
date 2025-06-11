package handlers

import (
    "crypto/sha1"
    "encoding/hex"
    "log"
    "net/http"

    "github.com/gorilla/websocket"
    "github.com/gin-gonic/gin"
    
    "yt-mp3-api/utils"
)

func DownloadHandler(c *gin.Context) {
    url := c.Query("url")
    format := c.Query("format")

    if !isValidRequest(url, format) {
        c.JSON(http.StatusBadRequest, gin.H{
            "error":   "Parâmetros inválidos",
            "message": "Use 'url' e 'format' (mp3 ou mp4). Ex: /download?url=...&format=mp3",
        })
        return
    }

    ws, err := utils.Upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        log.Printf("Falha ao atualizar para websocket: %v", err)
        return
    }
    defer ws.Close()

    id := generateDownloadID(url, format)
    log.Printf("Requisição de download recebida para ID %s (URL: %s, Formato: %s)", id, url, format)

    statusInfo := utils.GetDownloadStatus(id)
    
    if err := handleInitialStatus(ws, id, statusInfo); err != nil {
        return
    }

    if statusInfo.ShouldStartProcessing {
        utils.SetProcessingStatus(id)
        go utils.ProcessDownload(url, format, id)
    }

    monitorDone := utils.MonitorStatusAndNotify(ws, id, statusInfo.CurrentStatusForClient, c)
    utils.KeepWebSocketAlive(ws, id)
    
    <-monitorDone
    log.Printf("Handler de Download para %s finalizado.", id)
}

func isValidRequest(url, format string) bool {
    return url != "" && (format == "mp3" || format == "mp4")
}

func generateDownloadID(url, format string) string {
    hash := sha1.New()
    hash.Write([]byte(url + format))
    return hex.EncodeToString(hash.Sum(nil))
}

func handleInitialStatus(ws *websocket.Conn, id string, statusInfo utils.StatusInfo) error {
    var msg utils.WebSocketMessage
    
    switch statusInfo.Status {
    case "completed":
        msg = utils.WebSocketMessage{
            Status:  "completed",
            ID:      id,
            Message: "Conteúdo já processado e pronto. Use /result?id=" + id,
        }
    case "processing":
        msg = utils.WebSocketMessage{
            Status:  "processing",
            ID:      id,
            Message: "Processamento já em andamento.",
        }
    default:
        if statusInfo.ShouldStartProcessing {
            msg = utils.WebSocketMessage{
                Status:  "processing",
                ID:      id,
                Message: "Iniciando processamento.",
            }
        }
    }

    if err := ws.WriteJSON(msg); err != nil {
        log.Printf("Erro ao enviar status inicial para %s: %v", id, err)
        return err
    }
    
    return nil
}