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
)

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// Test constants
const (
	testToolName1        = "test-tool-1"
	testToolName2        = "test-tool-2"
	testToolName3        = "test-tool-3"
	testToolDescription1 = "First test tool"
	testToolDescription2 = "Second test tool"
	testToolDescription3 = "Third test tool"
	testClusterName      = "test-cluster"
	testRequestPath      = "/api/test"
	testRequestMethod    = "GET"
	testArgName          = "test-param"
	testArgType          = "string"
	testArgIn            = "query"
	testArgDescription   = "Test parameter"
	defaultTimeout       = "30s"
)

// Helper function to create test tool configurations
func createTestToolConfig(name, description string) model.ToolConfig {
	return model.ToolConfig{
		Name:        name,
		Description: description,
		Cluster:     testClusterName,
		Request: model.RequestConfig{
			Method:  testRequestMethod,
			Path:    testRequestPath,
			Timeout: defaultTimeout,
		},
		Args: []model.ArgConfig{
			{
				Name:        testArgName,
				Type:        testArgType,
				In:          testArgIn,
				Description: testArgDescription,
				Required:    true,
			},
		},
	}
}

// Helper function to create test MCP server configuration
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

// Helper function to reset singleton state for testing
func resetSingletons() {
	globalRegistry = nil
	globalDynamic = nil
	registryOnce = sync.Once{}
	dynamicOnce = sync.Once{}
}

func TestGetOrInitRegistry(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{
			name: "should return singleton registry instance",
			testFunc: func(t *testing.T) {
				resetSingletons()

				registry1 := GetOrInitRegistry()
				registry2 := GetOrInitRegistry()

				assert.NotNil(t, registry1)
				assert.NotNil(t, registry2)
				assert.Same(t, registry1, registry2, "GetOrInitRegistry should return the same instance")
			},
		},
		{
			name: "should initialize registry only once",
			testFunc: func(t *testing.T) {
				resetSingletons()

				registry := GetOrInitRegistry()
				require.NotNil(t, registry)

				// Add a tool to verify it's the same instance
				testTool := createTestToolConfig(testToolName1, testToolDescription1)
				registry.RegisterTool(testTool)

				// Get registry again and verify it contains the tool
				registry2 := GetOrInitRegistry()
				tools := registry2.ListTools()
				assert.Len(t, tools, 1)
				assert.Equal(t, testToolName1, tools[0].Name)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// TestGetOrInitDynamic tests the singleton dynamic consumer instance
func TestGetOrInitDynamic(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{
			name: "should return singleton dynamic consumer instance",
			testFunc: func(t *testing.T) {
				resetSingletons()

				dynamic1 := GetOrInitDynamic()
				dynamic2 := GetOrInitDynamic()

				assert.NotNil(t, dynamic1)
				assert.NotNil(t, dynamic2)
				assert.Same(t, dynamic1, dynamic2, "GetOrInitDynamic should return the same instance")
			},
		},
		{
			name: "should use same registry instance",
			testFunc: func(t *testing.T) {
				resetSingletons()

				registry := GetOrInitRegistry()
				dynamic := GetOrInitDynamic()

				assert.NotNil(t, dynamic)
				assert.Same(t, registry, dynamic.registry, "DynamicConsumer should use the singleton registry")
			},
		},
		{
			name: "should initialize dynamic consumer only once",
			testFunc: func(t *testing.T) {
				resetSingletons()

				dynamic1 := GetOrInitDynamic()
				dynamic2 := GetOrInitDynamic()

				// Verify they share the same registry
				assert.Same(t, dynamic1.registry, dynamic2.registry)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}

// TestGetOrInitRegistry_Concurrent Concurrent tests to ensure thread safety
func TestGetOrInitRegistry_Concurrent(t *testing.T) {
	resetSingletons()

	const numGoroutines = 100
	registries := make([]*ToolRegistry, numGoroutines)
	var wg sync.WaitGroup

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			defer wg.Done()
			registries[index] = GetOrInitRegistry()
		}(i)
	}

	wg.Wait()

	// Verify all registries are the same instance
	firstRegistry := registries[0]
	assert.NotNil(t, firstRegistry)

	for i := 1; i < numGoroutines; i++ {
		assert.Same(t, firstRegistry, registries[i], "All registries should be the same instance")
	}
}

// TestGetOrInitDynamic_Concurrent Concurrent tests to ensure thread safety
func TestGetOrInitDynamic_Concurrent(t *testing.T) {
	resetSingletons()

	const numGoroutines = 100
	dynamics := make([]*DynamicConsumer, numGoroutines)
	var wg sync.WaitGroup

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			defer wg.Done()
			dynamics[index] = GetOrInitDynamic()
		}(i)
	}

	wg.Wait()

	// Verify all dynamic consumers are the same instance
	firstDynamic := dynamics[0]
	assert.NotNil(t, firstDynamic)

	for i := 1; i < numGoroutines; i++ {
		assert.Same(t, firstDynamic, dynamics[i], "All dynamic consumers should be the same instance")
		assert.Same(t, firstDynamic.registry, dynamics[i].registry, "All registries should be the same instance")
	}
}

// TestNewDynamicConsumer tests the constructor of DynamicConsumer
func TestNewDynamicConsumer(t *testing.T) {
	tests := []struct {
		name     string
		registry *ToolRegistry
		expected bool
	}{
		{
			name:     "should create dynamic consumer with valid registry",
			registry: NewToolRegistry(),
			expected: true,
		},
		{
			name:     "should create dynamic consumer with nil registry",
			registry: nil,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			consumer := NewDynamicConsumer(tt.registry)

			if tt.expected {
				assert.NotNil(t, consumer)
				assert.Equal(t, tt.registry, consumer.registry)
			} else {
				assert.Nil(t, consumer)
			}
		})
	}
}

