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
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/filter/mcp/mcpserver/router"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// buildSelectionContext assembles the router input from the MCP request context.
// It extracts session, method, the requested tool name (for tools/call), and the
// JWT claims propagated by the auth/mcp filter (with sub/tenant promoted for
// convenient policy access).
func (f *MCPServerFilter) buildSelectionContext(ctx *MCPContext, method, requested string) router.SelectionContext {
	sc := router.SelectionContext{
		SessionID:         ctx.SessionID(),
		SessionGeneration: ctx.SessionGeneration(),
		Method:            method,
		Requested:         requested,
	}

	if claims := mcpAuthClaims(ctx); claims != nil {
		sc.Claims = claims
		sc.UserID = claimStr(claims, "sub")
		sc.Tenant = claimStr(claims, "tenant")
	}

	return sc
}

func (f *MCPServerFilter) buildSelectionContextWithCatalog(ctx *MCPContext, method, requested, catalogVersion string) router.SelectionContext {
	sc := f.buildSelectionContext(ctx, method, requested)
	sc.CatalogVersion = catalogVersion
	return sc
}

// mcpAuthClaims returns the JWT claims propagated by the auth/mcp filter, or
// nil when no claims are present (e.g. auth filter not in the chain).
func mcpAuthClaims(ctx *MCPContext) map[string]any {
	if ctx.Params == nil {
		return nil
	}
	claims, _ := ctx.Params[constant.MCPAuthClaimsParamKey].(map[string]any)
	return claims
}

// claimStr resolves a claim to its string form, returning "" when absent.
func claimStr(claims map[string]any, key string) string {
	if v, ok := claims[key]; ok {
		return fmt.Sprintf("%v", v)
	}
	return ""
}

// filterByPlan returns the tools/list subset whose names appear in the plan's
// visible view, preserving plan order for stable client behavior.
func filterByPlan(toolCfgs []model.ToolConfig, plan *router.SelectionPlan) []model.ToolConfig {
	if plan == nil {
		return toolCfgs
	}
	byName := make(map[string]model.ToolConfig, len(toolCfgs))
	for _, t := range toolCfgs {
		byName[t.Name] = t
	}
	visibleNames := plan.VisibleNames()
	out := make([]model.ToolConfig, 0, len(visibleNames))
	for _, name := range visibleNames {
		if t, ok := byName[name]; ok {
			out = append(out, t)
		}
	}
	return out
}

func (f *MCPServerFilter) handleSelectionFailure(ctx *MCPContext, sc router.SelectionContext, candidates []model.ToolConfig, cause error) ([]model.ToolConfig, error) {
	if handler, ok := f.selector.(router.SelectionFailureHandler); ok {
		plan, err := handler.HandleSelectionFailure(ctx.Ctx, sc, candidates, cause)
		if err != nil {
			return nil, err
		}
		return filterByPlan(candidates, plan), nil
	}
	return nil, cause
}
