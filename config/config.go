package config

import (
	"log"

	"github.com/spf13/viper"
)

type AppConfig struct {
	App struct {
		Port int
	}
	Database struct {
		Host     string
		Port     int
		Username string
		Password string
	}

	Url struct {
		Imgurl string
	}
}

var AppConfigInstance AppConfig

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	// 添加多个可能的配置文件路径
	viper.AddConfigPath("./config/")     // 当前目录下的 config
	viper.AddConfigPath("../config/")    // 上级目录的 config
	viper.AddConfigPath("../../config/") // 上上级目录的 config

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	if err := viper.Unmarshal(&AppConfigInstance); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}
}
