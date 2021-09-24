package dubboregistry

import (
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/dubboregistry/registry"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/adapter"
	"github.com/apache/dubbo-go-pixiu/pkg/filter/http/apiconfig/api"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
	"github.com/pkg/errors"
)

func init() {
	adapter.RegisterAdapterPlugin(&Plugin{})
}

var (
	_ adapter.AdapterPlugin = new(Plugin)
	_ adapter.Adapter       = new(Adapter)
)

type AdaptorConfig struct {
	AppointedListener string                    `yaml:"appointed_listener" json:"appointed_listener" :"appointed_listener" mapstructure: "appointed_listener"`
	Registries        map[string]model.Registry `yaml:"registries" json:"registries" mapstructure:"registries" :"registries"`
}

// Plugin to monitor dubbo services on registry center
type Plugin struct{
	adapterInstance *Adapter
}

// Kind returns the identifier of the plugin
func (p Plugin) Kind() string {
	return constant.DubboRegistryCenterAdapter
}

// CreateAdapter returns the dubbo registry center adapter
func (p *Plugin) CreateAdapter(config interface{}, bs *model.Bootstrap) (adapter.Adapter, error) {
	c, ok := config.(AdaptorConfig)
	listenerService := server.GetServer().GetListenerManager().GetListenerService(c.AppointedListener)
	if listenerService == nil {
		return nil, errors.New("Appointed Listener not found")
	}
	if !ok {
		return nil, errors.New("Configuration incorrect")
	}
	adapter := &Adapter{cfg: c, boundListenerService: c.AppointedListener}
	p.adapterInstance = adapter
	return adapter, nil
}

// RegisterAPIDiscovery register the api discovery service to the plugin so that
//   the plugin can update the api config
func RegisterAPIDiscovery(apiDiscovery api.APIDiscoveryService) {
	//adapter := adapter.GetAdapterPlugin(constant.DubboRegistryCenterAdapter)
	//adapter.
}

// Adapter to monitor dubbo services on registry center
type Adapter struct {
	cfg                  AdaptorConfig
	registries           map[string]registry.Registry
	boundListenerService string
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
		a.registries[k], err = registry.GetRegistry(k, registryConfig)
		a.registries[k].SetPixiuListenerName(a.boundListenerService)
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
