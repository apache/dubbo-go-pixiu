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

package remote

import (
	"bytes"
	"net/http"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

import (
	clientdubbo "github.com/apache/dubbo-go-pixiu/pkg/client/dubbo"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/config"
	"github.com/apache/dubbo-go-pixiu/pkg/router"
)

func newTestAPI(ir config.IntegrationRequest, pattern string) router.API {
	if pattern == "" {
		pattern = "/users/:id"
	}
	if ir.RequestType == "" {
		ir.RequestType = constant.DubboRequest
	}
	return router.API{
		URLPattern: pattern,
		Method: config.Method{
			Enable:             true,
			HTTPVerb:           http.MethodPost,
			IntegrationRequest: ir,
		},
	}
}

func TestBuildOutboundMapsQueryHeaderBodyAndURI(t *testing.T) {
	handler := &DubboHandler{}
	api := newTestAPI(config.IntegrationRequest{
		RequestType: constant.DubboRequest,
		DubboBackendConfig: config.DubboBackendConfig{
			Interface: "com.demo.UserService",
			Method:    "SayHello",
			Group:     "demo-group",
			Version:   "1.0.0",
		},
	}, "/users/:id")
	api.MappingParams = []config.MappingParam{
		{Name: "queryStrings.page", MapTo: "0", MapType: "java.lang.Integer"},
		{Name: "headers.x-user", MapTo: "1", MapType: clientdubbo.JavaStringClassName},
		{Name: "requestBody.profile.age", MapTo: "2", MapType: "java.lang.Integer"},
		{Name: "uri.id", MapTo: "3"},
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"http://example.com/users/42?page=7",
		bytes.NewBufferString(`{"profile":{"age":"18"}}`),
	)
	require.NoError(t, err)
	req.Header.Set("x-user", "alice")

	outbound, err := handler.BuildOutbound(req, api)
	require.NoError(t, err)
	require.NotNil(t, outbound)
	assert.Equal(t, "com.demo.UserService", outbound.Service)
	assert.Equal(t, "SayHello", outbound.Method)
	assert.Equal(t, "demo-group", outbound.Group)
	assert.Equal(t, "1.0.0", outbound.Version)
	assert.Equal(t, "dubbo", outbound.Protocol)
	assert.Empty(t, outbound.Serialization)
	assert.Equal(t, []any{7, "alice", 18, "42"}, outbound.Arguments)
	assert.Equal(t, []string{
		"java.lang.Integer",
		clientdubbo.JavaStringClassName,
		"java.lang.Integer",
		clientdubbo.JavaStringClassName,
	}, outbound.ParamTypes)
}

func TestBuildOutboundParameterTypesPriority(t *testing.T) {
	handler := &DubboHandler{}

	t.Run("integration request parameter types override opt types", func(t *testing.T) {
		api := newTestAPI(config.IntegrationRequest{
			RequestType: constant.DubboRequest,
			DubboBackendConfig: config.DubboBackendConfig{
				Interface:      "com.demo.UserService",
				Method:         "SayHello",
				Group:          "demo-group",
				Version:        "1.0.0",
				ParameterTypes: []string{"java.lang.Integer"},
			},
			MappingParams: []config.MappingParam{
				{Name: "requestBody.values", MapTo: "opt.values"},
				{Name: "requestBody.types", MapTo: "opt.types"},
			},
		}, "/users/:id")

		req, err := http.NewRequest(
			http.MethodPost,
			"http://example.com/users/42",
			bytes.NewBufferString(`{"values":["7"],"types":"java.lang.String"}`),
		)
		require.NoError(t, err)

		outbound, err := handler.BuildOutbound(req, api)
		require.NoError(t, err)
		assert.Equal(t, []string{"java.lang.Integer"}, outbound.ParamTypes)
		assert.Equal(t, []any{7}, outbound.Arguments)
	})

	t.Run("opt types override inferred types", func(t *testing.T) {
		api := newTestAPI(config.IntegrationRequest{
			RequestType: constant.DubboRequest,
			DubboBackendConfig: config.DubboBackendConfig{
				Interface: "com.demo.UserService",
				Method:    "SayHello",
				Group:     "demo-group",
				Version:   "1.0.0",
			},
			MappingParams: []config.MappingParam{
				{Name: "requestBody.values", MapTo: "opt.values"},
				{Name: "requestBody.types", MapTo: "opt.types"},
			},
		}, "/users/:id")

		req, err := http.NewRequest(
			http.MethodPost,
			"http://example.com/users/42",
			bytes.NewBufferString(`{"values":["7"],"types":"java.lang.Integer"}`),
		)
		require.NoError(t, err)

		outbound, err := handler.BuildOutbound(req, api)
		require.NoError(t, err)
		assert.Equal(t, []string{"java.lang.Integer"}, outbound.ParamTypes)
		assert.Equal(t, []any{7}, outbound.Arguments)
	})

	t.Run("infer parameter types when none declared", func(t *testing.T) {
		api := newTestAPI(config.IntegrationRequest{
			RequestType: constant.DubboRequest,
			DubboBackendConfig: config.DubboBackendConfig{
				Interface: "com.demo.UserService",
				Method:    "SayHello",
				Group:     "demo-group",
				Version:   "1.0.0",
			},
			MappingParams: []config.MappingParam{
				{Name: "queryStrings.name", MapTo: "0"},
			},
		}, "/users/:id")

		req, err := http.NewRequest(
			http.MethodPost,
			"http://example.com/users/42?name=alice",
			bytes.NewBufferString(`{}`),
		)
		require.NoError(t, err)

		outbound, err := handler.BuildOutbound(req, api)
		require.NoError(t, err)
		assert.Equal(t, []string{clientdubbo.JavaStringClassName}, outbound.ParamTypes)
		assert.Equal(t, []any{"alice"}, outbound.Arguments)
	})
}

func TestBuildOutboundRejectsMixedPositionalAndOptValues(t *testing.T) {
	handler := &DubboHandler{}
	api := newTestAPI(config.IntegrationRequest{
		RequestType: constant.DubboRequest,
		DubboBackendConfig: config.DubboBackendConfig{
			Interface: "com.demo.UserService",
			Method:    "SayHello",
			Group:     "demo-group",
			Version:   "1.0.0",
		},
		MappingParams: []config.MappingParam{
			{Name: "queryStrings.name", MapTo: "0"},
			{Name: "requestBody.values", MapTo: "opt.values"},
		},
	}, "/users/:id")

	req, err := http.NewRequest(
		http.MethodPost,
		"http://example.com/users/42?name=alice",
		bytes.NewBufferString(`{"values":["bob"]}`),
	)
	require.NoError(t, err)

	outbound, err := handler.BuildOutbound(req, api)
	assert.Nil(t, outbound)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "mutually exclusive")
}

