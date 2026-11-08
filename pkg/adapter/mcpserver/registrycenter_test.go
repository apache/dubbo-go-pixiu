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

package mcpserver

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/mcpserver/registry"
	filtermcp "github.com/apache/dubbo-go-pixiu/pkg/filter/mcp/mcpserver"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

type testPublicationSink struct {
	runtimeID string
	applied   int
	removed   int
	bySource  []filtermcp.ServerSource
	err       error
	onApply   func(source filtermcp.ServerSource, cfg *model.McpServerConfig)
}

type testController struct {
	runCount   atomic.Int32
	closeCount atomic.Int32
}

type concurrentPublicationSink struct {
	runtimeID string
	mu        sync.Mutex
	applied   int
	removed   int
}

func (c *testController) Run(ctx context.Context, _ time.Duration) error {
	c.runCount.Add(1)
	<-ctx.Done()
	return nil
}

func (c *testController) Close() error {
	c.closeCount.Add(1)
	return nil
}

func (s *testPublicationSink) RuntimeID() string {
	return s.runtimeID
}

func (s *testPublicationSink) ApplyMcpServerConfigByServer(_ string, _ *model.McpServerConfig) error {
	s.applied++
	return s.err
}

func (s *testPublicationSink) ApplyMcpServerConfigBySource(source filtermcp.ServerSource, cfg *model.McpServerConfig) error {
	s.applied++
	s.bySource = append(s.bySource, source.Normalize())
	if s.onApply != nil {
		s.onApply(source.Normalize(), cfg)
	}
	return s.err
}

func (s *testPublicationSink) RemoveAllMcpServerConfigs() error {
	s.removed++
	return nil
}

func (s *concurrentPublicationSink) RuntimeID() string {
	return s.runtimeID
}

func (s *concurrentPublicationSink) ApplyMcpServerConfigByServer(_ string, _ *model.McpServerConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.applied++
	return nil
}

func (s *concurrentPublicationSink) ApplyMcpServerConfigBySource(_ filtermcp.ServerSource, _ *model.McpServerConfig) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.applied++
	return nil
}

func (s *concurrentPublicationSink) RemoveAllMcpServerConfigs() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.removed++
	return nil
}

func (s *concurrentPublicationSink) counts() (int, int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.applied, s.removed
}

func withPublicationSinkLookup(t *testing.T, fn func() (filtermcp.ServerPublicationSink, error)) {
	t.Helper()
	old := serverPublicationSinkForSingleRuntime
	serverPublicationSinkForSingleRuntime = fn
	t.Cleanup(func() {
		serverPublicationSinkForSingleRuntime = old
	})
}

func TestAdapterBindPublicationSinkRequiresRuntime(t *testing.T) {
	filtermcp.ResetGlobalState()
	defer filtermcp.ResetGlobalState()

	a := &Adapter{}
	sink, err := a.bindPublicationSink()

	assert.ErrorIs(t, err, filtermcp.ErrDynamicConsumerUnavailable)
	assert.Nil(t, sink)
	assert.Nil(t, a.sink)
}

func TestAdapterApplyServerConfigEventNoBoundSinkDoesNothing(t *testing.T) {
	filtermcp.ResetGlobalState()
	defer filtermcp.ResetGlobalState()

	sink := &recordingEndpointSink{}
	a := &Adapter{endpoints: newEndpointReconciler(sink)}

	a.applyServerConfigEvent("nacos", "server-a", &model.McpServerConfig{
		Tools: []model.ToolConfig{{Name: "tool", Cluster: "cluster", BackendURL: "http://127.0.0.1:8080"}},
	})

	assert.Empty(t, a.endpoints.published)
	assert.Empty(t, sink.ops)
}

