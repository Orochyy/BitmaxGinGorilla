package controller

import (
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
	Delete(context *gin.Context)
}

type subscribeController struct {
	subscribeService service.SubscribeService
	jwtService       service.JWTService
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
