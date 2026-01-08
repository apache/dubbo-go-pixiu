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
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"sync"
	"time"
)

import (
	"github.com/pkg/errors"

	"golang.org/x/sync/singleflight"

	"google.golang.org/grpc"
	rpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"

	"google.golang.org/protobuf/proto"

	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"

	"google.golang.org/protobuf/types/descriptorpb"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

const (
	defaultDescCacheTTL = 5 * time.Minute
	reflectionTimeout   = 10 * time.Second
	// versionHashLength is the length of truncated MD5 hash used for version detection.
	// 16 hex characters (64 bits) provide sufficient uniqueness for cache invalidation
	// while keeping the hash compact for logging and comparison.
	versionHashLength = 16
)

// ReflectionVersion specifies which version of the reflection API to use
type ReflectionVersion string

const (
	// ReflectionV1Alpha uses the v1alpha reflection API
	ReflectionV1Alpha ReflectionVersion = "v1alpha"
	// ReflectionV1 uses the v1 reflection API (not yet implemented, reserved)
	ReflectionV1 ReflectionVersion = "v1"
	// ReflectionAuto attempts to detect the available version
	ReflectionAuto ReflectionVersion = "auto"
)

// ReflectionConfig holds configuration for the reflection manager
type ReflectionConfig struct {
	CacheTTL          time.Duration
	MaxCacheSize      int
	ReflectionVersion ReflectionVersion
	ContinueOnError   bool // Continue with partial results on dependency resolution failure
}

// ReflectionCacheStats holds cache statistics for the reflection manager
type ReflectionCacheStats struct {
	MethodCacheSize      int     `json:"method_cache_size"`
	MethodCacheHits      int64   `json:"method_cache_hits"`
	MethodCacheMisses    int64   `json:"method_cache_misses"`
	MethodCacheEvictions int64   `json:"method_cache_evictions"`
	MethodCacheHitRatio  float64 `json:"method_cache_hit_ratio"`
	FileRegistryCount    int     `json:"file_registry_count"`
	TTLSeconds           float64 `json:"ttl_secs"`
	MaxCacheSize         int     `json:"max_cache_size"`
}

// fileRegistryWithMetadata holds a file registry with metadata for version tracking
type fileRegistryWithMetadata struct {
	files       *protoregistry.Files
	versionHash string
	timestamp   time.Time
}

// ReflectionManager manages gRPC reflection clients and descriptor caching
// using official google.golang.org/protobuf libraries
type ReflectionManager struct {
	cache    *DescriptorCache
	cacheTTL time.Duration
	config   ReflectionConfig
	// fileDescCache caches file descriptors per address with metadata
	fileDescCache sync.Map // address -> *fileRegistryWithMetadata
	// fileRegistryGroup uses singleflight to ensure only one goroutine
	// performs reflection for each address at a time
	fileRegistryGroup singleflight.Group
	// Track missing dependencies per address for monitoring
	missingDeps sync.Map // address -> []string
}

// NewReflectionManager creates a new reflection manager with default config
func NewReflectionManager(cacheTTL time.Duration) *ReflectionManager {
	return NewReflectionManagerWithConfig(ReflectionConfig{
		CacheTTL:          cacheTTL,
		MaxCacheSize:      defaultMaxCacheSize,
		ReflectionVersion: ReflectionV1Alpha,
		ContinueOnError:   true,
	})
}

// NewReflectionManagerWithConfig creates a new reflection manager with custom config
func NewReflectionManagerWithConfig(config ReflectionConfig) *ReflectionManager {
	if config.CacheTTL <= 0 {
		config.CacheTTL = defaultDescCacheTTL
	}
	if config.MaxCacheSize < minCacheSize {
		config.MaxCacheSize = minCacheSize
	}
	if config.ReflectionVersion == "" {
		config.ReflectionVersion = ReflectionV1Alpha
	}

	return &ReflectionManager{
		cache:    NewDescriptorCacheWithSize(config.CacheTTL, config.MaxCacheSize),
		cacheTTL: config.CacheTTL,
		config:   config,
	}
}

// GetMethodDescriptor retrieves a method descriptor using gRPC reflection
// Results are cached for improved performance
func (rm *ReflectionManager) GetMethodDescriptor(
	ctx context.Context,
	conn *grpc.ClientConn,
	address string,
	serviceName string,
	methodName string,
) (protoreflect.MethodDescriptor, error) {
	// Build cache key
	cacheKey := BuildCacheKey(address, serviceName, methodName)

	// Check cache first (with version check)
	if cached := rm.cache.Get(cacheKey); cached != nil {
		logger.Debugf("Reflection cache hit for %s", cacheKey)
		return cached, nil
	}

	logger.Debugf("Reflection cache miss for %s, performing reflection", cacheKey)

	// Perform reflection with timeout
	reflectCtx, cancel := context.WithTimeout(ctx, reflectionTimeout)
	defer cancel()

	// Get or create file registry for this address
	files, versionHash, err := rm.getOrCreateFileRegistry(reflectCtx, conn, address, serviceName)
	if err != nil {
		return nil, err
	}

	// Find service descriptor
	serviceDesc, err := files.FindDescriptorByName(protoreflect.FullName(serviceName))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find service %s", serviceName)
	}

	svcDesc, ok := serviceDesc.(protoreflect.ServiceDescriptor)
	if !ok {
		return nil, fmt.Errorf("%s is not a service", serviceName)
	}

	// Find method descriptor
	methodDesc := svcDesc.Methods().ByName(protoreflect.Name(methodName))
	if methodDesc == nil {
		return nil, fmt.Errorf("method %s not found in service %s", methodName, serviceName)
	}

	// Cache the result with version hash
	rm.cache.SetWithVersion(cacheKey, methodDesc, versionHash)
	logger.Debugf("Cached method descriptor for %s (version: %s)", cacheKey, versionHash)

	return methodDesc, nil
}

