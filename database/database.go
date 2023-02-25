package database

import (
	"fmt"
	"sweng_backend/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {

	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}

	db.AutoMigrate(&model.Announcement{}, &model.Detail{}, &model.PA{}, &model.Product{}, &model.Solution{}, &model.Vertical{}, &model.Type{})
	db.AutoMigrate(&model.User{}, &model.Category{}, &model.Tag{}, &model.Project{})

	DB = db
	return db
}

func GetDB() *gorm.DB {
	return DB
}
