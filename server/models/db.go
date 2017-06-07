package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

func InitDB(args string) error {
	var err error
	db, err = gorm.Open("postgres", args)
	//db.AutoMigrate(&User{})
	//db.AutoMigrate(&Room{})
	//db.AutoMigrate(&Gamer{})
	//db.AutoMigrate(&Game{})
	//db.AutoMigrate(&Stats{})
	return err
}
