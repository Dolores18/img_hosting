package controllers

import (
	"fmt"
	"img_hosting/models"
	"img_hosting/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserController 用户控制器
type UserController struct {
	userService *services.UserService
}

func NewUserController() *UserController {
	return &UserController{
		userService: &services.UserService{},
	}
}

// GetProfile godoc
// @Summary 获取用户个人信息
// @Description 获取当前登录用户的详细信息
// @Tags 用户管理
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.Response{data=models.UserInfo}
// @Failure 401 {object} models.Response
// @Router /users/profile [get]
func (uc *UserController) GetProfile(c *gin.Context) {
	// 从 JWT 中获取用户 ID
	userID, exists := c.Get("user_id")
	if !exists {
		fmt.Printf("未找到用户ID\n")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权访问"})
		return
	}

	fmt.Printf("获取到用户ID: %v\n", userID)

	// 确保类型转换正确
	uid, ok := userID.(uint)
	if !ok {
		fmt.Printf("用户ID类型转换失败: %T\n", userID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "用户ID类型错误"})
		return
	}

	user, err := uc.userService.GetUserProfile(uid)
	if err != nil {
		fmt.Printf("获取用户信息失败: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// UpdateProfileRequest 更新用户信息请求
type UpdateProfileRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Age      int    `json:"age"`
}

// UpdateProfile godoc
// @Summary 更新用户信息
// @Description 更新当前登录用户的个人信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param request body UpdateProfileRequest true "用户信息更新请求"
// @Security BearerAuth
// @Success 200 {object} models.Response
// @Failure 400,401 {object} models.Response
// @Router /users/profile [put]
func (uc *UserController) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	var updates map[string]interface{}

	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	if err := uc.userService.UpdateUserProfile(userID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新用户信息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户信息更新成功"})
}

// ListUsersResponse 用户列表响应
type ListUsersResponse struct {
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
	Users    []models.UserInfo `json:"users"`
}

// ListUsers godoc
// @Summary 获取用户列表
// @Description 获取系统中的用户列表，支持分页和搜索
// @Tags 用户管理
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param search query string false "搜索关键词"
// @Security BearerAuth
// @Success 200 {object} ListUsersResponse
// @Failure 401,403 {object} models.Response
// @Router /users [get]
func (uc *UserController) ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	search := c.Query("search")

	users, total, err := uc.userService.ListUsers(page, pageSize, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users":     users,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// DeleteUser godoc
// @Summary 删除用户
// @Description 删除指定的用户
// @Tags 用户管理
// @Produce json
// @Param id path int true "用户ID"
// @Security BearerAuth
// @Success 200 {object} models.Response
// @Failure 400,401,403 {object} models.Response
// @Router /users/{id} [delete]
func (uc *UserController) DeleteUser(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	if err := uc.userService.DeleteUser(uint(userID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除用户失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户已删除"})
}

// UpdateStatus godoc
// @Summary 更新用户状态
// @Description 更新指定用户的状态（激活/禁用）
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param request body models.StatusUpdateRequest true "状态更新请求"
// @Security BearerAuth
// @Success 200 {object} models.Response
// @Failure 400,401,403 {object} models.Response
// @Router /users/{id}/status [put]
func (uc *UserController) UpdateStatus(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	if err := uc.userService.UpdateUserStatus(uint(userID), req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户状态已更新"})
}

// ManageRoles godoc
// @Summary 管理用户角色
// @Description 为指定用户分配角色
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param request body models.RoleAssignRequest true "角色分配请求"
// @Security BearerAuth
// @Success 200 {object} models.Response
// @Failure 400,401,403 {object} models.Response
// @Router /users/{id}/roles [post]
func (uc *UserController) ManageRoles(c *gin.Context) {
	fmt.Println("开始处理管理角色请求")

	// 获取用户ID
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		fmt.Printf("用户ID解析失败: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}
	fmt.Printf("目标用户ID: %d\n", userID)

	var req struct {
		Roles []string `json:"roles" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("请求数据绑定失败: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	fmt.Printf("请求的角色列表: %v\n", req.Roles)

	// 调用 service 层处理角色分配
	if err := uc.userService.AssignRoles(uint(userID), req.Roles); err != nil {
		fmt.Printf("角色分配失败: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("角色分配成功: userID=%d, roles=%v\n", userID, req.Roles)
	c.JSON(http.StatusOK, gin.H{
		"message": "角色分配成功",
		"user_id": userID,
		"roles":   req.Roles,
	})
}

// GetRoles godoc
// @Summary 获取用户角色
// @Description 获取指定用户的所有角色
// @Tags 用户管理
// @Produce json
// @Param id path int true "用户ID"
// @Security BearerAuth
// @Success 200 {object} models.Response{data=[]string}
// @Failure 400,401,403 {object} models.Response
// @Router /users/{id}/roles [get]
func (uc *UserController) GetRoles(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	roles, err := uc.userService.GetUserRoles(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户角色失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"roles": roles})
}
