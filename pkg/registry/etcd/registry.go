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

package etcd

import (
	"github.com/apache/dubbo-go-pixiu/pkg/registry"
	"github.com/apache/dubbo-go/common"
	"github.com/apache/dubbo-go/remoting/etcdv3"
)

const (
	// RegistryETCDV3Client etcd client name
	RegistryETCDV3Client = "etcd registry"
	rootPath             = "/dubbo"
)

func init() {
	// extension.SetRegistry(constant.Zookeeper, newZKRegistry)
}

type etcdV3Registry struct {
	*registry.BaseRegistry
	client *etcdv3.Client
}

func (r *etcdV3Registry) GetClient() *etcdv3.Client {
	return r.client
}

func (r etcdV3Registry) LoadServices() {
	panic("implement me")
}

func (r etcdV3Registry) Subscribe(url common.URL) error {
	panic("implement me")
}

func (r etcdV3Registry) Unsubscribe(url common.URL) error {
	panic("implement me")
}

//func newETCDV3Registry(regConfig model.Registry) (registry.Registry, error) {
//	baseReg := registry.NewBaseRegistry()
//	timeout, err := time.ParseDuration(regConfig.Timeout)
//	if err != nil {
//		return nil, errors.Errorf("Incorrect timeout configuration: %s", regConfig.Timeout)
//	}
//
//	defaultHeartbeat := 1
//	client, err := etcdv3.NewClient(RegistryETCDV3Client, strings.Split(regConfig.Address, ","), timeout, defaultHeartbeat)
//	if err != nil {
//		return nil, errors.Errorf("Initialize zookeeper client failed: %s", err.Error())
//	}
//
//	//client.RegisterHandler(eventChan)
//
//	return &etcdV3Registry{
//		BaseRegistry: baseReg,
//		client:       client,
//	}, nil
//}
