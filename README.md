# YT API - YouTube Downloader Package

## 📋 Visão Geral

Um package Go para download de vídeos do YouTube com suporte a conversão para MP3 e MP4. O projeto utiliza WebSockets para fornecer atualizações de status em tempo real e pode ser usado tanto como biblioteca importável quanto como servidor standalone.

## 🏗️ Arquitetura

### Stack Tecnológica

- **Backend**: Go (Gin Framework)
- **Cache/Status**: Redis
- **WebSockets**: Gorilla WebSocket
- **Download Engine**: yt-dlp
- **Compressão**: Archive/zip nativo

### Estrutura do Projeto

```
├── server.go              # Funções principais exportáveis
├── cmd/server/main.go     # Exemplo de servidor standalone
├── handlers/              # Handlers HTTP/WebSocket
│   ├── download.go        # Handler de download via WebSocket
│   ├── result.go          # Handler para servir arquivos
│   └── redishealth.go     # Health check
└── utils/                 # Utilitários (Redis, WebSocket, etc.)
```

## 🚀 Instalação e Uso

### Pré-requisitos

```bash
# Instalar yt-dlp
pip install yt-dlp

# Instalar Redis
# Ubuntu/Debian:
sudo apt install redis-server

# macOS:
brew install redis

# Windows: Download do site oficial
```

### Como Package Importável

#### 1. Instalação

```bash
go get github.com/Arthur-Scaratti/yt-api
```

#### 2. Uso Básico (QuickStart)

```go
package main

import (
    "log"
    ytapi "github.com/Arthur-Scaratti/yt-api"
)

func main() {
    // Inicia servidor com configurações padrão
    if err := ytapi.QuickStart(); err != nil {
        log.Fatal("Erro ao iniciar servidor:", err)
    }
}
```

#### 3. Uso com Context (Graceful Shutdown)

```go
package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"
    ytapi "github.com/Arthur-Scaratti/yt-api"
)

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Captura sinais para shutdown graceful
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        <-sigChan
        log.Println("Sinal recebido, iniciando shutdown...")
        cancel()
    }()

    // Inicia servidor com context
    if err := ytapi.QuickStartWithContext(ctx); err != nil {
        log.Printf("Servidor finalizado: %v", err)
    }
}
```

#### 4. Configuração Customizada

```go
package main

import (
    "time"
    ytapi "github.com/Arthur-Scaratti/yt-api"
)

func main() {
    // Configuração do servidor
    serverConfig := &ytapi.ServerConfig{
        Port:         ":3000",
        Debug:        true,
        ReadTimeout:  30 * time.Second,
        WriteTimeout: 30 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    // Configuração do Redis
    redisConfig := &ytapi.RedisConfig{
        Address:  "localhost:6379",
        Password: "",
        DB:       0,
    }

    // Inicia com configurações customizadas
    if err := ytapi.Start(serverConfig, redisConfig); err != nil {
        log.Fatal("Erro ao iniciar servidor:", err)
    }
}
```

### Como Servidor Standalone

```bash
# Clone o repositório
git clone https://github.com/Arthur-Scaratti/yt-api.git
cd yt-api

# Instale dependências
go mod tidy

# Execute o servidor
go run cmd/server/main.go

# Ou compile e execute
go build -o ytapi cmd/server/main.go
./ytapi
```

#### Variáveis de Ambiente

```bash
# Porta do servidor (padrão: 8080)
export PORT=3000

# Modo debug (padrão: false)
export DEBUG=true

# Execute o servidor
go run cmd/server/main.go
```

## 📡 API Endpoints

### 1. `/download` - WebSocket para Download

**Método**: `GET` (upgrade para WebSocket)

**Parâmetros**:
- `url`: URL do vídeo do YouTube
- `format`: `mp3` ou `mp4`

**Exemplo**:
```
ws://localhost:8080/download?url=https://youtube.com/watch?v=VIDEO_ID&format=mp3
```

**Fluxo WebSocket**:
1. Cliente conecta via WebSocket
2. Servidor valida parâmetros e inicia download
3. Servidor envia atualizações de status em tempo real
4. Cliente recebe ID do download quando concluído

**Mensagens WebSocket**:
```json
{
  "status": "processing|completed|error",
  "id": "hash_do_download",
  "message": "Descrição do status atual"
}
```

### 2. `/result` - Obter Arquivo Processado

**Método**: `GET`

**Parâmetros**:
- `id`: ID do download (recebido via WebSocket)

