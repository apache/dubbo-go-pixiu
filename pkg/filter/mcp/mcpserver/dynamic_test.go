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
	"sync"
	"testing"
	"time"
)

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/filter/mcp/mcpserver/transport"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// =============================================================================
// Test Utilities
// =============================================================================
// createTestToolConfig creates a simple test tool configuration
func createTestToolConfig(name, description string) model.ToolConfig {
	return model.ToolConfig{
		Name:        name,
		Description: description,
		Cluster:     "test-cluster",
		Request: model.RequestConfig{
			Method:  "GET",
			Path:    "/api/test",
			Timeout: "30s",
		},
		Args: []model.ArgConfig{
			{
				Name:        "param",
				Type:        "string",
				In:          "query",
				Description: "Test parameter",
				Required:    true,
			},
		},
	}
}

// createTestMcpServerConfig creates a test MCP server configuration
func createTestMcpServerConfig(tools []model.ToolConfig) *model.McpServerConfig {
	return &model.McpServerConfig{
		ServerInfo: model.ServerInfo{
			Name:        "Test MCP Server",
			Version:     "1.0.0",
			Description: "Test server for unit testing",
		},
		Tools: tools,
	}
}

func toolConfigNames(tools []model.ToolConfig) []string {
	names := make([]string, len(tools))
	for i, tool := range tools {
		names[i] = tool.Name
	}
	return names
}

// =============================================================================
// Singleton Tests
// =============================================================================

func TestSingletonInstances(t *testing.T) {
	ResetGlobalState()

	t.Run("Registry singleton", func(t *testing.T) {
		registry1 := GetOrInitRegistry()
		registry2 := GetOrInitRegistry()
		assert.Same(t, registry1, registry2)
	})

	t.Run("Dynamic consumer singleton", func(t *testing.T) {
		dynamic1 := GetOrInitDynamicConsumer()
		dynamic2 := GetOrInitDynamicConsumer()
		assert.Same(t, dynamic1, dynamic2)
		assert.Same(t, dynamic1.registry, GetOrInitRegistry())
	})
}

func TestSingletonConcurrency(t *testing.T) {
	ResetGlobalState()

	const numGoroutines = 50
	var wg sync.WaitGroup

	t.Run("Registry concurrent access", func(t *testing.T) {
		registries := make([]*ToolRegistry, numGoroutines)
		wg.Add(numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			go func(index int) {
				defer wg.Done()
				registries[index] = GetOrInitRegistry()
			}(i)
		}
		wg.Wait()

		// All should be the same instance
		for i := 1; i < numGoroutines; i++ {
			assert.Same(t, registries[0], registries[i])
		}
	})
}

// =============================================================================
// Configuration Application Tests
// =============================================================================

func TestApplyMcpServerConfig(t *testing.T) {
	t.Run("Basic configuration application", func(t *testing.T) {
		registry := NewToolRegistry()
		sm := transport.NewSessionManager()
		defer sm.Stop()
		sseHandler := transport.NewSSEHandler(sm)
		consumer := NewDynamicConsumer(registry, sm, sseHandler)

		// Test nil config
		err := consumer.ApplyMcpServerConfigByServer("default", nil)
		assert.NoError(t, err)
		assert.Empty(t, registry.ListTools())

		// Test empty config
		config := createTestMcpServerConfig([]model.ToolConfig{})
		err = consumer.ApplyMcpServerConfigByServer("default", config)
		assert.NoError(t, err)
		assert.Empty(t, registry.ListTools())

		// Test with tools
		consumer.ResetDebounceState()
		config = createTestMcpServerConfig([]model.ToolConfig{
			createTestToolConfig("tool1", "First tool"),
			createTestToolConfig("tool2", "Second tool"),
		})
		err = consumer.ApplyMcpServerConfigByServer("default", config)
		assert.NoError(t, err)

		tools := registry.ListTools()
		assert.Len(t, tools, 2)
		toolNames := []string{tools[0].Name, tools[1].Name}
		assert.ElementsMatch(t, []string{"tool1", "tool2"}, toolNames)
	})

	t.Run("Configuration replacement", func(t *testing.T) {
		registry := NewToolRegistry()
		sm := transport.NewSessionManager()
		defer sm.Stop()
		sseHandler := transport.NewSSEHandler(sm)
		consumer := NewDynamicConsumer(registry, sm, sseHandler)

		// Apply first config
		config1 := createTestMcpServerConfig([]model.ToolConfig{
			createTestToolConfig("tool1", "First tool"),
		})
		err := consumer.ApplyMcpServerConfigByServer("default", config1)
		assert.NoError(t, err)
		assert.Len(t, registry.ListTools(), 1)

		// Replace with second config
		consumer.ResetDebounceState()
		config2 := createTestMcpServerConfig([]model.ToolConfig{
			createTestToolConfig("tool2", "Second tool"),
			createTestToolConfig("tool3", "Third tool"),
		})
		err = consumer.ApplyMcpServerConfigByServer("default", config2)
		assert.NoError(t, err)

		tools := registry.ListTools()
		assert.Len(t, tools, 2)
		toolNames := []string{tools[0].Name, tools[1].Name}
		assert.ElementsMatch(t, []string{"tool2", "tool3"}, toolNames)
	})
}

