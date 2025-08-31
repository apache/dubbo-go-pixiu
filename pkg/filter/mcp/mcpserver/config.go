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
	"regexp"
)

type (
	// Config MCP Server Filter configuration
	Config struct {
		ServerInfo        ServerInfo               `yaml:"server_info" json:"server_info"`
		Endpoint          string                   `yaml:"endpoint" json:"endpoint" default:"/mcp"`
		Tools             []ToolConfig             `yaml:"tools,omitempty" json:"tools,omitempty"`
		Resources         []ResourceConfig         `yaml:"resources,omitempty" json:"resources,omitempty"`
		ResourceTemplates []ResourceTemplateConfig `yaml:"resource_templates,omitempty" json:"resource_templates,omitempty"`
		Prompts           []PromptConfig           `yaml:"prompts,omitempty" json:"prompts,omitempty"`
	}

	// ServerInfo server information
	ServerInfo struct {
		Name         string `yaml:"name" json:"name" default:"Pixiu MCP Server"`
		Version      string `yaml:"version" json:"version" default:"1.0.0"`
		Description  string `yaml:"description,omitempty" json:"description,omitempty" default:"MCP Server powered by Apache Dubbo-go-pixiu"`
		Instructions string `yaml:"instructions,omitempty" json:"instructions,omitempty" default:"Use the provided tools to interact with backend services."`
	}

	// ToolConfig tool configuration
	ToolConfig struct {
		Name        string        `yaml:"name" json:"name"`
		Description string        `yaml:"description" json:"description"`
		Cluster     string        `yaml:"cluster" json:"cluster"`
		Request     RequestConfig `yaml:"request" json:"request"`
		Args        []ArgConfig   `yaml:"args,omitempty" json:"args,omitempty"`
	}

	// RequestConfig request configuration
	RequestConfig struct {
		Method  string            `yaml:"method" json:"method" default:"GET"`
		Path    string            `yaml:"path" json:"path"`
		Headers map[string]string `yaml:"headers,omitempty" json:"headers,omitempty"`
		Timeout string            `yaml:"timeout,omitempty" json:"timeout,omitempty" default:"30s"`
	}

	// ArgConfig parameter configuration (simplified)
	ArgConfig struct {
		Name        string   `yaml:"name" json:"name"`
		Type        string   `yaml:"type" json:"type" default:"string"`
		In          string   `yaml:"in" json:"in"`
		Description string   `yaml:"description,omitempty" json:"description,omitempty"`
		Required    bool     `yaml:"required,omitempty" json:"required,omitempty"`
		Default     any      `yaml:"default,omitempty" json:"default,omitempty"`
		Enum        []string `yaml:"enum,omitempty" json:"enum,omitempty"`
	}

	// ResourceConfig resource configuration
	ResourceConfig struct {
		Name        string         `yaml:"name" json:"name"`
		URI         string         `yaml:"uri" json:"uri"`
		Description string         `yaml:"description,omitempty" json:"description,omitempty"`
		MIMEType    string         `yaml:"mime_type,omitempty" json:"mime_type,omitempty"`
		Source      ResourceSource `yaml:"source" json:"source"`
	}

	// ResourceSource resource source configuration (simplified)
	ResourceSource struct {
		Type     string `yaml:"type" json:"type"`
		Path     string `yaml:"path,omitempty" json:"path,omitempty"`         // for file type
		URL      string `yaml:"url,omitempty" json:"url,omitempty"`           // for url type
		Content  string `yaml:"content,omitempty" json:"content,omitempty"`   // for inline type
		Template string `yaml:"template,omitempty" json:"template,omitempty"` // for template type
	}

	// ResourceTemplateConfig resource template configuration
	ResourceTemplateConfig struct {
		Name        string                       `yaml:"name" json:"name"`
		URITemplate string                       `yaml:"uri_template" json:"uri_template"`
		Title       string                       `yaml:"title,omitempty" json:"title,omitempty"`
		Description string                       `yaml:"description,omitempty" json:"description,omitempty"`
		MIMEType    string                       `yaml:"mime_type,omitempty" json:"mime_type,omitempty"`
		Parameters  []ResourceTemplateParameter  `yaml:"parameters,omitempty" json:"parameters,omitempty"`
		Annotations *ResourceTemplateAnnotations `yaml:"annotations,omitempty" json:"annotations,omitempty"`
	}

	// ResourceTemplateParameter resource template parameter
	ResourceTemplateParameter struct {
		Name        string   `yaml:"name" json:"name"`
		Type        string   `yaml:"type" json:"type" default:"string"`
		Description string   `yaml:"description,omitempty" json:"description,omitempty"`
		Required    bool     `yaml:"required,omitempty" json:"required,omitempty" default:"false"`
		Enum        []string `yaml:"enum,omitempty" json:"enum,omitempty"`
		Default     any      `yaml:"default,omitempty" json:"default,omitempty"`
	}

	// ResourceTemplateAnnotations resource template annotations
	ResourceTemplateAnnotations struct {
		Audience     []string `yaml:"audience,omitempty" json:"audience,omitempty"`
		Priority     *float64 `yaml:"priority,omitempty" json:"priority,omitempty"`
		LastModified string   `yaml:"last_modified,omitempty" json:"last_modified,omitempty"`
	}

	// ComputedParameter computed parameter (for internal processing, simplified)
	ComputedParameter struct {
		Name        string
		Type        string
		In          string
		Description string
		Required    bool
		Enum        []string
		Default     any
	}
)

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
			// Simplified: removed Pattern and Format fields
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
