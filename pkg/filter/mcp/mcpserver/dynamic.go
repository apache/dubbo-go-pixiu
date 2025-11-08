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
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"time"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/filter/mcp/mcpserver/transport"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

const (
	// DefaultDebounceTime default debounce interval
	DefaultDebounceTime = 500 * time.Millisecond

	// EmptyFingerprint fingerprint value for empty configuration
	EmptyFingerprint = "00000000"
)

// ServerToolConfig tool configuration for a single server
type ServerToolConfig struct {
	Tools       []model.ToolConfig
	Fingerprint string
	LastApplied time.Time
}

// DynamicConsumer applies dynamic MCP configurations into the registry
type DynamicConsumer struct {
	registry       *ToolRegistry
	sessionManager *transport.SessionManager
	sseHandler     *transport.SSEHandler

	// Tool configuration management grouped by server
	mu            sync.RWMutex
	serverConfigs map[string]*ServerToolConfig // serverId -> server tool configuration
	debounceTime  time.Duration
}

func NewDynamicConsumer(reg *ToolRegistry, sm *transport.SessionManager, sseHandler *transport.SSEHandler) *DynamicConsumer {
	return &DynamicConsumer{
		registry:       reg,
		sessionManager: sm,
		sseHandler:     sseHandler,
		serverConfigs:  make(map[string]*ServerToolConfig),
		debounceTime:   DefaultDebounceTime,
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

	// Notify all connected clients about tools list change
	d.notifyToolsListChanged()

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
		_, _ = fmt.Fprintf(hash, "name:%s;cluster:%s;args:%d;", tool.Name, tool.Cluster, len(tool.Args))

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

// notifyToolsListChanged sends notifications/tools/list_changed to all connected clients
func (d *DynamicConsumer) notifyToolsListChanged() {
	if d.sessionManager == nil {
		logger.Debugf("[dubbo-go-pixiu] mcp server session manager not available, skip tools list_changed notification")
		return
	}

	// Get all active sessions
	sessionIDs := d.sessionManager.AllSessionIDs()
	if len(sessionIDs) == 0 {
		logger.Debugf("[dubbo-go-pixiu] mcp server no active sessions, skip tools list_changed notification")
		return
	}

	// Send notification to each session
	successCount := 0
	for _, sessionID := range sessionIDs {
		if err := d.sendToolsListChangedNotification(sessionID); err != nil {
			logger.Warnf("[dubbo-go-pixiu] mcp server failed to send tools list_changed to session %s: %v", sessionID, err)
		} else {
			successCount++
		}
	}

	logger.Infof("[dubbo-go-pixiu] mcp server sent tools/list_changed notification to %d/%d sessions", successCount, len(sessionIDs))
}

// sendToolsListChangedNotification sends notification to a specific session
func (d *DynamicConsumer) sendToolsListChangedNotification(sessionID string) error {
	session, exists := d.sessionManager.Session(sessionID)
	if !exists {
		return fmt.Errorf("session not found")
	}

	if session.PipeWriter == nil {
		return fmt.Errorf("SSE pipe not established")
	}

	// Build tools/list_changed notification (no params needed)
	notification := map[string]any{
		"jsonrpc": "2.0",
		"method":  "notifications/tools/list_changed",
	}

	messageJSON, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	sseData := d.sseHandler.FormatSSEMessage(string(messageJSON))

	if _, err := session.PipeWriter.Write([]byte(sseData)); err != nil {
		return fmt.Errorf("failed to write to SSE pipe: %w", err)
	}

	session.LastActivity = time.Now()
	logger.Debugf("[dubbo-go-pixiu] mcp server sent tools/list_changed to session: %s", sessionID)
	return nil
}
