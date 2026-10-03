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
	"math"
	"strconv"
	"testing"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

var (
	benchPlan    *SelectionPlan
	benchReceipt *AuthorizationReceipt
	benchResult  CallSuccessResult
	benchOK      bool
	benchErr     error
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

func benchSelector(b *testing.B, store *SessionPlanStore) *CompositeSelector {
	cfg := &model.RouterConfig{
		Fallback: FallbackFailClosed,
		Policy: model.PolicyConfig{Rules: []model.PolicyRule{
			{Name: "acme", When: model.PolicyMatch{Claim: "tenant", Equals: "acme"}, AllowTags: []string{"acme", "shared"}},
		}},
	}
	sel, err := Build(cfg, store)
	if err != nil {
		b.Fatal(err)
	}
	return sel.(*CompositeSelector)
}

func benchProgressiveSelector(b *testing.B, store *SessionPlanStore, tools []model.ToolConfig) *CompositeSelector {
	acmeTools := make([]string, 0, len(tools)/3+1)
	for _, tool := range tools {
		if tool.Meta != nil && len(tool.Meta.Tags) > 0 && tool.Meta.Tags[0] == "acme" {
			acmeTools = append(acmeTools, tool.Name)
		}
	}
	cfg := &model.RouterConfig{
		Fallback: FallbackFailClosed,
		Policy: model.PolicyConfig{Rules: []model.PolicyRule{
			{Name: "acme", When: model.PolicyMatch{Claim: "tenant", Equals: "acme"}, AllowTags: []string{"acme", "shared"}},
		}},
		Stages: model.RouterStages{Progressive: true},
		Workflows: []model.WorkflowConfig{
			{Name: "initial", Tools: acmeTools},
		},
		Progressive: model.ProgressiveConfig{
			InitialBundle:    "initial",
			ExpandAfterCalls: math.MaxInt32,
		},
	}
	sel, err := Build(cfg, store)
	if err != nil {
		b.Fatal(err)
	}
	return sel.(*CompositeSelector)
}

func BenchmarkCompositeSelector(b *testing.B) {
	for _, size := range []int{50, 1000, 10000} {
		b.Run("warm_"+strconv.Itoa(size), func(b *testing.B) {
			benchmarkCompositeSelectorWarm(b, size)
		})

		b.Run("call_"+strconv.Itoa(size), func(b *testing.B) {
			benchmarkCompositeSelectorCall(b, size)
		})

		b.Run("cold_"+strconv.Itoa(size), func(b *testing.B) {
			benchmarkCompositeSelectorCold(b, size)
		})
	}
}

func benchmarkCompositeSelectorWarm(b *testing.B, size int) {
	store := NewSessionPlanStore()
	defer store.Stop()
	cs := benchSelector(b, store)
	tools := benchTools(size)
	sc := SelectionContext{SessionID: "warm", Tenant: "acme", CatalogVersion: "bench-" + strconv.Itoa(size)}
	plan, err := cs.Select(context.Background(), sc, tools)
	if err != nil || plan == nil {
		b.Fatalf("prewarm failed: plan=%v err=%v", plan, err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchPlan, benchErr = cs.Select(context.Background(), sc, tools)
	}
	b.StopTimer()
	requireBenchPlan(b, "warm", benchPlan, benchErr)
	requireBenchStoreSize(b, "warm", store, 1)
}

func benchmarkCompositeSelectorCall(b *testing.B, size int) {
	store := NewSessionPlanStore()
	defer store.Stop()
	tools := benchTools(size)
	cs := benchProgressiveSelector(b, store, tools)
	sc := SelectionContext{SessionID: "call", Tenant: "acme", Requested: "acme_tool_0", CatalogVersion: "bench-" + strconv.Itoa(size)}
	plan, err := cs.Select(context.Background(), sc, tools)
	if err != nil || plan == nil {
		b.Fatalf("prewarm failed: plan=%v err=%v", plan, err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		benchReceipt, benchErr = cs.AuthorizeCall(context.Background(), sc, tools)
		if benchErr == nil {
			benchResult, benchErr = cs.RecordCallSuccess(context.Background(), *benchReceipt)
		}
	}
	b.StopTimer()
	if benchErr != nil || benchReceipt == nil {
		b.Fatalf("call benchmark failed: receipt=%v result=%v err=%v", benchReceipt, benchResult, benchErr)
	}
	requireBenchStoreSize(b, "call", store, 1)
}

func benchmarkCompositeSelectorCold(b *testing.B, size int) {
	store := NewSessionPlanStore()
	defer store.Stop()
	cs := benchSelector(b, store)
	tools := benchTools(size)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sessionID := "cold"
		sc := SelectionContext{SessionID: sessionID, Tenant: "acme", CatalogVersion: "bench-" + strconv.Itoa(size)}
		benchPlan, benchErr = cs.Select(context.Background(), sc, tools)
		store.Delete(cs.planKey(sessionID))
	}
	b.StopTimer()
	requireBenchPlan(b, "cold", benchPlan, benchErr)
	requireBenchStoreSize(b, "cold", store, 0)
}

func requireBenchPlan(b *testing.B, name string, plan *SelectionPlan, err error) {
	b.Helper()
	if err != nil || plan == nil {
		b.Fatalf("%s select failed: plan=%v err=%v", name, plan, err)
	}
}

func requireBenchStoreSize(b *testing.B, name string, store *SessionPlanStore, want int) {
	b.Helper()
	if store.Len() != want {
		b.Fatalf("%s benchmark store size = %d, want %d", name, store.Len(), want)
	}
}

func BenchmarkSessionPlanStore_Set(b *testing.B) {
	store := NewSessionPlanStore()
	defer store.Stop()
	key := NewPlanKey("bench", "s")
	plan := &SelectionPlan{SessionID: "s", ToolNames: []string{"a", "b"}, Version: "v"}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		store.Set(key, plan, SelectionContext{SessionID: key.SessionID})
	}
}

func BenchmarkSessionPlanStore_Get(b *testing.B) {
	store := NewSessionPlanStore()
	defer store.Stop()
	key := NewPlanKey("bench", "s")
	store.Set(key, &SelectionPlan{SessionID: "s", ToolNames: []string{"a", "b"}, Version: "v"}, SelectionContext{SessionID: key.SessionID})

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchPlan, benchOK = store.Get(key)
	}
	if !benchOK || benchPlan == nil {
		b.Fatal("expected cached plan")
	}
}

