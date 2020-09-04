package extension

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/context"
)

var (
	filterFuncCacheMap = make(map[string]func(ctx context.Context))
)

// SetFilterFunc will store the @filter and @name
func SetFilterFunc(name string, filter context.FilterFunc) {
	filterFuncCacheMap[name] = filter
}

// GetMustFilterFunc will return the proxy.FilterFunc
// if not found, it will panic
func GetMustFilterFunc(name string) context.FilterFunc {
	if filter, ok := filterFuncCacheMap[name]; ok {
		return filter
	}

	panic("filter func for " + name + " is not existing!")
}
