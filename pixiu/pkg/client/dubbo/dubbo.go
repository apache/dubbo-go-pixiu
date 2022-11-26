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
	"encoding/json"
	"strings"
	"sync"
	"time"
)

import (
	_ "dubbo.apache.org/dubbo-go/v3/cluster/cluster/failover"
	_ "dubbo.apache.org/dubbo-go/v3/cluster/loadbalance/random"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	_ "dubbo.apache.org/dubbo-go/v3/common/proxy/proxy_factory"
	dg "dubbo.apache.org/dubbo-go/v3/config"
	"dubbo.apache.org/dubbo-go/v3/config/generic"
	_ "dubbo.apache.org/dubbo-go/v3/filter/generic"
	_ "dubbo.apache.org/dubbo-go/v3/filter/graceful_shutdown"
	_ "dubbo.apache.org/dubbo-go/v3/metadata/service/local"
	"dubbo.apache.org/dubbo-go/v3/protocol/dubbo"
	_ "dubbo.apache.org/dubbo-go/v3/registry/protocol"
	_ "dubbo.apache.org/dubbo-go/v3/registry/zookeeper"
	hessian "github.com/apache/dubbo-go-hessian2"
	fc "github.com/dubbo-go-pixiu/pixiu-api/pkg/api/config"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/client"
	cst "github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/config"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
)

// TODO java class name elem
const (
	JavaStringClassName = "java.lang.String"
	JavaLangClassName   = "java.lang.Long"
)

const (
	defaultDubboProtocol = "zookeeper"

	traceNameDubbogoClient = "dubbogo-client"
	spanNameDubbogoClient  = "DUBBOGO CLIENT"

	spanTagMethod = "method"
	spanTagType   = "type"
	spanTagValues = "values"
)

var (
	dubboClient        *Client
	onceClient         = sync.Once{}
	defaultApplication = &dg.ApplicationConfig{
		Organization: "dubbo-go-pixiu",
		Name:         "Dubbogo Pixiu",
		Module:       "dubbogo Pixiu",
		Version:      config.Version,
		Owner:        "Dubbogo Pixiu",
		Environment:  "dev",
	}
)

