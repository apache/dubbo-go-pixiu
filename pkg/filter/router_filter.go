package filter

import (
	nh "net/http"
	"strings"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/constant"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/extension"
	"github.com/dubbogo/dubbo-go-proxy/pkg/context"
	"github.com/dubbogo/dubbo-go-proxy/pkg/context/http"
	"github.com/dubbogo/dubbo-go-proxy/pkg/model"
)

func init() {
	extension.SetFilterFunc(constant.HttpRouterFilter, HttpRouting())
}

// HttpRouting http router filter
func HttpRouting() context.FilterFunc {
	return func(c context.Context) {
		routingFilter(c.(*http.HttpContext))
	}
}

// routingFilter
func routingFilter(c *http.HttpContext) {
	result := true
	for _, v := range c.HttpConnectionManager.RouteConfig.Routes {
		result = routeMatch(c, v)
		if result {
			httpHeaderCorsHandler(c, v)
			break
		}
	}

	if !result {
		c.WriteWithStatus(nh.StatusForbidden, constant.Default403Body)
		c.Abort()
	}
}

// routeMatch will match router with request, only true or false way
func routeMatch(c *http.HttpContext, r model.Router) bool {
	result := true
	if len(r.Match.Headers) > 0 {
		for _, v := range r.Match.Headers {
			result = http.HttpHeaderMatch(c, v)
			if !result {
				break
			}
		}
	}

	if !result {
		return result
	}

	result = http.HttpRouteMatch(c, r.Match)

	if !result {
		return result
	}

	return http.HttpRouteActionMatch(c, r.Route)
}

// httpHeaderCorsHandler will set cors, handler mean can do c.Next()
func httpHeaderCorsHandler(c *http.HttpContext, r model.Router) {
	var acao string
	if r.Route.Cors.Enabled {
		acao = strings.Join(r.Route.Cors.AllowOrigin, "|")
	}

	c.Next()

	if acao != "" {
		c.AddHeader(constant.HeaderKeyAccessControlAllowOrigin, acao)
	}
}
