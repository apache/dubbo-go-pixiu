//go:build dubbo_legacy_call_tests
// +build dubbo_legacy_call_tests

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
	"net/http"
	"testing"
	"time"
)

import (
	dclient "dubbo.apache.org/dubbo-go/v3/client"
	dubboConstant "dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/filter/generic"
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

func (fixedPropagator) Inject(_ context.Context, _ propagation.TextMapCarrier) {
	// Intentionally left blank for tests that need a no-op injector.
}

func (fixedPropagator) Extract(ctx context.Context, _ propagation.TextMapCarrier) context.Context {
	return ctx
}

func (fixedPropagator) Fields() []string {
	return nil
}

type attachmentPropagator struct{}

const (
	invokePath                = "/invoke"
	traceparentKey            = "traceparent"
	traceparentValue          = "00-test-traceparent"
	otelScopeKey              = "x-otel-scope"
	otelScopeValue            = "pixiu"
	directDubboSerialization  = "hessian2"
	directTripleSerialization = "hessian2"
	helloRequestType          = "benchmark.HelloRequest"
	requestBodyValuesSource   = "requestBody.values"
	requestBodyTypesSource    = "requestBody.types"
	optValuesTarget           = "opt.values"
	optTypesTarget            = "opt.types"
	mappedAppName             = "mapped-app"
	demoAppName               = "demo-app"
	orderServiceInterface     = "com.example.OrderService"
	userServiceInterface      = "com.example.UserService"
	directDubboURL            = "dubbo://127.0.0.1:20880"
	directTripleURL           = "triple://127.0.0.1:50051"
	userAttachmentKey         = "user-key"
	userAttachmentValue       = "user-value"
	autoResolveRequestPath    = "/" + demoAppName + "/" + userServiceInterface + "/GetUser"
)

func (attachmentPropagator) Inject(_ context.Context, carrier propagation.TextMapCarrier) {
	carrier.Set(traceparentKey, traceparentValue)
	carrier.Set(otelScopeKey, otelScopeValue)
}

func (attachmentPropagator) Extract(ctx context.Context, _ propagation.TextMapCarrier) context.Context {
	return ctx
}

func (attachmentPropagator) Fields() []string {
	return []string{traceparentKey, otelScopeKey}
}

func TestCallDirectUsesConfiguredParameterTypesInsteadOfMappedTypes(t *testing.T) {
	dc := NewDubboClient()
	restorePropagator(t, fixedPropagator{})

	req := newDirectCallRequest(
		t,
		context.Background(),
		invokePath,
		invokePath,
		`{"values":["123"],"types":"int"}`,
		[]string{JavaStringClassName},
		directDubboSerialization,
		[]config.MappingParam{
			{Name: "headers.X-App", MapTo: "opt.application"},
			{Name: "headers.X-Interface", MapTo: "opt.interface"},
			{Name: "headers.X-Method", MapTo: "opt.method"},
			{Name: requestBodyValuesSource, MapTo: optValuesTarget},
			{Name: requestBodyTypesSource, MapTo: optTypesTarget},
		},
	)
	req.IngressRequest.Header.Set("X-App", mappedAppName)
	req.IngressRequest.Header.Set("X-Interface", orderServiceInterface)
	req.IngressRequest.Header.Set("X-Method", "CreateOrder")
	req.API.ApplicationName = "initial-app"
	req.API.Interface = "initial.Interface"
	req.API.Method.Method = "initialMethod"

	expectedIR := req.API.IntegrationRequest
	expectedIR.ApplicationName = mappedAppName
	expectedIR.Interface = orderServiceInterface
	expectedIR.Method = "CreateOrder"

	cached := false
	cacheServiceForRequest(t, dc, expectedIR, &generic.GenericService{
		Invoke: func(ctx context.Context, methodName string, types []string, args []hessian.Object) (any, error) {
			cached = true
			assert.Equal(t, "CreateOrder", methodName)
			assert.Equal(t, []string{JavaStringClassName}, types)
			require.Len(t, args, 1)
			assert.Equal(t, "123", args[0])
			return "ok", nil
		},
	})

	res, err := dc.Call(req)
	require.NoError(t, err)
	assert.Equal(t, "ok", res)
	assert.True(t, cached)
	assert.Equal(t, mappedAppName, req.API.ApplicationName)
	assert.Equal(t, orderServiceInterface, req.API.Interface)
	assert.Equal(t, "CreateOrder", req.API.Method.Method)
}

func TestCallAutoResolveSnapshotUsesMappedFields(t *testing.T) {
	dc := NewDubboClient()
	restorePropagator(t, fixedPropagator{})

	req := newDirectCallRequest(
		t,
		context.Background(),
		"/:application/:interface/:method",
		autoResolveRequestPath,
		`{"values":["u-1"]}`,
		[]string{JavaStringClassName},
		directDubboSerialization,
		[]config.MappingParam{
			{Name: requestBodyValuesSource, MapTo: optValuesTarget},
			{Name: "uri.application", MapTo: "opt.application"},
			{Name: "uri.interface", MapTo: "opt.interface"},
			{Name: "uri.method", MapTo: "opt.method"},
			{Name: "headers.X-Group", MapTo: "opt.group"},
			{Name: "headers.X-Version", MapTo: "opt.version"},
		},
	)
	req.IngressRequest.Header.Set("X-Group", "gray")
	req.IngressRequest.Header.Set("X-Version", "1.0.0")

	expectedIR := req.API.IntegrationRequest
	expectedIR.ApplicationName = demoAppName
	expectedIR.Interface = userServiceInterface
	expectedIR.Method = "GetUser"
	expectedIR.Group = "gray"
	expectedIR.Version = "1.0.0"

	cacheServiceForRequest(t, dc, expectedIR, &generic.GenericService{
		Invoke: func(ctx context.Context, methodName string, types []string, args []hessian.Object) (any, error) {
			assert.Equal(t, "GetUser", methodName)
			assert.Equal(t, []string{JavaStringClassName}, types)
			return "resolved", nil
		},
	})

	res, err := dc.Call(req)
	require.NoError(t, err)
	assert.Equal(t, "resolved", res)
	assert.Equal(t, demoAppName, req.API.ApplicationName)
	assert.Equal(t, userServiceInterface, req.API.Interface)
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
		Mode:                   "direct",
		URL:                    directDubboURL,
		Interface:              userServiceInterface,
		EffectiveProtocol:      "dubbo",
		EffectiveSerialization: directDubboSerialization,
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
			URL: "  " + directTripleURL + "  ",
		},
		DubboBackendConfig: config.DubboBackendConfig{
			ApplicationName: demoAppName,
			Interface:       userServiceInterface,
			Group:           "gray",
			Version:         "1.0.0",
			Serialization:   directTripleSerialization,
		},
	})
	require.NoError(t, err)
	assert.Equal(t, "direct", spec.Mode)
	assert.Equal(t, directTripleURL, spec.URL)
	assert.Equal(t, "tri", spec.EffectiveProtocol)
	assert.Equal(t, directTripleSerialization, spec.EffectiveSerialization)
	assert.Equal(t, userServiceInterface, spec.Interface)
	assert.Equal(t, "gray", spec.Group)
	assert.Equal(t, "1.0.0", spec.Version)
}

