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
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

import (
	dclient "dubbo.apache.org/dubbo-go/v3/client"
	_ "dubbo.apache.org/dubbo-go/v3/cluster/loadbalance/consistenthashing"
	dubboCommon "dubbo.apache.org/dubbo-go/v3/common"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/common/extension"
	dconfig "dubbo.apache.org/dubbo-go/v3/config"
	"dubbo.apache.org/dubbo-go/v3/config/generic"
	"dubbo.apache.org/dubbo-go/v3/global"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	dregistry "dubbo.apache.org/dubbo-go/v3/registry"

	hessian "github.com/apache/dubbo-go-hessian2"

	"github.com/pkg/errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
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
	lock                     sync.RWMutex
	GenericServicePool       map[string]*generic.GenericService
	dubboProxyConfig         *DubboProxyConfig
	registries               map[string]*global.RegistryConfig
	dubboClient              *dclient.Client
	registryProviderResolver func(config.IntegrationRequest, resolvedReferSpec) ([]providerProtocolView, error)
}

type resolvedConsumerDefaults struct {
	Cluster        string
	LoadBalance    string
	Retries        string
	RequestTimeout time.Duration
	Check          *bool
	Filter         string
	Serialization  string
	Sticky         *bool
	Params         map[string]string
	GenericType    string
}

type resolvedReferSpec struct {
	Mode              string
	Interface         string
	Group             string
	Version           string
	URL               string
	RegistryIDs       []string
	EffectiveProtocol string
	UseNacosWarmup    bool
	ConsumerDefaults  resolvedConsumerDefaults
}

type providerProtocolView struct {
	URLProtocol       string
	MetadataProtocol  string
	EndpointProtocols []string
	ServiceProtocols  []string
}