// TestApplyMcpServerConfig tests applying MCP server configurations
func TestDynamicConsumer_ApplyMcpServerConfig(t *testing.T) {
	tests := []struct {
		name          string
		initialTools  []model.ToolConfig
		configToApply *model.McpServerConfig
		expectedTools []string
		expectError   bool
	}{
		{
			name:          "should handle nil config",
			initialTools:  []model.ToolConfig{},
			configToApply: nil,
			expectedTools: []string{},
			expectError:   false,
		},
		{
			name: "should replace all tools with new configuration",
			initialTools: []model.ToolConfig{
				createTestToolConfig(testToolName1, testToolDescription1),
			},
			configToApply: createTestMcpServerConfig([]model.ToolConfig{
				createTestToolConfig(testToolName2, testToolDescription2),
				createTestToolConfig(testToolName3, testToolDescription3),
			}),
			expectedTools: []string{testToolName2, testToolName3},
			expectError:   false,
		},
		{
			name: "should clear all tools when empty config provided",
			initialTools: []model.ToolConfig{
				createTestToolConfig(testToolName1, testToolDescription1),
				createTestToolConfig(testToolName2, testToolDescription2),
			},
			configToApply: createTestMcpServerConfig([]model.ToolConfig{}),
			expectedTools: []string{},
			expectError:   false,
		},
		{
			name:         "should handle empty initial tools",
			initialTools: []model.ToolConfig{},
			configToApply: createTestMcpServerConfig([]model.ToolConfig{
				createTestToolConfig(testToolName1, testToolDescription1),
			}),
			expectedTools: []string{testToolName1},
			expectError:   false,
		},
		{
			name:         "should handle multiple tools with same name (last one wins)",
			initialTools: []model.ToolConfig{},
			configToApply: createTestMcpServerConfig([]model.ToolConfig{
				createTestToolConfig(testToolName1, testToolDescription1),
				createTestToolConfig(testToolName1, testToolDescription2), // Same name
			}),
			expectedTools: []string{testToolName1}, // Only one should remain
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create fresh registry and consumer for each test
			registry := NewToolRegistry()
			consumer := NewDynamicConsumer(registry)

			// Setup initial state
			for _, tool := range tt.initialTools {
				registry.RegisterTool(tool)
			}

			// Apply configuration
			err := consumer.ApplyMcpServerConfig(tt.configToApply)

			// Verify error expectation
			if tt.expectError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			// Verify tools state
			actualTools := registry.ListTools()
			actualToolNames := make([]string, len(actualTools))
			for i, tool := range actualTools {
				actualToolNames[i] = tool.Name
			}

			assert.ElementsMatch(t, tt.expectedTools, actualToolNames)
		})
	}
}