func TestToolRegistryListToolsPreservesReplaceOrder(t *testing.T) {
	registry := NewToolRegistry()
	registry.ReplaceAllTools([]model.ToolConfig{
		createTestToolConfig("tool2", "Second tool"),
		createTestToolConfig("tool1", "First tool"),
		createTestToolConfig("tool3", "Third tool"),
	})

	assert.Equal(t, []string{"tool2", "tool1", "tool3"}, toolConfigNames(registry.ListTools()))
}

func TestDynamicConsumerMergedToolsStableByServerIDAndConfigOrder(t *testing.T) {
	registry := NewToolRegistry()
	sm := transport.NewSessionManager()
	defer sm.Stop()
	sseHandler := transport.NewSSEHandler(sm)
	consumer := NewDynamicConsumer(registry, sm, sseHandler)
	consumer.SetDebounceTime(0)

	err := consumer.ApplyMcpServerConfigByServer("server-b", createTestMcpServerConfig([]model.ToolConfig{
		createTestToolConfig("b2", "B2"),
		createTestToolConfig("b1", "B1"),
	}))
	require.NoError(t, err)

	err = consumer.ApplyMcpServerConfigByServer("server-a", createTestMcpServerConfig([]model.ToolConfig{
		createTestToolConfig("a1", "A1"),
		createTestToolConfig("a2", "A2"),
	}))
	require.NoError(t, err)

	assert.Equal(t, []string{"a1", "a2", "b2", "b1"}, toolConfigNames(registry.ListTools()))
}

func TestApplyMcpServerConfigConcurrent(t *testing.T) {
	ResetGlobalState()
	registry := GetOrInitRegistry()
	consumer := GetOrInitDynamicConsumer()

	const numGoroutines = 10
	var wg sync.WaitGroup

	configs := make([]*model.McpServerConfig, numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		configs[i] = createTestMcpServerConfig([]model.ToolConfig{
			createTestToolConfig("tool-"+string(rune('A'+i)), "Tool "+string(rune('A'+i))),
		})
	}

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			defer wg.Done()
			consumer.ApplyMcpServerConfigByServer("default", configs[index])
		}(i)
	}
	wg.Wait()

	// One config should have won
	tools := registry.ListTools()
	assert.Len(t, tools, 1)
}

// =============================================================================
// Debounce Functionality Tests
// =============================================================================

