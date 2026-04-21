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
	"testing"
	"time"
)

import (
	dubboConstant "dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/filter/generic"

	hessian "github.com/apache/dubbo-go-hessian2"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

import (
	cst "github.com/apache/dubbo-go-pixiu/pkg/common/constant"
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

func restorePropagator(t *testing.T, propagator propagation.TextMapPropagator) {
	t.Helper()

	original := otel.GetTextMapPropagator()
	otel.SetTextMapPropagator(propagator)
	t.Cleanup(func() {
		otel.SetTextMapPropagator(original)
	})
}

func cacheServiceForOutbound(t *testing.T, dc *Client, req *DubboOutboundRequest, service *generic.GenericService) {
	t.Helper()

	spec := dc.resolveFromOutbound(req)
	key, err := spec.cacheKey()
	require.NoError(t, err)
	dc.GenericServicePool[key] = service
}

func TestResolveFromOutboundDirectMode(t *testing.T) {
	dc := NewDubboClient()
	dc.SetConfig(&DubboProxyConfig{
		LoadBalance: "roundrobin",
		Retries:     "5",
		Timeout: &model.TimeoutConfig{
			RequestTimeoutStr: "4s",
		},
	})

	spec := dc.resolveFromOutbound(&DubboOutboundRequest{
		Service:       "com.example.UserService",
		Group:         "gray",
		Version:       "1.0.0",
		Address:       "127.0.0.1:20880",
		Protocol:      "dubbo",
		Serialization: "hessian2",
	})

	assert.Equal(t, "direct", spec.Mode)
	assert.Equal(t, "com.example.UserService", spec.Interface)
	assert.Equal(t, "gray", spec.Group)
	assert.Equal(t, "1.0.0", spec.Version)
	assert.Equal(t, "dubbo://127.0.0.1:20880", spec.URL)
	assert.Equal(t, "dubbo", spec.EffectiveProtocol)
	assert.Equal(t, "hessian2", spec.EffectiveSerialization)
	assert.Empty(t, spec.RegistryIDs)
	assert.False(t, spec.UseNacosWarmup)
	assert.Equal(t, "failover", spec.ConsumerDefaults.Cluster)
	assert.Equal(t, "roundrobin", spec.ConsumerDefaults.LoadBalance)
	assert.Equal(t, "5", spec.ConsumerDefaults.Retries)
	assert.Equal(t, 4*time.Second, spec.ConsumerDefaults.RequestTimeout)
}

func TestResolveFromOutboundRegistryMode(t *testing.T) {
	dc := NewDubboClient()
	dc.SetConfig(&DubboProxyConfig{
		Registries: map[string]model.Registry{
			"zk": {
				Protocol: "zookeeper",
				Address:  "127.0.0.1:2181",
			},
			"nacos-main": {
				Protocol: "nacos",
				Address:  "127.0.0.1:8848",
			},
		},
	})
	require.NoError(t, dc.Apply())

	spec := dc.resolveFromOutbound(&DubboOutboundRequest{
		Service:       "com.example.UserService",
		Group:         "gray",
		Version:       "1.0.0",
		Protocol:      "tri",
		Serialization: "protobuf",
	})

	assert.Equal(t, "registry", spec.Mode)
	assert.Equal(t, "com.example.UserService", spec.Interface)
	assert.Equal(t, []string{"nacos-main", "zk"}, spec.RegistryIDs)
	assert.True(t, spec.UseNacosWarmup)
	assert.Equal(t, "tri", spec.EffectiveProtocol)
	assert.Equal(t, "protobuf", spec.EffectiveSerialization)
	assert.Equal(t, "failover", spec.ConsumerDefaults.Cluster)
	assert.Equal(t, "3", spec.ConsumerDefaults.Retries)
	assert.Equal(t, cst.DefaultReqTimeout, spec.ConsumerDefaults.RequestTimeout)
}

func TestPreparePayloadRejectsLengthMismatch(t *testing.T) {
	dc := NewDubboClient()

	types, vals, finalValues, err := dc.preparePayload(&DubboOutboundRequest{
		Arguments:  []any{"only-one"},
		ParamTypes: []string{"java.lang.String", "int"},
	})

	assert.Nil(t, types)
	assert.Nil(t, vals)
	assert.Nil(t, finalValues)
	assert.EqualError(t, err, "arguments/paramTypes length mismatch: 1 vs 2")
}

func TestCacheKeyIncludesConsumerDefaults(t *testing.T) {
	baseSpec := resolvedReferSpec{
		Mode:                   "direct",
		Interface:              "com.example.UserService",
		URL:                    "dubbo://127.0.0.1:20880",
		EffectiveProtocol:      "dubbo",
		EffectiveSerialization: "hessian2",
		ConsumerDefaults: resolvedConsumerDefaults{
			Cluster:        "failover",
			LoadBalance:    "roundrobin",
			Retries:        "3",
			RequestTimeout: time.Second,
		},
	}

	changedSpec := baseSpec
	changedSpec.ConsumerDefaults.LoadBalance = "random"
	changedSpec.ConsumerDefaults.Retries = "5"
	changedSpec.ConsumerDefaults.RequestTimeout = 2 * time.Second

	baseKey, err := baseSpec.cacheKey()
	require.NoError(t, err)
	changedKey, err := changedSpec.cacheKey()
	require.NoError(t, err)

	assert.NotEqual(t, baseKey, changedKey)

	var decoded genericServiceKey
	require.NoError(t, json.Unmarshal([]byte(changedKey), &decoded))
	assert.Equal(t, "random", decoded.LoadBalance)
	assert.Equal(t, "5", decoded.Retries)
	assert.Equal(t, "2s", decoded.RequestTimeout)
}

func TestCallUsesOutboundOnly(t *testing.T) {
	dc := NewDubboClient()
	restorePropagator(t, fixedPropagator{})

	req := &DubboOutboundRequest{
		Service:       "com.example.UserService",
		Method:        "GetUser",
		Group:         "gray",
		Version:       "1.0.0",
		Address:       "127.0.0.1:20880",
		Protocol:      "dubbo",
		Serialization: "hessian2",
		Arguments:     []any{"user-1", 2},
		ParamTypes:    []string{"java.lang.String", "int"},
		Attachments: map[string]any{
			"user-key": "user-value",
		},
	}

	cacheServiceForOutbound(t, dc, req, &generic.GenericService{
		Invoke: func(ctx context.Context, methodName string, types []string, args []hessian.Object) (any, error) {
			assert.Equal(t, "GetUser", methodName)
			assert.Equal(t, []string{"java.lang.String", "int"}, types)
			require.Len(t, args, 2)
			assert.Equal(t, "user-1", args[0])
			assert.Equal(t, 2, args[1])

			attachments, ok := ctx.Value(dubboConstant.AttachmentKey).(map[string]any)
			require.True(t, ok)
			assert.Equal(t, "user-value", attachments["user-key"])

			return "ok", nil
		},
	})

	res, err := dc.Call(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, "ok", res)
}

func TestCallAppliesTimeout(t *testing.T) {
	dc := NewDubboClient()
	restorePropagator(t, fixedPropagator{})

	req := &DubboOutboundRequest{
		Service:       "com.example.UserService",
		Method:        "GetUser",
		Address:       "127.0.0.1:20880",
		Protocol:      "dubbo",
		Serialization: "hessian2",
		Arguments:     []any{"user-1"},
		ParamTypes:    []string{"java.lang.String"},
		Timeout:       80 * time.Millisecond,
	}

	cacheServiceForOutbound(t, dc, req, &generic.GenericService{
		Invoke: func(ctx context.Context, methodName string, types []string, args []hessian.Object) (any, error) {
			deadline, ok := ctx.Deadline()
			require.True(t, ok)
			assert.WithinDuration(t, time.Now().Add(req.Timeout), deadline, 150*time.Millisecond)
			return "ok", nil
		},
	})

	res, err := dc.Call(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, "ok", res)
}
