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

package nacos

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync"
)

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	nacosmodel "github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// Nacos configuration group constants
const (
	McpServerVersionGroup          = "mcp-server-versions"
	McpServerSpecGroup             = "mcp-server"
	McpToolSpecGroup               = "mcp-tools"
	DefaultNacosListConfigMode     = "blur"
	DefaultNacosListConfigPageSize = 50
	ListMcpServerConfigIdPattern   = "*mcp-versions.json"
)

// Configuration key name constants
const (
	SystemConfigKeyPrefix     = "system::"
	CredentialConfigKeyPrefix = "cred::"
	ConfigKeySeparator        = "::"
)

// Supported MCP protocol types
var supportedMcpProtocols = map[string]bool{
	"mcp-sse":        false,
	"mcp-streamable": true,
}

var templateRegex = regexp.MustCompile(`\$\{nacos\.([a-zA-Z0-9-_:\\.]+/[a-zA-Z0-9-_:\\.]+)}`)

// ==================== Data Structure Definitions ====================

// ServerSpecInfo server specification information
type ServerSpecInfo struct {
	RemoteServerConfig *RemoteServerConfig `json:"remoteServerConfig"`
}

// RemoteServerConfig remote server configuration
type RemoteServerConfig struct {
	ServiceRef *ServiceRef `json:"serviceRef"`
}

// ServiceRef service reference information
type ServiceRef struct {
	ServiceName string `json:"serviceName"`
	GroupName   string `json:"groupName"`
	NamespaceId string `json:"namespaceId"`
}

// NacosRegistryClient Nacos registry client
// Responsible for communicating with Nacos configuration center and service discovery
type NacosRegistryClient struct {
	mu           sync.RWMutex // Read-write lock, protects concurrent access to servers map
	namespaceId  string
	configClient config_client.IConfigClient
	namingClient naming_client.INamingClient
	servers      map[string]*ServerContext
}

// VersionedMcpServerInfo MCP server information with version
type VersionedMcpServerInfo struct {
	serverInfo *BasicMcpServerInfo
	version    string
}

// ServerContext server context, manages all states of a single MCP server
type ServerContext struct {
	mu                     sync.Mutex
	id                     string
	versionedMcpServerInfo *VersionedMcpServerInfo
	serverChangeListener   McpServerListener
	configsMap             map[string]*ConfigListenerWrap
	serviceInfo            *nacosmodel.Service
	namingCallback         func(services []nacosmodel.Instance, err error)
}

// McpServerConfig MCP server configuration
type McpServerConfig struct {
	ServerSpecConfig string
	ToolsSpecConfig  string
	ServiceInfo      *nacosmodel.Service
	Credentials      map[string]any
	ToolConfigs      []model.ToolConfig // Converted tool configurations
}

// ConfigListenerWrap configuration listener wrapper
type ConfigListenerWrap struct {
	dataId   string
	group    string
	data     string
	listener func(namespace, group, dataId, data string)
}

// BasicMcpServerInfo basic MCP server information
type BasicMcpServerInfo struct {
	Name          string `json:"name"`
	Id            string `json:"id"`
	FrontProtocol string `json:"frontProtocol"`
	Protocol      string `json:"protocol"`
}

// VersionsMcpServerInfo MCP server information with version details
type VersionsMcpServerInfo struct {
	BasicMcpServerInfo
	LatestPublishedVersion string           `json:"latestPublishedVersion"`
	Versions               []*VersionDetail `json:"versionDetails"`
}

// VersionDetail version detail
type VersionDetail struct {
	Version  string `json:"version"`
	IsLatest bool   `json:"is_latest"`
}

// McpServerListener MCP server change listener
type McpServerListener func(info *McpServerConfig)

// ==================== Constructors ====================

