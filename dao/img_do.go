package dao

import (
	"img_hosting/models"
	"log"
	"time"

	"gorm.io/gorm"
)

// 创建图片
func CreateImg(db2 *gorm.DB, userid uint, image_url string, image_name string, image_extension string, hash_image string, image_size int64, image_type string) uint {
	img := models.Image{UserID: userid, ImageURL: image_url, ImageName: image_name, Imageextenion: image_extension, HashImage: hash_image, ImageSize: image_size, ImageType: image_type, UploadTime: time.Now()}
	err := db2.Create(&img).Error
	if err != nil {
		log.Fatalf("failed to create user: %v", err)
	}
	log.Printf("Image created with ID: %d", img.ImageID)
	return img.ImageID
}
