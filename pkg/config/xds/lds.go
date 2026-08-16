/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package xds

import (
	"encoding/json"
	"strconv"
)

import (
	"gopkg.in/yaml.v3"

	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/config/xds/apiclient"
	xdsmodel "github.com/apache/dubbo-go-pixiu/pkg/config/xds/model"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server/controls"
)

type LdsManager struct {
	DiscoverApi
	listenerMg controls.ListenerManager
}

// Fetch overwrite DiscoverApi.Fetch.
func (l *LdsManager) Fetch() error {
	r, err := l.DiscoverApi.Fetch("") //todo use local version
	if err != nil {
		return err
	}
	listeners := make([]*xdsmodel.Listener, 0, len(r))
	for _, one := range r {
		listener := &xdsmodel.PixiuExtensionListeners{}
		if err := one.To(listener); err != nil {
			logger.Errorf("unknown resource of %s, expect Listener", one.GetName())
			continue
		}
		logger.Infof("listener xds server %v", listener)
		listeners = append(listeners, listener.Listeners...)
	}
	return l.setupListeners(listeners)
}

func (l *LdsManager) Delta() error {
	readCh, err := l.DiscoverApi.Delta()
	if err != nil {
		return err
	}
	go l.asyncHandler(readCh)
	return nil
}

func (l *LdsManager) asyncHandler(read chan *apiclient.DeltaResources) {
	for delta := range read {
		err := l.applyDelta(delta)
		delta.Complete(err)
		if err != nil {
			logger.Errorf("can not setup listener: %v", err)
		}
	}
}

func (l *LdsManager) applyDelta(delta *apiclient.DeltaResources) error {
	if delta == nil {
		return nil
	}

	listeners := make([]*xdsmodel.Listener, 0, len(delta.NewResources))
	for _, resource := range delta.NewResources {
		listener := &xdsmodel.PixiuExtensionListeners{}
		if err := resource.To(listener); err != nil {
			return errors.Wrapf(err, "unknown resource %q, expect Listener", resource.GetName())
		}
		logger.Infof("listener xds server %v", listener)
		listeners = append(listeners, listener.Listeners...)
	}

	for _, name := range delta.RemovedResources {
		if name == constant.ListenerType {
			l.listenerMg.RemoveXDSListeners(l.listenerMg.XDSListenerNames())
		}
	}
	if len(delta.NewResources) == 0 {
		return nil
	}

	return l.setupListeners(listeners)
}

func (l *LdsManager) makeSocketAddress(address *xdsmodel.SocketAddress) model.SocketAddress {
	if address == nil {
		return model.SocketAddress{}
	}
	return model.SocketAddress{
		Address:      address.Address,
		Port:         int(address.Port),
		ResolverName: address.ResolverName,
		//Domains:      _l.Address.do, todo add the domains
		//CertsDir: _l.Address.SocketAddress"", //todo add the domains
	}
}

func (l *LdsManager) removeListeners(toRemoveHash map[string]struct{}) {
	names := make([]string, 0, len(toRemoveHash))
	for name := range toRemoveHash {
		names = append(names, name)
	}
	l.listenerMg.RemoveXDSListeners(names)
}

// setupListeners setup listeners accord to dynamic resource
func (l *LdsManager) setupListeners(listeners []*xdsmodel.Listener) error {
	//Make sure each one has a unique name like "host-port-protocol"
	for _, v := range listeners {
		v.Name = resolveListenerName(v.Address.SocketAddress.Address, int(v.Address.SocketAddress.Port), v.Protocol.String())
	}

	laterApplies := make([]func() error, 0, len(listeners))
	toRemoveHash := make(map[string]struct{}, len(listeners))

	lm := l.listenerMg
	for _, name := range lm.XDSListenerNames() {
		toRemoveHash[name] = struct{}{}
	}

	for _, listener := range listeners {
		delete(toRemoveHash, listener.Name)

		modelListener := l.makeListener(listener)
		// add or update later after removes
		laterApplies = append(laterApplies, func() error {
			return lm.UpsertXDSListener(&modelListener)
		})
	}
	// remove the listeners first to prevent tcp port conflict
	l.removeListeners(toRemoveHash)
	//do update and add new cluster.
	for _, fn := range laterApplies {
		if err := fn(); err != nil {
			return errors.Wrap(err, "can not modify listener")
		}
	}
	return nil
}

func resolveListenerName(host string, port int, protocol string) string {
	return host + "-" + strconv.Itoa(port) + "-" + protocol
}

func (l *LdsManager) makeListener(listener *xdsmodel.Listener) model.Listener {
	return model.Listener{
		Name:        listener.Name,
		ProtocolStr: listener.Protocol.String(),
		Protocol:    model.ProtocolType(model.ProtocolTypeValue[listener.Protocol.String()]),
		Address:     l.makeAddress(listener.Address),
		FilterChain: l.makeFilterChain(listener.FilterChain),
		Config:      nil, // todo set the additional config
	}
}

func (l *LdsManager) makeFilterChain(fChain *xdsmodel.FilterChain) model.FilterChain {
	return model.FilterChain{
		Filters: l.makeFilters(fChain.Filters),
	}
}

func (l *LdsManager) makeFilters(filters []*xdsmodel.NetworkFilter) []model.NetworkFilter {
	result := make([]model.NetworkFilter, 0, len(filters))
	for _, filter := range filters {
		result = append(result, model.NetworkFilter{
			Name: filter.Name,
			//Config: filter., todo define the config of filter
			Config: l.makeConfig(filter),
		})
	}
	return result
}

func (l *LdsManager) makeConfig(filter *xdsmodel.NetworkFilter) (m map[string]any) {
	switch cfg := filter.Config.(type) {
	case *xdsmodel.NetworkFilter_Yaml:
		if err := yaml.Unmarshal([]byte(cfg.Yaml.Content), &m); err != nil {
			logger.Errorf("can not make yaml from filter.Config: %s", cfg.Yaml.Content, err)
		}
	case *xdsmodel.NetworkFilter_Json:
		if err := json.Unmarshal([]byte(cfg.Json.Content), &m); err != nil {
			logger.Errorf("can not make json from filter.Config: %s", cfg.Json.Content, err)
		}
	case *xdsmodel.NetworkFilter_Struct:
		m = cfg.Struct.AsMap()
	default:
		logger.Errorf("can not get filter config of %s", filter.Name)
	}
	return
}

func (l *LdsManager) makeAddress(addr *xdsmodel.Address) model.Address {
	if addr == nil {
		return model.Address{}
	}
	return model.Address{
		SocketAddress: l.makeSocketAddress(addr.SocketAddress),
		Name:          addr.Name,
	}
}
