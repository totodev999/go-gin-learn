package middlewares

import (
	"crypto/rand"
	"encoding/hex"
	"free-market/utils"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		b := make([]byte, 8)
		if _, err := rand.Read(b); err != nil {
			panic("failed to generate random ID: " + err.Error())
		}
		reqId := hex.EncodeToString(b)
		ctx.Set("reqId", reqId)

		utils.Logger(utils.RequestStart, ctx)

		ctx.Next()
		status := ctx.Writer.Status()
		utils.Logger(utils.RequestEnd, ctx, status)
	}
}
