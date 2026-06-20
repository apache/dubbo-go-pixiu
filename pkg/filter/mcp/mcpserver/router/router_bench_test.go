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
	"strconv"
	"testing"
	"time"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

var (
	benchPlan *SelectionPlan
	benchOK   bool
	benchErr  error
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
		Enabled:  true,
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

func BenchmarkCompositeSelector(b *testing.B) {
	for _, size := range []int{50, 1000, 10000} {
		b.Run("warm_"+strconv.Itoa(size), func(b *testing.B) {
			store := NewSessionPlanStoreWithTTL(time.Hour)
			defer store.Stop()
			cs := benchSelector(b, store)
			tools := benchTools(size)
			sc := SelectionContext{SessionID: "warm", Tenant: "acme"}
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
			if benchErr != nil || benchPlan == nil {
				b.Fatalf("warm select failed: plan=%v err=%v", benchPlan, benchErr)
			}
			if store.Len() != 1 {
				b.Fatalf("warm benchmark store size = %d, want 1", store.Len())
			}
		})

		b.Run("cold_"+strconv.Itoa(size), func(b *testing.B) {
			store := NewSessionPlanStoreWithTTL(time.Hour)
			defer store.Stop()
			cs := benchSelector(b, store)
			tools := benchTools(size)

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				sessionID := "cold"
				sc := SelectionContext{SessionID: sessionID, Tenant: "acme"}
				benchPlan, benchErr = cs.Select(context.Background(), sc, tools)
				store.Delete(sessionID)
			}
			b.StopTimer()
			if benchErr != nil || benchPlan == nil {
				b.Fatalf("cold select failed: plan=%v err=%v", benchPlan, benchErr)
			}
			if store.Len() != 0 {
				b.Fatalf("cold benchmark store size = %d, want 0", store.Len())
			}
		})
	}
}

func BenchmarkSessionPlanStore_Set(b *testing.B) {
	store := NewSessionPlanStoreWithTTL(time.Hour)
	defer store.Stop()
	plan := &SelectionPlan{SessionID: "s", ToolNames: []string{"a", "b"}, Version: "v"}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		store.Set(plan)
	}
}

func BenchmarkSessionPlanStore_Get(b *testing.B) {
	store := NewSessionPlanStoreWithTTL(time.Hour)
	defer store.Stop()
	store.Set(&SelectionPlan{SessionID: "s", ToolNames: []string{"a", "b"}, Version: "v"})

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		benchPlan, benchOK = store.Get("s")
	}
	if !benchOK || benchPlan == nil {
		b.Fatal("expected cached plan")
	}
}

func BenchmarkSessionPlanStore_Delete(b *testing.B) {
	store := NewSessionPlanStoreWithTTL(time.Hour)
	defer store.Stop()
	plan := &SelectionPlan{SessionID: "s", ToolNames: []string{"a"}, Version: "v"}

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		store.Set(plan)
		store.Delete("s")
	}
	if store.Len() != 0 {
		b.Fatalf("store size = %d, want 0", store.Len())
	}
}

func BenchmarkSessionPlanStore_RecordCallSuccess(b *testing.B) {
	store := NewSessionPlanStoreWithTTL(time.Hour)
	defer store.Stop()
	store.Set(&SelectionPlan{SessionID: "s", ToolNames: []string{"a"}, Version: "v"})

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = store.RecordCallSuccess("s", "a", b.N+1)
	}
}

func BenchmarkSessionPlanStore_ThresholdTransition(b *testing.B) {
	store := NewSessionPlanStoreWithTTL(time.Hour)
	defer store.Stop()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		store.Set(&SelectionPlan{SessionID: "s", ToolNames: []string{"a"}, Version: "v"})
		_ = store.RecordCallSuccess("s", "a", 1)
		store.Delete("s")
	}
	if store.Len() != 0 {
		b.Fatalf("store size = %d, want 0", store.Len())
	}
}

func BenchmarkSessionPlanStore_ConcurrentGetSet(b *testing.B) {
	store := NewSessionPlanStoreWithTTL(time.Hour)
	defer store.Stop()
	plan := &SelectionPlan{SessionID: "s", ToolNames: []string{"a"}, Version: "v"}

	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			store.Set(plan)
			benchPlan, benchOK = store.Get("s")
		}
	})
	if !benchOK || benchPlan == nil {
		b.Fatal("expected cached plan")
	}
}

func BenchmarkSessionPlanStore_CapacityEviction(b *testing.B) {
	store := NewSessionPlanStoreWithOptions(SessionPlanStoreOptions{
		TTL:        time.Hour,
		MaxEntries: 64,
	})
	defer store.Stop()

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		store.Set(&SelectionPlan{SessionID: fmt.Sprintf("s-%d", i), ToolNames: []string{"a"}, Version: "v"})
	}
	if store.Len() > 64 {
		b.Fatalf("store size = %d, want <= 64", store.Len())
	}
}

func BenchmarkSessionPlanStore_TTLCleanup(b *testing.B) {
	now := time.Unix(1000, 0)
	store := NewSessionPlanStoreWithOptions(SessionPlanStoreOptions{
		TTL:             time.Second,
		Now:             func() time.Time { return now },
		CleanupInterval: time.Hour,
	})
	defer store.Stop()

	for i := 0; i < 128; i++ {
		store.Set(&SelectionPlan{SessionID: fmt.Sprintf("s-%d", i), ToolNames: []string{"a"}, Version: "v"})
	}
	now = now.Add(2 * time.Second)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		store.EvictExpired()
	}
	if store.Len() != 0 {
		b.Fatalf("store size = %d, want 0", store.Len())
	}
}
