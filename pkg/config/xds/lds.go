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
	"github.com/pkg/errors"

	"gopkg.in/yaml.v3"
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

type xdsListenerReplacer interface {
	ReplaceXDSListeners(listeners []*model.Listener) error
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

	if len(delta.NewResources) == 0 && !containsResource(delta.RemovedResources, constant.ListenerType) {
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

// setupListeners setup listeners accord to dynamic resource
func (l *LdsManager) setupListeners(listeners []*xdsmodel.Listener) error {
	//Make sure each one has a unique name like "host-port-protocol"
	for _, v := range listeners {
		if v == nil || v.Address == nil || v.Address.SocketAddress == nil {
			return errors.New("xDS listener must have a socket address")
		}
		if v.Address.SocketAddress.Address == "" || v.Address.SocketAddress.Port <= 0 || v.Address.SocketAddress.Port > 65535 {
			return errors.Errorf("xDS listener has invalid socket address %q:%d", v.Address.SocketAddress.Address, v.Address.SocketAddress.Port)
		}
		if v.FilterChain == nil {
			return errors.Errorf("xDS listener %q has no filter chain", v.Name)
		}
		protocol, err := validateListenerProtocol(v.Protocol)
		if err != nil {
			return err
		}
		v.Name = resolveListenerName(v.Address.SocketAddress.Address, int(v.Address.SocketAddress.Port), protocol)
	}

	converted := make([]*model.Listener, 0, len(listeners))
	for _, listener := range listeners {
		modelListener, err := l.makeListener(listener)
		if err != nil {
			return err
		}
		converted = append(converted, &modelListener)
	}
	replacer, ok := l.listenerMg.(xdsListenerReplacer)
	if !ok {
		return errors.New("listener manager does not support transactional xDS replacement")
	}
	return errors.Wrap(replacer.ReplaceXDSListeners(converted), "can not replace xDS listeners")
}

func resolveListenerName(host string, port int, protocol string) string {
	return host + "-" + strconv.Itoa(port) + "-" + protocol
}

func (l *LdsManager) makeListener(listener *xdsmodel.Listener) (model.Listener, error) {
	protocol, err := validateListenerProtocol(listener.Protocol)
	if err != nil {
		return model.Listener{}, err
	}
	filterChain, err := l.makeFilterChain(listener.FilterChain)
	if err != nil {
		return model.Listener{}, errors.Wrapf(err, "listener %q", listener.Name)
	}
	return model.Listener{
		Name:        listener.Name,
		ProtocolStr: protocol,
		Protocol:    model.ProtocolType(model.ProtocolTypeValue[protocol]),
		Address:     l.makeAddress(listener.Address),
		FilterChain: filterChain,
		Config:      nil, // todo set the additional config
	}, nil
}

func validateListenerProtocol(protocol xdsmodel.Listener_Protocols) (string, error) {
	name, declared := xdsmodel.Listener_Protocols_name[int32(protocol)]
	if !declared {
		return "", errors.Errorf("unsupported xDS listener protocol value %d", protocol)
	}
	if _, supported := model.ProtocolTypeValue[name]; !supported {
		return "", errors.Errorf("unsupported xDS listener protocol %q", name)
	}
	return name, nil
}

func (l *LdsManager) makeFilterChain(fChain *xdsmodel.FilterChain) (model.FilterChain, error) {
	if fChain == nil || len(fChain.Filters) == 0 {
		return model.FilterChain{}, errors.New("filter chain must contain at least one network filter")
	}
	filters, err := l.makeFilters(fChain.Filters)
	if err != nil {
		return model.FilterChain{}, err
	}
	return model.FilterChain{Filters: filters}, nil
}

func (l *LdsManager) makeFilters(filters []*xdsmodel.NetworkFilter) ([]model.NetworkFilter, error) {
	result := make([]model.NetworkFilter, 0, len(filters))
	for index, filter := range filters {
		if filter == nil || filter.Name == "" {
			return nil, errors.Errorf("network filter %d has an empty name", index)
		}
		config, err := l.makeConfig(filter)
		if err != nil {
			return nil, errors.Wrapf(err, "network filter %q", filter.Name)
		}
		result = append(result, model.NetworkFilter{
			Name:   filter.Name,
			Config: config,
		})
	}
	return result, nil
}

func (l *LdsManager) makeConfig(filter *xdsmodel.NetworkFilter) (map[string]any, error) {
	var m map[string]any
	switch cfg := filter.Config.(type) {
	case *xdsmodel.NetworkFilter_Yaml:
		if cfg.Yaml == nil {
			return nil, errors.New("YAML config is nil")
		}
		if err := yaml.Unmarshal([]byte(cfg.Yaml.Content), &m); err != nil {
			return nil, errors.Wrap(err, "decode YAML config")
		}
	case *xdsmodel.NetworkFilter_Json:
		if cfg.Json == nil {
			return nil, errors.New("JSON config is nil")
		}
		if err := json.Unmarshal([]byte(cfg.Json.Content), &m); err != nil {
			return nil, errors.Wrap(err, "decode JSON config")
		}
	case *xdsmodel.NetworkFilter_Struct:
		if cfg.Struct == nil {
			return nil, errors.New("Struct config is nil")
		}
		m = cfg.Struct.AsMap()
	default:
		return nil, errors.New("config is missing")
	}
	return m, nil
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
