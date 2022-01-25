package migrations

import (
	"BitmaxGinGorilla/config"
	"BitmaxGinGorilla/entity"
	"gorm.io/gorm"
)

var db *gorm.DB = config.SetupDatabaseConnection()

func DbMigrate() {
	db.AutoMigrate(&entity.User{})
}
