package extension

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/registry"
)

var registryMap = make(map[string]func(model.Registry) (registry.Registry, error), 8)

// SetFilterFunc will store the @filter and @name
func SetRegistry(name string, newRegFunc func(model.Registry) (registry.Registry, error)) {
	registryMap[name] = newRegFunc
}

// GetMustFilterFunc will return the pixiu.FilterFunc
// if not found, it will panic
func GetRegistry(name string, regConfig model.Registry) registry.Registry {
	if registry, ok := registryMap[name]; ok {
		reg, err := registry(regConfig)
		if err != nil {
			panic("Initialize Registry" + name + "failed due to: " + err.Error())
		}
		return reg
	}
	panic("Registry " + name + " does not support yet")
}
