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

package zookeeper

import (
	"strings"
	"time"
)

import (
	dubboCommon "dubbo.apache.org/dubbo-go/v3/common"

	hessian "github.com/apache/dubbo-go-hessian2"

	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/dubboregistry/common"
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/dubboregistry/registry"
	baseRegistry "github.com/apache/dubbo-go-pixiu/pkg/adapter/dubboregistry/registry/base"
	zk "github.com/apache/dubbo-go-pixiu/pkg/adapter/dubboregistry/remoting/zookeeper"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

var (
	_ baseRegistry.FacadeRegistry = new(ZKRegistry)
)

const (
	// RegistryZkClient zk client name
	RegistryZkClient = "zk registry"
	MaxFailTimes     = 2
	ConnDelay        = 3 * time.Second
	defaultTTL       = 10 * time.Minute
)

func init() {
	registry.SetRegistry(constant.Zookeeper, newZKRegistry)
	hessian.RegisterPOJO(&dubboCommon.MetadataInfo{})
	hessian.RegisterPOJO(&dubboCommon.ServiceInfo{})
	hessian.RegisterPOJO(&dubboCommon.URL{})
}

type ZKRegistry struct {
	*baseRegistry.BaseRegistry
	zkListeners map[registry.RegisteredType]registry.Listener
	client      *zk.ZooKeeperClient
}

var _ registry.Registry = new(ZKRegistry)

func newZKRegistry(regConfig model.Registry, adapterListener common.RegistryEventListener) (registry.Registry, error) {
	var zkReg = &ZKRegistry{}
	baseReg := baseRegistry.NewBaseRegistry(zkReg, adapterListener, registry.RegisterTypeFromName(regConfig.RegistryType))
	timeout, err := time.ParseDuration(regConfig.Timeout)
	if err != nil {
		return nil, errors.Errorf("Incorrect timeout configuration: %s", regConfig.Timeout)
	}
	client, eventChan, err := zk.NewZooKeeperClient(RegistryZkClient, strings.Split(regConfig.Address, ","), timeout)
	if err != nil {
		return nil, errors.Errorf("Initialize zookeeper client failed: %s", err.Error())
	}
	client.RegisterHandler(eventChan)
	zkReg.BaseRegistry = baseReg
	zkReg.client = client
	zkReg.zkListeners = make(map[registry.RegisteredType]registry.Listener)
	switch zkReg.RegisteredType {
	case registry.RegisteredTypeInterface:
		zkReg.zkListeners[zkReg.RegisteredType] = newZKIntfListener(zkReg.client, zkReg, zkReg.AdapterListener)
	case registry.RegisteredTypeApplication:
		zkReg.zkListeners[zkReg.RegisteredType] = newZkAppListener(zkReg.client, zkReg, zkReg.AdapterListener)
	default:
		return nil, errors.Errorf("Unsupported registry type: %s", regConfig.RegistryType)
	}
	return zkReg, nil
}

func (r *ZKRegistry) GetClient() *zk.ZooKeeperClient {
	return r.client
}

// DoSubscribe is the implementation of subscription on the target registry.
func (r *ZKRegistry) DoSubscribe() error {
	if err := r.interfaceSubscribe(); err != nil {
		return err
	}
	return nil
}

// To subscribe service level service discovery
func (r *ZKRegistry) interfaceSubscribe() error {
	intfListener, ok := r.zkListeners[r.RegisteredType]
	if !ok {
		return errors.New("Listener for interface level registration does not initialized")
	}
	go intfListener.WatchAndHandle()
	return nil
}

// DoUnsubscribe stops monitoring the target registry.
func (r *ZKRegistry) DoUnsubscribe() error {
	intfListener, ok := r.zkListeners[r.RegisteredType]
	if !ok {
		return errors.New("Listener for interface level registration does not initialized")
	}
	intfListener.Close()
	for k, l := range r.GetAllSvcListener() {
		l.Close()
		r.RemoveSvcListener(k)
	}
	return nil
}
