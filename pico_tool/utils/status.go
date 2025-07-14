package utils

type Status struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Desc string `json:"desc"`
}

type StatusWithData[T any] struct {
	Status
	Data T `json:"data"`
}
