package db

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var GormDatabase *gorm.DB

func init() {
	var err error
	GormDatabase, err = gorm.Open(sqlite.Open("./db/db.sqlite"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println("Database connection successful!")
}
