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

package model

import (
	"regexp"
	"time"
)

// McpServerConfig MCP Server Filter configuration
type McpServerConfig struct {
	ServerInfo        ServerInfo               `yaml:"server_info" json:"server_info"`
	Endpoint          string                   `yaml:"endpoint" json:"endpoint" default:"/mcp"`
	Tools             []ToolConfig             `yaml:"tools,omitempty" json:"tools,omitempty"`
	Resources         []ResourceConfig         `yaml:"resources,omitempty" json:"resources,omitempty"`
	ResourceTemplates []ResourceTemplateConfig `yaml:"resource_templates,omitempty" json:"resource_templates,omitempty"`
	Prompts           []PromptConfig           `yaml:"prompts,omitempty" json:"prompts,omitempty"`
}

// ServerInfo server information
type ServerInfo struct {
	Name         string `yaml:"name" json:"name" default:"Pixiu MCP Server"`
	Version      string `yaml:"version" json:"version" default:"1.0.0"`
	Description  string `yaml:"description,omitempty" json:"description,omitempty" default:"MCP Server powered by Apache Dubbo-go-pixiu"`
	Instructions string `yaml:"instructions,omitempty" json:"instructions,omitempty" default:"Use the provided tools to interact with backend services."`
}

// ToolConfig tool configuration
type ToolConfig struct {
	Name        string        `yaml:"name" json:"name"`
	Description string        `yaml:"description" json:"description"`
	Cluster     string        `yaml:"cluster" json:"cluster"`
	BackendURL  string        `yaml:"backend_url,omitempty" json:"backend_url,omitempty"`
	Request     RequestConfig `yaml:"request" json:"request"`
	Args        []ArgConfig   `yaml:"args,omitempty" json:"args,omitempty"`
}

// RequestConfig request configuration
type RequestConfig struct {
	Method  string            `yaml:"method" json:"method" default:"GET"`
	Path    string            `yaml:"path" json:"path"`
	Headers map[string]string `yaml:"headers,omitempty" json:"headers,omitempty"`
	Timeout string            `yaml:"timeout,omitempty" json:"timeout,omitempty" default:"30s"`
}

// ArgConfig parameter configuration (simplified)
type ArgConfig struct {
	Name        string   `yaml:"name" json:"name"`
	Type        string   `yaml:"type" json:"type" default:"string"`
	In          string   `yaml:"in" json:"in"`
	Description string   `yaml:"description,omitempty" json:"description,omitempty"`
	Required    bool     `yaml:"required,omitempty" json:"required,omitempty"`
	Default     any      `yaml:"default,omitempty" json:"default,omitempty"`
	Enum        []string `yaml:"enum,omitempty" json:"enum,omitempty"`
}

// ResourceConfig resource configuration
type ResourceConfig struct {
	Name        string         `yaml:"name" json:"name"`
	URI         string         `yaml:"uri" json:"uri"`
	Description string         `yaml:"description,omitempty" json:"description,omitempty"`
	MIMEType    string         `yaml:"mime_type,omitempty" json:"mime_type,omitempty"`
	Source      ResourceSource `yaml:"source" json:"source"`
}

// ResourceSource resource source configuration (simplified)
type ResourceSource struct {
	Type     string `yaml:"type" json:"type"`
	Path     string `yaml:"path,omitempty" json:"path,omitempty"`         // for file type
	URL      string `yaml:"url,omitempty" json:"url,omitempty"`           // for url type
	Content  string `yaml:"content,omitempty" json:"content,omitempty"`   // for inline type
	Template string `yaml:"template,omitempty" json:"template,omitempty"` // for template type
}

// ResourceTemplateConfig resource template configuration
type ResourceTemplateConfig struct {
	Name        string                       `yaml:"name" json:"name"`
	URITemplate string                       `yaml:"uri_template" json:"uri_template"`
	Title       string                       `yaml:"title,omitempty" json:"title,omitempty"`
	Description string                       `yaml:"description,omitempty" json:"description,omitempty"`
	MIMEType    string                       `yaml:"mime_type,omitempty" json:"mime_type,omitempty"`
	Parameters  []ResourceTemplateParameter  `yaml:"parameters,omitempty" json:"parameters,omitempty"`
	Annotations *ResourceTemplateAnnotations `yaml:"annotations,omitempty" json:"annotations,omitempty"`
}

// ResourceTemplateParameter resource template parameter
type ResourceTemplateParameter struct {
	Name        string   `yaml:"name" json:"name"`
	Type        string   `yaml:"type" json:"type" default:"string"`
	Description string   `yaml:"description,omitempty" json:"description,omitempty"`
	Required    bool     `yaml:"required,omitempty" json:"required,omitempty" default:"false"`
	Enum        []string `yaml:"enum,omitempty" json:"enum,omitempty"`
	Default     any      `yaml:"default,omitempty" json:"default,omitempty"`
}