func TestAdapterApplyServerConfigEventUsesBoundRuntime(t *testing.T) {
	sink := &testPublicationSink{runtimeID: "runtime-1"}
	withPublicationSinkLookup(t, func() (filtermcp.ServerPublicationSink, error) {
		return sink, nil
	})
	endpointSink := &recordingEndpointSink{}
	a := &Adapter{endpoints: newEndpointReconciler(endpointSink)}

	a.applyServerConfigEvent("nacos", "server-a", &model.McpServerConfig{
		Tools: []model.ToolConfig{{Name: "tool", Cluster: "cluster", BackendURL: "http://127.0.0.1:8080"}},
	})

	assert.Equal(t, 1, sink.applied)
	assert.Equal(t, []filtermcp.ServerSource{filtermcp.NewServerSource("nacos", "server-a")}, sink.bySource)
	require.Len(t, endpointSink.ops, 1)
	assert.Equal(t, "mcp/nacos/server-a/tool", endpointSink.ops[0].id)
}

func TestAdapterApplyBeforeRuntimeThenEventAfterRuntimePublishes(t *testing.T) {
	runtimeErr := filtermcp.ErrDynamicConsumerUnavailable
	publishedSink := &testPublicationSink{runtimeID: "runtime-1"}
	withPublicationSinkLookup(t, func() (filtermcp.ServerPublicationSink, error) {
		if runtimeErr != nil {
			return nil, runtimeErr
		}
		return publishedSink, nil
	})
	endpointSink := &recordingEndpointSink{}
	a := &Adapter{endpoints: newEndpointReconciler(endpointSink)}

	a.applyServerConfigEvent("nacos", "server-a", &model.McpServerConfig{
		Tools: []model.ToolConfig{{Name: "before", Cluster: "cluster", BackendURL: "http://127.0.0.1:8080"}},
	})
	assert.Empty(t, endpointSink.ops)
	assert.Equal(t, 0, publishedSink.applied)

	runtimeErr = nil
	a.applyServerConfigEvent("nacos", "server-a", &model.McpServerConfig{
		Tools: []model.ToolConfig{{Name: "after", Cluster: "cluster", BackendURL: "http://127.0.0.1:8081"}},
	})

	require.Len(t, endpointSink.ops, 1)
	assert.Equal(t, "mcp/nacos/server-a/after", endpointSink.ops[0].id)
	assert.Equal(t, 1, publishedSink.applied)
}

func TestAdapterApplyServerConfigEventAmbiguousRuntimeFailsClosed(t *testing.T) {
	withPublicationSinkLookup(t, func() (filtermcp.ServerPublicationSink, error) {
		return nil, filtermcp.ErrDynamicConsumerAmbiguous
	})

	endpointSink := &recordingEndpointSink{}
	a := &Adapter{endpoints: newEndpointReconciler(endpointSink)}
	a.applyServerConfigEvent("nacos", "server-a", &model.McpServerConfig{
		Tools: []model.ToolConfig{{Name: "tool", Cluster: "cluster", BackendURL: "http://127.0.0.1:8080"}},
	})

	assert.Empty(t, endpointSink.ops)
}

func TestAdapterApplyServerConfigEventValidationFailsAtomically(t *testing.T) {
	sink := &testPublicationSink{runtimeID: "runtime-1"}
	withPublicationSinkLookup(t, func() (filtermcp.ServerPublicationSink, error) {
		return sink, nil
	})

	endpointSink := &recordingEndpointSink{}
	a := &Adapter{endpoints: newEndpointReconciler(endpointSink)}
	a.applyServerConfigEvent("nacos", "server-a", &model.McpServerConfig{
		Tools: []model.ToolConfig{{Name: "bad", Cluster: "cluster", BackendURL: "://bad-url"}},
	})

	assert.Empty(t, endpointSink.ops)
	assert.Equal(t, 0, sink.applied)
}

