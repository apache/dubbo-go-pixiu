package client

// Response response from endpoint
type Response struct {
	Data interface{}
}

// NewResponse create response
func NewResponse(data interface{}) *Response {
	return &Response{Data: data}
}

// empty Response
var EmptyResponse = &Response{}
