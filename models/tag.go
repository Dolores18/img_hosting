package models

type Tag struct {
	TagID   uint   `gorm:"primaryKey;column:tag_id"`
	UserID  uint   `gorm:"column:user_id;not null"`
	TagName string `gorm:"column:tag_name;not null;unique"` // 添加unique约束
}
