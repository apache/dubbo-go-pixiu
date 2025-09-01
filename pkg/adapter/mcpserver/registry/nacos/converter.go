package nacos

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// ConvertNacosToolsToToolConfig 转换 Nacos 工具格式为 Filter 的 ToolConfig
func ConvertNacosToolsToToolConfig(toolsSpec *ToolsSpec) ([]model.ToolConfig, error) {
	var toolConfigs []model.ToolConfig

	for _, nacosTool := range toolsSpec.Tools {
		meta := toolsSpec.ToolsMeta[nacosTool.Name]
		if !meta.Enabled {
			continue
		}

		// 提取 json-go-template
		templateData, ok := meta.Templates["json-go-template"]
		if !ok {
			logger.Warnf("Tool %s has no json-go-template, skipping", nacosTool.Name)
			continue
		}

		toolConfig, err := convertSingleTool(nacosTool, templateData)
		if err != nil {
			return nil, fmt.Errorf("failed to convert tool %s: %w", nacosTool.Name, err)
		}

		toolConfigs = append(toolConfigs, toolConfig)
	}

	return toolConfigs, nil
}

func convertSingleTool(nacosTool NacosTool, templateData interface{}) (model.ToolConfig, error) {
	// 解析模板数据
	templateBytes, err := json.Marshal(templateData)
	if err != nil {
		return model.ToolConfig{}, err
	}

	var template JsonGoTemplate
	if err := json.Unmarshal(templateBytes, &template); err != nil {
		return model.ToolConfig{}, err
	}

	toolConfig := model.ToolConfig{
		Name:        nacosTool.Name,
		Description: nacosTool.Description,
		Cluster:     extractClusterFromURL(template.RequestTemplate.URL),
		Request: model.RequestConfig{
			Method:  template.RequestTemplate.Method,
			Path:    extractPathFromURL(template.RequestTemplate.URL),
			Headers: convertHeaders(template.RequestTemplate.Headers),
		},
		Args: func() []model.ArgConfig {
			args, err := convertInputSchemaToArgs(nacosTool.InputSchema, template.RequestTemplate)
			if err != nil {
				logger.Warnf("Failed to convert args for tool %s: %v", nacosTool.Name, err)
				return []model.ArgConfig{}
			}
			return args
		}(),
	}

	return toolConfig, nil
}

func extractClusterFromURL(url string) string {
	// 从 URL 中提取集群名
	if strings.HasPrefix(url, "http:/") {
		return strings.TrimPrefix(url, "http:/")
	}
	if strings.HasPrefix(url, "https:/") {
		return strings.TrimPrefix(url, "https:/")
	}
	return url
}

func extractPathFromURL(url string) string {
	// 提取路径部分
	if idx := strings.Index(url, "/"); idx != -1 {
		return url[idx:]
	}
	return "/"
}

func convertHeaders(headers []map[string]string) map[string]string {
	result := make(map[string]string)
	for _, header := range headers {
		if key, ok := header["key"]; ok {
			if value, ok := header["value"]; ok {
				result[key] = value
			}
		}
	}
	return result
}

func convertInputSchemaToArgs(inputSchema map[string]interface{}, requestTemplate RequestTemplate) ([]model.ArgConfig, error) {
	var args []model.ArgConfig

	properties, ok := inputSchema["properties"].(map[string]interface{})
	if !ok {
		return args, nil
	}

	required, _ := inputSchema["required"].([]interface{})
	requiredMap := make(map[string]bool)
	for _, req := range required {
		if reqStr, ok := req.(string); ok {
			requiredMap[reqStr] = true
		}
	}

	for name, prop := range properties {
		if propMap, ok := prop.(map[string]interface{}); ok {
			arg := model.ArgConfig{
				Name:        name,
				Type:        getString(propMap, "type", "string"),
				Description: getString(propMap, "description", ""),
				Required:    requiredMap[name],
				In:          determineArgLocation(name, requestTemplate),
			}
			args = append(args, arg)
		}
	}

	return args, nil
}

func getString(m map[string]interface{}, key, defaultValue string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return defaultValue
}

func determineArgLocation(argName string, requestTemplate RequestTemplate) string {
	// 根据模板配置确定参数位置
	if requestTemplate.ArgsToJsonBody {
		return "body"
	}
	if requestTemplate.ArgsToUrlParam {
		return "query"
	}

	// 检查 URL 中是否有路径参数
	if strings.Contains(requestTemplate.URL, "{{.args."+argName+"}}") {
		return "path"
	}

	return "body" // 默认值
}
