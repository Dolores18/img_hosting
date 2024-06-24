package services

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"img_hosting/pkg/logger"
	"regexp"
)

// 自定义验证器函数
func signValid(fl validator.FieldLevel) bool {
	name := fl.Field().Interface().(string)
	return name == "fengfeng"
}
func PasswordValid(fl validator.FieldLevel) bool {
	password := fl.Field().Interface().(string)
	// 长度大于8，小于20
	if len(password) < 8 || len(password) > 20 {
		return false
	}
	// 包含字母、数字、特殊符号
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[\W_]`).MatchString(password)
	return hasLetter && hasNumber && hasSpecial
}

func InitValidator() {
	log := logger.GetLogger()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("sign", signValid)
		if err != nil {
			log.Info("自定义用户名注册失败")
			return
		}
		err2 := v.RegisterValidation("password", PasswordValid)
		if err2 != nil {
			log.Info("自定义密码验证注册失败")
			return
		}
	}

}

func FormatValidationError(errs validator.ValidationErrors) map[string]string {
	errors := make(map[string]string)
	for _, err := range errs {
		fieldName := err.Field()
		switch err.Tag() {
		case "required":
			errors[fieldName] = fmt.Sprintf("The %s field is required.", fieldName)
		case "sign":
			errors[fieldName] = fmt.Sprintf("The value of %s field is invalid.", fieldName)
		case "gte":
			errors[fieldName] = fmt.Sprintf("The %s field must be greater than or equal to %s.", fieldName, err.Param())
		case "password":
			errors[fieldName] = fmt.Sprintf("The %s field must be 8-20 characters long and contain letters, numbers, and special characters.", fieldName)
		default:
			errors[fieldName] = fmt.Sprintf("Validation error for %s: %s.", fieldName, err.Tag())
		}
	}
	return errors
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

/*
// RegisterUser 保存用户到数据库
func RegisterUser(user *models.User) error {
	// 实现用户保存的逻辑，例如调用 DAO 层的方法保存用户到数据库
	return dao.SaveUser(user)
}
*/
