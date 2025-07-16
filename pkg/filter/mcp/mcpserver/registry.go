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
	"sync"

	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

// ToolRegistry 工具注册表，线程安全
type ToolRegistry struct {
	mu        sync.RWMutex
	tools     map[string]ToolConfig
	resources map[string]ResourceConfig
	prompts   map[string]PromptConfig
}

// NewToolRegistry 创建新的工具注册表
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools:     make(map[string]ToolConfig),
		resources: make(map[string]ResourceConfig),
		prompts:   make(map[string]PromptConfig),
	}
}

// RegisterTool 注册工具
func (r *ToolRegistry) RegisterTool(tool ToolConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tools[tool.Name]; exists {
		return fmt.Errorf("tool %s already exists", tool.Name)
	}

	r.tools[tool.Name] = tool
	return nil
}

// RegisterResource 注册资源
func (r *ToolRegistry) RegisterResource(resource ResourceConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.resources[resource.Name]; exists {
		return fmt.Errorf("resource %s already exists", resource.Name)
	}

	r.resources[resource.Name] = resource
	return nil
}

// GetTool 获取工具配置
func (r *ToolRegistry) GetTool(name string) (ToolConfig, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tool, exists := r.tools[name]
	return tool, exists
}

// GetResource 获取资源配置
func (r *ToolRegistry) GetResource(name string) (ResourceConfig, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	resource, exists := r.resources[name]
	return resource, exists
}

// ListTools 列出所有工具
func (r *ToolRegistry) ListTools() []ToolConfig {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tools := make([]ToolConfig, 0, len(r.tools))
	for _, tool := range r.tools {
		tools = append(tools, tool)
	}
	return tools
}

// ListResources 列出所有资源
func (r *ToolRegistry) ListResources() []ResourceConfig {
	r.mu.RLock()
	defer r.mu.RUnlock()

	resources := make([]ResourceConfig, 0, len(r.resources))
	for _, resource := range r.resources {
		resources = append(resources, resource)
	}
	return resources
}

// ToMCPTools 将工具配置转换为工具列表
func (r *ToolRegistry) ToMCPTools() ([]map[string]any, error) {
	tools := r.ListTools()
	mcpTools := make([]map[string]any, 0, len(tools))

	for _, tool := range tools {
		// 根据 MCP 协议规范构建工具
		mcpTool := map[string]any{
			"name":        tool.Name,
			"description": tool.Description,
			"inputSchema": r.convertToInputSchema(tool),
		}
		mcpTools = append(mcpTools, mcpTool)
	}

	return mcpTools, nil
}

// convertToInputSchema 将工具参数转换为 MCP inputSchema 格式
func (r *ToolRegistry) convertToInputSchema(tool ToolConfig) map[string]any {
	allParams, err := tool.GetAllParameters()
	if err != nil {
		logger.Errorf("failed to get parameters for tool %s: %v", tool.Name, err)
		return map[string]any{
			"type":       "object",
			"properties": map[string]any{},
		}
	}

	properties := make(map[string]any)
	required := make([]string, 0)

	for _, param := range allParams {
		propSchema := map[string]any{
			"type":        param.Type,
			"description": param.Description,
		}

		if len(param.Enum) > 0 {
			propSchema["enum"] = param.Enum
		}

		if param.Default != nil {
			propSchema["default"] = param.Default
		}

		properties[param.Name] = propSchema

		if param.Required {
			required = append(required, param.Name)
		}
	}

	schema := map[string]any{
		"type":       "object",
		"properties": properties,
	}

	if len(required) > 0 {
		schema["required"] = required
	}

	return schema
}

// ToMCPResources 将资源配置转换为资源列表
func (r *ToolRegistry) ToMCPResources() ([]map[string]any, error) {
	resources := r.ListResources()
	mcpResources := make([]map[string]any, 0, len(resources))

	for _, resource := range resources {
		mcpResource := map[string]any{
			"uri":         resource.URI,
			"name":        resource.Name,
			"description": resource.Description,
			"mimeType":    resource.MIMEType,
		}
		mcpResources = append(mcpResources, mcpResource)
	}

	return mcpResources, nil
}

// Count 返回注册的工具、资源和提示词数量
func (r *ToolRegistry) Count() (int, int, int) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.tools), len(r.resources), len(r.prompts)
}

// RegisterPrompt 注册提示词
func (r *ToolRegistry) RegisterPrompt(prompt PromptConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.prompts[prompt.Name]; exists {
		return fmt.Errorf("prompt %s already exists", prompt.Name)
	}

	r.prompts[prompt.Name] = prompt
	return nil
}

// GetPrompt 获取提示词
func (r *ToolRegistry) GetPrompt(name string) (PromptConfig, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	prompt, exists := r.prompts[name]
	return prompt, exists
}

// ListPrompts 列出所有提示词
func (r *ToolRegistry) ListPrompts() []PromptConfig {
	r.mu.RLock()
	defer r.mu.RUnlock()

	prompts := make([]PromptConfig, 0, len(r.prompts))
	for _, prompt := range r.prompts {
		prompts = append(prompts, prompt)
	}

	return prompts
}

// ToMCPPrompts 将提示词配置转换为 MCP 提示词列表
func (r *ToolRegistry) ToMCPPrompts() ([]map[string]any, error) {
	prompts := r.ListPrompts()
	mcpPrompts := make([]map[string]any, 0, len(prompts))

	for _, prompt := range prompts {
		mcpPrompt := map[string]any{
			"name":        prompt.Name,
			"description": prompt.Description,
		}

		if prompt.Title != "" {
			mcpPrompt["title"] = prompt.Title
		}

		if len(prompt.Arguments) > 0 {
			args := make([]map[string]any, 0, len(prompt.Arguments))
			for _, arg := range prompt.Arguments {
				argMap := map[string]any{
					"name":        arg.Name,
					"description": arg.Description,
					"required":    arg.Required,
				}
				args = append(args, argMap)
			}
			mcpPrompt["arguments"] = args
		}

		mcpPrompts = append(mcpPrompts, mcpPrompt)
	}

	return mcpPrompts, nil
}
