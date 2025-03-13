package models

import (
	"time"

	"gorm.io/gorm"
)

type UserInput struct {
	Name  string `json:"name" binding:"required,sign"`
	Age   int    `json:"age"`
	Email string `json:"email" binding:"required,email"`
	Psd   string `json:"psd" binding:"required,password"`
}

type UserLoginInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Psd   string `json:"psd" binding:"required,password"`
}

type UserInfo struct {
	UserID      uint           `gorm:"primaryKey;column:user_id" json:"user_id"`
	Name        string         `gorm:"column:name" json:"name"`
	Email       string         `gorm:"column:email" json:"email"`
	Password    string         `gorm:"column:psd" json:"-"`
	Phone       string         `gorm:"column:phone" json:"phone"`
	Age         int            `gorm:"column:age" json:"age"`
	Status      string         `gorm:"column:status" json:"status"`
	LastLoginAt time.Time      `gorm:"column:last_login_at" json:"last_login_at"`
	LastLoginIP string         `gorm:"column:last_login_ip" json:"last_login_ip"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at" json:"-"`
	Roles       []Roles        `gorm:"many2many:user_roles;foreignKey:UserID;joinForeignKey:UserID;References:RoleID;joinReferences:RoleID;constraint:OnDelete:CASCADE" json:"roles"`
}

// UserStatus 用户状态常量
const (
	UserStatusActive   = "active"
	UserStatusInactive = "inactive"
	UserStatusBanned   = "banned"
)
