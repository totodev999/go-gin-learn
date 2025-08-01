package controllers

import (
	"errors"
	"flea-market/dto"
	"flea-market/models"
	"flea-market/utils"
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

type IItemService interface {
	FindAll() (*[]models.Item, error)
	FindById(itemId uint, userId uint) (*models.Item, error)
	Create(createItemInput dto.CreateItemInput, userId uint) (*models.Item, error)
	Update(itemId uint, updateItemInput dto.UpdateItemInput, userId uint) (*models.Item, error)
	Delete(itemId uint, userId uint) error
}

type ItemController struct {
	service IItemService
}

func NewItemController(service IItemService) IItemController {
	return &ItemController{service: service}
}

func (c *ItemController) FindAll(ctx *gin.Context) {
	items, err := c.service.FindAll()
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": items})
}

func (c *ItemController) FindById(ctx *gin.Context) {
	userId, err := getUserId(ctx)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	itemId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		_ = ctx.Error(utils.NewBadRequestError("can't get id from path", err))
		return
	}

	item, err := c.service.FindById(uint(itemId), *userId)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": item})
}

func (c *ItemController) Create(ctx *gin.Context) {
	userId, err := getUserId(ctx)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	var input dto.CreateItemInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		_ = ctx.Error(utils.NewBadRequestError("Input data is invalid", err))
		return
	}

	newItem, err := c.service.Create(input, *userId)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"data": newItem})

}

func (c *ItemController) Update(ctx *gin.Context) {
	userId, err := getUserId(ctx)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		_ = ctx.Error(utils.NewBadRequestError("can't get id from path", err))
		return
	}

	var input dto.UpdateItemInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		_ = ctx.Error(utils.NewBadRequestError("Input data is invalid", err))
		return
	}

	updatedItem, err := c.service.Update(uint(id), input, *userId)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": updatedItem})

}

func (c *ItemController) Delete(ctx *gin.Context) {
	userId, err := getUserId(ctx)
	if err != nil {
		_ = ctx.Error(err)
		return
	}

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		_ = ctx.Error(utils.NewBadRequestError("can't get id from path", err))
		return
	}

	err = c.service.Delete(uint(id), *userId)
	if err != nil {
		_ = ctx.Error(err)
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
