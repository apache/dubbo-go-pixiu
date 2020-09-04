package filter

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/constant"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/extension"
	"github.com/dubbogo/dubbo-go-proxy/pkg/context"
	"github.com/dubbogo/dubbo-go-proxy/pkg/context/http"
)

func init() {
	extension.SetFilterFunc(constant.HttpDomainFilter, Domain())
}

// Domain
// https :authority
// http Host
func Domain() context.FilterFunc {
	return func(c context.Context) {
		if MatchDomainFilter(c.(*http.HttpContext)) {
			c.Next()
		}
	}
}

func MatchDomainFilter(c *http.HttpContext) bool {
	for _, v := range c.Listener.FilterChains {
		for _, d := range v.FilterChainMatch.Domains {
			if d == c.GetHeader(constant.Http1HeaderKeyHost) {
				return true
			}
		}
	}

	return false
}
