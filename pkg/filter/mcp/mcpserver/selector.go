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
	"fmt"
	"net/http"
	"strings"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/filter/mcp/mcpserver/router"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// routerAdminPathPrefix is the base path for the router plan inspection endpoint.
const routerAdminPathPrefix = "/__mcp/router/plan/"

// buildSelectionContext assembles the router input from the MCP request
// context. It extracts session, method, the requested tool name (for
// tools/call), a sanitized header snapshot, and the JWT claims propagated by
// the auth/mcp filter (with sub/tenant promoted for convenient policy access).
func (f *MCPServerFilter) buildSelectionContext(ctx *MCPContext, method, requested string) router.SelectionContext {
	sc := router.SelectionContext{
		SessionID: ctx.SessionID(),
		Method:    method,
		Requested: requested,
		Headers:   sanitizeHeaders(ctx),
	}

	if claims := mcpAuthClaims(ctx); claims != nil {
		sc.Claims = claims
		sc.UserID = claimStr(claims, "sub")
		sc.Tenant = claimStr(claims, "tenant")
	}

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

// sensitiveHeaders are never copied into the SelectionContext to avoid leaking
// credentials into decision logs or plan state.
var sensitiveHeaders = map[string]struct{}{
	"Authorization": {},
	"Cookie":        {},
	"Set-Cookie":    {},
}

// sanitizeHeaders returns a copy of the request headers with sensitive entries
// removed. It returns nil when there are no headers to copy.
func sanitizeHeaders(ctx *MCPContext) map[string]string {
	if ctx.Request == nil || len(ctx.Request.Header) == 0 {
		return nil
	}
	out := make(map[string]string, len(ctx.Request.Header))
	for k, v := range ctx.Request.Header {
		if _, blocked := sensitiveHeaders[k]; blocked {
			continue
		}
		if len(v) > 0 {
			out[k] = v[0]
		}
	}
	return out
}

// filterByPlan returns the subset of toolCfgs whose names appear in the plan,
// preserving the plan's ordering for stable client behavior.
func filterByPlan(toolCfgs []model.ToolConfig, plan *router.SelectionPlan) []model.ToolConfig {
	if plan == nil {
		return toolCfgs
	}
	byName := make(map[string]model.ToolConfig, len(toolCfgs))
	for _, t := range toolCfgs {
		byName[t.Name] = t
	}
	out := make([]model.ToolConfig, 0, len(plan.ToolNames))
	for _, name := range plan.ToolNames {
		if t, ok := byName[name]; ok {
			out = append(out, t)
		}
	}
	return out
}

// isRouterAdminRequest reports whether the request targets the router plan
// inspection endpoint.
func (f *MCPServerFilter) isRouterAdminRequest(ctx *contexthttp.HttpContext) bool {
	return ctx.Request != nil && ctx.Request.URL != nil &&
		strings.HasPrefix(ctx.Request.URL.Path, routerAdminPathPrefix)
}

// handleRouterAdmin serves GET /__mcp/router/plan/{session_id}. It is gated by
// audit.payload_logging: when that is off (or the router is disabled) it returns
// 404 so the endpoint's existence is not observable in production by default.
func (f *MCPServerFilter) handleRouterAdmin(ctx *contexthttp.HttpContext) filter.FilterStatus {
	if !f.routerAuditEnabled() {
		ctx.SendLocalReply(http.StatusNotFound, []byte("not found"))
		return filter.Stop
	}

	if ctx.Request.Method != constant.Get {
		ctx.Writer.Header().Set("Allow", "GET")
		ctx.SendLocalReply(http.StatusMethodNotAllowed, []byte("method not allowed"))
		return filter.Stop
	}

	inspector, ok := f.selector.(router.PlanInspector)
	if !ok {
		ctx.SendLocalReply(http.StatusNotFound, []byte("not found"))
		return filter.Stop
	}

	sessionID := strings.TrimPrefix(ctx.Request.URL.Path, routerAdminPathPrefix)
	if sessionID == "" {
		ctx.SendLocalReply(http.StatusBadRequest, []byte("missing session id"))
		return filter.Stop
	}

	plan, found := inspector.InspectPlan(sessionID)
	if !found {
		ctx.SendLocalReply(http.StatusNotFound, []byte(fmt.Sprintf("no plan for session %s", sessionID)))
		return filter.Stop
	}

	body, err := json.Marshal(plan)
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp router admin failed to marshal plan: %v", err)
		ctx.SendLocalReply(http.StatusInternalServerError, []byte("internal error"))
		return filter.Stop
	}

	ctx.Writer.Header().Set(constant.HeaderKeyContextType, constant.HeaderValueApplicationJson)
	ctx.SendLocalReply(http.StatusOK, body)
	return filter.Stop
}

// routerAuditEnabled reports whether the router is active with payload logging
// opted in, which is the precondition for exposing plan internals.
func (f *MCPServerFilter) routerAuditEnabled() bool {
	return f.selector != nil &&
		f.cfg.Router != nil &&
		f.cfg.Router.Enabled &&
		f.cfg.Router.Audit.PayloadLogging
}
