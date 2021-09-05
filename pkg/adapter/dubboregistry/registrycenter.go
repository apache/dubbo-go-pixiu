package dubboregistry

import (
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/dubboregistry/registry"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/adapter"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/pkg/errors"
)

func init() {
	adapter.RegisterAdapterPlugin(&Plugin{})
}

var (
	_ adapter.AdapterPlugin = new(Plugin)
	_ adapter.Adapter = new(Adapter)
)

type AdaptorConfig = map[string]model.Registry

// Plugin to monitor dubbo services on registry center
type Plugin struct {}

// Kind returns the identifier of the plugin
func (p Plugin) Kind() string {
	return constant.DubboRegistryCenterAdapter
}

// CreateAdapter returns the dubbo registry center adapter
func (p Plugin) CreateAdapter(config interface{}, bs *model.Bootstrap) (adapter.Adapter, error) {
	c, ok := config.(AdaptorConfig)
	if !ok {
		return nil, errors.New("Configuration incorrect")
	}
	return &Adapter{cfg: c}, nil
}

// Adapter to monitor dubbo services on registry center
type Adapter struct{
	cfg AdaptorConfig
	registries map[string]registry.Registry
}

// Start starts the adaptor
func (a Adapter) Start() {
	for _, reg := range a.registries {
		reg.Subscribe()
	}
}

// Stop stops the adaptor
func (a *Adapter) Stop() {
	for _, reg := range a.registries {
		reg.Unsubscribe()
	}
}

// Apply inits the registries according to the configuration
func (a *Adapter) Apply() error {
	// create registry per config
	for k, registryConfig := range a.cfg {
		a.registries[k] = registry.GetRegistry(k, registryConfig)
	}
	panic("implement me")
}

// Config returns the config of the adaptor
func (a Adapter) Config() interface{} {
	return a.cfg
}
