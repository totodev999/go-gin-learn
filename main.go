package main

import (
	"free-market/controllers"
	"free-market/infra"
	"free-market/middlewares"
	"free-market/repositories"
	"free-market/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setUpRouter(db *gorm.DB) *gin.Engine {
	itemRepository := repositories.NewItemRepository(db)
	itemService := services.NewItemService(itemRepository)
	itemController := controllers.NewItemController(itemService)

	authRepository := repositories.NewAuthRepository(db)
	authService := services.NewAuthService(authRepository)
	authController := controllers.NewAuthController(authService)

	router := gin.New()
	router.Use(middlewares.LoggerMiddleware())
	router.Use(gin.Recovery())

	itemRouter := router.Group("/items")

	itemRouterWithAuth := router.Group("/items", middlewares.AuthMiddleware(authService))

	authRouter := router.Group("/auth")

	itemRouter.GET("", itemController.FindAll)
	itemRouterWithAuth.GET("/:id", itemController.FIndById)
	itemRouterWithAuth.POST("", itemController.Create)
	itemRouterWithAuth.PUT("/:id", itemController.Update)
	itemRouterWithAuth.DELETE("/:id", itemController.Delete)

	authRouter.POST("/signup", authController.Signup)
	authRouter.POST("/login", authController.Login)

	return router
}

func main() {
	infra.Initializer()
	db := infra.SetupDB()

	router := setUpRouter(db)

	router.Run() // 0.0.0.0:8080 でサーバーを立てます。
}