func TestDebounceFeatures(t *testing.T) {
	t.Run("Content debounce - skip identical configs", func(t *testing.T) {
		registry := NewToolRegistry()
		sm := transport.NewSessionManager()
		defer sm.Stop()
		sseHandler := transport.NewSSEHandler(sm)
		consumer := NewDynamicConsumer(registry, sm, sseHandler)

		config := createTestMcpServerConfig([]model.ToolConfig{
			createTestToolConfig("tool1", "Test tool"),
		})

		// First application
		err := consumer.ApplyMcpServerConfigByServer("default", config)
		assert.NoError(t, err)
		assert.Len(t, registry.ListTools(), 1)

		// Second application with same config - should be skipped
		err = consumer.ApplyMcpServerConfigByServer("default", config)
		assert.NoError(t, err)
		assert.Len(t, registry.ListTools(), 1)

		// Verify debounce info
		info := consumer.GetDebounceInfo()
		assert.Equal(t, 1, info["server_count"])
	})

	t.Run("Time debounce - skip rapid calls", func(t *testing.T) {
		registry := NewToolRegistry()
		sm := transport.NewSessionManager()
		defer sm.Stop()
		sseHandler := transport.NewSSEHandler(sm)
		consumer := NewDynamicConsumer(registry, sm, sseHandler)

		config1 := createTestMcpServerConfig([]model.ToolConfig{
			createTestToolConfig("tool1", "First tool"),
		})
		config2 := createTestMcpServerConfig([]model.ToolConfig{
			createTestToolConfig("tool2", "Second tool"),
		})

		// First application
		err := consumer.ApplyMcpServerConfigByServer("default", config1)
		assert.NoError(t, err)
		tools := registry.ListTools()
		require.Len(t, tools, 1)
		assert.Equal(t, "tool1", tools[0].Name)

		// Immediate second application - should be debounced
		err = consumer.ApplyMcpServerConfigByServer("default", config2)
		assert.NoError(t, err)
		tools = registry.ListTools()
		require.Len(t, tools, 1)
		assert.Equal(t, "tool1", tools[0].Name, "Should still have first tool due to time debounce")
	})

	t.Run("Empty configuration handling", func(t *testing.T) {
		registry := NewToolRegistry()
		sm := transport.NewSessionManager()
		defer sm.Stop()
		sseHandler := transport.NewSSEHandler(sm)
		consumer := NewDynamicConsumer(registry, sm, sseHandler)

		// Add tool first
		config := createTestMcpServerConfig([]model.ToolConfig{
			createTestToolConfig("tool1", "Test tool"),
		})
		err := consumer.ApplyMcpServerConfigByServer("default", config)
		assert.NoError(t, err)
		assert.Len(t, registry.ListTools(), 1)

		// Apply empty config
		consumer.ResetDebounceState()
		emptyConfig := createTestMcpServerConfig([]model.ToolConfig{})
		err = consumer.ApplyMcpServerConfigByServer("default", emptyConfig)
		assert.NoError(t, err)
		assert.Empty(t, registry.ListTools())

		// Verify empty configuration
		info := consumer.GetDebounceInfo()
		assert.Equal(t, 1, info["server_count"])
	})
}

func TestDebounceConfiguration(t *testing.T) {
	registry := NewToolRegistry()
	sm := transport.NewSessionManager()
	defer sm.Stop()
	sseHandler := transport.NewSSEHandler(sm)
	consumer := NewDynamicConsumer(registry, sm, sseHandler)

	// Test default debounce time
	info := consumer.GetDebounceInfo()
	assert.Equal(t, DefaultDebounceTime.String(), info["debounce_time"])

	// Test custom debounce time
	customTime := 1000 * time.Millisecond
	consumer.SetDebounceTime(customTime)
	info = consumer.GetDebounceInfo()
	assert.Equal(t, customTime.String(), info["debounce_time"])

	// Test invalid debounce time (negative)
	consumer.SetDebounceTime(-100 * time.Millisecond)
	info = consumer.GetDebounceInfo()
	assert.Equal(t, customTime.String(), info["debounce_time"], "Negative time should be ignored")

	// Test reset debounce state
	config := createTestMcpServerConfig([]model.ToolConfig{
		createTestToolConfig("tool1", "Test tool"),
	})
	consumer.ApplyMcpServerConfigByServer("default", config)

	info = consumer.GetDebounceInfo()
	assert.Equal(t, 1, info["server_count"])

	consumer.ResetDebounceState()
	info = consumer.GetDebounceInfo()
	assert.Equal(t, 0, info["server_count"])
}

