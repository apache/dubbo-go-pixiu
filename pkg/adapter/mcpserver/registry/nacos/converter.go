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
	"strings"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/mcpserver/common/util"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// ConvertNacosToolsToToolConfig converts Nacos Tools to the Filter's ToolConfig
func ConvertNacosToolsToToolConfig(toolsSpec *ToolsSpec) ([]model.ToolConfig, error) {
	var toolConfigs []model.ToolConfig

	for _, nacosTool := range toolsSpec.Tools {
		meta := toolsSpec.ToolsMeta[nacosTool.Name]
		if !meta.Enabled {
			continue
		}

		// Extract json-go-template
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

func convertSingleTool(nacosTool NacosTool, templateData any) (model.ToolConfig, error) {
	// Parse template data
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
		Cluster:     nacosTool.Name, // Directly use the tool name as the cluster name
		BackendURL:  template.RequestTemplate.URL,
		Request: model.RequestConfig{
			Method:  template.RequestTemplate.Method,
			Path:    util.ExtractPathFromURL(template.RequestTemplate.URL),
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

func convertInputSchemaToArgs(inputSchema map[string]any, requestTemplate RequestTemplate) ([]model.ArgConfig, error) {
	var args []model.ArgConfig

	properties, ok := inputSchema["properties"].(map[string]any)
	if !ok {
		return args, nil
	}

	required, _ := inputSchema["required"].([]any)
	requiredMap := make(map[string]bool)
	for _, req := range required {
		if reqStr, ok := req.(string); ok {
			requiredMap[reqStr] = true
		}
	}

	for name, prop := range properties {
		if propMap, ok := prop.(map[string]any); ok {
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

func getString(m map[string]any, key, defaultValue string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return defaultValue
}

func determineArgLocation(argName string, requestTemplate RequestTemplate) string {
	// Determine parameter location based on template configuration
	if requestTemplate.ArgsToJsonBody {
		return "body"
	}
	if requestTemplate.ArgsToUrlParam {
		return "query"
	}

	// Check if URL contains path parameters
	if strings.Contains(requestTemplate.URL, "{{.args."+argName+"}}") {
		return "path"
	}

	return "body" // default
}
