package models

import (
	"log"
	"sync"
	"time"

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
		// 添加 SQLite 优化配置
		dsn := "test.db?_busy_timeout=5000&_journal_mode=WAL&_synchronous=NORMAL"
		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{
			PrepareStmt: true,
		})
		if err != nil {
			log.Fatalf("failed to connect database: %v", err)
		}

		// 设置连接池
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("failed to get db instance: %v", err)
		}

		// 设置连接池参数
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(time.Hour)

		// 自动迁移所有模型
		err = db.AutoMigrate(
			&UserInfo{},
			&Image{},
			&Roles{},
			&Permissions{},
			&UserRole{},
			&RolePermission{},
			&Tag{},
			&ImageTag{},
			&Token{},
			&File{},

			&PrivateFile{},
		)
		if err != nil {
			log.Printf("数据库迁移失败: %v", err)
		}
	})
	return db
}
