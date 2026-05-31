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
	"context"
	"fmt"
	"testing"
	"time"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// benchTools generates n tools spread across 3 tenants, each tagged for policy.
func benchTools(n int) []model.ToolConfig {
	tenants := []string{"acme", "globex", "initech"}
	tools := make([]model.ToolConfig, n)
	for i := 0; i < n; i++ {
		tenant := tenants[i%len(tenants)]
		tools[i] = model.ToolConfig{
			Name:        fmt.Sprintf("%s_tool_%d", tenant, i),
			Description: fmt.Sprintf("tool %d for %s", i, tenant),
			Cluster:     "c",
			Meta:        &model.ToolMeta{Tags: []string{tenant}, Risk: "low"},
		}
	}
	return tools
}

func benchSelector(b *testing.B) (*CompositeSelector, *SessionPlanStore) {
	cfg := &model.RouterConfig{
		Enabled:  true,
		Fallback: FallbackFailClosed,
		Policy: model.PolicyConfig{Rules: []model.PolicyRule{
			{Name: "acme", When: model.PolicyMatch{Claim: "tenant", Equals: "acme"}, AllowTags: []string{"acme", "shared"}},
		}},
	}
	store := NewSessionPlanStoreWithTTL(time.Hour)
	sel, err := Build(cfg, store)
	if err != nil {
		b.Fatal(err)
	}
	return sel.(*CompositeSelector), store
}

// BenchmarkComposite_Policy_50tools measures a cold (uncached) selection over
// 50 tools with a single policy rule. Acceptance target: well under 3ms.
func BenchmarkComposite_Policy_50tools(b *testing.B) {
	cs, store := benchSelector(b)
	defer store.Stop()
	tools := benchTools(50)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Unique session id each iteration to force recomputation (cold path).
		sc := SelectionContext{SessionID: fmt.Sprintf("s-%d", i), Tenant: "acme"}
		_, _ = cs.Select(context.Background(), sc, tools)
	}
}

// BenchmarkComposite_Policy_1ktools measures the cold selection path at 1k tools.
func BenchmarkComposite_Policy_1ktools(b *testing.B) {
	cs, store := benchSelector(b)
	defer store.Stop()
	tools := benchTools(1000)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sc := SelectionContext{SessionID: fmt.Sprintf("s-%d", i), Tenant: "acme"}
		_, _ = cs.Select(context.Background(), sc, tools)
	}
}

// BenchmarkComposite_CachedReuse measures the warm path where a session's plan
// is reused without recomputation.
func BenchmarkComposite_CachedReuse(b *testing.B) {
	cs, store := benchSelector(b)
	defer store.Stop()
	tools := benchTools(1000)
	sc := SelectionContext{SessionID: "warm", Tenant: "acme"}
	_, _ = cs.Select(context.Background(), sc, tools)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = cs.Select(context.Background(), sc, tools)
	}
}