func BenchmarkSessionPlanStore_Delete(b *testing.B) {
	store := NewSessionPlanStore()
	defer store.Stop()
	key := NewPlanKey("bench", "s")
	plan := &SelectionPlan{SessionID: "s", ToolNames: []string{"a"}, Version: "v"}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		store.Set(key, plan, SelectionContext{SessionID: key.SessionID})
		store.Delete(key)
	}
	if store.Len() != 0 {
		b.Fatalf("store size = %d, want 0", store.Len())
	}
}

func BenchmarkSessionPlanStore_RecordCallSuccess(b *testing.B) {
	store := NewSessionPlanStore()
	defer store.Stop()
	key := NewPlanKey("bench", "s")
	store.Set(key, &SelectionPlan{SessionID: "s", ToolNames: []string{"a"}, Version: "v"}, SelectionContext{SessionID: key.SessionID})

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		plan, _ := store.Get(key)
		receipt, _ := store.IssueReceipt(key, "a", plan, key.RouterID)
		_ = store.RecordCallSuccess(receipt, b.N+1)
	}
}

func BenchmarkSessionPlanStore_ThresholdTransition(b *testing.B) {
	store := NewSessionPlanStore()
	defer store.Stop()
	key := NewPlanKey("bench", "s")

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		store.Set(key, &SelectionPlan{SessionID: "s", ToolNames: []string{"a"}, Version: "v"}, SelectionContext{SessionID: key.SessionID})
		plan, _ := store.Get(key)
		receipt, _ := store.IssueReceipt(key, "a", plan, key.RouterID)
		_ = store.RecordCallSuccess(receipt, 1)
		store.Delete(key)
	}
	if store.Len() != 0 {
		b.Fatalf("store size = %d, want 0", store.Len())
	}
}

func BenchmarkSessionPlanStore_ConcurrentGetSet(b *testing.B) {
	store := NewSessionPlanStore()
	defer store.Stop()
	key := NewPlanKey("bench", "s")
	plan := &SelectionPlan{SessionID: "s", ToolNames: []string{"a"}, Version: "v"}

	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			store.Set(key, plan, SelectionContext{SessionID: key.SessionID})
			benchPlan, benchOK = store.Get(key)
		}
	})
	if !benchOK || benchPlan == nil {
		b.Fatal("expected cached plan")
	}
}

func BenchmarkSessionPlanStore_CapacityBoundedInsert(b *testing.B) {
	store := NewSessionPlanStoreWithOptions(SessionPlanStoreOptions{
		MaxEntries: 64,
	})
	defer store.Stop()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sessionID := fmt.Sprintf("s-%d", i)
		key := NewPlanKey("bench", sessionID)
		err := store.Set(key, &SelectionPlan{SessionID: sessionID, ToolNames: []string{"a"}, Version: "v"}, SelectionContext{SessionID: sessionID})
		if err != nil && err != ErrPlanStoreFull {
			b.Fatal(err)
		}
	}
	if store.Len() > 64 {
		b.Fatalf("store size = %d, want <= 64", store.Len())
	}
}

func BenchmarkSessionPlanStore_CapacityReject(b *testing.B) {
	store := NewSessionPlanStoreWithOptions(SessionPlanStoreOptions{
		MaxEntries: 128,
	})
	defer store.Stop()

	for i := 0; i < 128; i++ {
		sessionID := fmt.Sprintf("s-%d", i)
		key := NewPlanKey("bench", sessionID)
		if err := store.Set(key, &SelectionPlan{SessionID: sessionID, ToolNames: []string{"a"}, Version: "v"}, SelectionContext{SessionID: sessionID}); err != nil {
			b.Fatal(err)
		}
	}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		key := NewPlanKey("bench", fmt.Sprintf("overflow-%d", i))
		if err := store.Set(key, &SelectionPlan{SessionID: key.SessionID, ToolNames: []string{"a"}, Version: "v"}, SelectionContext{SessionID: key.SessionID}); err != ErrPlanStoreFull {
			b.Fatalf("expected ErrPlanStoreFull, got %v", err)
		}
	}
	if store.Len() != 128 {
		b.Fatalf("store size = %d, want 128", store.Len())
	}
}
