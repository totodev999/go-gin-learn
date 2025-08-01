package middlewares

import (
	// カスタムエラー型のパッケージ
	"flea-market/utils"

	"github.com/gin-gonic/gin"
)

func APIErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) > 0 {
			err := ctx.Errors.Last().Err

			if apiErr, ok := err.(*utils.APIError); ok {
				utils.Logger(apiErr.MessageCode, "", "", "", err.Error())
				ctx.JSON(apiErr.StatusCode, gin.H{"error": apiErr.MessageCode})
				return
			} else {
				utils.Logger(utils.UnknownError, "", "", "", err.Error())
				ctx.JSON(500, gin.H{"error": "Internal server error"})
				return
			}
		}
	}
}