// NewMcpRegistryClient creates a new Nacos MCP registry client
func NewMcpRegistryClient(clientConfig *constant.ClientConfig, serverConfig []constant.ServerConfig, namespaceId string) (*NacosRegistryClient, error) {
	logger.Infof("[dubbo-go-pixiu] nacos registry creating MCP registry client for namespace: %s", namespaceId)

	clientParam := vo.NacosClientParam{
		ClientConfig:  clientConfig,
		ServerConfigs: serverConfig,
	}

	configClient, err := clients.NewConfigClient(clientParam)
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] nacos registry failed to create config client: %v", err)
		return nil, fmt.Errorf("failed to create config client: %w", err)
	}

	namingClient, err := clients.NewNamingClient(clientParam)
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] nacos registry failed to create naming client: %v", err)
		return nil, fmt.Errorf("failed to create naming client: %w", err)
	}

	client := &NacosRegistryClient{
		namespaceId:  namespaceId,
		configClient: configClient,
		namingClient: namingClient,
		servers:      make(map[string]*ServerContext),
	}

	return client, nil
}

// ==================== Public API ====================

// ListMcpServer lists all MCP servers from Nacos MCP registry
func (n *NacosRegistryClient) ListMcpServer() ([]BasicMcpServerInfo, error) {

	configs, err := n.listMcpServerConfigs()
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] nacos registry failed to list MCP server configs: %v", err)
		return nil, err
	}

	result := make([]BasicMcpServerInfo, 0, len(configs))
	skippedCount := 0

	for _, config := range configs {

		mcpServerBasicConfig, err := n.configClient.GetConfig(vo.ConfigParam{
			Group:  McpServerVersionGroup,
			DataId: config.DataId,
		})
		if err != nil {
			logger.Errorf("[dubbo-go-pixiu] nacos registry failed to get MCP server version config (dataId: %s): %v", config.DataId, err)
			continue
		}

		if mcpServerBasicConfig == "" {
			skippedCount++
			continue
		}

		mcpServer := BasicMcpServerInfo{}
		if err := json.Unmarshal([]byte(mcpServerBasicConfig), &mcpServer); err != nil {
			logger.Errorf("[dubbo-go-pixiu] nacos registry failed to parse MCP server version config (dataId: %s): %v", config.DataId, err)
			skippedCount++
			continue
		}

		if !isMcpServerShouldBeDiscoveryForGateway(mcpServer) {
			skippedCount++
			continue
		}

		result = append(result, mcpServer)
	}

	logger.Debugf("[dubbo-go-pixiu] nacos registry successfully listed %d MCP servers, skipped %d invalid/unsupported servers",
		len(result), skippedCount)
	return result, nil
}

