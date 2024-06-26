package services

import (
	"golang.org/x/crypto/bcrypt"
)

// CheckPasswordHash 验证哈希密码是否与明文密码一致
func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
