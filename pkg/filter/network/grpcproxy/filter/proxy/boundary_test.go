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
	"testing"
	"time"
)

import (
	"go.uber.org/zap"

	"google.golang.org/protobuf/types/descriptorpb"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

// TestCacheStats tests cache statistics tracking
func TestCacheStats(t *testing.T) {
	cache := NewDescriptorCache(5 * time.Minute)

	// Initially no hits or misses
	stats := cache.GetStats()
	if stats.Hits != 0 || stats.Misses != 0 {
		t.Errorf("Expected zero stats, got hits=%d, misses=%d", stats.Hits, stats.Misses)
	}

	// Generate some misses
	cache.Get("nonexistent")
	cache.Get("also-nonexistent")

	stats = cache.GetStats()
	if stats.Misses != 2 {
		t.Errorf("Expected 2 misses, got %d", stats.Misses)
	}

	// Check max cache size
	if stats.MaxSize != defaultMaxCacheSize {
		t.Errorf("Expected max_cache_size %d, got %d",
			defaultMaxCacheSize, stats.MaxSize)
	}

	// Check TTL
	if stats.TTL.Seconds() != 300 {
		t.Errorf("Expected TTL 300s, got %v", stats.TTL.Seconds())
	}

	// Reset stats
	cache.ResetStats()
	stats = cache.GetStats()
	if stats.Hits != 0 || stats.Misses != 0 {
		t.Errorf("Expected stats to be reset, got hits=%d, misses=%d", stats.Hits, stats.Misses)
	}
}

// TestDynamicCacheSizeAdjustment tests dynamic cache size adjustment
func TestDynamicCacheSizeAdjustment(t *testing.T) {
	cache := NewDescriptorCache(5 * time.Minute)

	// Setting below minimum should adjust to minimum
	cache.SetMaxSize(10)
	stats := cache.GetStats()
	if stats.MaxSize != minCacheSize {
		t.Errorf("Expected MaxSize %d (minCacheSize), got %d", minCacheSize, stats.MaxSize)
	}

	// Setting above minimum should work
	cache.SetMaxSize(500)
	stats = cache.GetStats()
	if stats.MaxSize != 500 {
		t.Errorf("Expected MaxSize 500, got %d", stats.MaxSize)
	}
}

// TestInvalidateByVersion tests version-based cache invalidation API
func TestInvalidateByVersion(t *testing.T) {
	cache := NewDescriptorCache(5 * time.Minute)

	// Test with non-existent version hash
	count := cache.InvalidateByVersion("non-existent-version")
	if count != 0 {
		t.Errorf("Expected 0 invalidations for non-existent version, got %d", count)
	}
}

// TestReflectionManagerConfig tests configuration updates
func TestReflectionManagerConfig(t *testing.T) {
	config := ReflectionConfig{
		CacheTTL:          10 * time.Minute,
		MaxCacheSize:      500,
		ReflectionVersion: ReflectionV1Alpha,
		ContinueOnError:   true,
	}

	rm := NewReflectionManagerWithConfig(config)

	// Verify initial config
	initialConfig := rm.GetConfig()
	if initialConfig.MaxCacheSize != 500 {
		t.Errorf("Expected MaxCacheSize 500, got %d", initialConfig.MaxCacheSize)
	}

	// Update config
	newConfig := ReflectionConfig{
		MaxCacheSize:    1000,
		ContinueOnError: false,
	}
	rm.SetConfig(newConfig)

	updatedConfig := rm.GetConfig()
	if updatedConfig.MaxCacheSize != 1000 {
		t.Errorf("Expected MaxCacheSize 1000, got %d", updatedConfig.MaxCacheSize)
	}
	if updatedConfig.ContinueOnError != false {
		t.Error("Expected ContinueOnError to be false")
	}
}

// TestReflectionManagerCacheStats tests enhanced cache statistics
func TestReflectionManagerCacheStats(t *testing.T) {
	rm := NewReflectionManager(5 * time.Minute)

	stats := rm.GetCacheStats()

	// Verify stats struct has expected values
	if stats.MaxCacheSize != defaultMaxCacheSize {
		t.Errorf("Expected MaxCacheSize %d, got %d",
			defaultMaxCacheSize, stats.MaxCacheSize)
	}

	if stats.TTLSeconds != 300 {
		t.Errorf("Expected TTLSeconds 300, got %v", stats.TTLSeconds)
	}

	// Initial stats should be zero
	if stats.MethodCacheHits != 0 {
		t.Errorf("Expected MethodCacheHits 0, got %d", stats.MethodCacheHits)
	}

	if stats.FileRegistryCount != 0 {
		t.Errorf("Expected FileRegistryCount 0, got %d", stats.FileRegistryCount)
	}
}

// TestGetMissingDependencies tests missing dependency tracking
func TestGetMissingDependencies(t *testing.T) {
	rm := NewReflectionManagerWithConfig(ReflectionConfig{
		ContinueOnError: true,
	})

	// Initially no missing dependencies
	missing := rm.GetMissingDependencies("test-address")
	if missing != nil {
		t.Errorf("Expected nil missing dependencies, got %v", missing)
	}

	// Get all missing dependencies should return empty map
	allMissing := rm.GetAllMissingDependencies()
	if len(allMissing) != 0 {
		t.Errorf("Expected empty map, got %d entries", len(allMissing))
	}
}

// TestReflectionManagerContinueOnError tests the ContinueOnError behavior
func TestReflectionManagerContinueOnError(t *testing.T) {
	tests := []struct {
		name            string
		continueOnError bool
	}{
		{
			name:            "continue on error enabled",
			continueOnError: true,
		},
		{
			name:            "continue on error disabled",
			continueOnError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := ReflectionConfig{
				CacheTTL:          5 * time.Minute,
				ContinueOnError:   tt.continueOnError,
				ReflectionVersion: ReflectionV1Alpha,
			}

			rm := NewReflectionManagerWithConfig(config)

			if rm.GetConfig().ContinueOnError != tt.continueOnError {
				t.Errorf("Expected ContinueOnError %v, got %v",
					tt.continueOnError, rm.GetConfig().ContinueOnError)
			}
		})
	}
}