type genericServiceKey struct {
	Mode              string   `json:"mode"`
	URL               string   `json:"url"`
	RegistryIDs       []string `json:"registry_ids"`
	Cluster           string   `json:"cluster"`
	Interface         string   `json:"interface"`
	Version           string   `json:"version"`
	Group             string   `json:"group"`
	EffectiveProtocol string   `json:"effective_protocol"`
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

	spec, err := dc.resolveReferSpec(req.API.IntegrationRequest)
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

	invokeCtx := req.Context
	if invokeCtx == nil {
		invokeCtx = context.Background()
	}
	var cancel context.CancelFunc
	if req.Timeout > 0 {
		invokeCtx, cancel = context.WithTimeout(invokeCtx, req.Timeout)
		defer cancel()
	}

	tr := otel.Tracer(traceNameDubbogoClient)
	ctx, span := tr.Start(invokeCtx, spanNameDubbogoClient)
	span.SetAttributes(attribute.Key(spanTagMethod).String(method))
	span.SetAttributes(attribute.Key(spanTagType).StringSlice(types))
	span.SetAttributes(attribute.Key(spanTagValues).String(string(finalValues)))
	defer span.End()

	ctxWithAttachment := withAttachments(ctx)

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

func (dc *Client) resolveReferSpec(irequest config.IntegrationRequest) (resolvedReferSpec, error) {
	spec := resolvedReferSpec{
		Interface:        irequest.Interface,
		Group:            irequest.Group,
		Version:          irequest.Version,
		ConsumerDefaults: dc.resolveConsumerDefaults(irequest),
	}

	if len(dc.registries) > 0 {
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
		effectiveProtocol, err := dc.resolveRegistryProtocol(irequest, spec)
		if err != nil {
			return resolvedReferSpec{}, err
		}
		spec.EffectiveProtocol = effectiveProtocol
		return spec, nil
	}

	canonicalURL := strings.TrimSpace(irequest.URL)
	if canonicalURL == "" {
		return resolvedReferSpec{}, errors.New("dubbo refer mode invalid: no registry configured and no direct url provided")
	}
	directProtocol, err := directURLProtocol(canonicalURL)
	if err != nil {
		return resolvedReferSpec{}, err
	}

	spec.Mode = "direct"
	spec.URL = canonicalURL
	spec.EffectiveProtocol = directProtocol
	return spec, nil
}

func (dc *Client) resolveConsumerDefaults(irequest config.IntegrationRequest) resolvedConsumerDefaults {
	defaults := resolvedConsumerDefaults{
		Cluster:     "failover",
		Retries:     "3",
		GenericType: "true",
	}

	if dc.dubboProxyConfig == nil {
		if strings.TrimSpace(irequest.Retries) != "" {
			defaults.Retries = strings.TrimSpace(irequest.Retries)
		}
		defaults.RequestTimeout = cst.DefaultReqTimeout
		return defaults
	}

	defaults.Cluster = dc.dubboProxyConfig.GetCluster()
	defaults.LoadBalance = dc.dubboProxyConfig.LoadBalance
	defaults.Check = dc.dubboProxyConfig.GetCheck()
	defaults.Filter = dc.dubboProxyConfig.Filter
	defaults.Serialization = dc.dubboProxyConfig.Serialization
	defaults.Sticky = dc.dubboProxyConfig.Sticky
	defaults.GenericType = dc.dubboProxyConfig.GetGenericType()
	defaults.Params = dc.dubboProxyConfig.Params
	if strings.TrimSpace(dc.dubboProxyConfig.Retries) != "" {
		defaults.Retries = strings.TrimSpace(dc.dubboProxyConfig.Retries)
	} else if strings.TrimSpace(irequest.Retries) != "" {
		defaults.Retries = strings.TrimSpace(irequest.Retries)
	}

	if dc.dubboProxyConfig.Timeout != nil {
		if timeout, err := time.ParseDuration(dc.dubboProxyConfig.Timeout.RequestTimeoutStr); err == nil {
			defaults.RequestTimeout = timeout
		} else {
			defaults.RequestTimeout = cst.DefaultReqTimeout
		}
	} else {
		defaults.RequestTimeout = cst.DefaultReqTimeout
	}

	return defaults
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
	registryIDs := append([]string(nil), spec.RegistryIDs...)
	sort.Strings(registryIDs)
	return genericServiceKey{
		Mode:              spec.Mode,
		URL:               spec.URL,
		RegistryIDs:       registryIDs,
		Cluster:           spec.ConsumerDefaults.Cluster,
		Interface:         spec.Interface,
		Version:           spec.Version,
		Group:             spec.Group,
		EffectiveProtocol: spec.EffectiveProtocol,
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

	switch spec.Mode {
	case "registry":
		registryIDs := append([]string(nil), spec.RegistryIDs...)
		sort.Strings(registryIDs)
		opts = append(opts, dclient.WithRegistryIDs(registryIDs...))
	case "direct":
		opts = append(opts, dclient.WithURL(spec.URL))
	}

	switch defaults.Cluster {
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
		opts = append(opts, dclient.WithCluster(defaults.Cluster))
	}

	switch spec.EffectiveProtocol {
	case "tri", "triple":
		opts = append(opts, dclient.WithProtocolTriple())
	case "dubbo":
		opts = append(opts, dclient.WithProtocolDubbo())
	case "jsonrpc":
		opts = append(opts, dclient.WithProtocolJsonRPC())
	default:
		opts = append(opts, dclient.WithProtocol(spec.EffectiveProtocol))
	}

	if defaults.LoadBalance != "" {
		switch defaults.LoadBalance {
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
			opts = append(opts, dclient.WithLoadBalance(defaults.LoadBalance))
		}
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

	if defaults.Check != nil {
		opts = append(opts, withReferenceCheck(*defaults.Check))
	}
	if defaults.Filter != "" {
		opts = append(opts, dclient.WithFilter(defaults.Filter))
	}
	if defaults.Serialization != "" {
		switch defaults.Serialization {
		case "json":
			opts = append(opts, dclient.WithSerializationJSON())
		default:
			opts = append(opts, dclient.WithSerialization(defaults.Serialization))
		}
	}

	opts = append(opts, dclient.WithGenericType(defaults.GenericType))

	if defaults.Sticky != nil && *defaults.Sticky {
		opts = append(opts, dclient.WithSticky())
	}
	if len(defaults.Params) > 0 {
		opts = append(opts, dclient.WithParams(defaults.Params))
	}

	return opts, nil
}

func (dc *Client) resolveRegistryProtocol(irequest config.IntegrationRequest, spec resolvedReferSpec) (string, error) {
	resolver := dc.registryProviderResolver
	if resolver == nil {
		resolver = dc.resolveRegistryProviderViews
	}

	views, err := resolver(irequest, spec)
	if err != nil {
		return "", err
	}

	return resolveProviderProtocolViews(views)
}

func resolveProviderProtocolViews(views []providerProtocolView) (string, error) {
	candidates := make(map[string]struct{})
	for _, view := range views {
		addProtocolCandidate(candidates, view.URLProtocol)
		addProtocolCandidate(candidates, view.MetadataProtocol)
		for _, protocol := range view.EndpointProtocols {
			addProtocolCandidate(candidates, protocol)
		}
		for _, protocol := range view.ServiceProtocols {
			addProtocolCandidate(candidates, protocol)
		}
	}

	if len(candidates) == 0 {
		return "", errors.New("dubbo refer mode invalid: registry provider protocol metadata is missing")
	}
	if len(candidates) > 1 {
		values := make([]string, 0, len(candidates))
		for protocol := range candidates {
			values = append(values, protocol)
		}
		sort.Strings(values)
		return "", errors.Errorf("dubbo refer mode invalid: registry provider protocols are ambiguous: %v", values)
	}
	for protocol := range candidates {
		return protocol, nil
	}
	return "", errors.New("dubbo refer mode invalid: registry provider protocol metadata is missing")
}

func addProtocolCandidate(candidates map[string]struct{}, protocol string) {
	if normalized := normalizeReferenceProtocol(protocol); normalized != "" {
		candidates[normalized] = struct{}{}
	}
}

func (dc *Client) resolveRegistryProviderViews(irequest config.IntegrationRequest, spec resolvedReferSpec) ([]providerProtocolView, error) {
	views := make([]providerProtocolView, 0, len(spec.RegistryIDs))
	for _, registryID := range spec.RegistryIDs {
		registryCfg := dc.registries[registryID]
		if registryCfg == nil {
			continue
		}

		registryInstance, err := newRegistryInstance(registryID, registryCfg)
		if err != nil {
			return nil, err
		}
		if registryInstance == nil {
			continue
		}

		registryViews, resolveErr := dc.resolveRegistryViewsFromInstance(irequest, registryID, registryInstance)
		registryInstance.Destroy()
		if resolveErr != nil {
			return nil, resolveErr
		}
		views = append(views, registryViews...)
	}

	return views, nil
}

func newRegistryInstance(registryID string, registryCfg *global.RegistryConfig) (dregistry.Registry, error) {
	consumerCfg := &dconfig.RegistryConfig{
		Protocol:     registryCfg.Protocol,
		Timeout:      registryCfg.Timeout,
		Group:        registryCfg.Group,
		Namespace:    registryCfg.Namespace,
		TTL:          registryCfg.TTL,
		Address:      registryCfg.Address,
		Username:     registryCfg.Username,
		Password:     registryCfg.Password,
		Simplified:   registryCfg.Simplified,
		Preferred:    registryCfg.Preferred,
		Zone:         registryCfg.Zone,
		Weight:       registryCfg.Weight,
		Params:       registryCfg.Params,
		RegistryType: registryCfg.RegistryType,
	}

	registryInstance, err := consumerCfg.GetInstance(dubboCommon.CONSUMER)
	if err != nil {
		return nil, errors.Wrapf(err, "create registry instance %q", registryID)
	}
	return registryInstance, nil
}

func (dc *Client) resolveRegistryViewsFromInstance(irequest config.IntegrationRequest, registryID string, registryInstance dregistry.Registry) ([]providerProtocolView, error) {
	if discoveryProvider, ok := registryInstance.(interface {
		GetServiceDiscovery() dregistry.ServiceDiscovery
	}); ok {
		return dc.resolveServiceDiscoveryProviderViews(irequest, registryID, discoveryProvider.GetServiceDiscovery())
	}

	consumerURL, err := buildRegistryLookupConsumerURL(irequest)
	if err != nil {
		return nil, err
	}

	listener := &providerEventCollector{}
	if err := registryInstance.LoadSubscribeInstances(consumerURL, listener); err != nil {
		return nil, errors.Wrapf(err, "load provider instances from registry %q", registryID)
	}
	return listener.views, nil
}

func (dc *Client) resolveServiceDiscoveryProviderViews(irequest config.IntegrationRequest, registryID string, discovery dregistry.ServiceDiscovery) ([]providerProtocolView, error) {
	serviceNames, err := dc.resolveServiceDiscoveryNames(irequest)
	if err != nil {
		return nil, errors.Wrapf(err, "resolve service discovery names from registry %q", registryID)
	}

	views := make([]providerProtocolView, 0)
	for _, serviceName := range serviceNames {
		instances := discovery.GetInstances(serviceName)
		for _, instance := range instances {
			views = append(views, buildProviderViewFromInstance(irequest, instance))
		}
	}
	return views, nil
}

func (dc *Client) resolveServiceDiscoveryNames(irequest config.IntegrationRequest) ([]string, error) {
	if application := strings.TrimSpace(irequest.ApplicationName); application != "" {
		return []string{application}, nil
	}

	consumerURL, err := buildRegistryLookupConsumerURL(irequest)
	if err != nil {
		return nil, err
	}
	services, err := extension.GetGlobalServiceNameMapping().Get(consumerURL, nil)
	if err != nil {
		return nil, err
	}
	values := services.Values()
	names := make([]string, 0, len(values))
	for _, service := range values {
		if name, ok := service.(string); ok && strings.TrimSpace(name) != "" {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	return names, nil
}

func buildRegistryLookupConsumerURL(irequest config.IntegrationRequest) (*dubboCommon.URL, error) {
	params := url.Values{}
	if irequest.Interface != "" {
		params.Set(constant.InterfaceKey, irequest.Interface)
	}
	if irequest.Group != "" {
		params.Set(constant.GroupKey, irequest.Group)
	}
	if irequest.Version != "" {
		params.Set(constant.VersionKey, irequest.Version)
	}
	if irequest.ApplicationName != "" {
		params.Set(constant.ApplicationKey, irequest.ApplicationName)
	}
	params.Set(constant.SideKey, constant.SideConsumer)

	return dubboCommon.NewURL("consumer://127.0.0.1/"+irequest.Interface, dubboCommon.WithParams(params))
}

func buildProviderViewFromInstance(irequest config.IntegrationRequest, instance dregistry.ServiceInstance) providerProtocolView {
	view := providerProtocolView{
		MetadataProtocol: instance.GetMetadata()["protocol"],
	}

	for _, endpoint := range instance.GetEndPoints() {
		view.EndpointProtocols = append(view.EndpointProtocols, endpoint.Protocol)
	}

	if metadata := instance.GetServiceMetadata(); metadata != nil {
		for _, service := range metadata.Services {
			if service == nil || !serviceMatchesRequest(service, irequest) {
				continue
			}
			view.ServiceProtocols = append(view.ServiceProtocols, service.Protocol)
			if service.URL != nil {
				view.ServiceProtocols = append(view.ServiceProtocols, service.URL.Protocol)
			}
		}
	}

	return view
}

func serviceMatchesRequest(service interface {
	GetServiceKey() string
}, irequest config.IntegrationRequest) bool {
	expectedServiceKey := dubboCommon.ServiceKey(irequest.Interface, irequest.Group, irequest.Version)
	return service.GetServiceKey() == expectedServiceKey
}

type providerEventCollector struct {
	views []providerProtocolView
}

func (c *providerEventCollector) Notify(event *dregistry.ServiceEvent) {
	if event == nil || event.Service == nil {
		return
	}
	c.views = append(c.views, providerProtocolView{
		URLProtocol: event.Service.Protocol,
	})
}

func (c *providerEventCollector) NotifyAll(events []*dregistry.ServiceEvent, callback func()) {
	for _, event := range events {
		c.Notify(event)
	}
	if callback != nil {
		callback()
	}
}

func directURLProtocol(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", errors.Wrapf(err, "parse direct url %q", rawURL)
	}
	return normalizeReferenceProtocol(parsedURL.Scheme), nil
}

func normalizeReferenceProtocol(protocol string) string {
	normalized := strings.ToLower(strings.TrimSpace(protocol))
	switch normalized {
	case "triple":
		return "tri"
	default:
		return normalized
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
	otel.GetTextMapPropagator().Inject(ctx, carrier)
	for key, val := range carrier {
		attachments[key] = val
	}

	return context.WithValue(ctx, constant.AttachmentKey, attachments)
}

func withReferenceCheck(check bool) dclient.ReferenceOption {
	return func(opts *dclient.ReferenceOptions) {
		v := check
		opts.Reference.Check = &v
	}
}
