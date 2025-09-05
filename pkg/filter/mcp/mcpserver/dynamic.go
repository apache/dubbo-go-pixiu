package mcpserver

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

var (
	globalRegistry *ToolRegistry
	globalDynamic  *DynamicConsumer
)

// GetOrInitRegistry returns a singleton ToolRegistry
func GetOrInitRegistry() *ToolRegistry {
	if globalRegistry == nil {
		globalRegistry = NewToolRegistry()
	}
	return globalRegistry
}

// GetOrInitDynamic returns a singleton DynamicConsumer
func GetOrInitDynamic() *DynamicConsumer {
	if globalDynamic == nil {
		globalDynamic = NewDynamicConsumer(GetOrInitRegistry())
	}
	return globalDynamic
}

// DynamicConsumer applies dynamic MCP configurations into the registry
type DynamicConsumer struct {
	registry *ToolRegistry
}

func NewDynamicConsumer(reg *ToolRegistry) *DynamicConsumer {
	return &DynamicConsumer{registry: reg}
}

// update tools from the remote config in nacos
func (d *DynamicConsumer) ApplyMcpServerConfig(cfg *model.McpServerConfig) error {
	if cfg == nil {
		return nil
	}

	// full sync tools
	d.registry.ReplaceAllTools(cfg.Tools)
	logger.Infof("[MCP Dynamic] applied config: tools synced=%d", len(cfg.Tools))
	return nil
}
