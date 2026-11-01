/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to you under the Apache License, Version 2.0
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
	stdhttp "net/http"
	"testing"
)

import (
	pixiuhttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func TestCombinedPathHeaderRouteRequiresEveryHeader(t *testing.T) {
	route := &model.Router{
		ID: "user-v1",
		Match: model.RouterMatch{
			Methods: []string{"GET"},
			Prefix:  "/user",
			Headers: []model.HeaderMatcher{
				{Name: "X-A", Values: []string{"t1"}},
				{Name: "X-B", Values: []string{"t4"}},
			},
		},
		Route: model.RouteAction{Cluster: "user-v1"},
	}
	rm := CreateRouterCoordinator(&model.RouteConfiguration{Routes: []*model.Router{route}})

	for _, tc := range []struct {
		name    string
		path    string
		method  string
		headers map[string]string
		want    string
	}{
		{name: "all conditions match", path: "/user/42", method: "GET", headers: map[string]string{"X-A": "t1", "X-B": "t4"}, want: "user-v1"},
		{name: "no headers", path: "/user/42", method: "GET"},
		{name: "missing second header", path: "/user/42", method: "GET", headers: map[string]string{"X-A": "t1"}},
		{name: "wrong header value", path: "/user/42", method: "GET", headers: map[string]string{"X-A": "t1", "X-B": "wrong"}},
		{name: "wrong path", path: "/other/42", method: "GET", headers: map[string]string{"X-A": "t1", "X-B": "t4"}},
		{name: "wrong method", path: "/user/42", method: "POST", headers: map[string]string{"X-A": "t1", "X-B": "t4"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req, err := stdhttp.NewRequest(tc.method, tc.path, nil)
			if err != nil {
				t.Fatal(err)
			}
			for name, value := range tc.headers {
				req.Header.Set(name, value)
			}
			action, err := rm.Route(&pixiuhttp.HttpContext{Request: req})
			if tc.want == "" {
				if err == nil || action != nil {
					t.Fatalf("Route() = (%#v, %v), want no matching route", action, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Route() error = %v", err)
			}
			if action.Cluster != tc.want {
				t.Fatalf("Route() cluster = %q, want %q", action.Cluster, tc.want)
			}
		})
	}
}

func TestCombinedRouteCandidatesAndPriority(t *testing.T) {
	routes := []*model.Router{
		combinedRoute("first", []string{"GET"}, "/variant", "", []model.HeaderMatcher{{Name: "X-Variant", Values: []string{"v1", "v1-alt"}}}, "first"),
		combinedRoute("second", []string{"GET"}, "/variant", "", []model.HeaderMatcher{{Name: "X-Variant", Values: []string{"v2"}}}, "second"),
		combinedRoute("prefix", []string{"GET"}, "", "/items", []model.HeaderMatcher{{Name: "X-Route", Values: []string{"prefix"}}}, "prefix"),
		combinedRoute("exact", []string{"GET"}, "/items/42", "", []model.HeaderMatcher{{Name: "X-Route", Values: []string{"exact"}}}, "exact"),
		combinedRoute("specific-method", []string{"GET"}, "/wildcard", "", []model.HeaderMatcher{{Name: "X-Route", Values: []string{"specific"}}}, "specific"),
		combinedRoute("any-method", []string{"*"}, "/wildcard", "", []model.HeaderMatcher{{Name: "X-Route", Values: []string{"any"}}}, "any"),
		combinedRoute("regex-and", []string{"GET"}, "/regex", "", []model.HeaderMatcher{
			{Name: "X-Route", Values: []string{"^t[0-9]+$"}, Regex: true},
			{Name: "X-Group", Values: []string{"a", "b"}},
			{Name: "X-Tenant", Values: []string{"north"}},
		}, "regex"),
		combinedRoute("path-priority", []string{"GET"}, "/priority", "", []model.HeaderMatcher{{Name: "X-Route", Values: []string{"same"}}}, "path"),
		combinedRoute("header-priority", []string{"GET"}, "", "", []model.HeaderMatcher{{Name: "X-Route", Values: []string{"same"}}}, "header-only"),
		combinedRoute("rpc", []string{"*"}, "/com.example.Service", "", []model.HeaderMatcher{{Name: "X-Route", Values: []string{"rpc"}}}, "rpc"),
	}
	rm := CreateRouterCoordinator(&model.RouteConfiguration{Routes: routes})

	for _, tc := range []struct {
		name, method, path, want string
		headers                  map[string]string
	}{
		{name: "same path preserves candidate order", method: "GET", path: "/variant", headers: map[string]string{"X-Variant": "v2"}, want: "second"},
		{name: "same path first matching candidate wins", method: "GET", path: "/variant", headers: map[string]string{"X-Variant": "v1"}, want: "first"},
		{name: "same path values are OR", method: "GET", path: "/variant", headers: map[string]string{"X-Variant": "v1-alt"}, want: "first"},
		{name: "exact rejection falls back to prefix", method: "GET", path: "/items/42", headers: map[string]string{"X-Route": "prefix"}, want: "prefix"},
		{name: "exact path wins when its condition matches", method: "GET", path: "/items/42", headers: map[string]string{"X-Route": "exact"}, want: "exact"},
		{name: "wildcard method checked after specific rejection", method: "GET", path: "/wildcard", headers: map[string]string{"X-Route": "any"}, want: "any"},
		{name: "specific method wins when both match", method: "GET", path: "/wildcard", headers: map[string]string{"X-Route": "specific"}, want: "specific"},
		{name: "regex and repeated matchers are AND", method: "GET", path: "/regex", headers: map[string]string{"X-Route": "t17", "X-Group": "b", "X-Tenant": "north"}, want: "regex"},
		{name: "repeated matcher missing rejects", method: "GET", path: "/regex", headers: map[string]string{"X-Route": "t17", "X-Group": "b"}},
		{name: "regex mismatch rejects", method: "GET", path: "/regex", headers: map[string]string{"X-Route": "no", "X-Group": "b", "X-Tenant": "north"}},
		{name: "header-only remains first", method: "GET", path: "/priority", headers: map[string]string{"X-Route": "same"}, want: "header-only"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assertCombinedRoute(t, rm, tc.method, tc.path, tc.headers, tc.want)
		})
	}

	// RouteByPathAndName has no header input and retains its path-only behavior.
	action, err := rm.RouteByPathAndName("/variant", "GET")
	if err != nil || action.Cluster != "first" {
		t.Fatalf("RouteByPathAndName() = (%#v, %v), want first candidate", action, err)
	}
	action, err = rm.RouteByPathAndName("/com.example.Service", "GetUser")
	if err != nil || action.Cluster != "rpc" {
		t.Fatalf("RPC RouteByPathAndName() = (%#v, %v), want wildcard-method route", action, err)
	}
}

func TestCombinedRouteDynamicSnapshotAddAndDelete(t *testing.T) {
	first := combinedRoute("dynamic-first", []string{"GET"}, "/dynamic", "", []model.HeaderMatcher{{Name: "X-Route", Values: []string{"first"}}}, "first")
	rm := CreateRouterCoordinator(&model.RouteConfiguration{Routes: []*model.Router{first}, Dynamic: true})
	rm.debounce = 0

	second := combinedRoute("dynamic-second", []string{"GET"}, "/dynamic", "", []model.HeaderMatcher{{Name: "X-Route", Values: []string{"second"}}}, "second")
	rm.OnAddRouter(second)
	assertCombinedRoute(t, rm, "GET", "/dynamic", map[string]string{"X-Route": "second"}, "second")

	rm.OnDeleteRouter(first)
	assertCombinedRoute(t, rm, "GET", "/dynamic", map[string]string{"X-Route": "first"}, "")
	assertCombinedRoute(t, rm, "GET", "/dynamic", map[string]string{"X-Route": "second"}, "second")
}

func TestCombinedEquivalentPathCandidates(t *testing.T) {
	for _, tc := range []struct {
		name, first, second, path string
		prefix                    bool
	}{
		{name: "renamed parameter", first: "/orders/:id", second: "/orders/:name", path: "/orders/42"},
		{name: "wildcard then parameter", first: "/orders/*", second: "/orders/:id", path: "/orders/42"},
		{name: "parameter then wildcard", first: "/orders/:id", second: "/orders/*", path: "/orders/42"},
		{name: "renamed prefix parameter", first: "/orders/:id", second: "/orders/:name", path: "/orders/42/details", prefix: true},
		{name: "renamed intermediate parameter", first: "/orders/:id/items", second: "/orders/:name/items", path: "/orders/42/items"},
		{name: "intermediate wildcard", first: "/orders/*/items", second: "/orders/:id/items", path: "/orders/42/items"},
		{name: "multiple renamed parameters", first: "/orgs/:org/orders/:id", second: "/orgs/:tenant/orders/:name", path: "/orgs/7/orders/42"},
		{name: "renamed intermediate prefix parameter", first: "/orders/:id/items", second: "/orders/:name/items", path: "/orders/42/items/details", prefix: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			first := combinedRoute("first", []string{"GET"}, tc.first, "", []model.HeaderMatcher{{Name: "X-Tenant", Values: []string{"a", "shared"}}}, "first")
			second := combinedRoute("second", []string{"GET"}, tc.second, "", []model.HeaderMatcher{{Name: "X-Tenant", Values: []string{"b", "shared"}}}, "second")
			if tc.prefix {
				first.Match.Path, first.Match.Prefix = "", tc.first
				second.Match.Path, second.Match.Prefix = "", tc.second
			}
			rm := CreateRouterCoordinator(&model.RouteConfiguration{Routes: []*model.Router{first, second}})
			assertCombinedRoute(t, rm, "GET", tc.path, map[string]string{"X-Tenant": "a"}, "first")
			assertCombinedRoute(t, rm, "GET", tc.path, map[string]string{"X-Tenant": "b"}, "second")
			assertCombinedRoute(t, rm, "GET", tc.path, map[string]string{"X-Tenant": "shared"}, "first")
			assertCombinedRoute(t, rm, "GET", tc.path, map[string]string{"X-Tenant": "other"}, "")

			action, err := rm.RouteByPathAndName(tc.path, "GET")
			if err != nil || action.Cluster != "first" {
				t.Fatalf("RouteByPathAndName() = (%#v, %v), want first candidate", action, err)
			}
		})
	}
}

func TestCombinedTrieIdentityBoundaries(t *testing.T) {
	routes := []*model.Router{
		combinedRoute("root", []string{"GET"}, "/", "", []model.HeaderMatcher{{Name: "X-Route", Values: []string{"root"}}}, "root"),
		combinedRoute("parameter", []string{"GET"}, "/:token", "", []model.HeaderMatcher{{Name: "X-Route", Values: []string{"parameter"}}}, "parameter"),
		combinedRoute("single-wildcard", []string{"GET"}, "/*", "", []model.HeaderMatcher{{Name: "X-Route", Values: []string{"single"}}}, "single"),
		combinedRoute("match-all", []string{"GET"}, "/**", "", []model.HeaderMatcher{{Name: "X-Route", Values: []string{"all"}}}, "all"),
		combinedRoute("post-parameter", []string{"POST"}, "/:other", "", []model.HeaderMatcher{{Name: "X-Route", Values: []string{"post"}}}, "post"),
	}
	rm := CreateRouterCoordinator(&model.RouteConfiguration{Routes: routes})

	assertCombinedRoute(t, rm, "GET", "/", map[string]string{"X-Route": "root"}, "root")
	assertCombinedRoute(t, rm, "GET", "/42?query=removed", map[string]string{"X-Route": "parameter"}, "parameter")
	assertCombinedRoute(t, rm, "GET", "/42", map[string]string{"X-Route": "single"}, "single")
	assertCombinedRoute(t, rm, "GET", "/a/b", map[string]string{"X-Route": "all"}, "all")
	assertCombinedRoute(t, rm, "GET", "/a/b", map[string]string{"X-Route": "single"}, "")
	assertCombinedRoute(t, rm, "GET", "/42", map[string]string{"X-Route": "post"}, "")
	assertCombinedRoute(t, rm, "POST", "/42", map[string]string{"X-Route": "post"}, "post")

	action, err := rm.RouteByPathAndName("/42?query=removed", "GET")
	if err != nil || action.Cluster != "parameter" {
		t.Fatalf("RouteByPathAndName() = (%#v, %v), want parameter after query removal", action, err)
	}
}

func TestCombinedEquivalentPathsKeepDistinctNodes(t *testing.T) {
	routes := []*model.Router{
		combinedRoute("first", []string{"GET"}, "/orders/:id/items", "", []model.HeaderMatcher{{Name: "X-Tenant", Values: []string{"a"}}}, "first"),
		combinedRoute("second", []string{"GET"}, "/orders/:name/items", "", []model.HeaderMatcher{{Name: "X-Tenant", Values: []string{"b"}}}, "second"),
		combinedRoute("other-suffix", []string{"GET"}, "/orders/:id/details", "", []model.HeaderMatcher{{Name: "X-Tenant", Values: []string{"c"}}}, "other-suffix"),
		combinedRoute("literal", []string{"GET"}, "/orders/special/items", "", []model.HeaderMatcher{{Name: "X-Tenant", Values: []string{"d"}}}, "literal"),
		combinedRoute("prefix", []string{"GET"}, "", "/orders", []model.HeaderMatcher{{Name: "X-Tenant", Values: []string{"fallback"}}}, "prefix"),
	}
	rm := CreateRouterCoordinator(&model.RouteConfiguration{Routes: routes})
	assertCombinedRoute(t, rm, "GET", "/orders/42/items", map[string]string{"X-Tenant": "b"}, "second")
	assertCombinedRoute(t, rm, "GET", "/orders/42/items", map[string]string{"X-Tenant": "c"}, "")
	assertCombinedRoute(t, rm, "GET", "/orders/42/items", map[string]string{"X-Tenant": "d"}, "")
	assertCombinedRoute(t, rm, "GET", "/orders/42/details", map[string]string{"X-Tenant": "b"}, "")
	assertCombinedRoute(t, rm, "GET", "/orders/special/items", map[string]string{"X-Tenant": "d"}, "literal")
	assertCombinedRoute(t, rm, "GET", "/orders/42/items", map[string]string{"X-Tenant": "fallback"}, "prefix")
}

func TestCombinedRouteOverlappingPriority(t *testing.T) {
	headers := []model.HeaderMatcher{{Name: "X-Route", Values: []string{"shared"}}}
	routes := []*model.Router{
		combinedRoute("first", []string{"GET"}, "/ordered", "", headers, "first"),
		combinedRoute("second", []string{"GET"}, "/ordered", "", headers, "second"),
		combinedRoute("wildcard", []string{"*"}, "/method", "", headers, "wildcard"),
		combinedRoute("specific", []string{"GET"}, "/method", "", headers, "specific"),
	}
	rm := CreateRouterCoordinator(&model.RouteConfiguration{Routes: routes})
	assertCombinedRoute(t, rm, "GET", "/ordered", map[string]string{"X-Route": "shared"}, "first")
	assertCombinedRoute(t, rm, "GET", "/method", map[string]string{"X-Route": "shared"}, "specific")
}

func TestCombinedEquivalentPathDynamicSnapshotAddAndDelete(t *testing.T) {
	first := combinedRoute("first", []string{"GET"}, "/orders/:id/items", "", []model.HeaderMatcher{{Name: "X-Tenant", Values: []string{"a", "shared"}}}, "first")
	second := combinedRoute("second", []string{"GET"}, "/orders/:name/items", "", []model.HeaderMatcher{{Name: "X-Tenant", Values: []string{"b", "shared"}}}, "second")
	third := combinedRoute("third", []string{"GET"}, "/orders/*/items", "", []model.HeaderMatcher{{Name: "X-Tenant", Values: []string{"c", "shared"}}}, "third")
	rm := CreateRouterCoordinator(&model.RouteConfiguration{Routes: []*model.Router{first}, Dynamic: true})
	rm.debounce = 0

	rm.OnAddRouter(second)
	rm.OnAddRouter(third)
	assertCombinedRoute(t, rm, "GET", "/orders/42/items", map[string]string{"X-Tenant": "a"}, "first")
	assertCombinedRoute(t, rm, "GET", "/orders/42/items", map[string]string{"X-Tenant": "b"}, "second")
	assertCombinedRoute(t, rm, "GET", "/orders/42/items", map[string]string{"X-Tenant": "c"}, "third")
	assertCombinedRoute(t, rm, "GET", "/orders/42/items", map[string]string{"X-Tenant": "shared"}, "first")

	rm.OnDeleteRouter(first)
	assertCombinedRoute(t, rm, "GET", "/orders/42/items", map[string]string{"X-Tenant": "a"}, "")
	assertCombinedRoute(t, rm, "GET", "/orders/42/items", map[string]string{"X-Tenant": "b"}, "second")
	assertCombinedRoute(t, rm, "GET", "/orders/42/items", map[string]string{"X-Tenant": "c"}, "third")
	assertCombinedRoute(t, rm, "GET", "/orders/42/items", map[string]string{"X-Tenant": "shared"}, "second")

	rm.OnDeleteRouter(second)
	assertCombinedRoute(t, rm, "GET", "/orders/42/items", map[string]string{"X-Tenant": "b"}, "")
	assertCombinedRoute(t, rm, "GET", "/orders/42/items", map[string]string{"X-Tenant": "shared"}, "third")
}

func combinedRoute(id string, methods []string, path, prefix string, headers []model.HeaderMatcher, cluster string) *model.Router {
	return &model.Router{
		ID: id,
		Match: model.RouterMatch{
			Methods: methods,
			Path:    path,
			Prefix:  prefix,
			Headers: headers,
		},
		Route: model.RouteAction{Cluster: cluster},
	}
}

func assertCombinedRoute(t *testing.T, rm *RouterCoordinator, method, path string, headers map[string]string, want string) {
	t.Helper()
	req, err := stdhttp.NewRequest(method, path, nil)
	if err != nil {
		t.Fatal(err)
	}
	for name, value := range headers {
		req.Header.Set(name, value)
	}
	action, err := rm.Route(&pixiuhttp.HttpContext{Request: req})
	if want == "" {
		if err == nil || action != nil {
			t.Fatalf("Route(%s %s) = (%#v, %v), want no matching route", method, path, action, err)
		}
		return
	}
	if err != nil {
		t.Fatalf("Route(%s %s) error = %v", method, path, err)
	}
	if action.Cluster != want {
		t.Fatalf("Route(%s %s) cluster = %q, want %q", method, path, action.Cluster, want)
	}
}
