package controllers

import (
	"fmt"
	"img_hosting/middleware"
	"img_hosting/models"
	"img_hosting/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// LoginRequest 登录请求结构
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Psd      string `json:"psd"`
}

// AuthController 认证控制器
type AuthController struct {
	authService *services.AuthService
}

func NewAuthController() *AuthController {
	return &AuthController{
		authService: &services.AuthService{},
	}
}

// Login godoc
// @Summary 用户登录
// @Description 处理用户登录请求，返回JWT令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录请求"
// @Success 200 {object} models.Response{data=models.LoginResponse}
// @Failure 400 {object} models.Response "无效的请求数据"
// @Failure 401 {object} models.Response "认证失败"
// @Router /auth/login [post]
func (ac *AuthController) Login(c *gin.Context) {
	fmt.Println("开始处理登录请求")
	var loginReq LoginRequest

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		fmt.Printf("请求数据绑定失败: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}

	fmt.Printf("登录请求原始数据: email=%s, password=%s, psd=%s\n",
		loginReq.Email, loginReq.Password, loginReq.Psd)

	// 优先使用 password 字段，如果为空则使用 psd 字段
	password := loginReq.Password
	if password == "" {
		password = loginReq.Psd
		fmt.Println("使用 psd 字段作为密码")
	} else {
		fmt.Println("使用 password 字段作为密码")
	}

	fmt.Printf("最终使用的密码: %s (长度: %d)\n", password, len(password))

	user, err := ac.authService.Login(loginReq.Email, password)
	if err != nil {
		fmt.Printf("登录失败: %v\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 生成 JWT token
	token, err := middleware.GenerateJWT(user.UserID)
	if err != nil {
		fmt.Printf("生成令牌失败: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}

	fmt.Printf("登录成功: userID=%d\n", user.UserID)
	c.JSON(http.StatusOK, gin.H{
		"message":   "登录成功",
		"token":     token,
		"user_id":   user.UserID,
		"user_name": user.Name,
	})
}

// Register godoc
// @Summary 用户注册
// @Description 处理新用户注册请求
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body models.UserInput true "注册信息"
// @Success 201 {object} models.Response{data=models.RegisterResponse}
// @Failure 400 {object} models.Response "注册失败"
// @Router /auth/register [post]
func (ac *AuthController) Register(c *gin.Context) {
	fmt.Println("开始处理注册请求")
	var input models.UserInput

	if err := c.ShouldBindJSON(&input); err != nil {
		fmt.Printf("请求数据绑定失败: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据"})
		return
	}
	fmt.Printf("注册信息: name=%s, email=%s\n", input.Name, input.Email)

	user, err := ac.authService.Register(&input)
	if err != nil {
		fmt.Printf("注册失败: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("注册成功: userID=%d\n", user.UserID)
	c.JSON(http.StatusCreated, gin.H{
		"message": "注册成功",
		"user_id": user.UserID,
	})
}
