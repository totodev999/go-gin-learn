package middlewares

import (
	"crypto/rand"
	"encoding/hex"
	"flea-market/utils"

	"github.com/gin-gonic/gin"
)

func LoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		b := make([]byte, 8)
		if _, err := rand.Read(b); err != nil {
			panic("failed to generate random ID: " + err.Error())
		}
		methodPath := ctx.Request.Method + " " + ctx.Request.URL.Path
		reqId := hex.EncodeToString(b)
		clientIP := ctx.ClientIP()

		// Set Context
		utils.SetGinContext(ctx, utils.ContextReqID, reqId)
		utils.SetGinContext(ctx, utils.ContextIP, clientIP)
		utils.SetGinContext(ctx, utils.ContextMethodPath, methodPath)

		utils.Logger(utils.RequestStart, methodPath, reqId, clientIP)

		ctx.Next()
		status := ctx.Writer.Status()
		utils.Logger(utils.RequestEnd, methodPath, reqId, clientIP, status)
	}
}
