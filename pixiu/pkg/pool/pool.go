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

package pool

import (
	"errors"
	"sync"
)

import (
	"github.com/dubbo-go-pixiu/pixiu-api/pkg/api/config"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/client/dubbo"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/client/http"
)

// ClientPool  a pool of client.
type ClientPool struct {
	poolMap map[config.RequestType]*sync.Pool
}

var (
	clientPool *ClientPool
	once       = sync.Once{}
)

func newClientPool() *ClientPool {
	clientPool := &ClientPool{
		poolMap: make(map[config.RequestType]*sync.Pool, 4),
	}
	// init default support request type
	clientPool.poolMap[config.DubboRequest] = &sync.Pool{
		New: func() interface{} {
			return dubbo.NewDubboClient()
		},
	}
	clientPool.poolMap[config.HTTPRequest] = &sync.Pool{
		New: func() interface{} {
			return http.NewHTTPClient()
		},
	}
	return clientPool
}

// SingletonPool singleton pool
func SingletonPool() *ClientPool {
	if clientPool == nil {
		once.Do(func() {
			clientPool = newClientPool()
		})
	}

	return clientPool
}

// GetClient  a factory method to get a client according to apiType .
func (pool *ClientPool) GetClient(t config.RequestType) (client.Client, error) {
	if pool.poolMap[t] != nil {
		return pool.poolMap[t].Get().(client.Client), nil
	}
	return nil, errors.New("protocol not supported yet")
}

// Put put client to pool.
func (pool *ClientPool) Put(t config.RequestType, c client.Client) error {
	if pool.poolMap[t] != nil {
		pool.poolMap[t].Put(c)
		return nil
	}

	return errors.New("protocol not supported yet")
}
