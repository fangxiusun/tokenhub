package router

import (
	"embed"
	"net/http"
	"strings"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

	"github.com/fangxiusun/tokenhub/common"
	"github.com/fangxiusun/tokenhub/controller"
	"github.com/fangxiusun/tokenhub/middleware"
	"github.com/fangxiusun/tokenhub/model"
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

		// User authentication (public)
		api.POST("/user/register", controller.Register)
		api.POST("/user/login", controller.Login)
		api.POST("/user/logout", controller.Logout)

		// User self-service (requires auth)
		user := api.Group("/user")
		user.Use(middleware.UserAuth())
		{
			user.GET("/self", controller.GetSelf)
			user.PUT("/self", controller.UpdateSelf)
			user.DELETE("/self", controller.DeleteSelf)
			user.POST("/access-token", controller.GenerateAccessToken)
			user.GET("/aff-code", controller.GetAffCode)

			// 2FA management
			user.POST("/2fa/enable", controller.EnableTwoFA)
			user.POST("/2fa/verify", controller.VerifyTwoFA)
			user.POST("/2fa/disable", controller.DisableTwoFA)

			// Passkey management
			user.POST("/passkey/enable", controller.EnablePasskey)
			user.POST("/passkey/verify", controller.VerifyPasskey)
			user.DELETE("/passkey", controller.DeletePasskey)
		}

		// 2FA login (no auth required)
		api.POST("/user/login/2fa", controller.VerifyTwoFALogin)

		// Passkey login (no auth required)
		api.POST("/user/passkey/login/begin", controller.PasskeyLoginBegin)
		api.POST("/user/passkey/login/complete", controller.PasskeyLoginComplete)

		// Admin routes (requires admin auth)
		admin := api.Group("/admin")
		admin.Use(middleware.AdminAuth())
		{
			admin.GET("/user", controller.GetAllUsers)
			admin.GET("/user/search", controller.SearchUsers)
			admin.GET("/user/:id", controller.GetUser)
			admin.POST("/user", controller.CreateUser)
			admin.PUT("/user/:id", controller.UpdateUser)
			admin.DELETE("/user/:id", controller.DeleteUser)
			admin.POST("/user/manage", controller.ManageUser)

			// Privilege group management
			admin.GET("/privilege-group", controller.GetAllPrivilegeGroups)
			admin.GET("/privilege-group/:id", controller.GetPrivilegeGroup)
			admin.POST("/privilege-group", controller.CreatePrivilegeGroup)
			admin.PUT("/privilege-group/:id", controller.UpdatePrivilegeGroup)
			admin.DELETE("/privilege-group/:id", controller.DeletePrivilegeGroup)
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
	relay.Use(middleware.ReadRequestBody())
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
		"setup":   model.RootUserExists(),
	})
}

func postSetup(c *gin.Context) {
	if model.RootUserExists() {
		c.JSON(http.StatusOK, gin.H{
			"success": false,
			"message": "Root user already exists",
		})
		return
	}

	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Username and password are required",
		})
		return
	}

	user := model.User{
		Username:    req.Username,
		Password:    req.Password,
		DisplayName: req.Username,
		Role:        100,
		Status:      1,
		Group:       "default",
	}

	if err := user.Insert(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create root user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Root user created successfully",
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

