// Copyright Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

import (
	"fmt"
	"sync"
	"time"
)

import (
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/golang-lru/simplelru"
	"google.golang.org/protobuf/testing/protocmp"
	"istio.io/pkg/monitoring"
)

import (
	"github.com/apache/dubbo-go-pixiu/pilot/pkg/features"
	"github.com/apache/dubbo-go-pixiu/pkg/config"
	"github.com/apache/dubbo-go-pixiu/pkg/util/sets"
)

func init() {
	monitoring.MustRegister(xdsCacheReads)
	monitoring.MustRegister(xdsCacheEvictions)
	monitoring.MustRegister(xdsCacheSize)
}

var (
	xdsCacheReads = monitoring.NewSum(
		"xds_cache_reads",
		"Total number of xds cache xdsCacheReads.",
		monitoring.WithLabels(typeTag),
	)

	xdsCacheEvictions = monitoring.NewSum(
		"xds_cache_evictions",
		"Total number of xds cache evictions.",
	)

	xdsCacheSize = monitoring.NewGauge(
		"xds_cache_size",
		"Current size of xds cache",
	)

	xdsCacheHits   = xdsCacheReads.With(typeTag.Value("hit"))
	xdsCacheMisses = xdsCacheReads.With(typeTag.Value("miss"))
)

func hit() {
	if features.EnableXDSCacheMetrics {
		xdsCacheHits.Increment()
	}
}

func miss() {
	if features.EnableXDSCacheMetrics {
		xdsCacheMisses.Increment()
	}
}

func size(cs int) {
	if features.EnableXDSCacheMetrics {
		xdsCacheSize.Record(float64(cs))
	}
}

func indexConfig(configIndex map[ConfigKey]sets.Set, k string, dependentConfigs []ConfigKey) {
	for _, cfg := range dependentConfigs {
		if configIndex[cfg] == nil {
			configIndex[cfg] = sets.New()
		}
		configIndex[cfg].Insert(k)
	}
}

func clearIndexConfig(configIndex map[ConfigKey]sets.Set, k string, dependentConfigs []ConfigKey) {
	for _, cfg := range dependentConfigs {
		index := configIndex[cfg]
		if index != nil {
			index.Delete(k)
			if index.IsEmpty() {
				delete(configIndex, cfg)
			}
		}
	}
}

func indexType(typeIndex map[config.GroupVersionKind]sets.Set, k string, dependentTypes []config.GroupVersionKind) {
	for _, t := range dependentTypes {
		if typeIndex[t] == nil {
			typeIndex[t] = sets.New()
		}
		typeIndex[t].Insert(k)
	}
}

func clearIndexType(typeIndex map[config.GroupVersionKind]sets.Set, k string, dependentTypes []config.GroupVersionKind) {
	for _, t := range dependentTypes {
		index := typeIndex[t]
		if index != nil {
			index.Delete(k)
			if index.IsEmpty() {
				delete(typeIndex, t)
			}
		}
	}
}

// XdsCacheEntry interface defines functions that should be implemented by
// resources that can be cached.
type XdsCacheEntry interface {
	// Key is the key to be used in cache.
	Key() string
	// DependentTypes are config types that this cache key is dependant on.
	// Whenever any configs of this type changes, we should invalidate this cache entry.
	// Note: DependentConfigs should be preferred wherever possible.
	DependentTypes() []config.GroupVersionKind
	// DependentConfigs is config items that this cache key is dependent on.
	// Whenever these configs change, we should invalidate this cache entry.
	DependentConfigs() []ConfigKey
	// Cacheable indicates whether this entry is valid for cache. For example
	// for EDS to be cacheable, the Endpoint should have corresponding service.
	Cacheable() bool
}

type CacheToken uint64

