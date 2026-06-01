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

// Package router implements intelligent MCP tool routing for issue #937.
//
// It provides a pluggable ToolSelector that the MCP server filter invokes at
// three hookpoints: OnInitialize (session metadata capture), Select (tools/list
// trimming), and AuthorizeCall (tools/call enforcement). When no router is
// configured the filter keeps a nil selector and behaves exactly as before.
package router

import (
	"context"
	"errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// ErrToolNotAuthorized is returned by AuthorizeCall when a tool is not part of
// the session's selection plan and enforce_on_call is enabled.
var ErrToolNotAuthorized = errors.New("tool not authorized for this session")

// Selection mode labels used in plans, logs and metrics.
const (
	ModeHybrid         = "hybrid"
	ModeFallbackBundle = "fallback_bundle"
	ModeFailClosed     = "fail_closed"
)

// Decision stage labels.
const (
	StagePolicy      = "policy"
	StageWorkflow    = "workflow"
	StageProgressive = "progressive"
)

// SelectionContext is the input to a selection decision, extracted from the
// MCP context and request body. Fields are best-effort; an empty field simply
// means the corresponding signal was unavailable.
type SelectionContext struct {
	SessionID string         // Mcp-Session-Id
	Method    string         // "initialize" | "tools/list" | "tools/call"
	AgentID   string         // from initialize.clientInfo.name or header
	UserID    string         // from claims.sub
	Tenant    string         // from claims.tenant
	Claims    map[string]any // JWT claims already validated by the auth/mcp filter
	Requested string         // target tool name on tools/call
}

// SelectionPlan is the final result of one selection, bound to a session.
type SelectionPlan struct {
	SessionID        string          `json:"session_id"`
	ToolNames        []string        `json:"tool_names"`                   // authorized tool names, stable order
	VisibleToolNames []string        `json:"visible_tool_names,omitempty"` // tools/list names; nil means same as ToolNames
	Mode             string          `json:"mode"`                         // ModeHybrid | ModeFallbackBundle | ModeFailClosed
	Reasons          []DecisionTrace `json:"reasons,omitempty"`            // per-candidate keep/drop reasons, for audit
	Version          string          `json:"version"`                      // metadata snapshot version (registry + config hash)
	CreatedAt        int64           `json:"created_at"`                   // unix nano
	ExpiresAt        int64           `json:"expires_at,omitempty"`         // 0 = session-lifetime valid
}

// Contains reports whether the plan authorizes the named tool.
func (p *SelectionPlan) Contains(tool string) bool {
	if p == nil {
		return false
	}
	for _, t := range p.ToolNames {
		if t == tool {
			return true
		}
	}
	return false
}

// VisibleNames returns the stable tools/list view for this plan. A nil
// VisibleToolNames preserves compatibility with older in-memory plans/tests.
func (p *SelectionPlan) VisibleNames() []string {
	if p == nil {
		return nil
	}
	if p.VisibleToolNames != nil {
		return p.VisibleToolNames
	}
	return p.ToolNames
}

// DecisionTrace records why a single candidate was kept or dropped at a stage.
// It must never contain PII (no prompt text, no argument values).
type DecisionTrace struct {
	Tool   string `json:"tool,omitempty"`
	Kept   bool   `json:"kept"`
	Stage  string `json:"stage,omitempty"`  // StagePolicy | StageWorkflow | StageProgressive
	Rule   string `json:"rule,omitempty"`   // matched rule / workflow name
	Detail string `json:"detail,omitempty"` // short, non-PII explanation
}

// ToolSelector is the unified entry point invoked by the MCP server filter.
// Implementations may be composite selectors or test mocks.
type ToolSelector interface {
	// Select produces (or reuses) a SelectionPlan at tools/list time. The
	// implementation is responsible for consulting the SessionPlanStore and
	// recomputing when the plan is missing or stale.
	Select(ctx context.Context, sc SelectionContext, candidates []model.ToolConfig) (*SelectionPlan, error)

	// AuthorizeCall enforces the plan at tools/call time against the current
	// candidate catalog. A nil return allows the call. When enforce_on_call is
	// disabled the implementation may allow calls even without a session plan.
	AuthorizeCall(ctx context.Context, sc SelectionContext, candidates []model.ToolConfig) error

	// OnInitialize gives the selector a chance to capture session metadata.
	// Implementations may treat this as a no-op.
	OnInitialize(ctx context.Context, sc SelectionContext, candidates []model.ToolConfig) error

	// Name returns the implementation name for logs and metrics.
	Name() string
}

// CallSuccessRecorder is implemented by selectors that track successful
// tools/call completions separately from authorization checks.
type CallSuccessRecorder interface {
	RecordCallSuccess(ctx context.Context, sc SelectionContext) error
}

// PlanInspector is optionally implemented by selectors that can surface a
// session's current plan for the admin debug endpoint.
type PlanInspector interface {
	InspectPlan(sessionID string) (*SelectionPlan, bool)
}

// toolNames extracts the stable ordered name slice from a candidate set.
func toolNames(tools []model.ToolConfig) []string {
	names := make([]string, len(tools))
	for i, t := range tools {
		names[i] = t.Name
	}
	return names
}

// visibleToolNames extracts the tools/list view from a candidate set. Tools
// with meta.discovery_visibility=false remain authorized in the plan but are
// intentionally omitted from discovery.
func visibleToolNames(tools []model.ToolConfig) []string {
	names := make([]string, 0, len(tools))
	for _, t := range tools {
		if t.Meta != nil && t.Meta.DiscoveryVisibility != nil && !*t.Meta.DiscoveryVisibility {
			continue
		}
		names = append(names, t.Name)
	}
	return names
}
