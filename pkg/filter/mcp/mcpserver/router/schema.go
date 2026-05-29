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
	"sort"
	"strings"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// Default schema scoring weights, used when a weight is left at zero.
const (
	defaultTagWeight         = 2.0
	defaultCapabilityWeight  = 3.0
	defaultDescriptionWeight = 1.0
)

// SchemaMatcher ranks tools by lexical relevance to the user prompt and
// optionally truncates to the top K. It only reorders and trims by rank; it
// never drops a tool for authorization reasons (that is policy/workflow's job).
type SchemaMatcher struct {
	tagW  float64
	capW  float64
	descW float64
	topK  int
}

// NewSchemaMatcher builds a matcher from config, applying sensible defaults.
func NewSchemaMatcher(cfg model.SchemaConfig) *SchemaMatcher {
	m := &SchemaMatcher{
		tagW:  cfg.Weights.TagMatch,
		capW:  cfg.Weights.CapabilityMatch,
		descW: cfg.Weights.DescriptionMatch,
		topK:  cfg.TopK,
	}
	if m.tagW == 0 {
		m.tagW = defaultTagWeight
	}
	if m.capW == 0 {
		m.capW = defaultCapabilityWeight
	}
	if m.descW == 0 {
		m.descW = defaultDescriptionWeight
	}
	return m
}

// Rank scores each tool against the prompt tokens, sorts descending by score
// (stable on ties to preserve registry order), and truncates to topK.
//
// When there is no user prompt the input order is preserved and only topK
// truncation applies, so schema ranking is safe to enable without prompts.
func (m *SchemaMatcher) Rank(tools []model.ToolConfig, sc SelectionContext) ([]model.ToolConfig, []DecisionTrace) {
	tokens := tokenize(sc.UserPrompt)

	if len(tokens) > 0 {
		type scored struct {
			tool  model.ToolConfig
			score float64
			idx   int
		}
		items := make([]scored, len(tools))
		for i, t := range tools {
			items[i] = scored{tool: t, score: m.score(t, tokens), idx: i}
		}
		sort.SliceStable(items, func(a, b int) bool {
			if items[a].score != items[b].score {
				return items[a].score > items[b].score
			}
			return items[a].idx < items[b].idx
		})
		tools = make([]model.ToolConfig, len(items))
		for i, it := range items {
			tools[i] = it.tool
		}
	}

	if m.topK > 0 && len(tools) > m.topK {
		dropped := tools[m.topK:]
		traces := make([]DecisionTrace, 0, len(dropped))
		for _, t := range dropped {
			traces = append(traces, DecisionTrace{
				Tool:   t.Name,
				Kept:   false,
				Stage:  StageSchema,
				Detail: "below_top_k",
			})
		}
		return tools[:m.topK], traces
	}

	return tools, nil
}

// score computes a tool's weighted relevance to the prompt tokens.
func (m *SchemaMatcher) score(tool model.ToolConfig, tokens map[string]struct{}) float64 {
	var s float64
	if tool.Meta != nil {
		for _, tag := range tool.Meta.Tags {
			if _, ok := tokens[strings.ToLower(tag)]; ok {
				s += m.tagW
			}
		}
		for _, cap := range tool.Meta.Capabilities {
			if matchAnyToken(cap, tokens) {
				s += m.capW
			}
		}
	}
	for tok := range tokenize(tool.Description) {
		if _, ok := tokens[tok]; ok {
			s += m.descW
		}
	}
	return s
}

// matchAnyToken reports whether any dot/underscore-separated part of a
// capability identifier appears in the token set (e.g. "user.read" matches
// prompts mentioning "user" or "read").
func matchAnyToken(capability string, tokens map[string]struct{}) bool {
	for _, part := range strings.FieldsFunc(strings.ToLower(capability), func(r rune) bool {
		return r == '.' || r == '_' || r == '-'
	}) {
		if _, ok := tokens[part]; ok {
			return true
		}
	}
	return false
}

// tokenize lowercases and splits text into a set of word tokens.
func tokenize(text string) map[string]struct{} {
	if text == "" {
		return nil
	}
	fields := strings.FieldsFunc(strings.ToLower(text), func(r rune) bool {
		return !(r >= 'a' && r <= 'z') && !(r >= '0' && r <= '9')
	})
	if len(fields) == 0 {
		return nil
	}
	set := make(map[string]struct{}, len(fields))
	for _, f := range fields {
		set[f] = struct{}{}
	}
	return set
}
