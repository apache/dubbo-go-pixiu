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

package dubbo

import (
	"context"
	"reflect"
	"strings"
	"sync"
	"time"
)

import (
	"github.com/apache/dubbo-go/common/constant"
	dg "github.com/apache/dubbo-go/config"
	"github.com/apache/dubbo-go/protocol/dubbo"
	"github.com/pkg/errors"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/client"
	"github.com/dubbogo/dubbo-go-proxy/pkg/config"
	"github.com/dubbogo/dubbo-go-proxy/pkg/logger"
)

// TODO java class name elem
const (
	JavaStringClassName = "java.lang.String"
	JavaLangClassName   = "java.lang.Long"
)

var mappers = map[string]client.ParamMapper{
	"queryStrings": queryStringsMapper{},
	"headers":      headerMapper{},
	"requestBody":  bodyMapper{},
}

var (
	dubboClient        *Client
	onceClient         = sync.Once{}
	dgCfg              dg.ConsumerConfig
	defaultApplication = &dg.ApplicationConfig{
		Organization: "dubbo-go-proxy",
		Name:         "Dubbogo Proxy",
		Module:       "dubbogo proxy",
		Version:      "0.1.0",
		Owner:        "Dubbogo Proxy",
		Environment:  "dev",
	}
)

// Client client to generic invoke dubbo
type Client struct {
	lock               sync.RWMutex
	GenericServicePool map[string]*dg.GenericService
}

// SingletonDubboClient singleton dubbo clent
func SingletonDubboClient() *Client {
	if dubboClient == nil {
		onceClient.Do(func() {
			dubboClient = NewDubboClient()
		})
	}

	return dubboClient
}

// NewDubboClient create dubbo client
func NewDubboClient() *Client {
	return &Client{
		lock:               sync.RWMutex{},
		GenericServicePool: make(map[string]*dg.GenericService, 4),
	}
}

// Init init dubbo, config mapping can do here
func (dc *Client) Init() error {
	dc.GenericServicePool = make(map[string]*dg.GenericService, 4)

	cls := config.GetBootstrap().StaticResources.Clusters

	// dubbogo comsumer config
	dgCfg = dg.ConsumerConfig{
		Check:      new(bool),
		Registries: make(map[string]*dg.RegistryConfig, 4),
	}
	dgCfg.ApplicationConfig = defaultApplication
	for i := range cls {
		c := cls[i]
		dgCfg.Request_Timeout = c.RequestTimeoutStr
		dgCfg.Connect_Timeout = c.ConnectTimeoutStr
		for k, v := range c.Registries {
			dgCfg.Registries[k] = &dg.RegistryConfig{
				Protocol:   k,
				Address:    v.Address,
				TimeoutStr: v.Timeout,
				Username:   v.Username,
				Password:   v.Password,
			}
		}
	}

	initDubbogo()

	return nil
}

func initDubbogo() {
	dg.SetConsumerConfig(dgCfg)
	dubbo.SetClientConf(dubbo.GetDefaultClientConfig())
	dg.Load()
}

// Close clear GenericServicePool.
func (dc *Client) Close() error {
	dc.lock.Lock()
	defer dc.lock.Unlock()
	for k := range dc.GenericServicePool {
		delete(dc.GenericServicePool, k)
	}
	return nil
}

// Call invoke service
func (dc *Client) Call(req *client.Request) (resp *client.Response, err error) {
	dm := req.API.Method.IntegrationRequest
	types, values, err := dc.MappingParams(req)
	if err != nil {
		return *client.EmptyResponse, err
	}

	method := dm.Method
	logger.Debugf("[dubbo-go-proxy] invoke, method:%s, types:%s, reqData:%v", method, types, values)

	gs := dc.Get(dm)

	rst, err := gs.Invoke(req.Context, []interface{}{method, types, values})

	if err != nil {
		return &client.Response{}, err
	}

	logger.Debugf("[dubbo-go-proxy] dubbo client resp:%v", rst)

	if rst == nil {
		return &client.Response{}, nil
	}

	return &client.Response{Data: rst}, nil
}

