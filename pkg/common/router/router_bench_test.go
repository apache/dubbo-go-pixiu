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
	"math/rand"
	stdHttp "net/http"
	"reflect"
	"strconv"
	"testing"
	"time"
)

import (
	oldrouter "github.com/apache/dubbo-go-pixiu/pkg/common/router/mock"
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

/*
==============================
this is the benchmark for router
contrast oldrouter and newrouter
oldrouter: "github.com/apache/dubbo-go-pixiu/pkg/common/router/mock"
newrouter: "github.com/apache/dubbo-go-pixiu/pkg/common/router"
==============================
*/
type benchShape struct {
	NRoutes         int     // router number
	PrefixRatio     float64 // prefix router ratio（others are accurate path）
	HeaderOnlyRatio float64 // header only router ratio（without path/prefix）
	Methods         []string
}

func buildOldCoordinator(routes []*model.Router) *oldrouter.RouterCoordinator {
	cfg := &model.RouteConfiguration{
		Routes:  routes,
		Dynamic: false,
	}
	return oldrouter.CreateRouterCoordinator(cfg)
}

func genRoutes(sh benchShape) []*model.Router {
	routes := make([]*model.Router, 0, sh.NRoutes)
	if len(sh.Methods) == 0 {
		sh.Methods = []string{"GET", "POST"}
	}
	nHeader := int(float64(sh.NRoutes) * sh.HeaderOnlyRatio)
	nPrefix := int(float64(sh.NRoutes-nHeader) * sh.PrefixRatio)
	nPath := sh.NRoutes - nHeader - nPrefix

	// 1) Header-only
	for i := 0; i < nHeader; i++ {
		id := "hdr-" + strconv.Itoa(i)
		r := &model.Router{
			ID: id,
			Match: model.RouterMatch{
				Methods: sh.Methods,
				Headers: []model.HeaderMatcher{
					{Name: "X-Env", Values: []string{"prod"}, Regex: false},
				},
			},
			Route: model.RouteAction{Cluster: "c-h-" + id},
		}
		routes = append(routes, r)
	}
	// 2) Prefix routes
	for i := 0; i < nPrefix; i++ {
		id := "pre-" + strconv.Itoa(i)
		p := "/api/v1/service" + strconv.Itoa(i%50) + "/"
		r := &model.Router{
			ID: id,
			Match: model.RouterMatch{
				Methods: sh.Methods,
				Prefix:  p,
			},
			Route: model.RouteAction{Cluster: "c-p-" + id},
		}
		routes = append(routes, r)
	}
	// 3) Exact path
	for i := 0; i < nPath; i++ {
		id := "pth-" + strconv.Itoa(i)
		pp := "/api/v1/item/" + strconv.Itoa(i)
		r := &model.Router{
			ID: id,
			Match: model.RouterMatch{
				Methods: sh.Methods,
				Path:    pp,
			},
			Route: model.RouteAction{Cluster: "c-x-" + id},
		}
		routes = append(routes, r)
	}
	return routes
}

func buildNewCoordinator(routes []*model.Router) *RouterCoordinator {
	cfg := &model.RouteConfiguration{
		Routes:  routes,
		Dynamic: false,
	}
	return CreateRouterCoordinator(cfg)
}

func buildDelta(base []*model.Router, seed int64) []*model.Router {
	cp := make([]*model.Router, len(base))
	copy(cp, base)
	rnd := rand.New(rand.NewSource(seed))
	k := len(cp) / 100 // 1%
	out := make([]*model.Router, 0, k)
	for i := 0; i < k; i++ {
		idx := rnd.Intn(len(cp)) // NOSONAR
		old := cp[idx]
		newPath := "/api/v1/item/" + strconv.Itoa(rnd.Intn(100000)) // NOSONAR
		nr := &model.Router{
			ID: old.ID,
			Match: model.RouterMatch{
				Methods: old.Match.Methods,
				Path:    newPath,
				Headers: old.Match.Headers,
			},
			Route: old.Route,
		}
		out = append(out, nr)
	}
	return out
}

func genRequests(n int) []*stdHttp.Request {
	reqs := make([]*stdHttp.Request, 0, n)
	methods := []string{"GET", "POST"}
	for i := 0; i < n; i++ {
		var path string
		switch i % 3 {
		case 0:
			path = "/api/v1/item/" + strconv.Itoa(i%10000)
		case 1:
			path = "/api/v1/service" + strconv.Itoa(i%50) + "/foo/bar"
		default:
			path = "/unknown/" + strconv.Itoa(i)
		}
		req, _ := stdHttp.NewRequest(methods[i%len(methods)], path, nil)
		if i%5 == 0 { // trigger header-only route
			req.Header.Set("X-Env", "prod")
		}
		reqs = append(reqs, req)
	}
	return reqs
}

// helper: assert Route behavior of old/new is the same on a set of requests
func assertRouteSame(b testing.TB, oldc *oldrouter.RouterCoordinator, newc *RouterCoordinator, reqs []*stdHttp.Request) {
	b.Helper()
	for i, r := range reqs {
		ctxOld := http.HttpContext{Request: r}
		ctxNew := http.HttpContext{Request: r}

		oldRes, oldErr := oldc.Route(&ctxOld)
		newRes, newErr := newc.Route(&ctxNew)

		if (oldErr != nil && newErr == nil) || (oldErr == nil && newErr != nil) {
			b.Fatalf("route error text mismatch on #%d path=%s method=%s: oldErr=%v newErr=%v",
				i, r.URL.Path, r.Method, oldErr, newErr)
		}
		if !reflect.DeepEqual(oldRes, newRes) {
			b.Fatalf("route result mismatch on #%d path=%s method=%s: old=%#v new=%#v",
				i, r.URL.Path, r.Method, oldRes, newRes)
		}
	}
}

// helper: assert RouteByPathAndName behavior is the same
func assertRouteByPathAndNameSame(b testing.TB, oldc *oldrouter.RouterCoordinator, newc *RouterCoordinator, paths []string, method string) {
	b.Helper()
	for i, p := range paths {
		oldRes, oldErr := oldc.RouteByPathAndName(p, method)
		newRes, newErr := newc.RouteByPathAndName(p, method)

		if (oldErr != nil && newErr == nil) || (oldErr == nil && newErr != nil) {
			b.Fatalf("route error text mismatch on #%d path=%s method=%s: oldRes=%v oldErr=%v  newRes=%v newErr=%v",
				i, p, method, oldRes, oldErr, newRes, newErr)
		}
		if !reflect.DeepEqual(oldRes, newRes) {
			b.Fatalf("RouteByPathAndName result mismatch on #%d path=%s method=%s: old=%#v new=%#v",
				i, p, method, oldRes, newRes)
		}
	}
}

// ============= Bench 1：read throughput (one goroutine) =============

func BenchmarkRouteReadThroughput(b *testing.B) {
	shape := benchShape{NRoutes: 30000, PrefixRatio: 0.4, HeaderOnlyRatio: 0.1, Methods: []string{"GET", "POST"}}

	oldRoutes := genRoutes(shape)
	newRoutes := genRoutes(shape)
	reqs := genRequests(4096)

	oldc := buildOldCoordinator(oldRoutes)
	newc := buildNewCoordinator(newRoutes)

	assertRouteSame(b, oldc, newc, reqs)

	b.Run("old/locked-read-30k", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			r := reqs[i%len(reqs)]
			httpContext := http.HttpContext{
				Request: r,
			}
			_, _ = oldc.Route(&httpContext)
		}
	})

	b.Run("new/rcu-read-30k", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			r := reqs[i%len(reqs)]
			httpContext := http.HttpContext{
				Request: r,
			}
			_, _ = newc.Route(&httpContext)
		}
	})
}

