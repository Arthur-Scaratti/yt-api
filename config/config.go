package config

import (
    "log"
    "os"
    "strconv"
    
    "github.com/joho/godotenv"
)

type Config struct {
    // Server
    GinMode    string
    Port       string
    Host       string
    
    // Handlers
    DownloadHandler  string
    PlaylistHandler  string
    WebSocketHandler string
    
    // Download
    DownloadDir     string
    FilePermissions os.FileMode
    
    // yt-dlp
    ConcurrentFragments int
    FragmentRetries     int
    Retries            int
    ExtractorRetries   int
    DefaultQualityYTDLP int
    
    // Templates
    OutputTemplateSingle   string
    OutputTemplatePlaylist string
    ProgressTemplate       string
    
    // Defaults
    DefaultFormat   string
    DefaultQuality  string
    DefaultPlaylist string
    DefaultIndex    string
}

func Load() *Config {
    // Carrega o arquivo .env se existir
    if err := godotenv.Load(); err != nil {
        log.Println("Arquivo .env não encontrado, usando variáveis de ambiente do sistema")
    }
    
    return &Config{
        // Server
        GinMode: getEnv("GIN_MODE"),
        Port:    getEnv("PORT"),
        Host:    getEnv("HOST"),
        
        // Handlers
        DownloadHandler:  getEnv("DOWNLOAD_HANDLER"),
        PlaylistHandler:  getEnv("PLAYLIST_HANDLER"),
        WebSocketHandler: getEnv("WEBSOCKET_HANDLER"),
        
        // Download
        DownloadDir:     getEnv("DOWNLOAD_DIR"),
        FilePermissions: os.FileMode(getEnvInt("FILE_PERMISSIONS")),
        
        // yt-dlp
        ConcurrentFragments: getEnvInt("YTDLP_CONCURRENT_FRAGMENTS"),
        FragmentRetries:     getEnvInt("YTDLP_FRAGMENT_RETRIES"),
        Retries:            getEnvInt("YTDLP_RETRIES"),
        ExtractorRetries:   getEnvInt("YTDLP_EXTRACTOR_RETRIES"),
        DefaultQualityYTDLP: getEnvInt("YTDLP_DEFAULT_QUALITY"),
        
        // Templates
        OutputTemplateSingle:   getEnv("OUTPUT_TEMPLATE_SINGLE"),
        OutputTemplatePlaylist: getEnv("OUTPUT_TEMPLATE_PLAYLIST"),
        ProgressTemplate:       getEnv("PROGRESS_TEMPLATE"),
        
        // Defaults
        DefaultFormat:   getEnv("DEFAULT_FORMAT"),
        DefaultQuality:  getEnv("DEFAULT_QUALITY"),
        DefaultPlaylist: getEnv("DEFAULT_PLAYLIST"),
        DefaultIndex:    getEnv("DEFAULT_INDEX"),
    }
}

func getEnv(key string) string {
    return os.Getenv(key)
}

func getEnvInt(key string) int {
    if value := os.Getenv(key); value != "" {
        if intValue, err := strconv.Atoi(value); err == nil {
            return intValue
        }
    }
    return 0
}