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

package mcpexecutor

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/apache/dubbo-go-pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	h "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/mark3labs/mcp-go/mcp"
	mcp_server "github.com/mark3labs/mcp-go/server"
)

type (
	// FilterFactory is the factory of mcp_executor filter.
	FilterFactory struct {
		cfg *Config
	}

	// McpExecutorFilter is a filter that executes mcp logic.
	McpExecutorFilter struct {
		factory *FilterFactory
	}

	// Simplified structure to decode the /tools/call request body.
	callToolRequestBody struct {
		ToolName  string         `json:"tool_name"`
		Arguments map[string]any `json:"arguments"`
	}
)

// Config returns the config of the mcp_executor filter.
func (f *FilterFactory) Config() any {
	return f.cfg
}

// Apply prepares the mcp.Server and tool schemas.
func (f *FilterFactory) Apply() error {
	s := mcp_server.NewMCPServer(
		"pixiu-mcp-gateway",
		"1.0.0",
	)

	for i := range f.cfg.Tools {
		tc := f.cfg.Tools[i]

		var toolOptions []mcp.ToolOption
		toolOptions = append(toolOptions, mcp.WithDescription(tc.Description))

		for _, p := range tc.Parameters {
			var paramOptions []mcp.PropertyOption
			if p.Required {
				paramOptions = append(paramOptions, mcp.Required())
			}
			if p.Description != "" {
				paramOptions = append(paramOptions, mcp.Description(p.Description))
			}
			if len(p.Enum) > 0 {
				paramOptions = append(paramOptions, mcp.Enum(p.Enum...))
			}

			switch p.Type {
			case "string":
				toolOptions = append(toolOptions, mcp.WithString(p.Name, paramOptions...))
			case "number", "integer": // Fallback integer to number
				toolOptions = append(toolOptions, mcp.WithNumber(p.Name, paramOptions...))
			case "boolean":
				toolOptions = append(toolOptions, mcp.WithBoolean(p.Name, paramOptions...))
			default:
				return fmt.Errorf("unsupported parameter type for tool %s: %s", tc.Name, p.Type)
			}
		}

		toolSchema := mcp.NewTool(tc.Name, toolOptions...)

		noopHandler := func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			return mcp.NewToolResultError("this handler should not be called"), nil
		}
		s.AddTool(toolSchema, noopHandler)
	}

	return nil
}

// PrepareFilterChain adds the mcp_executor filter to the filter chain.
func (factory *FilterFactory) PrepareFilterChain(ctx *h.HttpContext, chain filter.FilterChain) error {
	f := &McpExecutorFilter{
		factory: factory,
	}
	chain.AppendDecodeFilters(f)
	chain.AppendEncodeFilters(f)
	return nil
}

// Decode prepares the downstream request based on the MCP tool call.
func (f *McpExecutorFilter) Decode(ctx *h.HttpContext) filter.FilterStatus {
	if !strings.HasSuffix(ctx.Request.URL.Path, "/tools/call") {
		return filter.Continue
	}

	bodyBytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		logger.Errorf("mcp_executor: failed to read request body: %v", err)
		ctx.SendLocalReply(http.StatusBadRequest, []byte("cannot read request body"))
		return filter.Stop
	}
	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	var reqBody callToolRequestBody
	if err := json.Unmarshal(bodyBytes, &reqBody); err != nil {
		logger.Errorf("mcp_executor: failed to unmarshal request body: %v", err)
		ctx.SendLocalReply(http.StatusBadRequest, []byte("invalid json format"))
		return filter.Stop
	}

	var toolConfig *Tool
	for i := range f.factory.cfg.Tools {
		if f.factory.cfg.Tools[i].Name == reqBody.ToolName {
			toolConfig = &f.factory.cfg.Tools[i]
			break
		}
	}

	if toolConfig == nil {
		msg := fmt.Sprintf("tool not found: %s", reqBody.ToolName)
		logger.Warnf("mcp_executor: %s", msg)
		ctx.SendLocalReply(http.StatusNotFound, []byte(msg))
		return filter.Stop
	}

	if ctx.Params == nil {
		ctx.Params = make(map[string]any)
	}
	ctx.Params["mcp_tool_config"] = toolConfig
	ctx.Params["mcp_arguments"] = reqBody.Arguments

	path := toolConfig.Request.PathTemplate
	for key, val := range reqBody.Arguments {
		if strVal, ok := val.(string); ok {
			path = strings.Replace(path, "{"+key+"}", strVal, -1)
		}
	}

	ctx.Request.Method = toolConfig.Request.Method
	ctx.Request.URL.Path = path

	if ctx.Route == nil && ctx.GetRouteEntry() != nil {
		ctx.Route = ctx.GetRouteEntry()
	}
	if ctx.Route != nil {
		ctx.Route.Cluster = toolConfig.Cluster
	} else {
		logger.Errorf("mcp_executor: no route entry found for tool %s", toolConfig.Name)
		ctx.SendLocalReply(http.StatusInternalServerError, []byte("internal server error: no route found"))
		return filter.Stop
	}

	return filter.Continue
}

// Encode transforms the downstream response into an MCP tool result.
func (f *McpExecutorFilter) Encode(ctx *h.HttpContext) filter.FilterStatus {
	if _, ok := ctx.Params["mcp_tool_config"]; !ok {
		return filter.Continue
	}

	var finalResult *mcp.CallToolResult

	if ctx.SourceResp == nil {
		finalResult = mcp.NewToolResultError("downstream request failed: no response from server")
	} else if resp, ok := ctx.SourceResp.(*http.Response); ok {
		if resp.StatusCode >= 400 {
			finalResult = mcp.NewToolResultError(fmt.Sprintf("downstream request failed with status code %d", resp.StatusCode))
		} else {
			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				finalResult = mcp.NewToolResultError(fmt.Sprintf("failed to read downstream response: %v", err))
			} else {
				finalResult = mcp.NewToolResultText(string(respBody))
			}
			_ = resp.Body.Close()
		}
	} else if unaryResp, ok := ctx.SourceResp.(*client.UnaryResponse); ok {
		finalResult = mcp.NewToolResultText(string(unaryResp.Data))
	} else {
		finalResult = mcp.NewToolResultError(fmt.Sprintf("downstream request failed: unexpected response type %T", ctx.SourceResp))
	}

	finalRespBytes, err := json.Marshal(finalResult)
	if err != nil {
		finalResult = mcp.NewToolResultError(fmt.Sprintf("failed to marshal final tool result: %v", err))
		finalRespBytes, _ = json.Marshal(finalResult)
	}

	ctx.TargetResp = &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(bytes.NewBuffer(finalRespBytes)),
	}

	return filter.Continue
}
