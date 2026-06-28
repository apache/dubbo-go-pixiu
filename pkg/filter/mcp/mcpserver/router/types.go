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

// Package router implements deterministic MCP tool governance.
//
// It provides a ToolSelector that the MCP server filter invokes at tools/list
// and tools/call when the MCP server configuration contains a router block.
package router

import (
	"context"
	"errors"
	"hash/fnv"
	"strconv"
	"sync/atomic"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// ErrToolNotAuthorized is returned by AuthorizeCall when a tool is not part of
// the session's selection plan.
var ErrToolNotAuthorized = errors.New("tool not authorized for this session")

// ErrInvalidPlanKey reports a programmer/configuration error at the plan store
// boundary. A plan key must identify both the router instance and MCP session.
var ErrInvalidPlanKey = errors.New("invalid mcp router session plan key")

// ErrInvalidSelectionPlan reports a nil or self-inconsistent plan passed to the
// store. The store never treats invalid input as a successful write.
var ErrInvalidSelectionPlan = errors.New("invalid mcp router selection plan")

// ErrSessionPlanNotActive reports that the transport session generation that
// produced the plan is no longer the active generation for that session ID.
var ErrSessionPlanNotActive = errors.New("mcp router session plan is not active")

// ErrSessionPlanStale reports that a conditional plan commit observed a newer
// plan than the one it started from.
var ErrSessionPlanStale = errors.New("mcp router session plan is stale")

// Selection mode labels used in plans, logs and metrics.
const (
	ModeSelected       = "selected"
	ModeFallbackBundle = "fallback_bundle"
	ModeFailClosed     = "fail_closed"
)

// Decision stage labels.
const (
	StageInternal    = "internal"
	StagePolicy      = "policy"
	StageWorkflow    = "workflow"
	StageProgressive = "progressive"
)

// maxDecisionTraceSamples bounds in-memory drop samples per stage.
const maxDecisionTraceSamples = 32

// Selection outcome labels describe why the final selected set has its shape.
const (
	SelectionOutcomeSelected             = "selected"
	SelectionOutcomeNoMatch              = "no_match"
	SelectionOutcomeExplicitDeny         = "explicit_deny"
	SelectionOutcomeEmptyByConfiguration = "empty_by_configuration"
	SelectionOutcomeInternalError        = "internal_error"
)

// SelectionContext is the input to a selection decision, extracted from the
// MCP context and request body. Fields are best-effort; an empty field simply
// means the corresponding signal was unavailable.
type SelectionContext struct {
	SessionID         string         // Mcp-Session-Id
	SessionGeneration uint64         // transport session generation validated for this request
	Method            string         // "tools/list" | "tools/call"
	UserID            string         // from claims.sub
	Tenant            string         // from claims.tenant
	Claims            map[string]any // JWT claims already validated by the auth/mcp filter
	Requested         string         // target tool name on tools/call
	CatalogVersion    string         // immutable registry snapshot version for this request
}

// StageCount records bounded per-stage cardinality without retaining one trace
// per candidate tool.
type StageCount struct {
	Input  int `json:"input"`
	Output int `json:"output"`
}

// SelectionPlan is the final result of one selection, bound to a session.
type SelectionPlan struct {
	SessionID          string                `json:"session_id"`
	ToolNames          []string              `json:"tool_names"`                   // authorized tool names, stable order
	VisibleToolNames   []string              `json:"visible_tool_names,omitempty"` // tools/list names; nil means same as ToolNames
	VisibleFingerprint string                `json:"-"`                            // tools/list visible definition fingerprint
	Mode               string                `json:"mode"`                         // ModeSelected | ModeFallbackBundle | ModeFailClosed
	Outcome            string                `json:"outcome,omitempty"`            // SelectionOutcome*
	StageCounts        map[string]StageCount `json:"stage_counts,omitempty"`
	Reasons            []DecisionTrace       `json:"reasons,omitempty"`  // bounded dropped-tool samples
	Version            string                `json:"version"`            // metadata snapshot version (registry + config hash)
	CreatedAt          int64                 `json:"created_at"`         // unix nano
	IdentityHash       string                `json:"-"`                  // validated-claims fingerprint, never logged
	ProgressiveHash    string                `json:"-"`                  // progressive config fingerprint
	ConfigHash         string                `json:"-"`                  // normalized authorization config fingerprint
	CatalogVersion     string                `json:"-"`                  // immutable catalog version used for this plan
	Generation         uint64                `json:"-"`                  // monotonic per-session plan generation
	Expanded           bool                  `json:"expanded,omitempty"` // progressive state for tests/log-free inspection
	toolSet            map[string]struct{}   `json:"-"`
}

// VisibleFingerprintFunc hashes the client-visible tools/list definition for a
// selected tool set. The MCP server layer injects the implementation that uses
// the exact mcp-go response projection; router keeps a small default for tests
// and non-server callers without importing the protocol package.
type VisibleFingerprintFunc func([]model.ToolConfig) string

func DefaultVisibleFingerprint(tools []model.ToolConfig) string {
	return VisibleToolNamesFingerprint(visibleToolNames(tools))
}

func VisibleToolNamesFingerprint(names []string) string {
	if len(names) == 0 {
		return "00000000"
	}
	h := fnv.New64a()
	for _, name := range names {
		_, _ = h.Write([]byte(name))
		_, _ = h.Write([]byte{0})
	}
	return strconv.FormatUint(h.Sum64(), 16)
}

// Contains reports whether the plan authorizes the named tool.
func (p *SelectionPlan) Contains(tool string) bool {
	if p == nil {
		return false
	}
	if p.toolSet != nil {
		_, ok := p.toolSet[tool]
		return ok
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

// DecisionTrace records why a single candidate was dropped at a stage.
// It must never contain PII (no prompt text, no argument values).
type DecisionTrace struct {
	Tool   string `json:"tool,omitempty"`
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
	// catalog and returns a receipt bound to the authorized plan generation.
	AuthorizeCall(ctx context.Context, sc SelectionContext, candidates []model.ToolConfig) (*AuthorizationReceipt, error)
}

// SelectionFailureHandler is optionally implemented by selectors that can
// produce a safe fallback plan after Select returns an error.
type SelectionFailureHandler interface {
	HandleSelectionFailure(ctx context.Context, sc SelectionContext, candidates []model.ToolConfig, cause error) (*SelectionPlan, error)
}

// PlanRefreshSelector is implemented by selectors that can refresh an existing
// session plan with an atomic stale-generation guard.
type PlanRefreshSelector interface {
	RefreshPlan(ctx context.Context, base SessionPlanContext, candidates []model.ToolConfig) (*SelectionPlan, bool, error)
}

// CallSuccessResult describes the progressive-disclosure effect of one
// successful tools/call.
type CallSuccessResult struct {
	Count        int64
	Transitioned bool
}

type ReceiptOutcome int

const (
	ReceiptSucceeded ReceiptOutcome = iota
	ReceiptAborted
)

// CallSuccessRecorder is implemented by selectors that track successful
// tools/call completions separately from authorization checks.
type CallSuccessRecorder interface {
	RecordCallSuccess(ctx context.Context, receipt AuthorizationReceipt) (CallSuccessResult, error)
}

// ReceiptFinalizer is implemented by selectors that can release an issued
// receipt on both successful and failed tools/call completion paths.
type ReceiptFinalizer interface {
	FinalizeReceipt(ctx context.Context, receipt AuthorizationReceipt, outcome ReceiptOutcome) (CallSuccessResult, error)
}

// AuthorizationReceipt binds one tools/call authorization to the plan and
// metadata generation observed at authorization time. It is safe to persist in
// request context, but must not be logged because it contains identity hashes.
type AuthorizationReceipt struct {
	RouterInstanceID string
	SessionID        string
	ToolName         string
	PlanGeneration   uint64
	IdentityHash     string
	ConfigHash       string
	CatalogVersion   string
	ProgressiveHash  string
	ReceiptID        uint64
	state            *receiptState
}

type receiptState struct {
	consumed atomic.Bool
}

func (s *receiptState) consume() bool {
	return s != nil && s.consumed.CompareAndSwap(false, true)
}

// toolNames extracts the stable ordered name slice from a candidate set.
func toolNames(tools []model.ToolConfig) []string {
	names := make([]string, len(tools))
	for i, t := range tools {
		names[i] = t.Name
	}
	return names
}

func toolNameSet(names []string) map[string]struct{} {
	if len(names) == 0 {
		return nil
	}
	set := make(map[string]struct{}, len(names))
	for _, name := range names {
		set[name] = struct{}{}
	}
	return set
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
