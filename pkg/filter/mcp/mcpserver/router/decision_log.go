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
	cryptorand "crypto/rand"
	"encoding/binary"
	"encoding/json"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

// maxDeniedSamples caps dropped tool names in a decision log.
const maxDeniedSamples = 10

// DecisionLogger emits a structured, PII-safe record for each selection. It
// never logs session IDs, identity claims, prompt text, tool arguments, or
// request text details. Only counts, the mode, and per-stage drop tallies are
// emitted by default. Denied tool samples require explicit decision detail logging.
type DecisionLogger struct {
	sampleRate    float64
	detailLogging bool
}

// NewDecisionLogger builds a logger. A sampleRate of 0 disables decision logs;
// values in (0,1] sample probabilistically/all. Detailed denied samples are
// emitted only when decision detail logging is explicitly active.
func NewDecisionLogger(sampleRate float64, detailLogging bool) *DecisionLogger {
	return &DecisionLogger{sampleRate: sampleRate, detailLogging: detailLogging}
}

// decisionRecord is the JSON shape emitted to the log.
type decisionRecord struct {
	Event          string         `json:"event"`
	Method         string         `json:"method"`
	CatalogVersion string         `json:"catalog_version,omitempty"`
	ConfigVersion  string         `json:"config_version,omitempty"`
	Candidates     int            `json:"candidates"`
	Selected       int            `json:"selected"`
	Mode           string         `json:"mode"`
	Stages         map[string]int `json:"stages,omitempty"`
	DeniedSamples  []string       `json:"denied_samples,omitempty"`
}

// Log emits a decision record for the given selection, subject to sampling.
func (d *DecisionLogger) Log(sc SelectionContext, plan *SelectionPlan, candidates int) {
	if d == nil || plan == nil {
		return
	}
	if !d.shouldSample() {
		return
	}

	rec := d.record(sc, plan, candidates)

	data, err := json.Marshal(rec)
	if err != nil {
		return
	}
	logger.Infof("[dubbo-go-pixiu] %s", string(data))
}

func (d *DecisionLogger) record(sc SelectionContext, plan *SelectionPlan, candidates int) decisionRecord {
	rec := decisionRecord{
		Event:          "mcp_router_decision",
		Method:         sc.Method,
		CatalogVersion: plan.CatalogVersion,
		ConfigVersion:  plan.ConfigHash,
		Candidates:     candidates,
		Selected:       len(plan.ToolNames),
		Mode:           plan.Mode,
		Stages:         stageDropCountsFromPlan(plan),
	}
	if d.detailLogging {
		rec.DeniedSamples = deniedSamples(plan.Reasons)
	}
	return rec
}

// shouldSample reports whether this decision should be logged.
func (d *DecisionLogger) shouldSample() bool {
	if d.sampleRate <= 0 {
		return false
	}
	if d.sampleRate >= 1 {
		return true
	}
	sample, ok := secureRandomFloat64()
	return ok && sample < d.sampleRate
}

func secureRandomFloat64() (float64, bool) {
	var b [8]byte
	if _, err := cryptorand.Read(b[:]); err != nil {
		return 0, false
	}
	// Keep the top 53 bits so the value fits exactly in a float64 mantissa.
	n := binary.BigEndian.Uint64(b[:]) >> 11
	return float64(n) / (1 << 53), true
}

// stageDropCounts tallies how many tools each stage dropped.
func stageDropCounts(traces []DecisionTrace) map[string]int {
	if len(traces) == 0 {
		return nil
	}
	counts := make(map[string]int)
	for _, t := range traces {
		counts[t.Stage]++
	}
	if len(counts) == 0 {
		return nil
	}
	return counts
}

func stageDropCountsFromPlan(plan *SelectionPlan) map[string]int {
	if plan == nil {
		return nil
	}
	if len(plan.StageCounts) == 0 {
		return stageDropCounts(plan.Reasons)
	}
	counts := make(map[string]int, len(plan.StageCounts))
	for stage, count := range plan.StageCounts {
		dropped := count.Input - count.Output
		if dropped > 0 {
			counts[stage] = dropped
		}
	}
	if len(counts) == 0 {
		return nil
	}
	return counts
}

// deniedSamples returns up to maxDeniedSamples names of dropped tools. Tool
// names can disclose governance intent, so callers must only include them when
// decision detail logging is explicitly active.
func deniedSamples(traces []DecisionTrace) []string {
	var out []string
	for _, t := range traces {
		out = append(out, t.Tool)
		if len(out) >= maxDeniedSamples {
			break
		}
	}
	return out
}
