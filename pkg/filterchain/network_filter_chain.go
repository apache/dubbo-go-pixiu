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

package filterchain

import (
	"context"
	"net/http"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/common/yaml"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

type NetworkFilterChain struct {
	filtersArray []filter.NetworkFilter
	config       model.FilterChain
}

// ServeHTTP handle http request
func (fc *NetworkFilterChain) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// todo: only one filter will exist for now, needs change when more than one
	for _, filter := range fc.filtersArray {
		filter.ServeHTTP(w, r)
	}
}

// OnDecode decode bytes received from getty listener
func (fc *NetworkFilterChain) OnDecode(data []byte) (any, int, error) {
	// todo: only one filter will exist for now, needs change when more than one
	for _, filter := range fc.filtersArray {
		return filter.OnDecode(data)
	}
	return nil, 0, errors.Errorf("filterChain don't have network filter")
}

// OnEncode encode struct to bytes sent to getty listener
func (fc *NetworkFilterChain) OnEncode(p any) ([]byte, error) {
	// todo: only one filter will exist for now, needs change when more than one
	for _, filter := range fc.filtersArray {
		return filter.OnEncode(p)
	}
	return nil, errors.Errorf("filterChain don't have network filter")
}

// OnData handle dubbo rpc invocation
func (fc *NetworkFilterChain) OnData(data any) (any, error) {
	// todo: only one filter will exist for now, needs change when more than one
	for _, filter := range fc.filtersArray {
		return filter.OnData(data)
	}
	return nil, errors.Errorf("filterChain don't have network filter")
}

// OnTripleData handle triple rpc invocation
func (fc *NetworkFilterChain) OnTripleData(ctx context.Context, methodName string, arguments []any) (any, error) {
	// todo: only one filter will exist for now, needs change when more than one
	for _, filter := range fc.filtersArray {
		return filter.OnTripleData(ctx, methodName, arguments)
	}
	return nil, errors.Errorf("filterChain don't have network filter")
}

// OnUnaryRPC handles a unary RPC call.
func (fc *NetworkFilterChain) OnUnaryRPC(ctx context.Context, fullMethod string, req any) (any, error) {
	// todo: only one filter will exist for now, needs change when more than one
	for _, filter := range fc.filtersArray {
		return filter.OnUnaryRPC(ctx, fullMethod, req)
	}
	return nil, errors.Errorf("filterChain don't have unary filter")
}

// OnStreamRPC handles a streaming RPC call.
func (fc *NetworkFilterChain) OnStreamRPC(stream model.RPCStream, info *model.RPCStreamInfo) error {
	// todo: only one filter will exist for now, needs change when more than one
	for _, filter := range fc.filtersArray {
		return filter.OnStreamRPC(stream, info)
	}
	return errors.Errorf("filterChain don't have gRPC stream filter")
}

// Close closes the filter chain and all filters in it.
func (fc *NetworkFilterChain) Close() error {
	var firstErr error
	for _, f := range fc.filtersArray {
		if err := closeNetworkFilter(f); err != nil {
			if firstErr == nil {
				firstErr = err
			}
			logger.Warnf("Failed to close filter: %v", err)
		}
	}
	return firstErr
}

func closeNetworkFilter(networkFilter filter.NetworkFilter) (err error) {
	if networkFilter == nil {
		return nil
	}
	defer func() {
		if recovered := recover(); recovered != nil {
			err = errors.Errorf("close network filter panicked: %v", recovered)
		}
	}()
	return networkFilter.Close()
}

// CreateNetworkFilterChain preserves the original best-effort public API used
// by out-of-tree listeners. xDS callers must use BuildNetworkFilterChain so a
// broken filter rejects the complete resource instead of being skipped.
func CreateNetworkFilterChain(config model.FilterChain) *NetworkFilterChain {
	var filters []filter.NetworkFilter
	for _, f := range config.Filters {
		p, err := filter.GetNetworkFilterPlugin(f.Name)
		if err != nil {
			logger.Errorf("CreateNetworkFilterChain %s getNetworkFilterPlugin error %s", f.Name, err)
			continue
		}
		filterConfig := p.Config()
		if err := yaml.ParseConfig(filterConfig, f.Config); err != nil {
			logger.Errorf("CreateNetworkFilterChain %s parse config error %s", f.Name, err)
			continue
		}
		networkFilter, err := p.CreateFilter(filterConfig)
		if err != nil {
			logger.Errorf("CreateNetworkFilterChain %s createFilter error %s", f.Name, err)
			continue
		}
		filters = append(filters, networkFilter)
	}
	return &NetworkFilterChain{filtersArray: filters, config: config}
}

// BuildNetworkFilterChain creates a complete network filter chain. A caller
// must not publish the chain unless every configured filter was constructed;
// silently omitting a broken filter changes the listener's security and
// routing semantics.
func BuildNetworkFilterChain(config model.FilterChain) (_ *NetworkFilterChain, err error) {
	var filters []filter.NetworkFilter
	defer func() {
		if err == nil {
			return
		}
		for _, networkFilter := range filters {
			if closeErr := closeNetworkFilter(networkFilter); closeErr != nil {
				logger.Warnf("Failed to close rejected filter: %v", closeErr)
			}
		}
	}()

	for index, f := range config.Filters {
		if f.Name == "" {
			return nil, errors.Errorf("network filter %d has an empty name", index)
		}
		p, err := filter.GetNetworkFilterPlugin(f.Name)
		if err != nil {
			return nil, errors.Wrapf(err, "get network filter plugin %q", f.Name)
		}

		filterConfig := p.Config()
		if err := yaml.ParseConfig(filterConfig, f.Config); err != nil {
			return nil, errors.Wrapf(err, "parse network filter %q config", f.Name)
		}

		networkFilter, err := p.CreateFilter(filterConfig)
		if err != nil {
			return nil, errors.Wrapf(err, "create network filter %q", f.Name)
		}
		if networkFilter == nil {
			return nil, errors.Errorf("network filter plugin %q returned nil", f.Name)
		}
		filters = append(filters, networkFilter)
	}

	return &NetworkFilterChain{
		filtersArray: filters,
		config:       config,
	}, nil
}