// ListenToMcpServer listens to MCP server configuration and backend services
func (n *NacosRegistryClient) ListenToMcpServer(id string, listener McpServerListener) error {
	logger.Infof("[dubbo-go-pixiu] nacos registry starting to listen to MCP server: %s", id)

	versionConfigId := fmt.Sprintf("%s-mcp-versions.json", id)

	// Get initial version configuration (network call, may block)
	serverVersionConfig, err := n.configClient.GetConfig(vo.ConfigParam{
		Group:  McpServerVersionGroup,
		DataId: versionConfigId,
	})
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] nacos registry failed to get initial version config for MCP server %s: %v", id, err)
	}

	// Version callback implementation: parse, update ctx.version and handle version changes asynchronously
	versionConfigCallBack := func(namespace, group, dataId, content string) {
		var info VersionsMcpServerInfo
		if err := json.Unmarshal([]byte(content), &info); err != nil {
			logger.Errorf("[dubbo-go-pixiu] nacos registry failed to parse version config for MCP server %s: %v", id, err)
			return
		}

		// Get ctx pointer (short-term global lock)
		n.mu.RLock()
		ctx := n.servers[id]
		n.mu.RUnlock()
		if ctx == nil {
			logger.Errorf("[dubbo-go-pixiu] nacos registry server context not found for MCP server %s", id)
			return
		}

		// Update ctx.versionedMcpServerInfo (local lock)
		ctx.mu.Lock()
		if ctx.versionedMcpServerInfo == nil {
			ctx.versionedMcpServerInfo = &VersionedMcpServerInfo{}
		}
		changed := ctx.versionedMcpServerInfo.version != info.LatestPublishedVersion
		ctx.versionedMcpServerInfo.serverInfo = &info.BasicMcpServerInfo
		if changed {
			ctx.versionedMcpServerInfo.version = info.LatestPublishedVersion
		}
		ctx.mu.Unlock()

		if changed {
			// Asynchronously execute time-consuming/network operations to avoid blocking SDK callbacks
			go func() {
				n.onServerVersionChanged(ctx)
				n.triggerMcpServerChange(id)
			}()
		}
	}

	// Create server context (short-term global lock)
	ctx := &ServerContext{
		id:                   id,
		serverChangeListener: listener,
		configsMap:           make(map[string]*ConfigListenerWrap),
	}
	// Use unified system key to store version configuration
	ctx.configsMap[buildSystemConfigKey(id, McpServerVersionGroup)] = &ConfigListenerWrap{
		dataId:   versionConfigId,
		group:    McpServerVersionGroup,
		listener: versionConfigCallBack,
	}

	n.mu.Lock()
	n.servers[id] = ctx
	n.mu.Unlock()

	// Manually trigger initial callback (callback will self-lock and asynchronously handle onServerVersionChanged)
	versionConfigCallBack(n.namespaceId, McpServerVersionGroup, versionConfigId, serverVersionConfig)

	// Start listening to configuration changes (network call)
	err = n.configClient.ListenConfig(vo.ConfigParam{
		Group:    McpServerVersionGroup,
		DataId:   versionConfigId,
		OnChange: versionConfigCallBack,
	})

	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] nacos registry failed to start listening to MCP server %s: %v", id, err)
		// If ListenConfig fails, clean up the ctx that was just added
		n.mu.Lock()
		delete(n.servers, id)
		n.mu.Unlock()
		return err
	}

	return nil
}

// CancelListenToServer cancels listening to MCP server
func (n *NacosRegistryClient) CancelListenToServer(id string) error {
	logger.Infof("[dubbo-go-pixiu] nacos registry canceling listen to MCP server: %s", id)

	// First briefly take out and delete ctx to avoid holding lock for long time
	n.mu.Lock()
	server, exist := n.servers[id]
	if !exist || server == nil {
		n.mu.Unlock()
		logger.Warnf("[dubbo-go-pixiu] nacos registry server context not found for MCP server %s", id)
		return nil
	}
	delete(n.servers, id)
	n.mu.Unlock()

	// Cancel all configuration listeners (not within global lock)
	// collect wraps under ctx.mu to avoid concurrent map mutation
	server.mu.Lock()
	wraps := make([]*ConfigListenerWrap, 0, len(server.configsMap))
	for _, wrap := range server.configsMap {
		if wrap != nil {
			wraps = append(wraps, wrap)
		}
	}
	server.configsMap = make(map[string]*ConfigListenerWrap) // clear map
	oldService := server.serviceInfo
	oldCallback := server.namingCallback
	server.serviceInfo = nil
	server.namingCallback = nil
	server.mu.Unlock()

	for _, wrap := range wraps {
		n.cancelListenToConfig(wrap)
	}

	// Cancel service subscription (network operation)
	if oldService != nil && oldCallback != nil {
		n.namingClient.Unsubscribe(&vo.SubscribeParam{
			GroupName:         oldService.GroupName,
			ServiceName:       oldService.Name,
			SubscribeCallback: oldCallback,
		})
	}

	return nil
}

// CloseClient closes the client and cleans up all resources
func (n *NacosRegistryClient) CloseClient() {
	logger.Infof("[dubbo-go-pixiu] nacos registry closing MCP registry client")

	// First collect all server IDs and release lock
	n.mu.RLock()
	ids := make([]string, 0, len(n.servers))
	for id := range n.servers {
		ids = append(ids, id)
	}
	n.mu.RUnlock()

	// Now safely cancel all server listeners
	for _, id := range ids {
		_ = n.CancelListenToServer(id)
	}

	// Close Nacos clients
	if n.namingClient != nil {
		n.namingClient.CloseClient()
	}
	if n.configClient != nil {
		n.configClient.CloseClient()
	}
}

// ==================== Internal Helper Methods ====================

