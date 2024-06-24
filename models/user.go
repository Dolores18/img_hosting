package models

import (
	"gorm.io/gorm"
	"time"
)

type UserInput struct {
	Name string `json:"name" binding:"required,sign"`
	Age  int    `json:"age" binding:"required,gte=18"`
	Psd  string `json:"psd" binding:"required,password"`
}
type UserInfo struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name" binding:"required,sign"`
	Age       int            `gorm:"not null" json:"age" binding:"required,gte=18"`
	Psd       string         `gorm:"not null" json:"psd" binding:"required,password"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
