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
)

import (
	"github.com/mark3labs/mcp-go/mcp"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

// ToolRegistry tool registry, thread-safe (optimized with single indexing)
type ToolRegistry struct {
	mu                sync.RWMutex
	tools             map[string]ToolConfig
	resources         map[string]ResourceConfig         // indexed by URI
	resourceTemplates map[string]ResourceTemplateConfig // indexed by name
	prompts           map[string]PromptConfig

	// TODO: Dynamic update support - add when integrating with Nacos
	// changeListeners   []ChangeListener              // change listeners
	// nacosClient      nacos.ConfigClient            // Nacos config client
	// serviceDiscovery nacos.NamingClient            // Nacos service discovery client
}

// NewToolRegistry creates a new tool registry
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools:             make(map[string]ToolConfig),
		resources:         make(map[string]ResourceConfig),
		resourceTemplates: make(map[string]ResourceTemplateConfig),
		prompts:           make(map[string]PromptConfig),
	}
}

// RegisterTool registers a tool
func (r *ToolRegistry) RegisterTool(tool ToolConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tools[tool.Name]; exists {
		return fmt.Errorf("tool %s already exists", tool.Name)
	}

	r.tools[tool.Name] = tool

	// TODO: Dynamic update notification - enable when integrating with Nacos
	// r.notifyToolsListChanged()

	return nil
}

// RegisterResource registers a resource (indexed by URI as per MCP specification)
func (r *ToolRegistry) RegisterResource(resource ResourceConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.resources[resource.URI]; exists {
		return fmt.Errorf("resource with URI %s already exists", resource.URI)
	}

	// Register resource by URI as per MCP specification
	r.resources[resource.URI] = resource
	return nil
}

// GetTool gets tool configuration
func (r *ToolRegistry) GetTool(name string) (ToolConfig, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tool, exists := r.tools[name]
	return tool, exists
}

// GetResourceByURI gets resource configuration (by URI, O(1) lookup)
func (r *ToolRegistry) GetResourceByURI(uri string) (ResourceConfig, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Direct O(1) lookup by URI
	resource, exists := r.resources[uri]
	return resource, exists
}

// RegisterResourceTemplate registers a resource template
func (r *ToolRegistry) RegisterResourceTemplate(template ResourceTemplateConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.resourceTemplates[template.Name]; exists {
		return fmt.Errorf("resource template %s already exists", template.Name)
	}

	r.resourceTemplates[template.Name] = template
	return nil
}

// GetResourceTemplate gets resource template configuration
func (r *ToolRegistry) GetResourceTemplate(name string) (ResourceTemplateConfig, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	template, exists := r.resourceTemplates[name]
	return template, exists
}

// ListResourceTemplates lists all resource templates
func (r *ToolRegistry) ListResourceTemplates() []ResourceTemplateConfig {
	r.mu.RLock()
	defer r.mu.RUnlock()

	templates := make([]ResourceTemplateConfig, 0, len(r.resourceTemplates))
	for _, template := range r.resourceTemplates {
		templates = append(templates, template)
	}
	return templates
}

// ListTools lists all tools
func (r *ToolRegistry) ListTools() []ToolConfig {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tools := make([]ToolConfig, 0, len(r.tools))
	for _, tool := range r.tools {
		tools = append(tools, tool)
	}
	return tools
}

// ListResources lists all resources
func (r *ToolRegistry) ListResources() []ResourceConfig {
	r.mu.RLock()
	defer r.mu.RUnlock()

	resources := make([]ResourceConfig, 0, len(r.resources))
	for _, resource := range r.resources {
		resources = append(resources, resource)
	}
	return resources
}

// ToMCPTools converts tool configurations to tool list
func (r *ToolRegistry) ToMCPTools() ([]map[string]any, error) {
	tools := r.ListTools()
	mcpTools := make([]map[string]any, 0, len(tools))

	for _, tool := range tools {
		// Build tool according to MCP protocol specification
		mcpTool := map[string]any{
			"name":        tool.Name,
			"description": tool.Description,
			"inputSchema": r.convertToInputSchema(tool),
		}
		mcpTools = append(mcpTools, mcpTool)
	}

	return mcpTools, nil
}

// convertToInputSchema converts tool parameters to MCP inputSchema format
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

// ToMCPResources converts resource configurations to MCP resource list using mcp-go structures
func (r *ToolRegistry) ToMCPResources() ([]mcp.Resource, error) {
	resources := r.ListResources()
	mcpResources := make([]mcp.Resource, 0, len(resources))

	for _, resource := range resources {
		// Use mcp-go Resource structure
		mcpResource := mcp.Resource{
			URI:         resource.URI,
			Name:        resource.Name,
			Description: resource.Description,
			MIMEType:    resource.MIMEType,
		}
		mcpResources = append(mcpResources, mcpResource)
	}

	return mcpResources, nil
}

// ToMCPResourceTemplates converts resource template configurations to MCP resource template list
func (r *ToolRegistry) ToMCPResourceTemplates() ([]map[string]any, error) {
	templates := r.ListResourceTemplates()
	mcpTemplates := make([]map[string]any, 0, len(templates))

	for _, template := range templates {
		mcpTemplate := map[string]any{
			"uriTemplate": template.URITemplate,
			"name":        template.Name,
			"description": template.Description,
			"mimeType":    template.MIMEType,
		}

		// Add optional fields
		if template.Title != "" {
			mcpTemplate["title"] = template.Title
		}

		// Add annotations
		if template.Annotations != nil {
			annotations := make(map[string]any)
			if len(template.Annotations.Audience) > 0 {
				annotations["audience"] = template.Annotations.Audience
			}
			if template.Annotations.Priority != nil {
				annotations["priority"] = *template.Annotations.Priority
			}
			if template.Annotations.LastModified != "" {
				annotations["lastModified"] = template.Annotations.LastModified
			}
			if len(annotations) > 0 {
				mcpTemplate["annotations"] = annotations
			}
		}

		mcpTemplates = append(mcpTemplates, mcpTemplate)
	}

	return mcpTemplates, nil
}

// RegisterPrompt registers a prompt
func (r *ToolRegistry) RegisterPrompt(prompt PromptConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.prompts[prompt.Name]; exists {
		return fmt.Errorf("prompt %s already exists", prompt.Name)
	}

	r.prompts[prompt.Name] = prompt
	return nil
}

// GetPrompt gets a prompt
func (r *ToolRegistry) GetPrompt(name string) (PromptConfig, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	prompt, exists := r.prompts[name]
	return prompt, exists
}

// ListPrompts lists all prompts
func (r *ToolRegistry) ListPrompts() []PromptConfig {
	r.mu.RLock()
	defer r.mu.RUnlock()

	prompts := make([]PromptConfig, 0, len(r.prompts))
	for _, prompt := range r.prompts {
		prompts = append(prompts, prompt)
	}

	return prompts
}

// ToMCPPrompts converts prompt configurations to MCP prompt list
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

// TODO: Dynamic update functionality - implement when integrating with Nacos