// listMcpServerConfigs paginates to get all MCP server configurations
func (n *NacosRegistryClient) listMcpServerConfigs() ([]nacosmodel.ConfigItem, error) {

	currentPageNum := 1
	result := make([]nacosmodel.ConfigItem, 0)

	for currentPageNum <= 100 {

		configPage, err := n.configClient.SearchConfig(vo.SearchConfigParam{
			Search:   DefaultNacosListConfigMode,
			DataId:   ListMcpServerConfigIdPattern,
			Group:    McpServerVersionGroup,
			PageNo:   currentPageNum,
			PageSize: DefaultNacosListConfigPageSize,
		})

		if err != nil {
			logger.Errorf("[dubbo-go-pixiu] nacos registry failed to list MCP server configs for page %d: %v", currentPageNum, err)
			return result, err
		}

		if configPage == nil {
			break
		}

		result = append(result, configPage.PageItems...)

		if configPage.PageNumber >= configPage.PagesAvailable {
			break
		}

		currentPageNum++
	}

	return result, nil
}

// onServerVersionChanged handles server version changes
func (n *NacosRegistryClient) onServerVersionChanged(ctx *ServerContext) {
	// Read necessary information
	ctx.mu.Lock()
	if ctx.versionedMcpServerInfo == nil || ctx.versionedMcpServerInfo.serverInfo == nil {
		ctx.mu.Unlock()
		logger.Errorf("[dubbo-go-pixiu] nacos registry missing version/serverInfo for ctx")
		return
	}
	id := ctx.versionedMcpServerInfo.serverInfo.Id
	version := ctx.versionedMcpServerInfo.version
	ctx.mu.Unlock()

	logger.Infof("[dubbo-go-pixiu] nacos registry processing version change for MCP server %s to version %s", id, version)

	configsMap := map[string]string{
		McpServerSpecGroup: fmt.Sprintf("%s-%s-mcp-server.json", id, version),
		McpToolSpecGroup:   fmt.Sprintf("%s-%s-mcp-tools.json", id, version),
	}

	// For each group: cancel old wrap (if any), create new Listener (network) and update ctx.configsMap under ctx.mu
	for group, dataId := range configsMap {
		configsKey := buildSystemConfigKey(id, group)

		// If there is an old wrap, remove reference and cancel after unlocking
		var oldWrap *ConfigListenerWrap
		ctx.mu.Lock()
		if w, ok := ctx.configsMap[configsKey]; ok && w != nil {
			oldWrap = w
			delete(ctx.configsMap, configsKey)
		}
		ctx.mu.Unlock()

		if oldWrap != nil {
			n.cancelListenToConfig(oldWrap)
		}

		// Create new version listener (network, may block)
		newWrap, err := n.ListenToConfig(ctx, dataId, group)
		if err != nil {
			logger.Errorf("[dubbo-go-pixiu] nacos registry failed to listen to config %s for MCP server %s: %v", dataId, id, err)
			continue
		}

		// Save new wrap under ctx.mu
		ctx.mu.Lock()
		ctx.configsMap[configsKey] = newWrap
		ctx.mu.Unlock()
	}
}

// triggerMcpServerChange triggers MCP server change notification
func (n *NacosRegistryClient) triggerMcpServerChange(id string) {
	// Get ctx pointer
	n.mu.RLock()
	ctx := n.servers[id]
	n.mu.RUnlock()
	if ctx == nil {
		logger.Warnf("[dubbo-go-pixiu] nacos registry server context not found for MCP server %s when triggering change", id)
		return
	}

	// Generate configuration using ctx.mu protection
	cfg := mapConfigMapToServerConfig(ctx) // Internally acquires ctx.mu
	if cfg != nil {
		// Call user callback after releasing lock
		ctx.mu.Lock()
		listener := ctx.serverChangeListener
		ctx.mu.Unlock()

		if listener != nil {
			listener(cfg)
		}
	} else {
		logger.Warnf("[dubbo-go-pixiu] nacos registry failed to generate config for MCP server %s", id)
	}
}

