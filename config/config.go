package config

import (
	"log"

	"github.com/spf13/viper"
)

type LogConfig struct {
	Path       string `mapstructure:"path"`
	Filename   string `mapstructure:"filename"`
	Level      string `mapstructure:"level"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
	Compress   bool   `mapstructure:"compress"`
}

type AppConfig struct {
	App struct {
		Port int
	}

	Upload struct {
		Path           string `mapstructure:"path"`
		ThumbnailsPath string `mapstructure:"thumbnails_path"`
		MaxSize        int64  `mapstructure:"max_size"`
	}

	PrivateFiles struct {
		Path         string `mapstructure:"path"`
		MaxSize      int64  `mapstructure:"max_size"`
		AllowedTypes string `mapstructure:"allowed_types"`
	} `mapstructure:"private_files"`

	Database struct {
		Host     string
		Port     int
		Username string
		Password string
	}

	Url struct {
		Imgurl string
	}

	Permissions struct {
		Routes map[string][]string `mapstructure:"routes"`
		Roles  map[string][]string `mapstructure:"roles"`
	} `mapstructure:"permissions"`

	Log LogConfig `mapstructure:"log"`
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

// GetConfig 获取配置实例
func GetConfig() *AppConfig {
	return &AppConfigInstance
}