// MappingParams param mapping to api.
func (dc *Client) MappingParams(req *client.Request) ([]string, []interface{}, error) {
	r := req.API.Method.IntegrationRequest
	var values []interface{}
	for _, mappingParam := range r.MappingParams {
		source, _, err := client.ParseMapSource(mappingParam.Name)
		if err != nil {
			return nil, nil, err
		}
		if mapper, ok := mappers[source]; ok {
			if err := mapper.Map(mappingParam, *req, &values); err != nil {
				return nil, nil, err
			}
		}
	}
	return req.API.IntegrationRequest.ParamTypes, values, nil
}

func (dc *Client) get(key string) *dg.GenericService {
	dc.lock.RLock()
	defer dc.lock.RUnlock()
	return dc.GenericServicePool[key]
}

func (dc *Client) check(key string) bool {
	dc.lock.RLock()
	defer dc.lock.RUnlock()
	if _, ok := dc.GenericServicePool[key]; ok {
		return true
	}
	return false
}

// Get find a dubbo GenericService
func (dc *Client) Get(ir config.IntegrationRequest) *dg.GenericService {
	key := apiKey(&ir)
	if dc.check(key) {
		return dc.get(key)
	}

	return dc.create(key, ir)
}

func apiKey(ir *config.IntegrationRequest) string {
	dbc := ir.DubboBackendConfig
	return strings.Join([]string{dbc.ClusterName, dbc.ApplicationName, dbc.Interface, dbc.Version, dbc.Group}, "_")
}

func (dc *Client) create(key string, irequest config.IntegrationRequest) *dg.GenericService {
	referenceConfig := dg.NewReferenceConfig(irequest.Interface, context.TODO())
	referenceConfig.InterfaceName = irequest.Interface
	referenceConfig.Cluster = constant.DEFAULT_CLUSTER
	var registers []string
	for k := range dgCfg.Registries {
		registers = append(registers, k)
	}
	referenceConfig.Registry = strings.Join(registers, ",")

	if len(irequest.DubboBackendConfig.Protocol) == 0 {
		referenceConfig.Protocol = dubbo.DUBBO
	} else {
		referenceConfig.Protocol = irequest.DubboBackendConfig.Protocol
	}

	referenceConfig.Version = irequest.DubboBackendConfig.Version
	referenceConfig.Group = irequest.Group
	referenceConfig.Generic = true
	if len(irequest.DubboBackendConfig.Retries) == 0 {
		referenceConfig.Retries = "3"
	} else {
		referenceConfig.Retries = irequest.DubboBackendConfig.Retries
	}
	dc.lock.Lock()
	defer dc.lock.Unlock()
	referenceConfig.GenericLoad(key)
	time.Sleep(200 * time.Millisecond) //sleep to wait invoker create
	clientService := referenceConfig.GetRPCService().(*dg.GenericService)

	dc.GenericServicePool[key] = clientService
	return clientService
}

func validateTarget(target interface{}) (reflect.Value, error) {
	rv := reflect.ValueOf(target)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return rv, errors.New("Target params must be a non-nil pointer")
	}
	if _, ok := target.(*[]interface{}); !ok {
		return rv, errors.New("Target params for dubbo backend must be *[]interface{}")
	}
	return rv, nil
}

func setTarget(rv reflect.Value, pos int, value interface{}) {
	if rv.Kind() != reflect.Ptr && rv.Type().Name() != "" && rv.CanAddr() {
		rv = rv.Addr()
	} else {
		rv = rv.Elem()
	}

	tempValue := rv.Interface().([]interface{})
	if len(tempValue) <= pos {
		list := make([]interface{}, pos+1-len(tempValue))
		tempValue = append(tempValue, list...)
	}
	tempValue[pos] = value
	rv.Set(reflect.ValueOf(tempValue))
}
