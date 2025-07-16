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
	"regexp"
	"strings"

	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

type (
	// Config MCP Server Filter 配置
	Config struct {
		ServerInfo ServerInfo       `yaml:"server_info" json:"server_info"`
		Endpoint   string           `yaml:"endpoint" json:"endpoint" default:"/mcp"`
		Tools      []ToolConfig     `yaml:"tools,omitempty" json:"tools,omitempty"`
		Resources  []ResourceConfig `yaml:"resources,omitempty" json:"resources,omitempty"`
		Prompts    []PromptConfig   `yaml:"prompts,omitempty" json:"prompts,omitempty"`
	}

	// ServerInfo 服务器信息
	ServerInfo struct {
		Name         string `yaml:"name" json:"name" default:"Pixiu MCP Server"`
		Version      string `yaml:"version" json:"version" default:"1.0.0"`
		Description  string `yaml:"description,omitempty" json:"description,omitempty" default:"MCP Server powered by Apache Dubbo-go-pixiu"`
		Instructions string `yaml:"instructions,omitempty" json:"instructions,omitempty" default:"Use the provided tools to interact with backend services."`
	}

	// ToolConfig 工具配置
	ToolConfig struct {
		Name        string          `yaml:"name" json:"name"`
		Description string          `yaml:"description" json:"description"`
		Cluster     string          `yaml:"cluster" json:"cluster"`
		Request     RequestConfig   `yaml:"request" json:"request"`
		Args        []ArgConfig     `yaml:"args,omitempty" json:"args,omitempty"`
		Response    *ResponseConfig `yaml:"response,omitempty" json:"response,omitempty"`
	}

	// RequestConfig 请求配置
	RequestConfig struct {
		Method  string            `yaml:"method" json:"method" default:"GET"`
		Path    string            `yaml:"path" json:"path"`
		Headers map[string]string `yaml:"headers,omitempty" json:"headers,omitempty"`
		Timeout string            `yaml:"timeout,omitempty" json:"timeout,omitempty" default:"30s"`
	}

	// ArgConfig 参数配置
	ArgConfig struct {
		Name        string `yaml:"name" json:"name"`
		Type        string `yaml:"type" json:"type" default:"string"`
		In          string `yaml:"in" json:"in"`
		Description string `yaml:"description,omitempty" json:"description,omitempty"`
		Required    bool   `yaml:"required,omitempty" json:"required,omitempty" default:"false"`
		Default     any    `yaml:"default,omitempty" json:"default,omitempty"`

		// 验证选项
		Enum      []string `yaml:"enum,omitempty" json:"enum,omitempty"`
		Pattern   string   `yaml:"pattern,omitempty" json:"pattern,omitempty"`
		Format    string   `yaml:"format,omitempty" json:"format,omitempty"`
		MinLength *int     `yaml:"min_length,omitempty" json:"min_length,omitempty"`
		MaxLength *int     `yaml:"max_length,omitempty" json:"max_length,omitempty"`
		Minimum   *float64 `yaml:"minimum,omitempty" json:"minimum,omitempty"`
		Maximum   *float64 `yaml:"maximum,omitempty" json:"maximum,omitempty"`

		// 高级验证规则
		Validate *ValidateConfig `yaml:"validate,omitempty" json:"validate,omitempty"`
	}

	// ValidateConfig 验证规则配置
	ValidateConfig struct {
		Required bool     `yaml:"required,omitempty" json:"required,omitempty"`
		Enum     []string `yaml:"enum,omitempty" json:"enum,omitempty"`
		Pattern  string   `yaml:"pattern,omitempty" json:"pattern,omitempty"`
		Format   string   `yaml:"format,omitempty" json:"format,omitempty"`
		Min      *float64 `yaml:"min,omitempty" json:"min,omitempty"`
		Max      *float64 `yaml:"max,omitempty" json:"max,omitempty"`
	}

	// ResponseConfig 响应配置
	ResponseConfig struct {
		Format      string `yaml:"format,omitempty" json:"format,omitempty" default:"json"`
		Description string `yaml:"description,omitempty" json:"description,omitempty"`
		PrependText string `yaml:"prepend_text,omitempty" json:"prepend_text,omitempty"`
		AppendText  string `yaml:"append_text,omitempty" json:"append_text,omitempty"`
		Transform   string `yaml:"transform,omitempty" json:"transform,omitempty" default:"none"`
	}

	// ResourceConfig 资源配置
	ResourceConfig struct {
		Name        string         `yaml:"name" json:"name"`
		URI         string         `yaml:"uri" json:"uri"`
		Description string         `yaml:"description,omitempty" json:"description,omitempty"`
		MIMEType    string         `yaml:"mime_type,omitempty" json:"mime_type,omitempty"`
		Source      ResourceSource `yaml:"source" json:"source"`
	}

	// ResourceSource 资源来源配置
	ResourceSource struct {
		Type   string         `yaml:"type" json:"type"`
		Config map[string]any `yaml:"config,omitempty" json:"config,omitempty"`
	}

	// ComputedParameter 计算得出的参数（用于内部处理）
	ComputedParameter struct {
		Name        string
		Type        string
		In          string
		Description string
		Required    bool
		Enum        []string
		Default     any
		MinLength   *int
		MaxLength   *int
		Minimum     *float64
		Maximum     *float64
		Pattern     string
		Format      string
		Example     string
	}
)