func TestResolveReferSpecURLSelectsDirectEvenWithRegistries(t *testing.T) {
	dc := NewDubboClient()
	dc.registries = map[string]*global.RegistryConfig{
		"zk": {Protocol: "zookeeper"},
	}

	spec, err := dc.resolveReferSpec(config.IntegrationRequest{
		HTTPBackendConfig: config.HTTPBackendConfig{
			URL: "  " + directDubboURL + "  ",
		},
		DubboBackendConfig: config.DubboBackendConfig{
			Interface:     userServiceInterface,
			Protocol:      "dubbo",
			Serialization: directDubboSerialization,
		},
	})
	require.NoError(t, err)
	assert.Equal(t, "direct", spec.Mode)
	assert.Equal(t, directDubboURL, spec.URL)
	assert.Empty(t, spec.RegistryIDs)
	assert.False(t, spec.UseNacosWarmup)
}

func TestResolveReferSpecRejectsDirectProtocolSchemeMismatch(t *testing.T) {
	dc := NewDubboClient()

	_, err := dc.resolveReferSpec(config.IntegrationRequest{
		HTTPBackendConfig: config.HTTPBackendConfig{
			URL: directDubboURL,
		},
		DubboBackendConfig: config.DubboBackendConfig{
			Protocol:      "tri",
			Serialization: directDubboSerialization,
		},
	})
	assert.EqualError(t, err, "direct protocol mismatch: url=dubbo protocol=tri")
}

