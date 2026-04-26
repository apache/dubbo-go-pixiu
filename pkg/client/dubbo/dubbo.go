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
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

import (
	dclient "dubbo.apache.org/dubbo-go/v3/client"
	_ "dubbo.apache.org/dubbo-go/v3/cluster/loadbalance/consistenthashing"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/filter/generic"
	"dubbo.apache.org/dubbo-go/v3/global"
	_ "dubbo.apache.org/dubbo-go/v3/imports"

	hessian "github.com/apache/dubbo-go-hessian2"

	"github.com/pkg/errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
)

import (
	cst "github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
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

type resolvedConsumerDefaults struct {
	Cluster        string
	LoadBalance    string
	Retries        string
	RequestTimeout time.Duration
}

type resolvedReferSpec struct {
	Mode                   string
	Interface              string
	Group                  string
	Version                string
	URL                    string
	RegistryIDs            []string
	EffectiveProtocol      string
	EffectiveSerialization string
	UseNacosWarmup         bool
	ConsumerDefaults       resolvedConsumerDefaults
}

type genericServiceKey struct {
	Mode              string   `json:"mode"`
	URL               string   `json:"url"`
	RegistryIDs       []string `json:"registry_ids"`
	Cluster           string   `json:"cluster"`
	LoadBalance       string   `json:"load_balance"`
	Retries           string   `json:"retries"`
	RequestTimeout    string   `json:"request_timeout"`
	Interface         string   `json:"interface"`
	Version           string   `json:"version"`
	Group             string   `json:"group"`
	EffectiveProtocol string   `json:"effective_protocol"`
	Serialization     string   `json:"serialization"`
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
				Protocol:     v.Protocol,
				Address:      v.Address,
				Timeout:      v.Timeout,
				Username:     v.Username,
				Password:     v.Password,
				Namespace:    v.Namespace,
				Group:        v.Group,
				RegistryType: v.RegistryType,
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

// Call invoke service.
func (dc *Client) Call(ctx context.Context, req *DubboOutboundRequest) (any, error) {
	if req == nil {
		return nil, errors.New("dubbo outbound request is nil")
	}

	spec := dc.resolveFromOutbound(req)
	types, vals, finalValues, err := dc.preparePayload(req)
	if err != nil {
		return nil, err
	}

	gs, err := dc.Get(spec)
	if err != nil {
		return nil, err
	}
	if gs == nil {
		return nil, errors.New("dubbo generic service is nil")
	}

	invokeCtx, cancel := prepareInvokeContext(ctx, req.Timeout)
	if cancel != nil {
		defer cancel()
	}

	spanCtx, span := otel.Tracer(traceNameDubbogoClient).Start(invokeCtx, spanNameDubbogoClient)
	defer span.End()
	span.SetAttributes(
		attribute.String(spanTagMethod, req.Method),
		attribute.StringSlice(spanTagType, types),
		attribute.String(spanTagValues, string(finalValues)),
	)

	spanCtx = context.WithValue(spanCtx, constant.AttachmentKey, mergeOutboundAttachments(spanCtx, req.Attachments))
	ctxWithAttachment := withAttachments(spanCtx)
	rst, err := gs.Invoke(ctxWithAttachment, req.Method, types, vals)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	logger.Debugf("[dubbo-go-pixiu] dubbo invoke result:%v", rst)
	return rst, nil
}

func (dc *Client) resolveFromOutbound(req *DubboOutboundRequest) resolvedReferSpec {
	spec := resolvedReferSpec{
		Interface:              req.Service,
		Group:                  req.Group,
		Version:                req.Version,
		EffectiveProtocol:      req.Protocol,
		EffectiveSerialization: req.Serialization,
		ConsumerDefaults:       dc.resolveGlobalConsumerDefaults(),
	}

	if strings.TrimSpace(req.Address) != "" {
		spec.Mode = "direct"
		spec.URL = req.Protocol + "://" + req.Address
		return spec
	}

	registryIDs := make([]string, 0, len(dc.registries))
	useNacosWarmup := false
	for id, registry := range dc.registries {
		registryIDs = append(registryIDs, id)
		if registry != nil && registry.Protocol == "nacos" {
			useNacosWarmup = true
		}
	}
	sort.Strings(registryIDs)

	spec.Mode = "registry"
	spec.RegistryIDs = registryIDs
	spec.UseNacosWarmup = useNacosWarmup
	return spec
}

func (dc *Client) resolveGlobalConsumerDefaults() resolvedConsumerDefaults {
	defaults := resolvedConsumerDefaults{
		Cluster:        "failover",
		Retries:        "3",
		RequestTimeout: cst.DefaultReqTimeout,
	}

	if dc.dubboProxyConfig == nil {
		return defaults
	}

	defaults.LoadBalance = dc.dubboProxyConfig.LoadBalance
	if strings.TrimSpace(dc.dubboProxyConfig.Retries) != "" {
		defaults.Retries = strings.TrimSpace(dc.dubboProxyConfig.Retries)
	}
	if dc.dubboProxyConfig.Timeout != nil {
		if timeout, err := time.ParseDuration(dc.dubboProxyConfig.Timeout.RequestTimeoutStr); err == nil {
			defaults.RequestTimeout = timeout
		}
	}

	return defaults
}

func (dc *Client) preparePayload(req *DubboOutboundRequest) ([]string, []hessian.Object, []byte, error) {
	if len(req.Arguments) == 0 && len(req.ParamTypes) == 0 {
		return []string{}, []hessian.Object{}, []byte("[]"), nil
	}
	if len(req.Arguments) != len(req.ParamTypes) {
		return nil, nil, nil, errors.Errorf("arguments/paramTypes length mismatch: %d vs %d", len(req.Arguments), len(req.ParamTypes))
	}

	types := append([]string(nil), req.ParamTypes...)
	vals := make([]hessian.Object, len(req.Arguments))
	for i, arg := range req.Arguments {
		vals[i] = arg
	}

	finalValues, err := json.Marshal(vals)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "marshal dubbo arguments")
	}

	return types, vals, finalValues, nil
}

