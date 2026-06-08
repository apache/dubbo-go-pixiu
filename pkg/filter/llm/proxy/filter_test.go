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
	"fmt"
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

func TestFilterFactoriesShareRuntimeCooldownStore(t *testing.T) {
	plugin := &Plugin{}
	firstFactory, err := plugin.CreateFilterFactory()
	if !assert.NoError(t, err) {
		return
	}
	secondFactory, err := plugin.CreateFilterFactory()
	if !assert.NoError(t, err) {
		return
	}

	firstStore := firstFactory.(*FilterFactory).cooldownStore()
	secondStore := secondFactory.(*FilterFactory).cooldownStore()

	assert.NotNil(t, firstStore)
	assert.Same(t, firstStore, secondStore)

	clusterName := "llm-shared-runtime-cooldown"
	endpoint := testLLMEndpoint("ep-1", 18089)
	firstExecutor := &RequestExecutor{
		clusterName: clusterName,
		cooldowns:   firstStore,
	}
	secondExecutor := &RequestExecutor{
		clusterName: clusterName,
		cooldowns:   secondStore,
	}

	firstExecutor.markEndpointCooldown(endpoint)

	assert.True(t, secondExecutor.endpointInCooldown(endpoint))
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
		lastFailure, _, ok := filter.cooldowns.lastFailureWithCurrentTTL(clusterName, endpoint)
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

	_, _, ok := store.lastFailureWithCurrentTTL(clusterName, endpoint)
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

	lastFailure, _, ok := store.lastFailureWithCurrentTTL(clusterName, endpoint)
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
	_, _, ok := store.lastFailureWithCurrentTTL(clusterName, movedEndpoint)
	assert.False(t, ok)
}

func TestRequestExecutorCooldownIsIsolatedByEndpointCredential(t *testing.T) {
	clusterName := "llm-credential-cooldown"
	oldEndpoint := testLLMEndpoint("ep-1", 18088)
	replacement := testLLMEndpoint("ep-1", 18088)
	replacement.LLMMeta.APIKey = "fixed-key"
	store := newCooldownStore()
	executor := &RequestExecutor{
		clusterName: clusterName,
		cooldowns:   store,
	}
	executor.markEndpointCooldown(oldEndpoint)

	assert.True(t, executor.endpointInCooldown(oldEndpoint))
	assert.False(t, executor.endpointInCooldown(replacement))
	_, _, ok := store.lastFailureWithCurrentTTL(clusterName, replacement)
	assert.False(t, ok)
}

func TestRequestExecutorCooldownSurvivesNonIdentityLLMConfigChanges(t *testing.T) {
	clusterName := "llm-policy-cooldown"
	oldEndpoint := testLLMEndpoint("ep-1", 18089)
	oldEndpoint.LLMMeta.APIKey = "same-key"
	oldEndpoint.LLMMeta.HealthCheckInterval = 10
	replacement := testLLMEndpoint("ep-1", 18089)
	replacement.Name = "renamed-endpoint"
	replacement.Metadata["dynamic"] = "value"
	replacement.LLMMeta.APIKey = "same-key"
	replacement.LLMMeta.HealthCheckInterval = 60000
	replacement.LLMMeta.RetryPolicy = model.RetryPolicy{
		Name: model.RetryerCountBased,
		Config: map[string]any{
			"attempts": 2,
		},
	}
	store := newCooldownStore()
	executor := &RequestExecutor{
		clusterName: clusterName,
		cooldowns:   store,
	}
	store.markFailure(clusterName, oldEndpoint, time.Now().Add(-50*time.Millisecond))

	assert.True(t, executor.endpointInCooldown(replacement))
}

func TestCooldownStoreLazySweepKeepsEntryAfterEndpointIntervalExtends(t *testing.T) {
	clusterName := "llm-extended-cooldown"
	oldEndpoint := testLLMEndpoint("ep-1", 18090)
	oldEndpoint.LLMMeta.APIKey = "same-key"
	oldEndpoint.LLMMeta.HealthCheckInterval = 10
	replacement := testLLMEndpoint("ep-1", 18090)
	replacement.LLMMeta.APIKey = "same-key"
	replacement.LLMMeta.HealthCheckInterval = 60000
	activeEndpoint := testLLMEndpoint("ep-2", 18091)
	store := newCooldownStore()
	executor := &RequestExecutor{
		clusterName: clusterName,
		cooldowns:   store,
	}
	store.markFailure(clusterName, oldEndpoint, time.Now().Add(-50*time.Millisecond))

	assert.True(t, executor.endpointInCooldown(replacement))
	_, _, _ = store.lastFailureWithCurrentTTL(clusterName, activeEndpoint)
	store.mu.Lock()
	_, ok := store.lastFailureByEndpoint[newCooldownKey(clusterName, replacement)]
	store.mu.Unlock()
	assert.True(t, ok)
}

