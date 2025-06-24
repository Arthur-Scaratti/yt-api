# YT API - YouTube Downloader Package

## 📋 Visão Geral

Um package Go moderno e eficiente para download de vídeos do YouTube com suporte completo a playlists, conversão para múltiplos formatos e atualizações de progresso em tempo real via WebSockets. O projeto foi refatorado para ser mais conciso, prático e sem dependências externas como Redis.

## 🏗️ Arquitetura

### Stack Tecnológica

- **Backend**: Go (Gin Framework)
- **WebSockets**: Gorilla WebSocket para progresso em tempo real
- **Download Engine**: yt-dlp
- **Compressão**: Archive/zip nativo
- **Cache**: Sistema de arquivos local

### Estrutura do Projeto

```
├── app/main.go            # Servidor standalone
├── start.go               # Função exportável para uso como package
├── config/config.go       # Configurações via variáveis de ambiente
├── handlers/              # Handlers HTTP/WebSocket
│   ├── download.go        # Handler principal de downloads
│   ├── playlist.go        # Servir arquivos de playlist
│   ├── playlistdl.go      # Download de playlists com progresso
│   └── websocket.go       # Gerenciamento de WebSockets
└── utils/                 # Utilitários
    ├── check.go           # Verificação de cache e arquivos
    ├── cleanup.go         # Limpeza automática de arquivos
    ├── sanitize.go        # Sanitização de nomes de arquivo
    ├── startcleanup.go    # Inicialização da limpeza automática
    └── tracking.go        # Rastreamento de acesso aos arquivos
```

## 🚀 Instalação e Uso

### Pré-requisitos

- Go 1.21+
- yt-dlp instalado e disponível no PATH

### Como Package Importável

```go
package main

import (
    "github.com/Arthur-Scaratti/yt-api"
)

func main() {
    // Inicia o servidor e retorna a instância do Gin
    router := ytapi.Start()
    // O servidor já está rodando
}
```

### Como Servidor Standalone

```bash
# Clone o repositório
git clone https://github.com/Arthur-Scaratti/yt-api.git
cd yt-api

# Execute o servidor
go run app/main.go
```

### Configuração via Variáveis de Ambiente

Crie um arquivo `.env` ou configure as variáveis de ambiente:

```env
# Servidor
GIN_MODE=debug
HOST=localhost
PORT=8080

# Handlers (rotas)
DOWNLOAD_HANDLER=/download
PLAYLIST_HANDLER=/playlist
WEBSOCKET_HANDLER=/ws

# Diretórios
DOWNLOAD_DIR=./downloads
FILE_PERMISSIONS=0755

# Configurações yt-dlp
YTDLP_CONCURRENT_FRAGMENTS=4
YTDLP_FRAGMENT_RETRIES=10
YTDLP_RETRIES=10
YTDLP_EXTRACTOR_RETRIES=3
YTDLP_DEFAULT_QUALITY=720

# Templates de saída
OUTPUT_TEMPLATE_SINGLE=%(title)s.%(ext)s
OUTPUT_TEMPLATE_PLAYLIST=%(playlist_index)s - %(title)s.%(ext)s
PROGRESS_TEMPLATE=download:[%(info.id)s] %(info.title)s

# Padrões
DEFAULT_FORMAT=mp4
DEFAULT_QUALITY=720p
DEFAULT_PLAYLIST=false
DEFAULT_INDEX=
```

## 📡 API Endpoints

### 1. Download de Vídeo/Playlist

```http
GET /download?url={URL}&format={FORMAT}&quality={QUALITY}&playlist={BOOLEAN}&index={NUMBER}
```

**Parâmetros:**
- `url` (obrigatório): URL do YouTube
- `format` (opcional): mp3, mp4, mkv, webm (padrão: mp4)
- `quality` (opcional): 144p, 240p, 360p, 480p, 720p, 1080p (padrão: 720p)
- `playlist` (opcional): true/false (padrão: false)
- `index` (opcional): índice específico da playlist

**Respostas:**

*Download único ou item específico:*
```json
// Retorna o arquivo diretamente
```

*Playlist completa (primeira requisição):*
```json
{
  "id": "dl_abc123",
  "progressUrl": "/ws?id=dl_abc123"
}
```

*Playlist já processada:*
```json
{
  "status": "Ready",
  "id": "dl_abc123",
  "count": 15,
  "files": [
    {
      "index": "1",
      "title": "Título do Vídeo",
      "filename": "1 - Título do Vídeo.mp4"
    }
  ],
  "download": "/playlist?id=dl_abc123&index=N"
}
```

