package io_server

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
	Data   any    `json:"data,inline"`
}

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

func OK(data interface{}) Response {
	return Response{
		Status: StatusOK,
		Data:   data,
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}
