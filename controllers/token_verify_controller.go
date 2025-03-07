package controllers

import (
	"img_hosting/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TokenVerifyController struct {
	tokenService *services.TokenService
}

func NewTokenVerifyController() *TokenVerifyController {
	return &TokenVerifyController{
		tokenService: services.NewTokenService(),
	}
}

// VerifyToken 验证token并返回用户权限
// @Summary 验证token
// @Description 用于Nginx回源验证token的有效性和权限
// @Accept  json
// @Produce json
// @Param   token    query     string  true  "访问Token"
// @Param   path     query     string  true  "请求的文件路径"
// @Success 200 {object} gin.H
// @Failure 401 {object} gin.H
// @Router  /api/verify-token [get]
func (tc *TokenVerifyController) VerifyToken(c *gin.Context) {
	token := c.Query("token")
	path := c.Query("path")

	if token == "" || path == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "缺少token或path参数",
		})
		return
	}

	// 验证token并获取用户信息
	user, err := tc.tokenService.ValidateToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "无效的token",
		})
		return
	}

	// 检查文件访问权限
	if hasAccess, err := tc.tokenService.CheckFileAccess(user.UserID, path); err != nil || !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "没有访问权限",
		})
		return
	}

	// 更新token使用记录
	if err := tc.tokenService.UpdateTokenUsage(token); err != nil {
		// 仅记录错误，不影响验证结果
		c.Error(err)
	}

	// 返回成功
	c.JSON(http.StatusOK, gin.H{
		"status":  "valid",
		"user_id": user.UserID,
	})
}
