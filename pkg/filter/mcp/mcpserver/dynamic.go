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
	"context"
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

type serverConfigApplySkipReason string

const (
	serverConfigApplySkippedUnchanged serverConfigApplySkipReason = "unchanged"
)

// ServerToolConfig tool configuration for a single server
type ServerToolConfig struct {
	Tools       []model.ToolConfig
	Fingerprint string
	LastApplied time.Time
}

type serverConfigApplyResult struct {
	SkippedReason serverConfigApplySkipReason
	Fingerprint   string
	Elapsed       time.Duration
	ServerCount   int
	MergedCount   int
}

// DynamicConsumer applies dynamic MCP configurations into the registry
type DynamicConsumer struct {
	registry       *ToolRegistry
	sessionManager *transport.SessionManager
	sseHandler     *transport.SSEHandler
	selector       router.ToolSelector
	plans          *router.SessionPlanStore
	governance     bool
	runtimeID      string

	// Tool configuration management grouped by server
	mu            sync.RWMutex
	serverConfigs map[ServerSource]*ServerToolConfig
	debounceTime  time.Duration
}

func NewDynamicConsumer(reg *ToolRegistry, sm *transport.SessionManager, sseHandler *transport.SSEHandler) *DynamicConsumer {
	return &DynamicConsumer{
		registry:       reg,
		sessionManager: sm,
		sseHandler:     sseHandler,
		serverConfigs:  make(map[ServerSource]*ServerToolConfig),
		debounceTime:   DefaultDebounceTime,
	}
}

func (d *DynamicConsumer) RuntimeID() string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.runtimeID
}

func (d *DynamicConsumer) setRuntimeID(id string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.runtimeID = id
}

func (d *DynamicConsumer) SetGovernance(selector router.ToolSelector, plans *router.SessionPlanStore) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.selector = selector
	d.plans = plans
	d.governance = selector != nil && plans != nil
}

// ApplyMcpServerConfigByServer applies a dynamic tool catalog update by server
// ID. Router configuration is intentionally not applied on this path; a dynamic
// payload that contains a router section is rejected before mutating tools so
// the runtime cannot enter a tools=new/router=old partial state.
func (d *DynamicConsumer) ApplyMcpServerConfigByServer(serverId string, cfg *model.McpServerConfig) error {
	return d.ApplyMcpServerConfigBySource(NewServerSource(DefaultRegistryName, serverId), cfg)
}

// ApplyMcpServerConfigBySource applies a dynamic tool catalog update by stable
// registry source identity. Router configuration is rebuild-only and rejected on
// this path before mutating catalog state.
func (d *DynamicConsumer) ApplyMcpServerConfigBySource(source ServerSource, cfg *model.McpServerConfig) error {
	source = source.Normalize()
	if cfg == nil {
		return d.removeServerConfig(source)
	}
	if err := d.validateDynamicConfig(cfg); err != nil {
		return err
	}

	result, err := d.applyServerTools(source, cfg.Tools)
	if err != nil {
		return err
	}
	if d.logSkippedServerConfigApply(source, result) {
		return nil
	}

	logger.Infof("[dubbo-go-pixiu] mcp server source %s config applied: %d tools, total sources: %d, merged tools: %d",
		source, len(cfg.Tools), result.ServerCount, result.MergedCount)

	d.notifyToolsListChanged()

	return nil
}

func (d *DynamicConsumer) validateDynamicConfig(cfg *model.McpServerConfig) error {
	if cfg.Router != nil {
		return errDynamicRouterUpdateUnsupported
	}
	if !d.governanceEnabled() {
		return nil
	}
	if err := router.ValidateTools(cfg.Tools); err != nil {
		return fmt.Errorf("invalid mcp tool router metadata: %w", err)
	}
	return nil
}

func (d *DynamicConsumer) applyServerTools(source ServerSource, tools []model.ToolConfig) (serverConfigApplyResult, error) {
	source = source.Normalize()
	result := serverConfigApplyResult{Fingerprint: d.calculateFingerprint(tools)}
	now := time.Now()

	d.mu.Lock()
	defer d.mu.Unlock()

	existingConfig := d.serverConfigs[source]
	if existingConfig != nil && existingConfig.Fingerprint == result.Fingerprint {
		result.SkippedReason = serverConfigApplySkippedUnchanged
		result.Elapsed = now.Sub(existingConfig.LastApplied)
		return result, nil
	}

	oldConfig := d.serverConfigs[source]
	serverConfig := &ServerToolConfig{
		Tools:       make([]model.ToolConfig, len(tools)),
		Fingerprint: result.Fingerprint,
		LastApplied: now,
	}
	for i := range tools {
		serverConfig.Tools[i] = *tools[i].DeepCopy()
	}
	d.serverConfigs[source] = serverConfig

	mergedTools := d.calculateCurrentMergedTools()
	mergedFingerprint := d.calculateFingerprint(mergedTools)
	if err := d.applyMergedConfig(mergedTools, mergedFingerprint); err != nil {
		d.restoreServerConfig(source, oldConfig)
		return result, err
	}
	result.ServerCount = len(d.serverConfigs)
	result.MergedCount = len(mergedTools)
	return result, nil
}

func (d *DynamicConsumer) restoreServerConfig(source ServerSource, oldConfig *ServerToolConfig) {
	if oldConfig != nil {
		d.serverConfigs[source] = oldConfig
		return
	}
	delete(d.serverConfigs, source)
}

