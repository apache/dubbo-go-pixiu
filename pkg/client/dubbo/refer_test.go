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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"
)

import (
	dclient "dubbo.apache.org/dubbo-go/v3/client"
	dubboConstant "dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/config/generic"
	"dubbo.apache.org/dubbo-go/v3/global"

	hessian "github.com/apache/dubbo-go-hessian2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pkg/common/mock"
	"github.com/apache/dubbo-go-pixiu/pkg/config"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

type fixedPropagator struct{}

func (fixedPropagator) Inject(_ context.Context, _ propagation.TextMapCarrier) {}

func (fixedPropagator) Extract(ctx context.Context, _ propagation.TextMapCarrier) context.Context {
	return ctx
}

func (fixedPropagator) Fields() []string {
	return nil
}

type attachmentPropagator struct{}

func (attachmentPropagator) Inject(_ context.Context, carrier propagation.TextMapCarrier) {
	carrier.Set("traceparent", "00-test-traceparent")
	carrier.Set("x-otel-scope", "pixiu")
}

func (attachmentPropagator) Extract(ctx context.Context, _ propagation.TextMapCarrier) context.Context {
	return ctx
}

func (attachmentPropagator) Fields() []string {
	return []string{"traceparent", "x-otel-scope"}
}

func TestCallUsesMapParamsBeforeRefer(t *testing.T) {
	dc := NewDubboClient()
	restorePropagator(t, fixedPropagator{})

	req := newDirectCallRequest(t, context.Background(), "/invoke", "/invoke", `{"values":["123"],"types":"java.lang.String"}`, []config.MappingParam{
		{Name: "headers.X-App", MapTo: "opt.application"},
		{Name: "headers.X-Interface", MapTo: "opt.interface"},
		{Name: "headers.X-Method", MapTo: "opt.method"},
		{Name: "requestBody.values", MapTo: "opt.values"},
		{Name: "requestBody.types", MapTo: "opt.types"},
	})
	req.IngressRequest.Header.Set("X-App", "mapped-app")
	req.IngressRequest.Header.Set("X-Interface", "com.example.OrderService")
	req.IngressRequest.Header.Set("X-Method", "CreateOrder")
	req.API.ApplicationName = "initial-app"
	req.API.Interface = "initial.Interface"
	req.API.Method.Method = "initialMethod"

	expectedIR := req.API.IntegrationRequest
	expectedIR.ApplicationName = "mapped-app"
	expectedIR.Interface = "com.example.OrderService"
	expectedIR.Method = "CreateOrder"

	cached := false
	cacheServiceForRequest(t, dc, expectedIR, &generic.GenericService{
		Invoke: func(ctx context.Context, methodName string, types []string, args []hessian.Object) (any, error) {
			cached = true
			assert.Equal(t, "CreateOrder", methodName)
			assert.Equal(t, []string{"java.lang.String"}, types)
			require.Len(t, args, 1)
			assert.Equal(t, "123", args[0])
			return "ok", nil
		},
	})

	res, err := dc.Call(req)
	require.NoError(t, err)
	assert.Equal(t, "ok", res)
	assert.True(t, cached)
	assert.Equal(t, "mapped-app", req.API.ApplicationName)
	assert.Equal(t, "com.example.OrderService", req.API.Interface)
	assert.Equal(t, "CreateOrder", req.API.Method.Method)
}

func TestCallAutoResolveSnapshotUsesMappedFields(t *testing.T) {
	dc := NewDubboClient()
	restorePropagator(t, fixedPropagator{})

	req := newDirectCallRequest(t, context.Background(), "/:application/:interface/:method", "/demo-app/com.example.UserService/GetUser", `{"values":["u-1"],"types":"java.lang.String"}`, []config.MappingParam{
		{Name: "requestBody.values", MapTo: "opt.values"},
		{Name: "requestBody.types", MapTo: "opt.types"},
		{Name: "uri.application", MapTo: "opt.application"},
		{Name: "uri.interface", MapTo: "opt.interface"},
		{Name: "uri.method", MapTo: "opt.method"},
		{Name: "headers.X-Group", MapTo: "opt.group"},
		{Name: "headers.X-Version", MapTo: "opt.version"},
	})
	req.IngressRequest.Header.Set("X-Group", "gray")
	req.IngressRequest.Header.Set("X-Version", "1.0.0")

	expectedIR := req.API.IntegrationRequest
	expectedIR.ApplicationName = "demo-app"
	expectedIR.Interface = "com.example.UserService"
	expectedIR.Method = "GetUser"
	expectedIR.Group = "gray"
	expectedIR.Version = "1.0.0"

	cacheServiceForRequest(t, dc, expectedIR, &generic.GenericService{
		Invoke: func(ctx context.Context, methodName string, types []string, args []hessian.Object) (any, error) {
			assert.Equal(t, "GetUser", methodName)
			assert.Equal(t, []string{"java.lang.String"}, types)
			return "resolved", nil
		},
	})

	res, err := dc.Call(req)
	require.NoError(t, err)
	assert.Equal(t, "resolved", res)
	assert.Equal(t, "demo-app", req.API.ApplicationName)
	assert.Equal(t, "com.example.UserService", req.API.Interface)
	assert.Equal(t, "GetUser", req.API.Method.Method)
	assert.Equal(t, "gray", req.API.Group)
	assert.Equal(t, "1.0.0", req.API.Version)
}

func TestGetAndCreateReturnErrorsInsteadOfPanic(t *testing.T) {
	dc := NewDubboClient()

	assert.NotPanics(t, func() {
		service, err := dc.Get(resolvedReferSpec{Mode: "invalid"})
		assert.Nil(t, service)
		assert.EqualError(t, err, "dubbo refer mode invalid: invalid")
	})

	spec := resolvedReferSpec{
		Mode:              "direct",
		URL:               "dubbo://127.0.0.1:20880",
		Interface:         "com.example.UserService",
		EffectiveProtocol: "dubbo",
		ConsumerDefaults: resolvedConsumerDefaults{
			Cluster: "failover",
		},
	}

	assert.NotPanics(t, func() {
		service, err := dc.create(spec)
		assert.Nil(t, service)
		assert.EqualError(t, err, "dubbo client is not initialized, call Apply() first")
	})
}

func TestBuildReferSnapshotRejectsMissingDirectURL(t *testing.T) {
	dc := NewDubboClient()

	_, err := dc.resolveReferSpec(config.IntegrationRequest{})
	assert.EqualError(t, err, "dubbo refer mode invalid: no registry configured and no direct url provided")
}

func TestResolveReferSpecDirectModeResolvesEffectiveProtocol(t *testing.T) {
	dc := NewDubboClient()

	spec, err := dc.resolveReferSpec(config.IntegrationRequest{
		HTTPBackendConfig: config.HTTPBackendConfig{
			URL: "  triple://127.0.0.1:50051  ",
		},
		DubboBackendConfig: config.DubboBackendConfig{
			ApplicationName: "demo-app",
			Interface:       "com.example.UserService",
			Group:           "gray",
			Version:         "1.0.0",
		},
	})
	require.NoError(t, err)
	assert.Equal(t, "direct", spec.Mode)
	assert.Equal(t, "triple://127.0.0.1:50051", spec.URL)
	assert.Equal(t, "tri", spec.EffectiveProtocol)
	assert.Equal(t, "com.example.UserService", spec.Interface)
	assert.Equal(t, "gray", spec.Group)
	assert.Equal(t, "1.0.0", spec.Version)
}

func TestResolvedReferSpecCacheKeyIncludesEffectiveProtocol(t *testing.T) {
	dubboSpec := resolvedReferSpec{
		Mode:              "registry",
		Interface:         "com.example.UserService",
		Group:             "gray",
		Version:           "1.0.0",
		RegistryIDs:       []string{"demo"},
		EffectiveProtocol: "dubbo",
	}
	triSpec := dubboSpec
	triSpec.EffectiveProtocol = "tri"

	dubboKey, err := dubboSpec.cacheKey()
	require.NoError(t, err)
	triKey, err := triSpec.cacheKey()
	require.NoError(t, err)

	assert.NotEqual(t, dubboKey, triKey)
}

func TestBuildReferenceOptionsUsesResolvedDirectSpec(t *testing.T) {
	dc := NewDubboClient()

	spec := resolvedReferSpec{
		Mode:              "direct",
		Interface:         "com.example.UserService",
		URL:               "triple://127.0.0.1:50051",
		EffectiveProtocol: "tri",
		ConsumerDefaults: resolvedConsumerDefaults{
			Cluster: "failover",
		},
	}

	refOpts := applyResolvedReferenceOptions(t, dc, spec)
	assert.Equal(t, "triple://127.0.0.1:50051", refOpts.Reference.URL)
	assert.Equal(t, "tri", refOpts.Reference.Protocol)
	assert.Empty(t, refOpts.Reference.RegistryIDs)
}

func TestDirectModeUsesCanonicalURLForOptionsAndCacheKey(t *testing.T) {
	dc := NewDubboClient()

	spec, err := dc.resolveReferSpec(config.IntegrationRequest{
		HTTPBackendConfig: config.HTTPBackendConfig{
			URL: "  dubbo://127.0.0.1:20880  ",
		},
		DubboBackendConfig: config.DubboBackendConfig{
			ClusterName:     "demo-cluster",
			ApplicationName: "demo-app",
			Interface:       "com.example.UserService",
			Group:           "gray",
			Version:         "1.0.0",
		},
	})
	require.NoError(t, err)
	assert.Equal(t, "direct", spec.Mode)
	assert.Equal(t, "dubbo://127.0.0.1:20880", spec.URL)
	assert.Equal(t, "dubbo", spec.EffectiveProtocol)

	refOpts := applyResolvedReferenceOptions(t, dc, spec)
	assert.Equal(t, "dubbo://127.0.0.1:20880", refOpts.Reference.URL)
	assert.Empty(t, refOpts.Reference.RegistryIDs)
	assert.Equal(t, dubboConstant.DubboProtocol, refOpts.Reference.Protocol)

	key := decodeResolvedGenericServiceKey(t, spec)
	assert.Equal(t, "dubbo://127.0.0.1:20880", key.URL)
	assert.Equal(t, "direct", key.Mode)
}

func TestDirectModeUsesURLProtocolInsteadOfConfiguredProtocol(t *testing.T) {
	dc := NewDubboClient()
	dc.SetConfig(&DubboProxyConfig{
		Protocol: "tri",
	})

	spec, err := dc.resolveReferSpec(config.IntegrationRequest{
		HTTPBackendConfig: config.HTTPBackendConfig{
			URL: "dubbo://127.0.0.1:20880",
		},
		DubboBackendConfig: config.DubboBackendConfig{
			Interface: "com.example.UserService",
		},
	})
	require.NoError(t, err)
	assert.Equal(t, "dubbo", spec.EffectiveProtocol)

	refOpts := applyResolvedReferenceOptions(t, dc, spec)
	assert.Equal(t, dubboConstant.DubboProtocol, refOpts.Reference.Protocol)
}

func TestBuildReferenceOptionsNormalizesTripleDirectProtocol(t *testing.T) {
	dc := NewDubboClient()
	dc.SetConfig(&DubboProxyConfig{
		Protocol: "tri",
	})

	spec, err := dc.resolveReferSpec(config.IntegrationRequest{
		HTTPBackendConfig: config.HTTPBackendConfig{
			URL: "triple://127.0.0.1:50051",
		},
		DubboBackendConfig: config.DubboBackendConfig{
			Interface: "com.example.UserService",
		},
	})
	require.NoError(t, err)
	assert.Equal(t, "tri", spec.EffectiveProtocol)

	refOpts := applyResolvedReferenceOptions(t, dc, spec)
	assert.Equal(t, "triple://127.0.0.1:50051", refOpts.Reference.URL)
	assert.Equal(t, "tri", refOpts.Reference.Protocol)
}

func TestRegistryModeIgnoresURLAndSortsRegistryIDs(t *testing.T) {
	dc := NewDubboClient()
	dc.registries = map[string]*global.RegistryConfig{
		"z-reg": {Protocol: "zookeeper"},
		"a-reg": {Protocol: "nacos"},
	}
	dc.registryProviderResolver = func(irequest config.IntegrationRequest, spec resolvedReferSpec) ([]providerProtocolView, error) {
		return []providerProtocolView{
			{URLProtocol: "dubbo"},
		}, nil
	}

	spec, err := dc.resolveReferSpec(config.IntegrationRequest{
		HTTPBackendConfig: config.HTTPBackendConfig{
			URL: "dubbo://should-be-ignored",
		},
		DubboBackendConfig: config.DubboBackendConfig{
			ApplicationName: "demo-app",
			Interface:       "com.example.UserService",
		},
	})
	require.NoError(t, err)
	assert.Equal(t, "registry", spec.Mode)
	assert.Equal(t, []string{"a-reg", "z-reg"}, spec.RegistryIDs)
	assert.True(t, spec.UseNacosWarmup)
	assert.Empty(t, spec.URL)
	assert.Equal(t, "dubbo", spec.EffectiveProtocol)
}

func TestBuildReferenceOptionsRejectsRegistryModeWithoutResolvedProtocol(t *testing.T) {
	dc := NewDubboClient()
	dc.registries = map[string]*global.RegistryConfig{
		"zk": {Protocol: "zookeeper"},
	}
	dc.registryProviderResolver = func(irequest config.IntegrationRequest, spec resolvedReferSpec) ([]providerProtocolView, error) {
		return nil, errors.New("dubbo refer mode invalid: registry provider protocol metadata is missing")
	}

	spec, err := dc.resolveReferSpec(config.IntegrationRequest{
		DubboBackendConfig: config.DubboBackendConfig{
			Interface: "com.example.UserService",
		},
	})
	assert.EqualError(t, err, "dubbo refer mode invalid: registry provider protocol metadata is missing")
	assert.Equal(t, resolvedReferSpec{}, spec)
}

func TestResolveReferSpecRegistryModeUsesProviderURLProtocol(t *testing.T) {
	dc := NewDubboClient()
	dc.registries = map[string]*global.RegistryConfig{
		"zk": {Protocol: "zookeeper"},
	}
	dc.registryProviderResolver = func(irequest config.IntegrationRequest, spec resolvedReferSpec) ([]providerProtocolView, error) {
		return []providerProtocolView{
			{URLProtocol: "dubbo"},
		}, nil
	}

	spec, err := dc.resolveReferSpec(config.IntegrationRequest{
		RequestType: "triple",
		DubboBackendConfig: config.DubboBackendConfig{
			ApplicationName: "demo-app",
			Interface:       "com.example.UserService",
		},
	})
	require.NoError(t, err)
	assert.Equal(t, "registry", spec.Mode)
	assert.Equal(t, "dubbo", spec.EffectiveProtocol)
}

func TestResolveReferSpecRegistryModeUsesProviderMetadataProtocol(t *testing.T) {
	dc := NewDubboClient()
	dc.registries = map[string]*global.RegistryConfig{
		"zk": {Protocol: "zookeeper"},
	}
	dc.registryProviderResolver = func(irequest config.IntegrationRequest, spec resolvedReferSpec) ([]providerProtocolView, error) {
		return []providerProtocolView{
			{MetadataProtocol: "triple"},
		}, nil
	}

	spec, err := dc.resolveReferSpec(config.IntegrationRequest{
		RequestType: "dubbo",
		DubboBackendConfig: config.DubboBackendConfig{
			ApplicationName: "demo-app",
			Interface:       "com.example.UserService",
		},
	})
	require.NoError(t, err)
	assert.Equal(t, "tri", spec.EffectiveProtocol)
}

func TestResolveReferSpecRegistryModeRejectsMissingProviderProtocol(t *testing.T) {
	dc := NewDubboClient()
	dc.registries = map[string]*global.RegistryConfig{
		"zk": {Protocol: "zookeeper"},
	}
	dc.registryProviderResolver = func(irequest config.IntegrationRequest, spec resolvedReferSpec) ([]providerProtocolView, error) {
		return []providerProtocolView{
			{},
		}, nil
	}

	_, err := dc.resolveReferSpec(config.IntegrationRequest{
		DubboBackendConfig: config.DubboBackendConfig{
			ApplicationName: "demo-app",
			Interface:       "com.example.UserService",
		},
	})
	assert.EqualError(t, err, "dubbo refer mode invalid: registry provider protocol metadata is missing")
}

func TestResolveReferSpecRegistryModeRejectsMixedProviderProtocols(t *testing.T) {
	dc := NewDubboClient()
	dc.registries = map[string]*global.RegistryConfig{
		"zk": {Protocol: "zookeeper"},
	}
	dc.registryProviderResolver = func(irequest config.IntegrationRequest, spec resolvedReferSpec) ([]providerProtocolView, error) {
		return []providerProtocolView{
			{URLProtocol: "dubbo"},
			{MetadataProtocol: "tri"},
		}, nil
	}

	_, err := dc.resolveReferSpec(config.IntegrationRequest{
		DubboBackendConfig: config.DubboBackendConfig{
			ApplicationName: "demo-app",
			Interface:       "com.example.UserService",
		},
	})
	assert.EqualError(t, err, "dubbo refer mode invalid: registry provider protocols are ambiguous: [dubbo tri]")
}

func TestBuildReferenceOptionsUsesResolvedRegistryProtocolInsteadOfProxyConfig(t *testing.T) {
	dc := NewDubboClient()
	dc.SetConfig(&DubboProxyConfig{
		Protocol: "tri",
	})
	dc.registries = map[string]*global.RegistryConfig{
		"zk": {Protocol: "zookeeper"},
	}
	dc.registryProviderResolver = func(irequest config.IntegrationRequest, spec resolvedReferSpec) ([]providerProtocolView, error) {
		return []providerProtocolView{
			{URLProtocol: "dubbo"},
		}, nil
	}

	spec, err := dc.resolveReferSpec(config.IntegrationRequest{
		RequestType: "triple",
		DubboBackendConfig: config.DubboBackendConfig{
			ApplicationName: "demo-app",
			Interface:       "com.example.UserService",
		},
	})
	require.NoError(t, err)

	refOpts := applyResolvedReferenceOptions(t, dc, spec)
	assert.Equal(t, []string{"zk"}, refOpts.Reference.RegistryIDs)
	assert.Equal(t, dubboConstant.DubboProtocol, refOpts.Reference.Protocol)
}

func TestCallWithZeroTimeoutDoesNotCreateImmediateDeadline(t *testing.T) {
	dc := NewDubboClient()
	restorePropagator(t, fixedPropagator{})

	req := newDirectCallRequest(t, context.Background(), "/invoke", "/invoke", `{"values":["1"],"types":"java.lang.String"}`, []config.MappingParam{
		{Name: "requestBody.values", MapTo: "opt.values"},
		{Name: "requestBody.types", MapTo: "opt.types"},
	})
	req.Timeout = 0

	sawDeadline := false
	cacheServiceForRequest(t, dc, req.API.IntegrationRequest, &generic.GenericService{
		Invoke: func(ctx context.Context, methodName string, types []string, args []hessian.Object) (any, error) {
			_, sawDeadline = ctx.Deadline()
			return "ok", nil
		},
	})

	res, err := dc.Call(req)
	require.NoError(t, err)
	assert.Equal(t, "ok", res)
	assert.False(t, sawDeadline)
}

func TestCallWithTimeoutCreatesPerCallDeadline(t *testing.T) {
	dc := NewDubboClient()
	restorePropagator(t, fixedPropagator{})

	req := newDirectCallRequest(t, context.Background(), "/invoke", "/invoke", `{"values":["1"],"types":"java.lang.String"}`, []config.MappingParam{
		{Name: "requestBody.values", MapTo: "opt.values"},
		{Name: "requestBody.types", MapTo: "opt.types"},
	})
	req.Timeout = 80 * time.Millisecond

	var deadline time.Time
	cacheServiceForRequest(t, dc, req.API.IntegrationRequest, &generic.GenericService{
		Invoke: func(ctx context.Context, methodName string, types []string, args []hessian.Object) (any, error) {
			var ok bool
			deadline, ok = ctx.Deadline()
			assert.True(t, ok)
			return "ok", nil
		},
	})

	start := time.Now()
	res, err := dc.Call(req)
	require.NoError(t, err)
	assert.Equal(t, "ok", res)
	assert.False(t, deadline.IsZero())
	assert.WithinDuration(t, start.Add(req.Timeout), deadline, 100*time.Millisecond)
}

func TestBuildReferenceOptionsUsesStaticRequestTimeoutAndConfigPriority(t *testing.T) {
	check := false
	sticky := true
	dc := NewDubboClient()
	dc.SetConfig(&DubboProxyConfig{
		Timeout: &model.TimeoutConfig{
			RequestTimeoutStr: "6s",
		},
		Cluster:       "failfast",
		Protocol:      "dubbo",
		LoadBalance:   "consistenthash",
		Retries:       "7",
		Check:         &check,
		Filter:        "tracing,metrics",
		Serialization: "json",
		Sticky:        &sticky,
		Params: map[string]string{
			"env": "test",
		},
		GenericType: "true",
	})

	spec, err := dc.resolveReferSpec(config.IntegrationRequest{
		HTTPBackendConfig: config.HTTPBackendConfig{
			URL: "dubbo://127.0.0.1:20880",
		},
		DubboBackendConfig: config.DubboBackendConfig{
			Interface: "com.example.UserService",
			Retries:   "2",
		},
	})
	require.NoError(t, err)

	refOpts := applyResolvedReferenceOptions(t, dc, spec)
	assert.Equal(t, dubboConstant.ClusterKeyFailfast, refOpts.Reference.Cluster)
	assert.Equal(t, dubboConstant.DubboProtocol, refOpts.Reference.Protocol)
	assert.Equal(t, dubboConstant.LoadBalanceKeyConsistentHashing, refOpts.Reference.Loadbalance)
	assert.Equal(t, "7", refOpts.Reference.Retries)
	assert.Equal(t, "6s", refOpts.Reference.RequestTimeout)
	require.NotNil(t, refOpts.Reference.Check)
	assert.False(t, *refOpts.Reference.Check)
	assert.Equal(t, "tracing,metrics", refOpts.Reference.Filter)
	assert.Equal(t, dubboConstant.JSONSerialization, refOpts.Reference.Serialization)
	assert.Equal(t, "true", refOpts.Reference.Generic)
	assert.True(t, refOpts.Reference.Sticky)
	assert.Equal(t, map[string]string{"env": "test"}, refOpts.Reference.Params)
}

func TestCallFlattensAttachmentsIntoMapAny(t *testing.T) {
	dc := NewDubboClient()
	restorePropagator(t, attachmentPropagator{})

	baseCtx := context.WithValue(context.Background(), dubboConstant.AttachmentKey, map[string]any{
		"user-key": "user-value",
	})
	req := newDirectCallRequest(t, baseCtx, "/invoke", "/invoke", `{"values":["1"],"types":"java.lang.String"}`, []config.MappingParam{
		{Name: "requestBody.values", MapTo: "opt.values"},
		{Name: "requestBody.types", MapTo: "opt.types"},
	})

	cacheServiceForRequest(t, dc, req.API.IntegrationRequest, &generic.GenericService{
		Invoke: func(ctx context.Context, methodName string, types []string, args []hessian.Object) (any, error) {
			attachments, ok := ctx.Value(dubboConstant.AttachmentKey).(map[string]any)
			require.True(t, ok)
			assert.Equal(t, "user-value", attachments["user-key"])
			assert.Equal(t, "00-test-traceparent", attachments["traceparent"])
			assert.Equal(t, "pixiu", attachments["x-otel-scope"])
			assert.NotContains(t, attachments, "carrier")
			return "ok", nil
		},
	})

	res, err := dc.Call(req)
	require.NoError(t, err)
	assert.Equal(t, "ok", res)
}

func TestCallPreservesStringMapAttachments(t *testing.T) {
	dc := NewDubboClient()
	restorePropagator(t, attachmentPropagator{})

	baseCtx := context.WithValue(context.Background(), dubboConstant.AttachmentKey, map[string]string{
		"user-key": "user-value",
	})
	req := newDirectCallRequest(t, baseCtx, "/invoke", "/invoke", `{"values":["1"],"types":"java.lang.String"}`, []config.MappingParam{
		{Name: "requestBody.values", MapTo: "opt.values"},
		{Name: "requestBody.types", MapTo: "opt.types"},
	})

	cacheServiceForRequest(t, dc, req.API.IntegrationRequest, &generic.GenericService{
		Invoke: func(ctx context.Context, methodName string, types []string, args []hessian.Object) (any, error) {
			attachments, ok := ctx.Value(dubboConstant.AttachmentKey).(map[string]any)
			require.True(t, ok)
			assert.Equal(t, "user-value", attachments["user-key"])
			assert.Equal(t, "00-test-traceparent", attachments["traceparent"])
			assert.Equal(t, "pixiu", attachments["x-otel-scope"])
			return "ok", nil
		},
	})

	res, err := dc.Call(req)
	require.NoError(t, err)
	assert.Equal(t, "ok", res)
}

func newDirectCallRequest(t *testing.T, ctx context.Context, urlPattern string, rawURL string, body string, mapping []config.MappingParam) *client.Request {
	t.Helper()

	request, err := http.NewRequest(http.MethodPost, rawURL, bytes.NewBufferString(body))
	require.NoError(t, err)

	api := mock.GetMockAPI(http.MethodPost, urlPattern)
	api.MappingParams = mapping
	api.IntegrationRequest.URL = "  dubbo://127.0.0.1:20880  "
	api.Interface = "com.example.UserService"
	api.Method.Method = "SayHello"
	api.ApplicationName = "demo-app"
	api.Group = "demo-group"
	api.Version = "1.0.0"

	return client.NewReq(ctx, request, api)
}

func cacheServiceForRequest(t *testing.T, dc *Client, irequest config.IntegrationRequest, service *generic.GenericService) {
	t.Helper()

	spec, err := dc.resolveReferSpec(irequest)
	require.NoError(t, err)
	key, err := spec.cacheKey()
	require.NoError(t, err)
	dc.GenericServicePool[key] = service
}

func applyResolvedReferenceOptions(t *testing.T, dc *Client, spec resolvedReferSpec) *dclient.ReferenceOptions {
	t.Helper()

	opts, err := dc.buildReferenceOptions(spec)
	require.NoError(t, err)

	refOpts := &dclient.ReferenceOptions{
		Reference:   global.DefaultReferenceConfig(),
		Application: global.DefaultApplicationConfig(),
		Shutdown:    global.DefaultShutdownConfig(),
		Metrics:     global.DefaultMetricsConfig(),
		Otel:        global.DefaultOtelConfig(),
		TLS:         global.DefaultTLSConfig(),
		Protocols:   make(map[string]*global.ProtocolConfig),
		Registries:  make(map[string]*global.RegistryConfig),
	}
	for _, opt := range opts {
		opt(refOpts)
	}
	return refOpts
}

func decodeResolvedGenericServiceKey(t *testing.T, spec resolvedReferSpec) genericServiceKey {
	t.Helper()

	keyRaw, err := spec.cacheKey()
	require.NoError(t, err)

	var key genericServiceKey
	require.NoError(t, json.Unmarshal([]byte(keyRaw), &key))
	return key
}

func restorePropagator(t *testing.T, propagator propagation.TextMapPropagator) {
	t.Helper()

	original := otel.GetTextMapPropagator()
	otel.SetTextMapPropagator(propagator)
	t.Cleanup(func() {
		otel.SetTextMapPropagator(original)
	})
}