// ============= Bench 2：read throughput (parallel) =============

func BenchmarkRouteReadParallel(b *testing.B) {
	shape := benchShape{NRoutes: 30000, PrefixRatio: 0.4, HeaderOnlyRatio: 0.1, Methods: []string{"GET", "POST"}}

	oldRoutes := genRoutes(shape)
	newRoutes := genRoutes(shape)
	reqs := genRequests(8192)

	oldc := buildOldCoordinator(oldRoutes)
	newc := buildNewCoordinator(newRoutes)

	assertRouteSame(b, oldc, newc, reqs)

	b.Run("old/parallel-30k", func(b *testing.B) {
		b.ReportAllocs()
		b.SetParallelism(40)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			i := rand.Int() // NOSONAR
			for pb.Next() {
				r := reqs[i%len(reqs)]
				httpContext := http.HttpContext{
					Request: r,
				}
				_, _ = oldc.Route(&httpContext)
				i++
			}
		})
	})

	b.Run("new/parallel-30k", func(b *testing.B) {
		b.ReportAllocs()
		b.SetParallelism(40)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			i := rand.Int() // NOSONAR
			for pb.Next() {
				r := reqs[i%len(reqs)]
				httpContext := http.HttpContext{
					Request: r,
				}
				_, _ = newc.Route(&httpContext)
				i++
			}
		})
	})
}

// ============= Bench 3：read and write（1% write） =============

