package middleware

import (
	"fmt"
	"img_hosting/config"
	"img_hosting/pkg/logger"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Define a secret key for signing the JWT tokens
var jwtKey = []byte("my_secret_key")

// Claims defines the structure for JWT claims
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

// GenerateJWT generates a JWT token for a given username
func GenerateJWT(userID uint) (string, error) {
	fmt.Printf("开始生成JWT令牌: userID=%d\n", userID)

	// Set expiration time for the token
	expirationTime := time.Now().Add(30 * 24 * time.Hour)
	claims := &Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Create the token with the specified claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		fmt.Printf("生成JWT令牌失败: %v\n", err)
		return "", err
	}

	fmt.Printf("JWT令牌生成成功: token=%s\n", tokenString)
	return tokenString, nil
}

// ExcludedPaths stores the paths that do not require authentication
var ExcludedPaths = map[string]bool{
	"/":                          true,
	"/statics/html/":             true,
	"/signin":                    true,
	"/login":                     true,
	"/register":                  true,
	"/statics/css/":              true,
	"/statics/":                  true,
	"/statics/css/style.css":     true,
	"/statics/imgs/example.jpg ": true,
	"/statics/html/index.html":   true,
	"/favicon.ico":               true,
}

// AuthMiddleware is a middleware for JWT authentication
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.GetLogger()
		path := c.FullPath()
		fmt.Printf("当前访问路径: %s\n", path)

		// 验证 Authorization header
		authHeader := c.GetHeader("Authorization")
		fmt.Printf("Authorization头: %s\n", authHeader)

		if authHeader == "" {
			fmt.Println("缺少认证信息")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少认证信息"})
			c.Abort()
			return
		}

		// 解析 token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的认证格式"})
			c.Abort()
			return
		}

		claims, err := ParseAndValidateToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的 token"})
			c.Abort()
			return
		}

		// 设置用户信息到上下文
		c.Set("user_id", claims.UserID)
		log.WithField("user_id", claims.UserID).Info("用户认证成功")

		// 检查是否是公开路由
		if permissions, exists := config.AppConfigInstance.Permissions.Routes[path]; exists {
			fmt.Printf("路由权限检查: path=%s, permissions=%v, exists=%v\n", path, permissions, exists)
			if len(permissions) == 0 {
				fmt.Println("公开路由，直接放行")
			}
		}

		c.Next()
	}
}

func validateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
}

// ParseAndValidateToken parses and validates the JWT token from the request header
func ParseAndValidateToken(c *gin.Context) (*Claims, error) {
	authHeader := c.GetHeader("Authorization")
	fmt.Printf("验证令牌 - Authorization头: %s\n", authHeader)

	if authHeader == "" {
		return nil, fmt.Errorf("authorization header is missing")
	}

	tokenString := strings.Split(authHeader, " ")[1]
	fmt.Printf("验证令牌 - 令牌字符串: %s\n", tokenString)

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil {
		fmt.Printf("令牌解析失败: %v\n", err)
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		fmt.Printf("令牌声明无效: ok=%v, valid=%v\n", ok, token.Valid)
		return nil, fmt.Errorf("invalid token claims")
	}

	fmt.Printf("令牌验证成功: userID=%d\n", claims.UserID)
	return claims, nil
}
