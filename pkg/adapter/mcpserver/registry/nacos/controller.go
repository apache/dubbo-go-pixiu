package nacos

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

// McpController 是 MCP Server 在 Nacos 中的配置同步器。
// 它负责发现、监听、转换和应用配置。
type McpController struct {
	client   *NacosRegistryClient
	onChange func(cfg *McpServerConfig)
	watched  map[string]bool
	mu       sync.RWMutex
}

// NewMcpController 创建新的 MCP 控制器
func NewMcpController(client *NacosRegistryClient, onChange func(cfg *McpServerConfig)) *McpController {
	return &McpController{
		client:   client,
		onChange: onChange,
		watched:  make(map[string]bool),
	}
}

// Run 启动控制器，定期发现和监听 MCP 服务
func (c *McpController) Run(ctx context.Context, interval time.Duration) error {
	logger.Infof("Starting MCP controller with interval: %v", interval)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// 立即执行一次
	if err := c.reconcile(); err != nil {
		logger.Errorf("Initial reconcile failed: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			logger.Infof("MCP controller stopped")
			return nil
		case <-ticker.C:
			if err := c.reconcile(); err != nil {
				logger.Errorf("Reconcile failed: %v", err)
			}
		}
	}
}

// reconcile 执行协调逻辑：发现服务、计算差异、绑定监听
func (c *McpController) reconcile() error {
	// 获取所有 MCP 服务
	servers, err := c.client.ListMcpServer()
	if err != nil {
		return fmt.Errorf("failed to list MCP servers: %w", err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// 计算需要监听的服务
	currentWatched := make(map[string]bool)
	for _, server := range servers {
		currentWatched[server.Id] = true

		// 如果还没有监听，则开始监听
		if !c.watched[server.Id] {
			logger.Infof("Starting to watch MCP server: %s (%s)", server.Name, server.Id)

			err := c.client.ListenToMcpServer(server.Id, c.wrapListener(server.Id))
			if err != nil {
				logger.Errorf("Failed to listen to server %s: %v", server.Id, err)
				continue
			}

			c.watched[server.Id] = true
		}
	}

	// 取消不再存在的服务监听
	for serverId := range c.watched {
		if !currentWatched[serverId] {
			logger.Infof("Stopping watch for MCP server: %s", serverId)

			err := c.client.CancelListenToServer(serverId)
			if err != nil {
				logger.Errorf("Failed to cancel listen for server %s: %v", serverId, err)
			}

			delete(c.watched, serverId)
		}
	}

	return nil
}

// wrapListener 包装监听器回调
func (c *McpController) wrapListener(serverId string) McpServerListener {
	return func(cfg *McpServerConfig) {
		logger.Infof("Received config update for server: %s", serverId)

		if c.onChange != nil {
			c.onChange(cfg)
		}
	}
}

// Close 关闭控制器
func (c *McpController) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 取消所有监听
	for serverId := range c.watched {
		err := c.client.CancelListenToServer(serverId)
		if err != nil {
			logger.Errorf("Failed to cancel listen for server %s: %v", serverId, err)
		}
	}

	// 清空监听列表
	c.watched = make(map[string]bool)

	return nil
}
