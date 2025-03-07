package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// 保留原有的验证函数
func SignValid(fl validator.FieldLevel) bool {
	name := fl.Field().Interface().(string)
	return !regexp.MustCompile(`[\W_]`).MatchString(name)
}

func PasswordValid(fl validator.FieldLevel) bool {
	password := fl.Field().Interface().(string)
	if len(password) < 8 || len(password) > 20 {
		return false
	}
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	hasSpecial := regexp.MustCompile(`[\W_]`).MatchString(password)
	return hasLetter && hasNumber && hasSpecial
}
