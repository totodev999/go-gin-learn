package app

import (
	"flea-market/infra"

	"github.com/gin-gonic/gin"
)

type App struct {
	engine *gin.Engine
}

func NewApp() *App {
	infra.Initializer()
	db := infra.SetupDB()
	engine := NewRouter(db)
	return &App{engine: engine}
}

func (a *App) Run() error {
	return a.engine.Run()
}
