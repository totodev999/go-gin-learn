package utils

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger(messageId MessageCode, ctx *gin.Context, msg ...any) {
	logTemplate := "[%s]-[Back]-[%s]-[%s]-[%s]-[%s]-[%s] %s\n"

	firstLetter := string([]rune(messageId)[0])
	var logLevel string

	switch firstLetter {
	case "I":
		logLevel = "INFO"
	case "W":
		logLevel = "WARN"
	case "E":
		logLevel = "ERROR"
	default:
		logLevel = "INFO"
	}

	message := fmt.Sprintf(Messages[messageId], msg...)

	reqId, exist := ctx.Get("reqId")
	if !exist {
		reqId = "not set"
	}
	// 文字列化しておく
	reqIdStr := fmt.Sprintf("%v", reqId)

	fmt.Printf(
		logTemplate,
		time.Now().Format(time.DateTime),
		logLevel,
		messageId,
		ctx.Request.Method+" "+ctx.Request.URL.Path,
		reqIdStr,
		ctx.ClientIP(),
		message,
	)
}
