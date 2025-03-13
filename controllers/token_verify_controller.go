package controllers

import (
	"fmt"
	"img_hosting/models"
	"img_hosting/services"
	"net/http"
	"strings"
	"sync"
	"time"

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

// 缓存相关代码
var (
	// 用户信息缓存
	tokenCache = make(map[string]*models.UserInfo)
	// 文件访问权限缓存 - token+path -> 是否有权限
	accessCache = make(map[string]bool)
	// 缓存互斥锁
	cacheMutex = &sync.RWMutex{}
	// 缓存过期时间 (10分钟)
	cacheExpiration = 10 * time.Minute
)

// VerifyToken godoc
// @Summary 验证Token
// @Description 验证token并返回用户权限信息
// @Tags 认证
// @Accept json
// @Produce json
// @Param token query string false "访问Token"
// @Param Authorization header string false "Bearer Token"
// @Param X-Token header string false "Token"
// @Param path query string true "请求的文件路径"
// @Success 200 {object} models.TokenVerifyResponse
// @Failure 401 {object} models.Response
// @Failure 403 {object} models.Response
// @Router /api/verify-token [get]
func (tc *TokenVerifyController) VerifyToken(c *gin.Context) {
	// 从查询参数获取token
	token := c.Query("token")

	// 如果查询参数中没有token，尝试从Authorization头获取
	if token == "" {
		authHeader := c.GetHeader("Authorization")
		// 支持 "Bearer token" 格式
		if len(authHeader) > 7 && strings.ToUpper(authHeader[0:7]) == "BEARER " {
			token = authHeader[7:]
		} else {
			// 也支持直接传递token
			token = authHeader
		}
	}

	// 从X-Token头获取（Nginx可能使用这种方式传递）
	if token == "" {
		token = c.GetHeader("X-Token")
	}

	// 获取路径
	path := c.Query("path")
	if path == "" {
		path = c.GetHeader("X-Path")
		// 如果还是空，尝试从X-Original-URI获取
		if path == "" {
			path = c.GetHeader("X-Original-URI")
		}
	}
	/*
		// 如果没有token，尝试验证URL签名
		if token == "" {
			signature := c.Query("signature")
			expires := c.Query("expires")

			if signature != "" && expires != "" {
				// 验证签名和过期时间
				if tc.tokenService.ValidateURLSignature(path, signature, expires) {
					c.JSON(http.StatusOK, gin.H{
						"status": "valid",
						"access_type": "signed_url",
					})
					return
				}
			}

			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "无效的访问凭证",
			})
			return
		}

		if token == "" || path == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "缺少token或path参数",
			})
			return
		}
	*/
	// 检查访问权限缓存
	accessCacheKey := token + ":" + path
	cacheMutex.RLock()
	if hasAccess, exists := accessCache[accessCacheKey]; exists {
		// 从用户缓存获取用户信息
		cachedUser, userExists := tokenCache[token]
		cacheMutex.RUnlock()

		if hasAccess && userExists {
			// 使用缓存的权限和用户信息
			c.JSON(http.StatusOK, gin.H{
				"status":  "valid",
				"user_id": cachedUser.UserID,
				"cached":  true,
			})
			c.Header("X-User-ID", fmt.Sprintf("%d", cachedUser.UserID))
			c.Header("X-User-Name", cachedUser.Name)
			c.Header("X-Cache-Hit", "true")
			return
		}
	} else {
		cacheMutex.RUnlock()
	}

	// 缓存未命中，验证token并获取用户信息
	user, err := tc.tokenService.ValidateToken(token)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return
	}

	// 检查文件访问权限
	hasAccess, err := tc.tokenService.CheckFileAccess(user.UserID, path)
	if err != nil || !hasAccess {
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

	// 将结果加入缓存
	cacheMutex.Lock()
	tokenCache[token] = user
	accessCache[accessCacheKey] = true

	// 设置缓存过期
	go func(tokenKey, accessKey string) {
		time.Sleep(cacheExpiration)
		cacheMutex.Lock()
		delete(tokenCache, tokenKey)
		delete(accessCache, accessKey)
		cacheMutex.Unlock()
	}(token, accessCacheKey)

	cacheMutex.Unlock()

	// 返回成功
	c.JSON(http.StatusOK, gin.H{
		"status":  "valid",
		"user_id": user.UserID,
		"cached":  false,
	})
	c.Header("X-User-ID", fmt.Sprintf("%d", user.UserID))
	c.Header("X-User-Name", user.Name)
	c.Header("X-Cache-Hit", "true")
}

// ClearTokenCache 清除指定token的缓存
func ClearTokenCache(token string) {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// 删除token用户缓存
	delete(tokenCache, token)

	// 删除相关的访问权限缓存
	for key := range accessCache {
		if strings.HasPrefix(key, token+":") {
			delete(accessCache, key)
		}
	}
}

// ClearAllCache 清除所有缓存
func ClearAllCache() {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	// 重置缓存
	tokenCache = make(map[string]*models.UserInfo)
	accessCache = make(map[string]bool)
}
