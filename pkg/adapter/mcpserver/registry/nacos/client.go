package nacos

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	nacosmodel "github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

const (
	McpServerVersionGroup          = "mcp-server-versions"
	McpServerSpecGroup             = "mcp-server"
	McpToolSpecGroup               = "mcp-tools"
	SystemConfigIdPrefix           = "system-"
	CredentialPrefix               = "credentials-"
	DefaultNacosListConfigMode     = "blur"
	DefaultNacosListConfigPageSize = 50
	ListMcpServerConfigIdPattern   = "*mcp-versions.json"
)

type ServerSpecInfo struct {
	RemoteServerConfig *RemoteServerConfig `json:"remoteServerConfig"`
}

type RemoteServerConfig struct {
	ServiceRef *ServiceRef `json:"serviceRef"`
}

type ServiceRef struct {
	ServiceName string `json:"serviceName"`
	GroupName   string `json:"groupName"`
	NamespaceId string `json:"namespaceId"`
}

type NacosRegistryClient struct {
	namespaceId  string
	configClient config_client.IConfigClient
	namingClient naming_client.INamingClient
	servers      map[string]*ServerContext
	mu           sync.RWMutex
}

type VersionedMcpServerInfo struct {
	serverInfo *BasicMcpServerInfo
	version    string
}

type ServerContext struct {
	id                     string
	versionedMcpServerInfo *VersionedMcpServerInfo
	serverChangeListener   McpServerListener
	configsMap             map[string]*ConfigListenerWrap
	serviceInfo            *nacosmodel.Service
	namingCallback         func(services []nacosmodel.Instance, err error)
	mu                     sync.RWMutex
}

type McpServerConfig struct {
	ServerSpecConfig string
	ToolsSpecConfig  string
	ServiceInfo      *nacosmodel.Service
	Credentials      map[string]interface{}
	ToolConfigs      []model.ToolConfig // 新增：转换后的工具配置
}

type ConfigListenerWrap struct {
	dataId   string
	group    string
	data     string
	listener func(namespace, group, dataId, data string)
}

type BasicMcpServerInfo struct {
	Name          string `json:"name"`
	Id            string `json:"id"`
	FrontProtocol string `json:"frontProtocol"`
	Protocol      string `json:"protocol"`
}

type VersionsMcpServerInfo struct {
	BasicMcpServerInfo
	LatestPublishedVersion string           `json:"latestPublishedVersion"`
	Versions               []*VersionDetail `json:"versionDetails"`
}

type VersionDetail struct {
	Version  string `json:"version"`
	IsLatest bool   `json:"is_latest"`
}

type McpServerListener func(info *McpServerConfig)

func NewMcpRegistryClient(clientConfig *constant.ClientConfig, serverConfig []constant.ServerConfig, namespaceId string) (*NacosRegistryClient, error) {
	clientParam := vo.NacosClientParam{
		ClientConfig:  clientConfig,
		ServerConfigs: serverConfig,
	}

	configClient, err := clients.NewConfigClient(clientParam)
	if err != nil {
		return nil, fmt.Errorf("failed to create config client: %w", err)
	}

	namingClient, err := clients.NewNamingClient(clientParam)
	if err != nil {
		return nil, fmt.Errorf("failed to create naming client: %w", err)
	}

	return &NacosRegistryClient{
		namespaceId:  namespaceId,
		configClient: configClient,
		namingClient: namingClient,
		servers:      make(map[string]*ServerContext),
	}, nil
}

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
			logger.Errorf("Failed to list mcp server configs for page %d: %v", currentPageNum, err)
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

