package header

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/context"
	"github.com/dubbogo/dubbo-go-proxy/pkg/context/http"
	"strings"
)

type headerFilter struct {
}

func New() *headerFilter {
	return &headerFilter{}
}

func (h *headerFilter) Do() context.FilterFunc {
	return func(c context.Context) {

		api := c.GetAPI()
		headers := api.Headers
		if len(headers) <= 0 {
			c.Next()
			return
		}
		switch c.(type) {
		case *http.HttpContext:
			hc := c.(*http.HttpContext)
			urlHeaders := hc.AllHeaders()
			if len(urlHeaders) <= 0 {
				c.Abort()
				return
			}

			for headerName, headerValue := range headers {
				urlHeaderValues := urlHeaders.Values(strings.ToLower(headerName))
				if urlHeaderValues == nil {
					c.Abort()
					return
				}
				for _, urlHeaderValue := range urlHeaderValues {
					if urlHeaderValue == headerValue {
						goto FOUND
					}
				}
				c.Abort()
			FOUND:
				continue
			}
			break
		}
	}
}
