package controllers

import (
	"github.com/gin-gonic/gin"
	_ "gorm.io/gorm"
	"img_hosting/dao"
	"img_hosting/middleware"
	"img_hosting/models"
	"img_hosting/pkg/logger"
	"img_hosting/services"
	"net/http"
	"regexp"
)

func isEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}
func Sinngin(c *gin.Context) {
	db := models.GetDB()
	log := logger.GetLogger()

	var loginRequest struct {
		Identifier struct {
			Identifier string `json:"identifier"`
		} `json:"identifier"`
		Psd string `json:"psd"`
	}

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		log.Error("解析JSON失败: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "解析JSON失败"})
		return
	}

	log.Info("接收到的登录请求: ", loginRequest)

	identifierValue := loginRequest.Identifier.Identifier

	psd := loginRequest.Psd
	println("密码是：%s", psd)

	if identifierValue == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名或邮箱不能为空"})
		return
	}

	if psd == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "密码不能为空"})
		return
	}

	log.Info("用户登录尝试，标识符：", identifierValue)
	log.Info("输入格式正确")
	log.Info(identifierValue)
	var identifierField string
	if isEmail(identifierValue) {
		identifierField = "email"
	} else {
		identifierField = "name"
	}
	user, err := dao.FindUserByFields(db, map[string]interface{}{identifierField: identifierValue})
	queryParams := map[string]interface{}{identifierField: identifierValue}

	// 打印查询参数
	log.Printf("查询参数: %+v", queryParams)
	println(identifierField, identifierValue)
	println("查询参数: %+v", queryParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "没有找到密码"})
		log.Info("没有找到密码")

		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		log.Info("用户不存在")

		return
	}

	err = services.CheckPasswordHash(psd, user.Psd)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户密码错误"})
		return
	}

	// 生成 JWT 令牌
	token, err := middleware.GenerateJWT(user.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "用户认证成功",
		"token":     token,
		"userid":    user.UserID,
		"user_name": user.Name,
	})
}

/*
curl -X POST ^
http://127.0.0.1:8080/searchimg ^
-H "Content-Type: application/json" ^
-d "{\"name\": \"m3\", \"age\":25\"psd\": \"5d65466678@\"}"

*/
//"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Im1pbWkyMyIsImV4cCI6MTcxOTM1NTU1NX0.oesiJ5wkaSKQtjAuP3vJzK-EYUMdfbKEFL0hWK6HOSg"