func TestBuildOutboundRejectsUnknownOptMapping(t *testing.T) {
	handler := &DubboHandler{}
	api := newTestAPI(config.IntegrationRequest{
		RequestType: constant.DubboRequest,
		DubboBackendConfig: config.DubboBackendConfig{
			Interface: "com.demo.UserService",
			Method:    "SayHello",
			Group:     "demo-group",
			Version:   "1.0.0",
		},
		MappingParams: []config.MappingParam{
			{Name: "queryStrings.name", MapTo: "opt.unknown"},
		},
	}, "/users/:id")

	req, err := http.NewRequest(
		http.MethodPost,
		"http://example.com/users/42?name=alice",
		bytes.NewBufferString(`{}`),
	)
	require.NoError(t, err)

	outbound, err := handler.BuildOutbound(req, api)
	assert.Nil(t, outbound)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown opt mapping")
}

func TestBuildOutboundRejectsDeprecatedOptApplication(t *testing.T) {
	handler := &DubboHandler{}
	api := newTestAPI(config.IntegrationRequest{
		RequestType: constant.DubboRequest,
		DubboBackendConfig: config.DubboBackendConfig{
			Interface: "com.demo.UserService",
			Method:    "SayHello",
			Group:     "demo-group",
			Version:   "1.0.0",
		},
		MappingParams: []config.MappingParam{
			{Name: "queryStrings.app", MapTo: "opt.application"},
		},
	}, "/users/:id")

	req, err := http.NewRequest(
		http.MethodPost,
		"http://example.com/users/42?app=demo",
		bytes.NewBufferString(`{}`),
	)
	require.NoError(t, err)

	outbound, err := handler.BuildOutbound(req, api)
	assert.Nil(t, outbound)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "deprecated opt mapping")
}

func TestBuildOutboundRejectsNegativePositionalMapping(t *testing.T) {
	handler := &DubboHandler{}
	api := newTestAPI(config.IntegrationRequest{
		RequestType: constant.DubboRequest,
		DubboBackendConfig: config.DubboBackendConfig{
			Interface: "com.demo.UserService",
			Method:    "SayHello",
		},
		MappingParams: []config.MappingParam{
			{Name: "queryStrings.name", MapTo: "-1"},
		},
	}, "/users/:id")

	req, err := http.NewRequest(
		http.MethodPost,
		"http://example.com/users/42?name=alice",
		bytes.NewBufferString(`{}`),
	)
	require.NoError(t, err)

	outbound, err := handler.BuildOutbound(req, api)
	assert.Nil(t, outbound)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "Parameter mapping")
	assert.Contains(t, err.Error(), "incorrect")
}

