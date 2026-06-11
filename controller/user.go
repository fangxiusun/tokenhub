package controller

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"

	"github.com/fangxiusun/tokenhub/common"
	"github.com/fangxiusun/tokenhub/constant"
	"github.com/fangxiusun/tokenhub/model"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email"`
}

// Login handles user login
func Login(c *gin.Context) {
	// Check if password login is enabled
	if !common.PasswordLoginEnabled {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Password login is disabled",
		})
		return
	}

	var req LoginRequest
	if err := common.DecodeJson(c.Request.Body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
		})
		return
	}

	username := strings.TrimSpace(req.Username)
	password := strings.TrimSpace(req.Password)

	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Username and password are required",
		})
		return
	}

	user := &model.User{
		Username: username,
		Password: password,
	}

	// Validate credentials
	err := user.ValidateAndFill()
	if err != nil {
		switch {
		case errors.Is(err, model.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Invalid username or password",
			})
		case errors.Is(err, model.ErrUserEmptyCredentials):
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Username and password are required",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Internal server error",
			})
		}
		return
	}

	// Generate JWT token
	token, err := common.GenerateJWT(user.Id, user.Username, user.Role, user.Status, user.Group)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to generate token",
		})
		return
	}

	// Check if 2FA is enabled
	if user.TFAEnabled {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "",
			"data": gin.H{
				"require_2fa": true,
				"username":    user.Username,
			},
		})
		return
	}

	// Update last login time
	model.UpdateUserLastLoginAt(user.Id)

	// Setup session
	session := sessions.Default(c)
	session.Set("id", user.Id)
	session.Set("username", user.Username)
	session.Set("role", user.Role)
	session.Set("status", user.Status)
	session.Set("group", user.Group)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to save session",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
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

// Logout handles user logout
func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to clear session",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
	})
}

// Register handles user registration
func Register(c *gin.Context) {
	// Check if registration is enabled
	if !common.RegisterEnabled {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Registration is disabled",
		})
		return
	}

	// Check if password registration is enabled
	if !common.PasswordRegisterEnabled {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Password registration is disabled",
		})
		return
	}

	var req RegisterRequest
	if err := common.DecodeJson(c.Request.Body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
		})
		return
	}

	// Validate input
	if req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Username and password are required",
		})
		return
	}

	// Check if user exists
	exist, err := model.CheckUserExistOrDeleted(req.Username, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
		})
		return
	}
	if exist {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": "Username or email already exists",
		})
		return
	}

	// Create clean user with email
	cleanUser := model.User{
		Username:    strings.TrimSpace(req.Username),
		Password:    req.Password,
		DisplayName: req.Username,
		Email:       req.Email,
		Role:        constant.RoleCommonUser,
		Status:      constant.UserStatusEnabled,
	}

	if err := cleanUser.Insert(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Registration successful",
	})
}

// GetSelf returns the current user's profile
func GetSelf(c *gin.Context) {
	id, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Not authenticated",
		})
		return
	}

	user, err := model.GetUserById(id.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data": gin.H{
			"id":           user.Id,
			"username":     user.Username,
			"display_name": user.DisplayName,
			"role":         user.Role,
			"status":       user.Status,
			"email":        user.Email,
			"group":        user.Group,
			"quota":        user.Quota,
			"used_quota":   user.UsedQuota,
			"request_count": user.RequestCount,
			"aff_code":     user.AffCode,
			"aff_count":    user.AffCount,
			"aff_quota":    user.AffQuota,
			"inviter_id":   user.InviterId,
			"created_at":   user.CreatedAt,
			"last_login_at": user.LastLoginAt,
		},
	})
}

// UpdateSelf updates the current user's profile
func UpdateSelf(c *gin.Context) {
	var requestData map[string]interface{}
	if err := common.DecodeJson(c.Request.Body, &requestData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
		})
		return
	}

	id, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Not authenticated",
		})
		return
	}

	user, err := model.GetUserById(id.(int), true)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "User not found",
		})
		return
	}

	// Update display_name if provided
	if displayName, ok := requestData["display_name"].(string); ok {
		user.DisplayName = displayName
	}

	// Handle password update
	if newPassword, ok := requestData["password"].(string); ok && newPassword != "" {
		if originalPassword, ok := requestData["original_password"].(string); ok {
			if !common.ValidatePasswordAndHash(originalPassword, user.Password) {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "Original password is incorrect",
				})
				return
			}
			user.Password = newPassword
			if err := user.Update(true); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Failed to update password",
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "Password updated successfully",
			})
			return
		}
	}

	// Update user
	if err := user.Update(false); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Profile updated successfully",
	})
}

