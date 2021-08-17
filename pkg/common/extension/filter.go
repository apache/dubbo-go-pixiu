package extension

import (
	"fmt"
	http2 "github.com/apache/dubbo-go-pixiu/pkg/common/http"
	context2 "github.com/apache/dubbo-go-pixiu/pkg/context"
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

type (
	// HttpFilterPlugin describe plugin
	HttpFilterPlugin interface {
		// Kind returns the unique kind name to represent itself.
		Kind() string

		// CreateFilterFactory return the filter callback
		CreateFilter(hcm *http2.HttpConnectionManager, config interface{}, bs *model.Bootstrap) (HttpFilter, error)
	}

	// HttpFilter describe http filter
	HttpFilter interface {
		// PrepareFilterChain add filter into chain
		PrepareFilterChain(ctx *http.HttpContext) error

		// Handle filter hook function
		Handle(ctx context2.Context)
	}

	// NetworkFilter describe network filter plugin
	NetworkFilterPlugin interface {
		// Kind returns the unique kind name to represent itself.
		Kind() string
		// CreateFilterFactory return the filter callback
		CreateFilter(config interface{}, bs *model.Bootstrap) (NetworkFilter, error)
	}

	// NetworkFilter describe network filter
	NetworkFilter interface {
		// OnData handle the http context from worker
		OnData(hc *http.HttpContext) error
	}

	// FilterFunc filter func, filter
	FilterFunc func(ctx context2.Context)

	// FilterChain filter chain
	FilterChain []FilterFunc
)

var (
	httpFilterPluginRegistry    = map[string]HttpFilterPlugin{}
	networkFilterPluginRegistry = map[string]NetworkFilterPlugin{}
)

// Register registers filter.
func RegisterHttpFilter(f HttpFilterPlugin) {
	if f.Kind() == "" {
		panic(fmt.Errorf("%T: empty kind", f))
	}

	existedFilter, existed := httpFilterPluginRegistry[f.Kind()]
	if existed {
		panic(fmt.Errorf("%T and %T got same kind: %s", f, existedFilter, f.Kind()))
	}

	httpFilterPluginRegistry[f.Kind()] = f
}

// GetHttpFilterPlugin get plugin by kind
func GetHttpFilterPlugin(kind string) (HttpFilterPlugin, error) {
	existedFilter, existed := httpFilterPluginRegistry[kind]
	if existed {
		return existedFilter, nil
	} else {
		return nil, errors.Errorf("plugin not found %s", kind)
	}
}

// Register registers network filter.
func RegisterNetworkFilter(f NetworkFilterPlugin) {
	if f.Kind() == "" {
		panic(fmt.Errorf("%T: empty kind", f))
	}

	existedFilter, existed := networkFilterPluginRegistry[f.Kind()]
	if existed {
		panic(fmt.Errorf("%T and %T got same kind: %s", f, existedFilter, f.Kind()))
	}

	networkFilterPluginRegistry[f.Kind()] = f
}

// GetNetworkFilterPlugin get plugin by kind
func GetNetworkFilterPlugin(kind string) (NetworkFilterPlugin, error) {
	existedFilter, existed := networkFilterPluginRegistry[kind]
	if existed {
		return existedFilter, nil
	} else {
		return nil, errors.Errorf("plugin not found %s", kind)
	}
}
