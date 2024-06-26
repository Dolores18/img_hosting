package dao

import (
	"gorm.io/gorm"
	"img_hosting/models"
	"log"
	"time"
)

type Image struct {
	ImageID     uint      `gorm:"primaryKey" json:"id"`
	UserName    string    `gorm:"not null" json:"username"` // 用户名
	ImageURL    string    `json:"image_url"`                // 图片存储路径或URL
	ImageName   string    `json:"image_name"`               // 图片名称
	ImageSize   int64     `json:"image_size"`               // 图片大小（字节）
	ImageType   string    `json:"image_type"`               // 图片格式
	UploadTime  time.Time `gorm:"autoCreateTime"`           // 上传时间
	Description string    `json:"description"`              // 图片描述（可选）
}

// 创建图片
func CreateImg(db2 *gorm.DB, username string, image_url string, image_name string, image_size int64, image_type string) {
	img := models.Image{UserName: username, ImageURL: image_url, ImageName: image_name, ImageSize: image_size, ImageType: image_type, UploadTime: time.Now()}
	err := db2.Create(&img).Error
	if err != nil {
		log.Fatalf("failed to create user: %v", err)
	}
	log.Printf("Image created with ID: %d", img.ImageID)
}
