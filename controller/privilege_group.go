package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/fangxiusun/tokenhub/common"
	"github.com/fangxiusun/tokenhub/model"
)

// GetAllPrivilegeGroups returns all privilege groups
func GetAllPrivilegeGroups(c *gin.Context) {
	groups, err := model.GetAllPrivilegeGroups()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to get privilege groups",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    groups,
	})
}

// GetPrivilegeGroup returns a specific privilege group
func GetPrivilegeGroup(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid ID",
		})
		return
	}

	group, err := model.GetPrivilegeGroupById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Privilege group not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    group,
	})
}

// CreatePrivilegeGroup creates a new privilege group
func CreatePrivilegeGroup(c *gin.Context) {
	var group model.PrivilegeGroup
	if err := common.DecodeJson(c.Request.Body, &group); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
		})
		return
	}

	if group.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Name is required",
		})
		return
	}

	// Check if name already exists
	existing, _ := model.GetPrivilegeGroupByName(group.Name)
	if existing != nil {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": "Privilege group name already exists",
		})
		return
	}

	if err := model.CreatePrivilegeGroup(&group); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to create privilege group",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    group,
	})
}

// UpdatePrivilegeGroup updates a privilege group
func UpdatePrivilegeGroup(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid ID",
		})
		return
	}

	existing, err := model.GetPrivilegeGroupById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Privilege group not found",
		})
		return
	}

	var updates model.PrivilegeGroup
	if err := common.DecodeJson(c.Request.Body, &updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request body",
		})
		return
	}

	// Update fields
	if updates.DisplayName != "" {
		existing.DisplayName = updates.DisplayName
	}
	if updates.Description != "" {
		existing.Description = updates.Description
	}
	if updates.Quota > 0 {
		existing.Quota = updates.Quota
	}
	if updates.RateLimit > 0 {
		existing.RateLimit = updates.RateLimit
	}
	if updates.Models != "" {
		existing.Models = updates.Models
	}
	if updates.Status != 0 {
		existing.Status = updates.Status
	}

	if err := model.UpdatePrivilegeGroup(existing); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to update privilege group",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    existing,
	})
}

// DeletePrivilegeGroup deletes a privilege group
func DeletePrivilegeGroup(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid ID",
		})
		return
	}

	// Prevent deleting default group
	group, err := model.GetPrivilegeGroupById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "Privilege group not found",
		})
		return
	}

	if group.Name == "default" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Cannot delete default privilege group",
		})
		return
	}

	if err := model.DeletePrivilegeGroup(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to delete privilege group",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Privilege group deleted",
	})
}

