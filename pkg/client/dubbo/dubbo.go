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
	"strings"
	"sync"
	"time"

	"go.opentelemetry.io/otel/trace"

	"github.com/apache/dubbo-go/common/constant"
	dg "github.com/apache/dubbo-go/config"
	"github.com/apache/dubbo-go/protocol/dubbo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"

	"github.com/apache/dubbo-go-pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pkg/config"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	fc "github.com/dubbogo/dubbo-go-pixiu-filter/pkg/api/config"
	"github.com/pkg/errors"
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
	dgCfg              dg.ConsumerConfig
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
	staticResources := config.GetBootstrap().StaticResources
	cls := staticResources.Clusters
	tc := staticResources.TimeoutConfig

	// dubbogo consumer config
	dgCfg = dg.ConsumerConfig{
		Check:      new(bool),
		Registries: make(map[string]*dg.RegistryConfig, 4),
	}
	// timeout config
	dgCfg.Connect_Timeout = tc.ConnectTimeoutStr
	dgCfg.Request_Timeout = tc.RequestTimeoutStr
	dgCfg.ApplicationConfig = defaultApplication
	for i := range cls {
		c := cls[i]
		for k, v := range c.Registries {
			if len(v.Protocol) == 0 {
				logger.Warnf("can not find registry protocol config, use default type 'zookeeper'")
				v.Protocol = defaultDubboProtocol
			}
			dgCfg.Registries[k] = &dg.RegistryConfig{
				Protocol:   v.Protocol,
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
func (dc *Client) Call(req *client.Request) (res interface{}, err error) {
	values, err := dc.genericArgs(req)
	if err != nil {
		return nil, err
	}
	val, ok := values.(*dubboTarget)
	if !ok {
		return nil, errors.New("map parameters failed")
	}

	dm := req.API.Method.IntegrationRequest
	method := dm.Method

	logger.Debugf("[dubbo-go-pixiu] dubbo invoke, method:%s, types:%s, reqData:%v", method, val.Types, val.Values)

	gs := dc.Get(dm)
	tr := otel.Tracer(traceNameDubbogoClient)
	_, span := tr.Start(req.Context, spanNameDubbogoClient)
	trace.SpanFromContext(req.Context).SpanContext()
	span.SetAttributes(attribute.Key(spanTagMethod).String(method))
	span.SetAttributes(attribute.Key(spanTagType).Array(val.Types))
	span.SetAttributes(attribute.Key(spanTagValues).Array(val.Values))
	defer span.End()
	ctx := context.WithValue(req.Context, constant.TRACING_REMOTE_SPAN_CTX, trace.SpanFromContext(req.Context).SpanContext())
	rst, err := gs.Invoke(ctx, []interface{}{method, val.Types, val.Values})
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
func (dc *Client) Get(ir fc.IntegrationRequest) *dg.GenericService {
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

func (dc *Client) create(key string, irequest fc.IntegrationRequest) *dg.GenericService {
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
	//TODO: fix it later
	// sleep to wait invoker create
	time.Sleep(500 * time.Millisecond)
	clientService := referenceConfig.GetRPCService().(*dg.GenericService)

	dc.GenericServicePool[key] = clientService
	return clientService
}
