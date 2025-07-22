package controllers

import (
	"errors"
	"free-market/dto"
	"free-market/models"
	"free-market/services"
	"free-market/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type IItemController interface {
	FindAll(ctx *gin.Context)
	FindById(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type ItemController struct {
	service services.IItemService
}

func NewItemController(service services.IItemService) IItemController {
	return &ItemController{service: service}
}

func (c *ItemController) FindAll(ctx *gin.Context) {
	items, err := c.service.FindAll()
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": items})
}

func (c *ItemController) FindById(ctx *gin.Context) {
	userId, err := getUserId(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	itemId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Error(utils.NewBadRequestError("can't get id from path", err))
		return
	}

	item, err := c.service.FindById(uint(itemId), *userId)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": item})
}

func (c *ItemController) Create(ctx *gin.Context) {
	userId, err := getUserId(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	var input dto.CreateItemInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.Error(utils.NewBadRequestError("Input data is invalid", err))
		return
	}

	newItem, err := c.service.Create(input, *userId)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": newItem})

}

func (c *ItemController) Update(ctx *gin.Context) {
	userId, err := getUserId(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Error(utils.NewBadRequestError("can't get id from path", err))
		return
	}

	var input dto.UpdateItemInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.Error(utils.NewBadRequestError("Input data is invalid", err))
		return
	}

	updatedItem, err := c.service.Update(uint(id), input, *userId)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedItem})

}

func (c *ItemController) Delete(ctx *gin.Context) {
	userId, err := getUserId(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Error(utils.NewBadRequestError("can't get id from path", err))
		return
	}

	err = c.service.Delete(uint(id), *userId)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusOK)
}

func getUserId(ctx *gin.Context) (*uint, error) {
	user, exists := ctx.Get("user")
	if !exists {
		return nil, utils.NewUnauthorized("user is not set in request", errors.New("UnAuthorized"))
	}

	usr, ok := user.(*models.User)
	if !ok {
		return nil, utils.NewUnauthorized("user in context is invalid", errors.New("InvalidType"))
	}
	userId := usr.ID

	return &userId, nil
}
