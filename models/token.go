package models

import (
	"time"
)

type Token struct {
	Token      string    `gorm:"primaryKey;type:varchar(255)" json:"token"`
	UserID     uint      `gorm:"not null;index" json:"user_id"`
	ExpiresAt  time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	DeviceID   string    `gorm:"type:varchar(255)" json:"device_id"`
	IPAddress  string    `gorm:"type:varchar(255)" json:"ip_address"`
	Status     string    `gorm:"type:varchar(20);default:'active'" json:"status"` // active, inactive, expired, revoked
	LastUsedAt time.Time `json:"last_used_at"`
	User       UserInfo  `gorm:"foreignKey:UserID;references:UserID" json:"-"`
}

// TokenStatus 定义可能的令牌状态
const (
	TokenStatusActive   = "active"
	TokenStatusInactive = "inactive"
	TokenStatusExpired  = "expired"
	TokenStatusRevoked  = "revoked"
)