func TestBuildOutboundRejectsProtocolSchemeMismatch(t *testing.T) {
	handler := &DubboHandler{}
	api := newTestAPI(config.IntegrationRequest{
		RequestType: constant.DubboRequest,
		DubboBackendConfig: config.DubboBackendConfig{
			Interface:     "com.demo.UserService",
			Method:        "SayHello",
			Group:         "demo-group",
			Version:       "1.0.0",
			Protocol:      "tri",
			Serialization: "hessian2",
		},
		HTTPBackendConfig: config.HTTPBackendConfig{
			URL: "dubbo://127.0.0.1:20880",
		},
	}, "/users/:id")

	req, err := http.NewRequest(
		http.MethodPost,
		"http://example.com/users/42",
		bytes.NewBufferString(`{}`),
	)
	require.NoError(t, err)

	outbound, err := handler.BuildOutbound(req, api)
	assert.Nil(t, outbound)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "direct protocol mismatch")
}

func TestBuildOutboundRejectsDirectURLWithoutSerialization(t *testing.T) {
	handler := &DubboHandler{}
	api := newTestAPI(config.IntegrationRequest{
		RequestType: constant.DubboRequest,
		DubboBackendConfig: config.DubboBackendConfig{
			Interface: "com.demo.UserService",
			Method:    "SayHello",
			Group:     "demo-group",
			Version:   "1.0.0",
		},
		HTTPBackendConfig: config.HTTPBackendConfig{
			URL: "dubbo://127.0.0.1:20880",
		},
	}, "/users/:id")

	req, err := http.NewRequest(
		http.MethodPost,
		"http://example.com/users/42",
		bytes.NewBufferString(`{}`),
	)
	require.NoError(t, err)

	outbound, err := handler.BuildOutbound(req, api)
	assert.Nil(t, outbound)
	require.Error(t, err)
	assert.EqualError(t, err, "direct generic invoke requires serialization")
}

func TestBuildOutboundRejectsDirectURLWithoutParameterTypes(t *testing.T) {
	handler := &DubboHandler{}
	api := newTestAPI(config.IntegrationRequest{
		RequestType: constant.DubboRequest,
		DubboBackendConfig: config.DubboBackendConfig{
			Interface:     "com.demo.UserService",
			Method:        "SayHello",
			Serialization: "hessian2",
		},
		HTTPBackendConfig: config.HTTPBackendConfig{
			URL: "dubbo://127.0.0.1:20880",
		},
		MappingParams: []config.MappingParam{
			{Name: "requestBody.values", MapTo: "opt.values"},
			{Name: "requestBody.types", MapTo: "opt.types"},
		},
	}, "/users/:id")

	req, err := http.NewRequest(
		http.MethodPost,
		"http://example.com/users/42",
		bytes.NewBufferString(`{"values":["alice"],"types":"java.lang.String"}`),
	)
	require.NoError(t, err)

	outbound, err := handler.BuildOutbound(req, api)
	assert.Nil(t, outbound)
	require.Error(t, err)
	assert.EqualError(t, err, "direct generic invoke requires parameterTypes")
}

func TestBuildOutboundResolvesProtocolPriority(t *testing.T) {
	handler := &DubboHandler{}

	t.Run("integration request protocol wins", func(t *testing.T) {
		api := newTestAPI(config.IntegrationRequest{
			RequestType: constant.DubboRequest,
			DubboBackendConfig: config.DubboBackendConfig{
				Interface: "com.demo.UserService",
				Method:    "SayHello",
				Group:     "demo-group",
				Version:   "1.0.0",
				Protocol:  "triple",
			},
		}, "/users/:id")

		req, err := http.NewRequest(http.MethodPost, "http://example.com/users/42", bytes.NewBufferString(`{}`))
		require.NoError(t, err)

		outbound, err := handler.BuildOutbound(req, api)
		require.NoError(t, err)
		assert.Equal(t, "tri", outbound.Protocol)
	})

	t.Run("request type wins when protocol empty", func(t *testing.T) {
		api := newTestAPI(config.IntegrationRequest{
			RequestType: "triple",
			DubboBackendConfig: config.DubboBackendConfig{
				Interface: "com.demo.UserService",
				Method:    "SayHello",
				Group:     "demo-group",
				Version:   "1.0.0",
			},
		}, "/users/:id")

		req, err := http.NewRequest(http.MethodPost, "http://example.com/users/42", bytes.NewBufferString(`{}`))
		require.NoError(t, err)

		outbound, err := handler.BuildOutbound(req, api)
		require.NoError(t, err)
		assert.Equal(t, "tri", outbound.Protocol)
	})

	t.Run("default to dubbo", func(t *testing.T) {
		api := newTestAPI(config.IntegrationRequest{
			DubboBackendConfig: config.DubboBackendConfig{
				Interface: "com.demo.UserService",
				Method:    "SayHello",
				Group:     "demo-group",
				Version:   "1.0.0",
			},
		}, "/users/:id")

		req, err := http.NewRequest(http.MethodPost, "http://example.com/users/42", bytes.NewBufferString(`{}`))
		require.NoError(t, err)

		outbound, err := handler.BuildOutbound(req, api)
		require.NoError(t, err)
		assert.Equal(t, "dubbo", outbound.Protocol)
	})

	t.Run("direct url scheme wins when protocol empty", func(t *testing.T) {
		api := newTestAPI(config.IntegrationRequest{
			RequestType: constant.DubboRequest,
			DubboBackendConfig: config.DubboBackendConfig{
				Interface:      "com.demo.UserService",
				Method:         "SayHello",
				ParameterTypes: []string{},
				Serialization:  "hessian2",
			},
			HTTPBackendConfig: config.HTTPBackendConfig{
				URL: "tri://127.0.0.1:50051",
			},
		}, "/users/:id")

		req, err := http.NewRequest(http.MethodPost, "http://example.com/users/42", bytes.NewBufferString(`{}`))
		require.NoError(t, err)

		outbound, err := handler.BuildOutbound(req, api)
		require.NoError(t, err)
		assert.Equal(t, "tri", outbound.Protocol)
		assert.Empty(t, outbound.Arguments)
		assert.Empty(t, outbound.ParamTypes)
	})
}

