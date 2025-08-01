package controllers

import (
	"flea-market/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)

type IAPICallController interface {
	GetAllPosts(ctx *gin.Context)
}

type IAPICallService interface {
	GetAllPosts() (*[]repositories.Post, error)
}

type APICallController struct {
	service IAPICallService
}

func NewAPICallController(service IAPICallService) IAPICallController {
	return &APICallController{service: service}

}

func (c *APICallController) GetAllPosts(ctx *gin.Context) {
	data, err := c.service.GetAllPosts()
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})

}
