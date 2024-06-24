package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	_ "golang.org/x/crypto/bcrypt"
	"img_hosting/dao"
	"img_hosting/models"
	"img_hosting/pkg/logger"
	"img_hosting/services" // 导入 services 包
	"net/http"
)

func RegisterUser(c *gin.Context) {
	log := logger.GetLogger() //必须实例化

	var userInput models.UserInput
	if err := c.ShouldBindJSON(&userInput); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			c.JSON(http.StatusBadRequest, gin.H{"errors": services.FormatValidationError(validationErrors)})
			log.WithFields(logrus.Fields{"error": err.Error()}).Error("验证失败")

			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//检查用户是否存在
	isEmpty, err := services.IsUserEmpty(userInput.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check user existence"})
		return
	}

	if !isEmpty {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}
	//生成哈希密码
	hashedPassword, err := services.HashPassword(userInput.Psd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	fmt.Println(hashedPassword)
	userInput.Psd = hashedPassword
	db := models.GetDB()
	dao.CreateUser(db, userInput.Name, userInput.Age, hashedPassword)
	log.WithFields(logrus.Fields{"Name": userInput.Name}).Info("用户成功注册")

	c.JSON(http.StatusOK, gin.H{"user": userInput})
}

/*
curl -X POST ^
http://127.0.0.1:8080/register ^
-H "Content-Type: application/json" ^
-d "{\"name\": \"mimi23\", \"age\": 20, \"psd\": \"5d65466678@\"}"

*/
