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
	"sort"
	"sync"
	"time"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/filter/mcp/mcpserver/router"
	"github.com/apache/dubbo-go-pixiu/pkg/filter/mcp/mcpserver/transport"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

const (
	// DefaultDebounceTime default debounce interval
	DefaultDebounceTime = 500 * time.Millisecond
)

var errDynamicRouterUpdateUnsupported = fmt.Errorf("dynamic MCP updates support tool catalog changes only; router changes require filter rebuild")

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

// ApplyMcpServerConfigByServer applies a dynamic tool catalog update by server
// ID. Router configuration is intentionally not applied on this path; a dynamic
// payload that contains a router section is rejected before mutating tools so
// the runtime cannot enter a tools=new/router=old partial state.
func (d *DynamicConsumer) ApplyMcpServerConfigByServer(serverId string, cfg *model.McpServerConfig) error {
	if cfg == nil {
		return d.removeServerConfig(serverId)
	}
	if cfg.Router != nil {
		return errDynamicRouterUpdateUnsupported
	}
	if err := router.ValidateTools(cfg.Tools); err != nil {
		return fmt.Errorf("invalid mcp tool router metadata: %w", err)
	}

	d.mu.Lock()

	// 1. Calculate new configuration fingerprint
	fingerprint := d.calculateFingerprint(cfg.Tools)

	// 2. Check if the server's configuration really needs to be updated
	if existingConfig, exists := d.serverConfigs[serverId]; exists {
		if existingConfig.Fingerprint == fingerprint {
			d.mu.Unlock()
			logger.Debugf("[dubbo-go-pixiu] mcp server %s config unchanged (fp=%s), skipped", serverId, fingerprint)
			return nil
		}
	}

	// 3. Debounce check (based on this server's configuration change time)
	now := time.Now()
	if existingConfig, exists := d.serverConfigs[serverId]; exists {
		// Skip only if this server is within debounce time
		if !existingConfig.LastApplied.IsZero() && now.Sub(existingConfig.LastApplied) < d.debounceTime {
			d.mu.Unlock()
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
	mergedFingerprint := d.calculateFingerprint(mergedTools)
	if err := d.applyMergedConfig(mergedTools, mergedFingerprint); err != nil {
		// Rollback
		if oldConfig != nil {
			d.serverConfigs[serverId] = oldConfig
		} else {
			delete(d.serverConfigs, serverId)
		}
		d.mu.Unlock()
		return err
	}
	serverCount := len(d.serverConfigs)
	mergedCount := len(mergedTools)
	d.mu.Unlock()

	logger.Infof("[dubbo-go-pixiu] mcp server %s config applied: %d tools, total servers: %d, merged tools: %d",
		serverId, len(cfg.Tools), serverCount, mergedCount)

	// Notify all connected clients about tools list change
	d.notifyToolsListChanged()

	return nil
}

// calculateFingerprint returns the registry-level tool catalog fingerprint.
func (d *DynamicConsumer) calculateFingerprint(tools []model.ToolConfig) string {
	return toolCatalogFingerprint(tools)
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

// applyMergedConfig applies merged configuration to the registry.
func (d *DynamicConsumer) applyMergedConfig(tools []model.ToolConfig, fingerprint string) error {
	return d.registry.replaceAllToolsWithFingerprint(tools, fingerprint)
}

// removeServerConfig removes server configuration
func (d *DynamicConsumer) removeServerConfig(serverId string) error {
	d.mu.Lock()

	if _, exists := d.serverConfigs[serverId]; !exists {
		d.mu.Unlock()
		return nil // Already does not exist
	}

	oldConfig := d.serverConfigs[serverId]
	delete(d.serverConfigs, serverId)

	// Recalculate and apply merged configuration
	mergedTools := d.calculateCurrentMergedTools()
	mergedFingerprint := d.calculateFingerprint(mergedTools)
	if err := d.applyMergedConfig(mergedTools, mergedFingerprint); err != nil {
		d.serverConfigs[serverId] = oldConfig
		d.mu.Unlock()
		return err
	}
	serverCount := len(d.serverConfigs)
	d.mu.Unlock()

	logger.Infof("[dubbo-go-pixiu] mcp server %s config removed, remaining servers: %d",
		serverId, serverCount)
	d.notifyToolsListChanged()

	return nil
}

// calculateCurrentMergedTools calculates merged tools from all current servers
func (d *DynamicConsumer) calculateCurrentMergedTools() []model.ToolConfig {
	var allTools []model.ToolConfig

	serverIDs := make([]string, 0, len(d.serverConfigs))
	for serverID := range d.serverConfigs {
		serverIDs = append(serverIDs, serverID)
	}
	sort.Strings(serverIDs)

	for _, serverID := range serverIDs {
		allTools = append(allTools, d.serverConfigs[serverID].Tools...)
	}

	return allTools
}

// notifyToolsListChanged marks notifications/tools/list_changed for all active
// sessions. Offline sessions keep a bounded pending version for reconnect.
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

	markedCount := 0
	flushedCount := 0
	for _, sessionID := range sessionIDs {
		session, exists := d.sessionManager.GetSession(sessionID)
		if !exists {
			continue
		}
		session.MarkToolsListChangedPending()
		markedCount++
		if err := d.flushToolsListChangedNotification(session); err != nil {
			logger.Warnf("[dubbo-go-pixiu] mcp server failed to send tools list_changed: %v", err)
		} else {
			flushedCount++
		}
	}

	if markedCount == 0 {
		logger.Debugf("[dubbo-go-pixiu] mcp server no active sessions to mark for tools list_changed notification")
		return
	}
	logger.Infof("[dubbo-go-pixiu] mcp server marked tools/list_changed notification for %d sessions, flushed %d online sessions", markedCount, flushedCount)
}

// flushToolsListChangedNotification sends pending notifications to a specific
// session when an SSE stream is attached.
func (d *DynamicConsumer) flushToolsListChangedNotification(session *transport.MCPSession) error {
	if d.sseHandler == nil {
		return fmt.Errorf("SSE handler not configured")
	}
	for {
		version, pending := session.PendingToolsListChangedVersion()
		if !pending || !session.HasPipeWriter() {
			return nil
		}
		notification := map[string]any{
			"jsonrpc": "2.0",
			"method":  "notifications/tools/list_changed",
		}
		if err := d.sseHandler.SendSSEMessage(session, notification); err != nil {
			return err
		}
		session.MarkToolsListChangedNotified(version)
		logger.Debugf("[dubbo-go-pixiu] mcp server sent tools/list_changed")
	}
}