func TestResolvedReferSpecCacheKeyIncludesEffectiveProtocol(t *testing.T) {
	dubboSpec := resolvedReferSpec{
		Mode:              "registry",
		Interface:         userServiceInterface,
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
		Mode:                   "direct",
		Interface:              userServiceInterface,
		URL:                    directTripleURL,
		EffectiveProtocol:      "tri",
		EffectiveSerialization: directTripleSerialization,
		ConsumerDefaults: resolvedConsumerDefaults{
			Cluster: "failover",
		},
	}

	refOpts := applyResolvedReferenceOptions(t, dc, spec)
	assert.Equal(t, directTripleURL, refOpts.Reference.URL)
	assert.Equal(t, "tri", refOpts.Reference.Protocol)
	assert.Equal(t, directTripleSerialization, refOpts.Reference.Serialization)
	assert.Empty(t, refOpts.Reference.RegistryIDs)
}

func TestDirectModeUsesCanonicalURLForOptionsAndCacheKey(t *testing.T) {
	dc := NewDubboClient()

	spec, err := dc.resolveReferSpec(config.IntegrationRequest{
		HTTPBackendConfig: config.HTTPBackendConfig{
			URL: "  " + directDubboURL + "  ",
		},
		DubboBackendConfig: config.DubboBackendConfig{
			ClusterName:     "demo-cluster",
			ApplicationName: demoAppName,
			Interface:       userServiceInterface,
			Group:           "gray",
			Version:         "1.0.0",
			Serialization:   directDubboSerialization,
		},
	})
	require.NoError(t, err)
	assert.Equal(t, "direct", spec.Mode)
	assert.Equal(t, directDubboURL, spec.URL)
	assert.Equal(t, "dubbo", spec.EffectiveProtocol)

	refOpts := applyResolvedReferenceOptions(t, dc, spec)
	assert.Equal(t, directDubboURL, refOpts.Reference.URL)
	assert.Empty(t, refOpts.Reference.RegistryIDs)
	assert.Equal(t, dubboConstant.DubboProtocol, refOpts.Reference.Protocol)
	assert.Equal(t, directDubboSerialization, refOpts.Reference.Serialization)

	key := decodeResolvedGenericServiceKey(t, spec)
	assert.Equal(t, directDubboURL, key.URL)
	assert.Equal(t, "direct", key.Mode)
	assert.Equal(t, directDubboSerialization, key.Serialization)
}

func TestDirectCacheKeyIncludesSerialization(t *testing.T) {
	base := resolvedReferSpec{
		Mode:                   "direct",
		Interface:              userServiceInterface,
		URL:                    directTripleURL,
		EffectiveProtocol:      "tri",
		EffectiveSerialization: directTripleSerialization,
		ConsumerDefaults: resolvedConsumerDefaults{
			Cluster: "failover",
		},
	}
	other := base
	other.EffectiveSerialization = "json"

	baseKeyRaw, err := base.cacheKey()
	require.NoError(t, err)
	otherKeyRaw, err := other.cacheKey()
	require.NoError(t, err)
	assert.NotEqual(t, baseKeyRaw, otherKeyRaw)

	baseKey := decodeResolvedGenericServiceKey(t, base)
	otherKey := decodeResolvedGenericServiceKey(t, other)
	assert.Equal(t, directTripleSerialization, baseKey.Serialization)
	assert.Equal(t, "json", otherKey.Serialization)
}

