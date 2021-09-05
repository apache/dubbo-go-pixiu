package baseregistry

import (
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/dubboregistry/registry"
	"sync"
	//"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/router"
)

type FacadeRegistry interface {
	// DoSubscribe subscribes the registry cluster to monitor the changes.
	DoSubscribe() error
	// DoUnsubscribe unsubscribes the registry cluster.
	DoUnsubscribe() error
}

type SvcListeners struct {
	// listeners use url.ServiceKey as the index.
	listeners map[string]registry.Listener
	listenerLock sync.Mutex
}

// GetListener returns existing listener or nil
func (s *SvcListeners) GetListener(id string) registry.Listener {
	s.listenerLock.Lock()
	defer s.listenerLock.Unlock()
	listener, ok := s.listeners[id]
	if !ok {
		return nil
	}
	return listener
}

// SetListener set the listener to the registry
func (s *SvcListeners) SetListener(id string, listener registry.Listener) {
	s.listenerLock.Lock()
	defer s.listenerLock.Unlock()
	s.listeners[id] = listener
}

// RemoveListener removes the listener of the registry
func (s *SvcListeners) RemoveListener(id string) {
	s.listenerLock.Lock()
	defer s.listenerLock.Unlock()
	delete(s.listeners, id)
}

func (s *SvcListeners) GetAllListener() map[string]registry.Listener {
	s.listenerLock.Lock()
	defer s.listenerLock.Unlock()
	return s.listeners
}

type BaseRegistry struct {
	svcListeners   *SvcListeners
	facadeRegistry FacadeRegistry
}

func NewBaseRegistry(facade FacadeRegistry) *BaseRegistry {
	return &BaseRegistry{
		facadeRegistry: facade,
		svcListeners:   &SvcListeners{
			listeners: make(map[string]registry.Listener),
		},
	}
}

// GetSvcListener returns existing listener or nil
func (r *BaseRegistry) GetSvcListener(id string) registry.Listener {
	return r.svcListeners.GetListener(id)
}

// GetAllSvcListener get all the listener of the registry
func (r *BaseRegistry) GetAllSvcListener() map[string]registry.Listener {
	return r.svcListeners.GetAllListener()
}

// SetSvcListener set the listener to the registry
func (r *BaseRegistry) SetSvcListener(id string, listener registry.Listener) {
	r.svcListeners.SetListener(id, listener)
}

// RemoveSvcListener remove the listener of the registry
func (r *BaseRegistry) RemoveSvcListener(id string) {
	r.svcListeners.RemoveListener(id)
}

// Subscribe monitors the target registry.
func (r *BaseRegistry) Subscribe() error {
	return r.facadeRegistry.DoSubscribe()
}

// Unsubscribe stops monitoring the target registry.
func (r *BaseRegistry) Unsubscribe() error {
	return r.facadeRegistry.DoUnsubscribe()
}
