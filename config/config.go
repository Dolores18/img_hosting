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
}

var AppConfigInstance AppConfig

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config/")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	if err := viper.Unmarshal(&AppConfigInstance); err != nil {
		log.Fatalf("Unable to decode into struct, %v", err)
	}
}
