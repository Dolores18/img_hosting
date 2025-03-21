package models

import "time"

type Roles struct {
	RoleID      uint          `gorm:"primaryKey" json:"id"`
	RoleName    string        `gorm:"unique" json:"role_name"`
	Description string        `json:"description,omitempty"`
	IsActive    bool          `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time     `gorm:"autoCreateTime"`
	UpdatedAt   time.Time     `gorm:"autoUpdateTime"`
	Permissions []Permissions `gorm:"many2many:role_permissions;foreignKey:RoleID;joinForeignKey:RoleID;References:PermissionID;joinReferences:PermissionID;constraint:OnDelete:CASCADE"`
}
