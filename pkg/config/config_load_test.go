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

package config

import (
	"encoding/json"
	"log"
	"os"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/configcenter"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

var b model.Bootstrap

func TestMain(m *testing.M) {
	log.Println("Prepare Bootstrap")

	hcmc := model.HttpConnectionManagerConfig{
		RouteConfig: model.RouteConfiguration{
			Routes: []*model.Router{
				{
					ID: "1",
					Match: model.RouterMatch{
						Prefix: "/api/v1",
						Methods: []string{
							"POST",
						},
						Path: "",
					},
					Route: model.RouteAction{
						Cluster:                     "test_dubbo",
						ClusterNotFoundResponseCode: 505,
					},
				},
			},
			Dynamic: false,
		},
		HTTPFilters: []*model.HTTPFilter{
			{
				Name:   "dgp.filters.http.api",
				Config: nil,
			},
			{
				Name:   "dgp.filters.http.router",
				Config: nil,
			},
			{
				Name:   "dgp.filters.http.dubboproxy",
				Config: nil,
			},
		},
		ServerName:        "test_http_dubbo",
		GenerateRequestID: false,
		IdleTimeoutStr:    "100",
		TimeoutStr:        "10s",
	}

	var inInterface map[string]any
	inrec, _ := json.Marshal(hcmc)
	_ = json.Unmarshal(inrec, &inInterface)
	b = model.Bootstrap{
		StaticResources: model.StaticResources{
			Listeners: []*model.Listener{
				{
					Name:        "net/http",
					ProtocolStr: "HTTPS",
					Address: model.Address{
						SocketAddress: model.SocketAddress{
							Address: "0.0.0.0",
							Port:    443,
						},
					},
					Config: model.HttpConfig{
						IdleTimeoutStr:  "5s",
						WriteTimeoutStr: "5s",
						ReadTimeoutStr:  "5s",
					},
					FilterChain: model.FilterChain{
						Filters: []model.NetworkFilter{
							{
								Name:   "dgp.filter.httpconnectionmanager",
								Config: inInterface,
							},
						},
					},
				},
			},
			Clusters: []*model.ClusterConfig{
				{
					Name:    "test_dubbo",
					TypeStr: "EDS",
					Type:    model.EDS,
					LbStr:   "RoundRobin",
				},
			},
			ShutdownConfig: &model.ShutdownConfig{
				Timeout:      "60s",
				StepTimeout:  "10s",
				RejectPolicy: "immediacy",
			},
			PprofConf: model.PprofConf{
				Enable: true,
				Address: model.Address{
					SocketAddress: model.SocketAddress{
						Address: "0.0.0.0",
						Port:    6060,
					},
				},
			},
		},
	}
	retCode := m.Run()
	os.Exit(retCode)
}

func TestLoad(t *testing.T) {
	conf := Load("conf_test.yaml")
	assert.Equal(t, 1, len(conf.StaticResources.Listeners))
	assert.Equal(t, 1, len(conf.StaticResources.Clusters))
	assert.Nil(t, Adapter(&b))
	str1, _ := json.Marshal(conf.StaticResources)
	str2, _ := json.Marshal(b.StaticResources)
	assert.Equal(t, string(str1), string(str2))
}

func TestStruct2JSON(t *testing.T) {
	if bytes, err := json.Marshal(b); err != nil {
		t.Fatal(err)
	} else {
		t.Log(string(bytes))
	}
}

// closeRecorderLoad is a configcenter.Load stand-in that records Close calls
// so ConfigManager.Close can be asserted without a live remote config center.
type closeRecorderLoad struct {
	closeCount int
}

func (c *closeRecorderLoad) LoadConfigs(boot *model.Bootstrap, opts ...configcenter.Option) (*model.Bootstrap, error) {
	return nil, nil
}

func (c *closeRecorderLoad) ViewRemoteConfig() *model.Bootstrap { return nil }

func (c *closeRecorderLoad) Close() { c.closeCount++ }

// TestConfigManager_Close verifies the v2 close path: Close must forward to
// the underlying remote config loader so the Nacos v2 gRPC connection is
// released on shutdown. See AlexStocks' [P1] review on PR #982.
func TestConfigManager_Close(t *testing.T) {
	t.Run("forwards to the loader", func(t *testing.T) {
		recorder := &closeRecorderLoad{}
		m := &ConfigManager{load: recorder}

		assert.Equal(t, 0, recorder.closeCount, "loader must not be closed before Close")

		m.Close()

		assert.Equal(t, 1, recorder.closeCount, "Close must forward to the underlying loader")
	})

	t.Run("no-op when no remote config center", func(t *testing.T) {
		m := &ConfigManager{load: nil}
		assert.NotPanics(t, func() { m.Close() })
	})
}
