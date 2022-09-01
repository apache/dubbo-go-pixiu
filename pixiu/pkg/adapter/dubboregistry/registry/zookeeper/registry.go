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
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/adapter/dubboregistry/common"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/adapter/dubboregistry/registry"
	baseRegistry "github.com/apache/dubbo-go-pixiu/pixiu/pkg/adapter/dubboregistry/registry/base"
	zk "github.com/apache/dubbo-go-pixiu/pixiu/pkg/adapter/dubboregistry/remoting/zookeeper"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
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
}

type ZKRegistry struct {
	*baseRegistry.BaseRegistry
	zkListeners map[registry.RegisteredType]registry.Listener
	client      *zk.ZooKeeperClient
}

var _ registry.Registry = new(ZKRegistry)

func newZKRegistry(regConfig model.Registry, adapterListener common.RegistryEventListener) (registry.Registry, error) {
	var zkReg = &ZKRegistry{}
	baseReg := baseRegistry.NewBaseRegistry(zkReg, adapterListener)
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
	initZKListeners(zkReg)
	return zkReg, nil
}

func initZKListeners(reg *ZKRegistry) {
	reg.zkListeners = make(map[registry.RegisteredType]registry.Listener)
	reg.zkListeners[registry.RegisteredTypeInterface] = newZKIntfListener(reg.client, reg, reg.AdapterListener)
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
	intfListener, ok := r.zkListeners[registry.RegisteredTypeInterface]
	if !ok {
		return errors.New("Listener for interface level registration does not initialized")
	}
	go intfListener.WatchAndHandle()
	return nil
}

// DoUnsubscribe stops monitoring the target registry.
func (r *ZKRegistry) DoUnsubscribe() error {
	intfListener, ok := r.zkListeners[registry.RegisteredTypeInterface]
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
