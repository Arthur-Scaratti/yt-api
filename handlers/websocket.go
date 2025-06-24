package handlers

import (
    "log"
    "net/http"
    "sync"

    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
    "slices"
)

var upgrader = websocket.Upgrader{
    CheckOrigin: func(r *http.Request) bool {
        return true
    },
}
var (
    wsClients = make(map[string][]*websocket.Conn)
    wsMutex   sync.RWMutex
)

func WebSocketHandler(c *gin.Context) {
    id := c.Query("id")
    if id == "" {
        return
    }

    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        log.Printf("Erro no upgrade do websocket: %v", err)
        return
    }
    wsMutex.Lock()
    wsClients[id] = append(wsClients[id], conn)
    wsMutex.Unlock()

    defer func() {
        wsMutex.Lock()
        defer wsMutex.Unlock()

        connections := wsClients[id]
        for i, c := range connections {
            if c == conn {
                wsClients[id] = slices.Delete(connections, i, i+1)
                break
            }
        }
        if len(wsClients[id]) == 0 {
            delete(wsClients, id)
        }
        conn.Close()
        log.Printf("Cliente desconectado e removido para o id: %s", id)
    }()
    for {
        if _, _, err := conn.NextReader(); err != nil {
            break
        }
    }
}

func broadcastItem(id, title string) {
    wsMutex.RLock()
    connections := wsClients[id]
    wsMutex.RUnlock()
    
    for _, conn := range connections {
        conn.WriteJSON(gin.H{
            "id":    id,
            "title": title,
        })
    }
}

func closeWebSocketConnections(id string) {
    wsMutex.Lock()
    defer wsMutex.Unlock()
    
    connections := wsClients[id]
    for _, conn := range connections {
        conn.Close()
    }
    // Remove todas as conexões do mapa
    delete(wsClients, id)
    log.Printf("Todas as conexões WebSocket fechadas para o id: %s", id)
}