// ListenToConfig listens to configuration changes
func (n *NacosRegistryClient) ListenToConfig(ctx *ServerContext, dataId string, group string) (*ConfigListenerWrap, error) {

	wrap := &ConfigListenerWrap{
		dataId: dataId,
		group:  group,
	}

	// Create configListener (only use ctx.mu for local state updates, put time-consuming operations outside lock)
	configListener := func(namespace, group, dataId, data string) {
		// Quickly atomically update wrap.data (within ctx.mu)
		ctx.mu.Lock()
		changed := wrap.data != data
		if changed {
			wrap.data = data
		}
		ctx.mu.Unlock()

		if !changed {
			return
		}

		// Execute time-consuming or network-related logic outside lock
		switch group {
		case McpToolSpecGroup:
			n.resetNacosTemplateConfigs(ctx, wrap)
		case McpServerSpecGroup:
			n.refreshServiceListenerIfNeeded(ctx, data)
		}

		n.triggerMcpServerChange(ctx.id)
	}

	// Get initial configuration (network call)
	config, err := n.configClient.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] nacos registry failed to get initial config for MCP server %s (dataId: %s, group: %s): %v",
			ctx.id, dataId, group, err)
		return nil, err
	}

	wrap.listener = configListener
	wrap.data = config

	// Process initial configuration (placed outside lock, as these operations may perform network calls)
	switch group {
	case McpToolSpecGroup:
		n.resetNacosTemplateConfigs(ctx, wrap)
	case McpServerSpecGroup:
		n.refreshServiceListenerIfNeeded(ctx, config)
	}

	// Start listening to configuration changes (network call)
	err = n.configClient.ListenConfig(vo.ConfigParam{
		DataId:   dataId,
		Group:    group,
		OnChange: configListener,
	})
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] nacos registry failed to start listening to config for MCP server %s (dataId: %s, group: %s): %v",
			ctx.id, dataId, group, err)
		return nil, err
	}

	return wrap, nil
}

// cancelListenToConfig cancels configuration listening
func (n *NacosRegistryClient) cancelListenToConfig(wrap *ConfigListenerWrap) error {
	if wrap == nil {
		return nil
	}

	err := n.configClient.CancelListenConfig(vo.ConfigParam{
		DataId:   wrap.dataId,
		Group:    wrap.group,
		OnChange: wrap.listener,
	})

	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] nacos registry failed to cancel config listener (dataId: %s, group: %s): %v",
			wrap.dataId, wrap.group, err)
	}

	return err
}

// ==================== Configuration Conversion ====================

// mapConfigMapToServerConfig converts configuration map to MCP server configuration
// Internally acquires ctx.mu to protect reading ctx fields
func mapConfigMapToServerConfig(ctx *ServerContext) *McpServerConfig {

	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	result := &McpServerConfig{
		Credentials: make(map[string]any),
	}

	// Process system configurations
	for key, data := range ctx.configsMap {
		if data == nil {
			continue
		}

		if isSystemConfigKey(key) {
			_, _, group, _ := parseConfigKey(key)
			switch group {
			case McpServerSpecGroup:
				result.ServerSpecConfig = data.data
			case McpToolSpecGroup:
				result.ToolsSpecConfig = data.data
			}
		} else if isCredentialConfigKey(key) {
			_, _, group, dataId := parseConfigKey(key)
			credentialId := group + "_" + dataId

			var credData any
			if err := json.Unmarshal([]byte(data.data), &credData); err != nil {
				result.Credentials[credentialId] = data.data
			} else {
				result.Credentials[credentialId] = credData
			}
		}
	}

	// Directly reference serviceInfo (caller must note that this object may be modified)
	result.ServiceInfo = ctx.serviceInfo

	// Convert Nacos format to ToolConfig
	if result.ToolsSpecConfig != "" {
		toolsSpec := &ToolsSpec{}
		if err := json.Unmarshal([]byte(result.ToolsSpecConfig), toolsSpec); err == nil {
			toolConfigs, err := ConvertNacosToolsToToolConfig(toolsSpec)
			if err == nil {
				result.ToolConfigs = toolConfigs
			} else {
				logger.Errorf("[dubbo-go-pixiu] nacos registry failed to convert tools spec for MCP server %s: %v", ctx.id, err)
			}
		} else {
			logger.Errorf("[dubbo-go-pixiu] nacos registry failed to parse tools spec JSON for MCP server %s: %v", ctx.id, err)
		}
	}

	return result
}