func (d *DynamicConsumer) logSkippedServerConfigApply(source ServerSource, result serverConfigApplyResult) bool {
	switch result.SkippedReason {
	case serverConfigApplySkippedUnchanged:
		logger.Debugf("[dubbo-go-pixiu] mcp server source %s config unchanged (fp=%s, elapsed=%v), skipped", source, result.Fingerprint, result.Elapsed)
		return true
	default:
		return false
	}
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
	d.serverConfigs = make(map[ServerSource]*ServerToolConfig)
	logger.Debugf("[dubbo-go-pixiu] mcp dynamic debounce state reset")
}

// applyMergedConfig applies merged configuration to the registry.
func (d *DynamicConsumer) applyMergedConfig(tools []model.ToolConfig, fingerprint string) error {
	return d.registry.replaceAllToolsWithFingerprint(tools, fingerprint)
}

// removeServerConfig removes server configuration
func (d *DynamicConsumer) removeServerConfig(source ServerSource) error {
	source = source.Normalize()
	d.mu.Lock()

	if _, exists := d.serverConfigs[source]; !exists {
		d.mu.Unlock()
		return nil // Already does not exist
	}

	oldConfig := d.serverConfigs[source]
	delete(d.serverConfigs, source)

	// Recalculate and apply merged configuration
	mergedTools := d.calculateCurrentMergedTools()
	mergedFingerprint := d.calculateFingerprint(mergedTools)
	if err := d.applyMergedConfig(mergedTools, mergedFingerprint); err != nil {
		d.serverConfigs[source] = oldConfig
		d.mu.Unlock()
		return err
	}
	serverCount := len(d.serverConfigs)
	d.mu.Unlock()

	logger.Infof("[dubbo-go-pixiu] mcp server source %s config removed, remaining sources: %d",
		source, serverCount)
	d.notifyToolsListChanged()

	return nil
}

// RemoveAllMcpServerConfigs removes all dynamic registry publications from the
// runtime catalog. It is used by adapter shutdown/reconfiguration cleanup.
func (d *DynamicConsumer) RemoveAllMcpServerConfigs() error {
	d.mu.Lock()
	if len(d.serverConfigs) == 0 {
		d.mu.Unlock()
		return nil
	}
	oldConfigs := d.serverConfigs
	d.serverConfigs = make(map[ServerSource]*ServerToolConfig)
	if err := d.applyMergedConfig(nil, EmptyFingerprint); err != nil {
		d.serverConfigs = oldConfigs
		d.mu.Unlock()
		return err
	}
	d.mu.Unlock()

	logger.Infof("[dubbo-go-pixiu] mcp dynamic server configs removed")
	d.notifyToolsListChanged()
	return nil
}

// calculateCurrentMergedTools calculates merged tools from all current servers
func (d *DynamicConsumer) calculateCurrentMergedTools() []model.ToolConfig {
	var allTools []model.ToolConfig

	sources := make([]ServerSource, 0, len(d.serverConfigs))
	for source := range d.serverConfigs {
		sources = append(sources, source)
	}
	sort.Slice(sources, func(i, j int) bool {
		return sources[i].Key() < sources[j].Key()
	})

	for _, source := range sources {
		allTools = append(allTools, d.serverConfigs[source].Tools...)
	}

	return allTools
}

// notifyToolsListChanged marks notifications/tools/list_changed only for
// sessions whose visible tool set changed. Offline sessions keep one bounded
// pending version for reconnect.
func (d *DynamicConsumer) notifyToolsListChanged() {
	if d.sessionManager == nil {
		logger.Debugf("[dubbo-go-pixiu] mcp server session manager not available, skip tools list_changed notification")
		return
	}

	sessions := d.sessionManager.SnapshotSessions()
	if len(sessions) == 0 {
		logger.Debugf("[dubbo-go-pixiu] mcp server no active sessions, skip tools list_changed notification")
		return
	}

	markedCount := 0
	flushedCount := 0
	changed := d.sessionsWithChangedVisibleSet()
	for _, session := range sessions {
		if changed != nil {
			if _, ok := changed[session.ID]; !ok {
				continue
			}
		}
		marked, flushed, err := markToolsListChangedAndFlush(d.sseHandler, session)
		if !marked {
			continue
		}
		markedCount++
		if err != nil {
			logger.Warnf("[dubbo-go-pixiu] mcp server failed to send tools list_changed: %v", err)
		} else if flushed {
			flushedCount++
		}
	}

	if markedCount == 0 {
		logger.Debugf("[dubbo-go-pixiu] mcp server no active sessions to mark for tools list_changed notification")
		return
	}
	logger.Infof("[dubbo-go-pixiu] mcp server marked tools/list_changed notification for %d sessions, flushed %d online sessions", markedCount, flushedCount)
}

func (d *DynamicConsumer) sessionsWithChangedVisibleSet() router.StringSet {
	d.mu.RLock()
	governance := d.governance
	selector := d.selector
	plans := d.plans
	d.mu.RUnlock()
	if !governance || selector == nil || plans == nil {
		return nil
	}
	refresher, ok := selector.(router.PlanRefreshSelector)
	if !ok {
		return nil
	}
	snapshot := d.registry.toolCatalogSnapshotUnsafe()
	changed := make(router.StringSet)
	for _, item := range plans.SessionPlanContexts() {
		item.Context.CatalogVersion = snapshot.Version
		plan, committed, err := refresher.RefreshPlan(context.Background(), item, snapshot.orderedToolsUnsafe())
		if err != nil {
			logger.Warnf("[dubbo-go-pixiu] mcp server failed to refresh session plan after catalog update: %v", err)
			continue
		}
		if !committed {
			continue
		}
		if item.Plan.VisibleFingerprint != plan.VisibleFingerprint {
			changed[item.Key.SessionID] = struct{}{}
		}
	}
	return changed
}

func (d *DynamicConsumer) governanceEnabled() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.governance
}