// Validate 验证工具配置的有效性
func (tc *ToolConfig) Validate() error {
	// 验证基本字段
	if tc.Name == "" {
		return fmt.Errorf("tool name is required")
	}
	if tc.Cluster == "" {
		return fmt.Errorf("tool cluster is required")
	}

	// 验证请求配置
	if tc.Request.Method == "" {
		return fmt.Errorf("request method is required")
	}

	// 验证路径格式
	if !strings.HasPrefix(tc.Request.Path, "/") {
		return fmt.Errorf("request path must start with '/': %s", tc.Request.Path)
	}

	// 提取路径参数
	pathParams := GetPathParameterNames(tc.Request.Path)

	// 验证参数配置
	argNames := make(map[string]bool)
	pathArgNames := make(map[string]bool)

	for _, arg := range tc.Args {
		// 验证参数名称
		if arg.Name == "" {
			return fmt.Errorf("arg name is required")
		}
		if argNames[arg.Name] {
			return fmt.Errorf("duplicate arg name: %s", arg.Name)
		}
		argNames[arg.Name] = true

		// 验证参数位置
		if arg.In != "path" && arg.In != "query" && arg.In != "header" && arg.In != "body" {
			return fmt.Errorf("invalid arg location '%s' for arg '%s', must be one of: path, query, header, body",
				arg.In, arg.Name)
		}

		// 记录路径参数
		if arg.In == "path" {
			pathArgNames[arg.Name] = true
		}
	}

	// 验证路径参数是否都有定义
	for _, pathParam := range pathParams {
		if !pathArgNames[pathParam] {
			// 路径参数未在 args 中定义，但这是可以接受的（会自动推断）
			logger.Warnf("path parameter '%s' in path '%s' is not explicitly defined in args",
				pathParam, tc.Request.Path)
		}
	}

	return nil
}

// GetAllParameters 获取工具的所有参数
func (tc *ToolConfig) GetAllParameters() ([]ComputedParameter, error) {
	var allParams []ComputedParameter

	// 1. 自动提取路径参数
	pathParams := GetPathParameterNames(tc.Request.Path)
	for _, paramName := range pathParams {
		// 查找对应的 arg 配置
		var argConfig *ArgConfig
		for _, arg := range tc.Args {
			if arg.Name == paramName && arg.In == "path" {
				argConfig = &arg
				break
			}
		}

		// 创建计算参数
		computed := ComputedParameter{
			Name:     paramName,
			Type:     "string", // 默认类型
			In:       "path",
			Required: true, // 路径参数总是必需的
		}

		// 应用 arg 配置
		if argConfig != nil {
			computed.Type = argConfig.Type
			computed.Description = argConfig.Description
			computed.Pattern = argConfig.Pattern
			computed.Format = argConfig.Format
		}

		allParams = append(allParams, computed)
	}

	// 2. 添加非路径参数
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
				Pattern:     arg.Pattern,
				Format:      arg.Format,
				MinLength:   arg.MinLength,
				MaxLength:   arg.MaxLength,
				Minimum:     arg.Minimum,
				Maximum:     arg.Maximum,
			}
			allParams = append(allParams, computed)
		}
	}

	return allParams, nil
}

// GetPathParameterNames 获取路径模板中的所有参数名称
func GetPathParameterNames(pathTemplate string) []string {
	re := regexp.MustCompile(`\{([^}]+)\}`)
	matches := re.FindAllStringSubmatch(pathTemplate, -1)

	// 初始化为空切片而不是 nil
	names := []string{}
	for _, match := range matches {
		names = append(names, match[1])
	}

	return names
}

// PromptConfig 提示词配置
type PromptConfig struct {
	Name        string                 `yaml:"name" json:"name"`
	Title       string                 `yaml:"title,omitempty" json:"title,omitempty"`
	Description string                 `yaml:"description,omitempty" json:"description,omitempty"`
	Arguments   []PromptArgumentConfig `yaml:"arguments,omitempty" json:"arguments,omitempty"`
	Messages    []PromptMessageConfig  `yaml:"messages" json:"messages"`
}

// PromptArgumentConfig 提示词参数配置
type PromptArgumentConfig struct {
	Name        string `yaml:"name" json:"name"`
	Description string `yaml:"description,omitempty" json:"description,omitempty"`
	Required    bool   `yaml:"required,omitempty" json:"required,omitempty"`
}

// PromptMessageConfig 提示词消息配置
type PromptMessageConfig struct {
	Role    string `yaml:"role" json:"role"` // "user" or "assistant"
	Content string `yaml:"content" json:"content"`
}