func mergeOutboundAttachments(ctx context.Context, outbound map[string]any) map[string]any {
	attachments := make(map[string]any, len(outbound))
	if attaRaw := ctx.Value(constant.AttachmentKey); attaRaw != nil {
		switch userAtta := attaRaw.(type) {
		case map[string]any:
			for key, val := range userAtta {
				attachments[key] = val
			}
		case map[string]string:
			for key, val := range userAtta {
				attachments[key] = val
			}
		}
	}
	for key, val := range outbound {
		attachments[key] = val
	}
	return attachments
}

func prepareInvokeContext(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	if ctx == nil {
		ctx = context.Background()
	}
	if timeout <= 0 {
		return ctx, nil
	}
	return context.WithTimeout(ctx, timeout)
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

func (spec resolvedReferSpec) validate() error {
	switch spec.Mode {
	case "registry":
		if len(spec.RegistryIDs) == 0 {
			return errors.New("dubbo refer mode invalid: registry mode requires registry ids")
		}
		return nil
	case "direct":
		if strings.TrimSpace(spec.URL) == "" {
			return errors.New("dubbo refer mode invalid: direct mode requires direct url")
		}
		return nil
	default:
		return errors.Errorf("dubbo refer mode invalid: %s", spec.Mode)
	}
}

func (spec resolvedReferSpec) cacheKey() (string, error) {
	if err := spec.validate(); err != nil {
		return "", err
	}

	key := spec.genericServiceKey()
	raw, err := json.Marshal(key)
	if err != nil {
		return "", errors.Wrap(err, "marshal generic service key")
	}
	return string(raw), nil
}

func (spec resolvedReferSpec) genericServiceKey() genericServiceKey {
	// Cache key includes all fields that affect reference creation.
	registryIDs := append([]string(nil), spec.RegistryIDs...)
	sort.Strings(registryIDs)
	return genericServiceKey{
		Mode:              spec.Mode,
		URL:               spec.URL,
		RegistryIDs:       registryIDs,
		Cluster:           spec.ConsumerDefaults.Cluster,
		LoadBalance:       spec.ConsumerDefaults.LoadBalance,
		Retries:           spec.ConsumerDefaults.Retries,
		RequestTimeout:    spec.ConsumerDefaults.RequestTimeout.String(),
		Interface:         spec.Interface,
		Version:           spec.Version,
		Group:             spec.Group,
		EffectiveProtocol: spec.EffectiveProtocol,
		Serialization:     spec.EffectiveSerialization,
	}
}

// Get find a dubbo GenericService
func (dc *Client) Get(spec resolvedReferSpec) (*generic.GenericService, error) {
	key, err := spec.cacheKey()
	if err != nil {
		return nil, err
	}
	if dc.check(key) {
		return dc.get(key), nil
	}

	return dc.create(spec)
}

func (dc *Client) create(spec resolvedReferSpec) (*generic.GenericService, error) {
	if err := spec.validate(); err != nil {
		return nil, err
	}
	if dc.dubboClient == nil {
		return nil, errors.New("dubbo client is not initialized, call Apply() first")
	}

	key, err := spec.cacheKey()
	if err != nil {
		return nil, err
	}

	opts, err := dc.buildReferenceOptions(spec)
	if err != nil {
		return nil, err
	}

	dc.lock.Lock()
	defer dc.lock.Unlock()

	// Another request may have built the same GenericService while this one prepared options.
	if service, ok := dc.GenericServicePool[key]; ok {
		return service, nil
	}

	clientService, err := dc.dubboClient.NewGenericService(spec.Interface, opts...)
	if err != nil {
		return nil, err
	}

	if spec.Mode == "registry" && spec.UseNacosWarmup {
		time.Sleep(time.Second)
	}

	dc.GenericServicePool[key] = clientService

	return clientService, nil
}

// buildReferenceOptions builds a list of dubbo-go ReferenceOption using the official API.
func (dc *Client) buildReferenceOptions(spec resolvedReferSpec) ([]dclient.ReferenceOption, error) {
	if err := spec.validate(); err != nil {
		return nil, err
	}
	if strings.TrimSpace(spec.EffectiveProtocol) == "" {
		return nil, errors.New("dubbo refer mode invalid: effective protocol is required")
	}

	defaults := spec.ConsumerDefaults
	opts := make([]dclient.ReferenceOption, 0, 16)

	opts = append(opts, dclient.WithInterface(spec.Interface))
	if spec.Group != "" {
		opts = append(opts, dclient.WithGroup(spec.Group))
	}
	if spec.Version != "" {
		opts = append(opts, dclient.WithVersion(spec.Version))
	}

	// Mode selects either registry discovery or a direct provider URL.
	opts = appendModeReferenceOptions(opts, spec)
	opts = append(opts, clusterReferenceOption(defaults.Cluster))
	opts = append(opts, protocolReferenceOption(spec.EffectiveProtocol))
	if spec.EffectiveSerialization != "" {
		opts = append(opts, dclient.WithSerialization(spec.EffectiveSerialization))
	}

	if loadBalanceOpt := loadBalanceReferenceOption(defaults.LoadBalance); loadBalanceOpt != nil {
		opts = append(opts, loadBalanceOpt)
	}

	retries := 3
	if strings.TrimSpace(defaults.Retries) != "" {
		if resolvedRetries, err := strconv.Atoi(defaults.Retries); err == nil {
			retries = resolvedRetries
		}
	}
	opts = append(opts, dclient.WithRetries(retries))

	timeout := defaults.RequestTimeout
	if timeout <= 0 {
		timeout = cst.DefaultReqTimeout
	}
	opts = append(opts, dclient.WithRequestTimeout(timeout))

	opts = append(opts, dclient.WithGeneric())

	return opts, nil
}

func appendModeReferenceOptions(opts []dclient.ReferenceOption, spec resolvedReferSpec) []dclient.ReferenceOption {
	switch spec.Mode {
	case "registry":
		registryIDs := append([]string(nil), spec.RegistryIDs...)
		sort.Strings(registryIDs)
		return append(opts, dclient.WithRegistryIDs(registryIDs...))
	case "direct":
		return append(opts, dclient.WithURL(spec.URL))
	default:
		return opts
	}
}

func clusterReferenceOption(cluster string) dclient.ReferenceOption {
	switch cluster {
	case "failover":
		return dclient.WithClusterFailOver()
	case "failfast":
		return dclient.WithClusterFailFast()
	case "failsafe":
		return dclient.WithClusterFailSafe()
	case "failback":
		return dclient.WithClusterFailBack()
	case "broadcast":
		return dclient.WithClusterBroadcast()
	case "forking":
		return dclient.WithClusterForking()
	case "available":
		return dclient.WithClusterAvailable()
	case "zoneaware":
		return dclient.WithClusterZoneAware()
	case "adaptiveservice":
		return dclient.WithClusterAdaptiveService()
	default:
		return dclient.WithCluster(cluster)
	}
}

func protocolReferenceOption(protocol string) dclient.ReferenceOption {
	switch protocol {
	case "tri", "triple":
		return dclient.WithProtocolTriple()
	case "dubbo":
		return dclient.WithProtocolDubbo()
	case "jsonrpc":
		return dclient.WithProtocolJsonRPC()
	default:
		return dclient.WithProtocol(protocol)
	}
}

func loadBalanceReferenceOption(loadBalance string) dclient.ReferenceOption {
	switch loadBalance {
	case "":
		return nil
	case "random":
		return dclient.WithLoadBalanceRandom()
	case "roundrobin":
		return dclient.WithLoadBalanceRoundRobin()
	case "leastactive":
		return dclient.WithLoadBalanceLeastActive()
	case "consistenthash", "consistenthashing":
		return dclient.WithLoadBalanceConsistentHashing()
	case "p2c":
		return dclient.WithLoadBalanceP2C()
	default:
		return dclient.WithLoadBalance(loadBalance)
	}
}

func withAttachments(ctx context.Context) context.Context {
	attachments := make(map[string]any)
	if attaRaw := ctx.Value(constant.AttachmentKey); attaRaw != nil {
		switch userAtta := attaRaw.(type) {
		case map[string]any:
			for key, val := range userAtta {
				attachments[key] = val
			}
		case map[string]string:
			for key, val := range userAtta {
				attachments[key] = val
			}
		}
	}

	carrier := propagation.MapCarrier{}
	// Carry tracing headers as Dubbo attachments for the upstream invocation.
	otel.GetTextMapPropagator().Inject(ctx, carrier)
	for key, val := range carrier {
		attachments[key] = val
	}

	return context.WithValue(ctx, constant.AttachmentKey, attachments)
}
