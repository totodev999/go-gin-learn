package app

import (
	"flea-market/controllers"
	"flea-market/infra"
	"flea-market/middlewares"
	"flea-market/repositories"
	"flea-market/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB) *gin.Engine {
	itemRepository := repositories.NewItemRepository(db)
	itemService := services.NewItemService(itemRepository)
	itemController := controllers.NewItemController(itemService)

	authRepository := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepository)
	authController := controllers.NewAuthController(authService)

	apiClient := infra.NewBaseAPIClient()
	apiCallRepository := repositories.NewAPICallRepository(apiClient)
	apiCallService := services.NewAPICallService(apiCallRepository)
	apiCallController := controllers.NewAPICallController(apiCallService)

	router := gin.New()
	router.Use(middlewares.LoggerMiddleware())
	router.Use(middlewares.APIErrorHandler())
	router.Use(gin.Recovery())

	itemRouter := router.Group("/items")
	itemRouterWithAuth := router.Group("/items", middlewares.AuthMiddleware(authService))
	authRouter := router.Group("/auth")
	externalRouter := router.Group("/external")

	itemRouter.GET("", itemController.FindAll)
	itemRouterWithAuth.GET("/:id", itemController.FindById)
	itemRouterWithAuth.POST("", itemController.Create)
	itemRouterWithAuth.PUT("/:id", itemController.Update)
	itemRouterWithAuth.DELETE("/:id", itemController.Delete)

	authRouter.POST("/signup", authController.Signup)
	authRouter.POST("/login", authController.Login)

	externalRouter.GET("", apiCallController.GetAllPosts)

	return router
}