// DeleteSelf deletes the current user's account
func DeleteSelf(c *gin.Context) {
	id, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Not authenticated",
		})
		return
	}

	user, err := model.GetUserById(id.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "User not found",
		})
		return
	}

	// Prevent root user from deleting themselves
	if user.Role == constant.RoleRootUser {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Cannot delete root user",
		})
		return
	}

	if err := user.Delete(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete user",
		})
		return
	}

	// Clear session
	session := sessions.Default(c)
	session.Clear()
	_ = session.Save()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Account deleted successfully",
	})
}

// GetAllUsers returns all users (admin only)
func GetAllUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	users, total, err := model.GetAllUsers(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get users",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
			"items":     users,
		},
	})
}

// SearchUsers searches users (admin only)
func SearchUsers(c *gin.Context) {
	keyword := c.Query("keyword")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	users, total, err := model.SearchUsers(keyword, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to search users",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data": gin.H{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
			"items":     users,
		},
	})
}

// GetUser returns a specific user (admin only)
func GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid user ID",
		})
		return
	}

	user, err := model.GetUserById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    user,
	})
}

// CreateUser creates a new user (admin only)
func CreateUser(c *gin.Context) {
	var req struct {
		Username    string `json:"username"`
		Password    string `json:"password"`
		DisplayName string `json:"display_name"`
		Email       string `json:"email"`
		Role        int    `json:"role"`
	}
	if err := common.DecodeJson(c.Request.Body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
		})
		return
	}

	if req.Username == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Username and password are required",
		})
		return
	}

	// Check if user exists
	exist, err := model.CheckUserExistOrDeleted(req.Username, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Database error",
		})
		return
	}
	if exist {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": "Username or email already exists",
		})
		return
	}

	// Get admin role and check permissions
	adminRole := c.GetInt("role")
	if req.Role >= adminRole {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Cannot create user with equal or higher role",
		})
		return
	}

	cleanUser := model.User{
		Username:    strings.TrimSpace(req.Username),
		Password:    req.Password,
		DisplayName: req.DisplayName,
		Email:       req.Email,
		Role:        req.Role,
		Status:      constant.UserStatusEnabled,
	}

	if cleanUser.DisplayName == "" {
		cleanUser.DisplayName = cleanUser.Username
	}

	if err := cleanUser.Insert(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User created successfully",
	})
}

// UpdateUser updates a user (admin only)
func UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid user ID",
		})
		return
	}

	var updatedUser model.User
	if err := common.DecodeJson(c.Request.Body, &updatedUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
		})
		return
	}

	originUser, err := model.GetUserById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "User not found",
		})
		return
	}

	// Check permissions
	myRole := c.GetInt("role")
	if myRole <= originUser.Role {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Cannot modify user with equal or higher role",
		})
		return
	}

	// Update fields
	if updatedUser.Username != "" {
		originUser.Username = updatedUser.Username
	}
	if updatedUser.DisplayName != "" {
		originUser.DisplayName = updatedUser.DisplayName
	}
	if updatedUser.Group != "" {
		originUser.Group = updatedUser.Group
	}

	updatePassword := updatedUser.Password != ""
	if updatePassword {
		originUser.Password = updatedUser.Password
	}

	if err := originUser.Edit(updatePassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User updated successfully",
	})
}

// DeleteUser deletes a user (admin only)
func DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid user ID",
		})
		return
	}

	user, err := model.GetUserById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "User not found",
		})
		return
	}

	// Check permissions
	myRole := c.GetInt("role")
	if myRole <= user.Role {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Cannot delete user with equal or higher role",
		})
		return
	}

	if err := user.HardDelete(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User deleted successfully",
	})
}

