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

package adapter

import (
	"fmt"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/go-errors/errors"
)

type (
	// AdapterPlugin plugin for adapter
	AdapterPlugin interface {
		// Kind returns the unique kind name to represent itself.
		Kind() string

		// CreateAdapter return the Adapter callback
		CreateAdapter(config interface{}, bs *model.Bootstrap) (Adapter, error)
	}

	// Adapter adapter interface
	Adapter interface {
		Start()
		Stop()
	}
)

var (
	adapterPlugins = map[string]AdapterPlugin{}
)

// Register registers adapter plugin
func RegisterAdapterPlugin(p AdapterPlugin) {
	if p.Kind() == "" {
		panic(fmt.Errorf("%T: empty kind", p))
	}

	existedPlugin, existed := adapterPlugins[p.Kind()]
	if existed {
		panic(fmt.Errorf("%T and %T got same kind: %s", p, existedPlugin, p.Kind()))
	}

	adapterPlugins[p.Kind()] = p
}

// GetAdapterPlugin get plugin by kind
func GetAdapterPlugin(kind string) (AdapterPlugin, error) {
	existedAdapter, existed := adapterPlugins[kind]
	if existed {
		return existedAdapter, nil
	} else {
		return nil, errors.Errorf("plugin not found %s", kind)
	}
}
