package dao

import (
	"gorm.io/gorm"
	"img_hosting/models"
	"log"
	"time"
)

// 创建用户
func CreateUser(db *gorm.DB, name string, age int, password string) {
	user := models.UserInfo{Name: name, Age: age, Psd: password, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	err := db.Create(&user).Error
	if err != nil {
		log.Fatalf("failed to create user: %v", err)
	}
	log.Printf("User created with ID: %d", user.ID)
}

// FindUserByName 查找给定 name 的用户数据
func FindUserByName(db *gorm.DB, name string) (*models.UserInfo, error) {
	var user models.UserInfo
	result := db.Where("name = ?", name).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Printf("User with name %s not found", name)
			return nil, nil
		}
		log.Printf("Failed to find user: %v", result.Error)
		return nil, result.Error
	}
	return &user, nil
}

// 更新用户数据
func UpdateUser(db *gorm.DB, name string, updates map[string]interface{}) error {
	//执行更新
	if err := db.Model(&models.UserInfo{}).Where("name = ?", name).Updates(updates).Error; err != nil {
		return err
	}
	return nil
}

/*
updates = map[string]interface{}{
"name": ”jane"
"email“ ：”123“

更简单的方法
    // 解析 JSON 字符串
    err = json.Unmarshal([]byte(jsonString), &updates)
    if err != nil {
        fmt.Println("Error parsing JSON:", err)
        return
    }

    // 打印解析后的数据
    fmt.Println("Parsed JSON:", updates)
*/
