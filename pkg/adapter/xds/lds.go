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
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/xds/apiclient"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	xdspb "github.com/apache/dubbo-go-pixiu/pkg/model/xds/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
	"gopkg.in/yaml.v2"
)

type LdsManager struct {
	DiscoverApi
}

// Fetch overwrite DiscoverApi.Fetch.
func (l *LdsManager) Fetch() error {
	r, err := l.DiscoverApi.Fetch("") //todo use local version
	if err != nil {
		logger.Error("can not fetch lds", err)
		return err
	}
	listeners := make([]*xdspb.Listener, 0, len(r))
	for _, one := range r {
		_listener := &xdspb.Listener{}
		if err := one.To(_listener); err != nil {
			logger.Errorf("unknown resource of %s, expect Listener", one.GetName())
			continue
		}
		logger.Infof("listener xds server %v", _listener)
		listeners = append(listeners, _listener)
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
		l.removeListeners(delta.RemovedResource)
		listeners := make([]*xdspb.Listener, 0, len(delta.NewResources))
		for _, one := range delta.NewResources {
			_listener := &xdspb.Listener{}
			if err := one.To(_listener); err != nil {
				logger.Errorf("unknown resource of %s, expect Listener", one.GetName())
				continue
			}
			logger.Infof("listener xds server %v", _listener)
			listeners = append(listeners, _listener)
		}
		if len(listeners) > 0 {
			l.setupListeners(listeners)
		}
	}
}

func (l *LdsManager) makeSocketAddress(address *xdspb.SocketAddress) model.SocketAddress {
	if address == nil {
		return model.SocketAddress{}
	}
	return model.SocketAddress{
		ProtocolStr:  address.ProtocolStr,
		Protocol:     model.ProtocolType(address.Protocol), //todo
		Address:      address.Address,
		Port:         int(address.Port),
		ResolverName: address.ResolverName,
		//Domains:      _l.Address.do, todo add the domains
		//CertsDir: _l.Address.SocketAddress"", //todo add the domains
	}
}
func (l *LdsManager) removeListeners(names []string) {
	server.GetServer().GetListenerManager().RemoveListener(names)
}

// setupListeners setup listeners accord to dynamic resource
func (l *LdsManager) setupListeners(listeners []*xdspb.Listener) {
	for _, _l := range listeners {
		modelListener := l.makeListener(_l)
		server.GetServer().GetListenerManager().AddOrUpdateListener(&modelListener)
	}
}

func (l *LdsManager) makeListener(_l *xdspb.Listener) model.Listener {
	return model.Listener{
		Name:         _l.Name,
		Address:      l.makeAddress(_l.Address),
		FilterChains: l.makeFilterChain(_l.FilterChains),
		Config:       nil, // todo set the additional config
	}
}

func (l *LdsManager) makeFilterChain(fChain []*xdspb.FilterChain) []model.FilterChain {
	r := make([]model.FilterChain, 0, len(fChain))
	for _, one := range fChain {
		r = append(r, model.FilterChain{
			FilterChainMatch: model.FilterChainMatch{},
			Filters:          l.makeFilters(one.Filters),
		})
	}
	return r
}

func (l *LdsManager) makeFilters(filters []*xdspb.Filter) []model.Filter {
	_filters := make([]model.Filter, 0, len(filters))
	for _, _filter := range filters {
		_filters = append(_filters, model.Filter{
			Name: _filter.Name,
			//Config: _filter., todo define the config of filter
			Config: l.makeConfig(_filter),
		})
	}
	return _filters
}

func (l *LdsManager) makeConfig(filter *xdspb.Filter) (m map[string]interface{}) {
	if filter.GetYaml() != nil {
		//todo read yaml
		if err := yaml.Unmarshal([]byte(filter.GetYaml().Content), &m); err != nil {
			logger.Errorf("can not make yaml from filter.Config: %s", filter.GetYaml().Content, err)
		}
		return
	}
	logger.Errorf("can not get filter config of %s", filter.Name)
	return
}

func (l *LdsManager) makeAddress(addr *xdspb.Address) model.Address {
	if addr == nil {
		return model.Address{}
	}
	return model.Address{
		SocketAddress: l.makeSocketAddress(addr.SocketAddress),
		Name:          addr.Name,
	}
}
