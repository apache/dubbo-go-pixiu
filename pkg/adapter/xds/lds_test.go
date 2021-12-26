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
		filter *pixiupb.Filter
	}
	tests := []struct {
		name  string
		args  args
		wantM map[string]interface{}
	}{
		{
			name: "yaml",
			args: args{
				filter: &pixiupb.Filter{
					Name: "",
					Config: &pixiupb.Filter_Yaml{
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
				filter: &pixiupb.Filter{
					Name:   "",
					Config: &pixiupb.Filter_Struct{Struct: httpManagerConfigStruct},
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