**Exemplo**:
```
GET /result?id=abc123def456
```

**Respostas**:
- **200**: Download do arquivo
- **202**: Ainda processando
- **404**: Arquivo não encontrado

### 3. `/health` - Health Check

**Método**: `GET`

**Exemplo**:
```
GET /health
```

**Resposta**:
```json
{
  "status": "ok",
  "redis": "connected"
}
```

## 🧪 Testando Localmente

### 1. Teste Rápido com cURL

```bash
# Health check
curl http://localhost:8080/health

# Verificar se um download existe
curl "http://localhost:8080/result?id=exemplo123"
```

### 2. Teste WebSocket com JavaScript

```javascript
const ws = new WebSocket(`ws://localhost:8080/download?url=${encodeURIComponent('https://youtube.com/watch?v=VIDEO_ID')}&format=mp3`);

ws.onmessage = function(event) {
  const data = JSON.parse(event.data);
  console.log('Status:', data.status, '-', data.message);
  
  if (data.status === 'completed') {
    // Download do arquivo
    window.location.href = `/result?id=${data.id}`;
    ws.close();
  }
};

ws.onerror = function(error) {
  console.error('WebSocket error:', error);
};
```

### 3. Teste WebSocket com Python

```python
import websocket
import json
import requests

def on_message(ws, message):
    data = json.loads(message)
    print(f"Status: {data['status']} - {data.get('message', '')}")
    
    if data['status'] == 'completed':
        # Download do arquivo
        response = requests.get(f"http://localhost:8080/result?id={data['id']}")
        filename = f"download_{data['id']}.mp3"
        with open(filename, 'wb') as f:
            f.write(response.content)
        print(f"Arquivo salvo como: {filename}")
        ws.close()

def on_error(ws, error):
    print(f"Erro: {error}")

# Conectar ao WebSocket
url = "ws://localhost:8080/download?url=https://youtube.com/watch?v=VIDEO_ID&format=mp3"
ws = websocket.WebSocketApp(url, on_message=on_message, on_error=on_error)
ws.run_forever()
```

## 🔧 Funções Exportadas

### Configurações

```go
// Configurações padrão
serverConfig := ytapi.DefaultServerConfig()
redisConfig := ytapi.DefaultRedisConfig()

// Informações do servidor
info := ytapi.GetServerInfo()
```

### Inicialização

```go
// Início rápido
ytapi.QuickStart()
ytapi.QuickStartWithContext(ctx)

// Início customizado
ytapi.Start(serverConfig, redisConfig)
ytapi.StartWithContext(ctx, serverConfig, redisConfig)

// Apenas inicializar (sem iniciar servidor)
engine := ytapi.InitializeServer(serverConfig, redisConfig)
```

## 🔄 Fluxo de Funcionamento

1. **Conexão WebSocket**: Cliente conecta em `/download` com parâmetros
2. **Validação**: Servidor valida URL e formato
3. **ID Geração**: Hash SHA1 da URL + formato
4. **Status Check**: Verifica se download já existe no Redis
5. **Processamento**: Inicia download com `yt-dlp` se necessário
6. **Atualizações**: Envia status via WebSocket em tempo real
7. **Conclusão**: Cliente recebe ID para buscar arquivo em `/result`

## 📁 Estrutura de Dados

### ServerConfig
```go
type ServerConfig struct {
    Port         string        // Porta do servidor (ex: ":8080")
    Debug        bool          // Modo debug
    ReadTimeout  time.Duration // Timeout de leitura
    WriteTimeout time.Duration // Timeout de escrita
    IdleTimeout  time.Duration // Timeout de idle
}
```

### RedisConfig
```go
type RedisConfig struct {
    Address  string // Endereço do Redis (ex: "localhost:6379")
    Password string // Senha do Redis
    DB       int    // Número do banco Redis
}
```

## 🚨 Requisitos do Sistema

- **Go**: 1.22.4 ou superior
- **Redis**: Qualquer versão recente
- **yt-dlp**: Instalado e acessível via PATH
- **Espaço em disco**: Para armazenamento temporário dos downloads

## 📝 Notas Importantes

- Os arquivos são armazenados temporariamente em `downloads/`
- Cada download é identificado por um hash SHA1 único
- O sistema suporta múltiplos downloads simultâneos
- WebSockets mantêm conexão ativa durante todo o processo
- Arquivos múltiplos são automaticamente compactados em ZIP