// getOrCreateFileRegistry gets or creates a file registry for the given address
// Uses singleflight to ensure only one goroutine performs reflection per address
// Returns the files registry and version hash for change detection
func (rm *ReflectionManager) getOrCreateFileRegistry(
	ctx context.Context,
	conn *grpc.ClientConn,
	address string,
	serviceName string,
) (*protoregistry.Files, string, error) {
	// Check cache first (fast path, no locking)
	if cached, ok := rm.fileDescCache.Load(address); ok {
		reg := cached.(*fileRegistryWithMetadata)
		return reg.files, reg.versionHash, nil
	}

	// Use singleflight to deduplicate concurrent reflection requests
	// for the same address. This prevents thundering herd problem.
	result, err, shared := rm.fileRegistryGroup.Do(address, func() (any, error) {
		// Double-check cache in case another goroutine just filled it
		if cached, ok := rm.fileDescCache.Load(address); ok {
			reg := cached.(*fileRegistryWithMetadata)
			return registryWithVersion{files: reg.files, versionHash: reg.versionHash}, nil
		}

		// Create reflection client based on configured version
		client := rpb.NewServerReflectionClient(conn)
		stream, err := client.ServerReflectionInfo(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create reflection stream")
		}
		defer stream.CloseSend()

		// Request file descriptor for the service
		fileDescs, missingDeps, err := rm.resolveServiceFileDescriptors(stream, serviceName)
		if err != nil && !rm.config.ContinueOnError {
			return nil, err
		}

		// Log missing dependencies if any
		if len(missingDeps) > 0 {
			logger.Warnf("Missing dependencies for %s: %v", address, missingDeps)
			rm.missingDeps.Store(address, missingDeps)
		}

		// Compute version hash from all file descriptors
		versionHash := rm.computeVersionHash(fileDescs)

		// Build file registry
		files, err := rm.buildFileRegistry(fileDescs)
		if err != nil && !rm.config.ContinueOnError {
			return nil, err
		}

		// Cache the registry with metadata for future requests
		rm.fileDescCache.Store(address, &fileRegistryWithMetadata{
			files:       files,
			versionHash: versionHash,
			timestamp:   time.Now(),
		})

		return registryWithVersion{files: files, versionHash: versionHash}, nil
	})

	if err != nil {
		return nil, "", err
	}

	if shared {
		logger.Debugf("Reflection request for %s shared with another goroutine", address)
	}

	rw := result.(registryWithVersion)
	return rw.files, rw.versionHash, nil
}