func BenchmarkReloadLatency(b *testing.B) {
	shape := benchShape{NRoutes: 30000, PrefixRatio: 0.4, HeaderOnlyRatio: 0.1, Methods: []string{"GET", "POST"}}
	base := genRoutes(shape)

	oldc := buildOldCoordinator(base)
	newc := buildNewCoordinator(base)

	{
		checkOld := buildOldCoordinator(base)
		checkNew := buildNewCoordinator(base)
		deltaOld := buildDelta(base, 1)
		for i := range deltaOld {
			checkOld.OnAddRouter(deltaOld[i])
			checkNew.OnAddRouter(deltaOld[i])
		}
		reqs := genRequests(1024)
		time.Sleep(100 * time.Millisecond)
		assertRouteSame(b, checkOld, checkNew, reqs)
	}

	b.Run("old/reload-1percent-30k", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, r := range buildDelta(base, int64(i)) {
				oldc.OnAddRouter(r)
			}
		}
	})

	b.Run("new/reload-1percent-30k", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, r := range buildDelta(base, int64(i)) {
				newc.OnAddRouter(r)
			}
		}
	})
}

func BenchmarkRoute100kReadThroughput(b *testing.B) {
	shape := benchShape{
		NRoutes:         100_000,
		PrefixRatio:     0.4,
		HeaderOnlyRatio: 0.1,
		Methods:         []string{"GET", "POST"},
	}
	reqs := genRequests(16_384)

	oldc := buildOldCoordinator(genRoutes(shape))
	newc := buildNewCoordinator(genRoutes(shape))

	assertRouteSame(b, oldc, newc, reqs)

	b.Run("old/locked-read-100k", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			r := reqs[i%len(reqs)]
			httpContext := http.HttpContext{
				Request: r,
			}
			_, _ = oldc.Route(&httpContext)
		}
	})
	b.Run("new/rcu-read-100k", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			r := reqs[i%len(reqs)]
			httpContext := http.HttpContext{
				Request: r,
			}
			_, _ = newc.Route(&httpContext)
		}
	})
}

func BenchmarkRoute100kReadParallel(b *testing.B) {
	shape := benchShape{
		NRoutes:         100_000,
		PrefixRatio:     0.4,
		HeaderOnlyRatio: 0.1,
		Methods:         []string{"GET", "POST"},
	}
	reqs := genRequests(32_768)

	oldc := buildOldCoordinator(genRoutes(shape))
	newc := buildNewCoordinator(genRoutes(shape))

	assertRouteSame(b, oldc, newc, reqs)

	b.Run("old/parallel-100k", func(b *testing.B) {
		b.ReportAllocs()
		b.SetParallelism(4)
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				r := reqs[i&(len(reqs)-1)]
				httpContext := http.HttpContext{
					Request: r,
				}
				_, _ = oldc.Route(&httpContext)
				i++
			}
		})
	})
	b.Run("new/parallel-100k", func(b *testing.B) {
		b.ReportAllocs()
		b.SetParallelism(4)
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				r := reqs[i&(len(reqs)-1)]
				httpContext := http.HttpContext{
					Request: r,
				}
				_, _ = newc.Route(&httpContext)
				i++
			}
		})
	})
}

func BenchmarkReload100kLatency1Percent(b *testing.B) {
	shape := benchShape{
		NRoutes:         100_000,
		PrefixRatio:     0.4,
		HeaderOnlyRatio: 0.1,
		Methods:         []string{"GET", "POST"},
	}
	base := genRoutes(shape)

	oldc := buildOldCoordinator(genRoutes(shape))
	newc := buildNewCoordinator(genRoutes(shape))

	{
		checkOld := buildOldCoordinator(base)
		checkNew := buildNewCoordinator(base)
		deltaOld := buildDelta(base, 1)
		for i := range deltaOld {
			checkOld.OnAddRouter(deltaOld[i])
			checkNew.OnAddRouter(deltaOld[i])
		}
		reqs := genRequests(2048)
		time.Sleep(100 * time.Millisecond)
		assertRouteSame(b, checkOld, checkNew, reqs)
	}

	b.Run("old/reload-1percent-100k", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, r := range buildDelta(base, int64(i)) {
				oldc.OnAddRouter(r)
			}
		}
	})

	b.Run("new/reload-1percent-100k", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for _, r := range buildDelta(base, int64(i)) {
				newc.OnAddRouter(r)
			}
		}
	})
}

// ============= Bench 4：RouteByPathAndName（API behavior must same） =============

func BenchmarkRouteByPathAndName(b *testing.B) {
	shape := benchShape{
		NRoutes:         20000,
		PrefixRatio:     0.5,
		HeaderOnlyRatio: 0.0,
		Methods:         []string{"GET"},
	}

	oldc := buildOldCoordinator(genRoutes(shape))
	newc := buildNewCoordinator(genRoutes(shape))

	paths := []string{
		"/api/v1/item/12345",
		"/api/v1/service7/xxx/yyy",
		"/no/match/path",
	}
	method := "GET"

	assertRouteByPathAndNameSame(b, oldc, newc, paths, method)

	b.Run("old/RouteByPathAndName", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			path := paths[i%len(paths)]
			_, _ = oldc.RouteByPathAndName(path, method)
		}
	})

	b.Run("new/RouteByPathAndName", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			path := paths[i%len(paths)]
			_, _ = newc.RouteByPathAndName(path, method)
		}
	})
}
