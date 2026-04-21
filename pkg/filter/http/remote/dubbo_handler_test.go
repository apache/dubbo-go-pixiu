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
	"github.com/apache/dubbo-go-pixiu/pkg/common/mock"
	"github.com/apache/dubbo-go-pixiu/pkg/config"
	"github.com/apache/dubbo-go-pixiu/pkg/router"
)

func newTestAPI() router.API {
	api := mock.GetMockAPI(http.MethodPost, "/users/:id")
	api.IntegrationRequest = config.IntegrationRequest{
		RequestType: constant.DubboRequest,
		DubboBackendConfig: config.DubboBackendConfig{
			Interface: "com.demo.UserService",
			Method:    "SayHello",
			Group:     "demo-group",
			Version:   "1.0.0",
		},
	}
	return api
}

func TestBuildOutboundMapsQueryHeaderBodyAndURI(t *testing.T) {
	handler := &DubboHandler{}
	api := newTestAPI()
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
		api := newTestAPI()
		api.IntegrationRequest.ParameterTypes = []string{"java.lang.Integer"}
		api.MappingParams = []config.MappingParam{
			{Name: "requestBody.values", MapTo: "opt.values"},
			{Name: "requestBody.types", MapTo: "opt.types"},
		}

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
		api := newTestAPI()
		api.MappingParams = []config.MappingParam{
			{Name: "requestBody.values", MapTo: "opt.values"},
			{Name: "requestBody.types", MapTo: "opt.types"},
		}

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
		api := newTestAPI()
		api.MappingParams = []config.MappingParam{
			{Name: "queryStrings.name", MapTo: "0"},
		}

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
	api := newTestAPI()
	api.MappingParams = []config.MappingParam{
		{Name: "queryStrings.name", MapTo: "0"},
		{Name: "requestBody.values", MapTo: "opt.values"},
	}

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
	api := newTestAPI()
	api.MappingParams = []config.MappingParam{
		{Name: "queryStrings.name", MapTo: "opt.unknown"},
	}

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
	api := newTestAPI()
	api.MappingParams = []config.MappingParam{
		{Name: "queryStrings.app", MapTo: "opt.application"},
	}

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

func TestBuildOutboundRejectsProtocolSchemeMismatch(t *testing.T) {
	handler := &DubboHandler{}
	api := newTestAPI()
	api.IntegrationRequest.Protocol = "tri"
	api.IntegrationRequest.URL = "dubbo://127.0.0.1:20880"
	api.IntegrationRequest.Serialization = "hessian2"

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
	api := newTestAPI()
	api.IntegrationRequest.URL = "dubbo://127.0.0.1:20880"

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

func TestBuildOutboundResolvesProtocolPriority(t *testing.T) {
	handler := &DubboHandler{}

	t.Run("integration request protocol wins", func(t *testing.T) {
		api := newTestAPI()
		api.IntegrationRequest.Protocol = "triple"
		api.IntegrationRequest.RequestType = constant.DubboRequest

		req, err := http.NewRequest(http.MethodPost, "http://example.com/users/42", bytes.NewBufferString(`{}`))
		require.NoError(t, err)

		outbound, err := handler.BuildOutbound(req, api)
		require.NoError(t, err)
		assert.Equal(t, "tri", outbound.Protocol)
	})

	t.Run("request type wins when protocol empty", func(t *testing.T) {
		api := newTestAPI()
		api.IntegrationRequest.Protocol = ""
		api.IntegrationRequest.RequestType = "triple"

		req, err := http.NewRequest(http.MethodPost, "http://example.com/users/42", bytes.NewBufferString(`{}`))
		require.NoError(t, err)

		outbound, err := handler.BuildOutbound(req, api)
		require.NoError(t, err)
		assert.Equal(t, "tri", outbound.Protocol)
	})

	t.Run("default to dubbo", func(t *testing.T) {
		api := newTestAPI()
		api.IntegrationRequest.Protocol = ""
		api.IntegrationRequest.RequestType = ""

		req, err := http.NewRequest(http.MethodPost, "http://example.com/users/42", bytes.NewBufferString(`{}`))
		require.NoError(t, err)

		outbound, err := handler.BuildOutbound(req, api)
		require.NoError(t, err)
		assert.Equal(t, "dubbo", outbound.Protocol)
	})
}

func TestBuildOutboundValidatesDirectInvokeArity(t *testing.T) {
	handler := &DubboHandler{}
	api := newTestAPI()
	api.IntegrationRequest.URL = "dubbo://127.0.0.1:20880"
	api.IntegrationRequest.Serialization = "hessian2"
	api.IntegrationRequest.ParameterTypes = []string{"java.lang.Integer", clientdubbo.JavaStringClassName}
	api.MappingParams = []config.MappingParam{
		{Name: "requestBody.values", MapTo: "opt.values"},
	}

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
