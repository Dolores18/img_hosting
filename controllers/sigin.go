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
)

func IsExist(name string) (bool, error) {
	db := models.GetDB()
	user, err := dao.FindUserByName(db, name)
	if err != nil {
		// 返回错误
		return false, err
	}
	if user == nil {
		return false, nil
	}
	return true, nil
}

func Sigin(c *gin.Context) {
	db := models.GetDB()
	log := logger.GetLogger() //必须实例化

	var userInput models.UserInput
	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errs": err.Error()})
		log.Info("解析 JSON 失败")
		return
	}

	isExist, err := IsExist(userInput.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errs": "读取数据库失败"})
		return
	}

	if !isExist {
		c.JSON(http.StatusNotFound, gin.H{"errs": "用户不存在"})
		return
	}

	user, err := dao.FindUserByName(db, userInput.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errs": "没有找到密码"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"err": "用户不存在"})
		return
	}

	err = services.CheckPasswordHash(userInput.Psd, user.Psd)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"errs": "用户密码错误"})
		return
	}

	// 生成 JWT 令牌
	token, err := middleware.GenerateJWT(user.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "用户认证成功",
		"token":   token,
		"name":    user.Name,
	})
}

/*
curl -X POST ^
http://127.0.0.1:8080/sigin ^
-H "Content-Type: application/json" ^
-d "{\"name\": \"m3\", \"age\":25\"psd\": \"5d65466678@\"}"

*/
//"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6Im1pbWkyMyIsImV4cCI6MTcxOTM1NTU1NX0.oesiJ5wkaSKQtjAuP3vJzK-EYUMdfbKEFL0hWK6HOSg"
