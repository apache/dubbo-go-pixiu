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

package router

import (
	"encoding/json"
	"math/rand"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

// maxDeniedSamples caps how many dropped tool names appear in a decision log to
// keep log lines bounded regardless of catalog size.
const maxDeniedSamples = 10

// DecisionLogger emits a structured, PII-safe record for each selection. It
// never logs prompt text or tool arguments; only counts, the mode, per-stage
// in/out tallies, and a bounded sample of denied tool names.
type DecisionLogger struct {
	sampleRate float64
}

// NewDecisionLogger builds a logger. A sampleRate <= 0 logs every decision; a
// rate >= 1 also logs every decision; values in between sample probabilistically.
func NewDecisionLogger(sampleRate float64) *DecisionLogger {
	return &DecisionLogger{sampleRate: sampleRate}
}

// decisionRecord is the JSON shape emitted to the log.
type decisionRecord struct {
	Event           string         `json:"event"`
	SessionID       string         `json:"session_id"`
	AgentID         string         `json:"agent_id,omitempty"`
	Tenant          string         `json:"tenant,omitempty"`
	Method          string         `json:"method"`
	MetadataVersion string         `json:"metadata_version"`
	Candidates      int            `json:"candidates"`
	Selected        int            `json:"selected"`
	Mode            string         `json:"mode"`
	Stages          map[string]int `json:"stages,omitempty"`
	DeniedSamples   []string       `json:"denied_samples,omitempty"`
}

// Log emits a decision record for the given selection, subject to sampling.
func (d *DecisionLogger) Log(sc SelectionContext, plan *SelectionPlan, candidates int) {
	if d == nil || plan == nil {
		return
	}
	if !d.shouldSample() {
		return
	}

	rec := decisionRecord{
		Event:           "mcp_router_decision",
		SessionID:       sc.SessionID,
		AgentID:         sc.AgentID,
		Tenant:          sc.Tenant,
		Method:          sc.Method,
		MetadataVersion: plan.Version,
		Candidates:      candidates,
		Selected:        len(plan.ToolNames),
		Mode:            plan.Mode,
		Stages:          stageDropCounts(plan.Reasons),
		DeniedSamples:   deniedSamples(plan.Reasons),
	}

	payload, err := json.Marshal(rec)
	if err != nil {
		return
	}
	logger.Infof("[dubbo-go-pixiu] %s", string(payload))
}

// shouldSample reports whether this decision should be logged.
func (d *DecisionLogger) shouldSample() bool {
	if d.sampleRate <= 0 || d.sampleRate >= 1 {
		return true
	}
	return rand.Float64() < d.sampleRate
}

// stageDropCounts tallies how many tools each stage dropped.
func stageDropCounts(traces []DecisionTrace) map[string]int {
	if len(traces) == 0 {
		return nil
	}
	counts := make(map[string]int)
	for _, t := range traces {
		if !t.Kept {
			counts[t.Stage]++
		}
	}
	if len(counts) == 0 {
		return nil
	}
	return counts
}

// deniedSamples returns up to maxDeniedSamples names of dropped tools. Tool
// names are configuration identifiers, not PII, so they are safe to log.
func deniedSamples(traces []DecisionTrace) []string {
	var out []string
	for _, t := range traces {
		if t.Kept {
			continue
		}
		out = append(out, t.Tool)
		if len(out) >= maxDeniedSamples {
			break
		}
	}
	return out
}
