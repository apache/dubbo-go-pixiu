package client

type Response struct {
	Data interface{}
}

func NewResponse(data interface{}) *Response {
	return &Response{Data: data}
}

var EmptyResponse = &Response{}
