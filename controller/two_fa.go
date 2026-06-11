package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"

	"github.com/fangxiusun/tokenhub/common"
	"github.com/fangxiusun/tokenhub/model"
)

type TwoFARequest struct {
	Code string `json:"code"`
}

// EnableTwoFA enables 2FA for the current user
func EnableTwoFA(c *gin.Context) {
	id, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Not authenticated",
		})
		return
	}

	userId := id.(int)

	// Check if 2FA is already enabled
	existing, _ := model.GetTwoFactorAuthByUserId(userId)
	if existing != nil && existing.Enabled {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "2FA is already enabled",
		})
		return
	}

	// Generate new secret
	secret, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "TokenHub",
		AccountName: getUserUsername(userId),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to generate 2FA secret",
		})
		return
	}

	// Save 2FA config
	tfa := &model.TwoFactorAuth{
		UserId:  userId,
		Secret:  secret.Secret(),
		Enabled: false,
	}
	if err := model.CreateOrUpdateTwoFactorAuth(tfa); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to save 2FA config",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"secret":  secret.Secret(),
			"qr_code": secret.URL(),
		},
	})
}

// VerifyTwoFA verifies the 2FA code and enables 2FA
func VerifyTwoFA(c *gin.Context) {
	id, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Not authenticated",
		})
		return
	}

	var req TwoFARequest
	if err := common.DecodeJson(c.Request.Body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
		})
		return
	}

	if req.Code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Code is required",
		})
		return
	}

	userId := id.(int)

	// Get 2FA config
	tfa, err := model.GetTwoFactorAuthByUserId(userId)
	if err != nil || tfa == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "2FA not configured",
		})
		return
	}

	// Verify code
	valid := totp.Validate(req.Code, tfa.Secret)
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid code",
		})
		return
	}

	// Enable 2FA
	tfa.Enabled = true
	if err := model.CreateOrUpdateTwoFactorAuth(tfa); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to enable 2FA",
		})
		return
	}

	// Update user
	user, err := model.GetUserById(userId, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get user",
		})
		return
	}
	user.TFAEnabled = true
	if err := user.Update(false); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "2FA enabled successfully",
	})
}

// DisableTwoFA disables 2FA for the current user
func DisableTwoFA(c *gin.Context) {
	id, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Not authenticated",
		})
		return
	}

	var req TwoFARequest
	if err := common.DecodeJson(c.Request.Body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
		})
		return
	}

	if req.Code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Code is required",
		})
		return
	}

	userId := id.(int)

	// Get 2FA config
	tfa, err := model.GetTwoFactorAuthByUserId(userId)
	if err != nil || tfa == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "2FA not configured",
		})
		return
	}

	// Verify code
	valid := totp.Validate(req.Code, tfa.Secret)
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid code",
		})
		return
	}

	// Delete 2FA config
	if err := model.DeleteTwoFactorAuth(userId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete 2FA config",
		})
		return
	}

	// Update user
	user, err := model.GetUserById(userId, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get user",
		})
		return
	}
	user.TFAEnabled = false
	if err := user.Update(false); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "2FA disabled successfully",
	})
}

// VerifyTwoFALogin verifies 2FA code during login
func VerifyTwoFALogin(c *gin.Context) {
	var req struct {
		Code     string `json:"code"`
		Username string `json:"username"`
	}
	if err := common.DecodeJson(c.Request.Body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
		})
		return
	}

	if req.Code == "" || req.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Code and username are required",
		})
		return
	}

	// Get user
	user, err := model.GetUserByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid credentials",
		})
		return
	}

	// Get 2FA config
	tfa, err := model.GetTwoFactorAuthByUserId(user.Id)
	if err != nil || tfa == nil || !tfa.Enabled {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "2FA is not enabled",
		})
		return
	}

	// Verify code
	valid := totp.Validate(req.Code, tfa.Secret)
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid code",
		})
		return
	}

	// Generate JWT
	token, err := common.GenerateJWT(user.Id, user.Username, user.Role, user.Status, user.Group)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to generate token",
		})
		return
	}

	// Update last login
	model.UpdateUserLastLoginAt(user.Id)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"token":        token,
			"id":           user.Id,
			"username":     user.Username,
			"display_name": user.DisplayName,
			"role":         user.Role,
			"status":       user.Status,
			"group":        user.Group,
		},
	})
}

func getUserUsername(userId int) string {
	user, err := model.GetUserById(userId)
	if err != nil {
		return ""
	}
	return user.Username
}
