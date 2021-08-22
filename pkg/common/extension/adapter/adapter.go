package adapter

import (
	"fmt"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/go-errors/errors"
)

type (
	AdapterPlugin interface {
		// Kind returns the unique kind name to represent itself.
		Kind() string

		// CreateAdapter return the Adapter callback
		CreateAdapter(config interface{}, bs *model.Bootstrap) (Adapter, error)
	}

	Adapter interface {
		Start()
		Stop()
	}
)

var (
	adapterPlugins = map[string]AdapterPlugin{}
)

// Register registers adapter plugin
func RegisterAdapterPlugin(p AdapterPlugin) {
	if p.Kind() == "" {
		panic(fmt.Errorf("%T: empty kind", p))
	}

	existedPlugin, existed := adapterPlugins[p.Kind()]
	if existed {
		panic(fmt.Errorf("%T and %T got same kind: %s", p, existedPlugin, p.Kind()))
	}

	adapterPlugins[p.Kind()] = p
}

// GetAdapterPlugin get plugin by kind
func GetAdapterPlugin(kind string) (AdapterPlugin, error) {
	existedAdapter, existed := adapterPlugins[kind]
	if existed {
		return existedAdapter, nil
	} else {
		return nil, errors.Errorf("plugin not found %s", kind)
	}
}
