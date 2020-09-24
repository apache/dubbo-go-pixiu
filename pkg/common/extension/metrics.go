package extension

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/metrics"
)

var metricRegistryMap = map[string]metrics.Registry{}

func init() {
	metricRegistryMap = make(map[string]metrics.Registry)
}

// SetRegistry will store registry.
func SetRegistry(name string, registry metrics.Registry) {
	metricRegistryMap[name] = registry
}

// GetMustRegistry will return the metrics.Registry
// if not found, it will panic.
func GetMustRegistry(name string) metrics.Registry {
	if registry, ok := metricRegistryMap[name]; ok {
		return registry
	}

	panic("metric registry for " + name + " is not existing!")
}
