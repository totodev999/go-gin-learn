package controllers

import (
	"context"
	"flea-market/repositories"
	"flea-market/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IAPICallService interface {
	GetAllPosts(ctx context.Context) (*[]repositories.Post, error)
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
