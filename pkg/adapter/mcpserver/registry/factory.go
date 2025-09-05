package registry

import (
	"fmt"

	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

var providers = map[string]BuildFunc{}

// RegisterProvider registers a provider-specific controller builder.
func RegisterProvider(protocol string, fn BuildFunc) {
	providers[protocol] = fn
}

// BuildController builds a provider controller based on registry protocol.
func BuildController(reg model.Registry, onChange func(*model.McpServerConfig)) (Controller, error) {
	if fn, ok := providers[reg.Protocol]; ok {
		return fn(reg, onChange)
	}
	return nil, fmt.Errorf("no provider for protocol: %s", reg.Protocol)
}