func TestAdapterApplyServerConfigEventUsesPreparedEndpointPlan(t *testing.T) {
	sink := &testPublicationSink{runtimeID: "runtime-1"}
	sink.onApply = func(_ filtermcp.ServerSource, cfg *model.McpServerConfig) {
		cfg.Tools[0].BackendURL = "://mutated-after-catalog-apply"
	}
	withPublicationSinkLookup(t, func() (filtermcp.ServerPublicationSink, error) {
		return sink, nil
	})

	endpointSink := &recordingEndpointSink{}
	a := &Adapter{endpoints: newEndpointReconciler(endpointSink)}
	cfg := &model.McpServerConfig{
		Tools: []model.ToolConfig{{Name: "tool", Cluster: "cluster", BackendURL: "http://127.0.0.1:8080"}},
	}

	a.applyServerConfigEvent("nacos", "server-a", cfg)

	assert.Equal(t, 1, sink.applied)
	require.Len(t, endpointSink.ops, 1)
	assert.Equal(t, "set", endpointSink.ops[0].action)
	assert.Equal(t, "127.0.0.1:8080", endpointSink.ops[0].address)
	assert.NotEmpty(t, a.publishedSources, "source is tracked only after catalog and endpoint desired state are both committed")
}

func TestAdapterApplyServerConfigEventCatalogFailureSkipsEndpoints(t *testing.T) {
	sink := &testPublicationSink{runtimeID: "runtime-1", err: errors.New("boom")}
	withPublicationSinkLookup(t, func() (filtermcp.ServerPublicationSink, error) {
		return sink, nil
	})
	endpointSink := &recordingEndpointSink{}
	a := &Adapter{
		endpoints: newEndpointReconciler(endpointSink),
	}

	a.applyServerConfigEvent("nacos", "server-a", &model.McpServerConfig{
		Tools: []model.ToolConfig{{Name: "tool", Cluster: "cluster", BackendURL: "http://127.0.0.1:8080"}},
	})

	assert.Empty(t, endpointSink.ops)
	assert.Equal(t, 1, sink.applied)
}

func TestAdapterTombstoneRemovesCatalogAndEndpoints(t *testing.T) {
	sink := &testPublicationSink{runtimeID: "runtime-1"}
	withPublicationSinkLookup(t, func() (filtermcp.ServerPublicationSink, error) {
		return sink, nil
	})
	endpointSink := &recordingEndpointSink{}
	a := &Adapter{endpoints: newEndpointReconciler(endpointSink)}

	a.applyServerConfigEvent("nacos", "server-a", &model.McpServerConfig{
		Tools: []model.ToolConfig{{Name: "tool", Cluster: "cluster", BackendURL: "http://127.0.0.1:8080"}},
	})
	a.applyServerConfigEvent("nacos", "server-a", nil)

	require.Len(t, endpointSink.ops, 2)
	assert.Equal(t, "set", endpointSink.ops[0].action)
	assert.Equal(t, "delete", endpointSink.ops[1].action)
	assert.Equal(t, endpointSink.ops[0].id, endpointSink.ops[1].id)
	assert.Equal(t, 2, sink.applied)
	assert.Empty(t, a.publishedSources)
}

func TestAdapterStopRemovesPublishedDynamicState(t *testing.T) {
	sink := &testPublicationSink{runtimeID: "runtime-1"}
	withPublicationSinkLookup(t, func() (filtermcp.ServerPublicationSink, error) {
		return sink, nil
	})

	endpointSink := &recordingEndpointSink{}
	a := &Adapter{
		controllers:      map[string]registry.Controller{},
		endpoints:        newEndpointReconciler(endpointSink),
		publishedSources: map[string]filtermcp.ServerSource{},
	}
	a.applyServerConfigEvent("nacos", "server-a", &model.McpServerConfig{
		Tools: []model.ToolConfig{{Name: "tool", Cluster: "cluster", BackendURL: "http://127.0.0.1:8080"}},
	})
	require.Len(t, endpointSink.ops, 1)

	a.Stop()

	require.Len(t, endpointSink.ops, 2)
	assert.Equal(t, "delete", endpointSink.ops[1].action)
	assert.Equal(t, 1, sink.removed)
}