// ResourceTemplateAnnotations resource template annotations
type ResourceTemplateAnnotations struct {
	Audience     []string `yaml:"audience,omitempty" json:"audience,omitempty"`
	Priority     *float64 `yaml:"priority,omitempty" json:"priority,omitempty"`
	LastModified string   `yaml:"last_modified,omitempty" json:"last_modified,omitempty"`
}

// ComputedParameter computed parameter (for internal processing, simplified)
type ComputedParameter struct {
	Name        string
	Type        string
	In          string
	Description string
	Required    bool
	Enum        []string
	Default     any
}

// PromptConfig prompt configuration
type PromptConfig struct {
	Name        string                 `yaml:"name" json:"name"`
	Title       string                 `yaml:"title,omitempty" json:"title,omitempty"`
	Description string                 `yaml:"description,omitempty" json:"description,omitempty"`
	Arguments   []PromptArgumentConfig `yaml:"arguments,omitempty" json:"arguments,omitempty"`
	Messages    []PromptMessageConfig  `yaml:"messages" json:"messages"`
}

// PromptArgumentConfig prompt argument configuration
type PromptArgumentConfig struct {
	Name        string `yaml:"name" json:"name"`
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
	Required    bool   `yaml:"required,omitempty" json:"required,omitempty"`
}

// PromptMessageConfig prompt message configuration
type PromptMessageConfig struct {
	Role    string `yaml:"role" json:"role"` // "user" or "assistant"
	Content string `yaml:"content" json:"content"`
}

// RegistryConfig registry configuration
type RegistryConfig struct {
	ToolConfigs     map[string]ToolConfig     `yaml:"toolConfigs"`
	ResourceConfigs map[string]ResourceConfig `yaml:"resourceConfigs"`
	PromptConfigs   map[string]PromptConfig   `yaml:"promptConfigs"`
	LastUpdated     time.Time                 `yaml:"lastUpdated"`
}

// GetAllParameters gets all parameters of the tool
func (tc *ToolConfig) GetAllParameters() ([]ComputedParameter, error) {
	var allParams []ComputedParameter

	// 1. Automatically extract path parameters
	pathParams := GetPathParameterNames(tc.Request.Path)
	for _, paramName := range pathParams {
		// Find corresponding arg configuration
		var argConfig *ArgConfig
		for _, arg := range tc.Args {
			if arg.Name == paramName && arg.In == "path" {
				argConfig = &arg
				break
			}
		}

		// Create computed parameter
		computed := ComputedParameter{
			Name:     paramName,
			Type:     "string", // Default type
			In:       "path",
			Required: true, // Path parameters are always required
		}

		// Apply arg configuration
		if argConfig != nil {
			computed.Type = argConfig.Type
			computed.Description = argConfig.Description
		}

		allParams = append(allParams, computed)
	}

	// 2. Add non-path parameters
	for _, arg := range tc.Args {
		if arg.In != "path" {
			computed := ComputedParameter{
				Name:        arg.Name,
				Type:        arg.Type,
				In:          arg.In,
				Description: arg.Description,
				Required:    arg.Required,
				Enum:        arg.Enum,
				Default:     arg.Default,
				// Simplified: removed complex validation fields
			}
			allParams = append(allParams, computed)
		}
	}

	return allParams, nil
}

// GetPathParameterNames gets all parameter names in the path template
func GetPathParameterNames(pathTemplate string) []string {
	re := regexp.MustCompile(`\{([^}]+)}`)
	matches := re.FindAllStringSubmatch(pathTemplate, -1)

	// Initialize as empty slice instead of nil
	names := []string{}
	for _, match := range matches {
		names = append(names, match[1])
	}

	return names
}

// DeepCopy returns a new independent copy of Config
// Deep copy slices/maps to avoid sharing pointers with the factory
func (config *McpServerConfig) DeepCopy() *McpServerConfig {
	if config == nil {
		return nil
	}

	cpConfig := *config

	// Deep copy Tools
	if config.Tools != nil {
		cpConfig.Tools = make([]ToolConfig, len(config.Tools))
		for i := range config.Tools {
			cpConfig.Tools[i] = *config.Tools[i].DeepCopy()
		}
	}

	// Deep copy Resources
	if config.Resources != nil {
		cpConfig.Resources = make([]ResourceConfig, len(config.Resources))
		copy(cpConfig.Resources, config.Resources)
	}

	// Deep copy ResourceTemplates
	if config.ResourceTemplates != nil {
		cpConfig.ResourceTemplates = make([]ResourceTemplateConfig, len(config.ResourceTemplates))
		for i := range config.ResourceTemplates {
			cpConfig.ResourceTemplates[i] = *config.ResourceTemplates[i].DeepCopy()
		}
	}

	// Deep copy Prompts
	if config.Prompts != nil {
		cpConfig.Prompts = make([]PromptConfig, len(config.Prompts))
		for i := range config.Prompts {
			cpConfig.Prompts[i] = *config.Prompts[i].DeepCopy()
		}
	}

	return &cpConfig
}

