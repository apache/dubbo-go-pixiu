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
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// =============================================================================
// Test Utilities
// =============================================================================
// resetSingletons resets singleton state for testing
func resetSingletons() {
	globalRegistry = nil
	globalDynamic = nil
	registryOnce = sync.Once{}
	dynamicOnce = sync.Once{}
}

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

// =============================================================================
// Singleton Tests
// =============================================================================

func TestSingletonInstances(t *testing.T) {
	resetSingletons()

	t.Run("Registry singleton", func(t *testing.T) {
		registry1 := GetOrInitRegistry()
		registry2 := GetOrInitRegistry()
		assert.Same(t, registry1, registry2)
	})

	t.Run("Dynamic consumer singleton", func(t *testing.T) {
		dynamic1 := GetOrInitDynamic()
		dynamic2 := GetOrInitDynamic()
		assert.Same(t, dynamic1, dynamic2)
		assert.Same(t, dynamic1.registry, GetOrInitRegistry())
	})
}

func TestSingletonConcurrency(t *testing.T) {
	resetSingletons()

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
		consumer := NewDynamicConsumer(registry)

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
		consumer := NewDynamicConsumer(registry)

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

func TestApplyMcpServerConfigConcurrent(t *testing.T) {
	resetSingletons()
	registry := GetOrInitRegistry()
	consumer := GetOrInitDynamic()

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
		consumer := NewDynamicConsumer(registry)

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
		consumer := NewDynamicConsumer(registry)

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
		consumer := NewDynamicConsumer(registry)

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
	consumer := NewDynamicConsumer(registry)

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
	consumer := NewDynamicConsumer(registry)

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
}

// =============================================================================
// Integration Tests
// =============================================================================

func TestIntegration(t *testing.T) {
	resetSingletons()

	registry := GetOrInitRegistry()
	consumer := GetOrInitDynamic()

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
	consumer := NewDynamicConsumer(registry)

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
