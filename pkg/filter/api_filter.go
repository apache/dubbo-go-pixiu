package filter

import (
	"net/http"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/constant"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/extension"
	"github.com/dubbogo/dubbo-go-proxy/pkg/context"
	"github.com/dubbogo/dubbo-go-proxy/pkg/model"
)

func init() {
	extension.SetFilterFunc(constant.HttpApiFilter, ApiFilter())
}

func ApiFilter() context.FilterFunc {
	return func(c context.Context) {
		url := c.GetUrl()
		method := c.GetMethod()

		if api, ok := model.EmptyApi.FindApi(url); ok {
			if !api.MatchMethod(method) {
				c.WriteWithStatus(http.StatusMethodNotAllowed, constant.Default405Body)
				c.AddHeader(constant.HeaderKeyContextType, constant.HeaderValueTextPlain)
				c.Abort()
				return
			}

			if !api.IsOk(api.Name) {
				c.WriteWithStatus(http.StatusNotAcceptable, constant.Default406Body)
				c.AddHeader(constant.HeaderKeyContextType, constant.HeaderValueTextPlain)
				c.Abort()
				return
			}

			c.Api(api)
			c.Next()
		} else {
			c.WriteWithStatus(http.StatusNotFound, constant.Default404Body)
			c.AddHeader(constant.HeaderKeyContextType, constant.HeaderValueTextPlain)
			c.Abort()
		}
	}
}
