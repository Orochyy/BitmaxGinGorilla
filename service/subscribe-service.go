package service

import (
	"BitmaxGinGorilla/dto"
	"BitmaxGinGorilla/entity"
	"BitmaxGinGorilla/repository"
	"fmt"
	"github.com/mashingan/smapping"
	"log"
)

type SubscribeService interface {
	Create(s dto.Subscribe) entity.Subsсribe
	Delete(s entity.Subsсribe)
	IsAllowedToEdit(userID string, subID uint64) bool
}

type subscribeService struct {
	subscribeRepository repository.SubscribeRepository
}

func NewSubService(subscribeRepo repository.SubscribeRepository) SubscribeService {
	return &subscribeService{
		subscribeRepository: subscribeRepo,
	}
}

func (service *subscribeService) Create(s dto.Subscribe) entity.Subsсribe {
	subscribe := entity.Subsсribe{}
	err := smapping.FillStruct(&subscribe, smapping.MapFields(&s))
	if err != nil {
		log.Fatalf("Failed map %v: ", err)
	}
	res := service.subscribeRepository.CreateSubscribe(subscribe)
	return res
}

func (service *subscribeService) Delete(s entity.Subsсribe) {
	service.subscribeRepository.KillSubscribe(s)
}
func (service *subscribeService) IsAllowedToEdit(userID string, subID uint64) bool {
	s := service.subscribeRepository.FindBookByID(subID)
	id := fmt.Sprintf("%v", s.UserID)
	return userID == id
}
