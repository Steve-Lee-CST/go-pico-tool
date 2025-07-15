package common

type CommonResponse[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data *T     `json:"data,omitempty"`
}
