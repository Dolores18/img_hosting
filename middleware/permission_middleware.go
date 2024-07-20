package middleware

import (
	"github.com/gin-gonic/gin"
	"img_hosting/dao"
	"net/http"
)

// 定义路由到权限的映射
var routePermissions = map[string][]string{
	"/imgupload":   {"upload_img"},
	"/user_manage": {"usermanage"},
	// 可以定义需要多个权限的路由
	"/admin_panel": {"admin", "view_panel"},
}

// PermissionMiddleware 创建一个检查用户权限的中间件
func PermissionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 1. 从上下文中获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		// 2. 获取当前路由所需的权限
		requiredPermissions := getRequiredPermissions(c.FullPath())
		if len(requiredPermissions) == 0 {
			c.Next()
			return
		}

		// 3. 检查用户权限
		permissions, err := dao.UserHasPermissions(userID.(uint), requiredPermissions)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check permissions"})
			c.Abort()
			return
		}

		// 4. 验证是否拥有所有必需的权限
		for _, perm := range requiredPermissions {
			if !permissions[perm] {
				c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

// getRequiredPermissions 根据路由路径获取所需权限
func getRequiredPermissions(path string) []string {
	return routePermissions[path]
}