// DeepCopy returns a new independent copy of Config
// Deep copy slices/maps to avoid sharing pointers with the factory
func (toolConfig *ToolConfig) DeepCopy() *ToolConfig {
	if toolConfig == nil {
		return nil
	}
	cpConfig := *toolConfig
	cpConfig.Request = *toolConfig.Request.DeepCopy()

	if toolConfig.Args != nil {
		cpConfig.Args = make([]ArgConfig, len(toolConfig.Args))
		for index := range toolConfig.Args {
			argPtr := &toolConfig.Args[index]
			cpConfig.Args[index] = *argPtr.DeepCopy()
		}
	}

	return &cpConfig
}

// DeepCopy returns a new independent copy of Config
// Deep copy slices/maps to avoid sharing pointers with the factory
func (config *RequestConfig) DeepCopy() *RequestConfig {
	if config == nil {
		return nil
	}
	cpConfig := *config
	if config.Headers != nil {
		cpConfig.Headers = make(map[string]string, len(config.Headers))
		for k, v := range config.Headers {
			cpConfig.Headers[k] = v
		}
	}

	return &cpConfig
}

// DeepCopy returns a new independent copy of Config
// Deep copy slices/maps to avoid sharing pointers with the factory
func (config *ArgConfig) DeepCopy() *ArgConfig {
	if config == nil {
		return nil
	}
	cpConfig := *config
	if config.Enum != nil {
		cpConfig.Enum = make([]string, len(config.Enum))
		copy(cpConfig.Enum, config.Enum)
	}
	return &cpConfig
}

// DeepCopy returns a new independent copy of Config
// Deep copy slices/maps to avoid sharing pointers with the factory
func (param *ResourceTemplateParameter) DeepCopy() *ResourceTemplateParameter {
	if param == nil {
		return nil
	}

	cpParam := *param

	if param.Enum != nil {
		cpParam.Enum = make([]string, len(param.Enum))
		copy(cpParam.Enum, param.Enum)
	}

	return &cpParam
}

// DeepCopy returns a new independent copy of Config
// Deep copy slices/maps to avoid sharing pointers with the factory
func (ann *ResourceTemplateAnnotations) DeepCopy() *ResourceTemplateAnnotations {
	if ann == nil {
		return nil
	}

	cpAnn := *ann

	// Deep copy Audience slice
	if ann.Audience != nil {
		cpAnn.Audience = make([]string, len(ann.Audience))
		copy(cpAnn.Audience, ann.Audience)
	}

	// Deep copy Priority pointer
	if ann.Priority != nil {
		p := *ann.Priority
		cpAnn.Priority = &p
	}

	return &cpAnn
}

// DeepCopy returns a new independent copy of Config
// Deep copy slices/maps to avoid sharing pointers with the factory
func (config *ResourceTemplateConfig) DeepCopy() *ResourceTemplateConfig {
	if config == nil {
		return nil
	}

	cpConfig := *config

	// Deep copy Parameters slice
	if config.Parameters != nil {
		cpConfig.Parameters = make([]ResourceTemplateParameter, len(config.Parameters))
		for i := range config.Parameters {
			paramPtr := &config.Parameters[i]
			cpConfig.Parameters[i] = *paramPtr.DeepCopy()
		}
	}

	// Deep copy Annotations pointer
	if config.Annotations != nil {
		cpConfig.Annotations = config.Annotations.DeepCopy()
	}

	return &cpConfig
}

// DeepCopy returns a new independent copy of Config
// Deep copy slices/maps to avoid sharing pointers with the factory
func (config *PromptConfig) DeepCopy() *PromptConfig {
	if config == nil {
		return nil
	}

	cp := *config

	// Deep copy Arguments slice
	if config.Arguments != nil {
		cp.Arguments = make([]PromptArgumentConfig, len(config.Arguments))
		copy(cp.Arguments, config.Arguments)
	}

	// Deep copy Messages slice
	if config.Messages != nil {
		cp.Messages = make([]PromptMessageConfig, len(config.Messages))
		copy(cp.Messages, config.Messages)
	}

	return &cp
}
