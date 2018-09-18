package ex

type Response struct {
	Code    int    `json:"code,omitempty" description:"code of the response"`
	Message string `json:"message,omitempty" description:"message of the response"`
}

func ReturnBadRequest() (int, Response) {
	ar := Response{
		Code:    INVALID_PARAMS,
		Message: GetMsg(INVALID_PARAMS),
	}
	return INVALID_PARAMS, ar
}

func ReturnInternalServerError(err error) (int, Response) {
	ar := Response{
		Code:    INTERNAL_ERROR,
		Message: GetMsg(INTERNAL_ERROR) + ": " + err.Error(),
	}
	return INTERNAL_ERROR, ar
}

func ReturnConfigExistError() (int, Response) {
	ar := Response{
		Code:    ERROR_CONFIG_EXIST,
		Message: GetMsg(ERROR_CONFIG_EXIST),
	}
	return ERROR_CONFIG_EXIST, ar
}

func ReturnConfigNotExistError() (int, Response) {
	ar := Response{
		Code:    ERROR_CONFIG_NOT_EXIST,
		Message: GetMsg(ERROR_CONFIG_NOT_EXIST),
	}
	return ERROR_CONFIG_NOT_EXIST, ar
}

func ReturnLimitError() (int, Response) {
	ar := Response{
		Code:    ERROR_LIMIT,
		Message: GetMsg(ERROR_LIMIT),
	}
	return ERROR_LIMIT, ar
}

func ReturnOK() (int, Response) {
	ar := Response{
		Code:    SUCCESS,
		Message: GetMsg(SUCCESS),
	}
	return SUCCESS, ar
}