// ListMcpServer 从 nacos mcp 注册中心列出所有 mcp 服务器
func (n *NacosRegistryClient) ListMcpServer() ([]BasicMcpServerInfo, error) {
	configs, err := n.listMcpServerConfigs()
	if err != nil {
		return nil, err
	}

	result := make([]BasicMcpServerInfo, 0, len(configs))

	for _, config := range configs {
		mcpServerBasicConfig, err := n.configClient.GetConfig(vo.ConfigParam{
			Group:  McpServerVersionGroup,
			DataId: config.DataId,
		})
		if err != nil {
			logger.Errorf("Failed to get mcp server version config (dataId: %s): %v", config.DataId, err)
			continue
		}

		if mcpServerBasicConfig == "" {
			logger.Infof("Empty mcp server version config (dataId: %s)", config.DataId)
			continue
		}

		mcpServer := BasicMcpServerInfo{}
		if err := json.Unmarshal([]byte(mcpServerBasicConfig), &mcpServer); err != nil {
			logger.Errorf("Failed to parse mcp server version config (dataId: %s): %v", config.DataId, err)
			continue
		}

		if !isMcpServerShouldBeDiscoveryForGateway(mcpServer) {
			logger.Debugf("MCP server %s (%s) skipped for gateway discovery", mcpServer.Name, mcpServer.Id)
			continue
		}

		result = append(result, mcpServer)
	}

	return result, nil
}

func isMcpServerShouldBeDiscoveryForGateway(info BasicMcpServerInfo) bool {
	// 支持的协议类型
	supportedProtocols := map[string]bool{
		"mcp-sse":        true,
		"mcp-streamable": true,
	}
	return supportedProtocols[info.FrontProtocol]
}

// ListenToMcpServer 监听 mcp 服务器配置和后端服务
func (n *NacosRegistryClient) ListenToMcpServer(id string, listener McpServerListener) error {
	versionConfigId := fmt.Sprintf("%s-mcp-versions.json", id)
	serverVersionConfig, err := n.configClient.GetConfig(vo.ConfigParam{
		Group:  McpServerVersionGroup,
		DataId: versionConfigId,
	})
	if err != nil {
		logger.Errorf("Get mcp server(id: %s) version config error: %v", id, err)
	} else {
		logger.Infof("Get mcp server(id: %s) version config success", id)
	}

	versionConfigCallBack := func(namespace, group, dataId, content string) {
		logger.Infof("Config callback for mcp server %s", id)

		info := VersionsMcpServerInfo{}
		if err := json.Unmarshal([]byte(content), &info); err != nil {
			logger.Errorf("Parse mcp server (id: %s) version config error: %v", id, err)
			return
		}

		latestVersion := info.LatestPublishedVersion
		ctx := n.servers[id]
		if ctx.versionedMcpServerInfo == nil {
			ctx.versionedMcpServerInfo = &VersionedMcpServerInfo{}
		}
		ctx.versionedMcpServerInfo.serverInfo = &info.BasicMcpServerInfo

		if ctx.versionedMcpServerInfo.version != latestVersion {
			ctx.versionedMcpServerInfo.version = latestVersion
			n.onServerVersionChanged(ctx)
			n.triggerMcpServerChange(id)
		}
	}

	n.servers[id] = &ServerContext{
		id:                   id,
		serverChangeListener: listener,
		configsMap: map[string]*ConfigListenerWrap{
			McpServerVersionGroup: {
				dataId:   versionConfigId,
				group:    McpServerVersionGroup,
				listener: versionConfigCallBack,
			},
		},
	}

	// 手动触发初始回调
	versionConfigCallBack(n.namespaceId, McpServerVersionGroup, versionConfigId, serverVersionConfig)

	// 开始监听配置变更
	err = n.configClient.ListenConfig(vo.ConfigParam{
		Group:    McpServerVersionGroup,
		DataId:   versionConfigId,
		OnChange: versionConfigCallBack,
	})

	return err
}

func (n *NacosRegistryClient) onServerVersionChanged(ctx *ServerContext) {
	id := ctx.versionedMcpServerInfo.serverInfo.Id
	version := ctx.versionedMcpServerInfo.version

	configsMap := map[string]string{
		McpServerSpecGroup: fmt.Sprintf("%s-%s-mcp-server.json", id, version),
		McpToolSpecGroup:   fmt.Sprintf("%s-%s-mcp-tools.json", id, version),
	}

	for group, dataId := range configsMap {
		configsKey := fmt.Sprintf(SystemConfigIdPrefix+"%s@@%s", id, group)

		// 取消旧版本的监听
		if data, exist := ctx.configsMap[configsKey]; exist {
			n.cancelListenToConfig(data)
		}

		configListenerWrap, err := n.ListenToConfig(ctx, dataId, group)
		if err != nil {
			logger.Errorf("Failed to listen to config %s: %v", dataId, err)
			continue
		}

		if configListenerWrap != nil {
			ctx.configsMap[configsKey] = configListenerWrap
		}
	}
}

