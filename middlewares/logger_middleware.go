package middlewares

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

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

		logTemplate := "[%s]-[%s]-[ID:%s]-[%s] %s\n"

		fmt.Printf(logTemplate, time.Now().Format(time.DateTime), "START", reqId, ctx.Request.Method+" "+ctx.Request.URL.Path, "開始")

		ctx.Next()
		status := ctx.Writer.Status()
		fmt.Printf(logTemplate, time.Now().Format(time.DateTime), "END", reqId, ctx.Request.Method+" "+ctx.Request.URL.Path, "終了"+strconv.Itoa(status))
	}
}