// Integration test to verify end-to-end behavior
func TestDynamicConsumer_ApplyMcpServerConfig_Integration(t *testing.T) {
	// Test integration with singleton pattern
	resetSingletons()

	// Get singleton instances
	registry := GetOrInitRegistry()
	consumer := GetOrInitDynamic()

	// Verify initial state
	assert.Empty(t, registry.ListTools())

	// Apply first configuration
	config1 := createTestMcpServerConfig([]model.ToolConfig{
		createTestToolConfig(testToolName1, testToolDescription1),
	})

	err := consumer.ApplyMcpServerConfig(config1)
	assert.NoError(t, err)

	tools := registry.ListTools()
	assert.Len(t, tools, 1)
	assert.Equal(t, testToolName1, tools[0].Name)

	// Apply second configuration (should replace first)
	config2 := createTestMcpServerConfig([]model.ToolConfig{
		createTestToolConfig(testToolName2, testToolDescription2),
		createTestToolConfig(testToolName3, testToolDescription3),
	})

	err = consumer.ApplyMcpServerConfig(config2)
	assert.NoError(t, err)

	tools = registry.ListTools()
	assert.Len(t, tools, 2)

	toolNames := []string{tools[0].Name, tools[1].Name}
	assert.ElementsMatch(t, []string{testToolName2, testToolName3}, toolNames)
}

// Concurrent test to ensure thread safety during configuration application
func TestDynamicConsumer_ApplyMcpServerConfig_Concurrent(t *testing.T) {
	resetSingletons()

	// Get singleton instances
	registry := GetOrInitRegistry()
	consumer := GetOrInitDynamic()

	const numGoroutines = 10
	var wg sync.WaitGroup

	configs := make([]*model.McpServerConfig, numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		configs[i] = createTestMcpServerConfig([]model.ToolConfig{
			createTestToolConfig(
				"tool-"+string(rune('A'+i)),
				"Description for tool "+string(rune('A'+i)),
			),
		})
	}

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(index int) {
			defer wg.Done()
			err := consumer.ApplyMcpServerConfig(configs[index])
			assert.NoError(t, err)
		}(i)
	}

	wg.Wait()

	// Verify final state (exactly one configuration should have won)
	tools := registry.ListTools()
	assert.Len(t, tools, 1, "Only one configuration should remain after concurrent updates")
}

// Benchmark tests
func BenchmarkGetOrInitRegistry(b *testing.B) {
	resetSingletons()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetOrInitRegistry()
	}
}

func BenchmarkGetOrInitDynamic(b *testing.B) {
	resetSingletons()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GetOrInitDynamic()
	}
}

func BenchmarkApplyMcpServerConfig(b *testing.B) {
	registry := NewToolRegistry()
	consumer := NewDynamicConsumer(registry)

	config := createTestMcpServerConfig([]model.ToolConfig{
		createTestToolConfig(testToolName1, testToolDescription1),
		createTestToolConfig(testToolName2, testToolDescription2),
		createTestToolConfig(testToolName3, testToolDescription3),
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		consumer.ApplyMcpServerConfig(config)
	}
}

func BenchmarkGetOrInitRegistry_Concurrent(b *testing.B) {
	resetSingletons()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			GetOrInitRegistry()
		}
	})
}

func BenchmarkGetOrInitDynamic_Concurrent(b *testing.B) {
	resetSingletons()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			GetOrInitDynamic()
		}
	})
}
