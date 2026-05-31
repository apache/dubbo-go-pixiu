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
	"net"
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

var routerAdminNotFoundBody = []byte("not found")

// buildSelectionContext assembles the router input from the MCP request context.
// It extracts session, method, the requested tool name (for tools/call), and the
// JWT claims propagated by the auth/mcp filter (with sub/tenant promoted for
// convenient policy access).
func (f *MCPServerFilter) buildSelectionContext(ctx *MCPContext, method, requested string) router.SelectionContext {
	sc := router.SelectionContext{
		SessionID: ctx.SessionID(),
		Method:    method,
		Requested: requested,
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

// isRouterAdminRequest reports whether the request targets the router plan
// inspection endpoint.
func (f *MCPServerFilter) isRouterAdminRequest(ctx *contexthttp.HttpContext) bool {
	return ctx.Request != nil && ctx.Request.URL != nil &&
		strings.HasPrefix(ctx.Request.URL.Path, routerAdminPathPrefix)
}

// handleRouterAdmin serves GET /__mcp/router/plan/{session_id}. It is gated by
// audit.payload_logging (off by default) and restricted to loopback clients only,
// so the endpoint is not reachable from real clients even when logging is enabled.
// In proxy/sidecar deployments where RemoteAddr is the proxy's loopback address,
// additional routing-layer restrictions should be applied.
func (f *MCPServerFilter) handleRouterAdmin(ctx *contexthttp.HttpContext) filter.FilterStatus {
	if !f.routerAuditEnabled() {
		ctx.SendLocalReply(http.StatusNotFound, routerAdminNotFoundBody)
		return filter.Stop
	}

	// Restrict to loopback clients only. This prevents accidental exposure when
	// payload_logging is enabled, while preserving local debugging and kubectl
	// port-forward use cases. We trust only the TCP peer address (RemoteAddr),
	// never X-Forwarded-For, to avoid trivial bypass.
	if !isLoopback(ctx.Request.RemoteAddr) {
		ctx.SendLocalReply(http.StatusNotFound, routerAdminNotFoundBody)
		return filter.Stop
	}

	if ctx.Request.Method != constant.Get {
		ctx.Writer.Header().Set("Allow", "GET")
		ctx.SendLocalReply(http.StatusMethodNotAllowed, []byte("method not allowed"))
		return filter.Stop
	}

	inspector, ok := f.selector.(router.PlanInspector)
	if !ok {
		ctx.SendLocalReply(http.StatusNotFound, routerAdminNotFoundBody)
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

// isLoopback reports whether the remote address is a loopback (localhost) peer.
// It parses the host portion of "host:port" and checks for IPv4 127.0.0.0/8,
// IPv6 ::1, or the literal "localhost". Returns false on parse errors to fail
// closed (deny non-parseable addresses).
func isLoopback(remoteAddr string) bool {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		// RemoteAddr should always be "host:port", but if parsing fails treat
		// it as non-loopback to fail closed.
		return false
	}
	if host == "localhost" {
		return true
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return false
	}
	return ip.IsLoopback()
}
