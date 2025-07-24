package controllers

import (
	"free-market/dto"
	"free-market/services"
	"free-market/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IAuthController interface {
	Signup(ctx *gin.Context)
	Login(ctx *gin.Context)
}

type AuthController struct {
	service services.IAuthService
}

func (c *AuthController) Signup(ctx *gin.Context) {
	var input dto.SignupInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		_ = ctx.Error(utils.NewBadRequestError("Input data is invalid", err))
		return
	}

	err := c.service.Signup(input.Email, input.Password)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.Status(http.StatusCreated)
}

func (c *AuthController) Login(ctx *gin.Context) {
	var input dto.LoginInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		_ = ctx.Error(utils.NewBadRequestError("Input data is invalid", err))
		return
	}

	token, err := c.service.Login(input.Email, input.Password)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func NewAuthController(service services.IAuthService) IAuthController {
	return &AuthController{service: service}
}
