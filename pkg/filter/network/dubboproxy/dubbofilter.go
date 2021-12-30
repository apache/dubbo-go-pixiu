package dubboproxy

import "github.com/apache/dubbo-go-pixiu/pkg/context/http"

// DubboFilter describe http filter
type DubboFilter interface {
		// PrepareFilterChain add filter into chain
		PrepareFilterChain(ctx *http.HttpContext) error

		// Handle filter hook function
		Handle(ctx *http.HttpContext)

		Apply() error

		// Config get config for filter
		Config() interface{}
}