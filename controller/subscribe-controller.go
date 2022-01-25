package controller

import "BitmaxGinGorilla/service"

type SubscribeController interface {
	//Create
	//Delete
}

type subscribeController struct {
	subscribeController service.SubscribeService
	jwtService          service.JWTService
}