func TestAdapterStartStopManagesAllControllers(t *testing.T) {
	ctrlA := &testController{}
	ctrlB := &testController{}
	sink := &testPublicationSink{runtimeID: "runtime-1"}
	withPublicationSinkLookup(t, func() (filtermcp.ServerPublicationSink, error) {
		return sink, nil
	})
	a := &Adapter{
		controllers: map[string]registry.Controller{
			"registry-a": ctrlA,
			"registry-b": ctrlB,
		},
		endpoints: newEndpointReconciler(&recordingEndpointSink{}),
	}

	a.Start()
	require.Eventually(t, func() bool {
		return ctrlA.runCount.Load() == 1 && ctrlB.runCount.Load() == 1
	}, time.Second, 10*time.Millisecond)

	a.Stop()

	assert.Equal(t, int32(1), ctrlA.closeCount.Load())
	assert.Equal(t, int32(1), ctrlB.closeCount.Load())
	assert.Equal(t, 1, sink.removed)
}

func TestAdapterApplyServerConfigEventResolvesSinkLazily(t *testing.T) {
	filtermcp.ResetGlobalState()
	defer filtermcp.ResetGlobalState()

	runtimeErr := filtermcp.ErrDynamicConsumerUnavailable
	sink := &testPublicationSink{runtimeID: "runtime-1"}
	withPublicationSinkLookup(t, func() (filtermcp.ServerPublicationSink, error) {
		if runtimeErr != nil {
			return nil, runtimeErr
		}
		return sink, nil
	})
	endpointSink := &recordingEndpointSink{}
	a := &Adapter{endpoints: newEndpointReconciler(endpointSink)}

	a.applyServerConfigEvent("nacos", "server-a", &model.McpServerConfig{
		Tools: []model.ToolConfig{{Name: "tool", Cluster: "cluster", BackendURL: "http://127.0.0.1:8080"}},
	})
	assert.Empty(t, a.endpoints.published)
	assert.Empty(t, endpointSink.ops)

	runtimeErr = nil
	a.applyServerConfigEvent("nacos", "server-a", &model.McpServerConfig{
		Tools: []model.ToolConfig{{Name: "tool", Cluster: "cluster", BackendURL: "http://127.0.0.1:8080"}},
	})

	assert.Equal(t, 1, sink.applied)
	assert.NotEmpty(t, a.endpoints.published)
	assert.NotEmpty(t, endpointSink.ops)
}

func TestAdapterApplyStopAndRegistryEventConcurrent(t *testing.T) {
	sink := &concurrentPublicationSink{runtimeID: "runtime-1"}
	withPublicationSinkLookup(t, func() (filtermcp.ServerPublicationSink, error) {
		return sink, nil
	})

	a := &Adapter{
		id:               "adapter-test",
		cfg:              &AdapterConfig{Registries: map[string]model.Registry{}},
		controllers:      map[string]registry.Controller{},
		endpoints:        newEndpointReconciler(&recordingEndpointSink{}),
		publishedSources: map[string]filtermcp.ServerSource{},
	}
	cfg := &model.McpServerConfig{
		Tools: []model.ToolConfig{{Name: "tool", Cluster: "cluster", BackendURL: "http://127.0.0.1:8080"}},
	}

	var wg sync.WaitGroup
	errCh := make(chan error, 25)
	for i := 0; i < 25; i++ {
		wg.Add(3)
		go func() {
			defer wg.Done()
			a.applyServerConfigEvent("nacos", "server-a", cfg)
		}()
		go func() {
			defer wg.Done()
			errCh <- a.Apply()
		}()
		go func() {
			defer wg.Done()
			a.Stop()
		}()
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		require.NoError(t, err)
	}

	applied, removed := sink.counts()
	assert.Greater(t, applied, 0)
	assert.Greater(t, removed, 0)
}
