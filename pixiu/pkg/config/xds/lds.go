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
	xdsModel "github.com/dubbo-go-pixiu/pixiu-api/pkg/xds/model"
	"gopkg.in/yaml.v2"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/config/xds/apiclient"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/server/controls"
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
	listeners := make([]*xdsModel.Listener, 0, len(r))
	for _, one := range r {
		listener := &xdsModel.PixiuExtensionListeners{}
		if err := one.To(listener); err != nil {
			logger.Errorf("unknown resource of %s, expect Listener", one.GetName())
			continue
		}
		logger.Infof("listener xds server %v", listener)
		listeners = append(listeners, listener.Listeners...)
	}
	l.setupListeners(listeners)
	return nil
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
		listeners := make([]*xdsModel.Listener, 0, len(delta.NewResources))
		for _, one := range delta.NewResources {
			listener := &xdsModel.PixiuExtensionListeners{}
			if err := one.To(listener); err != nil {
				logger.Errorf("unknown resource of %s, expect Listener", one.GetName())
				continue
			}
			logger.Infof("listener xds server %v", listener)
			listeners = append(listeners, listener.Listeners...)
		}

		l.setupListeners(listeners)
	}
}

func (l *LdsManager) makeSocketAddress(address *xdsModel.SocketAddress) model.SocketAddress {
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
	l.listenerMg.RemoveListener(names)
}

// setupListeners setup listeners accord to dynamic resource
func (l *LdsManager) setupListeners(listeners []*xdsModel.Listener) {
	//Make sure each one has a unique name like "host-port-protocol"
	for _, v := range listeners {
		v.Name = resolveListenerName(v.Address.SocketAddress.Address, int(v.Address.SocketAddress.Port), v.Protocol.String())
	}

	laterApplies := make([]func() error, 0, len(listeners))
	toRemoveHash := make(map[string]struct{}, len(listeners))

	lm := l.listenerMg
	activeListeners, err := lm.CloneXdsControlListener()
	if err != nil {
		logger.Errorf("Clone Xds Control Listener fail: %s", err)
		return
	}
	//put all current listeners to $toRemoveHash
	for _, v := range activeListeners {
		//Make sure each one has a unique name like "host-port-protocol"
		v.Name = resolveListenerName(v.Address.SocketAddress.Address, v.Address.SocketAddress.Port, v.ProtocolStr)
		toRemoveHash[v.Name] = struct{}{}
	}

	for _, listener := range listeners {
		delete(toRemoveHash, listener.Name)

		modelListener := l.makeListener(listener)
		// add or update later after removes
		switch {
		case lm.HasListener(modelListener.Name):
			laterApplies = append(laterApplies, func() error {
				return lm.UpdateListener(&modelListener)
			})
		default:
			laterApplies = append(laterApplies, func() error {
				return lm.AddListener(&modelListener)
			})
		}
	}
	// remove the listeners first to prevent tcp port conflict
	l.removeListeners(toRemoveHash)
	//do update and add new cluster.
	for _, fn := range laterApplies {
		if err := fn(); err != nil {
			logger.Errorf("can not modify listener", err)
		}
	}
}

func resolveListenerName(host string, port int, protocol string) string {
	return host + "-" + strconv.Itoa(port) + "-" + protocol
}

func (l *LdsManager) makeListener(listener *xdsModel.Listener) model.Listener {
	return model.Listener{
		Name:        listener.Name,
		ProtocolStr: listener.Protocol.String(),
		Protocol:    model.ProtocolType(model.ProtocolTypeValue[listener.Protocol.String()]),
		Address:     l.makeAddress(listener.Address),
		FilterChain: l.makeFilterChain(listener.FilterChain),
		Config:      nil, // todo set the additional config
	}
}

func (l *LdsManager) makeFilterChain(fChain *xdsModel.FilterChain) model.FilterChain {
	return model.FilterChain{
		Filters: l.makeFilters(fChain.Filters),
	}
}

func (l *LdsManager) makeFilters(filters []*xdsModel.NetworkFilter) []model.NetworkFilter {
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

func (l *LdsManager) makeConfig(filter *xdsModel.NetworkFilter) (m map[string]interface{}) {
	switch cfg := filter.Config.(type) {
	case *xdsModel.NetworkFilter_Yaml:
		if err := yaml.Unmarshal([]byte(cfg.Yaml.Content), &m); err != nil {
			logger.Errorf("can not make yaml from filter.Config: %s", cfg.Yaml.Content, err)
		}
	case *xdsModel.NetworkFilter_Json:
		if err := json.Unmarshal([]byte(cfg.Json.Content), &m); err != nil {
			logger.Errorf("can not make json from filter.Config: %s", cfg.Json.Content, err)
		}
	case *xdsModel.NetworkFilter_Struct:
		m = cfg.Struct.AsMap()
	default:
		logger.Errorf("can not get filter config of %s", filter.Name)
	}
	return
}

func (l *LdsManager) makeAddress(addr *xdsModel.Address) model.Address {
	if addr == nil {
		return model.Address{}
	}
	return model.Address{
		SocketAddress: l.makeSocketAddress(addr.SocketAddress),
		Name:          addr.Name,
	}
}
