package controllers

import (
	"fmt"
	"img_hosting/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// TokenCreateRequest 创建令牌请求
type TokenCreateRequest struct {
	DeviceID  string `json:"device_id"`
	IPAddress string `json:"ip_address"`
}

// TokenController 处理 token 相关的请求
type TokenController struct{}

// CreateToken godoc
// @Summary 创建新的访问令牌
// @Description 为当前用户创建一个新的访问令牌
// @Tags 令牌管理
// @Accept json
// @Produce json
// @Param request body TokenCreateRequest true "令牌创建请求"
// @Security BearerAuth
// @Success 200 {object} models.Response{data=models.Token}
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/tokens [post]
func (tc *TokenController) CreateToken(c *gin.Context) {
	fmt.Println("开始处理创建令牌请求")

	var req TokenCreateRequest

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

// ListTokens godoc
// @Summary 获取用户的所有令牌
// @Description 获取当前用户的所有访问令牌列表
// @Tags 令牌管理
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.Response{data=[]models.Token}
// @Failure 500 {object} models.Response
// @Router /api/tokens [get]
func (tc *TokenController) ListTokens(c *gin.Context) {
	userID := c.GetUint("user_id")

	tokens, err := services.ListUserTokens(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取 token 列表失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"tokens": tokens})
}

// RevokeToken godoc
// @Summary 撤销指定令牌
// @Description 撤销指定的访问令牌
// @Tags 令牌管理
// @Produce json
// @Param token path string true "要撤销的令牌"
// @Security BearerAuth
// @Success 200 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/tokens/{token} [delete]
func (tc *TokenController) RevokeToken(c *gin.Context) {
	userID := c.GetUint("user_id")
	tokenStr := c.Param("token")

	if err := services.RevokeToken(userID, tokenStr); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "撤销 token 失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "token 已撤销"})
}
