package migrations

import (
	"BitmaxGinGorilla/config"
	"GinApiGormMysqlElif/entity"
	"gorm.io/gorm"
)

var db *gorm.DB = config.SetupDatabaseConnection()

func DbMigrate() {
	db.AutoMigrate(&entity.User{})
}
