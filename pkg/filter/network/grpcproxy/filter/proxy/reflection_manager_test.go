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
	"google.golang.org/protobuf/proto"

	"google.golang.org/protobuf/types/descriptorpb"
)

func TestNewReflectionManager(t *testing.T) {
	tests := []struct {
		name     string
		ttl      time.Duration
		expected time.Duration
	}{
		{
			name:     "positive TTL",
			ttl:      10 * time.Minute,
			expected: 10 * time.Minute,
		},
		{
			name:     "zero TTL uses default",
			ttl:      0,
			expected: defaultDescCacheTTL,
		},
		{
			name:     "negative TTL uses default",
			ttl:      -1 * time.Minute,
			expected: defaultDescCacheTTL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rm := NewReflectionManager(tt.ttl)
			defer rm.Close()

			if rm == nil {
				t.Fatal("NewReflectionManager returned nil")
			}
			if rm.cache == nil {
				t.Error("cache is nil")
			}
			if rm.cacheTTL != tt.expected {
				t.Errorf("cacheTTL = %v, want %v", rm.cacheTTL, tt.expected)
			}
		})
	}
}

func TestReflectionManager_InvalidateCache(t *testing.T) {
	rm := NewReflectionManager(5 * time.Minute)
	defer rm.Close()

	// Store something in the file descriptor cache
	rm.fileDescCache.Store("test-address:50051", "dummy-value")

	// Verify it's stored
	if _, ok := rm.fileDescCache.Load("test-address:50051"); !ok {
		t.Fatal("value not stored in fileDescCache")
	}

	// Invalidate
	rm.InvalidateCache("test-address:50051")

	// Verify it's gone
	if _, ok := rm.fileDescCache.Load("test-address:50051"); ok {
		t.Error("value still exists after InvalidateCache")
	}
}

func TestReflectionManager_ClearCache(t *testing.T) {
	rm := NewReflectionManager(5 * time.Minute)
	defer rm.Close()

	// Store multiple entries
	rm.fileDescCache.Store("addr1:50051", "value1")
	rm.fileDescCache.Store("addr2:50052", "value2")
	rm.fileDescCache.Store("addr3:50053", "value3")

	// Also set some method cache entries
	rm.cache.Set("key1", nil)
	rm.cache.Set("key2", nil)

	// Clear all
	rm.ClearCache()

	// Verify file descriptor cache is cleared
	count := 0
	rm.fileDescCache.Range(func(_, _ any) bool {
		count++
		return true
	})
	if count != 0 {
		t.Errorf("fileDescCache has %d entries after ClearCache, want 0", count)
	}

	// Verify method cache is cleared
	if size := rm.cache.Size(); size != 0 {
		t.Errorf("method cache has %d entries after ClearCache, want 0", size)
	}
}

func TestReflectionManager_Close(t *testing.T) {
	rm := NewReflectionManager(5 * time.Minute)

	// Should not panic
	rm.Close()

	// Close again should also not panic
	rm.Close()
}

func TestReflectionManager_GetCacheStats(t *testing.T) {
	rm := NewReflectionManager(10 * time.Minute)
	defer rm.Close()

	// Initially empty
	stats := rm.GetCacheStats()

	if stats.MethodCacheSize != 0 {
		t.Errorf("MethodCacheSize = %v, want 0", stats.MethodCacheSize)
	}
	if stats.FileRegistryCount != 0 {
		t.Errorf("FileRegistryCount = %v, want 0", stats.FileRegistryCount)
	}
	if stats.TTLSeconds != 600.0 {
		t.Errorf("TTLSeconds = %v, want 600", stats.TTLSeconds)
	}

	// Add some entries
	rm.fileDescCache.Store("addr1", "value1")
	rm.fileDescCache.Store("addr2", "value2")
	rm.cache.Set("method1", nil)

	stats = rm.GetCacheStats()
	if stats.MethodCacheSize != 1 {
		t.Errorf("MethodCacheSize = %v, want 1", stats.MethodCacheSize)
	}
	if stats.FileRegistryCount != 2 {
		t.Errorf("FileRegistryCount = %v, want 2", stats.FileRegistryCount)
	}
}

