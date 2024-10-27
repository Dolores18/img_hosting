package middleware

import (
	"img_hosting/dao"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// 定义路由到权限的映射
var routePermissions = map[string][]string{
	"/imgupload":    {"upload_img"},
	"/user_manage":  {"usermanage"},
	"/admin_panel":  {"admin", "view_panel"},
	"/searchbytag":  {"search_img"}, // 添加搜索图片的权限要求
	"/createtag":    {"createtag"},  // 添加创建标签的权限要求
	"/getalltag":    {"search_img"}, // 添加获取所有标签的权限要求
	"/searchimg":    {"search_img"}, // 添加搜索图片的权限要求
	"/searchAllimg": {"search_img"}, // 改为使用 search_img 权限
}

// PermissionMiddleware 创建一个检查用户权限的中间件
func PermissionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 添加路径调试信息
		fullPath := c.FullPath()
		println("请求路径:", fullPath)

		// 1. 从上下文中获取用户ID
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}
		// 修改用户ID的输出格式
		println("中间件获取的用户id是:", uint64(userID.(uint)))

		// 2. 获取当前路由所需的权限
		requiredPermissions := getRequiredPermissions(fullPath)
		// 修改权限数组的输出格式
		println("需要的权限:", strings.Join(requiredPermissions, ", "))

		// 3. 检查用户权限
		permissions, err := dao.UserHasPermissions(userID.(uint), requiredPermissions)
		if err != nil {
			println("权限检查错误:", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check permissions"})
			c.Abort()
			return
		}
		// 添加权限检查结果的详细输出
		println("用户权限检查结果:")
		for perm, has := range permissions {
			println(perm, ":", has)
		}

		// 4. 验证是否拥有所有必需的权限
		for _, perm := range requiredPermissions {
			if !permissions[perm] {
				c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
				c.Abort()
				return
			}
		}
		println("中间件检查权限完毕")
		c.Next()
		println("中间件执行完毕")
	}
}

// getRequiredPermissions 根据路由路径获取所需权限
func getRequiredPermissions(path string) []string {
	return routePermissions[path]
}
