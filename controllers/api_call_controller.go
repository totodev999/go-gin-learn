package controllers

import (
	"context"
	"flea-market/repositories"
	"flea-market/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type IAPICallService interface {
	GetAllPosts(ctx context.Context) (*[]repositories.Post, error)
	GetUserAndPosts(ctx context.Context, userId uint) (*repositories.UserAndPosts, error)
}

type APICallController struct {
	service IAPICallService
}

func NewAPICallController(service IAPICallService) *APICallController {
	return &APICallController{service: service}

}

func (c *APICallController) GetAllPosts(ctx *gin.Context) {
	reqCtx := utils.GinToGoContext(ctx)
	data, err := c.service.GetAllPosts(reqCtx)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})

}

func (c *APICallController) GetUserAndPosts(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Param("userId"))
	if err != nil {
		_ = ctx.Error(utils.NewBadRequestError("can't get userId from path", err))
		return
	}
	reqCtx := utils.GinToGoContext(ctx)

	data, err := c.service.GetUserAndPosts(reqCtx, uint(userId))

	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})

}
