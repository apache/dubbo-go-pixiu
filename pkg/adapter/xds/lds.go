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
)

import (
	model2 "github.com/dubbogo/dubbo-go-pixiu-filter/pkg/xds/model"

	"gopkg.in/yaml.v2"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/xds/apiclient"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
)

type LdsManager struct {
	DiscoverApi
}

// Fetch overwrite DiscoverApi.Fetch.
func (l *LdsManager) Fetch() error {
	r, err := l.DiscoverApi.Fetch("") //todo use local version
	if err != nil {
		return err
	}
	listeners := make([]*model2.Listener, 0, len(r))
	for _, one := range r {
		_listener := &model2.PixiuExtensionListeners{}
		if err := one.To(_listener); err != nil {
			logger.Errorf("unknown resource of %s, expect Listener", one.GetName())
			continue
		}
		logger.Infof("listener xds server %v", _listener)
		listeners = append(listeners, _listener.Listeners...)
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
		listeners := make([]*model2.Listener, 0, len(delta.NewResources))
		for _, one := range delta.NewResources {
			_listener := &model2.PixiuExtensionListeners{}
			if err := one.To(_listener); err != nil {
				logger.Errorf("unknown resource of %s, expect Listener", one.GetName())
				continue
			}
			logger.Infof("listener xds server %v", _listener)
			listeners = append(listeners, _listener.Listeners...)
		}

		l.setupListeners(listeners)
	}
}

func (l *LdsManager) makeSocketAddress(address *model2.SocketAddress) model.SocketAddress {
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
		server.GetServer().GetListenerManager().RemoveListener(names)
	}
}

// setupListeners setup listeners accord to dynamic resource
func (l *LdsManager) setupListeners(listeners []*model2.Listener) {
	laterApplies := make([]func() error, 0, len(listeners))
	toRemoveHash := make(map[string]struct{}, len(listeners))

	for _, _l := range listeners {
		toRemoveHash[_l.Name] = struct{}{}
	}

	for _, _l := range listeners {
		delete(toRemoveHash, _l.Name)
		modelListener := l.makeListener(_l)
		// apply add or update later after removes
		laterApplies = append(laterApplies, func() error {
			err := server.GetServer().GetListenerManager().AddOrUpdateListener(&modelListener)
			if err != nil {
				logger.Errorf("can not add/update listener config=> %v", modelListener)
			}
			return nil
		})
	}
	l.removeListeners(toRemoveHash)
	for _, fn := range laterApplies { //do update and add new cluster.
		if err := fn(); err != nil {
			logger.Errorf("can not modify listener", err)
		}
	}
}

func (l *LdsManager) makeListener(_l *model2.Listener) model.Listener {
	return model.Listener{
		Name:        _l.Name,
		Address:     l.makeAddress(_l.Address),
		FilterChain: l.makeFilterChain(_l.FilterChain),
		Config:      nil, // todo set the additional config
	}
}

func (l *LdsManager) makeFilterChain(fChain *model2.FilterChain) model.FilterChain {
	return model.FilterChain{
		Filters: l.makeFilters(fChain.Filters),
	}
}

func (l *LdsManager) makeFilters(filters []*model2.NetworkFilter) []model.NetworkFilter {
	_filters := make([]model.NetworkFilter, 0, len(filters))
	for _, _filter := range filters {
		_filters = append(_filters, model.NetworkFilter{
			Name: _filter.Name,
			//Config: _filter., todo define the config of filter
			Config: l.makeConfig(_filter),
		})
	}
	return _filters
}

func (l *LdsManager) makeConfig(filter *model2.NetworkFilter) (m map[string]interface{}) {
	switch cfg := filter.Config.(type) {
	case *model2.NetworkFilter_Yaml:
		if err := yaml.Unmarshal([]byte(cfg.Yaml.Content), &m); err != nil {
			logger.Errorf("can not make yaml from filter.Config: %s", cfg.Yaml.Content, err)
		}
	case *model2.NetworkFilter_Json:
		if err := json.Unmarshal([]byte(cfg.Json.Content), &m); err != nil {
			logger.Errorf("can not make json from filter.Config: %s", cfg.Json.Content, err)
		}
	case *model2.NetworkFilter_Struct:
		m = cfg.Struct.AsMap()
	default:
		logger.Errorf("can not get filter config of %s", filter.Name)
	}
	return
}

func (l *LdsManager) makeAddress(addr *model2.Address) model.Address {
	if addr == nil {
		return model.Address{}
	}
	return model.Address{
		SocketAddress: l.makeSocketAddress(addr.SocketAddress),
		Name:          addr.Name,
	}
}
