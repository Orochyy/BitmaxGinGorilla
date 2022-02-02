package controller

import (
	"BitmaxGinGorilla/dto"
	"BitmaxGinGorilla/entity"
	"BitmaxGinGorilla/helper"
	"BitmaxGinGorilla/service"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type SubscribeController interface {
	Insert(context *gin.Context)
	Delete(context *gin.Context)
}

type subscribeController struct {
	subscribeService service.SubscribeService
	jwtService       service.JWTService
}

func NewSubController(subServ service.SubscribeService, jwtServ service.JWTService) SubscribeController {
	return &subscribeController{
		subscribeService: subServ,
		jwtService:       jwtServ,
	}
}

func (c *subscribeController) Insert(context *gin.Context) {
	var subCreateDTO dto.SubCreateDTO
	errDTO := context.ShouldBind(&subCreateDTO)
	if errDTO != nil {
		res := helper.BuildErrorResponse("Failed to process request", errDTO.Error(), helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, res)
	} else {
		authHeader := context.GetHeader("Authorization")
		userID := c.getUserIDByToken(authHeader)
		convertedUserID, err := strconv.ParseUint(userID, 10, 64)
		if err == nil {
			subCreateDTO.UserID = convertedUserID
		}
		result := c.subscribeService.Create(subCreateDTO)
		response := helper.BuildResponse(true, "OK", result)
		context.JSON(http.StatusCreated, response)
	}
}

func (c *subscribeController) Delete(context *gin.Context) {
	var sub entity.Subsсribe
	id, err := strconv.ParseUint(context.Param("id"), 0, 0)
	if err != nil {
		response := helper.BuildErrorResponse("Failed tou get id", "No param id were found", helper.EmptyObj{})
		context.JSON(http.StatusBadRequest, response)
	}
	sub.ID = id
	authHeader := context.GetHeader("Authorization")
	token, errToken := c.jwtService.ValidateToken(authHeader)
	if errToken != nil {
		panic(errToken.Error())
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := fmt.Sprintf("%v", claims["user_id"])
	if c.subscribeService.IsAllowedToEdit(userID, sub.ID) {
		c.subscribeService.Delete(sub)
		res := helper.BuildResponse(true, "Deleted", helper.EmptyObj{})
		context.JSON(http.StatusOK, res)
	} else {
		response := helper.BuildErrorResponse("You dont have permission", "You are not the owner", helper.EmptyObj{})
		context.JSON(http.StatusForbidden, response)
	}
}

func (c *subscribeController) getUserIDByToken(token string) string {
	aToken, err := c.jwtService.ValidateToken(token)
	if err != nil {
		panic(err.Error())
	}
	claims := aToken.Claims.(jwt.MapClaims)
	id := fmt.Sprintf("%v", claims["user_id"])
	return id
}
