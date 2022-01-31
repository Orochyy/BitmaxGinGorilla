package repository

import (
	"BitmaxGinGorilla/entity"
	"gorm.io/gorm"
)

type SubscribeRepository interface {
	CreateSubscribe(user entity.Subsсribe) entity.Subsсribe
	KillSubscribe(s entity.Subsсribe)
	FindBookByID(subID uint64) entity.Subsсribe
}

type subscribeConnection struct {
	connection *gorm.DB
}

func (db *subscribeConnection) CreateSubscribe(s entity.Subsсribe) entity.Subsсribe {
	db.connection.Save(&s)
	db.connection.Preload("User").Find(&s)
	return s
}

func (db *subscribeConnection) KillSubscribe(s entity.Subsсribe) {
	db.connection.Delete(&s)
}

func (db *subscribeConnection) FindBookByID(subID uint64) entity.Subsсribe {
	var sub entity.Subsсribe
	db.connection.Preload("User").Find(&sub, subID)
	return sub
}
