package dubboregistry

import (
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/dubboregistry/registry"
	_ "github.com/apache/dubbo-go-pixiu/pkg/adapter/dubboregistry/registry/zookeeper"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/adapter"
	"github.com/apache/dubbo-go-pixiu/pkg/filter/http/apiconfig/api"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/api/config"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/router"
)

func init() {
	adapter.RegisterAdapterPlugin(&Plugin{})
}

var (
	_ adapter.AdapterPlugin = new(Plugin)
	_ adapter.Adapter       = new(Adapter)
)

type (
	// Plugin to monitor dubbo services on registry center
	Plugin struct {
	}

	AdaptorConfig struct {
		Registries map[string]model.Registry `yaml:"registries" json:"registries" mapstructure:"registries" :"registries"`
	}
)

// Kind returns the identifier of the plugin
func (p Plugin) Kind() string {
	return constant.DubboRegistryCenterAdapter
}

// CreateAdapter returns the dubbo registry center adapter
func (p *Plugin) CreateAdapter(a *model.Adapter, bs *model.Bootstrap) (adapter.Adapter, error) {
	adapter := &Adapter{id: a.ID,
		registries: make(map[string]registry.Registry),
		cfg:        AdaptorConfig{Registries: make(map[string]model.Registry)}}
	return adapter, nil
}

// Adapter to monitor dubbo services on registry center
type Adapter struct {
	id             string
	cfg            AdaptorConfig
	registries     map[string]registry.Registry
	apiDiscoveries api.APIDiscoveryService
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
	for k, registryConfig := range a.cfg.Registries {
		var err error
		a.registries[k], err = registry.GetRegistry(k, registryConfig, a)
		if err != nil {
			return err
		}
	}

	return nil
}

// Config returns the config of the adaptor
func (a Adapter) Config() interface{} {
	return a.cfg
}

func (a *Adapter) OnAddAPI(r router.API) error {
	acm := server.GetApiConfigManager()
	return acm.AddAPI(a.id, r)
}

func (a *Adapter) OnDeleteRouter(r config.Resource) error {
	acm := server.GetApiConfigManager()
	return acm.DeleteAPI(a.id, r)
}
