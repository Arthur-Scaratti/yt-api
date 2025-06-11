
# YouTube Downloader API

## 📋 Visão Geral

Uma API REST moderna para download de vídeos do YouTube com suporte a conversão para MP3 e MP4. O projeto utiliza WebSockets para fornecer atualizações de status em tempo real, eliminando a necessidade de polling e melhorando significativamente a experiência do usuário.

## 🏗️ Arquitetura e Decisões Técnicas

### Stack Tecnológica

- **Backend**: Go (Gin Framework)
- **Cache/Status**: Redis
- **WebSockets**: Gorilla WebSocket
- **Download Engine**: yt-dlp
- **Compressão**: Archive/zip nativo

### Principais Decisões Arquiteturais

### 1. **WebSockets vs Polling**

- ✅ **Escolhido**: WebSockets para atualizações em tempo real
- ❌ **Rejeitado**: Polling via endpoints `/status`
- **Razão**: Reduz latência, melhora UX, diminui carga no servidor

### 2. **Cache de Status com Redis**

- Persiste o estado dos downloads entre reinicializações
- Evita reprocessamento desnecessário
- Permite múltiplas conexões para o mesmo download

### 3. **ID Único Baseado em Hash**

- Combinação de URL + formato gera SHA1
- Evita downloads duplicados
- Facilita cache e organização de arquivos

## 🚀 Instalação e Configuração

### Pré-requisitos

```sh
# Instalar yt-dlp
pip install yt-dlp

# Instalar e configurar Redis

# Ubuntu/Debian:
sudo apt install redis-server

# macOS:
brew install redis

# Windows: Download do site oficial
```

### Instalação do Projeto

```sh
# Clone o repositório
git clone <repository-url>
cd yt-downloader

# Instalar dependências Go
go mod tidy

# Instalar dependência WebSocket
go get github.com/gorilla/websocket

# Executar
go run main.go
```

### Variáveis de Ambiente

```sh
# Redis (opcional - padrões localhost:6379)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# Servidor
PORT=8080
```

## 📡 API Endpoints

### 1. WebSocket Download Endpoint

```http
GET /download?url=<youtube_url>&format=<mp3|mp4>
Upgrade: websocket
```

**Parâmetros:**

- `url` (required): URL válida do YouTube
- `format` (required): `mp3` ou `mp4`

**Fluxo WebSocket:**

1. Cliente conecta via WebSocket
2. Servidor retorna status inicial
3. Atualizações automáticas durante processamento
4. Notificação quando concluído

**Mensagens WebSocket:**

```json
// Status inicial/em progresso
{
  "status": "processing",
  "id": "abc123...",
  "message": "Iniciando processamento."
}

// Concluído
{
  "status": "completed",
  "id": "abc123...",
  "message": "Download concluído. Use /result?id=abc123..."
}

// Erro
{
  "status": "error",
  "id": "abc123...",
  "message": "Erro durante o processamento."
}
```

### 2. Result Endpoint

```http
GET /result?id=<download_id>
```

### 3. Health Check

```http
GET /health
```

```json
{
  "status": "ok",
  "redis": "connected"
}
```

## 💻 Como Usar

### Exemplo com JavaScript (Cliente Web)

```javascript
const ws = new WebSocket(`ws://localhost:8080/download?url=${encodeURIComponent(youtubeUrl)}&format=mp3`);

ws.onmessage = function(event) {
  const data = JSON.parse(event.data);
  switch(data.status) {
    case 'processing':
      updateUI('Processando...', data.message);
      break;
    case 'completed':
      updateUI('Concluído!', data.message);
      window.location.href = `/result?id=${data.id}`;
      ws.close();
      break;
    case 'error':
      updateUI('Erro', data.message);
      ws.close();
      break;
  }
};

ws.onerror = function(error) {
  console.error('WebSocket error:', error);
};
```

### Exemplo com Python

```python
import websocket
import json
import requests

def on_message(ws, message):
    data = json.loads(message)
    print(f"Status: {data['status']} - {data.get('message', '')}")
    if data['status'] == 'completed':
        response = requests.get(f"http://localhost:8080/result?id={data['id']}")
        with open(f"download_{data['id']}.mp3", 'wb') as f:
            f.write(response.content)
        ws.close()

