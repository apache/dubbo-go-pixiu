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
	pixiupb "github.com/dubbo-go-pixiu/pixiu-api/pkg/xds/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	structpb2 "google.golang.org/protobuf/types/known/structpb"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
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
	configMap := map[string]interface{}{
		"route_config": map[string]interface{}{
			"routes": []interface{}{
				map[string]interface{}{
					"match": map[string]interface{}{
						"prefix": "/",
					},
					"route": map[string]interface{}{
						"cluster":                         "http_bin",
						"cluster_not_found_response_code": "505",
					},
				},
			},
		},
		"http_filters": []interface{}{
			map[string]interface{}{
				"name":   "dgp.filter.http.httpproxy",
				"config": nil,
			},
			map[string]interface{}{
				"name":   "dgp.filter.http.response",
				"config": nil,
			},
		},
	}
	httpManagerConfigStruct, _ := structpb2.NewStruct(configMap)

	oneConfigMap := map[string]interface{}{
		"route_config": map[interface{}]interface{}{
			"routes": []interface{}{
				map[interface{}]interface{}{
					"match": map[interface{}]interface{}{
						"prefix": "/",
					},
					"route": map[interface{}]interface{}{
						"cluster":                         "http_bin",
						"cluster_not_found_response_code": "505",
					},
				},
			},
		},
		"http_filters": []interface{}{
			map[interface{}]interface{}{
				"name":   "dgp.filter.http.httpproxy",
				"config": nil,
			},
			map[interface{}]interface{}{
				"name":   "dgp.filter.http.response",
				"config": nil,
			},
		},
	}

	type args struct {
		filter *pixiupb.NetworkFilter
	}
	tests := []struct {
		name  string
		args  args
		wantM map[string]interface{}
	}{
		{
			name: "yaml",
			args: args{
				filter: &pixiupb.NetworkFilter{
					Name: "",
					Config: &pixiupb.NetworkFilter_Yaml{
						Yaml: &pixiupb.Config{
							Content: httpManagerConfigYaml,
						}},
				},
			},
			wantM: oneConfigMap,
		},
		{
			name: "struct",
			args: args{
				filter: &pixiupb.NetworkFilter{
					Name:   "",
					Config: &pixiupb.NetworkFilter_Struct{Struct: httpManagerConfigStruct},
				},
			},
			wantM: configMap,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LdsManager{}
			r := l.makeConfig(tt.args.filter)
			assert := require.New(t)
			assert.Equal(tt.wantM, r)
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
	l := &pixiupb.Listener{}
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
	m map[string]*model.Listener
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

func TestSetupListeners(t *testing.T) {
	mock := &mockListenerManager{m: map[string]*model.Listener{}}
	lm := &LdsManager{listenerMg: mock}

	listeners := []*pixiupb.Listener{
		{
			Protocol: pixiupb.Listener_HTTP,
			Address: &pixiupb.Address{
				SocketAddress: &pixiupb.SocketAddress{
					Address: "0.0.0.0",
					Port:    8080,
				},
			},
			FilterChain: &pixiupb.FilterChain{},
		},
		{
			Protocol: pixiupb.Listener_TRIPLE,
			Address: &pixiupb.Address{
				SocketAddress: &pixiupb.SocketAddress{
					Address: "0.0.0.0",
					Port:    8081,
				},
			},
			FilterChain: &pixiupb.FilterChain{},
		},
	}
	lm.setupListeners(listeners)
	for _, v := range listeners {
		assert.Equal(t, v.Protocol.String(), model.ProtocolTypeName[int32(mock.m[v.Name].Protocol)])
		assert.Equal(t, v.Address.SocketAddress.Address, mock.m[v.Name].Address.SocketAddress.Address)
		assert.Equal(t, int(v.Address.SocketAddress.Port), mock.m[v.Name].Address.SocketAddress.Port)
	}

	newListeners := []*pixiupb.Listener{
		{
			Protocol: pixiupb.Listener_HTTP,
			Address: &pixiupb.Address{
				SocketAddress: &pixiupb.SocketAddress{
					Address: "0.0.0.0",
					Port:    8080,
				},
			},
			FilterChain: &pixiupb.FilterChain{},
		},
	}
	lm.setupListeners(newListeners)
	assert.Equal(t, 1, len(mock.m))
}
