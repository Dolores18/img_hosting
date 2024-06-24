package controllers

import (
	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	_ "golang.org/x/crypto/bcrypt"
	"img_hosting/models"
	"img_hosting/pkg/logger"
	"img_hosting/services" // 导入 services 包
	"net/http"
)

func RegisterUser(c *gin.Context) {
	log := logger.GetLogger() //必须实例化

	var user models.UserInfo
	if err := c.ShouldBindJSON(&user); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			c.JSON(http.StatusBadRequest, gin.H{"errors": services.FormatValidationError(validationErrors)})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := services.HashPassword(user.Psd)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	log.WithFields(logrus.Fields{"Name": user.Name}).Info("用户成功注册")

	c.JSON(http.StatusOK, gin.H{"user": user})
}

/*
curl -X POST ^
http://127.0.0.1:8080/register ^
-H "Content-Type: application/json" ^
-d "{\"name\": \"fengfen\", \"age\": 20, \"psd\": \"5d65466678@\"}"

*/