func TestReflectionManager_TopologicalSort(t *testing.T) {
	rm := NewReflectionManager(5 * time.Minute)
	defer rm.Close()

	tests := []struct {
		name      string
		fileDescs []*descriptorpb.FileDescriptorProto
		wantOrder []string // Expected order (dependencies before dependents)
	}{
		{
			name: "single file no deps",
			fileDescs: []*descriptorpb.FileDescriptorProto{
				{Name: proto.String("a.proto")},
			},
			wantOrder: []string{"a.proto"},
		},
		{
			name: "linear dependency chain",
			fileDescs: []*descriptorpb.FileDescriptorProto{
				{Name: proto.String("c.proto"), Dependency: []string{"b.proto"}},
				{Name: proto.String("b.proto"), Dependency: []string{"a.proto"}},
				{Name: proto.String("a.proto")},
			},
			wantOrder: []string{"a.proto", "b.proto", "c.proto"},
		},
		{
			name: "diamond dependency",
			fileDescs: []*descriptorpb.FileDescriptorProto{
				{Name: proto.String("d.proto"), Dependency: []string{"b.proto", "c.proto"}},
				{Name: proto.String("b.proto"), Dependency: []string{"a.proto"}},
				{Name: proto.String("c.proto"), Dependency: []string{"a.proto"}},
				{Name: proto.String("a.proto")},
			},
			// a must come first, b and c can be in any order, d must come last
			wantOrder: nil, // We'll check constraints instead
		},
		{
			name: "external dependency (not in list)",
			fileDescs: []*descriptorpb.FileDescriptorProto{
				{Name: proto.String("b.proto"), Dependency: []string{"external.proto"}},
				{Name: proto.String("a.proto")},
			},
			wantOrder: nil, // Just verify no error
		},
		{
			name:      "empty list",
			fileDescs: []*descriptorpb.FileDescriptorProto{},
			wantOrder: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := rm.topologicalSort(tt.fileDescs)
			if err != nil {
				t.Fatalf("topologicalSort() error = %v", err)
			}

			if tt.wantOrder != nil {
				// Exact order check
				if len(result) != len(tt.wantOrder) {
					t.Errorf("got %d files, want %d", len(result), len(tt.wantOrder))
					return
				}
				for i, fd := range result {
					if fd.GetName() != tt.wantOrder[i] {
						t.Errorf("position %d: got %s, want %s", i, fd.GetName(), tt.wantOrder[i])
					}
				}
			} else {
				// Just verify length matches and no error
				if len(result) != len(tt.fileDescs) {
					t.Errorf("got %d files, want %d", len(result), len(tt.fileDescs))
				}
			}
		})
	}
}

func TestReflectionManager_TopologicalSort_DependencyOrder(t *testing.T) {
	rm := NewReflectionManager(5 * time.Minute)
	defer rm.Close()

	// Create a complex dependency graph
	fileDescs := []*descriptorpb.FileDescriptorProto{
		{Name: proto.String("leaf1.proto"), Dependency: []string{"mid.proto"}},
		{Name: proto.String("leaf2.proto"), Dependency: []string{"mid.proto"}},
		{Name: proto.String("mid.proto"), Dependency: []string{"root.proto"}},
		{Name: proto.String("root.proto")},
	}

	result, err := rm.topologicalSort(fileDescs)
	if err != nil {
		t.Fatalf("topologicalSort() error = %v", err)
	}

	// Build position map
	position := make(map[string]int)
	for i, fd := range result {
		position[fd.GetName()] = i
	}

	// Verify dependency constraints
	// root should come before mid
	if position["root.proto"] >= position["mid.proto"] {
		t.Error("root.proto should come before mid.proto")
	}
	// mid should come before leaf1 and leaf2
	if position["mid.proto"] >= position["leaf1.proto"] {
		t.Error("mid.proto should come before leaf1.proto")
	}
	if position["mid.proto"] >= position["leaf2.proto"] {
		t.Error("mid.proto should come before leaf2.proto")
	}
}

func TestReflectionManager_BuildFileRegistry(t *testing.T) {
	rm := NewReflectionManager(5 * time.Minute)
	defer rm.Close()

	// Create a simple proto file descriptor
	fileDescs := []*descriptorpb.FileDescriptorProto{
		{
			Name:    proto.String("test.proto"),
			Package: proto.String("test"),
			MessageType: []*descriptorpb.DescriptorProto{
				{
					Name: proto.String("TestMessage"),
					Field: []*descriptorpb.FieldDescriptorProto{
						{
							Name:   proto.String("value"),
							Number: proto.Int32(1),
							Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
							Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						},
					},
				},
			},
		},
	}

	registry, err := rm.buildFileRegistry(fileDescs)
	if err != nil {
		t.Fatalf("buildFileRegistry() error = %v", err)
	}

	if registry == nil {
		t.Fatal("buildFileRegistry() returned nil registry")
	}

	// Verify the file was registered
	fd, err := registry.FindFileByPath("test.proto")
	if err != nil {
		t.Fatalf("FindFileByPath() error = %v", err)
	}

	if fd.Package().Name() != "test" {
		t.Errorf("package = %q, want %q", fd.Package().Name(), "test")
	}
}

