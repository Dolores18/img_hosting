package models

import (
	"gorm.io/gorm"
	"time"
)

type UserInput struct {
	Name  string `json:"name" binding:"required,sign"`
	Age   int    `json:"age"`
	Email string `json:"email" binding:"required,email"`
	Psd   string `json:"psd" binding:"required,password"`
}
type UserLoginInput struct {
	Name  string `json:"name" `
	Email string `json:"email"`
	Psd   string `json:"psd" binding:"required,password"`
}

type UserInfo struct {
	UserID    uint           `gorm:"primaryKey" json:"user_id"`
	Name      string         `gorm:"type:varchar(255);not null;uniqueIndex" json:"name" binding:"required,sign"`
	Age       int            `gorm:"not null" json:"age" binding:"required,gte=18"`
	Email     string         `gorm:"type:varchar(255);not null;uniqueIndex" json:"email" binding:"required,email"`
	Psd       string         `gorm:"not null" json:"psd" binding:"required,password"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
