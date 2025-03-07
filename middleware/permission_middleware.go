package middleware

import (
	"img_hosting/dao"
	"net/http"
	"strings"

	"img_hosting/config"

	"fmt"

	"github.com/gin-gonic/gin"
)

// 定义路由到权限的映射
var routePermissions = config.AppConfigInstance.Permissions.Routes

// PermissionMiddleware 创建一个检查用户权限的中间件
func PermissionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fullPath := c.FullPath()
		fmt.Printf("权限中间件 - 当前路径: %s\n", fullPath)

		// 检查是否是公开路由
		if permissions, exists := config.AppConfigInstance.Permissions.Routes[fullPath]; exists {
			fmt.Printf("权限检查: path=%s, permissions=%v, exists=%v\n", fullPath, permissions, exists)
			if len(permissions) == 0 {
				fmt.Println("公开路由，无需权限检查")
				c.Next()
				return
			}
		}

		// 获取用户ID
		userID, exists := c.Get("user_id")
		fmt.Printf("用户ID检查: userID=%v, exists=%v\n", userID, exists)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未认证的用户"})
			c.Abort()
			return
		}

		// 获取所需权限
		requiredPermissions := getRequiredPermissions(fullPath)
		fmt.Printf("所需权限: %v\n", requiredPermissions)
		if len(requiredPermissions) == 0 {
			c.Next()
			return
		}

		// 检查用户权限
		permissions, err := dao.UserHasPermissions(userID.(uint), requiredPermissions)
		fmt.Printf("用户权限检查结果: permissions=%v, err=%v\n", permissions, err)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "权限检查失败"})
			c.Abort()
			return
		}

		// 验证权限
		for _, perm := range requiredPermissions {
			if !permissions[perm] {
				c.JSON(http.StatusForbidden, gin.H{
					"error":               "权限不足",
					"required_permission": perm,
				})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// getRequiredPermissions 根据路由路径获取所需权限
func getRequiredPermissions(path string) []string {
	// 处理带参数的路径
	if strings.Contains(path, "/users/") {
		if strings.HasSuffix(path, "/status") {
			return routePermissions["/users/:id/status"]
		}
		if strings.HasSuffix(path, "/roles") {
			return routePermissions["/users/:id/roles"]
		}
		if strings.Count(path, "/") == 2 {
			return routePermissions["/users/:id"]
		}
	}

	if strings.Contains(path, "/private-files/") {
		if strings.HasSuffix(path, "/download") {
			return routePermissions["/private-files/:id/download"]
		}
		return routePermissions["/private-files/:id"]
	}

	return routePermissions[path]
}
