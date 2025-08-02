package utils

import (
	"fmt"
	"time"
)

func orDefault(val, def string) string {
	if val == "" {
		return def
	}
	return val
}

func Logger(messageId MessageCode, methodPath, reqID, clientIP string, msg ...any) {
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

	fmt.Printf(
		logTemplate,
		time.Now().Format(time.DateTime),
		logLevel,
		messageId,
		orDefault(methodPath, "N/A"),
		orDefault(reqID, "N/A"),
		orDefault(clientIP, "N/A"),
		message,
	)
}
