package baseregistry

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension"
	"github.com/apache/dubbo-go-pixiu/pkg/registry"
	"sync"
	//"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/router"
)

type FacadeRegistry interface {
	// LoadInterfaces loads the dubbo services from interfaces level
	LoadInterfaces() ([]registry.RegisteredAPI, []error)
	// LoadApplications loads the dubbo services from application level
	LoadApplications() ([]registry.RegisteredAPI, []error)
	// DoSubscribe subscribes the registries to monitor the changes.
	DoSubscribe(registry.RegisteredAPI) error
	// DoSDSubscribe subscribes the service discovery registries to monitor the changes.
	DoSDSubscribe(registry.RegisteredAPI) error
	// DoUnsubscribe unsubscribes the registries.
	DoUnsubscribe(registry.RegisteredAPI) error
	// DoSDUnsubscribe unsubscribes the registries.
	DoSDUnsubscribe(registry.RegisteredAPI) error
}

type SrvListeners struct {
	// listeners use url.ServiceKey as the index.
	listeners map[string]registry.Listener
	listenerLock sync.Mutex
}

// GetListener returns existing listener or nil
func (s *SrvListeners) GetListener(id string) registry.Listener {
	listener, ok := s.listeners[id]
	if !ok {
		return nil
	}
	return listener
}

// SetListener set the listener to the registry and start to read process.
func (s *SrvListeners) SetListener(id string, listener registry.Listener) {
	s.listenerLock.Lock()
	defer s.listenerLock.Unlock()
	if _, ok := s.listeners[id]; !ok {
		s.listeners[id] = listener
	}
}


type BaseRegistry struct {
	srvListeners *SrvListeners
	facadeRegistry FacadeRegistry
}

func NewBaseRegistry(facade FacadeRegistry) *BaseRegistry {
	return &BaseRegistry{
		facadeRegistry: facade,
		srvListeners:   &SrvListeners{
			listeners: make(map[string]registry.Listener),
		},
	}
}

// LoadServices loads all the registered Dubbo services from registry
func (r *BaseRegistry) LoadServices() error {
	interfaceAPIs, _ := r.facadeRegistry.LoadInterfaces()
	applicationAPIs, _ := r.facadeRegistry.LoadApplications()
	apis := r.deduplication(append(interfaceAPIs, applicationAPIs...))
	localAPIDiscSrv := extension.GetMustAPIDiscoveryService(constant.LocalMemoryApiDiscoveryService)
	for i := range apis {
		err := localAPIDiscSrv.AddAPI(apis[i].API)
		if err != nil {
			return err
		}
		r.Subscribe(apis[i])
	}
	return nil
}

func (r *BaseRegistry) deduplication(apis []registry.RegisteredAPI) []registry.RegisteredAPI {
	urlPatternMap := make(map[string]bool)
	var rstAPIs []registry.RegisteredAPI
	for i := range apis {
		if _, ok := urlPatternMap[apis[i].API.URLPattern]; !ok {
			rstAPIs = append(rstAPIs,
				registry.RegisteredAPI{
					API:            apis[i].API,
					RegisteredPath: apis[i].RegisteredPath,
					RegisteredType: apis[i].RegisteredType,
				})
			urlPatternMap[apis[i].API.URLPattern] = true
		}
	}
	return rstAPIs
}

// GetListener returns existing listener or nil
func (r *BaseRegistry) GetListener(id string) registry.Listener {
	return r.srvListeners.GetListener(id)
}

// SetListener set the listener to the registry and start to read process.
func (r *BaseRegistry) SetListener(id string, listener registry.Listener) {
	r.srvListeners.SetListener(id, listener)
}

// Subscribe monitors the target registry.
func (r *BaseRegistry) Subscribe(service registry.RegisteredAPI) error {
	return r.facadeRegistry.DoSubscribe(service)
}

// Unsubscribe stops monitoring the target registry.
func (r *BaseRegistry) Unsubscribe(service registry.RegisteredAPI) error {
	return r.facadeRegistry.DoUnsubscribe(service)
}
