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
	"fmt"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

type (
	// HttpFilterPlugin describe plugin
	HttpFilterPlugin interface {
		// Kind returns the unique kind name to represent itself.
		Kind() string

		// CreateFilter return the filter callback
		CreateFilter() (HttpFilter, error)
	}

	// HttpFilter describe http filter
	HttpFilter interface {
		// PrepareFilterChain add filter into chain
		PrepareFilterChain(ctx *http.HttpContext) error

		// Handle filter hook function
		Handle(ctx *http.HttpContext)

		Apply() error

		// Config get config for filter
		Config() interface{}
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
		CreateFilter(config interface{}, bs *model.Bootstrap) (NetworkFilter, error)
	}

	// NetworkFilter describe network filter
	NetworkFilter interface {
		// OnData handle the http context from worker
		OnData(hc *http.HttpContext) error
	}
)

var (
	httpFilterPluginRegistry    = map[string]HttpFilterPlugin{}
	networkFilterPluginRegistry = map[string]NetworkFilterPlugin{}
)

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

// Register registers network filter.
func RegisterNetworkFilter(f NetworkFilterPlugin) {
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
