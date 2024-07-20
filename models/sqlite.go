package models

import (
	"log"
	"sync"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

// GetDB 返回一个单例的数据库连接对象
func GetDB() *gorm.DB {
	once.Do(func() {
		var err error
		db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
		if err != nil {
			log.Fatalf("failed to connect database: %v", err)
		}

		// 自动迁移
		db.AutoMigrate(&UserInfo{}, &Image{}, &Roles{}, &Permissions{}, &UserRole{}, &RolePermission{})

	})
	return db
}
