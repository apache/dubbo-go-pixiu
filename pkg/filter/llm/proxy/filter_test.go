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

func TestFilterFactoryCooldownStoreIsShared(t *testing.T) {
	factory := &FilterFactory{cfg: &Config{}}

	store := factory.cooldownStore()

	assert.NotNil(t, store)
	assert.Same(t, store, factory.cooldownStore())
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
		scheme:    "http",
		cooldowns: newCooldownStore(),
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
					cooldowns:      filter.cooldowns,
				})
			}
		}()
	}

	close(start)
	wg.Wait()

	for _, endpoint := range endpoints {
		assert.Equal(t, map[string]string{"static": "value"}, endpoint.Metadata)
		lastFailure, ok := filter.cooldowns.lastFailure(clusterName, endpoint)
		assert.True(t, ok)
		assert.False(t, lastFailure.IsZero())
	}
}

func TestRequestExecutorEndpointInCooldownClearsExpiredCooldownFromProxyStore(t *testing.T) {
	clusterName := "llm-expired-cooldown"
	endpoint := testLLMEndpoint("ep-1", 18082)
	store := newCooldownStore()
	executor := &RequestExecutor{
		clusterName: clusterName,
		cooldowns:   store,
	}
	store.markFailure(clusterName, endpoint, time.Now().Add(-time.Hour))

	assert.False(t, executor.endpointInCooldown(endpoint))

	_, ok := store.lastFailure(clusterName, endpoint)
	assert.False(t, ok)
}

func TestRequestExecutorMarkEndpointCooldownStoresCooldownInProxyStore(t *testing.T) {
	clusterName := "llm-mark-cooldown"
	endpoint := testLLMEndpoint("ep-1", 18083)
	store := newCooldownStore()
	executor := &RequestExecutor{
		clusterName: clusterName,
		cooldowns:   store,
	}

	executor.markEndpointCooldown(endpoint)

	lastFailure, ok := store.lastFailure(clusterName, endpoint)
	assert.True(t, ok)
	assert.False(t, lastFailure.IsZero())
}

func TestRequestExecutorCooldownIsIsolatedByEndpointAddress(t *testing.T) {
	clusterName := "llm-address-cooldown"
	oldEndpoint := testLLMEndpoint("ep-1", 18084)
	movedEndpoint := testLLMEndpoint("ep-1", 18085)
	store := newCooldownStore()
	executor := &RequestExecutor{
		clusterName: clusterName,
		cooldowns:   store,
	}

	executor.markEndpointCooldown(oldEndpoint)

	assert.True(t, executor.endpointInCooldown(oldEndpoint))
	assert.False(t, executor.endpointInCooldown(movedEndpoint))
	_, ok := store.lastFailure(clusterName, movedEndpoint)
	assert.False(t, ok)
}

func TestRequestExecutorCooldownUsesCurrentEndpointInterval(t *testing.T) {
	clusterName := "llm-current-interval-cooldown"
	oldEndpoint := testLLMEndpoint("ep-1", 18088)
	oldEndpoint.LLMMeta.HealthCheckInterval = 10
	replacement := testLLMEndpoint("ep-1", 18088)
	replacement.LLMMeta.HealthCheckInterval = 60000
	store := newCooldownStore()
	executor := &RequestExecutor{
		clusterName: clusterName,
		cooldowns:   store,
	}
	store.markFailure(clusterName, oldEndpoint, time.Now().Add(-50*time.Millisecond))

	assert.True(t, executor.endpointInCooldown(replacement))
}

func TestCooldownStoreLazySweepRemovesExpiredChurnedEndpointEntry(t *testing.T) {
	clusterName := "llm-churn-cooldown"
	oldEndpoint := testLLMEndpoint("ep-1", 18084)
	movedEndpoint := testLLMEndpoint("ep-1", 18085)
	activeEndpoint := testLLMEndpoint("ep-2", 18086)
	store := newCooldownStore()

	store.markFailure(clusterName, oldEndpoint, time.Now().Add(-time.Hour))
	store.markFailure(clusterName, movedEndpoint, time.Now())

	store.mu.Lock()
	_, oldExistsAfterMove := store.lastFailureByEndpoint[newCooldownKey(clusterName, oldEndpoint)]
	_, movedExists := store.lastFailureByEndpoint[newCooldownKey(clusterName, movedEndpoint)]
	store.mu.Unlock()
	assert.False(t, oldExistsAfterMove)
	assert.True(t, movedExists)

	store.markFailure(clusterName, movedEndpoint, time.Now().Add(-time.Hour))
	store.lastFailure(clusterName, activeEndpoint)

	store.mu.Lock()
	_, movedExistsAfterDeletedStyleSweep := store.lastFailureByEndpoint[newCooldownKey(clusterName, movedEndpoint)]
	store.mu.Unlock()
	assert.False(t, movedExistsAfterDeletedStyleSweep)
}

func TestStrategyExecuteIgnoresUnhealthyPreferredEndpoint(t *testing.T) {
	clusterName := "llm-preferred-health"
	healthyEndpoint := testLLMEndpoint("ep-1", 18086)
	preferredEndpoint := testLLMEndpoint("ep-2", 18087)
	preferredEndpoint.UnHealthy = true
	clusterManager := server.CreateDefaultClusterManager(&model.Bootstrap{
		StaticResources: model.StaticResources{
			Clusters: []*model.ClusterConfig{{
				Name:      clusterName,
				LbStr:     model.LoadBalancerRoundRobin,
				Endpoints: []*model.Endpoint{healthyEndpoint, preferredEndpoint},
			}},
		},
	})

	var attemptedHost string
	filter := &Filter{
		client: http.Client{
			Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				attemptedHost = req.URL.Host
				return &http.Response{
					StatusCode: http.StatusOK,
					Status:     "200 OK",
					Body:       http.NoBody,
				}, nil
			}),
		},
		scheme:    "http",
		cooldowns: newCooldownStore(),
	}
	req, err := http.NewRequest(http.MethodPost, "http://example.com/v1/chat/completions", http.NoBody)
	if !assert.NoError(t, err) {
		return
	}
	hc := &contexthttp.HttpContext{
		Request: req,
		Params: map[string]any{
			llmPreferredEndpointIDKey: preferredEndpoint.ID,
		},
	}

	resp, err := (&Strategy{}).Execute(&RequestExecutor{
		hc:             hc,
		filter:         filter,
		clusterName:    clusterName,
		clusterManager: clusterManager,
		cooldowns:      filter.cooldowns,
	})

	assert.NoError(t, err)
	if assert.NotNil(t, resp) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
	assert.Equal(t, healthyEndpoint.Address.GetAddress(), attemptedHost)
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
