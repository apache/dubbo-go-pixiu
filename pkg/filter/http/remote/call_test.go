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
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/client/dubbo"
	clienthttp "github.com/apache/dubbo-go-pixiu/pkg/client/http"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	extfilter "github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/config"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/router"
)

type recordingDubboClient struct {
	req           *dubbo.DubboOutboundRequest
	res           any
	err           error
	contextMarker any
}

type testContextKey struct{}

func (c *recordingDubboClient) Apply() error {
	return nil
}

func (c *recordingDubboClient) Close() error {
	return nil
}

func (c *recordingDubboClient) Call(ctx context.Context, req *dubbo.DubboOutboundRequest) (any, error) {
	c.req = req
	c.contextMarker = ctx.Value(testContextKey{})
	return c.res, c.err
}

func TestMatchClientRoutesHTTPToHTTPClient(t *testing.T) {
	filter := &Filter{conf: filterConfig{DubboProxyConfig: &dubbo.DubboProxyConfig{}}}

	cli, err := filter.matchHTTPClient(constant.HTTPRequest)
	require.NoError(t, err)
	assert.Same(t, clienthttp.SingletonHTTPClient(), cli)
}

func TestDecodeRoutesDubboAndTripleThroughOutboundClient(t *testing.T) {
	for _, requestType := range []string{constant.DubboRequest, constant.TripleRequest} {
		t.Run(requestType, func(t *testing.T) {
			resp := map[string]string{"ok": requestType}

			marker := struct{}{}
			req := httptest.NewRequest(http.MethodPost, "http://example.com/users/42?name=alice", nil)
			req = req.WithContext(context.WithValue(req.Context(), testContextKey{}, marker))
			recorder := &recordingDubboClient{res: resp}
			ctx := &contexthttp.HttpContext{
				Timeout: 150 * time.Millisecond,
				Request: req,
				Writer:  httptest.NewRecorder(),
			}
			ctx.API(router.API{
				Method: config.Method{
					Timeout:  time.Second,
					HTTPVerb: http.MethodPost,
					IntegrationRequest: config.IntegrationRequest{
						RequestType: requestType,
						DubboBackendConfig: config.DubboBackendConfig{
							Interface: "com.demo.UserService",
							Method:    "SayHello",
						},
						MappingParams: []config.MappingParam{
							{Name: "queryStrings.name", MapTo: "0"},
						},
					},
				},
			})

			f := &Filter{
				conf:        filterConfig{DubboProxyConfig: &dubbo.DubboProxyConfig{}},
				dubboClient: recorder,
			}
			status := f.Decode(ctx)

			require.Equal(t, extfilter.Continue, status)
			assert.Equal(t, resp, ctx.SourceResp)
			require.NotNil(t, recorder.req)
			assert.Equal(t, marker, recorder.contextMarker)
			assert.Equal(t, "com.demo.UserService", recorder.req.Service)
			assert.Equal(t, "SayHello", recorder.req.Method)
			assert.Equal(t, []any{"alice"}, recorder.req.Arguments)
			assert.Equal(t, 150*time.Millisecond, recorder.req.Timeout)
		})
	}
}

func TestMatchClientRejectsUnknownRequestType(t *testing.T) {
	filter := &Filter{conf: filterConfig{DubboProxyConfig: &dubbo.DubboProxyConfig{}}}

	cli, err := filter.matchHTTPClient("grpc")
	assert.Nil(t, cli)
	assert.EqualError(t, err, "not support")
}

func TestDecodeStopsWithLocalReplyWhenBuildOutboundFails(t *testing.T) {
	recorder := &recordingDubboClient{res: "should not be called"}

	writer := httptest.NewRecorder()
	ctx := &contexthttp.HttpContext{
		Timeout: time.Second,
		Request: httptest.NewRequest(http.MethodPost, "http://example.com/users/42?app=demo", nil),
		Writer:  writer,
	}
	ctx.API(router.API{
		Method: config.Method{
			Timeout:  time.Second,
			HTTPVerb: http.MethodPost,
			IntegrationRequest: config.IntegrationRequest{
				RequestType: constant.DubboRequest,
				DubboBackendConfig: config.DubboBackendConfig{
					Interface: "com.demo.UserService",
					Method:    "SayHello",
				},
				MappingParams: []config.MappingParam{
					{Name: "queryStrings.app", MapTo: "opt.application"},
				},
			},
		},
	})

	f := &Filter{
		conf:        filterConfig{DubboProxyConfig: &dubbo.DubboProxyConfig{}},
		dubboClient: recorder,
	}
	status := f.Decode(ctx)

	assert.Equal(t, extfilter.Stop, status)
	assert.Nil(t, recorder.req)
	assert.True(t, ctx.LocalReply())
	assert.Equal(t, http.StatusInternalServerError, writer.Code)
	assert.Contains(t, writer.Body.String(), "client call error")
}

func TestFilterFactoryApplyOnlyInitializesDubboClient(t *testing.T) {
	originalDubboInit := initDubboClient
	t.Cleanup(func() {
		initDubboClient = originalDubboInit
	})

	dubboInitCalls := 0
	initDubboClient = func(conf *dubbo.DubboProxyConfig) {
		dubboInitCalls++
	}

	factory := &FilterFactory{
		conf: &filterConfig{
			DubboProxyConfig: &dubbo.DubboProxyConfig{},
		},
	}

	require.NoError(t, factory.Apply())
	assert.Equal(t, 1, dubboInitCalls)
}

func TestFilterFactoryApplyRejectsDeprecatedAutoResolve(t *testing.T) {
	trueVal := true
	factory := &FilterFactory{
		conf: &filterConfig{
			DubboProxyConfig: &dubbo.DubboProxyConfig{
				AutoResolve: &trueVal,
			},
		},
	}
	err := factory.Apply()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "auto_resolve is no longer supported")
}