func TestBuildOutboundValidatesDirectInvokeArity(t *testing.T) {
	handler := &DubboHandler{}
	api := newTestAPI(config.IntegrationRequest{
		RequestType: constant.DubboRequest,
		DubboBackendConfig: config.DubboBackendConfig{
			Interface:      "com.demo.UserService",
			Method:         "SayHello",
			Group:          "demo-group",
			Version:        "1.0.0",
			ParameterTypes: []string{"java.lang.Integer", clientdubbo.JavaStringClassName},
			Serialization:  "hessian2",
		},
		HTTPBackendConfig: config.HTTPBackendConfig{
			URL: "dubbo://127.0.0.1:20880",
		},
		MappingParams: []config.MappingParam{
			{Name: "requestBody.values", MapTo: "opt.values"},
		},
	}, "/users/:id")

	req, err := http.NewRequest(
		http.MethodPost,
		"http://example.com/users/42",
		bytes.NewBufferString(`{"values":["7"]}`),
	)
	require.NoError(t, err)

	outbound, err := handler.BuildOutbound(req, api)
	assert.Nil(t, outbound)
	require.Error(t, err)
	assert.EqualError(t, err, "direct generic invoke requires values to match parameterTypes")
}

func TestBuildOutboundPreservesOptValuesInlineTypeCoercion(t *testing.T) {
	handler := &DubboHandler{}
	api := newTestAPI(config.IntegrationRequest{
		RequestType: constant.DubboRequest,
		DubboBackendConfig: config.DubboBackendConfig{
			Interface: "com.demo.UserService",
			Method:    "SayHello",
			Group:     "demo-group",
			Version:   "1.0.0",
		},
		MappingParams: []config.MappingParam{
			{Name: "requestBody.values", MapTo: "opt.values", MapType: "java.lang.Integer,java.lang.String"},
		},
	}, "/users/:id")

	req, err := http.NewRequest(
		http.MethodPost,
		"http://example.com/users/42",
		bytes.NewBufferString(`{"values":["7","alice"]}`),
	)
	require.NoError(t, err)

	outbound, err := handler.BuildOutbound(req, api)
	require.NoError(t, err)
	assert.Equal(t, []any{7, "alice"}, outbound.Arguments)
	assert.Equal(t, []string{"java.lang.Integer", "java.lang.String"}, outbound.ParamTypes)
}

func TestNormalizeOptTypes(t *testing.T) {
	handler := &DubboHandler{}

	t.Run("nil", func(t *testing.T) {
		types, err := handler.normalizeOptTypes(nil)
		require.NoError(t, err)
		assert.Nil(t, types)
	})

	t.Run("string slice", func(t *testing.T) {
		types, err := handler.normalizeOptTypes([]string{" java.lang.Integer ", "java.lang.String"})
		require.NoError(t, err)
		assert.Equal(t, []string{"java.lang.Integer", "java.lang.String"}, types)
	})
}

func TestNormalizeOptValues(t *testing.T) {
	handler := &DubboHandler{}

	t.Run("nil", func(t *testing.T) {
		values, err := handler.normalizeOptValues(nil)
		require.NoError(t, err)
		assert.Nil(t, values)
	})

	t.Run("string slice", func(t *testing.T) {
		values, err := handler.normalizeOptValues([]string{"1", "2"})
		require.NoError(t, err)
		assert.Equal(t, []any{"1", "2"}, values)
	})

	t.Run("default single value", func(t *testing.T) {
		values, err := handler.normalizeOptValues(7)
		require.NoError(t, err)
		assert.Equal(t, []any{7}, values)
	})
}
