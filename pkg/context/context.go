package context

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/client"
)

// Context run context
type Context interface {
	Next()
	Abort()

	AppendFilterFunc(ff ...FilterFunc)

	Status(code int)
	StatusCode() int
	WriteWithStatus(int, []byte) (int, error)
	Write([]byte) (int, error)
	AddHeader(k, v string)
	GetHeader(k string) string
	GetUrl() string
	GetMethod() string

	BuildFilters()

	Api(api *client.Api)
	GetApi() *client.Api

	WriteErr(p interface{})
}
