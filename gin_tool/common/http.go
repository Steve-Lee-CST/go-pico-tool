package common

type CommonResponse[T any] struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    *T     `json:"data,omitempty"`
}

func SuccessResponse[T any](data *T) CommonResponse[T] {
	return CommonResponse[T]{
		Code:    0,
		Message: "success",
		Data:    data,
	}
}

func ErrorResponse[T any](code int, message string, data *T) CommonResponse[T] {
	return CommonResponse[T]{
		Code:    code,
		Message: message,
		Data:    data,
	}
}
