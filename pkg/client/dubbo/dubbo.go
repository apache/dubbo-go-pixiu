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
	"strconv"
	"strings"
	"sync"
	"time"
)

import (
	dclient "dubbo.apache.org/dubbo-go/v3/client"
	_ "dubbo.apache.org/dubbo-go/v3/cluster/loadbalance/consistenthashing"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/config/generic"
	"dubbo.apache.org/dubbo-go/v3/global"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	hessian "github.com/apache/dubbo-go-hessian2"

	"github.com/pkg/errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/client"
	cst "github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/config"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

const (
	JavaStringClassName = "java.lang.String"
	JavaLangClassName   = "java.lang.Long"
)

func javaClassNameElem(values []hessian.Object) []string {
	types := make([]string, len(values))
	for i, val := range values {
		if _, ok := val.(string); ok {
			types[i] = JavaStringClassName
			continue
		}
		types[i] = JavaLangClassName
	}
	return types
}

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
	defaultApplication = &global.ApplicationConfig{
		Organization: "dubbo-go-pixiu",
		Name:         "Dubbogo Pixiu",
		Module:       "dubbogo Pixiu",
		Owner:        "Dubbogo Pixiu",
		Environment:  "dev",
	}
)

// Client client to generic invoke dubbo
type Client struct {
	lock               sync.RWMutex
	GenericServicePool map[string]*generic.GenericService
	dubboProxyConfig   *DubboProxyConfig
	registries         map[string]*global.RegistryConfig
	dubboClient        *dclient.Client
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
	// Build registry configurations
	registries := make(map[string]*global.RegistryConfig)
	if dc.dubboProxyConfig != nil && dc.dubboProxyConfig.Registries != nil {
		for k, v := range dc.dubboProxyConfig.Registries {
			if len(v.Protocol) == 0 {
				logger.Warnf("can not find registry protocol config, use default type 'zookeeper'")
				v.Protocol = defaultDubboProtocol
			}
			registries[k] = &global.RegistryConfig{
				Protocol:  v.Protocol,
				Address:   v.Address,
				Timeout:   v.Timeout,
				Username:  v.Username,
				Password:  v.Password,
				Namespace: v.Namespace,
				Group:     v.Group,
			}
		}
	}
	dc.registries = registries

	// Create dubbo client with registries and application config
	var err error
	dc.dubboClient, err = dclient.NewClient(
		dclient.SetClientApplication(defaultApplication),
		dclient.SetClientRegistries(registries),
	)
	if err != nil {
		return err
	}

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
func (dc *Client) Call(req *client.Request) (res any, err error) {
	// if GET with no args, values would be nil
	values, err := dc.genericArgs(req)
	if err != nil {
		return nil, err
	}
	target, ok := values.(*dubboTarget)
	if !ok {
		return nil, errors.New("map parameters failed")
	}

	dm := req.API.IntegrationRequest
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
		if len(types) == 0 {
			types = javaClassNameElem(vals)
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
	ctx, span := tr.Start(req.Context, spanNameDubbogoClient)
	trace.SpanFromContext(req.Context).SpanContext()
	span.SetAttributes(attribute.Key(spanTagMethod).String(method))
	span.SetAttributes(attribute.Key(spanTagType).StringSlice(types))
	span.SetAttributes(attribute.Key(spanTagValues).String(string(finalValues)))
	defer span.End()

	// tracing inject manually;
	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, carrier)
	ctxWithAttachment := context.WithValue(ctx, constant.AttachmentKey, map[string]string(carrier))

	rst, err := gs.Invoke(ctxWithAttachment, method, types, vals)
	if err != nil {
		// TODO statusCode I don’t know what dubbo will return when it times out, so I will return it directly. I will judge it when I call it.
		span.RecordError(err)
		return nil, err
	}

	logger.Debugf("[dubbo-go-pixiu] dubbo client resp:%v", rst)

	return rst, nil
}

func (dc *Client) genericArgs(req *client.Request) (any, error) {
	values, err := dc.MapParams(req)
	if err != nil {
		return nil, err
	}

	return values, nil
}

// MapParams params mapping to api.
func (dc *Client) MapParams(req *client.Request) (any, error) {
	r := req.API.IntegrationRequest
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

func buildOption(conf config.MappingParam) client.RequestOption {
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
func (dc *Client) Get(ir config.IntegrationRequest) *generic.GenericService {
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

func (dc *Client) create(key string, irequest config.IntegrationRequest) *generic.GenericService {
	useNacosRegister := false
	registerIds := make([]string, 0)
	for k, v := range dc.registries {
		registerIds = append(registerIds, k)
		if v.Protocol == "nacos" {
			useNacosRegister = true
		}
	}

	dc.lock.Lock()
	defer dc.lock.Unlock()

	if service, ok := dc.GenericServicePool[key]; ok {
		return service
	}

	// Build ReferenceOptions using dubbo-go v3.3.1 client API
	opts := dc.buildReferenceOptions(irequest, registerIds)

	// Create ReferenceOptions with all required fields initialized.
	//
	// Why we manually initialize ReferenceOptions instead of using Client.Dial():
	// 1. Pixiu needs generic.GenericService for generic invocation (no IDL required)
	// 2. Client.Dial() returns Connection which is designed for typed clients
	// 3. Connection doesn't expose GenericService - we need direct access to ReferenceOptions
	// 4. This approach mirrors what Client.dial() does internally, following dubbo-go's design
	//
	// This initialization pattern matches dubbo-go's internal defaultReferenceOptions():
	// - All fields are initialized with default configs to prevent nil pointer dereference
	// - User options are applied afterwards to override defaults (functional options pattern)
	//
	// TODO: Refactor to use official dubbo-go API when generic invocation is supported
	// Currently, we manually construct ReferenceOptions because:
	// - Client.Dial() returns Connection (for typed clients only)
	// - Connection doesn't expose generic.GenericService
	// When dubbo-go adds official support for generic invocation via Client.Dial() or similar API,
	// we should migrate to that approach instead of manually initializing ReferenceOptions.
	// This will eliminate the need to track dubbo-go's internal structure changes.
	refOpts := &dclient.ReferenceOptions{
		Reference:   global.DefaultReferenceConfig(),
		Application: defaultApplication,
		Consumer:    global.DefaultConsumerConfig(),
		Shutdown:    global.DefaultShutdownConfig(),
		Metrics:     global.DefaultMetricsConfig(),
		Otel:        global.DefaultOtelConfig(),
		TLS:         global.DefaultTLSConfig(),
		Protocols:   make(map[string]*global.ProtocolConfig),
		Registries:  dc.registries,
	}

	// Apply user-provided options to override defaults
	for _, opt := range opts {
		opt(refOpts)
	}

	// Set generic mode for generic invocation
	refOpts.Reference.Generic = "true"

	// Log dubbo client configuration
	dc.logClientConfig(refOpts.Reference)

	// Call Refer to initialize the service reference
	refOpts.Refer()

	// sleep when first call to fetch enough service meta data from nacos
	// todo: Refer should guarantee it
	if useNacosRegister {
		time.Sleep(1000 * time.Millisecond)
	}

	clientService := refOpts.GetRPCService().(*generic.GenericService)
	dc.GenericServicePool[key] = clientService

	return clientService
}

// buildReferenceOptions builds a list of dubbo-go ReferenceOption using the official API
// Priority: DubboProxyConfig > IntegrationRequest > Default
func (dc *Client) buildReferenceOptions(irequest config.IntegrationRequest, registerIds []string) []dclient.ReferenceOption {
	opts := make([]dclient.ReferenceOption, 0, 16)

	// 1. Service Identity
	opts = append(opts, dclient.WithInterface(irequest.Interface))
	if irequest.Group != "" {
		opts = append(opts, dclient.WithGroup(irequest.Group))
	}
	if irequest.Version != "" {
		opts = append(opts, dclient.WithVersion(irequest.Version))
	}

	// 2. Registry Configuration
	if len(registerIds) > 0 {
		opts = append(opts, dclient.WithRegistryIDs(registerIds...))
	}

	// 3. Cluster Strategy (from DubboProxyConfig)
	cluster := dc.dubboProxyConfig.GetCluster()
	switch cluster {
	case "failover":
		opts = append(opts, dclient.WithClusterFailOver())
	case "failfast":
		opts = append(opts, dclient.WithClusterFailFast())
	case "failsafe":
		opts = append(opts, dclient.WithClusterFailSafe())
	case "failback":
		opts = append(opts, dclient.WithClusterFailBack())
	case "broadcast":
		opts = append(opts, dclient.WithClusterBroadcast())
	case "forking":
		opts = append(opts, dclient.WithClusterForking())
	case "available":
		opts = append(opts, dclient.WithClusterAvailable())
	case "zoneaware":
		opts = append(opts, dclient.WithClusterZoneAware())
	case "adaptiveservice":
		opts = append(opts, dclient.WithClusterAdaptiveService())
	default:
		opts = append(opts, dclient.WithCluster(cluster))
	}

	// 4. Protocol (from DubboProxyConfig)
	protocol := dc.dubboProxyConfig.GetProtocol()
	switch protocol {
	case "tri", "triple":
		opts = append(opts, dclient.WithProtocolTriple())
	case "dubbo":
		opts = append(opts, dclient.WithProtocolDubbo())
	case "jsonrpc":
		opts = append(opts, dclient.WithProtocolJsonRPC())
	default:
		opts = append(opts, dclient.WithProtocol(protocol))
	}

	// 5. Load Balance (from DubboProxyConfig)
	if dc.dubboProxyConfig.LoadBalance != "" {
		lb := dc.dubboProxyConfig.LoadBalance
		switch lb {
		case "random":
			opts = append(opts, dclient.WithLoadBalanceRandom())
		case "roundrobin":
			opts = append(opts, dclient.WithLoadBalanceRoundRobin())
		case "leastactive":
			opts = append(opts, dclient.WithLoadBalanceLeastActive())
		case "consistenthash", "consistenthashing":
			opts = append(opts, dclient.WithLoadBalanceConsistentHashing())
		case "p2c":
			opts = append(opts, dclient.WithLoadBalanceP2C())
		default:
			opts = append(opts, dclient.WithLoadBalance(lb))
		}
	}

	// 6. Retries (Priority: DubboProxyConfig > IntegrationRequest > Default)
	var retries int
	if dc.dubboProxyConfig.Retries != "" {
		if r, err := strconv.Atoi(dc.dubboProxyConfig.Retries); err == nil {
			retries = r
		} else {
			retries = 3
		}
	} else if irequest.Retries != "" {
		if r, err := strconv.Atoi(irequest.Retries); err == nil {
			retries = r
		} else {
			retries = 3
		}
	} else {
		retries = 3
	}
	opts = append(opts, dclient.WithRetries(retries))

	// 7. Request Timeout (from DubboProxyConfig.Timeout)
	var timeout time.Duration
	if dc.dubboProxyConfig.Timeout != nil {
		if t, err := time.ParseDuration(dc.dubboProxyConfig.Timeout.RequestTimeoutStr); err == nil {
			timeout = t
		} else {
			timeout = cst.DefaultReqTimeout
		}
	} else {
		timeout = cst.DefaultReqTimeout
	}
	opts = append(opts, dclient.WithRequestTimeout(timeout))

	// 8. Check (from DubboProxyConfig) - startup provider availability check
	if check := dc.dubboProxyConfig.GetCheck(); check != nil {
		if *check {
			opts = append(opts, dclient.WithCheck())
		}
		// If check=false, don't add WithCheck() - dubbo-go defaults to not checking
	}

	// 9. Filter (from DubboProxyConfig) - filter chain configuration
	if dc.dubboProxyConfig.Filter != "" {
		opts = append(opts, dclient.WithFilter(dc.dubboProxyConfig.Filter))
	}

	// 10. Serialization (from DubboProxyConfig) - serialization protocol
	if dc.dubboProxyConfig.Serialization != "" {
		switch dc.dubboProxyConfig.Serialization {
		case "json":
			opts = append(opts, dclient.WithSerializationJSON())
		default:
			opts = append(opts, dclient.WithSerialization(dc.dubboProxyConfig.Serialization))
		}
	}

	// 11. Sticky (from DubboProxyConfig) - sticky connection
	if dc.dubboProxyConfig.Sticky != nil && *dc.dubboProxyConfig.Sticky {
		opts = append(opts, dclient.WithSticky())
	}

	// 12. Params (from DubboProxyConfig) - custom parameters
	if len(dc.dubboProxyConfig.Params) > 0 {
		opts = append(opts, dclient.WithParams(dc.dubboProxyConfig.Params))
	}

	return opts
}

// logClientConfig logs the effective dubbo client configuration for debugging
func (dc *Client) logClientConfig(refConf *global.ReferenceConfig) {
	check := dc.dubboProxyConfig.GetCheck()
	checkStr := "nil (use dubbo-go default)"
	if check != nil {
		if *check {
			checkStr = "true"
		} else {
			checkStr = "false"
		}
	}

	sticky := dc.dubboProxyConfig.Sticky
	stickyStr := "false"
	if sticky != nil && *sticky {
		stickyStr = "true"
	}

	logger.Debugf("[dubbo-go-pixiu] Dubbo client config: interface=%s, group=%s, version=%s, cluster=%s, protocol=%s, loadbalance=%s, retries=%d, check=%s, timeout=%s, filter=%s, serialization=%s, sticky=%s, params=%v",
		refConf.InterfaceName,
		refConf.Group,
		refConf.Version,
		refConf.Cluster,
		refConf.Protocol,
		refConf.Loadbalance,
		refConf.Retries,
		checkStr,
		refConf.RequestTimeout,
		refConf.Filter,
		refConf.Serialization,
		stickyStr,
		refConf.Params,
	)
}
