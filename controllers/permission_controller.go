package controllers

import (
	"img_hosting/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PermissionController struct {
	permService *services.PermissionService
}

func NewPermissionController() *PermissionController {
	return &PermissionController{
		permService: &services.PermissionService{},
	}
}

// GetAllPermissions 获取所有权限
func (pc *PermissionController) GetAllPermissions(c *gin.Context) {
	permissions, err := pc.permService.GetAllPermissions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取权限列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"permissions": permissions})
}

// GetAllRoles 获取所有角色
func (pc *PermissionController) GetAllRoles(c *gin.Context) {
	roles, err := pc.permService.GetAllRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取角色列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"roles": roles})
}

// GetRolePermissions 获取角色的权限
func (pc *PermissionController) GetRolePermissions(c *gin.Context) {
	roleName := c.Param("role")

	permissions, err := pc.permService.GetRolePermissions(roleName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"role":        roleName,
		"permissions": permissions,
	})
}

// UpdateRolePermissions 更新角色的权限
func (pc *PermissionController) UpdateRolePermissions(c *gin.Context) {
	roleName := c.Param("role")

	var req struct {
		Permissions []string `json:"permissions" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	if err := pc.permService.UpdateRolePermissions(roleName, req.Permissions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "角色权限更新成功",
		"role":        roleName,
		"permissions": req.Permissions,
	})
}

// CreatePermission 创建新权限
func (pc *PermissionController) CreatePermission(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	if err := pc.permService.CreatePermission(req.Name, req.Description); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "权限创建成功",
		"permission": req.Name,
	})
}

// CreateRole 创建新角色
func (pc *PermissionController) CreateRole(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	if err := pc.permService.CreateRole(req.Name, req.Description); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "角色创建成功",
		"role":    req.Name,
	})
}

// SyncConfigPermissions 同步配置文件中的权限到数据库
func (pc *PermissionController) SyncConfigPermissions(c *gin.Context) {
	if err := pc.permService.SyncConfigPermissions(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "同步权限失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "权限同步成功"})
}

// GetUserPermissions 获取用户的所有权限
func (pc *PermissionController) GetUserPermissions(c *gin.Context) {
	// 获取用户ID
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	// 获取用户权限
	permissions, err := pc.permService.GetUserPermissions(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户权限失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":     userID,
		"permissions": permissions,
	})
}

// GetCurrentUserPermissions 获取当前用户的所有权限
func (pc *PermissionController) GetCurrentUserPermissions(c *gin.Context) {
	// 从上下文中获取当前用户ID
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	// 获取用户权限
	permissions, err := pc.permService.GetUserPermissions(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户权限失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":     userID,
		"permissions": permissions,
	})
}