func TestResolveReferSpecRegistryModeUsesDeclaredProtocolWithoutProbe(t *testing.T) {
	dc := NewDubboClient()
	dc.registries = map[string]*global.RegistryConfig{
		"zk": {Protocol: "zookeeper"},
	}

	spec, err := dc.resolveReferSpec(config.IntegrationRequest{
		RequestType: "dubbo",
		DubboBackendConfig: config.DubboBackendConfig{
			Interface: userServiceInterface,
			Protocol:  "dubbo",
		},
	})
	require.NoError(t, err)
	assert.Equal(t, "registry", spec.Mode)
	assert.Equal(t, []string{"zk"}, spec.RegistryIDs)
	assert.Equal(t, "dubbo", spec.EffectiveProtocol)
}

func TestResolveReferSpecRegistryModeUsesTripleRequestTypeWithoutProbe(t *testing.T) {
	dc := NewDubboClient()
	dc.registries = map[string]*global.RegistryConfig{
		"zk": {Protocol: "zookeeper"},
	}

	spec, err := dc.resolveReferSpec(config.IntegrationRequest{
		RequestType: "triple",
		DubboBackendConfig: config.DubboBackendConfig{
			Interface: userServiceInterface,
		},
	})
	require.NoError(t, err)
	assert.Equal(t, "tri", spec.EffectiveProtocol)
}

func TestBuildReferenceOptionsUsesRegistryProtocolFromRequest(t *testing.T) {
	dc := NewDubboClient()

	spec := resolvedReferSpec{
		Mode:              "registry",
		Interface:         userServiceInterface,
		RegistryIDs:       []string{"zk"},
		EffectiveProtocol: "dubbo",
		ConsumerDefaults: resolvedConsumerDefaults{
			Cluster:        "failover",
			Retries:        "7",
			LoadBalance:    "consistenthash",
			RequestTimeout: 6 * time.Second,
		},
	}

	refOpts := applyResolvedReferenceOptions(t, dc, spec)
	assert.Equal(t, []string{"zk"}, refOpts.Reference.RegistryIDs)
	assert.Equal(t, dubboConstant.DubboProtocol, refOpts.Reference.Protocol)
	assert.Equal(t, dubboConstant.LoadBalanceKeyConsistentHashing, refOpts.Reference.Loadbalance)
	assert.Equal(t, "7", refOpts.Reference.Retries)
	assert.Equal(t, "6s", refOpts.Reference.RequestTimeout)
	assert.Equal(t, "true", refOpts.Reference.Generic)
}

func TestCallWithZeroTimeoutDoesNotCreateImmediateDeadline(t *testing.T) {
	dc := NewDubboClient()
	restorePropagator(t, fixedPropagator{})

	req := newDirectCallRequest(
		t,
		context.Background(),
		invokePath,
		invokePath,
		`{"values":["1"]}`,
		[]string{JavaStringClassName},
		directDubboSerialization,
		[]config.MappingParam{
			{Name: requestBodyValuesSource, MapTo: optValuesTarget},
		},
	)
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

	req := newDirectCallRequest(
		t,
		context.Background(),
		invokePath,
		invokePath,
		`{"values":["1"]}`,
		[]string{JavaStringClassName},
		directDubboSerialization,
		[]config.MappingParam{
			{Name: requestBodyValuesSource, MapTo: optValuesTarget},
		},
	)
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
	dc := NewDubboClient()
	dc.SetConfig(&DubboProxyConfig{
		Timeout: &model.TimeoutConfig{
			RequestTimeoutStr: "6s",
		},
		LoadBalance: "consistenthash",
		Retries:     "7",
	})

	spec, err := dc.resolveReferSpec(config.IntegrationRequest{
		HTTPBackendConfig: config.HTTPBackendConfig{
			URL: directDubboURL,
		},
		DubboBackendConfig: config.DubboBackendConfig{
			Interface:     userServiceInterface,
			Retries:       "2",
			Serialization: directDubboSerialization,
		},
	})
	require.NoError(t, err)

	refOpts := applyResolvedReferenceOptions(t, dc, spec)
	assert.Equal(t, dubboConstant.ClusterKeyFailover, refOpts.Reference.Cluster)
	assert.Equal(t, dubboConstant.DubboProtocol, refOpts.Reference.Protocol)
	assert.Equal(t, dubboConstant.LoadBalanceKeyConsistentHashing, refOpts.Reference.Loadbalance)
	assert.Equal(t, "7", refOpts.Reference.Retries)
	assert.Equal(t, "6s", refOpts.Reference.RequestTimeout)
	assert.Equal(t, "true", refOpts.Reference.Generic)
	assert.Equal(t, directDubboSerialization, refOpts.Reference.Serialization)
}