// XdsCache interface defines a store for caching XDS responses.
// All operations are thread safe.
type XdsCache interface {
	// Add adds the given XdsCacheEntry with the value for the given pushContext to the cache.
	// If the cache has been updated to a newer push context, the write will be dropped silently.
	// This ensures stale data does not overwrite fresh data when dealing with concurrent
	// writers.
	Add(entry XdsCacheEntry, pushRequest *PushRequest, value *discovery.Resource)
	// Get retrieves the cached value if it exists. The boolean indicates
	// whether the entry exists in the cache.
	Get(entry XdsCacheEntry) (*discovery.Resource, bool)
	// Clear removes the cache entries that are dependent on the configs passed.
	Clear(map[ConfigKey]struct{})
	// ClearAll clears the entire cache.
	ClearAll()
	// Keys returns all currently configured keys. This is for testing/debug only
	Keys() []string
	// Snapshot returns a snapshot of all keys and values. This is for testing/debug only
	Snapshot() map[string]*discovery.Resource
}

// NewXdsCache returns an instance of a cache.
func NewXdsCache() XdsCache {
	cache := &lruCache{
		enableAssertions: features.EnableUnsafeAssertions,
		configIndex:      map[ConfigKey]sets.Set{},
		typesIndex:       map[config.GroupVersionKind]sets.Set{},
	}
	cache.store = newLru(cache.evict)

	return cache
}

// NewLenientXdsCache returns an instance of a cache that does not validate token based get/set and enable assertions.
func NewLenientXdsCache() XdsCache {
	cache := &lruCache{
		enableAssertions: false,
		configIndex:      map[ConfigKey]sets.Set{},
		typesIndex:       map[config.GroupVersionKind]sets.Set{},
	}
	cache.store = newLru(cache.evict)

	return cache
}

type lruCache struct {
	enableAssertions bool
	store            simplelru.LRUCache
	// token stores the latest token of the store, used to prevent stale data overwrite.
	// It is refreshed when Clear or ClearAll are called
	token       CacheToken
	mu          sync.RWMutex
	configIndex map[ConfigKey]sets.Set
	typesIndex  map[config.GroupVersionKind]sets.Set
}

var _ XdsCache = &lruCache{}

func newLru(evictCallback simplelru.EvictCallback) simplelru.LRUCache {
	sz := features.XDSCacheMaxSize
	if sz <= 0 {
		sz = 20000
	}
	l, err := simplelru.NewLRU(sz, evictCallback)
	if err != nil {
		panic(fmt.Errorf("invalid lru configuration: %v", err))
	}
	return l
}

func (l *lruCache) evict(k interface{}, v interface{}) {
	if features.EnableXDSCacheMetrics {
		xdsCacheEvictions.Increment()
	}

	key := k.(string)
	value := v.(cacheValue)

	// we don't need to acquire locks, since this function is called when we write to the store
	clearIndexConfig(l.configIndex, key, value.dependentConfigs)
	clearIndexType(l.typesIndex, key, value.dependentTypes)
}

// assertUnchanged checks that a cache entry is not changed. This helps catch bad cache invalidation
// We should never have a case where we overwrite an existing item with a new change. Instead, when
// config sources change, Clear/ClearAll should be called. At this point, we may get multiple writes
// because multiple writers may get cache misses concurrently, but they ought to generate identical
// configuration. This also checks that our XDS config generation is deterministic, which is a very
// important property.
func (l *lruCache) assertUnchanged(key string, existing *discovery.Resource, replacement *discovery.Resource) {
	if l.enableAssertions {
		if existing == nil {
			// This is a new addition, not an update
			return
		}
		// Record time so that we can correlate when the error actually happened, since the async reporting
		// may be delayed
		t0 := time.Now()
		// This operation is really slow, which makes tests fail for unrelated reasons, so we process it async.
		go func() {
			if !cmp.Equal(existing, replacement, protocmp.Transform()) {
				warning := fmt.Errorf("assertion failed at %v, cache entry changed but not cleared for key %v: %v\n%v\n%v",
					t0, key, cmp.Diff(existing, replacement, protocmp.Transform()), existing, replacement)
				panic(warning)
			}
		}()
	}
}

