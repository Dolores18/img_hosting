package dao

import (
	"img_hosting/models"
	"time"

	"gorm.io/gorm"
)

// CreateToken 创建新的 token
func CreateToken(db *gorm.DB, token string, userID uint, expiresAt time.Time, deviceID, ipAddress string) error {
	tokenModel := models.Token{
		Token:      token,
		UserID:     userID,
		ExpiresAt:  expiresAt,
		DeviceID:   deviceID,
		IPAddress:  ipAddress,
		Status:     "active",
		LastUsedAt: time.Now(),
	}
	return db.Create(&tokenModel).Error
}

// GetTokenByString 通过 token 字符串获取 token
func GetTokenByString(db *gorm.DB, token string) (*models.Token, error) {
	var tokenModel models.Token
	err := db.Where("token = ? AND status = ?", token, "active").First(&tokenModel).Error
	if err != nil {
		return nil, err
	}
	return &tokenModel, nil
}

// ListUserTokens 获取用户的所有 token
func ListUserTokens(db *gorm.DB, userID uint) ([]models.Token, error) {
	var tokens []models.Token
	err := db.Where("user_id = ?", userID).Find(&tokens).Error
	return tokens, err
}

// RevokeToken 撤销 token
func RevokeToken(db *gorm.DB, token string) error {
	return db.Model(&models.Token{}).
		Where("token = ?", token).
		Update("status", "revoked").Error
}

// UpdateTokenLastUsed 更新 token 最后使用时间
func UpdateTokenLastUsed(db *gorm.DB, token string) error {
	return db.Model(&models.Token{}).
		Where("token = ?", token).
		Update("last_used_at", time.Now()).Error
}

// CleanExpiredTokens 清理过期的 token
func CleanExpiredTokens(db *gorm.DB) error {
	return db.Where("expires_at < ?", time.Now()).
		Or("status = ?", "expired").
		Delete(&models.Token{}).Error
}
