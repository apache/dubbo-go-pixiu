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
	"testing"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// TestMethodWildcard_Dubbo verifies the original bug: methods: ["*"] should
// match arbitrary RPC method names (Dubbo $invoke / Triple GetUser / etc.)
// when looked up via RouteByPathAndName.
func TestMethodWildcard_Dubbo(t *testing.T) {
	r1 := newTestRouter("r1", "*", "/com.example.UserService", "cluster-user")

	cfg := &model.RouteConfiguration{
		Routes:  []*model.Router{r1},
		Dynamic: false,
	}
	rm := CreateRouterCoordinator(cfg)
	rm.debounce = 0

	mustRouteByPathAndNameCluster(t, rm, "$invoke", "/com.example.UserService", "cluster-user")
	mustRouteByPathAndNameCluster(t, rm, "GetUser", "/com.example.UserService", "cluster-user")
	mustRouteByPathAndNameCluster(t, rm, "UpdateOrder", "/com.example.UserService", "cluster-user")
}

// TestMethodWildcard_HTTP verifies methods: ["*"] on the HTTP path so HTTP
// configurations using "*" continue to match any verb.
func TestMethodWildcard_HTTP(t *testing.T) {
	r1 := newTestRouter("r1", "*", "/foo", "cluster-foo")

	cfg := &model.RouteConfiguration{
		Routes:  []*model.Router{r1},
		Dynamic: false,
	}
	rm := CreateRouterCoordinator(cfg)
	rm.debounce = 0

	mustRouteCluster(t, rm, "GET", "/foo", "cluster-foo")
	mustRouteCluster(t, rm, "POST", "/foo", "cluster-foo")
	mustRouteCluster(t, rm, "DELETE", "/foo", "cluster-foo")
}

// TestMethodWildcard_FallbackWhenSpecificMisses is the discriminator for
// Option B: when MethodTries[method] exists but does NOT match the path,
// the lookup must still fall back to MethodTries["*"].
//
// Without the fix, GET /bar would fail because MethodTries["GET"] only
// contains /foo, and the original code returned an error without ever
// consulting MethodTries["*"].
func TestMethodWildcard_FallbackWhenSpecificMisses(t *testing.T) {
	rGet := newTestRouter("r1", "GET", "/foo", "cluster-foo")
	rStar := newTestRouter("r2", "*", "/bar", "cluster-bar")

	cfg := &model.RouteConfiguration{
		Routes:  []*model.Router{rGet, rStar},
		Dynamic: false,
	}
	rm := CreateRouterCoordinator(cfg)
	rm.debounce = 0

	// specific bucket hits as before
	mustRouteCluster(t, rm, "GET", "/foo", "cluster-foo")

	// specific bucket exists but path is not in it → fall back to "*"
	mustRouteCluster(t, rm, "GET", "/bar", "cluster-bar")
	mustRouteByPathAndNameCluster(t, rm, "GET", "/bar", "cluster-bar")
}

// TestMethodWildcard_SpecificBeatsWildcard guards against an over-eager
// fallback: when both a specific and a wildcard route match the SAME path,
// the specific one must win.
func TestMethodWildcard_SpecificBeatsWildcard(t *testing.T) {
	rGet := newTestRouter("r1", "GET", "/foo", "cluster-get")
	rStar := newTestRouter("r2", "*", "/foo", "cluster-star")

	cfg := &model.RouteConfiguration{
		Routes:  []*model.Router{rGet, rStar},
		Dynamic: false,
	}
	rm := CreateRouterCoordinator(cfg)
	rm.debounce = 0

	// GET hits the specific GET route, not the wildcard
	mustRouteCluster(t, rm, "GET", "/foo", "cluster-get")
	// POST has no specific bucket → falls back to wildcard
	mustRouteCluster(t, rm, "POST", "/foo", "cluster-star")
}

// TestMethodWildcard_MethodAllowed exercises MethodAllowed directly, which
// is what header-only routes use.
func TestMethodWildcard_MethodAllowed(t *testing.T) {
	if !model.MethodAllowed([]string{"*"}, "GET") {
		t.Fatalf("MethodAllowed([*], GET) should be true")
	}
	if !model.MethodAllowed([]string{"*"}, "$invoke") {
		t.Fatalf("MethodAllowed([*], $invoke) should be true")
	}
	if !model.MethodAllowed([]string{"GET", "*"}, "POST") {
		t.Fatalf("MethodAllowed([GET, *], POST) should be true")
	}
	if model.MethodAllowed([]string{"GET"}, "POST") {
		t.Fatalf("MethodAllowed([GET], POST) should be false")
	}
}