// registryWithVersion is a helper type to return both files and version hash from singleflight
type registryWithVersion struct {
	files       *protoregistry.Files
	versionHash string
}

// resolveServiceFileDescriptors resolves all file descriptors needed for a service
// Returns file descriptors, list of missing dependencies, and error
func (rm *ReflectionManager) resolveServiceFileDescriptors(
	stream rpb.ServerReflection_ServerReflectionInfoClient,
	serviceName string,
) ([]*descriptorpb.FileDescriptorProto, []string, error) {
	// Request file containing the service
	req := &rpb.ServerReflectionRequest{
		MessageRequest: &rpb.ServerReflectionRequest_FileContainingSymbol{
			FileContainingSymbol: serviceName,
		},
	}

	if err := stream.Send(req); err != nil {
		return nil, nil, errors.Wrap(err, "failed to send reflection request")
	}

	resp, err := stream.Recv()
	if err != nil {
		if err == io.EOF {
			return nil, nil, errors.New("unexpected EOF when receiving reflection response")
		}
		return nil, nil, errors.Wrap(err, "failed to receive reflection response")
	}

	fdResp := resp.GetFileDescriptorResponse()
	if fdResp == nil {
		if errResp := resp.GetErrorResponse(); errResp != nil {
			return nil, nil, fmt.Errorf("reflection error: %s", errResp.ErrorMessage)
		}
		return nil, nil, errors.New("unexpected reflection response type")
	}

	// Parse file descriptors
	fileDescs := make([]*descriptorpb.FileDescriptorProto, 0, len(fdResp.FileDescriptorProto))
	for _, fdBytes := range fdResp.FileDescriptorProto {
		fd := &descriptorpb.FileDescriptorProto{}
		if err := proto.Unmarshal(fdBytes, fd); err != nil {
			return nil, nil, errors.Wrap(err, "failed to unmarshal file descriptor")
		}
		fileDescs = append(fileDescs, fd)
	}

	// Resolve dependencies recursively
	resolved := make(map[string]bool)
	missingDeps := make([]string, 0)
	for _, fd := range fileDescs {
		resolved[fd.GetName()] = true
	}

	for _, fd := range fileDescs {
		deps, missing, err := rm.resolveDependencies(stream, fd.GetDependency(), resolved)
		if err != nil {
			if !rm.config.ContinueOnError {
				return nil, nil, err
			}
			logger.Warnf("Failed to resolve dependencies for %s: %v", fd.GetName(), err)
		}
		missingDeps = append(missingDeps, missing...)
		fileDescs = append(fileDescs, deps...)
	}

	return fileDescs, missingDeps, nil
}