def on_error(ws, error):
    print(f"Erro: {error}")

url = "wss://localhost:8080/download?url=https://youtube.com/watch?v=...&format=mp3"
ws = websocket.WebSocketApp(url, on_message=on_message, on_error=on_error)
ws.run_forever()
```

## 📁 Estrutura do Projeto

```
yt-downloader/
├── main.go                 # Entrada principal e configuração de rotas
├── handlers/
│   ├── download.go         # Handler WebSocket para downloads
│   └── result.go           # Handler para entrega de arquivos
├── utils/
│   ├── redis.go            # Configuração e operações Redis
│   ├── cleanup.go          # Sistema de limpeza automática
│   ├── download.go         # Processamento de downloads
│   └── websocket.go        # Utilitários WebSocket
├── downloads/              # Diretório de arquivos temporários
└── go.mod                  # Dependências Go
```

## 🔧 Funcionalidades Avançadas

### Sistema de Limpeza Automática

- Remove arquivos antigos automaticamente
- Configurável por tamanho total (padrão: 2GB)
- Executado após cada download concluído

### Cache Inteligente

- Evita reprocessamento de URLs já baixadas
- Compartilha downloads entre múltiplos clientes
- Persiste entre reinicializações do servidor

### Compressão Automática

- Múltiplos arquivos → ZIP automático
- Arquivo único → Download direto
- Nomes de arquivo preservados

## ⚠️ Considerações de Segurança

### Implementadas

- Validação de parâmetros de entrada
- Sanitização de nomes de arquivos
- Limpeza automática de arquivos temporários

### Recomendadas para Produção

```go
// CORS restrito
var upgrader = websocket.Upgrader{
  CheckOrigin: func(r *http.Request) bool {
    origin := r.Header.Get("Origin")
    return origin == "https://meudominio.com"
  },
}

// Rate limiting
// Implementar middleware de rate limiting

// Autenticação
// Adicionar sistema de tokens/API keys

// HTTPS obrigatório
// Configurar TLS em produção
```

## 🚨 Limitações e Cuidados

### Limitações Técnicas

- **Dependência externa**: Requer `yt-dlp` instalado
- **Armazenamento**: Arquivos temporários consomem espaço
- **CPU intensivo**: Downloads podem sobrecarregar servidor
- **Rede**: Dependente da velocidade de internet

### Cuidados Operacionais

- **YouTube Terms**: Respeitar termos de uso do YouTube
- **Copyright**: Baixar apenas conteúdo autorizado
- **Recursos**: Monitorar uso de CPU, memória e disco
- **Logs**: Logs podem conter URLs sensíveis

### Monitoramento Recomendado

```sh
# Verificar espaço em disco
df -h downloads/

# Monitorar processos yt-dlp
ps aux | grep yt-dlp

# Status Redis
redis-cli ping
redis-cli info memory
```

## 🔮 Possibilidades de Expansão

### Funcionalidades Futuras

- **Playlist support**: Download de playlists completas
- **Formatos adicionais**: FLAC, OGG, diferentes qualidades
- **Metadata**: Extração e edição de metadados
- **Preview**: Thumbnail e informações antes do download
- **Agendamento**: Downloads programados
- **API de progresso**: Porcentagem detalhada de progresso

### Melhorias Técnicas

- **Database**: PostgreSQL para histórico detalhado
- **Queue system**: RabbitMQ para processamento assíncrono
- **Microservices**: Separar download engine
- **Container**: Docker para deploy facilitado
- **CDN**: Cache de arquivos populares
- **Analytics**: Dashboard de uso e performance

### Integrações Possíveis

```go
// Webhook notifications
type WebhookConfig struct {
  URL    string   `json:"url"`
  Events []string `json:"events"` // completed, error
}

// Cloud storage
func UploadToS3(filePath, bucket, key string) error {
  // Upload para AWS S3, Google Cloud, etc.
}

// Email notifications
func SendCompletionEmail(email, downloadID string) error {
  // Notificar usuário por email
}
```

## 🏃‍♂️ Teste Rápido

```sh
# Iniciar servidor
go run main.go

# Em outro terminal, testar health
curl http://localhost:8080/health

# Testar WebSocket (usando wscat)
# wscat -c ws://localhost:8080/download?url=...&format=mp3
```
