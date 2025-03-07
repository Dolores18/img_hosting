package controllers

import (
	"fmt"
	"img_hosting/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TokenController 处理 token 相关的请求
type TokenController struct{}

// CreateToken 创建新的 token
func (tc *TokenController) CreateToken(c *gin.Context) {
	fmt.Println("开始处理创建令牌请求")

	var req struct {
		DeviceID  string `json:"device_id"`
		IPAddress string `json:"ip_address"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("请求数据绑定失败: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	fmt.Printf("创建令牌请求: deviceID=%s, ipAddress=%s\n",
		req.DeviceID, req.IPAddress)

	userID := c.GetUint("user_id")

	// Device-ID 是可选的
	if req.DeviceID == "" {
		req.DeviceID = c.GetHeader("Device-ID")
	}

	token, err := services.CreateToken(userID, req.DeviceID, req.IPAddress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建 token 失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":      token,
		"expires_at": token.ExpiresAt,
		"device_id":  req.DeviceID, // 返回使用的 device_id
	})
}

// ListTokens 获取用户的所有 token
func (tc *TokenController) ListTokens(c *gin.Context) {
	userID := c.GetUint("user_id")

	tokens, err := services.ListUserTokens(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取 token 列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tokens": tokens})
}

// RevokeToken 撤销指定的 token
func (tc *TokenController) RevokeToken(c *gin.Context) {
	userID := c.GetUint("user_id")
	tokenStr := c.Param("token")

	if err := services.RevokeToken(userID, tokenStr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "撤销 token 失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "token 已撤销"})
}
