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

	// At "error_handler.go", error.Error() is called, so except for panic, basically string will be passed.
	for _, v := range msg {
		if err, ok := v.(error); ok {
			message += "errorStack=" + fmt.Sprintf("%+v", err)
		}
	}

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
