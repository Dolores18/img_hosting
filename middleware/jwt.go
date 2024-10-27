package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"img_hosting/pkg/logger"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Define a secret key for signing the JWT tokens
var jwtKey = []byte("my_secret_key")

// Claims defines the structure for JWT claims
type Claims struct {
	Username string `json:"username"`
	UserID   uint   `json:"user_id"`
	jwt.StandardClaims
}

// GenerateJWT generates a JWT token for a given username
func GenerateJWT(userid uint) (string, error) {
	// Set expiration time for the token
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		UserID: userid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	// Create the token with the specified claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Sign the token with the secret key
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

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
	"/imgupload":                 true,
	"/statics/html/index.html":   true,
	"/favicon.ico":               true,
}

// AuthMiddleware is a middleware for JWT authentication
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.GetLogger()
		/*
				// 获取请求的 IP 地址
				clientIP := c.ClientIP()

				// 检查请求路径是否在排除的路径中，或者请求来自本地 IP
			if _, ok := ExcludedPaths[c.Request.URL.Path]; ok || clientIP == "127.0.0.1" || clientIP == "::1" {
					c.Next()
					return
				}
		*/
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
		//提取user_id
		// 从 token 中提取用户 ID
		userIDFloat, ok := claims["user_id"].(float64)
		//println("从token中提取的用户id是", userIDFloat)

		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token or invalid type"})
			c.Abort()
			return
		}

		// 将 float64 转换为 uint
		userID := uint(userIDFloat)
		//fmt.Printf("用户id为:%d\n", userID)

		// 将用户 ID 存储在 Gin 上下文中
		c.Set("user_id", userID)

		log.WithField("username", claims["username"]).WithField("user_id", userID).Info("User authenticated with JWT")
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

// ParseAndValidateToken parses and validates the JWT token from the request header
func ParseAndValidateToken(c *gin.Context) (*Claims, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return nil, fmt.Errorf("authorization header is missing")
	}

	tokenString := strings.Split(authHeader, " ")[1]
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil
}
