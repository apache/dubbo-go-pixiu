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
	pixiupb "github.com/dubbogo/dubbo-go-pixiu-filter/pkg/xds/model"
	"github.com/stretchr/testify/require"
	structpb2 "google.golang.org/protobuf/types/known/structpb"
	"testing"
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

	_configMap := map[string]interface{}{
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
			wantM: _configMap,
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
