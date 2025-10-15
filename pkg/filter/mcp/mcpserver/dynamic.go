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
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
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

// ServerToolConfig tool configuration for a single server
type ServerToolConfig struct {
	Tools       []model.ToolConfig
	Fingerprint string
	LastApplied time.Time
}

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

	// Tool configuration management grouped by server
	mu            sync.RWMutex
	serverConfigs map[string]*ServerToolConfig // serverId -> server tool configuration
	debounceTime  time.Duration
}

func NewDynamicConsumer(reg *ToolRegistry) *DynamicConsumer {
	return &DynamicConsumer{
		registry:      reg,
		serverConfigs: make(map[string]*ServerToolConfig),
		debounceTime:  DefaultDebounceTime,
	}
}

// ApplyMcpServerConfigByServer applies configuration by server ID
func (d *DynamicConsumer) ApplyMcpServerConfigByServer(serverId string, cfg *model.McpServerConfig) error {
	if cfg == nil {
		return d.removeServerConfig(serverId)
	}

	d.mu.Lock()
	defer d.mu.Unlock()

	// 1. Calculate new configuration fingerprint
	fingerprint := d.calculateFingerprint(cfg.Tools)

	// 2. Check if the server's configuration really needs to be updated
	if existingConfig, exists := d.serverConfigs[serverId]; exists {
		if existingConfig.Fingerprint == fingerprint {
			logger.Debugf("[dubbo-go-pixiu] mcp server %s config unchanged (fp=%s), skipped", serverId, fingerprint)
			return nil
		}
	}

	// 3. Debounce check (based on this server's configuration change time)
	now := time.Now()
	if existingConfig, exists := d.serverConfigs[serverId]; exists {
		// Skip only if this server is within debounce time
		if !existingConfig.LastApplied.IsZero() && now.Sub(existingConfig.LastApplied) < d.debounceTime {
			logger.Debugf("[dubbo-go-pixiu] mcp server %s debounce active (elapsed=%v), skipped", serverId, now.Sub(existingConfig.LastApplied))
			return nil
		}
	}

	// 4. Fully replace the server's tool configuration
	oldConfig := d.serverConfigs[serverId]
	serverConfig := &ServerToolConfig{
		Tools:       make([]model.ToolConfig, len(cfg.Tools)),
		Fingerprint: fingerprint,
		LastApplied: now,
	}
	copy(serverConfig.Tools, cfg.Tools)
	d.serverConfigs[serverId] = serverConfig

	// 5. Recalculate merged tools from all servers and apply to registry
	mergedTools := d.calculateCurrentMergedTools()
	if err := d.applyMergedConfig(mergedTools); err != nil {
		// Rollback
		if oldConfig != nil {
			d.serverConfigs[serverId] = oldConfig
		} else {
			delete(d.serverConfigs, serverId)
		}
		return err
	}

	logger.Infof("[dubbo-go-pixiu] mcp server %s config applied: %d tools, total servers: %d, merged tools: %d",
		serverId, len(cfg.Tools), len(d.serverConfigs), len(mergedTools))

	return nil
}

// calculateFingerprint calculates a robust fingerprint for the configuration using SHA256
func (d *DynamicConsumer) calculateFingerprint(tools []model.ToolConfig) string {
	if len(tools) == 0 {
		return EmptyFingerprint
	}

	// Create a sorted list of tools for consistent hashing
	sortedTools := make([]model.ToolConfig, len(tools))
	copy(sortedTools, tools)
	sort.Slice(sortedTools, func(i, j int) bool {
		if sortedTools[i].Name != sortedTools[j].Name {
			return sortedTools[i].Name < sortedTools[j].Name
		}
		return sortedTools[i].Cluster < sortedTools[j].Cluster
	})

	// Build hash input string
	hash := sha256.New()
	for _, tool := range sortedTools {
		fmt.Fprintf(hash, "name:%s;cluster:%s;args:%d;", tool.Name, tool.Cluster, len(tool.Args))
	}

	// Return first 8 characters of hex encoded hash
	fullHash := hex.EncodeToString(hash.Sum(nil))
	return fullHash[:8]
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

// GetDebounceInfo returns debounce state information (for debugging/monitoring)
func (d *DynamicConsumer) GetDebounceInfo() map[string]any {
	d.mu.RLock()
	defer d.mu.RUnlock()

	return map[string]any{
		"debounce_time": d.debounceTime.String(),
		"server_count":  len(d.serverConfigs),
	}
}

// ResetDebounceState resets debounce state (mainly for testing)
func (d *DynamicConsumer) ResetDebounceState() {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Clear all server configurations
	d.serverConfigs = make(map[string]*ServerToolConfig)
	logger.Debugf("[dubbo-go-pixiu] mcp dynamic debounce state reset")
}

// applyMergedConfig applies merged configuration to the registry
func (d *DynamicConsumer) applyMergedConfig(tools []model.ToolConfig) error {
	d.registry.ReplaceAllTools(tools)
	return nil
}

// removeServerConfig removes server configuration
func (d *DynamicConsumer) removeServerConfig(serverId string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, exists := d.serverConfigs[serverId]; !exists {
		return nil // Already does not exist
	}

	delete(d.serverConfigs, serverId)

	// Recalculate and apply merged configuration
	mergedTools := d.calculateCurrentMergedTools()
	if err := d.applyMergedConfig(mergedTools); err != nil {
		return err
	}

	logger.Infof("[dubbo-go-pixiu] mcp server %s config removed, remaining servers: %d",
		serverId, len(d.serverConfigs))

	return nil
}

// calculateCurrentMergedTools calculates merged tools from all current servers
func (d *DynamicConsumer) calculateCurrentMergedTools() []model.ToolConfig {
	var allTools []model.ToolConfig

	// Simply accumulate tools from all servers
	for _, config := range d.serverConfigs {
		allTools = append(allTools, config.Tools...)
	}

	return allTools
}