// TestReflectionManagerInvalidationCache tests cache invalidation
func TestReflectionManagerInvalidationCache(t *testing.T) {
	rm := NewReflectionManager(5 * time.Minute)

	// Invalidate non-existent address should not panic
	rm.InvalidateCache("non-existent-address")

	// Clear cache should not panic
	rm.ClearCache()

	// Close should not panic
	rm.Close()
}

// TestTopologicalSortCircularDependency tests handling of circular dependencies
func TestTopologicalSortCircularDependency(t *testing.T) {
	rm := NewReflectionManager(5 * time.Minute)

	// Create circular dependency: A depends on B, B depends on A
	fileDescs := []*descriptorpb.FileDescriptorProto{
		{
			Name:       strPtr("a.proto"),
			Dependency: []string{"b.proto"},
			Package:    strPtr("test"),
		},
		{
			Name:       strPtr("b.proto"),
			Dependency: []string{"a.proto"},
			Package:    strPtr("test"),
		},
	}

	// Should handle circular dependency gracefully (return all files)
	result, err := rm.topologicalSort(fileDescs)
	if err != nil {
		t.Errorf("Topological sort should handle circular dependencies, got error: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected 2 files in result, got %d", len(result))
	}
}

// TestTopologicalSortEmptyList tests empty file descriptor list
func TestTopologicalSortEmptyList(t *testing.T) {
	rm := NewReflectionManager(5 * time.Minute)

	result, err := rm.topologicalSort([]*descriptorpb.FileDescriptorProto{})
	if err != nil {
		t.Errorf("Expected no error for empty list, got: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected empty result, got %d files", len(result))
	}
}

// TestTopologicalSortWithExternalDependencies tests handling of external dependencies
func TestTopologicalSortWithExternalDependencies(t *testing.T) {
	rm := NewReflectionManager(5 * time.Minute)

	fileDescs := []*descriptorpb.FileDescriptorProto{
		{
			Name:       strPtr("a.proto"),
			Dependency: []string{"google/protobuf/empty.proto"},
			Package:    strPtr("test"),
		},
	}

	// Should handle external dependency gracefully
	result, err := rm.topologicalSort(fileDescs)
	if err != nil {
		t.Errorf("Expected no error with external dependency, got: %v", err)
	}

	if len(result) != 1 {
		t.Errorf("Expected 1 file in result, got %d", len(result))
	}
}

// TestReflectionVersionConstants tests reflection version constants
func TestReflectionVersionConstants(t *testing.T) {
	tests := []struct {
		name    string
		version ReflectionVersion
		valid   bool
	}{
		{"v1alpha", ReflectionV1Alpha, true},
		{"v1", ReflectionV1, true},
		{"auto", ReflectionAuto, true},
		{"invalid", ReflectionVersion("invalid"), true}, // Any string is valid
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.version == "" {
				t.Error("Reflection version should not be empty")
			}
		})
	}
}

// TestNewReflectionManagerWithInvalidConfig tests configuration validation
func TestNewReflectionManagerWithInvalidConfig(t *testing.T) {
	tests := []struct {
		name            string
		config          ReflectionConfig
		expectCacheSize int
		expectTTL       time.Duration
	}{
		{
			name: "zero TTL uses default",
			config: ReflectionConfig{
				CacheTTL:        0,
				MaxCacheSize:    100,
				ContinueOnError: true,
			},
			expectCacheSize: 100,
			expectTTL:       defaultDescCacheTTL,
		},
		{
			name: "negative TTL uses default",
			config: ReflectionConfig{
				CacheTTL:        -1,
				MaxCacheSize:    100,
				ContinueOnError: true,
			},
			expectCacheSize: 100,
			expectTTL:       defaultDescCacheTTL,
		},
		{
			name: "cache size below minimum",
			config: ReflectionConfig{
				CacheTTL:        5 * time.Minute,
				MaxCacheSize:    10,
				ContinueOnError: true,
			},
			expectCacheSize: minCacheSize,
			expectTTL:       5 * time.Minute,
		},
		{
			name: "empty reflection version",
			config: ReflectionConfig{
				CacheTTL:          5 * time.Minute,
				MaxCacheSize:      100,
				ReflectionVersion: "",
				ContinueOnError:   true,
			},
			expectCacheSize: 100,
			expectTTL:       5 * time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm := NewReflectionManagerWithConfig(tt.config)

			config := rm.GetConfig()
			if config.MaxCacheSize != tt.expectCacheSize {
				t.Errorf("Expected cache size %d, got %d", tt.expectCacheSize, config.MaxCacheSize)
			}

			if config.CacheTTL != tt.expectTTL {
				t.Errorf("Expected TTL %v, got %v", tt.expectTTL, config.CacheTTL)
			}
		})
	}
}

// Helper function for string pointers
func strPtr(s string) *string {
	return &s
}

// Init logger for tests
func init() {
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	logger.InitLogger(&cfg)
}
