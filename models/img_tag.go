package models

type ImageTag struct {
	ImageID uint  `gorm:"primaryKey"`
	TagID   uint  `gorm:"primaryKey"`
	Image   Image `gorm:"foreignKey:ImageID"`
	Tag     Tag   `gorm:"foreignKey:TagID"`
}
