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
	"sync"
	"time"
)

import (
	gxzookeeper "github.com/dubbogo/gost/database/kv/zk"
	perrors "github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
)

func ValidateZookeeperClient(container ZkClientFacade, zkName string) error {
	lock := container.ZkClientLock()
	config := container.GetConfig()

	lock.Lock()
	defer lock.Unlock()

	if container.ZkClient() == nil {
		newClient, cltErr := NewZkClient(zkName, config)
		if cltErr != nil {
			return perrors.WithMessagef(cltErr, "create zookeeper client (address:%+v)", config.Address)
		}
		container.SetZkClient(newClient)
	}
	return nil
}

func NewZkClient(zkName string, config *model.RemoteConfig) (*gxzookeeper.ZookeeperClient, error) {
	var (
		zkAddresses = strings.Split(config.Address, ",")
	)
	timeout, _ := time.ParseDuration(config.Timeout)

	newClient, cltErr := gxzookeeper.NewZookeeperClient(zkName, zkAddresses, true, gxzookeeper.WithZkTimeOut(timeout))
	if cltErr != nil {
		logger.Warnf("newZookeeperClient(name{%s}, zk address{%v}, timeout{%d}) = error{%v}",
			zkName, config.Address, timeout.String(), cltErr)
		return nil, perrors.WithMessagef(cltErr, "newZookeeperClient(address:%+v)", config.Address)
	}
	return newClient, nil
}

type Listener interface {
	// Close closes this listener
	Close()
	// WatchAndHandle watch the target path and handle the event
	WatchAndHandle()
}

type ZkClientFacade interface {
	Name() string
	ZkClient() *gxzookeeper.ZookeeperClient
	SetZkClient(*gxzookeeper.ZookeeperClient)
	ZkClientLock() *sync.Mutex
	WaitGroup() *sync.WaitGroup // for wait group control, zk client listener & zk client container
	Done() chan struct{}        // for registry destroy
	RestartCallBack() bool
	GetConfig() *model.RemoteConfig
}

// HandleClientRestart keeps the connection between client and server
// This method should be used only once. You can use handleClientRestart() in package registry.
func HandleClientRestart(r ZkClientFacade) {
	defer r.WaitGroup().Done()
	for {
		select {
		case <-r.ZkClient().Reconnect():
			r.RestartCallBack()
			time.Sleep(10 * time.Microsecond)
		case <-r.Done():
			logger.Warnf("receive registry destroy event, quit client restart handler")
			return
		}
	}
}