### 2. Servir Arquivos de Playlist

```http
GET /playlist?id={ID}&index={INDEX}
```

**Parâmetros:**
- `id` (obrigatório): ID da playlist
- `index` (opcional): índice específico ou vazio para ZIP completo

**Comportamento:**
- Sem `index`: retorna arquivo ZIP com toda a playlist
- Com `index`: retorna arquivo específico da playlist

### 3. WebSocket para Progresso

```http
GET /ws?id={ID}
```

**Mensagens recebidas:**
```json
{
  "id": "dl_abc123",
  "title": "Nome do arquivo sendo baixado"
}
```

```json
{
  "id": "dl_abc123",
  "title": "completed"
}
```

## 🔧 Funcionalidades

### Cache Inteligente
- Sistema de cache baseado em hash SHA256 dos parâmetros
- Reutilização automática de downloads existentes
- Rastreamento de último acesso para limpeza

### Suporte Completo a Playlists
- Download de playlists inteiras em background
- Progresso em tempo real via WebSockets
- Servir arquivos individuais ou ZIP completo
- Numeração automática dos arquivos

### Limpeza Automática
- Execução a cada 8 horas (500 minutos)
- Remove 50% dos arquivos mais antigos
- Baseado no último acesso aos arquivos

### Formatos Suportados
- **MP3**: Extração de áudio
- **MP4/MKV/WEBM**: Vídeo com merge automático
- Seleção inteligente de qualidade

### Sanitização de Arquivos
- Nomes de arquivo seguros para todos os sistemas
- Remoção de caracteres especiais
- Preservação da legibilidade

## 🔄 Fluxo de Funcionamento

### Download Único
1. Recebe requisição com URL
2. Gera hash único baseado nos parâmetros
3. Verifica cache existente
4. Se existe: retorna arquivo imediatamente
5. Se não existe: executa yt-dlp e retorna arquivo

### Download de Playlist
1. Recebe requisição com `playlist=true`
2. Gera hash único
3. Verifica cache existente
4. Se existe: retorna lista de arquivos
5. Se não existe: 
   - Retorna ID e URL do WebSocket
   - Inicia download em background
   - Envia progresso via WebSocket
   - Cliente pode acessar arquivos via `/playlist`

## 🛠️ Desenvolvimento

### Estrutura de Dados

```go
// Configuração principal
type Config struct {
    GinMode             string
    Port                string
    Host                string
    DownloadHandler     string
    PlaylistHandler     string
    WebSocketHandler    string
    DownloadDir         string
    FilePermissions     os.FileMode
    // ... outras configurações
}

// Informação de acesso para limpeza
type AccessInfo struct {
    LastAccessed time.Time `json:"last_accessed"`
}
```

### Principais Funções

- `DownloadHandler()`: Handler principal de downloads
- `PlaylistHandler()`: Servir arquivos de playlist
- `RunPlaylistDownload()`: Download de playlist com progresso
- `WebSocketHandler()`: Gerenciamento de conexões WebSocket
- `CheckExistingID()`: Verificação de cache
- `StartAutoCleanup()`: Limpeza automática

## 📝 Exemplos de Uso

### Download de vídeo único
```bash
curl "http://localhost:8080/download?url=https://youtube.com/watch?v=VIDEO_ID&format=mp4&quality=720p"
```

### Download de playlist completa
```bash
# Inicia o download
curl "http://localhost:8080/download?url=https://youtube.com/playlist?list=PLAYLIST_ID&playlist=true&format=mp3"

# Conecta ao WebSocket para progresso
# ws://localhost:8080/ws?id=RETURNED_ID

# Baixa arquivo específico
curl "http://localhost:8080/playlist?id=RETURNED_ID&index=1"

# Baixa ZIP completo
curl "http://localhost:8080/playlist?id=RETURNED_ID"
```

## 🔒 Segurança

- Sanitização automática de nomes de arquivo
- Validação de parâmetros de entrada
- Isolamento de arquivos por ID único
- Limpeza automática de arquivos antigos

## 📊 Performance

- Downloads paralelos com fragmentos concorrentes
- Cache eficiente baseado em hash
- Limpeza automática para gerenciamento de espaço
- WebSockets para comunicação eficiente

## 🤝 Contribuição

1. Fork o projeto
2. Crie uma branch para sua feature
3. Commit suas mudanças
4. Push para a branch
5. Abra um Pull Request

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE.txt) para detalhes.