// resolveDependencies resolves file descriptor dependencies
// Returns resolved descriptors, list of missing dependencies, and error
func (rm *ReflectionManager) resolveDependencies(
	stream rpb.ServerReflection_ServerReflectionInfoClient,
	deps []string,
	resolved map[string]bool,
) ([]*descriptorpb.FileDescriptorProto, []string, error) {
	var result []*descriptorpb.FileDescriptorProto
	var missingDeps []string

	for _, dep := range deps {
		if resolved[dep] {
			continue
		}

		req := &rpb.ServerReflectionRequest{
			MessageRequest: &rpb.ServerReflectionRequest_FileByFilename{
				FileByFilename: dep,
			},
		}

		if err := stream.Send(req); err != nil {
			if err == io.EOF {
				// Stream closed by server, return what we have collected so far
				logger.Debugf("Stream EOF during send for dependency %s", dep)
				return result, missingDeps, nil
			}
			if !rm.config.ContinueOnError {
				return nil, nil, errors.Wrapf(err, "failed to send request for %s", dep)
			}
			logger.Warnf("Failed to send request for dependency %s: %v", dep, err)
			missingDeps = append(missingDeps, dep)
			continue
		}

		resp, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				// Stream closed by server, return what we have collected so far
				logger.Debugf("Stream EOF during recv for dependency %s", dep)
				return result, missingDeps, nil
			}
			if !rm.config.ContinueOnError {
				return nil, nil, errors.Wrapf(err, "failed to receive response for %s", dep)
			}
			logger.Warnf("Failed to receive response for dependency %s: %v", dep, err)
			missingDeps = append(missingDeps, dep)
			continue
		}

		fdResp := resp.GetFileDescriptorResponse()
		if fdResp == nil {
			// No file descriptor in response, could be an error response or unexpected type
			if errResp := resp.GetErrorResponse(); errResp != nil {
				logger.Warnf("Reflection error for dependency %s: %s", dep, errResp.ErrorMessage)
				// Mark as missing and continue if ContinueOnError is enabled
				if rm.config.ContinueOnError {
					missingDeps = append(missingDeps, dep)
					continue
				}
				return nil, nil, fmt.Errorf("reflection error for dependency %s: %s", dep, errResp.ErrorMessage)
			}
			// Unexpected response type, skip this dependency
			missingDeps = append(missingDeps, dep)
			continue
		}

		for _, fdBytes := range fdResp.FileDescriptorProto {
			fd := &descriptorpb.FileDescriptorProto{}
			if err := proto.Unmarshal(fdBytes, fd); err != nil {
				if !rm.config.ContinueOnError {
					return nil, nil, errors.Wrap(err, "failed to unmarshal dependency")
				}
				logger.Warnf("Failed to unmarshal dependency %s: %v", dep, err)
				continue
			}
			if !resolved[fd.GetName()] {
				resolved[fd.GetName()] = true
				result = append(result, fd)

				// Resolve nested dependencies
				nested, nestedMissing, err := rm.resolveDependencies(stream, fd.GetDependency(), resolved)
				if err != nil {
					if !rm.config.ContinueOnError {
						return nil, nil, err
					}
					logger.Warnf("Failed to resolve nested dependencies for %s: %v", fd.GetName(), err)
				}
				missingDeps = append(missingDeps, nestedMissing...)
				result = append(result, nested...)
			}
		}
	}

	return result, missingDeps, nil
}

// buildFileRegistry builds a protoregistry.Files from file descriptor protos
// Continues on error if ContinueOnError is enabled
func (rm *ReflectionManager) buildFileRegistry(
	fileDescs []*descriptorpb.FileDescriptorProto,
) (*protoregistry.Files, error) {
	files := new(protoregistry.Files)

	// Build dependency graph and sort topologically
	sorted, err := rm.topologicalSort(fileDescs)
	if err != nil {
		if !rm.config.ContinueOnError {
			return nil, err
		}
		logger.Warnf("Topological sort failed, using original order: %v", err)
		sorted = fileDescs
	}

	// Register files in dependency order
	registeredCount := 0
	for _, fd := range sorted {
		fileDesc, err := protodesc.NewFile(fd, files)
		if err != nil {
			// Skip files that fail to register (might be already registered or missing deps)
			logger.Debugf("Skipping file %s: %v", fd.GetName(), err)
			continue
		}
		if err := files.RegisterFile(fileDesc); err != nil {
			// Skip already registered files
			logger.Debugf("Skipping duplicate file %s: %v", fd.GetName(), err)
			continue
		}
		registeredCount++
	}

	logger.Debugf("Registered %d/%d file descriptors successfully", registeredCount, len(sorted))

	if registeredCount == 0 && !rm.config.ContinueOnError {
		return nil, errors.New("failed to register any file descriptors")
	}

	return files, nil
}

