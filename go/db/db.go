package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDatabase() {
	var err error
	DB, err = gorm.Open(sqlite.Open("db.sqlite"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
}