func (l *lruCache) Add(entry XdsCacheEntry, pushReq *PushRequest, value *discovery.Resource) {
	if !entry.Cacheable() || pushReq == nil || pushReq.Start.Equal(time.Time{}) {
		return
	}
	// It will not overflow until year 2262
	token := CacheToken(pushReq.Start.UnixNano())
	k := entry.Key()
	l.mu.Lock()
	defer l.mu.Unlock()
	cur, f := l.store.Get(k)
	if f {
		// This is the stale resource
		if token < cur.(cacheValue).token || token < l.token {
			// entry may be stale, we need to drop it. This can happen when the cache is invalidated
			// after we call Get.
			return
		}
		if l.enableAssertions {
			l.assertUnchanged(k, cur.(cacheValue).value, value)
		}
	}

	if token < l.token {
		return
	}

	// we have to make sure we evict old entries with the same key
	// to prevent leaking in the index maps
	if old, ok := l.store.Get(k); ok {
		l.evict(k, old)
	}

	dependentConfigs := entry.DependentConfigs()
	dependentTypes := entry.DependentTypes()
	toWrite := cacheValue{value: value, token: token, dependentConfigs: dependentConfigs, dependentTypes: dependentTypes}
	l.store.Add(k, toWrite)
	l.token = token
	indexConfig(l.configIndex, k, dependentConfigs)
	indexType(l.typesIndex, k, dependentTypes)
	size(l.store.Len())
}

type cacheValue struct {
	value            *discovery.Resource
	token            CacheToken
	dependentConfigs []ConfigKey
	dependentTypes   []config.GroupVersionKind
}

func (l *lruCache) Get(entry XdsCacheEntry) (*discovery.Resource, bool) {
	if !entry.Cacheable() {
		return nil, false
	}
	k := entry.Key()
	l.mu.Lock()
	defer l.mu.Unlock()
	val, ok := l.store.Get(k)
	if !ok {
		miss()
		return nil, false
	}
	cv := val.(cacheValue)
	if cv.value == nil {
		miss()
		return nil, false
	}
	hit()
	return cv.value, true
}

func (l *lruCache) Clear(configs map[ConfigKey]struct{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.token = CacheToken(time.Now().UnixNano())
	for ckey := range configs {
		referenced := l.configIndex[ckey]
		delete(l.configIndex, ckey)
		for key := range referenced {
			l.store.Remove(key)
		}
		tReferenced := l.typesIndex[ckey.Kind]
		delete(l.typesIndex, ckey.Kind)
		for key := range tReferenced {
			l.store.Remove(key)
		}
	}
	size(l.store.Len())
}

func (l *lruCache) ClearAll() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.token = CacheToken(time.Now().UnixNano())
	// Purge with an evict function would turn up to be pretty slow since
	// it runs the function for every key in the store, might be better to just
	// create a new store.
	l.store = newLru(l.evict)
	l.configIndex = map[ConfigKey]sets.Set{}
	l.typesIndex = map[config.GroupVersionKind]sets.Set{}
	size(l.store.Len())
}

func (l *lruCache) Keys() []string {
	l.mu.RLock()
	defer l.mu.RUnlock()
	iKeys := l.store.Keys()
	keys := make([]string, 0, len(iKeys))
	for _, ik := range iKeys {
		keys = append(keys, ik.(string))
	}
	return keys
}

func (l *lruCache) Snapshot() map[string]*discovery.Resource {
	l.mu.RLock()
	defer l.mu.RUnlock()
	iKeys := l.store.Keys()
	res := make(map[string]*discovery.Resource, len(iKeys))
	for _, ik := range iKeys {
		v, ok := l.store.Get(ik)
		if !ok {
			continue
		}

		res[ik.(string)] = v.(cacheValue).value
	}
	return res
}

// DisabledCache is a cache that is always empty
type DisabledCache struct{}

var _ XdsCache = &DisabledCache{}

func (d DisabledCache) Add(key XdsCacheEntry, pushReq *PushRequest, value *discovery.Resource) {}

func (d DisabledCache) Get(XdsCacheEntry) (*discovery.Resource, bool) {
	return nil, false
}

func (d DisabledCache) Clear(configsUpdated map[ConfigKey]struct{}) {}

func (d DisabledCache) ClearAll() {}

func (d DisabledCache) Keys() []string { return nil }

func (d DisabledCache) Snapshot() map[string]*discovery.Resource { return nil }