func TestFingerprintCalculation(t *testing.T) {
	registry := NewToolRegistry()
	sm := transport.NewSessionManager()
	defer sm.Stop()
	sseHandler := transport.NewSSEHandler(sm)
	consumer := NewDynamicConsumer(registry, sm, sseHandler)

	// Empty tools
	fingerprint1 := consumer.calculateFingerprint([]model.ToolConfig{})
	assert.Equal(t, EmptyFingerprint, fingerprint1)

	// Single tool
	tool1 := createTestToolConfig("tool1", "First tool")
	fingerprint2 := consumer.calculateFingerprint([]model.ToolConfig{tool1})
	assert.NotEqual(t, EmptyFingerprint, fingerprint2)
	assert.Len(t, fingerprint2, 8, "Fingerprint should be 8 characters")

	// Multiple tools - order should not matter
	tool2 := createTestToolConfig("tool2", "Second tool")
	fingerprint3 := consumer.calculateFingerprint([]model.ToolConfig{tool1, tool2})
	fingerprint4 := consumer.calculateFingerprint([]model.ToolConfig{tool2, tool1})
	assert.Equal(t, fingerprint3, fingerprint4, "Tool order should not affect fingerprint")

	// Different tools should have different fingerprints
	tool3 := createTestToolConfig("tool3", "Third tool")
	fingerprint5 := consumer.calculateFingerprint([]model.ToolConfig{tool3})
	assert.NotEqual(t, fingerprint2, fingerprint5, "Different tools should have different fingerprints")

	// Same tool should have same fingerprint
	fingerprint6 := consumer.calculateFingerprint([]model.ToolConfig{tool1})
	assert.Equal(t, fingerprint2, fingerprint6, "Same tool should have same fingerprint")

	// Duplicate identities should still be order-independent because the full
	// serialized body is part of the sort key.
	duplicateA := createTestToolConfig("dup", "First duplicate")
	duplicateB := createTestToolConfig("dup", "Second duplicate")
	assert.Equal(t,
		consumer.calculateFingerprint([]model.ToolConfig{duplicateA, duplicateB}),
		consumer.calculateFingerprint([]model.ToolConfig{duplicateB, duplicateA}),
		"Duplicate tool identities should still hash deterministically")

	// Metadata-only changes must affect the fingerprint so dynamic updates are
	// not skipped when routing tags/risk change.
	toolWithLowRisk := createTestToolConfig("tool1", "First tool")
	toolWithLowRisk.Meta = &model.ToolMeta{Risk: "low", Tags: []string{"safe"}}
	toolWithHighRisk := createTestToolConfig("tool1", "First tool")
	toolWithHighRisk.Meta = &model.ToolMeta{Risk: "high", Tags: []string{"admin"}}
	assert.NotEqual(t,
		consumer.calculateFingerprint([]model.ToolConfig{toolWithLowRisk}),
		consumer.calculateFingerprint([]model.ToolConfig{toolWithHighRisk}),
		"Meta changes should affect fingerprint")

	// Request and argument shape changes also affect backend behavior and must
	// be part of the dynamic configuration hash.
	toolWithPathA := createTestToolConfig("tool1", "First tool")
	toolWithPathB := createTestToolConfig("tool1", "First tool")
	toolWithPathB.Request.Path = "/api/other/{param}"
	assert.NotEqual(t,
		consumer.calculateFingerprint([]model.ToolConfig{toolWithPathA}),
		consumer.calculateFingerprint([]model.ToolConfig{toolWithPathB}),
		"Request path changes should affect fingerprint")

	toolWithArgA := createTestToolConfig("tool1", "First tool")
	toolWithArgB := createTestToolConfig("tool1", "First tool")
	toolWithArgB.Args[0].Required = false
	assert.NotEqual(t,
		consumer.calculateFingerprint([]model.ToolConfig{toolWithArgA}),
		consumer.calculateFingerprint([]model.ToolConfig{toolWithArgB}),
		"Argument changes should affect fingerprint")
}