// Client client to generic invoke dubbo
type Client struct {
	lock               sync.RWMutex
	GenericServicePool map[string]*generic.GenericService
	dubboProxyConfig   *DubboProxyConfig
	rootConfig         *dg.RootConfig
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

// InitDefaultDubboClient init default dubbo client
func InitDefaultDubboClient(dpc *DubboProxyConfig) {
	dubboClient = NewDubboClient()
	dubboClient.SetConfig(dpc)
	if err := dubboClient.Apply(); err != nil {
		logger.Warnf("dubbo client apply error %s", err)
	}
}

// NewDubboClient create dubbo client
func NewDubboClient() *Client {
	return &Client{
		lock:               sync.RWMutex{},
		GenericServicePool: make(map[string]*generic.GenericService, 4),
	}
}

// SetConfig set config
func (dc *Client) SetConfig(dpc *DubboProxyConfig) {
	dc.dubboProxyConfig = dpc
}

// Apply init dubbo, config mapping can do here
func (dc *Client) Apply() error {

	rootConfigBuilder := dg.NewRootConfigBuilder()
	if dc.dubboProxyConfig != nil && dc.dubboProxyConfig.Registries != nil {
		for k, v := range dc.dubboProxyConfig.Registries {
			if len(v.Protocol) == 0 {
				logger.Warnf("can not find registry protocol config, use default type 'zookeeper'")
				v.Protocol = defaultDubboProtocol
			}
			rootConfigBuilder.AddRegistry(k, &dg.RegistryConfig{
				Protocol:  v.Protocol,
				Address:   v.Address,
				Timeout:   v.Timeout,
				Username:  v.Username,
				Password:  v.Password,
				Namespace: v.Namespace,
				Group:     v.Group,
			})
		}
	}
	rootConfigBuilder.SetApplication(defaultApplication)
	rootConfig := rootConfigBuilder.Build()

	if err := dg.Load(dg.WithRootConfig(rootConfig)); err != nil {
		panic(err)
	}
	dc.rootConfig = rootConfig
	return nil
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
func (dc *Client) Call(req *client.Request) (res interface{}, err error) {
	// if GET with no args, values would be nil
	values, err := dc.genericArgs(req)
	if err != nil {
		return nil, err
	}
	target, ok := values.(*dubboTarget)
	if !ok {
		return nil, errors.New("map parameters failed")
	}

	dm := req.API.Method.IntegrationRequest
	method := dm.Method
	types := []string{}
	vals := []hessian.Object{}
	finalValues := []byte{}

	if target != nil {
		logger.Debugf("[dubbo-go-pixiu] dubbo invoke, method:%s, types:%s, reqData:%v", method, target.Types, target.Values)
		types = target.Types
		vals = make([]hessian.Object, len(target.Values))
		for i, v := range target.Values {
			vals[i] = v
		}
		var err error
		finalValues, err = json.Marshal(vals)
		if err != nil {
			logger.Warnf("[dubbo-go-pixiu] reqData convert to string failed: %v", err)
		}
	} else {
		logger.Debugf("[dubbo-go-pixiu] dubbo invoke, method:%s, types:%s, reqData:%v", method, nil, nil)
	}

	gs := dc.Get(dm)
	tr := otel.Tracer(traceNameDubbogoClient)
	_, span := tr.Start(req.Context, spanNameDubbogoClient)
	trace.SpanFromContext(req.Context).SpanContext()
	span.SetAttributes(attribute.Key(spanTagMethod).String(method))
	span.SetAttributes(attribute.Key(spanTagType).StringSlice(types))
	span.SetAttributes(attribute.Key(spanTagValues).String(string(finalValues)))
	defer span.End()
	ctx := context.WithValue(req.Context, constant.TracingRemoteSpanCtx, trace.SpanFromContext(req.Context).SpanContext())
	rst, err := gs.Invoke(ctx, method, types, vals)
	if err != nil {
		return nil, err
	}

	logger.Debugf("[dubbo-go-pixiu] dubbo client resp:%v", rst)

	return rst, nil
}

func (dc *Client) genericArgs(req *client.Request) (interface{}, error) {
	values, err := dc.MapParams(req)
	if err != nil {
		return nil, err
	}

	return values, nil
}

// MapParams params mapping to api.
func (dc *Client) MapParams(req *client.Request) (interface{}, error) {
	r := req.API.Method.IntegrationRequest
	values := newDubboTarget(r.MappingParams)
	if dc.dubboProxyConfig != nil && dc.dubboProxyConfig.IsDefaultMap {
		values = newDubboTarget(defaultMappingParams)
	}
	for _, mappingParam := range r.MappingParams {
		source, _, err := client.ParseMapSource(mappingParam.Name)
		if err != nil {
			return nil, err
		}
		if mapper, ok := mappers[source]; ok {
			if err := mapper.Map(mappingParam, req, values, buildOption(mappingParam)); err != nil {
				return nil, err
			}
		}
	}
	return values, nil
}

func buildOption(conf fc.MappingParam) client.RequestOption {
	var opt client.RequestOption
	isGeneric, mapToType := getGenericMapTo(conf.MapTo)
	if isGeneric {
		opt = DefaultMapOption[mapToType]
	}
	return opt
}

func (dc *Client) get(key string) *generic.GenericService {
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
func (dc *Client) Get(ir fc.IntegrationRequest) *generic.GenericService {
	key := apiKey(&ir)
	if dc.check(key) {
		return dc.get(key)
	}

	return dc.create(key, ir)
}

func apiKey(ir *fc.IntegrationRequest) string {
	dbc := ir.DubboBackendConfig
	return strings.Join([]string{dbc.ClusterName, dbc.ApplicationName, dbc.Interface, dbc.Version, dbc.Group}, "_")
}

func (dc *Client) create(key string, irequest fc.IntegrationRequest) *generic.GenericService {
	useNacosRegister := false
	registerIds := make([]string, 0)
	for k, v := range dc.rootConfig.Registries {
		registerIds = append(registerIds, k)
		if v.Protocol == "nacos" {
			useNacosRegister = true
		}
	}

	refConf := dg.ReferenceConfig{
		InterfaceName: irequest.Interface,
		Cluster:       constant.ClusterKeyFailover,
		RegistryIDs:   registerIds,
		Protocol:      dubbo.DUBBO,
		Generic:       "true",
		Version:       irequest.DubboBackendConfig.Version,
		Group:         irequest.Group,
	}

	if len(irequest.DubboBackendConfig.Retries) == 0 {
		refConf.Retries = "3"
	} else {
		refConf.Retries = irequest.DubboBackendConfig.Retries
	}

	if dc.dubboProxyConfig.Timeout != nil {
		refConf.RequestTimeout = dc.dubboProxyConfig.Timeout.RequestTimeoutStr
	} else {
		refConf.RequestTimeout = cst.DefaultReqTimeout.String()
	}
	logger.Debugf("[dubbo-go-pixiu] client dubbo timeout val %v", refConf.RequestTimeout)
	dc.lock.Lock()
	defer dc.lock.Unlock()

	if service, ok := dc.GenericServicePool[key]; ok {
		return service
	}

	if err := dg.Load(dg.WithRootConfig(dc.rootConfig)); err != nil {
		panic(err)
	}

	_ = refConf.Init(dc.rootConfig)
	refConf.GenericLoad(key)

	// sleep when first call to fetch enough service meta data from nacos
	// todo: GenericLoad should guarantee it
	if useNacosRegister {
		time.Sleep(1000 * time.Millisecond)
	}

	clientService := refConf.GetRPCService().(*generic.GenericService)
	dc.GenericServicePool[key] = clientService

	return clientService
}
