package dao

import (
	"errors"
	"fmt"
	"img_hosting/models"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"
)

func CreateUser(db *gorm.DB, name string, age int, email string, password string) {
	user := models.UserInfo{Name: name, Age: age, Email: email, Psd: password, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	err := db.Create(&user).Error
	if err != nil {
		log.Fatalf("failed to create user: %v", err)
	}
	log.Printf("User created with ID: %d", user.UserID)
}

func FindUserByFields(db *gorm.DB, fields map[string]interface{}) (*models.UserInfo, error) {
	var user models.UserInfo
	query := db.Model(&models.UserInfo{})

	// 创建一个切片来存储所有的条件
	var conditions []string
	var values []interface{}

	// 构建查询条件
	for field, value := range fields {
		conditions = append(conditions, field+" = ?")
		values = append(values, value)
	}

	// 如果有条件，则应用它们
	if len(conditions) > 0 {
		query = query.Where(strings.Join(conditions, " AND "), values...)
	}

	result := query.First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // 用户不存在
		}
		return nil, result.Error // 其他错误
	}

	return &user, nil // 返回找到的用户信息
}

func FindUserByFieldsOr(db *gorm.DB, fields map[string]interface{}, options ...bool) (*models.UserInfo, error) {
	var user models.UserInfo
	query := db.Model(&models.UserInfo{})

	// 默认使用"或"查询
	isOrQuery := true
	if len(options) > 0 {
		isOrQuery = options[0]
	}

	// 使用 scope 函数来构建查询
	query = query.Scopes(func(db *gorm.DB) *gorm.DB {
		var conditions []string
		var values []interface{}

		for field, value := range fields {
			conditions = append(conditions, field+" = ?")
			values = append(values, value)
		}

		if len(conditions) > 0 {
			if isOrQuery {
				return db.Where(strings.Join(conditions, " OR "), values...)
			} else {
				return db.Where(strings.Join(conditions, " AND "), values...)
			}
		}
		return db
	})

	result := query.First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // 用户不存在
		}
		return nil, fmt.Errorf("查询数据库时发生错误: %w", result.Error) // 包装错误以提供更多上下文
	}

	return &user, nil
}

func FindUserByEmail(db *gorm.DB, email string) (*models.UserInfo, error) {
	var user models.UserInfo
	result := db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			log.Printf("User with name %s not found", email)
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