// ManageUser manages user status/role (admin only)
func ManageUser(c *gin.Context) {
	var req struct {
		Id       int    `json:"id"`
		Action   string `json:"action"`
		Value    int    `json:"value"`
		Password string `json:"password"`
	}
	if err := common.DecodeJson(c.Request.Body, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
		})
		return
	}

	user, err := model.GetUserById(req.Id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "User not found",
		})
		return
	}

	// Check permissions
	myRole := c.GetInt("role")
	if myRole < user.Role {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Cannot manage user with higher role",
		})
		return
	}

	switch req.Action {
	case "disable":
		user.Status = constant.UserStatusDisabled
	case "enable":
		user.Status = constant.UserStatusEnabled
	case "delete":
		if user.Role == constant.RoleRootUser {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Cannot delete root user",
			})
			return
		}
		if err := user.Delete(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to delete user",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "User deleted",
		})
		return
	case "promote":
		// Root user can promote to admin
		// Admin can promote common user to admin
		if myRole == constant.RoleRootUser {
			if user.Role >= constant.RoleAdminUser {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "User is already admin or higher",
				})
				return
			}
			user.Role = constant.RoleAdminUser
		} else if myRole == constant.RoleAdminUser {
			// Admin can only promote common user to admin
			if user.Role != constant.RoleCommonUser {
				c.JSON(http.StatusForbidden, gin.H{
					"success": false,
					"message": "Admin can only promote common users",
				})
				return
			}
			user.Role = constant.RoleAdminUser
		}
	case "demote":
		// Cannot demote root user
		if user.Role == constant.RoleRootUser {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Cannot demote root user",
			})
			return
		}
		// Root user can demote anyone
		// Admin can demote common user (but admin cannot demote admin if they are not root)
		if myRole == constant.RoleRootUser {
			user.Role = constant.RoleCommonUser
		} else if myRole == constant.RoleAdminUser {
			if user.Role != constant.RoleAdminUser {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "Admin can only demote other admins",
				})
				return
			}
			user.Role = constant.RoleCommonUser
		}
	case "set_role":
		// Allow setting specific role (for super admin)
		if myRole != constant.RoleRootUser {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Only root user can set role",
			})
			return
		}
		if user.Role == constant.RoleRootUser && req.Value != constant.RoleRootUser {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Cannot change root user role",
			})
			return
		}
		if !constant.IsValidateRole(req.Value) {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Invalid role value",
			})
			return
		}
		user.Role = req.Value
	case "set_group":
		// Allow admin to set user's privilege group
		// The group name is passed as the Password field for simplicity
		if req.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Group name is required",
			})
			return
		}
		user.Group = req.Password
	case "reset_password":
		// Allow admin to reset password for lower-level users
		if req.Password == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "New password is required",
			})
			return
		}
		user.Password = req.Password
		if err := user.Update(true); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to reset password",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Password reset successfully",
		})
		return
	case "add_quota":
		if req.Value <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Quota value must be positive",
			})
			return
		}
		if err := model.IncreaseUserQuota(user.Id, req.Value); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to update quota",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Quota updated",
		})
		return
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid action",
		})
		return
	}

	if err := user.Update(false); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User updated",
		"data": gin.H{
			"role":   user.Role,
			"status": user.Status,
		},
	})
}

// GenerateAccessToken generates a new access token for the current user
func GenerateAccessToken(c *gin.Context) {
	id, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Not authenticated",
		})
		return
	}

	user, err := model.GetUserById(id.(int), true)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "User not found",
		})
		return
	}

	// Generate random access token
	token := common.GetRandomString(32)
	user.SetAccessToken(token)

	if err := user.Update(false); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to generate access token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    user.AccessToken,
	})
}

// GetAffCode returns the current user's affiliate code
func GetAffCode(c *gin.Context) {
	id, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Not authenticated",
		})
		return
	}

	user, err := model.GetUserById(id.(int))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "User not found",
		})
		return
	}

	if user.AffCode == "" {
		user.AffCode = common.GetRandomString(4)
		if err := user.Update(false); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Failed to generate affiliate code",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    user.AffCode,
	})
}

