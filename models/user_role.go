package models

type UserRole struct {
	UserID uint     `gorm:"primaryKey"`
	RoleID uint     `gorm:"primaryKey"`
	User   UserInfo `gorm:"foreignKey:UserID"`
	Role   Roles    `gorm:"foreignKey:RoleID"`
}
