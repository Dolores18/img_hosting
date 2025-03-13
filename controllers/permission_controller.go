package controllers

import (
	"img_hosting/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UpdateRolePermissionsRequest 更新角色权限请求
type UpdateRolePermissionsRequest struct {
	Permissions []string `json:"permissions" binding:"required"`
}

// CreatePermissionRequest 创建权限请求
type CreatePermissionRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// PermissionController 权限控制器
type PermissionController struct {
	permService *services.PermissionService
}

func NewPermissionController() *PermissionController {
	return &PermissionController{
		permService: &services.PermissionService{},
	}
}

// GetAllPermissions godoc
// @Summary 获取所有权限
// @Description 获取系统中所有可用的权限列表
// @Tags 权限管理
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.Response{data=[]models.Permissions}
// @Failure 500 {object} models.Response
// @Router /permissions/all [get]
func (pc *PermissionController) GetAllPermissions(c *gin.Context) {
	permissions, err := pc.permService.GetAllPermissions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取权限列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"permissions": permissions})
}

// GetAllRoles godoc
// @Summary 获取所有角色
// @Description 获取系统中所有可用的角色列表
// @Tags 权限管理
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.Response{data=[]string}
// @Failure 500 {object} models.Response
// @Router /permissions/roles [get]
func (pc *PermissionController) GetAllRoles(c *gin.Context) {
	roles, err := pc.permService.GetAllRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取角色列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"roles": roles})
}

// GetRolePermissions godoc
// @Summary 获取角色权限
// @Description 获取指定角色的所有权限
// @Tags 权限管理
// @Produce json
// @Param role path string true "角色名称"
// @Security BearerAuth
// @Success 200 {object} models.Response{data=[]string}
// @Failure 400 {object} models.Response
// @Router /permissions/roles/{role} [get]
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

// UpdateRolePermissions godoc
// @Summary 更新角色权限
// @Description 更新指定角色的权限列表
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param role path string true "角色名称"
// @Param request body UpdateRolePermissionsRequest true "权限列表"
// @Security BearerAuth
// @Success 200 {object} models.Response
// @Failure 400,500 {object} models.Response
// @Router /permissions/roles/{role} [put]
func (pc *PermissionController) UpdateRolePermissions(c *gin.Context) {
	roleName := c.Param("role")

	var req UpdateRolePermissionsRequest

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

// CreatePermission godoc
// @Summary 创建权限
// @Description 创建新的权限项
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param request body CreatePermissionRequest true "权限信息"
// @Security BearerAuth
// @Success 201 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /permissions/create [post]
func (pc *PermissionController) CreatePermission(c *gin.Context) {
	var req CreatePermissionRequest

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

// CreateRole godoc
// @Summary 创建角色
// @Description 创建新的角色
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param request body CreateRoleRequest true "角色信息"
// @Security BearerAuth
// @Success 201 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /permissions/roles/create [post]
func (pc *PermissionController) CreateRole(c *gin.Context) {
	var req CreateRoleRequest

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

// SyncConfigPermissions godoc
// @Summary 同步权限配置
// @Description 从配置文件同步权限到数据库
// @Tags 权限管理
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /permissions/sync [post]
func (pc *PermissionController) SyncConfigPermissions(c *gin.Context) {
	if err := pc.permService.SyncConfigPermissions(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "同步权限失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "权限同步成功"})
}

// GetUserPermissions godoc
// @Summary 获取用户权限
// @Description 获取指定用户的所有权限
// @Tags 权限管理
// @Produce json
// @Param id path int true "用户ID"
// @Security BearerAuth
// @Success 200 {object} models.Response{data=[]string}
// @Failure 400,500 {object} models.Response
// @Router /permissions/users/{id}/permissions [get]
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

// GetCurrentUserPermissions godoc
// @Summary 获取当前用户权限
// @Description 获取当前登录用户的所有权限
// @Tags 权限管理
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.Response{data=[]string}
// @Failure 401,500 {object} models.Response
// @Router /permissions/users/current/permissions [get]
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
