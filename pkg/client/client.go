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

package client

import (
	"errors"
	"github.com/dubbogo/dubbo-go-proxy/pkg/client/dubbo"
	"github.com/dubbogo/dubbo-go-proxy/pkg/client/httpclient"
	"github.com/dubbogo/dubbo-go-proxy/pkg/model"
	_ "github.com/gin-gonic/gin"
	"sync"
)

type Client interface {
	Init() error
	Close() error
	Call(req *Request) (resp Response, err error)
}

type ClientPool struct {
	poolMap map[model.ApiType]*sync.Pool
}

func (pool *ClientPool) getClient(t model.ApiType) (Client, error) {
	if pool.poolMap[t] != nil {
		return pool.poolMap[t].Get().(Client), nil
	}
	return nil, errors.New("protocol not supported yet")
}

func NewClientPool() *ClientPool {
	clientPool := &ClientPool{
		poolMap: make(map[model.ApiType]*sync.Pool),
	}
	clientPool.poolMap[model.DUBBO] = &sync.Pool{
		New: func() interface{} {
			return dubbo.NewDubboClient()
		},
	}
	clientPool.poolMap[model.REST] = &sync.Pool{
		New: func() interface{} {
			return httpclient.NewHttpClient()
		},
	}
	return clientPool
}
