package models

import "time"

type Permissions struct {
	PermissionID   uint      `gorm:"primaryKey" json:"id"`
	PermissionName string    `json:"permission_name"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime"`
}