// parseTemplatePlaceholders parses template placeholders and returns replaced content and required (dataId,group) list
func parseTemplatePlaceholders(content string) (string, [][2]string) {
	all := templateRegex.FindAllStringSubmatch(content, -1)
	if len(all) == 0 {
		return content, nil
	}
	placeholders := make([][2]string, 0, len(all))
	newContent := content
	seen := make(map[string]bool)
	for _, m := range all {
		if len(m) < 2 {
			continue
		}
		p := m[1] // dataId/group
		parts := strings.Split(p, "/")
		if len(parts) != 2 {
			continue
		}
		dataId := strings.TrimSpace(parts[0])
		group := strings.TrimSpace(parts[1])
		key := group + "::" + dataId
		if seen[key] {
			// replace duplicates as well
			newContent = strings.ReplaceAll(newContent, "${nacos."+p+"}", ".config.credentials."+group+"_"+dataId)
			continue
		}
		seen[key] = true
		placeholders = append(placeholders, [2]string{dataId, group})
		newContent = strings.ReplaceAll(newContent, "${nacos."+p+"}", ".config.credentials."+group+"_"+dataId)
	}
	return newContent, placeholders
}

// resetNacosTemplateConfigs resets template configurations and their referenced credential configurations
func (n *NacosRegistryClient) resetNacosTemplateConfigs(ctx *ServerContext, config *ConfigListenerWrap) {
	// Parse out placeholders and replaced text
	newContent, placeholders := parseTemplatePlaceholders(config.data)

	// For each placeholder, create ListenToConfig (network call), collect newWraps
	newWraps := make(map[string]*ConfigListenerWrap)
	for _, ph := range placeholders {
		dataId := ph[0]
		group := ph[1]
		wrap, err := n.ListenToConfig(ctx, dataId, group)
		if err != nil {
			logger.Warnf("[dubbo-go-pixiu] nacos registry failed to listen to credential config %s/%s: %v", dataId, group, err)
			continue
		}
		key := buildCredentialConfigKey(group, dataId)
		newWraps[key] = wrap
	}

	// Update ctx.configsMap: cancel old ones, add new ones. Execute deletion and insertion within ctx.mu (cancel network operations outside)
	// First find wraps that need to be canceled
	ctx.mu.Lock()
	toRemove := make([]*ConfigListenerWrap, 0)
	for key, wrap := range ctx.configsMap {
		if isCredentialConfigKey(key) {
			if _, ok := newWraps[key]; !ok {
				toRemove = append(toRemove, wrap)
				delete(ctx.configsMap, key)
			}
		}
	}
	// Add new wrap references
	for key, wrap := range newWraps {
		ctx.configsMap[key] = wrap
	}
	// Update original config data content
	config.data = newContent
	ctx.mu.Unlock()

	// Cancel old listeners (network call)
	for _, w := range toRemove {
		n.cancelListenToConfig(w)
	}
}

