package service

import (
	"BitmaxGinGorilla/dto"
	"BitmaxGinGorilla/entity"
	"BitmaxGinGorilla/repository"
	"github.com/mashingan/smapping"
	"log"
)

type SubscribeService interface {
	Create(s dto.Subsription) entity.Subsribe
	Delete(s dto.Subsription)
}

type subscribeService struct {
	subscribeRepository repository.SubscribeRepository
}

func NewSubService(subscribeRepo repository.SubscribeRepository) subscribeService {
	return &subscribeService{
		subscribeRepository: subscribeRepo,
	}
}

func (service *subscribeService) Create(s dto.Subsription) entity.Subsribe {
	subscribe := entity.Subsribe{}
	err := smapping.FillStruct(&subscribe, smapping.MapFields(&s))
	if err != nil {
		log.Fatalf("Failed map %v: ", err)
	}
	res := service.subscribeRepository.CreateSubscribe(subscribe)
	return res
}

func (service *subscribeService) Delete(s entity.Subsribe) {
	service.subscribeRepository.KillSubscribe(s)
}
