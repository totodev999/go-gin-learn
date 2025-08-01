package utils

type MessageCode string

const (
	RequestStart            MessageCode = "I001-00001"
	RequestEnd              MessageCode = "I001-00002"
	ExternalAPIRequestStart MessageCode = "I001-00003"
	ExternalAPIRequestEnd   MessageCode = "I001-00004"

	BadRequest     MessageCode = "I001-00010"
	NotFound       MessageCode = "I001-00011"
	UnAuthorized   MessageCode = "I001-00012"
	GenericMessage MessageCode = "I001-00020"

	DuplicateKeyError       MessageCode = "W001-00001"
	ExternalAPIReturnsError MessageCode = "W001-00010"

	DBError                    MessageCode = "E001-00001"
	ExternalAPIConnectionError MessageCode = "E001-00002"

	UnknownError MessageCode = "E001-00010"
)

var Messages = map[MessageCode]string{
	RequestStart:            "リクエスト開始",
	RequestEnd:              "リクエスト終了 status code:%v",
	ExternalAPIRequestStart: "外部APIリクエスト開始  Request:%v",
	ExternalAPIRequestEnd:   "外部APIリクエスト終了 %v",

	BadRequest:   "Bad request",
	NotFound:     "Not Found",
	UnAuthorized: "UnAuthorized",

	GenericMessage: "%v",

	DuplicateKeyError:       "Duplicate key",
	ExternalAPIReturnsError: "External API returns an error:%v",

	DBError:                    "DB error",
	ExternalAPIConnectionError: "Connection failed Error:%v",
	UnknownError:               "UnknownError Error Detail:%v",
}