func (n *NacosRegistryClient) triggerMcpServerChange(id string) {
	context, exist := n.servers[id]
	if !exist || context.serverChangeListener == nil {
		return
	}

	config := mapConfigMapToServerConfig(context)
	if config != nil {
		context.serverChangeListener(config)
	}
}

func mapConfigMapToServerConfig(ctx *ServerContext) *McpServerConfig {
	result := &McpServerConfig{
		Credentials: make(map[string]interface{}),
	}

	for key, data := range ctx.configsMap {
		if data == nil {
			continue
		}

		if strings.HasPrefix(key, SystemConfigIdPrefix) {
			parts := strings.Split(key, "@@")
			if len(parts) == 2 {
				group := parts[1]
				if group == McpServerSpecGroup {
					result.ServerSpecConfig = data.data
				} else if group == McpToolSpecGroup {
					result.ToolsSpecConfig = data.data
				}
			}
		} else if strings.HasPrefix(key, CredentialPrefix) {
			credentialId := strings.TrimPrefix(key, CredentialPrefix)

			var credData interface{}
			if err := json.Unmarshal([]byte(data.data), &credData); err != nil {
				result.Credentials[credentialId] = data.data
			} else {
				result.Credentials[credentialId] = credData
			}
		}
	}

	result.ServiceInfo = ctx.serviceInfo

	// 新增：转换 Nacos 格式为 ToolConfig
	if result.ToolsSpecConfig != "" {
		toolsSpec := &ToolsSpec{}
		if err := json.Unmarshal([]byte(result.ToolsSpecConfig), toolsSpec); err == nil {
			toolConfigs, err := ConvertNacosToolsToToolConfig(toolsSpec)
			if err == nil {
				result.ToolConfigs = toolConfigs
			} else {
				logger.Errorf("Failed to convert tools spec: %v", err)
			}
		}
	}

	return result
}

func (n *NacosRegistryClient) replaceTemplateAndExactConfigsItems(ctx *ServerContext, config *ConfigListenerWrap) map[string]*ConfigListenerWrap {
	result := make(map[string]*ConfigListenerWrap)
	compile := regexp.MustCompile(`\$\{nacos\.([a-zA-Z0-9-_:\\.]+/[a-zA-Z0-9-_:\\.]+)}`)
	allConfigs := compile.FindAllString(config.data, -1)

	newContent := config.data
	for _, configStr := range allConfigs {
		dataIdAndGroup := strings.TrimPrefix(configStr, "${nacos.")
		dataIdAndGroup = strings.TrimSuffix(dataIdAndGroup, "}")

		parts := strings.Split(dataIdAndGroup, "/")
		if len(parts) == 2 {
			dataId := strings.TrimSpace(parts[0])
			group := strings.TrimSpace(parts[1])

			configWrap, err := n.ListenToConfig(ctx, dataId, group)
			if err == nil && configWrap != nil {
				credentialKey := CredentialPrefix + group + "_" + dataId
				result[credentialKey] = configWrap
				newContent = strings.Replace(newContent, configStr, ".config.credentials."+group+"_"+dataId, -1)
			}
		}
	}

	config.data = newContent
	return result
}

func (n *NacosRegistryClient) resetNacosTemplateConfigs(ctx *ServerContext, config *ConfigListenerWrap) {
	newCredentials := n.replaceTemplateAndExactConfigsItems(ctx, config)

	// 取消旧的凭据配置监听
	for key, wrap := range ctx.configsMap {
		if strings.HasPrefix(key, CredentialPrefix) {
			if _, ok := newCredentials[key]; !ok {
				n.cancelListenToConfig(wrap)
				delete(ctx.configsMap, key)
			}
		}
	}

	// 添加新的凭据配置
	for key, data := range newCredentials {
		ctx.configsMap[key] = data
	}
}

