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

	// default values
	methodPath := "N/A"
	reqIdStr := "not set"
	clientIP := "N/A"

	if ctx != nil {
		methodPath = ctx.Request.Method + " " + ctx.Request.URL.Path
		if reqId, exist := ctx.Get("reqId"); exist {
			reqIdStr = fmt.Sprintf("%v", reqId)
		}
		clientIP = ctx.ClientIP()
	}

	fmt.Printf(
		logTemplate,
		time.Now().Format(time.DateTime),
		logLevel,
		messageId,
		methodPath,
		reqIdStr,
		clientIP,
		message,
	)
}
