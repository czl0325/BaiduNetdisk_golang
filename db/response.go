package db

type BaseResponse struct {
	Code    int
	Message string
	Data    interface{}
}
