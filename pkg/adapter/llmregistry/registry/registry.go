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

package registry

import (
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/llmregistry/common"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// A map to store registry creation functions by protocol name.
var registryMap = make(map[string]func(model.Registry, common.RegistryEventListener) (Registry, error), 8)

// Registry interface defines the basic features of a service registry.
type Registry interface {
	// Subscribe starts monitoring the target registry for service changes.
	Subscribe() error
	// Unsubscribe stops monitoring the target registry.
	Unsubscribe() error
}

// SetRegistry registers a factory function for creating a new registry client.
func SetRegistry(name string, newRegFunc func(model.Registry, common.RegistryEventListener) (Registry, error)) {
	registryMap[name] = newRegFunc
}

// GetRegistry creates and returns a new registry client based on the configuration.
// It panics if the registry client cannot be initialized.
func GetRegistry(regConfig model.Registry, listener common.RegistryEventListener) (Registry, error) {
	if newRegFunc, ok := registryMap[regConfig.Protocol]; ok {
		reg, err := newRegFunc(regConfig, listener)
		if err != nil {
			return nil, errors.New("Initialize Registry " + regConfig.Protocol + " failed due to: " + err.Error())
		}
		return reg, nil
	}
	return nil, errors.New("Registry protocol " + regConfig.Protocol + " is not supported")
}
