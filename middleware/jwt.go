package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"img_hosting/pkg/logger"
)

var jwtKey = []byte("my_secret_key")

// ExcludedPaths stores the paths that do not require authentication
var ExcludedPaths = map[string]bool{
	"/":                          true,
	"/statics/html/":             true,
	"/sigin":                     true,
	"/login":                     true,
	"/register":                  true,
	"/statics/css/":              true,
	"/statics/":                  true,
	"/statics/css/style.css":     true,
	"/statics/imgs/example.jpg ": true,
}

// AuthMiddleware is a middleware for JWT authentication
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.GetLogger()

		// 获取请求的 IP 地址
		clientIP := c.ClientIP()

		// 检查请求路径是否在排除的路径中，或者请求来自本地 IP
		if _, ok := ExcludedPaths[c.Request.URL.Path]; ok || clientIP == "127.0.0.1" || clientIP == "::1" {
			c.Next()
			return
		}

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		tokenString := strings.Split(authHeader, " ")[1]
		token, err := validateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		log.WithField("username", claims["username"]).Info("User authenticated with JWT")
		c.Next()
	}
}

func validateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
}
