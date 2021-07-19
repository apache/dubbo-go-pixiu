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

package extension

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/registry"
)

var registryMap = make(map[string]func(model.Registry) (registry.Registry, error), 8)

// SetFilterFunc will store the @filter and @name
func SetRegistry(name string, newRegFunc func(model.Registry) (registry.Registry, error)) {
	registryMap[name] = newRegFunc
}

// GetMustFilterFunc will return the pixiu.FilterFunc
// if not found, it will panic
func GetRegistry(name string, regConfig model.Registry) registry.Registry {
	if registry, ok := registryMap[name]; ok {
		reg, err := registry(regConfig)
		if err != nil {
			panic("Initialize Registry" + name + "failed due to: " + err.Error())
		}
		return reg
	}
	panic("Registry " + name + " does not support yet")
}
