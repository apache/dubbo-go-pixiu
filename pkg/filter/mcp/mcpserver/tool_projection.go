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
	"encoding/json"
)

import (
	"github.com/mark3labs/mcp-go/mcp"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/fingerprint"
	"github.com/apache/dubbo-go-pixiu/pkg/filter/mcp/mcpserver/router"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

const (
	inPath  = "path"
	inQuery = "query"
	inBody  = "body"

	typeString  = "string"
	typeInteger = "integer"
	typeNumber  = "number"
	typeBoolean = "boolean"
)

func BuildMCPTools(toolCfgs []model.ToolConfig) []mcp.Tool {
	tools := make([]mcp.Tool, 0, len(toolCfgs))
	for _, toolCfg := range toolCfgs {
		tools = append(tools, BuildMCPTool(toolCfg))
	}
	return tools
}

func BuildMCPTool(toolCfg model.ToolConfig) mcp.Tool {
	toolOptions := []mcp.ToolOption{
		mcp.WithDescription(toolCfg.Description),
	}

	for _, arg := range toolCfg.Args {
		opts := BuildToolParameterOptions(&arg)
		switch arg.Type {
		case typeString:
			toolOptions = append(toolOptions, mcp.WithString(arg.Name, opts...))
		case typeInteger, typeNumber:
			toolOptions = append(toolOptions, mcp.WithNumber(arg.Name, opts...))
		case typeBoolean:
			toolOptions = append(toolOptions, mcp.WithBoolean(arg.Name, opts...))
		}
	}

	return mcp.NewTool(toolCfg.Name, toolOptions...)
}

func (f *MCPServerFilter) buildMCPTools(toolCfgs []model.ToolConfig) []mcp.Tool {
	return BuildMCPTools(toolCfgs)
}

func (f *MCPServerFilter) buildMCPTool(toolCfg model.ToolConfig) mcp.Tool {
	return BuildMCPTool(toolCfg)
}

func BuildMCPToolMaps(toolCfgs []model.ToolConfig) ([]map[string]any, error) {
	tools := BuildMCPTools(toolCfgs)
	out := make([]map[string]any, 0, len(tools))
	for _, tool := range tools {
		item, err := MCPToolMap(tool)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, nil
}

func MCPToolMap(tool mcp.Tool) (map[string]any, error) {
	data, err := json.Marshal(tool)
	if err != nil {
		return nil, err
	}
	var out map[string]any
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// BuildToolParameterOptions builds the mcp.PropertyOption slice for a tool
// argument. It is the sole schema projection used by handlers and registry
// exports.
func BuildToolParameterOptions(arg *model.ArgConfig) []mcp.PropertyOption {
	opts := []mcp.PropertyOption{mcp.Description(arg.Description)}

	if arg.Required {
		opts = append(opts, mcp.Required())
	}

	if arg.Default != nil {
		switch arg.Type {
		case typeString:
			if defaultStr, ok := arg.Default.(string); ok {
				opts = append(opts, mcp.DefaultString(defaultStr))
			}
		case typeInteger, typeNumber:
			switch defaultVal := arg.Default.(type) {
			case float64:
				opts = append(opts, mcp.DefaultNumber(defaultVal))
			case int:
				opts = append(opts, mcp.DefaultNumber(float64(defaultVal)))
			case int64:
				opts = append(opts, mcp.DefaultNumber(float64(defaultVal)))
			}
		case typeBoolean:
			if defaultBool, ok := arg.Default.(bool); ok {
				opts = append(opts, mcp.DefaultBool(defaultBool))
			}
		}
	}

	if len(arg.Enum) > 0 && arg.Type == typeString {
		opts = append(opts, mcp.Enum(arg.Enum...))
	}

	return opts
}

func visibleToolsFingerprint(tools []model.ToolConfig) string {
	mcpTools := BuildMCPTools(router.DiscoverableTools(tools))
	if len(mcpTools) == 0 {
		return EmptyFingerprint
	}
	items := make([]string, 0, len(mcpTools))
	for _, tool := range mcpTools {
		data, err := json.Marshal(tool)
		if err != nil {
			items = append(items, tool.Name)
			continue
		}
		items = append(items, string(data))
	}
	return fingerprint.StringsSorted(items)
}
