package database

import (
	"log"
	"product-api/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"product-api/settings"
)

var DB *gorm.DB

func Connect(cfg *settings.Config) {
	dsn := cfg.DBUsername +
		":" + cfg.DBPassword +
		"@tcp(" + cfg.DBHost +
		":" + cfg.DBPort + ")/" +
		cfg.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	err = DB.AutoMigrate(&models.Product{})
	if err != nil {
		log.Fatal("Failed to perform database migrations:", err)
	}
}
