package middlewares

import (
	// カスタムエラー型のパッケージ
	"flea-market/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

func APIErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) > 0 {
			err := ctx.Errors.Last().Err

			if apiErr, ok := err.(*utils.APIError); ok {
				utils.Logger(apiErr.MessageCode, ctx, err.Error())
				fmt.Printf("utils.APIError  %v", apiErr.Message)
				ctx.JSON(apiErr.StatusCode, gin.H{"error": apiErr.MessageCode})
				return
			} else {
				utils.Logger(utils.UnknownError, ctx, err.Error())
				ctx.JSON(500, gin.H{"error": "Internal server error"})
				return
			}
		}
	}
}
