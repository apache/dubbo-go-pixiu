package client

type Request struct {
	Body   []byte
	Header map[string]string
	Api    *Api
}

func NewRequest(b []byte, api *Api) *Request {
	return &Request{
		Body: b,
		Api:  api,
	}
}
