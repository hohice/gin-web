package ex

import "net/http"

const (
	SUCCESS        = http.StatusOK
	INTERNAL_ERROR = http.StatusInternalServerError
	INVALID_PARAMS = http.StatusBadRequest
	ERROR_LIMIT    = http.StatusTooManyRequests

	ERROR_AUTH_CHECK_TOKEN_FAIL    = 10001
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT = 10002
	ERROR_AUTH_TOKEN               = 10003
	ERROR_AUTH                     = 10004

	ERROR_CONFIG_NOT_EXIST = http.StatusNotFound
	ERROR_CONFIG_EXIST     = http.StatusFound
)

var MsgFlags = map[int]string{
	SUCCESS:        "ok",
	INTERNAL_ERROR: "Internal Server error",
	INVALID_PARAMS: "Invalid Name supplied!",
	ERROR_LIMIT:    "API rate limit exceeded",

	ERROR_AUTH_CHECK_TOKEN_FAIL:    "Token auth failed",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT: "Token time out",
	ERROR_AUTH_TOKEN:               "Token generate failed",
	ERROR_AUTH:                     "Token error",

	ERROR_CONFIG_NOT_EXIST: "config is not exist",
	ERROR_CONFIG_EXIST:     "config already exist",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[INTERNAL_ERROR]
}
