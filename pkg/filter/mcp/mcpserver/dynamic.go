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
	"fmt"
	"sync"
	"time"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

const (
	// DefaultDebounceTime default debounce interval
	DefaultDebounceTime = 500 * time.Millisecond

	// EmptyFingerprint fingerprint value for empty configuration
	EmptyFingerprint = "00000000"
)

var (
	globalRegistry *ToolRegistry
	globalDynamic  *DynamicConsumer

	// sync.Once variables for thread-safe singleton initialization
	registryOnce sync.Once
	dynamicOnce  sync.Once
)

// GetOrInitRegistry returns a singleton ToolRegistry
func GetOrInitRegistry() *ToolRegistry {
	registryOnce.Do(func() {
		globalRegistry = NewToolRegistry()
	})
	return globalRegistry
}

// GetOrInitDynamic returns a singleton DynamicConsumer
func GetOrInitDynamic() *DynamicConsumer {
	dynamicOnce.Do(func() {
		globalDynamic = NewDynamicConsumer(GetOrInitRegistry())
	})
	return globalDynamic
}

// DynamicConsumer applies dynamic MCP configurations into the registry
type DynamicConsumer struct {
	registry *ToolRegistry

	// lightweight fingerprint debounce fields
	mu              sync.Mutex
	lastFingerprint string
	lastApplied     time.Time
	debounceTime    time.Duration
}

func NewDynamicConsumer(reg *ToolRegistry) *DynamicConsumer {
	return &DynamicConsumer{
		registry:     reg,
		debounceTime: DefaultDebounceTime,
	}
}

// update tools from the remote config in nacos
func (d *DynamicConsumer) ApplyMcpServerConfig(cfg *model.McpServerConfig) error {
	if cfg == nil {
		return nil
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	// 1. calculate configuration fingerprint
	fingerprint := d.calculateFingerprint(cfg.Tools)

	// 2. skip if fingerprint is the same (idempotency check)
	if d.lastFingerprint == fingerprint {
		logger.Debugf("[dubbo-go-pixiu] mcp dynamic config unchanged (fp=%s), skipped", fingerprint)
		return nil
	}

	// 3. time-based debounce check
	now := time.Now()
	if !d.lastApplied.IsZero() && now.Sub(d.lastApplied) < d.debounceTime {
		logger.Debugf("[dubbo-go-pixiu] mcp dynamic debounce active (elapsed=%v), skipped", now.Sub(d.lastApplied))
		return nil
	}

	// 4. apply configuration
	d.registry.ReplaceAllTools(cfg.Tools)

	// 5. update debounce state
	d.lastFingerprint = fingerprint
	d.lastApplied = now

	logger.Infof("[dubbo-go-pixiu] mcp dynamic applied config: %d tools, fingerprint=%s",
		len(cfg.Tools), fingerprint)
	return nil
}

// calculateFingerprint calculates a lightweight fingerprint for the configuration
func (d *DynamicConsumer) calculateFingerprint(tools []model.ToolConfig) string {
	if len(tools) == 0 {
		return EmptyFingerprint
	}

	var sum uint32
	count := uint32(len(tools))

	for _, tool := range tools {
		// tool name hash
		nameHash := d.simpleStringHash(tool.Name)
		// cluster hash
		clusterHash := d.simpleStringHash(tool.Cluster)
		// arguments count
		argsCount := uint32(len(tool.Args))

		// combine hashes
		toolHash := nameHash ^ (clusterHash << 8) ^ (argsCount << 16)
		sum ^= toolHash
	}

	// combine with tool count
	fingerprint := sum ^ (count << 24)

	return fmt.Sprintf("%08x", fingerprint)
}

// simpleStringHash simple and efficient string hash function
func (d *DynamicConsumer) simpleStringHash(s string) uint32 {
	var hash uint32 = 5381
	for _, c := range s {
		hash = ((hash << 5) + hash) + uint32(c)
	}
	return hash
}

// SetDebounceTime dynamically adjusts debounce time
func (d *DynamicConsumer) SetDebounceTime(duration time.Duration) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if duration >= 0 {
		d.debounceTime = duration
		logger.Infof("[dubbo-go-pixiu] mcp dynamic debounce time updated to %v", duration)
	}
}

// GetDebounceInfo gets debounce information (for debugging)
func (d *DynamicConsumer) GetDebounceInfo() map[string]interface{} {
	d.mu.Lock()
	defer d.mu.Unlock()

	return map[string]interface{}{
		"last_fingerprint": d.lastFingerprint,
		"last_applied":     d.lastApplied,
		"debounce_time":    d.debounceTime.String(),
	}
}

// ResetDebounceState resets debounce state (mainly for testing)
func (d *DynamicConsumer) ResetDebounceState() {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.lastFingerprint = ""
	d.lastApplied = time.Time{}
	logger.Debugf("[dubbo-go-pixiu] mcp dynamic debounce state reset")
}