// refreshServiceListenerIfNeeded refreshes service listener as needed
func (n *NacosRegistryClient) refreshServiceListenerIfNeeded(ctx *ServerContext, serverConfig string) {
	var serverInfo ServerSpecInfo
	if err := json.Unmarshal([]byte(serverConfig), &serverInfo); err != nil {
		logger.Errorf("[dubbo-go-pixiu] nacos registry failed to parse server config for MCP server %s: %v", ctx.id, err)
		return
	}

	if serverInfo.RemoteServerConfig == nil || serverInfo.RemoteServerConfig.ServiceRef == nil {
		return
	}
	ref := serverInfo.RemoteServerConfig.ServiceRef

	// Get old subscription information and cancel externally (read references within ctx.mu)
	ctx.mu.Lock()
	oldService := ctx.serviceInfo
	oldCallback := ctx.namingCallback
	ctx.mu.Unlock()

	if oldService != nil && oldCallback != nil {
		// Cancel old subscription (network)
		n.namingClient.Unsubscribe(&vo.SubscribeParam{
			GroupName:         oldService.GroupName,
			ServiceName:       oldService.Name,
			SubscribeCallback: oldCallback,
		})
	}

	// Get new service information (network)
	service, err := n.namingClient.GetService(vo.GetServiceParam{
		GroupName:   ref.GroupName,
		ServiceName: ref.ServiceName,
	})
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] nacos registry failed to get service for MCP server %s (groupName: %s, serviceName: %s): %v",
			ctx.id, ref.GroupName, ref.ServiceName, err)
		return
	}

	// Create naming callback: only use ctx.mu to update local serviceInfo.Hosts, then trigger change (callback does not hold global lock)
	namingCb := func(services []nacosmodel.Instance, err error) {
		if err != nil {
			logger.Errorf("[dubbo-go-pixiu] nacos registry service callback error for MCP server %s: %v", ctx.id, err)
			return
		}
		ctx.mu.Lock()
		if ctx.serviceInfo != nil {
			ctx.serviceInfo.Hosts = services
		}
		ctx.mu.Unlock()
		n.triggerMcpServerChange(ctx.id)
	}

	// Update serviceInfo and namingCallback in ctx (within local lock)
	ctx.mu.Lock()
	ctx.serviceInfo = &service
	ctx.namingCallback = namingCb
	ctx.mu.Unlock()

	// Subscribe to new service (network)
	if err := n.namingClient.Subscribe(&vo.SubscribeParam{
		GroupName:         ctx.serviceInfo.GroupName,
		ServiceName:       ctx.serviceInfo.Name,
		SubscribeCallback: ctx.namingCallback,
	}); err != nil {
		logger.Errorf("[dubbo-go-pixiu] nacos registry failed to subscribe to service for MCP server %s: %v", ctx.id, err)
	}
}

// ==================== Configuration Key Helper Functions ====================

// buildSystemConfigKey constructs system configuration key name
// Format: system::<id>::<group>
func buildSystemConfigKey(id, group string) string {
	return SystemConfigKeyPrefix + id + ConfigKeySeparator + group
}

// buildCredentialConfigKey constructs credential configuration key name
// Format: cred::<group>::<dataId>
func buildCredentialConfigKey(group, dataId string) string {
	return CredentialConfigKeyPrefix + group + ConfigKeySeparator + dataId
}

// parseConfigKey parses configuration key name
// Returns: keyType, id, group, dataId
func parseConfigKey(key string) (keyType, id, group, dataId string) {
	if strings.HasPrefix(key, SystemConfigKeyPrefix) {
		keyType = "system"
		parts := strings.Split(strings.TrimPrefix(key, SystemConfigKeyPrefix), ConfigKeySeparator)
		if len(parts) == 2 {
			id = parts[0]
			group = parts[1]
		}
	} else if strings.HasPrefix(key, CredentialConfigKeyPrefix) {
		keyType = "credential"
		parts := strings.Split(strings.TrimPrefix(key, CredentialConfigKeyPrefix), ConfigKeySeparator)
		if len(parts) == 2 {
			group = parts[0]
			dataId = parts[1]
		}
	}
	return
}

// isSystemConfigKey determines if it is a system configuration key
func isSystemConfigKey(key string) bool {
	return strings.HasPrefix(key, SystemConfigKeyPrefix)
}

// isCredentialConfigKey determines if it is a credential configuration key
func isCredentialConfigKey(key string) bool {
	return strings.HasPrefix(key, CredentialConfigKeyPrefix)
}

// ==================== Utility Functions ====================

// isMcpServerShouldBeDiscoveryForGateway checks if MCP server should be discovered by gateway
func isMcpServerShouldBeDiscoveryForGateway(info BasicMcpServerInfo) bool {
	v, ok := supportedMcpProtocols[info.FrontProtocol]
	return ok && v
}