func TestReflectionManager_BuildFileRegistry_WithDependencies(t *testing.T) {
	rm := NewReflectionManager(5 * time.Minute)
	defer rm.Close()

	// Create proto files with dependency
	fileDescs := []*descriptorpb.FileDescriptorProto{
		{
			Name:    proto.String("base.proto"),
			Package: proto.String("base"),
			MessageType: []*descriptorpb.DescriptorProto{
				{
					Name: proto.String("BaseMessage"),
					Field: []*descriptorpb.FieldDescriptorProto{
						{
							Name:   proto.String("id"),
							Number: proto.Int32(1),
							Type:   descriptorpb.FieldDescriptorProto_TYPE_INT64.Enum(),
							Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						},
					},
				},
			},
		},
		{
			Name:       proto.String("derived.proto"),
			Package:    proto.String("derived"),
			Dependency: []string{"base.proto"},
			MessageType: []*descriptorpb.DescriptorProto{
				{
					Name: proto.String("DerivedMessage"),
					Field: []*descriptorpb.FieldDescriptorProto{
						{
							Name:     proto.String("base"),
							Number:   proto.Int32(1),
							Type:     descriptorpb.FieldDescriptorProto_TYPE_MESSAGE.Enum(),
							TypeName: proto.String(".base.BaseMessage"),
							Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						},
					},
				},
			},
		},
	}

	registry, err := rm.buildFileRegistry(fileDescs)
	if err != nil {
		t.Fatalf("buildFileRegistry() error = %v", err)
	}

	// Verify both files were registered
	_, err = registry.FindFileByPath("base.proto")
	if err != nil {
		t.Errorf("FindFileByPath(base.proto) error = %v", err)
	}

	_, err = registry.FindFileByPath("derived.proto")
	if err != nil {
		t.Errorf("FindFileByPath(derived.proto) error = %v", err)
	}
}

func TestReflectionManager_BuildFileRegistry_Empty(t *testing.T) {
	rm := NewReflectionManager(5 * time.Minute)
	defer rm.Close()

	registry, err := rm.buildFileRegistry([]*descriptorpb.FileDescriptorProto{})
	if err != nil {
		t.Fatalf("buildFileRegistry() error = %v", err)
	}

	if registry == nil {
		t.Error("buildFileRegistry() returned nil for empty input")
	}
}

func TestReflectionManager_BuildFileRegistry_WithService(t *testing.T) {
	rm := NewReflectionManager(5 * time.Minute)
	defer rm.Close()

	fileDescs := []*descriptorpb.FileDescriptorProto{
		{
			Name:    proto.String("service.proto"),
			Package: proto.String("myservice"),
			MessageType: []*descriptorpb.DescriptorProto{
				{
					Name: proto.String("Request"),
					Field: []*descriptorpb.FieldDescriptorProto{
						{
							Name:   proto.String("query"),
							Number: proto.Int32(1),
							Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
							Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						},
					},
				},
				{
					Name: proto.String("Response"),
					Field: []*descriptorpb.FieldDescriptorProto{
						{
							Name:   proto.String("result"),
							Number: proto.Int32(1),
							Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
							Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
						},
					},
				},
			},
			Service: []*descriptorpb.ServiceDescriptorProto{
				{
					Name: proto.String("MyService"),
					Method: []*descriptorpb.MethodDescriptorProto{
						{
							Name:       proto.String("Search"),
							InputType:  proto.String(".myservice.Request"),
							OutputType: proto.String(".myservice.Response"),
						},
					},
				},
			},
		},
	}

	registry, err := rm.buildFileRegistry(fileDescs)
	if err != nil {
		t.Fatalf("buildFileRegistry() error = %v", err)
	}

	// Find service descriptor
	desc, err := registry.FindDescriptorByName("myservice.MyService")
	if err != nil {
		t.Fatalf("FindDescriptorByName() error = %v", err)
	}

	if desc.Name() != "MyService" {
		t.Errorf("service name = %q, want %q", desc.Name(), "MyService")
	}
}
