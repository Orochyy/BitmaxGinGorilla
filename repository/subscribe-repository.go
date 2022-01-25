package repository

import (
	"BitmaxGinGorilla/entity"
	"gorm.io/gorm"
)

type SubscribeRepository interface {
	CreateSubscribe(user entity.Subsribe) entity.Subsribe
	KillSubscribe(s entity.Subsribe)
}

type subscribeConnection struct {
	connection *gorm.DB
}

func (db *subscribeConnection) CreateSubscribe(s entity.Subsribe) entity.Subsribe {
	db.connection.Save(&s)
	db.connection.Preload("User").Find(&s)
	return s
}

func (db *subscribeConnection) KillSubscribe(s entity.Subsribe) {
	db.connection.Delete(&s)
}