func TestApplyMcpServerConfig_InvalidRiskRejected(t *testing.T) {
	registry := NewToolRegistry()
	sm := transport.NewSessionManager()
	defer sm.Stop()
	sseHandler := transport.NewSSEHandler(sm)
	consumer := NewDynamicConsumer(registry, sm, sseHandler)

	validTool := createTestToolConfig("tool1", "First tool")
	validTool.Meta = &model.ToolMeta{Risk: "low"}
	err := consumer.ApplyMcpServerConfigByServer("default", createTestMcpServerConfig([]model.ToolConfig{validTool}))
	require.NoError(t, err)
	require.Len(t, registry.ListTools(), 1)

	invalidTool := createTestToolConfig("tool2", "Second tool")
	invalidTool.Meta = &model.ToolMeta{Risk: "hihg"}
	err = consumer.ApplyMcpServerConfigByServer("default", createTestMcpServerConfig([]model.ToolConfig{invalidTool}))
	assert.ErrorContains(t, err, "invalid mcp tool router metadata")
	assert.ErrorContains(t, err, "unsupported risk")

	tools := registry.ListTools()
	require.Len(t, tools, 1)
	assert.Equal(t, "tool1", tools[0].Name)
}

func TestApplyMcpServerConfig_MetadataChangeIsNotSkipped(t *testing.T) {
	registry := NewToolRegistry()
	sm := transport.NewSessionManager()
	defer sm.Stop()
	sseHandler := transport.NewSSEHandler(sm)
	consumer := NewDynamicConsumer(registry, sm, sseHandler)
	consumer.SetDebounceTime(0)

	lowRiskTool := createTestToolConfig("tool1", "First tool")
	lowRiskTool.Meta = &model.ToolMeta{Risk: "low"}
	err := consumer.ApplyMcpServerConfigByServer("default", createTestMcpServerConfig([]model.ToolConfig{lowRiskTool}))
	require.NoError(t, err)

	highRiskTool := createTestToolConfig("tool1", "First tool")
	highRiskTool.Meta = &model.ToolMeta{Risk: "high"}
	err = consumer.ApplyMcpServerConfigByServer("default", createTestMcpServerConfig([]model.ToolConfig{highRiskTool}))
	require.NoError(t, err)

	tools := registry.ListTools()
	require.Len(t, tools, 1)
	require.NotNil(t, tools[0].Meta)
	assert.Equal(t, "high", tools[0].Meta.Risk)
}

// =============================================================================
// Integration Tests
// =============================================================================

func TestIntegration(t *testing.T) {
	ResetGlobalState()

	registry := GetOrInitRegistry()
	consumer := GetOrInitDynamicConsumer()

	// Verify initial state
	assert.Empty(t, registry.ListTools())

	// Apply first configuration
	config1 := createTestMcpServerConfig([]model.ToolConfig{
		createTestToolConfig("tool1", "First tool"),
	})
	err := consumer.ApplyMcpServerConfigByServer("default", config1)
	assert.NoError(t, err)

	tools := registry.ListTools()
	require.Len(t, tools, 1)
	assert.Equal(t, "tool1", tools[0].Name)

	// Apply second configuration
	consumer.ResetDebounceState()
	config2 := createTestMcpServerConfig([]model.ToolConfig{
		createTestToolConfig("tool2", "Second tool"),
		createTestToolConfig("tool3", "Third tool"),
	})
	err = consumer.ApplyMcpServerConfigByServer("default", config2)
	assert.NoError(t, err)

	tools = registry.ListTools()
	assert.Len(t, tools, 2)
	toolNames := []string{tools[0].Name, tools[1].Name}
	assert.ElementsMatch(t, []string{"tool2", "tool3"}, toolNames)
}

// =============================================================================
// Benchmark Tests
// =============================================================================

func BenchmarkApplyMcpServerConfig(b *testing.B) {
	registry := NewToolRegistry()
	sm := transport.NewSessionManager()
	defer sm.Stop()
	sseHandler := transport.NewSSEHandler(sm)
	consumer := NewDynamicConsumer(registry, sm, sseHandler)

	config := createTestMcpServerConfig([]model.ToolConfig{
		createTestToolConfig("tool1", "First tool"),
		createTestToolConfig("tool2", "Second tool"),
		createTestToolConfig("tool3", "Third tool"),
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		consumer.ApplyMcpServerConfigByServer("default", config)
	}
}
