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

package xds

import (
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"google.golang.org/protobuf/encoding/protojson"

	structpb2 "google.golang.org/protobuf/types/known/structpb"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/config/xds/apiclient"
	xdsmodel "github.com/apache/dubbo-go-pixiu/pkg/config/xds/model"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func TestLdsManager_makeConfig(t *testing.T) {
	var httpManagerConfigYaml = `
route_config:
  routes:
    - match:
        prefix: "/"
      route:
        cluster: "http_bin"
        cluster_not_found_response_code: "505"
http_filters:
  - name: dgp.filter.http.httpproxy
    config:
  - name: dgp.filter.http.response
    config:
`
	configMap := map[string]any{
		"route_config": map[string]any{
			"routes": []any{
				map[string]any{
					"match": map[string]any{
						"prefix": "/",
					},
					"route": map[string]any{
						"cluster":                         "http_bin",
						"cluster_not_found_response_code": "505",
					},
				},
			},
		},
		"http_filters": []any{
			map[string]any{
				"name":   "dgp.filter.http.httpproxy",
				"config": nil,
			},
			map[string]any{
				"name":   "dgp.filter.http.response",
				"config": nil,
			},
		},
	}
	httpManagerConfigStruct, _ := structpb2.NewStruct(configMap)

	type args struct {
		filter *xdsmodel.NetworkFilter
	}
	tests := []struct {
		name  string
		args  args
		wantM map[string]any
	}{
		{
			name: "yaml",
			args: args{
				filter: &xdsmodel.NetworkFilter{
					Name: "yaml_filter",
					Config: &xdsmodel.NetworkFilter_Yaml{
						Yaml: &xdsmodel.Config{
							Content: httpManagerConfigYaml,
						}},
				},
			},
			wantM: configMap,
		},
		{
			name: "struct",
			args: args{
				filter: &xdsmodel.NetworkFilter{
					Name:   "struct_filter",
					Config: &xdsmodel.NetworkFilter_Struct{Struct: httpManagerConfigStruct},
				},
			},
			wantM: configMap,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LdsManager{}
			gotM := l.makeConfig(tt.args.filter)
			assertions := require.New(t)

			assertions.Equal(tt.wantM, gotM)
		})
	}
}

func TestMakeListener(t *testing.T) {
	lm := &LdsManager{}
	json := `
{
      "name": "net/http",
      "address": {
        "socketAddress": {
          "address": "0.0.0.0",
          "port": "8080"
        }
      },
      "filterChain": {
        "filters": [
          {
            "name": "dgp.filter.httpconnectionmanager",
            "struct": {
              "http_filters": [
                {
                  "config": null,
                  "name": "dgp.filter.http.httpproxy"
                }
              ],
              "route_config": {
                "routes": [
                  {
                    "match": {
                      "prefix": "/"
                    },
                    "route": {
                      "cluster": "http_bin",
                      "cluster_not_found_response_code": 503
                    }
                  }
                ]
              }
            }
          }
        ]
      }
    }
`
	l := &xdsmodel.Listener{}
	if err := protojson.Unmarshal([]byte(json), l); err != nil {
		t.Fatal(err)
	}
	listener := lm.makeListener(l)
	assert.NotNil(t, listener)
	assert.Equal(t, "net/http", listener.Name)
	assert.Equal(t, "0.0.0.0", listener.Address.SocketAddress.Address)
	assert.Equal(t, 8080, listener.Address.SocketAddress.Port)
	assert.Equal(t, 1, len(listener.FilterChain.Filters))
}

type mockListenerManager struct {
	m          map[string]*model.Listener
	xdsManaged map[string]struct{}
}

func (m *mockListenerManager) AddListener(l *model.Listener) error {
	m.m[l.Name] = l
	return nil
}

func (m *mockListenerManager) UpdateListener(l *model.Listener) error {
	m.m[l.Name] = l
	return nil
}

func (m *mockListenerManager) RemoveListener(names []string) {
	for _, name := range names {
		delete(m.m, name)
	}
}

