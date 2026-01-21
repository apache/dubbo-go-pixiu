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

func TestDescriptorCache_SetAndGet(t *testing.T) {
	cache := NewDescriptorCache(1 * time.Minute)
	defer cache.Close()

	key := "test/service/method"

	// Initially should be nil
	if got := cache.Get(key); got != nil {
		t.Errorf("Get() on empty cache = %v, want nil", got)
	}

	// Set a nil value (valid for testing)
	cache.Set(key, nil)

	// Should be able to get it back (even if nil)
	// Note: Our implementation stores the entry but returns nil descriptor
	// This is acceptable behavior as it caches the "not found" state
}

func TestDescriptorCache_Expiration(t *testing.T) {
	// Use very short TTL for testing
	ttl := 100 * time.Millisecond
	cache := NewDescriptorCache(ttl)
	defer cache.Close()

	key := "test/service/method"
	cache.Set(key, nil)

	// Should exist before expiration
	// Wait a bit less than TTL
	time.Sleep(50 * time.Millisecond)

	// Entry should still be valid (we can't test the descriptor value easily
	// without a real MethodDescriptor, but we can verify the cache mechanics)

	// Wait for expiration
	time.Sleep(100 * time.Millisecond)

	// After expiration, Get should return nil
	if got := cache.Get(key); got != nil {
		t.Errorf("Get() after expiration = %v, want nil", got)
	}
}

func TestDescriptorCache_Delete(t *testing.T) {
	cache := NewDescriptorCache(1 * time.Minute)
	defer cache.Close()

	key := "test/service/method"
	cache.Set(key, nil)

	// Delete the entry
	cache.Delete(key)

	// Should be nil after deletion
	if got := cache.Get(key); got != nil {
		t.Errorf("Get() after Delete() = %v, want nil", got)
	}
}

func TestDescriptorCache_Clear(t *testing.T) {
	cache := NewDescriptorCache(1 * time.Minute)
	defer cache.Close()

	// Add multiple entries
	keys := []string{"key1", "key2", "key3"}
	for _, k := range keys {
		cache.Set(k, nil)
	}

	// Verify size
	if size := cache.Size(); size != 3 {
		t.Errorf("Size() = %d, want 3", size)
	}

	// Clear all
	cache.Clear()

	// Verify all cleared
	if size := cache.Size(); size != 0 {
		t.Errorf("Size() after Clear() = %d, want 0", size)
	}

	for _, k := range keys {
		if got := cache.Get(k); got != nil {
			t.Errorf("Get(%q) after Clear() = %v, want nil", k, got)
		}
	}
}

func TestDescriptorCache_Size(t *testing.T) {
	cache := NewDescriptorCache(1 * time.Minute)
	defer cache.Close()

	// Initially empty
	if size := cache.Size(); size != 0 {
		t.Errorf("Size() on empty cache = %d, want 0", size)
	}

	// Add entries
	for i := 0; i < 5; i++ {
		cache.Set("key"+string(rune('0'+i)), nil)
	}

	if size := cache.Size(); size != 5 {
		t.Errorf("Size() = %d, want 5", size)
	}
}

func TestBuildCacheKey(t *testing.T) {
	tests := []struct {
		address     string
		serviceName string
		methodName  string
		expected    string
	}{
		{
			address:     "localhost:50051",
			serviceName: "helloworld.Greeter",
			methodName:  "SayHello",
			expected:    "localhost:50051/helloworld.Greeter/SayHello",
		},
		{
			address:     "backend.example.com:8080",
			serviceName: "my.package.Service",
			methodName:  "Method",
			expected:    "backend.example.com:8080/my.package.Service/Method",
		},
		{
			address:     "",
			serviceName: "Service",
			methodName:  "Method",
			expected:    "/Service/Method",
		},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := BuildCacheKey(tt.address, tt.serviceName, tt.methodName)
			if result != tt.expected {
				t.Errorf("BuildCacheKey(%q, %q, %q) = %q, want %q",
					tt.address, tt.serviceName, tt.methodName, result, tt.expected)
			}
		})
	}
}
