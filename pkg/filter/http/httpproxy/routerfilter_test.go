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

package httpproxy

import (
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

import (
	"github.com/agiledragon/gomonkey/v2"

	"github.com/stretchr/testify/require"
)

import (
	_ "github.com/apache/dubbo-go-pixiu/pkg/cluster/loadbalancer/roundrobin"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
)

func TestForwardedRequestPreservesRawPath(t *testing.T) {
	upstream := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Upstream-Escaped-Path", r.URL.EscapedPath())
		w.Header().Set("X-Upstream-Raw-Query", r.URL.RawQuery)
		_, _ = io.WriteString(w, r.RequestURI)
	}))
	t.Cleanup(upstream.Close)
	host, portText, err := net.SplitHostPort(upstream.Listener.Addr().String())
	require.NoError(t, err)
	port, err := strconv.Atoi(portText)
	require.NoError(t, err)
	manager := server.CreateDefaultClusterManager(&model.Bootstrap{
		StaticResources: model.StaticResources{
			Clusters: []*model.ClusterConfig{{
				Name:  "raw-path-upstream",
				LbStr: model.LoadBalancerRoundRobin,
				Endpoints: []*model.Endpoint{{
					ID:      "upstream",
					Address: model.SocketAddress{Address: host, Port: port},
				}},
			}},
		},
	})
	factory, err := (&Plugin{}).CreateFilterFactory()
	require.NoError(t, err)
	factory.Config().(*Config).Scheme = "http"
	require.NoError(t, factory.Apply())
	// Change executable code only while this fixture has no active HTTP callbacks.
	// Substitute only the process-wide lookup; endpoint selection and HTTP I/O are real.
	patch := gomonkey.ApplyFunc(server.GetClusterManager, func() *server.ClusterManager { return manager })
	t.Cleanup(func() {
		factory.(*FilterFactory).client.CloseIdleConnections()
		upstream.Close()
		patch.Reset()
	})
	upstream.Start()

	cases := []struct {
		name    string
		uri     string
		path    string
		rawPath string
		want    string
	}{
		{name: "escaped slash", uri: "/files/a%2Fb?token=x%2Fy&token=z", want: "/files/a%2Fb?token=x%2Fy&token=z"},
		{name: "lowercase escaped slash", uri: "/files/a%2fb?encoded=%252F", want: "/files/a%2fb?encoded=%252F"},
		{name: "escaped dot", uri: "/files/dot%2Ename", want: "/files/dot%2Ename"},
		{name: "escaped tilde", uri: "/files/tilde%7euser?x=%7e", want: "/files/tilde%7euser?x=%7e"},
		{name: "percent", uri: "/files/100%25?empty=&q=a+b", want: "/files/100%25?empty=&q=a+b"},
		{name: "question mark", uri: "/files/question%3Fmark?x=1", want: "/files/question%3Fmark?x=1"},
		{name: "hash", uri: "/files/hash%23part?x=1", want: "/files/hash%23part?x=1"},
		{name: "double encoded slash", uri: "/files/%252F?encoded=%252F", want: "/files/%252F?encoded=%252F"},
		{name: "plain", uri: "/files/plain?x=1&x=2", want: "/files/plain?x=1&x=2"},
		{name: "double slash", uri: "/files//double?x=1", want: "/files//double?x=1"},
		{name: "unicode", uri: "/files/%E4%B8%AD%E6%96%87", want: "/files/%E4%B8%AD%E6%96%87"},
		{name: "stale raw path after rewrite", uri: "/files/a%2Fb?x=1", path: "/rewritten/path", want: "/rewritten/path?x=1"},
		{name: "invalid raw path", uri: "/files/plain?x=1", rawPath: "/files/%ZZ", want: "/files/plain?x=1"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, tc.uri, nil)
			if tc.path != "" {
				request.URL.Path = tc.path
			}
			if tc.rawPath != "" {
				request.URL.RawPath = tc.rawPath
			}
			ctx := &contexthttp.HttpContext{
				Request: request,
				Writer:  httptest.NewRecorder(),
				Route:   &model.RouteAction{Cluster: "raw-path-upstream"},
			}
			chain := filter.NewDefaultFilterChain()
			require.NoError(t, factory.PrepareFilterChain(ctx, chain))
			chain.OnDecode(ctx)
			response, ok := ctx.SourceResp.(*http.Response)
			require.True(t, ok, "expected upstream response")
			defer response.Body.Close()
			body, err := io.ReadAll(response.Body)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, response.StatusCode)
			require.Equal(t, request.URL.EscapedPath(), response.Header.Get("X-Upstream-Escaped-Path"))
			require.Equal(t, request.URL.RawQuery, response.Header.Get("X-Upstream-Raw-Query"))
			require.Equal(t, tc.want, string(body))
		})
	}
}