func (n *NacosRegistryClient) refreshServiceListenerIfNeeded(ctx *ServerContext, serverConfig string) {
	var serverInfo ServerSpecInfo
	if err := json.Unmarshal([]byte(serverConfig), &serverInfo); err != nil {
		logger.Errorf("Failed to parse server config: %v", err)
		return
	}

	if serverInfo.RemoteServerConfig == nil || serverInfo.RemoteServerConfig.ServiceRef == nil {
		return
	}

	ref := serverInfo.RemoteServerConfig.ServiceRef

	// 取消旧的服务订阅
	if ctx.serviceInfo != nil && ctx.namingCallback != nil {
		n.namingClient.Unsubscribe(&vo.SubscribeParam{
			GroupName:         ctx.serviceInfo.GroupName,
			ServiceName:       ctx.serviceInfo.Name,
			SubscribeCallback: ctx.namingCallback,
		})
	}

	// 获取新的服务信息
	service, err := n.namingClient.GetService(vo.GetServiceParam{
		GroupName:   ref.GroupName,
		ServiceName: ref.ServiceName,
	})
	if err != nil {
		logger.Errorf("Failed to get service (groupName: %s, serviceName: %s): %v",
			ref.GroupName, ref.ServiceName, err)
		return
	}

	ctx.serviceInfo = &service

	// 创建命名回调函数
	ctx.namingCallback = func(services []nacosmodel.Instance, err error) {
		if err == nil && ctx.serviceInfo != nil {
			ctx.serviceInfo.Hosts = services
			n.triggerMcpServerChange(ctx.id)
		}
	}

	// 订阅新的服务
	n.namingClient.Subscribe(&vo.SubscribeParam{
		GroupName:         ctx.serviceInfo.GroupName,
		ServiceName:       ctx.serviceInfo.Name,
		SubscribeCallback: ctx.namingCallback,
	})
}

func (n *NacosRegistryClient) ListenToConfig(ctx *ServerContext, dataId string, group string) (*ConfigListenerWrap, error) {
	wrap := &ConfigListenerWrap{
		dataId: dataId,
		group:  group,
	}

	configListener := func(namespace, group, dataId, data string) {
		if ctx.serverChangeListener != nil && wrap.data != data {
			wrap.data = data

			if group == McpToolSpecGroup {
				n.resetNacosTemplateConfigs(ctx, wrap)
			} else if group == McpServerSpecGroup {
				n.refreshServiceListenerIfNeeded(ctx, data)
			}

			n.triggerMcpServerChange(ctx.versionedMcpServerInfo.serverInfo.Id)
		}
	}

	config, err := n.configClient.GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
	if err != nil {
		return nil, err
	}

	wrap.listener = configListener
	wrap.data = config

	if group == McpToolSpecGroup {
		n.resetNacosTemplateConfigs(ctx, wrap)
	} else if group == McpServerSpecGroup {
		n.refreshServiceListenerIfNeeded(ctx, config)
	}

	err = n.configClient.ListenConfig(vo.ConfigParam{
		DataId:   dataId,
		Group:    group,
		OnChange: configListener,
	})
	if err != nil {
		return nil, err
	}

	return wrap, nil
}

func (n *NacosRegistryClient) cancelListenToConfig(wrap *ConfigListenerWrap) error {
	if wrap == nil {
		return nil
	}

	return n.configClient.CancelListenConfig(vo.ConfigParam{
		DataId:   wrap.dataId,
		Group:    wrap.group,
		OnChange: wrap.listener,
	})
}

func (n *NacosRegistryClient) CancelListenToServer(id string) error {
	server, exist := n.servers[id]
	if !exist || server == nil {
		return nil
	}

	delete(n.servers, id)

	// 取消所有配置监听
	for _, wrap := range server.configsMap {
		if wrap != nil {
			n.configClient.CancelListenConfig(vo.ConfigParam{
				DataId:   wrap.dataId,
				Group:    wrap.group,
				OnChange: wrap.listener,
			})
		}
	}

	// 取消服务订阅
	if server.serviceInfo != nil && server.namingCallback != nil {
		n.namingClient.Unsubscribe(&vo.SubscribeParam{
			GroupName:         server.serviceInfo.GroupName,
			ServiceName:       server.serviceInfo.Name,
			SubscribeCallback: server.namingCallback,
		})
	}

	return nil
}

func (n *NacosRegistryClient) CloseClient() {
	// 取消所有服务器监听
	for id := range n.servers {
		n.CancelListenToServer(id)
	}

	// 关闭 Nacos 客户端
	if n.namingClient != nil {
		n.namingClient.CloseClient()
	}
	if n.configClient != nil {
		n.configClient.CloseClient()
	}
}