func (m *mockListenerManager) HasListener(name string) bool {
	_, ok := m.m[name]
	return ok
}

func (m *mockListenerManager) CloneXdsControlListener() ([]*model.Listener, error) {
	var res []*model.Listener
	for _, v := range m.m {
		res = append(res, v)
	}
	return res, nil
}

func (m *mockListenerManager) UpsertXDSListener(listener *model.Listener) error {
	if m.xdsManaged == nil {
		m.xdsManaged = make(map[string]struct{})
	}
	m.m[listener.Name] = listener
	m.xdsManaged[listener.Name] = struct{}{}
	return nil
}

func (m *mockListenerManager) RemoveXDSListeners(names []string) {
	for _, name := range names {
		if _, owned := m.xdsManaged[name]; !owned {
			continue
		}
		delete(m.m, name)
		delete(m.xdsManaged, name)
	}
}

func (m *mockListenerManager) XDSListenerNames() []string {
	res := make([]string, 0, len(m.xdsManaged))
	for name := range m.xdsManaged {
		res = append(res, name)
	}
	return res
}

func TestSetupListeners(t *testing.T) {
	staticListener := &model.Listener{Name: "static-listener"}
	mock := &mockListenerManager{m: map[string]*model.Listener{"static-listener": staticListener}}
	lm := &LdsManager{listenerMg: mock}

	listeners := []*xdsmodel.Listener{
		{
			Protocol: xdsmodel.Listener_HTTP,
			Address: &xdsmodel.Address{
				SocketAddress: &xdsmodel.SocketAddress{
					Address: "0.0.0.0",
					Port:    8080,
				},
			},
			FilterChain: &xdsmodel.FilterChain{},
		},
		{
			Protocol: xdsmodel.Listener_TRIPLE,
			Address: &xdsmodel.Address{
				SocketAddress: &xdsmodel.SocketAddress{
					Address: "0.0.0.0",
					Port:    8081,
				},
			},
			FilterChain: &xdsmodel.FilterChain{},
		},
	}
	lm.setupListeners(listeners)
	for _, v := range listeners {
		assert.Equal(t, v.Protocol.String(), model.ProtocolTypeName[int32(mock.m[v.Name].Protocol)])
		assert.Equal(t, v.Address.SocketAddress.Address, mock.m[v.Name].Address.SocketAddress.Address)
		assert.Equal(t, int(v.Address.SocketAddress.Port), mock.m[v.Name].Address.SocketAddress.Port)
	}

	newListeners := []*xdsmodel.Listener{
		{
			Protocol: xdsmodel.Listener_HTTP,
			Address: &xdsmodel.Address{
				SocketAddress: &xdsmodel.SocketAddress{
					Address: "0.0.0.0",
					Port:    8080,
				},
			},
			FilterChain: &xdsmodel.FilterChain{},
		},
	}
	lm.setupListeners(newListeners)
	assert.Equal(t, 2, len(mock.m))
	assert.Same(t, staticListener, mock.m["static-listener"])
}

func TestLdsManager_ApplyDelta(t *testing.T) {
	staticListener := &model.Listener{Name: "static-listener"}
	dynamicListener := &model.Listener{Name: "dynamic-listener"}
	mock := &mockListenerManager{
		m: map[string]*model.Listener{
			"static-listener":  staticListener,
			"dynamic-listener": dynamicListener,
		},
		xdsManaged: map[string]struct{}{"dynamic-listener": {}},
	}
	manager := &LdsManager{listenerMg: mock}

	manager.applyDelta(&apiclient.DeltaResources{})
	assert.Same(t, dynamicListener, mock.m["dynamic-listener"])
	assert.Same(t, staticListener, mock.m["static-listener"])

	manager.applyDelta(&apiclient.DeltaResources{
		RemovedResources: []string{constant.ListenerType},
	})
	assert.NotContains(t, mock.m, "dynamic-listener")
	assert.Same(t, staticListener, mock.m["static-listener"])
}
