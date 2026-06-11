package router

import (
	"embed"
	"net/http"
	"strings"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

	"github.com/your-username/your-project/common"
	"github.com/your-username/your-project/controller"
	"github.com/your-username/your-project/middleware"
)

// ThemeAssets holds embedded frontend assets
type ThemeAssets struct {
	BuildFS   embed.FS
	IndexPage []byte
}

// SetRouter sets up all routes
func SetRouter(engine *gin.Engine, assets ThemeAssets) {
	// Setup web routes (frontend)
	SetWebRouter(engine, assets)

	// Setup API routes
	SetApiRouter(engine)

	// Setup relay routes (AI API proxy)
	SetRelayRouter(engine)
}

// SetApiRouter sets up API routes
func SetApiRouter(engine *gin.Engine) {
	api := engine.Group("/api")
	{
		// System
		api.GET("/status", getStatus)
		api.GET("/setup", getSetup)
		api.POST("/setup", postSetup)

		// User authentication
		api.POST("/user/register", register)
		api.POST("/user/login", login)
		api.POST("/user/logout", logout)

		// User self-service (requires auth)
		user := api.Group("/user")
		user.Use(middleware.UserAuth())
		{
			user.GET("/self", getSelf)
			user.PUT("/self", updateSelf)
			user.DELETE("/self", deleteSelf)
			user.GET("/token", getTokens)
			user.POST("/token", createToken)
			user.PUT("/token/:id", updateToken)
			user.DELETE("/token/:id", deleteToken)
		}

		// Admin routes (requires admin auth)
		admin := api.Group("/admin")
		admin.Use(middleware.AdminAuth())
		{
			admin.GET("/user", getUsers)
			admin.POST("/user", createUser)
			admin.PUT("/user/:id", updateUser)
			admin.DELETE("/user/:id", deleteUser)

			admin.GET("/channel", getChannels)
			admin.POST("/channel", createChannel)
			admin.PUT("/channel/:id", updateChannel)
			admin.DELETE("/channel/:id", deleteChannel)
			admin.POST("/channel/test", testChannel)
		}

		// Root routes (requires root auth)
		root := api.Group("/root")
		root.Use(middleware.RootAuth())
		{
			root.GET("/option", getOptions)
			root.PUT("/option", updateOption)
		}
	}
}

// SetRelayRouter sets up relay routes
func SetRelayRouter(engine *gin.Engine) {
	relay := engine.Group("/v1")
	relay.Use(middleware.TokenAuth())
	relay.Use(middleware.Distribute())
	{
		relay.GET("/models", getModels)
		relay.POST("/chat/completions", chatCompletions)
		relay.POST("/completions", completions)
		relay.POST("/embeddings", embeddings)
		relay.POST("/images/generations", imageGenerations)
	}
}

// SetWebRouter sets up web routes for frontend
func SetWebRouter(engine *gin.Engine, assets ThemeAssets) {
	// Create embedded filesystem
	frontendFS := common.EmbedFolder(assets.BuildFS, "web/default/dist")

	// Serve static files
	engine.Use(static.Serve("/", frontendFS))

	// Handle SPA routing (no route matched)
	engine.NoRoute(func(c *gin.Context) {
		// If it's an API request, return 404
		if strings.HasPrefix(c.Request.RequestURI, "/v1") ||
			strings.HasPrefix(c.Request.RequestURI, "/api") ||
			strings.HasPrefix(c.Request.RequestURI, "/assets") {
			controller.RelayNotFound(c)
			return
		}

		// Serve index.html for SPA routing
		c.Header("Cache-Control", "no-cache")
		c.Data(http.StatusOK, "text/html; charset=utf-8", assets.IndexPage)
	})
}

// Handler stubs
func getStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "OK",
		"version": "1.0.0",
	})
}

func getSetup(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"setup":   false,
	})
}

func postSetup(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Setup completed",
	})
}

func register(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Registration successful",
	})
}

func login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login successful",
	})
}

func logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logout successful",
	})
}

func getSelf(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

func updateSelf(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Update successful",
	})
}

func deleteSelf(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Delete successful",
	})
}

func getTokens(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    []interface{}{},
	})
}

func createToken(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Token created",
	})
}

func updateToken(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Token updated",
	})
}

func deleteToken(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Token deleted",
	})
}

func getUsers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    []interface{}{},
	})
}

func createUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User created",
	})
}

func updateUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User updated",
	})
}

func deleteUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User deleted",
	})
}

func getChannels(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    []interface{}{},
	})
}

func createChannel(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Channel created",
	})
}

func updateChannel(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Channel updated",
	})
}

func deleteChannel(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Channel deleted",
	})
}

func testChannel(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Channel test passed",
	})
}

func getOptions(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    map[string]interface{}{},
	})
}

func updateOption(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Option updated",
	})
}

func getModels(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    []interface{}{},
	})
}

func chatCompletions(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Chat completions endpoint",
	})
}

func completions(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Completions endpoint",
	})
}

func embeddings(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Embeddings endpoint",
	})
}

func imageGenerations(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Image generations endpoint",
	})
}
