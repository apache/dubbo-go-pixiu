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

package proxy

import (
	"errors"
	"net/http"
	"sync"
	"testing"
	"time"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	_ "github.com/apache/dubbo-go-pixiu/pkg/cluster/loadbalancer/roundrobin" // Register RoundRobin for LLM proxy cluster tests.
	_ "github.com/apache/dubbo-go-pixiu/pkg/cluster/retry/noretry"           // Register NoRetry for LLM proxy retry tests.
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func TestStrategyExecuteUsesRuntimeCooldownStateWithoutMutatingEndpointMetadata(t *testing.T) {
	clusterName := "llm-runtime-cooldown"
	endpoints := []*model.Endpoint{
		testLLMEndpoint("ep-1", 18080),
		testLLMEndpoint("ep-2", 18081),
	}
	clusterConfig := &model.ClusterConfig{
		Name:      clusterName,
		LbStr:     model.LoadBalancerRoundRobin,
		Endpoints: endpoints,
	}
	clusterManager := server.CreateDefaultClusterManager(&model.Bootstrap{
		StaticResources: model.StaticResources{
			Clusters: []*model.ClusterConfig{clusterConfig},
		},
	})
	filter := &Filter{
		client: http.Client{
			Transport: roundTripFunc(func(*http.Request) (*http.Response, error) {
				return nil, errors.New("upstream unavailable")
			}),
		},
		scheme: "http",
	}
	strategy := &Strategy{}

	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 16; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < 50; j++ {
				req, err := http.NewRequest(http.MethodPost, "http://example.com/v1/chat/completions", http.NoBody)
				if err != nil {
					t.Errorf("new request: %v", err)
					return
				}
				hc := &contexthttp.HttpContext{
					Request: req,
					Params:  map[string]any{},
				}
				_, _ = strategy.Execute(&RequestExecutor{
					hc:             hc,
					filter:         filter,
					clusterName:    clusterName,
					clusterManager: clusterManager,
				})
			}
		}()
	}

	close(start)
	wg.Wait()

	for _, endpoint := range endpoints {
		assert.Equal(t, map[string]string{"static": "value"}, endpoint.Metadata)
		state := clusterManager.GetEndpointRuntimeState(clusterName, endpoint.ID, endpoint.Address.GetAddress())
		if assert.NotNil(t, state) {
			unhealthy, ok := state.Load(LLMUnhealthyKey)
			assert.True(t, ok)
			assert.Equal(t, "true", unhealthy)
			_, ok = state.Load(HealthyCheckTimeKey)
			assert.True(t, ok)
		}
	}
}

func TestRequestExecutorEndpointInCooldownClearsExpiredCooldownAtomically(t *testing.T) {
	clusterName := "llm-expired-cooldown"
	endpoint := testLLMEndpoint("ep-1", 18082)
	clusterManager := server.CreateDefaultClusterManager(&model.Bootstrap{
		StaticResources: model.StaticResources{
			Clusters: []*model.ClusterConfig{{
				Name:      clusterName,
				LbStr:     model.LoadBalancerRoundRobin,
				Endpoints: []*model.Endpoint{endpoint},
			}},
		},
	})
	executor := &RequestExecutor{
		clusterName:    clusterName,
		clusterManager: clusterManager,
	}
	state := clusterManager.GetEndpointRuntimeState(clusterName, endpoint.ID, endpoint.Address.GetAddress())
	if !assert.NotNil(t, state) {
		return
	}
	state.StoreMany(map[string]string{
		LLMUnhealthyKey:     "true",
		HealthyCheckTimeKey: time.Now().Add(-time.Hour).Format(time.RFC3339),
	})

	assert.False(t, executor.endpointInCooldown(endpoint))

	values := state.LoadMany(LLMUnhealthyKey, HealthyCheckTimeKey)
	_, unhealthyExists := values[LLMUnhealthyKey]
	_, timeExists := values[HealthyCheckTimeKey]
	assert.False(t, unhealthyExists)
	assert.False(t, timeExists)
}

func TestRequestExecutorMarkEndpointCooldownStoresCooldownPairAtomically(t *testing.T) {
	clusterName := "llm-mark-cooldown"
	endpoint := testLLMEndpoint("ep-1", 18083)
	clusterManager := server.CreateDefaultClusterManager(&model.Bootstrap{
		StaticResources: model.StaticResources{
			Clusters: []*model.ClusterConfig{{
				Name:      clusterName,
				LbStr:     model.LoadBalancerRoundRobin,
				Endpoints: []*model.Endpoint{endpoint},
			}},
		},
	})
	executor := &RequestExecutor{
		clusterName:    clusterName,
		clusterManager: clusterManager,
	}

	executor.markEndpointCooldown(endpoint)

	state := clusterManager.GetEndpointRuntimeState(clusterName, endpoint.ID, endpoint.Address.GetAddress())
	if !assert.NotNil(t, state) {
		return
	}
	values := state.LoadMany(LLMUnhealthyKey, HealthyCheckTimeKey)
	assert.Equal(t, "true", values[LLMUnhealthyKey])
	assert.NotEmpty(t, values[HealthyCheckTimeKey])
}

func testLLMEndpoint(id string, port int) *model.Endpoint {
	return &model.Endpoint{
		ID:       id,
		Name:     "endpoint-" + id,
		Metadata: map[string]string{"static": "value"},
		Address: model.SocketAddress{
			Address: "127.0.0.1",
			Port:    port,
		},
		LLMMeta: &model.LLMMeta{
			RetryPolicy: model.RetryPolicy{
				Name: model.RetryerNoRetry,
			},
			Fallback:            true,
			HealthCheckInterval: 60000,
		},
	}
}