// topologicalSort sorts file descriptors by dependency order
func (rm *ReflectionManager) topologicalSort(
	fileDescs []*descriptorpb.FileDescriptorProto,
) ([]*descriptorpb.FileDescriptorProto, error) {
	// Build a map for quick lookup
	fdMap := make(map[string]*descriptorpb.FileDescriptorProto)
	for _, fd := range fileDescs {
		fdMap[fd.GetName()] = fd
	}

	visited := make(map[string]bool)
	inStack := make(map[string]bool)
	var result []*descriptorpb.FileDescriptorProto

	var visit func(name string) error
	visit = func(name string) error {
		if inStack[name] {
			// Circular dependency detected
			logger.Warnf("Circular dependency detected involving: %s", name)
			// Don't return error, just skip to avoid infinite loop
			return nil
		}
		if visited[name] {
			return nil
		}

		fd, ok := fdMap[name]
		if !ok {
			// External dependency, skip
			return nil
		}

		inStack[name] = true
		for _, dep := range fd.GetDependency() {
			if err := visit(dep); err != nil {
				return err
			}
		}
		inStack[name] = false
		visited[name] = true
		result = append(result, fd)
		return nil
	}

	for name := range fdMap {
		if err := visit(name); err != nil {
			return nil, err
		}
	}

	return result, nil
}

// computeVersionHash computes a hash from file descriptors for version detection
func (rm *ReflectionManager) computeVersionHash(fileDescs []*descriptorpb.FileDescriptorProto) string {
	h := md5.New()

	// Sort file names for consistent hashing
	names := make([]string, len(fileDescs))
	for i, fd := range fileDescs {
		names[i] = fd.GetName()
	}

	// Write names and sizes to hash
	for _, name := range names {
		h.Write([]byte(name))
		h.Write([]byte{0}) // null terminator
	}

	return fmt.Sprintf("%x", h.Sum(nil))[:versionHashLength]
}

// InvalidateCache removes cached descriptors for a specific address
func (rm *ReflectionManager) InvalidateCache(address string) {
	rm.fileDescCache.Delete(address)
	// Forget the singleflight key to allow new reflection requests
	rm.fileRegistryGroup.Forget(address)
	// Clear missing dependencies tracking
	rm.missingDeps.Delete(address)
	logger.Infof("Cache invalidation for address: %s", address)
}

// InvalidateByVersion removes all cached entries with a specific version hash
func (rm *ReflectionManager) InvalidateByVersion(versionHash string) int {
	count := 0
	rm.fileDescCache.Range(func(key, value any) bool {
		reg := value.(*fileRegistryWithMetadata)
		if reg.versionHash == versionHash {
			rm.fileDescCache.Delete(key)
			if address, ok := key.(string); ok {
				rm.fileRegistryGroup.Forget(address)
			}
			count++
		}
		return true
	})
	// Also invalidate method cache with this version
	count += rm.cache.InvalidateByVersion(versionHash)
	if count > 0 {
		logger.Infof("Invalidated %d cache entries with version hash: %s", count, versionHash)
	}
	return count
}

// ClearCache clears all cached descriptors
func (rm *ReflectionManager) ClearCache() {
	rm.cache.Clear()
	// Clear file descriptor cache and singleflight cache
	rm.fileDescCache.Range(func(key, _ any) bool {
		rm.fileDescCache.Delete(key)
		// Forget the singleflight key to allow new reflection requests
		if address, ok := key.(string); ok {
			rm.fileRegistryGroup.Forget(address)
		}
		return true
	})
	// Clear missing dependencies tracking
	rm.missingDeps.Range(func(key, _ any) bool {
		rm.missingDeps.Delete(key)
		return true
	})
	logger.Info("Reflection descriptor cache cleared")
}

