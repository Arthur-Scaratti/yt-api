package ytapi

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Arthur-Scaratti/yt-api/handlers"
	"github.com/Arthur-Scaratti/yt-api/utils"
	"github.com/gin-gonic/gin"
)

type ServerConfig struct {
	Port         string
	Debug        bool
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type RedisConfig struct {
	Address  string
	Password string
	DB       int
}

func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Port:         ":8080",
		Debug:        false,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

func DefaultRedisConfig() *RedisConfig {
	return &RedisConfig{
		Address:  "redis:6379",
		Password: "",
		DB:       0,
	}
}

func QuickStart() error {
	log.Println("🚀 Iniciando YT API com configurações padrão...")

	serverConfig := DefaultServerConfig()
	redisConfig := DefaultRedisConfig()

	return Start(serverConfig, redisConfig)
}

func QuickStartWithContext(ctx context.Context) error {
	log.Println("🚀 Iniciando YT API com configurações padrão e context...")

	serverConfig := DefaultServerConfig()
	redisConfig := DefaultRedisConfig()

	return StartWithContext(ctx, serverConfig, redisConfig)
}

func Start(serverConfig *ServerConfig, redisConfig *RedisConfig) error {
	r := InitializeServer(serverConfig, redisConfig)

	srv := &http.Server{
		Addr:         serverConfig.Port,
		Handler:      r,
		ReadTimeout:  serverConfig.ReadTimeout,
		WriteTimeout: serverConfig.WriteTimeout,
		IdleTimeout:  serverConfig.IdleTimeout,
	}

	log.Printf("🌐 Servidor iniciando na porta %s", serverConfig.Port)
	log.Printf("📊 Redis conectando em %s", redisConfig.Address)

	return srv.ListenAndServe()
}

func StartWithContext(ctx context.Context, serverConfig *ServerConfig, redisConfig *RedisConfig) error {
	r := InitializeServer(serverConfig, redisConfig)

	srv := &http.Server{
		Addr:         serverConfig.Port,
		Handler:      r,
		ReadTimeout:  serverConfig.ReadTimeout,
		WriteTimeout: serverConfig.WriteTimeout,
		IdleTimeout:  serverConfig.IdleTimeout,
	}

	serverErr := make(chan error, 1)

	go func() {
		log.Printf("🌐 Servidor iniciando na porta %s", serverConfig.Port)
		log.Printf("📊 Redis conectando em %s", redisConfig.Address)
		serverErr <- srv.ListenAndServe()
	}()

	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
		log.Println("⏹️  Iniciando graceful shutdown...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	}
}

// InitializeServer inicializa e configura o servidor
func InitializeServer(serverConfig *ServerConfig, redisConfig *RedisConfig) *gin.Engine {
	// Inicializa Redis com configuração customizada
	utils.InitRedis(redisConfig.Address, redisConfig.Password, redisConfig.DB)

	if !serverConfig.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	SetupRoutes(r)
	return r
}

func SetupRoutes(r *gin.Engine) {
	r.GET("/download", handlers.DownloadHandler)
	r.GET("/result", handlers.ResultHandler)
	r.GET("/health", handlers.HealthCheckHandler)
}

func GetServerInfo() map[string]interface{} {
	return map[string]interface{}{
		"name":        "YT MP3 Downloader API",
		"version":     "1.0.1",
		"status":      "running",
		"description": "API para download de vídeos/áudios do YouTube",
		"endpoints": map[string]string{
			"/download": "WebSocket - Inicia download (GET com ?url=...&format=mp3|mp4)",
			"/result":   "HTTP - Obtém arquivo processado (GET com ?id=...)",
			"/health":   "HTTP - Verificação de saúde da API",
		},
	}
}
