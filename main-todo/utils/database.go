package utils

import (
	"main-todo/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {

	DB, err := gorm.Open(sqlite.Open("sample.DB"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")

	}

	// Auto-migrate the Todo model to create the table
	DB.AutoMigrate(&models.Todo{}, &models.User{})
}

func GetDB() *gorm.DB {
	return DB
}
