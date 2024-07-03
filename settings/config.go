package settings

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	// Database configurations
	DBUsername string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string

	// Other configurations here
}

var AppConfig *Config

func InitConfig() {
	viper.SetConfigFile("config.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// Default configurations here
	viper.SetDefault("DBUsername", "apiuser")
	viper.SetDefault("DBPassword", "password")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading configuration: %v", err)
	}
	if err := viper.Unmarshal(&AppConfig); err != nil {
		log.Fatalf("Error when unmarshaling the configuration: %v", err)
	}
}
