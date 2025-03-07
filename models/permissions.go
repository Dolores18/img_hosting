package models

import "time"

type Permissions struct {
	PermissionID uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"column:permission_name" json:"permission_name"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

// TableName 指定表名
func (Permissions) TableName() string {
	return "permissions"
}