func TestCallFlattensAttachmentsIntoMapAny(t *testing.T) {
	dc := NewDubboClient()
	restorePropagator(t, attachmentPropagator{})

	baseCtx := context.WithValue(context.Background(), dubboConstant.AttachmentKey, map[string]any{
		userAttachmentKey: userAttachmentValue,
	})
	req := newDirectCallRequest(
		t,
		baseCtx,
		invokePath,
		invokePath,
		`{"values":["1"]}`,
		[]string{JavaStringClassName},
		directDubboSerialization,
		[]config.MappingParam{
			{Name: requestBodyValuesSource, MapTo: optValuesTarget},
		},
	)

	cacheServiceForRequest(t, dc, req.API.IntegrationRequest, &generic.GenericService{
		Invoke: func(ctx context.Context, methodName string, types []string, args []hessian.Object) (any, error) {
			attachments, ok := ctx.Value(dubboConstant.AttachmentKey).(map[string]any)
			require.True(t, ok)
			assert.Equal(t, userAttachmentValue, attachments[userAttachmentKey])
			assert.Equal(t, traceparentValue, attachments[traceparentKey])
			assert.Equal(t, otelScopeValue, attachments[otelScopeKey])
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
		userAttachmentKey: userAttachmentValue,
	})
	req := newDirectCallRequest(
		t,
		baseCtx,
		invokePath,
		invokePath,
		`{"values":["1"]}`,
		[]string{JavaStringClassName},
		directDubboSerialization,
		[]config.MappingParam{
			{Name: requestBodyValuesSource, MapTo: optValuesTarget},
		},
	)

	cacheServiceForRequest(t, dc, req.API.IntegrationRequest, &generic.GenericService{
		Invoke: func(ctx context.Context, methodName string, types []string, args []hessian.Object) (any, error) {
			attachments, ok := ctx.Value(dubboConstant.AttachmentKey).(map[string]any)
			require.True(t, ok)
			assert.Equal(t, userAttachmentValue, attachments[userAttachmentKey])
			assert.Equal(t, traceparentValue, attachments[traceparentKey])
			assert.Equal(t, otelScopeValue, attachments[otelScopeKey])
			return "ok", nil
		},
	})

	res, err := dc.Call(req)
	require.NoError(t, err)
	assert.Equal(t, "ok", res)
}

func TestCallDirectRejectsMissingSerialization(t *testing.T) {
	dc := NewDubboClient()
	restorePropagator(t, fixedPropagator{})

	req := newDirectCallRequest(
		t,
		context.Background(),
		invokePath,
		invokePath,
		`{"values":["1"]}`,
		[]string{JavaStringClassName},
		"",
		[]config.MappingParam{
			{Name: requestBodyValuesSource, MapTo: optValuesTarget},
		},
	)
	res, err := dc.Call(req)
	assert.Nil(t, res)
	assert.EqualError(t, err, "direct generic invoke requires serialization")
}

func TestCallDirectRejectsMissingParameterTypes(t *testing.T) {
	dc := NewDubboClient()
	restorePropagator(t, fixedPropagator{})

	req := newDirectCallRequest(
		t,
		context.Background(),
		invokePath,
		invokePath,
		`{"values":["1"],"types":"java.lang.String"}`,
		nil,
		directDubboSerialization,
		[]config.MappingParam{
			{Name: requestBodyValuesSource, MapTo: optValuesTarget},
			{Name: requestBodyTypesSource, MapTo: optTypesTarget},
		},
	)
	res, err := dc.Call(req)
	assert.Nil(t, res)
	assert.EqualError(t, err, "direct generic invoke requires parameterTypes")
}

func TestCallDirectConvertsMappedValuesUsingConfiguredParameterTypes(t *testing.T) {
	dc := NewDubboClient()
	restorePropagator(t, fixedPropagator{})

	req := newDirectCallRequest(
		t,
		context.Background(),
		invokePath,
		invokePath,
		`{"values":["123"]}`,
		[]string{"int"},
		directDubboSerialization,
		[]config.MappingParam{
			{Name: requestBodyValuesSource, MapTo: optValuesTarget},
		},
	)

	cacheServiceForRequest(t, dc, req.API.IntegrationRequest, &generic.GenericService{
		Invoke: func(ctx context.Context, methodName string, types []string, args []hessian.Object) (any, error) {
			assert.Equal(t, []string{"int"}, types)
			require.Len(t, args, 1)
			assert.Equal(t, 123, args[0])
			return "ok", nil
		},
	})

	res, err := dc.Call(req)
	require.NoError(t, err)
	assert.Equal(t, "ok", res)
}

func TestCallDirectPassesThroughComplexDeclaredTypes(t *testing.T) {
	dc := NewDubboClient()
	restorePropagator(t, fixedPropagator{})

	req := newDirectCallRequest(
		t,
		context.Background(),
		invokePath,
		invokePath,
		`{"values":[{"name":"alice"}]}`,
		[]string{helloRequestType},
		directTripleSerialization,
		[]config.MappingParam{
			{Name: requestBodyValuesSource, MapTo: optValuesTarget},
		},
	)
	req.API.IntegrationRequest.URL = directTripleURL

	cacheServiceForRequest(t, dc, req.API.IntegrationRequest, &generic.GenericService{
		Invoke: func(ctx context.Context, methodName string, types []string, args []hessian.Object) (any, error) {
			assert.Equal(t, []string{helloRequestType}, types)
			require.Len(t, args, 1)
			require.IsType(t, map[string]any{}, args[0])
			assert.Equal(t, map[string]any{"name": "alice"}, args[0])
			return "ok", nil
		},
	})

	res, err := dc.Call(req)
	require.NoError(t, err)
	assert.Equal(t, "ok", res)
}

func TestCallDirectSupportsExplicitZeroParameterMethod(t *testing.T) {
	dc := NewDubboClient()
	restorePropagator(t, fixedPropagator{})

	req := newDirectCallRequest(
		t,
		context.Background(),
		invokePath,
		invokePath,
		`{}`,
		[]string{},
		directDubboSerialization,
		nil,
	)
	req.API.Method.Method = "Ping"

	cacheServiceForRequest(t, dc, req.API.IntegrationRequest, &generic.GenericService{
		Invoke: func(ctx context.Context, methodName string, types []string, args []hessian.Object) (any, error) {
			assert.Equal(t, "Ping", methodName)
			assert.Empty(t, types)
			assert.Empty(t, args)
			return "pong", nil
		},
	})

	res, err := dc.Call(req)
	require.NoError(t, err)
	assert.Equal(t, "pong", res)
}

func newDirectCallRequest(
	t *testing.T,
	ctx context.Context,
	urlPattern string,
	rawURL string,
	body string,
	parameterTypes []string,
	serialization string,
	mapping []config.MappingParam,
) *client.Request {
	t.Helper()

	request, err := http.NewRequest(http.MethodPost, rawURL, bytes.NewBufferString(body))
	require.NoError(t, err)

	api := mock.GetMockAPI(http.MethodPost, urlPattern)
	api.MappingParams = mapping
	api.IntegrationRequest.URL = "  " + directDubboURL + "  "
	api.IntegrationRequest.ParameterTypes = cloneParameterTypes(parameterTypes)
	api.IntegrationRequest.Serialization = serialization
	api.Interface = userServiceInterface
	api.Method.Method = "SayHello"
	api.ApplicationName = demoAppName
	api.Group = "demo-group"
	api.Version = "1.0.0"

	return client.NewReq(ctx, request, api)
}

func cloneParameterTypes(parameterTypes []string) []string {
	if parameterTypes == nil {
		return nil
	}

	cloned := make([]string, len(parameterTypes))
	copy(cloned, parameterTypes)
	return cloned
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