func TestLegacyEndpointHealthMetadataKeysRemainExported(t *testing.T) {
	assert.Equal(t, "LLMUnhealthy", LLMUnhealthyKey)
	assert.Equal(t, "HealthyCheckTime", HealthyCheckTimeKey)
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
	assert.True(t, oldExistsAfterMove)
	assert.True(t, movedExists)

	store.mu.Lock()
	store.lastSweep = time.Now().Add(-cooldownStoreSweepAfter - time.Millisecond)
	store.mu.Unlock()
	store.lastFailureWithCurrentTTL(clusterName, activeEndpoint)

	store.mu.Lock()
	_, oldExistsAfterUnrelatedSweep := store.lastFailureByEndpoint[newCooldownKey(clusterName, oldEndpoint)]
	store.mu.Unlock()
	assert.False(t, oldExistsAfterUnrelatedSweep)
}

func TestCooldownStoreLazySweepRemovesExpiredEndpointFromDifferentCluster(t *testing.T) {
	expiredEndpoint := testLLMEndpoint("ep-1", 18092)
	activeEndpoint := testLLMEndpoint("ep-2", 18093)
	store := newCooldownStore()

	store.markFailure("old-cluster", expiredEndpoint, time.Now().Add(-time.Hour))
	store.mu.Lock()
	store.lastSweep = time.Now().Add(-cooldownStoreSweepAfter - time.Millisecond)
	store.mu.Unlock()
	store.lastFailureWithCurrentTTL("active-cluster", activeEndpoint)

	store.mu.Lock()
	_, expiredExists := store.lastFailureByEndpoint[newCooldownKey("old-cluster", expiredEndpoint)]
	store.mu.Unlock()
	assert.False(t, expiredExists)
}

func TestCooldownStoreEvictsOldestEntryWhenCapacityExceeded(t *testing.T) {
	store := newCooldownStore()
	now := time.Now()
	oldestEndpoint := testLLMEndpoint("ep-0", 19000)

	for i := 0; i < maxCooldownStoreEntries; i++ {
		endpoint := testLLMEndpoint(fmt.Sprintf("ep-%d", i), 19000+i)
		if i == 0 {
			oldestEndpoint = endpoint
		}
		store.markFailure("capacity-cluster", endpoint, now.Add(time.Duration(i)*time.Millisecond))
	}

	newestEndpoint := testLLMEndpoint("ep-new", 21000)
	store.markFailure("capacity-cluster", newestEndpoint, now.Add(time.Hour))

	store.mu.Lock()
	_, oldestExists := store.lastFailureByEndpoint[newCooldownKey("capacity-cluster", oldestEndpoint)]
	_, newestExists := store.lastFailureByEndpoint[newCooldownKey("capacity-cluster", newestEndpoint)]
	entryCount := len(store.lastFailureByEndpoint)
	store.mu.Unlock()

	assert.False(t, oldestExists)
	assert.True(t, newestExists)
	assert.Equal(t, maxCooldownStoreEntries, entryCount)
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

func BenchmarkCooldown_EndpointInCooldown(b *testing.B) {
	const clusterName = "llm-cooldown-bench"
	const endpointCount = 100
	// Keep the cooldown TTL far longer than any -benchtime run so every
	// iteration stays on the intended in-cooldown hot path. Otherwise entries
	// could expire mid-benchmark and shift measurement onto the delete+log path.
	const cooldownTTLMillis = int64(24 * time.Hour / time.Millisecond)
	store := newCooldownStore()
	executor := &RequestExecutor{
		clusterName: clusterName,
		cooldowns:   store,
	}
	endpoints := make([]*model.Endpoint, endpointCount)
	for i := range endpoints {
		endpoint := testLLMEndpoint(fmt.Sprintf("ep-%d", i), 19000+i)
		endpoint.LLMMeta.APIKey = fmt.Sprintf("api-key-%d", i)
		endpoint.LLMMeta.HealthCheckInterval = cooldownTTLMillis
		endpoints[i] = endpoint
		store.markFailure(clusterName, endpoint, time.Now())
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		executor.endpointInCooldown(endpoints[i%endpointCount])
	}
}
