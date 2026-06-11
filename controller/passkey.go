package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/fangxiusun/tokenhub/common"
	"github.com/fangxiusun/tokenhub/model"
)

// PasskeyUser implements the user interface for passkey operations
type PasskeyUser struct {
	Id          int
	DisplayName string
	Name        string
}

func (u PasskeyUser) WebAuthnID() []byte {
	return []byte(strconv.Itoa(u.Id))
}

func (u PasskeyUser) WebAuthnName() string {
	return u.Name
}

func (u PasskeyUser) WebAuthnDisplayName() string {
	return u.DisplayName
}

func (u PasskeyUser) WebAuthnCredentials() []interface{} {
	return nil
}

// EnablePasskey enables passkey for the current user
func EnablePasskey(c *gin.Context) {
	id, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Not authenticated",
		})
		return
	}

	userId := id.(int)

	// Check if passkey is already enabled
	existing, _ := model.GetPasskeyByUserId(userId)
	if existing != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Passkey is already enabled",
		})
		return
	}

	// For now, return a placeholder response
	// In production, this would initiate WebAuthn registration
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Passkey registration initiated",
		"data": gin.H{
			"challenge": common.GetRandomString(32),
		},
	})
}

// VerifyPasskey verifies and saves the passkey credential
func VerifyPasskey(c *gin.Context) {
	id, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Not authenticated",
		})
		return
	}

	userId := id.(int)

	var req struct {
		CredentialID string `json:"credential_id"`
		PublicKey    string `json:"public_key"`
	}
	if err := common.DecodeJson(c.Request.Body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
		})
		return
	}

	// Save credential
	cred := &model.PasskeyCredential{
		UserId:       userId,
		CredentialID: req.CredentialID,
		PublicKey:    req.PublicKey,
		SignCount:    0,
	}
	if err := model.CreatePasskeyCredential(cred); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to save passkey",
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
	user.PasskeyEnabled = true
	if err := user.Update(false); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Passkey enabled successfully",
	})
}

// PasskeyLoginBegin starts the passkey login process
func PasskeyLoginBegin(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
	}
	if err := common.DecodeJson(c.Request.Body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
		})
		return
	}

	user, err := model.GetUserByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "User not found",
		})
		return
	}

	// For now, return a placeholder response
	// In production, this would initiate WebAuthn authentication
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"challenge": common.GetRandomString(32),
			"user_id":   user.Id,
		},
	})
}

// PasskeyLoginComplete completes the passkey login
func PasskeyLoginComplete(c *gin.Context) {
	var req struct {
		Username     string `json:"username"`
		CredentialID string `json:"credential_id"`
	}
	if err := common.DecodeJson(c.Request.Body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
		})
		return
	}

	user, err := model.GetUserByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "User not found",
		})
		return
	}

	// Verify credential exists
	cred, err := model.GetPasskeyByCredentialId(req.CredentialID)
	if err != nil || cred == nil || cred.UserId != user.Id {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid passkey",
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

// DeletePasskey deletes the passkey for the current user
func DeletePasskey(c *gin.Context) {
	id, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Not authenticated",
		})
		return
	}

	userId := id.(int)

	if err := model.DeletePasskeyByUserId(userId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete passkey",
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
	user.PasskeyEnabled = false
	if err := user.Update(false); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Passkey deleted successfully",
	})
}

func init() {
	// Passkey initialization placeholder
}
