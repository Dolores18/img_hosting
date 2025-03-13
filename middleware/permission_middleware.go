package middleware

import (
	"img_hosting/dao"
	"img_hosting/services"
	"net/http"
	"strings"

	"img_hosting/config"

	"fmt"

	"img_hosting/pkg/cache"

	"github.com/gin-gonic/gin"
)

// 定义路由到权限的映射
var routePermissions = config.AppConfigInstance.Permissions.Routes

// PermissionMiddleware 创建一个检查用户权限的中间件
func PermissionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fullPath := c.FullPath()
		method := c.Request.Method
		fmt.Printf("权限中间件 - 当前路径: %s, 请求方法: %s\n", fullPath, method)

		// 检查是否是公开路由
		routePermissions := config.AppConfigInstance.Permissions.Routes
		if permissions, exists := routePermissions[fullPath]; exists {
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
		requiredPermissions := getRequiredPermissions(fullPath, method)
		fmt.Printf("所需权限: %v\n", requiredPermissions)

		// 特殊处理用户角色管理路由
		if len(requiredPermissions) == 0 && strings.Contains(fullPath, "/users/") && strings.HasSuffix(fullPath, "/roles") {
			requiredPermissions = []string{"manage_user_roles"}
			fmt.Printf("特殊处理用户角色路由，设置所需权限: %v\n", requiredPermissions)
		}

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
func getRequiredPermissions(path string, method string) []string {
	fmt.Printf("获取路径权限: path=%s, method=%s\n", path, method)
	routePermissions := config.AppConfigInstance.Permissions.Routes

	// 1. 尝试匹配大写方法的路由
	methodPath := fmt.Sprintf("%s %s", method, path)
	if perms, exists := routePermissions[methodPath]; exists {
		fmt.Printf("匹配到大写方法路由: %s, 权限: %v\n", methodPath, perms)
		return perms
	}

	// 2. 尝试匹配小写方法的路由
	methodPathLower := fmt.Sprintf("%s %s", strings.ToLower(method), path)
	if perms, exists := routePermissions[methodPathLower]; exists {
		fmt.Printf("匹配到小写方法路由: %s, 权限: %v\n", methodPathLower, perms)
		return perms
	}

	// 打印配置中的所有路由权限
	fmt.Println("配置中的路由权限:")
	for routePath, perms := range routePermissions {
		fmt.Printf("  %s: %v\n", routePath, perms)
	}

	// 处理带参数的路径
	if strings.Contains(path, "/users/") {
		if strings.HasSuffix(path, "/status") {
			perms := routePermissions["/users/:id/status"]
			fmt.Printf("匹配用户状态路由: path=%s, permissions=%v\n", path, perms)
			return perms
		}
		if strings.HasSuffix(path, "/roles") {
			perms := routePermissions["/users/:id/roles"]
			fmt.Printf("匹配用户角色路由: path=%s, permissions=%v\n", path, perms)
			return perms
		}
		if strings.Count(path, "/") == 2 {
			perms := routePermissions["/users/:id"]
			fmt.Printf("匹配用户ID路由: path=%s, permissions=%v\n", path, perms)
			return perms
		}
	}

	if strings.Contains(path, "/private-files/") {
		if strings.HasSuffix(path, "/download") {
			perms := routePermissions["/private-files/:id/download"]
			fmt.Printf("匹配文件下载路由: path=%s, permissions=%v\n", path, perms)
			return perms
		}
		perms := routePermissions["/private-files/:id"]
		fmt.Printf("匹配文件ID路由: path=%s, permissions=%v\n", path, perms)
		return perms
	}

	// 直接匹配
	perms := routePermissions[path]
	fmt.Printf("直接匹配: path=%s, permissions=%v\n", path, perms)
	return perms
}

// 获取用户权限（带缓存）
func getUserPermissions(userID uint) (map[string]bool, error) {
	if perms, exists := cache.GetUserPermissionCache(userID); exists {
		return perms, nil
	}

	permService := &services.PermissionService{}
	perms, err := permService.GetUserPermissionMap(userID)
	if err != nil {
		return nil, err
	}

	cache.SetUserPermissionCache(userID, perms)
	return perms, nil
}
