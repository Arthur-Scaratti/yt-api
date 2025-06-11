package utils

import (
    "log"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}

type WebSocketMessage struct {
    Status  string `json:"status"`
    ID      string `json:"id"`
    Message string `json:"message,omitempty"`
}

func CreateWebSocketMessages(status, id string) WebSocketMessage {
    msg := WebSocketMessage{Status: status, ID: id}
    
    switch status {
    case "completed":
        msg.Message = "Download concluído. Use /result?id=" + id
    case "error":
        msg.Message = "Erro durante o processamento."
    case "processing":
        msg.Message = "Processamento em andamento."
    }
    
    return msg
}

func MonitorStatusAndNotify(ws *websocket.Conn, id, initialStatus string, ctx *gin.Context) <-chan struct{} {
    done := make(chan struct{})
    
    go func() {
        defer close(done)
        defer log.Printf("Goroutine de monitoramento para %s finalizada.", id)

        ticker := time.NewTicker(2 * time.Second)
        defer ticker.Stop()

        lastKnownStatus := initialStatus

        for {
            select {
            case <-ticker.C:
                status, err := RedisClient.Get(Ctx, "media:"+id).Result()
                if err != nil {
                    log.Printf("Erro no Redis ao consultar status para %s: %v. Continuando.", id, err)
                    continue
                }

                if status != lastKnownStatus || status == "completed" || status == "error" {
                    log.Printf("Status consultado para %s: %s. Enviando atualização.", id, status)
                    
                    msg := CreateWebSocketMessages(status, id)
                    if status == "processing" && lastKnownStatus == "processing" {
                        msg.Message = "" // Evita mensagem repetitiva
                    }

                    if err := ws.WriteJSON(msg); err != nil {
                        log.Printf("Erro ao escrever mensagem no WebSocket para %s: %v. Cliente provavelmente desconectado.", id, err)
                        return
                    }

                    lastKnownStatus = status

                    if status == "completed" || status == "error" {
                        log.Printf("Status final %s para %s. Parando monitoramento.", status, id)
                        return
                    }
                }
            case <-ctx.Request.Context().Done():
                log.Printf("Cliente desconectado (contexto finalizado) para %s. Parando monitoramento.", id)
                return
            }
        }
    }()
    
    return done
}

func KeepWebSocketAlive(ws *websocket.Conn, id string) {
    for {
        if _, _, err := ws.ReadMessage(); err != nil {
            log.Printf("Erro de leitura do cliente %s (provavelmente desconectado): %v", id, err)
            break
        }
    }
}