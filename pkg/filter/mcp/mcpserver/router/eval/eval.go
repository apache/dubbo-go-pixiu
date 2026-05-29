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

// Package eval provides an offline evaluation harness for the MCP tool router's
// SchemaMatcher. Given a tool catalog and a labeled query set (query ->
// expected tool names), it computes Recall@K and Precision@K so config changes
// to schema weights or top_k can be validated without a running gateway.
package eval

import (
	"github.com/apache/dubbo-go-pixiu/pkg/filter/mcp/mcpserver/router"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// Case is a single labeled evaluation example.
type Case struct {
	Query    string   `json:"query" yaml:"query"`
	Expected []string `json:"expected_tools" yaml:"expected_tools"`
}

// Dataset bundles a tool catalog with labeled queries.
type Dataset struct {
	Tools []model.ToolConfig `json:"tools" yaml:"tools"`
	Cases []Case             `json:"cases" yaml:"cases"`
}

// Result holds aggregate metrics across all cases at a given K.
type Result struct {
	K             int     `json:"k"`
	Cases         int     `json:"cases"`
	MeanRecall    float64 `json:"mean_recall_at_k"`
	MeanPrecision float64 `json:"mean_precision_at_k"`
}

// Evaluate runs the SchemaMatcher over each case and aggregates Recall@K and
// Precision@K. The matcher is configured with the provided schema weights and
// a top_k of k. K must be > 0.
func Evaluate(ds Dataset, weights model.SchemaWeights, k int) Result {
	matcher := router.NewSchemaMatcher(model.SchemaConfig{Weights: weights, TopK: k})

	var sumRecall, sumPrecision float64
	for _, c := range ds.Cases {
		ranked, _ := matcher.Rank(ds.Tools, router.SelectionContext{UserPrompt: c.Query})
		topK := topKNames(ranked, k)
		recall, precision := score(topK, c.Expected)
		sumRecall += recall
		sumPrecision += precision
	}

	n := float64(len(ds.Cases))
	res := Result{K: k, Cases: len(ds.Cases)}
	if n > 0 {
		res.MeanRecall = sumRecall / n
		res.MeanPrecision = sumPrecision / n
	}
	return res
}

// topKNames extracts up to k tool names from a ranked slice.
func topKNames(ranked []model.ToolConfig, k int) []string {
	if k > len(ranked) {
		k = len(ranked)
	}
	names := make([]string, k)
	for i := 0; i < k; i++ {
		names[i] = ranked[i].Name
	}
	return names
}

// score computes recall and precision of predicted against expected.
//
//	recall    = |predicted ∩ expected| / |expected|
//	precision = |predicted ∩ expected| / |predicted|
func score(predicted, expected []string) (recall, precision float64) {
	if len(expected) == 0 {
		return 1, 1 // nothing expected: vacuously perfect
	}
	expSet := make(map[string]struct{}, len(expected))
	for _, e := range expected {
		expSet[e] = struct{}{}
	}
	hits := 0
	for _, p := range predicted {
		if _, ok := expSet[p]; ok {
			hits++
		}
	}
	recall = float64(hits) / float64(len(expected))
	if len(predicted) > 0 {
		precision = float64(hits) / float64(len(predicted))
	}
	return recall, precision
}
