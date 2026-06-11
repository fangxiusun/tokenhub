package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"

	"github.com/fangxiusun/tokenhub/common"
	"github.com/fangxiusun/tokenhub/i18n"
	"github.com/fangxiusun/tokenhub/middleware"
	"github.com/fangxiusun/tokenhub/model"
	"github.com/fangxiusun/tokenhub/router"
	"github.com/fangxiusun/tokenhub/service"
)

//go:embed web/default/dist
var buildFS embed.FS

//go:embed web/default/dist/index.html
var indexPage []byte

func main() {
	startTime := fmt.Sprintf("%d", os.Getpid())
	log.Printf("Server starting (PID: %s)", startTime)

	// Initialize resources
	if err := InitResources(); err != nil {
		log.Fatalf("Failed to initialize resources: %v", err)
	}

	// Setup HTTP server
	engine := setupHttpServer()

	// Start background tasks
	startBackgroundTasks()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		log.Println("Shutting down server...")
		model.Close()
		os.Exit(0)
	}()

	// Start HTTP server
	port := common.GetEnvOrDefault("PORT", "3000")
	log.Printf("Server starting on port %s", port)
	if err := engine.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func InitResources() error {
	// Initialize logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Initialize database
	if err := model.Init(); err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	// Initialize Redis (optional)
	common.InitRedis()

	// Initialize i18n
	if err := i18n.Init(); err != nil {
		log.Printf("Warning: Failed to initialize i18n: %v", err)
	}

	// Initialize services
	service.Init()

	log.Println("All resources initialized successfully")
	return nil
}

func setupHttpServer() *gin.Engine {
	if common.GetEnvOrDefault("GIN_MODE", "debug") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()

	// Custom recovery
	engine.Use(gin.CustomRecovery(func(c *gin.Context, err any) {
		log.Printf("Panic recovered: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": gin.H{
				"message": fmt.Sprintf("Internal server error: %v", err),
				"type":    "server_error",
			},
		})
	}))

	// Apply middleware
	engine.Use(middleware.RequestId())
	engine.Use(middleware.Logger())
	engine.Use(middleware.Cors())

	// Initialize session store
	store := cookie.NewStore([]byte(common.GetEnvOrDefault("SESSION_SECRET", "default-secret")))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   2592000, // 30 days
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
	})
	engine.Use(sessions.Sessions("session", store))

	// Setup routes
	router.SetRouter(engine, router.ThemeAssets{
		BuildFS:   buildFS,
		IndexPage: indexPage,
	})

	return engine
}

func startBackgroundTasks() {
	// Start background tasks here
	log.Println("Background tasks started")
}

