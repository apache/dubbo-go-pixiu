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
	"sync"
	"sync/atomic"
	"time"
)

import (
	"github.com/dubbogo/gost/container/gxlru"

	"google.golang.org/protobuf/reflect/protoreflect"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

// Constants for cache configuration
const (
	defaultMaxCacheSize = 1000 // Maximum number of cached descriptors
	minCacheSize        = 100  // Minimum cache size to prevent thrashing
	// Logger prefix for descriptor cache operations
	loggerPrefix = "[grpc-proxy-descriptor-cache] "
)

// cacheEntry represents a cached method descriptor with metadata
// Implements gxlru.Value interface (Size always returns 1 for counting)
type cacheEntry struct {
	key         string
	descriptor  protoreflect.MethodDescriptor
	expireAt    time.Time
	versionHash string // Hash of the descriptor for version detection
}

// Size implements gxlru.Value interface - returns 1 for per-entry counting
func (e *cacheEntry) Size() int {
	return 1
}

// CacheStats holds cache statistics for monitoring (returned by GetStats)
type CacheStats struct {
	Hits      int64
	Misses    int64
	Evictions int64
	Size      int
	MaxSize   int
	TTL       time.Duration
	HitRatio  float64
}

// DescriptorCache provides TTL-based caching with LRU eviction for gRPC method descriptors
type DescriptorCache struct {
	lru       *gxlru.LRUCache // LRU cache from gost
	entries   sync.Map        // key -> *cacheEntry (for TTL and version lookup)
	ttl       time.Duration
	stopCh    chan struct{}
	closeOnce sync.Once
	maxSize   atomic.Int32 // Maximum cache size (atomic for lock-free reads)
	// Atomic statistics for lock-free updates
	hits      atomic.Int64
	misses    atomic.Int64
	evictions atomic.Int64
}

// NewDescriptorCache creates a new descriptor cache with the specified TTL
func NewDescriptorCache(ttl time.Duration) *DescriptorCache {
	return NewDescriptorCacheWithSize(ttl, defaultMaxCacheSize)
}

// NewDescriptorCacheWithSize creates a new descriptor cache with custom size limit
func NewDescriptorCacheWithSize(ttl time.Duration, maxSize int) *DescriptorCache {
	if maxSize < minCacheSize {
		maxSize = minCacheSize
		logger.Infof("%sCache size adjusted to minimum: %d", loggerPrefix, maxSize)
	}

	cache := &DescriptorCache{
		lru:    gxlru.NewLRUCache(int64(maxSize)),
		ttl:    ttl,
		stopCh: make(chan struct{}),
	}
	cache.maxSize.Store(int32(maxSize))

	go cache.cleanupLoop()
	return cache
}

// Get retrieves a method descriptor from cache
// Returns nil if not found or expired
func (c *DescriptorCache) Get(key string) protoreflect.MethodDescriptor {
	entry, ok := c.entries.Load(key)
	if !ok {
		c.misses.Add(1)
		return nil
	}

	e := entry.(*cacheEntry)
	if time.Now().Before(e.expireAt) {
		// Valid cache hit
		c.hits.Add(1)
		return e.descriptor
	}

	// Entry expired, delete it
	c.entries.Delete(key)
	c.lru.Delete(key)
	c.misses.Add(1)
	return nil
}

// GetWithVersion retrieves a method descriptor and returns version mismatch status
// Returns (descriptor, true) if found and version matches
// Returns (descriptor, false) if found but version differs (stale)
// Returns (nil, false) if not found or expired
func (c *DescriptorCache) GetWithVersion(key string, versionHash string) (protoreflect.MethodDescriptor, bool) {
	entry, ok := c.entries.Load(key)
	if !ok {
		c.misses.Add(1)
		return nil, false
	}

	e := entry.(*cacheEntry)
	if time.Now().After(e.expireAt) {
		// Entry expired
		c.entries.Delete(key)
		c.lru.Delete(key)
		c.misses.Add(1)
		return nil, false
	}

	// Check version hash
	if e.versionHash != "" && versionHash != "" && e.versionHash != versionHash {
		// Version mismatch - stale cache
		logger.Debugf("%sCache version mismatch for %s: cached=%s, current=%s",
			loggerPrefix, key, e.versionHash, versionHash)
		c.entries.Delete(key)
		c.lru.Delete(key)
		c.misses.Add(1)
		return nil, false
	}

	// Valid cache hit
	c.hits.Add(1)
	return e.descriptor, true
}

// Set stores a method descriptor in the cache with optional version hash
func (c *DescriptorCache) Set(key string, descriptor protoreflect.MethodDescriptor) {
	c.SetWithVersion(key, descriptor, "")
}

// SetWithVersion stores a method descriptor with version hash for change detection
func (c *DescriptorCache) SetWithVersion(key string, descriptor protoreflect.MethodDescriptor, versionHash string) {
	newEntry := &cacheEntry{
		key:         key,
		descriptor:  descriptor,
		expireAt:    time.Now().Add(c.ttl),
		versionHash: versionHash,
	}

	// Add to LRU cache (handles eviction automatically)
	c.lru.Set(key, newEntry)

	// Store in entries map for TTL/version lookup (possibly replacing old entry)
	c.entries.Store(key, newEntry)
}

// Delete removes a specific entry from the cache
func (c *DescriptorCache) Delete(key string) {
	c.entries.Delete(key)
	c.lru.Delete(key)
}

// InvalidateByVersion removes all entries with a specific version hash
// Useful when service definitions are updated
func (c *DescriptorCache) InvalidateByVersion(versionHash string) int {
	count := 0
	c.entries.Range(func(key, value any) bool {
		entry := value.(*cacheEntry)
		if entry.versionHash == versionHash {
			c.entries.Delete(key)
			c.lru.Delete(key.(string))
			count++
		}
		return true
	})
	if count > 0 {
		logger.Infof("%sInvalidated %d cache entries with version hash: %s", loggerPrefix, count, versionHash)
	}
	return count
}

// Clear removes all entries from the cache
func (c *DescriptorCache) Clear() {
	c.entries.Range(func(key, value any) bool {
		c.entries.Delete(key)
		c.lru.Delete(key.(string))
		return true
	})
	c.lru.Clear()
}

// Size returns the current number of entries in the cache
func (c *DescriptorCache) Size() int {
	return int(c.lru.Length())
}

// GetStats returns current cache statistics
func (c *DescriptorCache) GetStats() CacheStats {
	hits := c.hits.Load()
	misses := c.misses.Load()
	total := hits + misses

	var hitRatio float64
	if total > 0 {
		hitRatio = float64(hits) / float64(total)
	}

	// Use gost's eviction count + our own tracking
	evictions := c.evictions.Load() + c.lru.Evictions()

	return CacheStats{
		Hits:      hits,
		Misses:    misses,
		Evictions: evictions,
		Size:      int(c.lru.Length()),
		MaxSize:   int(c.maxSize.Load()),
		TTL:       c.ttl,
		HitRatio:  hitRatio,
	}
}

// ResetStats clears cache statistics (atomic reset)
func (c *DescriptorCache) ResetStats() {
	c.hits.Store(0)
	c.misses.Store(0)
	c.evictions.Store(0)
	// Note: gost's internal eviction counter cannot be reset
}

// SetMaxSize dynamically adjusts the maximum cache size
func (c *DescriptorCache) SetMaxSize(size int) {
	if size < minCacheSize {
		size = minCacheSize
	}
	c.maxSize.Store(int32(size))
	c.lru.SetCapacity(int64(size))
}

// cleanupLoop periodically removes expired entries
func (c *DescriptorCache) cleanupLoop() {
	ticker := time.NewTicker(c.ttl / 2)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			now := time.Now()
			c.entries.Range(func(key, value any) bool {
				entry := value.(*cacheEntry)
				if now.After(entry.expireAt) {
					c.entries.Delete(key)
					c.lru.Delete(key.(string))
				}
				return true
			})
		case <-c.stopCh:
			return
		}
	}
}

// Close stops the cleanup goroutine (safe to call multiple times)
func (c *DescriptorCache) Close() {
	c.closeOnce.Do(func() {
		close(c.stopCh)
	})
}

// BuildCacheKey creates a cache key from service and method names
func BuildCacheKey(address, serviceName, methodName string) string {
	return address + "/" + serviceName + "/" + methodName
}
