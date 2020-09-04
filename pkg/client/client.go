package client

type Client interface {
	Init() error
	Close() error

	Call(req *Request) (resp Response, err error)
}
