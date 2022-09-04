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

package filter

import (
	"context"
	"fmt"
	stdHttp "net/http"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/context/dubbo"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/context/http"
)

type (
	// HttpFilterPlugin describe plugin
	HttpFilterPlugin interface {
		// Kind returns the unique kind name to represent itself.
		Kind() string

		// CreateFilterFactory return the filter factory
		CreateFilterFactory() (HttpFilterFactory, error)
	}

	// HttpFilterFactory describe http filter
	HttpFilterFactory interface {
		// Config Expose the config so that Filter Manger can inject it, so it must be a pointer
		Config() interface{}

		// Apply After the config is injected, check it or make it to default
		Apply() error

		// PrepareFilterChain create filter and append it to FilterChain
		//
		// Be Careful !!! Do not pass the Factory's config pointer to the Filter instance,
		// Factory's config may be updated by FilterManager
		PrepareFilterChain(ctx *http.HttpContext, chain FilterChain) error
	}

	// HttpDecodeFilter before invoke upstream, like add/remove Header, route mutation etc..
	//
	// if config like this:
	// - A
	// - B
	// - C
	// decode filters will be invoked in the config order: A、B、C, and decode filters will be
	// invoked in the reverse order: C、B、A
	HttpDecodeFilter interface {
		Decode(ctx *http.HttpContext) FilterStatus
	}

	// HttpEncodeFilter after invoke upstream,
	// decode filters will be invoked in the reverse order
	HttpEncodeFilter interface {
		Encode(ctx *http.HttpContext) FilterStatus
	}

	// NetworkFilter describe network filter plugin
	NetworkFilterPlugin interface {
		// Kind returns the unique kind name to represent itself.
		Kind() string
		// CreateFilterFactory return the filter callback
		CreateFilter(config interface{}) (NetworkFilter, error)
		// Config Expose the config so that Filter Manger can inject it, so it must be a pointer
		Config() interface{}
	}

	// NetworkFilter describe network filter
	NetworkFilter interface {
		// ServeHTTP handle http request and response
		ServeHTTP(w stdHttp.ResponseWriter, r *stdHttp.Request)
		// OnDecode decode bytes received from getty server
		OnDecode(data []byte) (interface{}, int, error)
		// OnEncode encode interface sent to getty server
		OnEncode(p interface{}) ([]byte, error)
		// OnData dubbo rpc invocation from getty server
		OnData(data interface{}) (interface{}, error)
		// OnTripleData triple rpc invocation from triple-server
		OnTripleData(ctx context.Context, methodName string, arguments []interface{}) (interface{}, error)
	}

	// EmptyNetworkFilter default empty network filter adapter which offers empty function implements
	EmptyNetworkFilter struct{}

	// DubboFilter describe dubbo filter
	DubboFilter interface {
		// Handle rpc invocation
		Handle(ctx *dubbo.RpcContext) FilterStatus
	}

	// NetworkFilter describe dubbo filter plugin
	DubboFilterPlugin interface {
		// Kind returns the unique kind name to represent itself.
		Kind() string
		// CreateFilter return the filter callback
		CreateFilter(config interface{}) (DubboFilter, error)
		// Config Expose the config so that Filter Manger can inject it, so it must be a pointer
		Config() interface{}
	}
)

var (
	httpFilterPluginRegistry    = map[string]HttpFilterPlugin{}
	networkFilterPluginRegistry = map[string]NetworkFilterPlugin{}
	dubboFilterPluginRegistry   = map[string]DubboFilterPlugin{}
)

// OnDecode empty implement
func (enf *EmptyNetworkFilter) OnDecode(data []byte) (interface{}, int, error) {
	panic("OnDecode is not implemented")
}

// OnEncode empty implement
func (enf *EmptyNetworkFilter) OnEncode(p interface{}) ([]byte, error) {
	panic("OnEncode is not implemented")
}

// OnData receive data from listener
func (enf *EmptyNetworkFilter) OnData(data interface{}) (interface{}, error) {
	panic("OnData is not implemented")
}

// OnTripleData empty implement
func (enf *EmptyNetworkFilter) OnTripleData(ctx context.Context, methodName string, arguments []interface{}) (interface{}, error) {
	panic("OnTripleData is not implemented")
}

// ServeHTTP empty implement
func (enf *EmptyNetworkFilter) ServeHTTP(w stdHttp.ResponseWriter, r *stdHttp.Request) {
	panic("ServeHTTP is not implemented")
}

// Register registers filter plugin.
func RegisterHttpFilter(f HttpFilterPlugin) {
	if f.Kind() == "" {
		panic(fmt.Errorf("%T: empty kind", f))
	}

	existedPlugin, existed := httpFilterPluginRegistry[f.Kind()]
	if existed {
		panic(fmt.Errorf("%T and %T got same kind: %s", f, existedPlugin, f.Kind()))
	}

	httpFilterPluginRegistry[f.Kind()] = f
}

// GetHttpFilterPlugin get plugin by kind
func GetHttpFilterPlugin(kind string) (HttpFilterPlugin, error) {
	existedFilter, existed := httpFilterPluginRegistry[kind]
	if existed {
		return existedFilter, nil
	}
	return nil, errors.Errorf("plugin not found %s", kind)
}

// RegisterNetworkFilter registers network filter.
func RegisterNetworkFilterPlugin(f NetworkFilterPlugin) {
	if f.Kind() == "" {
		panic(fmt.Errorf("%T: empty kind", f))
	}

	existedFilter, existed := networkFilterPluginRegistry[f.Kind()]
	if existed {
		panic(fmt.Errorf("%T and %T got same kind: %s", f, existedFilter, f.Kind()))
	}

	networkFilterPluginRegistry[f.Kind()] = f
}

// GetNetworkFilterPlugin get plugin by kind
func GetNetworkFilterPlugin(kind string) (NetworkFilterPlugin, error) {
	existedFilter, existed := networkFilterPluginRegistry[kind]
	if existed {
		return existedFilter, nil
	}
	return nil, errors.Errorf("plugin not found %s", kind)
}

// RegisterDubboFilterPlugin registers dubbo filter.
func RegisterDubboFilterPlugin(f DubboFilterPlugin) {
	if f.Kind() == "" {
		panic(fmt.Errorf("%T: empty kind", f))
	}

	existedFilter, existed := dubboFilterPluginRegistry[f.Kind()]
	if existed {
		panic(fmt.Errorf("%T and %T got same kind: %s", f, existedFilter, f.Kind()))
	}

	dubboFilterPluginRegistry[f.Kind()] = f
}

// GetDubboFilterPlugin get plugin by kind
func GetDubboFilterPlugin(kind string) (DubboFilterPlugin, error) {
	existedFilter, existed := dubboFilterPluginRegistry[kind]
	if existed {
		return existedFilter, nil
	}
	return nil, errors.Errorf("plugin not found %s", kind)
}
