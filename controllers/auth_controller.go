package controllers

import (
	"context"
	"flea-market/dto"
	"flea-market/models"
	"flea-market/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IAuthService interface {
	Signup(ctx context.Context, email string, password string) error
	Login(ctx context.Context, email string, password string) (*string, error)
	GetUserFromToken(toke string) (*models.User, error)
}

type AuthController struct {
	service IAuthService
}

func (c *AuthController) Signup(ctx *gin.Context) {
	reqCtx := utils.GinToGoContext(ctx)
	var input dto.SignupInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		_ = ctx.Error(utils.NewBadRequestError("Input data is invalid", err))
		return
	}

	err := c.service.Signup(reqCtx, input.Email, input.Password)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.Status(http.StatusCreated)
}

func (c *AuthController) Login(ctx *gin.Context) {
	reqCtx := utils.GinToGoContext(ctx)

	var input dto.LoginInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		_ = ctx.Error(utils.NewBadRequestError("Input data is invalid", err))
		return
	}

	token, err := c.service.Login(reqCtx, input.Email, input.Password)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func NewAuthController(service IAuthService) *AuthController {
	return &AuthController{service: service}
}
