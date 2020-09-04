package context

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/model"
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

	Api(api *model.Api)
	GetApi() *model.Api

	WriteErr(p interface{})
}
