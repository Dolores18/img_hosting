package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"img_hosting/dao"
	"img_hosting/models"
	"time"
)

// CreateToken 创建新的 token
func CreateToken(userID uint, deviceID, ipAddress string) (*models.Token, error) {
	// 生成随机 token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, err
	}
	tokenString := hex.EncodeToString(tokenBytes)

	// 如果没有提供 deviceID，生成一个随机标识符
	if deviceID == "" {
		deviceBytes := make([]byte, 8)
		rand.Read(deviceBytes)
		deviceID = "dev_" + hex.EncodeToString(deviceBytes)
	}

	// 设置过期时间（30天）
	expiresAt := time.Now().Add(30 * 24 * time.Hour)

	// 创建 token
	err := dao.CreateToken(models.GetDB(), tokenString, userID, expiresAt, deviceID, ipAddress)
	if err != nil {
		return nil, err
	}

	return dao.GetTokenByString(models.GetDB(), tokenString)
}

// ListUserTokens 获取用户的所有 token
func ListUserTokens(userID uint) ([]models.Token, error) {
	return dao.ListUserTokens(models.GetDB(), userID)
}

// RevokeToken 撤销指定的 token
func RevokeToken(userID uint, token string) error {
	// 验证 token 属于该用户
	tokenModel, err := dao.GetTokenByString(models.GetDB(), token)
	if err != nil {
		return err
	}
	if tokenModel.UserID != userID {
		return errors.New("unauthorized")
	}

	return dao.RevokeToken(models.GetDB(), token)
}

type TokenService struct{}

func NewTokenService() *TokenService {
	return &TokenService{}
}

// ValidateToken 验证token并返回用户信息
func (s *TokenService) ValidateToken(tokenStr string) (*models.UserInfo, error) {
	db := models.GetDB()

	// 使用现有的 dao 方法获取 token
	token, err := dao.GetTokenByString(db, tokenStr)
	if err != nil {
		return nil, errors.New("token无效")
	}

	// 检查过期时间
	if token.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("token已过期")
	}

	// 更新最后使用时间
	dao.UpdateTokenLastUsed(db, tokenStr)

	// 获取用户信息
	var user models.UserInfo
	if err := db.First(&user, token.UserID).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	return &user, nil
}

// CheckFileAccess 检查用户是否有权限访问文件
func (s *TokenService) CheckFileAccess(userID uint, path string) (bool, error) {
	db := models.GetDB()

	var file models.PrivateFile
	if err := db.Where("user_id = ? AND storage_path = ? AND status = ?",
		userID, path, models.FileStatusActive).First(&file).Error; err != nil {
		return false, nil
	}

	return true, nil
}

// UpdateTokenUsage 更新token使用记录
func (s *TokenService) UpdateTokenUsage(tokenStr string) error {
	return dao.UpdateTokenLastUsed(models.GetDB(), tokenStr)
}