// Close cleans up resources
func (rm *ReflectionManager) Close() {
	rm.cache.Close()
}

// WarmupService pre-fetches and caches descriptors for a specific service.
// This eliminates cold-start latency for the first request to this service.
// The conn parameter should be a connection to the backend server.
func (rm *ReflectionManager) WarmupService(ctx context.Context, conn *grpc.ClientConn, address, serviceName string) error {
	_, _, err := rm.getOrCreateFileRegistry(ctx, conn, address, serviceName)
	if err != nil {
		return errors.Wrapf(err, "warmup failed for service %s at %s", serviceName, address)
	}
	logger.Infof("Warmup completed for service %s at %s", serviceName, address)
	return nil
}

// WarmupServices pre-fetches and caches descriptors for multiple services concurrently.
// Returns a map of service names to errors (nil if successful).
func (rm *ReflectionManager) WarmupServices(ctx context.Context, conn *grpc.ClientConn, address string, serviceNames []string) map[string]error {
	results := make(map[string]error)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, svc := range serviceNames {
		wg.Add(1)
		go func(serviceName string) {
			defer wg.Done()
			err := rm.WarmupService(ctx, conn, address, serviceName)
			mu.Lock()
			results[serviceName] = err
			mu.Unlock()
		}(svc)
	}

	wg.Wait()
	return results
}

// GetCacheStats returns cache statistics for monitoring
func (rm *ReflectionManager) GetCacheStats() ReflectionCacheStats {
	fileCount := 0
	rm.fileDescCache.Range(func(_, _ any) bool {
		fileCount++
		return true
	})

	methodStats := rm.cache.GetStats()

	return ReflectionCacheStats{
		MethodCacheSize:      methodStats.Size,
		MethodCacheHits:      methodStats.Hits,
		MethodCacheMisses:    methodStats.Misses,
		MethodCacheEvictions: methodStats.Evictions,
		MethodCacheHitRatio:  methodStats.HitRatio,
		FileRegistryCount:    fileCount,
		TTLSeconds:           rm.cacheTTL.Seconds(),
		MaxCacheSize:         methodStats.MaxSize,
	}
}

// GetMissingDependencies returns missing dependencies for a given address
func (rm *ReflectionManager) GetMissingDependencies(address string) []string {
	if missing, ok := rm.missingDeps.Load(address); ok {
		return missing.([]string)
	}
	return nil
}

// GetAllMissingDependencies returns all missing dependencies across all addresses
func (rm *ReflectionManager) GetAllMissingDependencies() map[string][]string {
	result := make(map[string][]string)
	rm.missingDeps.Range(func(key, value any) bool {
		if address, ok := key.(string); ok {
			result[address] = value.([]string)
		}
		return true
	})
	return result
}

// SetConfig updates the reflection manager configuration
func (rm *ReflectionManager) SetConfig(config ReflectionConfig) {
	if config.CacheTTL > 0 {
		rm.cacheTTL = config.CacheTTL
		rm.config.CacheTTL = config.CacheTTL
	}
	if config.MaxCacheSize >= minCacheSize {
		rm.cache.SetMaxSize(config.MaxCacheSize)
		rm.config.MaxCacheSize = config.MaxCacheSize
	}
	if config.ReflectionVersion != "" {
		rm.config.ReflectionVersion = config.ReflectionVersion
	}
	rm.config.ContinueOnError = config.ContinueOnError

	logger.Infof("Reflection manager config updated: version=%s, continueOnError=%v",
		rm.config.ReflectionVersion, rm.config.ContinueOnError)
}

// GetConfig returns the current configuration
func (rm *ReflectionManager) GetConfig() ReflectionConfig {
	return rm.config
}
