package utils

type MessageCode string

const (
	RequestStart MessageCode = "I001-00001"
	RequestEnd   MessageCode = "I001-00002"
	BadRequest   MessageCode = "I001-00010"
	NotFound     MessageCode = "I001-00011"
	UnAuthorized MessageCode = "I001-00012"

	DuplicateKeyError MessageCode = "W001-00001"

	DBError      MessageCode = "E001-00001"
	UnknownError MessageCode = "E001-00002"
)

var Messages = map[MessageCode]string{
	RequestStart: "リクエスト開始",
	RequestEnd:   "リクエスト終了 status code:%v",

	BadRequest:   "Bad request",
	NotFound:     "Not Found",
	UnAuthorized: "UnAuthorized",

	DuplicateKeyError: "Duplicate key",

	DBError:      "DB ssserror",
	UnknownError: "UnknownError Error Detail:%v",
}
