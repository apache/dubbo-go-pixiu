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
	"bytes"
	"math/rand"
	stdHttp "net/http"
	"strconv"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	oldrouter "github.com/apache/dubbo-go-pixiu/pkg/common/router/mock"
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/context/mock"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func TestCreateRouterCoordinator(t *testing.T) {
	specs := []RouteSpec{
		// exact
		{ID: "test", Methods: []string{"POST"}, Path: "/api/v1/**", Cluster: "test_dubbo"},
	}

	coordinator := BuildNew(specs)

	request, err := stdHttp.NewRequest("POST", "http://www.dubbogopixiu.com/api/v1?name=tc", bytes.NewReader([]byte("{\"id\":\"12345\"}")))
	assert.NoError(t, err)
	c := mock.GetMockHTTPContext(request)
	a, err := coordinator.Route(c)
	assert.NoError(t, err)
	assert.Equal(t, a.Cluster, "test_dubbo")

	router := &model.Router{
		ID: "1",
		Match: model.RouterMatch{
			Prefix: "/user",
		},
		Route: model.RouteAction{
			Cluster: "test",
		},
	}

	coordinator.OnAddRouter(router)
	coordinator.OnDeleteRouter(router)
}

func TestRoute(t *testing.T) {
	const (
		Cluster1 = "test-cluster-1"
		Cluster2 = "test-cluster-2"
		Cluster3 = "test-cluster-3"
		Cluster4 = "test-cluster-4"
		Cluster5 = "test-cluster-5"
		Cluster6 = "test-cluster-6"
		Cluster7 = "test-cluster-7"
	)

	hcmc := model.HttpConnectionManagerConfig{
		RouteConfig: model.RouteConfiguration{
			Routes: []*model.Router{
				{
					ID: "1",
					Match: model.RouterMatch{
						Headers: []model.HeaderMatcher{
							{
								Name:   "A",
								Values: []string{"1", "2", "0"},
							},
						},
						Methods: []string{"GET", "POST"},
					},
					Route: model.RouteAction{
						Cluster:                     Cluster1,
						ClusterNotFoundResponseCode: 505,
					},
				},
				{
					ID: "2",
					Match: model.RouterMatch{
						Headers: []model.HeaderMatcher{
							{
								Name:   "A",
								Values: []string{"3", "4", "5"},
							},
						},
						Methods: []string{"GET", "POST"},
					},
					Route: model.RouteAction{
						Cluster:                     Cluster2,
						ClusterNotFoundResponseCode: 505,
					},
				},
				{
					ID: "3",
					Match: model.RouterMatch{
						Headers: []model.HeaderMatcher{
							{
								Name:   "B",
								Values: []string{"1"},
							},
						},
						Methods: []string{"GET", "POST"},
					},
					Route: model.RouteAction{
						Cluster:                     Cluster3,
						ClusterNotFoundResponseCode: 505,
					},
				},
				{
					ID: "4",
					Match: model.RouterMatch{
						Headers: []model.HeaderMatcher{
							{
								Name:   "normal-regex",
								Values: []string{"(k){2}"},
								Regex:  true,
							},
							{
								Name:   "broken-regex",
								Values: []string{"(t){2]]"},
								Regex:  true,
							},
						},
						Methods: []string{"GET", "POST"},
					},
					Route: model.RouteAction{
						Cluster:                     Cluster4,
						ClusterNotFoundResponseCode: 505,
					},
				},
				{
					ID: "5",
					Match: model.RouterMatch{
						Headers: []model.HeaderMatcher{
							{
								Name:   "broken-regex",
								Values: []string{"(t){2]]"},
								Regex:  true,
							},
						},
						Methods: []string{"GET", "POST"},
					},
					Route: model.RouteAction{
						Cluster:                     Cluster5,
						ClusterNotFoundResponseCode: 505,
					},
				},
				{
					ID: "6",
					Match: model.RouterMatch{
						Headers: []model.HeaderMatcher{
							{
								Name:   "C",
								Values: []string{"1", "2", "0"},
							},
							{
								Name:   "D",
								Values: []string{"3", "4", "5"},
							},
						},
						Methods: []string{"GET", "POST"},
					},
					Route: model.RouteAction{
						Cluster:                     Cluster6,
						ClusterNotFoundResponseCode: 505,
					},
				},
				{
					ID: "7",
					Match: model.RouterMatch{
						Headers: []model.HeaderMatcher{
							{
								Name:   "E",
								Values: []string{"1", "2", "0"},
							},
							{
								Name:   "normal-regex",
								Values: []string{"(k){2}"},
								Regex:  true,
							},
						},
						Methods: []string{"GET", "POST"},
					},
					Route: model.RouteAction{
						Cluster:                     Cluster7,
						ClusterNotFoundResponseCode: 505,
					},
				},
			},
			Dynamic: false,
		},
		HTTPFilters: []*model.HTTPFilter{
			{
				Name:   "test",
				Config: nil,
			},
		},
		ServerName:        "test_http_dubbo",
		GenerateRequestID: false,
		IdleTimeoutStr:    "100",
	}

	testCases := []struct {
		Name   string
		URL    string
		Method string
		Header map[string]string
		Expect string
	}{
		{
			Name: "one override header",
			URL:  "/user",
			Header: map[string]string{
				"A": "1",
			},
			Expect: Cluster1,
		},
		{
			Name: "one header matched",
			URL:  "/user",
			Header: map[string]string{
				"A": "3",
			},
			Expect: Cluster2,
		},
		{
			Name: "more header with one regex matched",
			URL:  "/user",
			Header: map[string]string{
				"A":            "5",
				"normal-regex": "kkkk",
			},
			Expect: Cluster2,
		},
		{
			Name:   "one header but wrong method",
			URL:    "/user",
			Method: "PUT",
			Header: map[string]string{
				"A": "0",
			},
			Expect: "route failed for PUT/user, no rules matched",
		},
		{
			Name: "one broken regex header",
			URL:  "/user",
			Header: map[string]string{
				"broken-regex": "tt",
			},
			Expect: "route failed for GET/user, no rules matched",
		},
		{
			Name: "one matched header 2",
			Header: map[string]string{
				"B": "1",
			},
			Expect: Cluster3,
		},
		{
			Name:   "only header but wrong method",
			Method: "DELETE",
			Header: map[string]string{
				"B": "1",
			},
			Expect: "route failed for DELETE, no rules matched",
		},
		{
			Name:   "only header but wrong method",
			Method: "DELETE",
			Header: map[string]string{
				"B": "1",
			},
			Expect: "route failed for DELETE, no rules matched",
		},
		{
			Name: "regex AND normal",
			URL:  "/user",
			Header: map[string]string{
				"E":            "0",
				"normal-regex": "kk",
			},
			Expect: Cluster7,
		},
		{
			Name: "normal AND normal",
			URL:  "/user",
			Header: map[string]string{
				"C": "1",
				"D": "3",
			},
			Expect: Cluster6,
		},
	}

	r := CreateRouterCoordinator(&hcmc.RouteConfig)

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			method := "GET"
			if len(tc.Method) > 0 {
				method = tc.Method
			}
			request, err := stdHttp.NewRequest(method, tc.URL, nil)
			assert.NoError(t, err)

			if tc.Header != nil {
				for k, v := range tc.Header {
					request.Header.Set(k, v)
				}
			}
			c := mock.GetMockHTTPContext(request)

			a, err := r.Route(c)
			if err != nil {
				assert.Equal(t, tc.Expect, err.Error())
			} else {
				assert.Equal(t, tc.Expect, a.Cluster)
			}
		})
	}
}

/* ==============================
   below are parity test between old and new router
   ============================== */

type HeaderSpec struct {
	Name   string
	Values []string
	Regex  bool
}

type RouteSpec struct {
	ID      string
	Methods []string
	Path    string
	Prefix  string
	Headers []HeaderSpec
	Cluster string
}

func (s RouteSpec) toRouter() *model.Router {
	h := make([]model.HeaderMatcher, 0, len(s.Headers))
	for _, x := range s.Headers {
		h = append(h, model.HeaderMatcher{Name: x.Name, Values: append([]string(nil), x.Values...), Regex: x.Regex})
	}
	return &model.Router{
		ID: s.ID,
		Match: model.RouterMatch{
			Methods: append([]string(nil), s.Methods...),
			Path:    s.Path,
			Prefix:  s.Prefix,
			Headers: h,
		},
		Route: model.RouteAction{Cluster: s.Cluster},
	}
}

func buildOld(specs []RouteSpec) *oldrouter.RouterCoordinator {
	rs := make([]*model.Router, 0, len(specs))
	for _, s := range specs {
		rs = append(rs, s.toRouter())
	}
	cfg := &model.RouteConfiguration{Routes: rs, Dynamic: false}
	return oldrouter.CreateRouterCoordinator(cfg)
}

func BuildNew(specs []RouteSpec) *RouterCoordinator {
	rs := make([]*model.Router, 0, len(specs))
	for _, s := range specs {
		rs = append(rs, s.toRouter())
	}
	cfg := &model.RouteConfiguration{Routes: rs, Dynamic: false}
	return CreateRouterCoordinator(cfg)
}

type res struct {
	ok      bool
	cluster string
	err     string
}

func call(cOld *oldrouter.RouterCoordinator, cNew *RouterCoordinator, method, path string, hdr map[string]string) (res, res) {
	req, _ := stdHttp.NewRequest(method, path, nil)
	httpContext := http.HttpContext{
		Request: req,
	}
	for k, v := range hdr {
		req.Header.Set(k, v)
	}
	oa, oe := cOld.Route(&httpContext)
	na, ne := cNew.Route(&httpContext)

	or := res{}
	if oe != nil || oa == nil {
		if oe != nil {
			or.err = oe.Error()
		}
	} else {
		or.ok = true
		or.cluster = oa.Cluster
	}

	nr := res{}
	if ne != nil || na == nil {
		if ne != nil {
			nr.err = ne.Error()
		}
	} else {
		nr.ok = true
		nr.cluster = na.Cluster
	}
	return or, nr
}

func assertSame(t *testing.T, oldc *oldrouter.RouterCoordinator, newc *RouterCoordinator,
	method, path string, hdr map[string]string, wantOK bool, wantCluster string) {

	ro, rn := call(oldc, newc, method, path, hdr)
	if ro.ok != rn.ok || ro.cluster != rn.cluster {
		t.Fatalf("mismatch: %s %s hdr=%v\n old={ok:%v cluster:%q err:%q}\n new={ok:%v cluster:%q err:%q}",
			method, path, hdr, ro.ok, ro.cluster, ro.err, rn.ok, rn.cluster, rn.err)
	}
	if ro.ok != wantOK || rn.ok != wantOK {
		t.Fatalf("ok mismatch: %s %s hdr=%v wantOK=%v oldOK=%v newOK=%v", method, path, hdr, wantOK, ro.ok, rn.ok)
	}
	if wantOK && wantCluster != "" && (ro.cluster != wantCluster || rn.cluster != wantCluster) {
		t.Fatalf("cluster mismatch: %s %s hdr=%v want=%q old=%q new=%q", method, path, hdr, wantCluster, ro.cluster, rn.cluster)
	}
}

type varSyntax struct {
	simplePattern func(seg string) string         // /users/:id
	digitsPattern func(seg string) (string, bool) // /users/:id(\d+)
	multiPattern  func(a, b string) string        // /shops/:a/orders/:b
}

func colonSyntax() varSyntax {
	return varSyntax{
		simplePattern: func(seg string) string {
			return "/users/:" + seg
		},
		digitsPattern: func(seg string) (string, bool) {
			return "/users/:" + seg + "(\\d+)", true
		},
		multiPattern: func(a, b string) string {
			return "/shops/:" + a + "/orders/:" + b
		},
	}
}

/* ==============================
   test cases (var/regex/priority/header/)
   ============================== */

func TestParitySimpleCases(t *testing.T) {
	syntax = colonSyntax()

	specs := []RouteSpec{
		// exact
		{ID: "exact", Methods: []string{"GET"}, Path: "/api/v1/item/100", Cluster: "c-exact"},
		// prefix（/**）
		{ID: "pre", Methods: []string{"GET"}, Prefix: "/api/v1/svc/", Cluster: "c-pre"},
		// var
		{ID: "var", Methods: []string{"GET"}, Path: syntax.simplePattern("id"), Cluster: "c-var"},
		// multi
		{ID: "multi", Methods: []string{"GET", "POST"}, Path: "/multi", Cluster: "c-multi"},
		// Header regex
		{ID: "hdr", Methods: []string{"GET"}, Headers: []HeaderSpec{{Name: "X-Env", Values: []string{"^prod|staging$"}, Regex: true}}, Cluster: "c-hdr"},
	}

	oldc := buildOld(specs)
	newc := BuildNew(specs)

	cases := []struct {
		name    string
		method  string
		path    string
		hdr     map[string]string
		ok      bool
		cluster string
	}{
		{"exact", "GET", "/api/v1/item/100", nil, true, "c-exact"},
		{"prefix.deep", "GET", "/api/v1/svc/a/b", nil, true, "c-pre"},
		{"var.hit", "GET", "/users/42", nil, true, "c-var"},
		{"var.not_deeper", "GET", "/users/42/extra", nil, false, ""},
		{"multi.get", "GET", "/multi", nil, true, "c-multi"},
		{"multi.post", "POST", "/multi", nil, true, "c-multi"},
		{"hdr.regex", "GET", "/whatever", map[string]string{"X-Env": "prod"}, true, "c-hdr"},
		{"miss", "GET", "/no/match", nil, false, ""},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assertSame(t, oldc, newc, tc.method, tc.path, tc.hdr, tc.ok, tc.cluster)
		})
	}
}

func TestPrioritySpecificOverWildcard(t *testing.T) {
	syntax = colonSyntax()

	specs := []RouteSpec{
		{ID: "wild", Methods: []string{"GET"}, Prefix: "/api/v1/**", Cluster: "c-wild"},
		{ID: "spec", Methods: []string{"GET"}, Path: "/api/v1/test-dubbo/user/name/" + syntax.simplePattern("name")[len("/users/"):], Cluster: "c-spec"},
		// equals to /api/v1/test-dubbo/user/name/:name
	}
	oldc := buildOld(specs)
	newc := BuildNew(specs)

	assertSame(t, oldc, newc, "GET",
		"/api/v1/test-dubbo/user/name/yqxu", nil, true, "c-spec")
}

func TestPriorityDeeperWins(t *testing.T) {
	specs := []RouteSpec{
		{ID: "shallow", Methods: []string{"GET"}, Prefix: "/api/v1/", Cluster: "c-shallow"},
		{ID: "deeper", Methods: []string{"GET"}, Prefix: "/api/v1/test-dubbo/", Cluster: "c-deeper"},
	}
	oldc := buildOld(specs)
	newc := BuildNew(specs)

	assertSame(t, oldc, newc, "GET",
		"/api/v1/test-dubbo/user/name/abc", nil, true, "c-deeper")
}

func TestPrioritySingleStarOverDoubleStar(t *testing.T) {
	// use var to express "/*"
	syntax = colonSyntax()
	specs := []RouteSpec{
		{ID: "multi", Methods: []string{"GET"}, Prefix: "/api/", Cluster: "c-**"},
		{ID: "single", Methods: []string{"GET"}, Path: "/api/" + syntax.simplePattern("seg")[len("/users/"):] + "/users", Cluster: "c-*"},
		// equals to /api/:seg/users
	}
	oldc := buildOld(specs)
	newc := BuildNew(specs)

	assertSame(t, oldc, newc, "GET", "/api/v1/users", nil, true, "c-*")
	assertSame(t, oldc, newc, "GET", "/api/v1/x/users", nil, true, "c-**")
}

func TestVariablesSingleAndMulti(t *testing.T) {
	syntax = colonSyntax()
	specs := []RouteSpec{
		{ID: "one", Methods: []string{"GET"}, Path: syntax.simplePattern("id"), Cluster: "c-one"},
		{ID: "two", Methods: []string{"GET"}, Path: syntax.multiPattern("shopId", "orderId"), Cluster: "c-two"},
		{ID: "pre", Methods: []string{"GET"}, Prefix: "/shops/", Cluster: "c-pre"},
	}
	oldc := buildOld(specs)
	newc := BuildNew(specs)

	assertSame(t, oldc, newc, "GET", "/users/777", nil, true, "c-one")
	assertSame(t, oldc, newc, "GET", syntax.multiPattern("12", "34"), nil, true, "c-two")
	assertSame(t, oldc, newc, "GET", syntax.multiPattern("12", "34")+"/extra", nil, true, "c-pre")
}

func TestHeaderRegexWithRoutes(t *testing.T) {
	specs := []RouteSpec{
		{ID: "hdr", Methods: []string{"GET"}, Headers: []HeaderSpec{{Name: "X-Env", Values: []string{"^prod|staging$"}, Regex: true}}, Cluster: "c-hdr"},
		{ID: "pre", Methods: []string{"GET"}, Prefix: "/api/", Cluster: "c-pre"},
	}
	oldc := buildOld(specs)
	newc := BuildNew(specs)

	assertSame(t, oldc, newc, "GET", "/whatever", map[string]string{"X-Env": "prod"}, true, "c-hdr")
	assertSame(t, oldc, newc, "GET", "/api/foo", map[string]string{"X-Env": "dev"}, true, "c-pre")
}

/* ==============================
   random data fuzz test
   ============================== */

func TestParityRandomized(t *testing.T) {
	syntax = colonSyntax()
	const (
		nRoutes           = 20000
		nRequests         = 10000
		prefixRatio       = 0.40
		headerRatio       = 0.10
		seed        int64 = 20250929
	)

	specs := genRandomSpecsWithVars(syntax, nRoutes, prefixRatio, headerRatio, seed)
	oldc := buildOld(specs)
	newc := BuildNew(specs)

	reqs := genRandomRequests(nRequests, seed+1)
	for i, req := range reqs {
		ro, rn := call(oldc, newc, req.Method, req.URL.Path, headerFromReq(req))
		if ro.ok != rn.ok || ro.cluster != rn.cluster {
			t.Fatalf("Randomized mismatch at #%d: %s %s old={ok:%v cluster:%q err:%q} new={ok:%v cluster:%q err:%q}",
				i, req.Method, req.URL.Path, ro.ok, ro.cluster, ro.err, rn.ok, rn.cluster, rn.err)
		}
	}
}

func headerFromReq(r *stdHttp.Request) map[string]string {
	if len(r.Header) == 0 {
		return nil
	}
	out := make(map[string]string)
	for k, vs := range r.Header {
		if len(vs) > 0 {
			out[k] = vs[0]
		}
	}
	return out
}

/* ==============================
   random data generation tools
   ============================== */

var syntax varSyntax

func genRandomSpecsWithVars(s varSyntax, n int, prefixRatio, headerOnlyRatio float64, seed int64) []RouteSpec {
	rnd := rand.New(rand.NewSource(seed))
	out := make([]RouteSpec, 0, n)

	nHeader := int(float64(n) * headerOnlyRatio)
	nPrefix := int(float64(n-nHeader) * prefixRatio)
	// preserve 20% for "variable path", the rest for exact path
	nVars := int(float64(n) * 0.20)
	nPath := n - nHeader - nPrefix - nVars
	if nPath < 0 {
		nPath = 0
	}

	// Header-only (regex + normal)
	for i := 0; i < nHeader; i++ {
		if i%5 == 0 {
			out = append(out, RouteSpec{
				ID:      "hdrx-" + strconv.Itoa(i),
				Methods: []string{"GET", "POST"},
				Headers: []HeaderSpec{{Name: "X-Trace", Values: []string{"^pixiu-[0-9a-f]{8}$"}, Regex: true}},
				Cluster: "c-hx-" + strconv.Itoa(i),
			})
		} else {
			out = append(out, RouteSpec{
				ID:      "hdr-" + strconv.Itoa(i),
				Methods: []string{"GET", "POST"},
				Headers: []HeaderSpec{{Name: "X-Env", Values: []string{"prod"}, Regex: false}},
				Cluster: "c-h-" + strconv.Itoa(i),
			})
		}
	}

	// Prefix
	for i := 0; i < nPrefix; i++ {
		base := "/api/v" + strconv.Itoa(1+rnd.Intn(3)) + "/svc" + strconv.Itoa(rnd.Intn(50)) + "/" // NOSONAR
		out = append(out, RouteSpec{
			ID:      "pre-" + strconv.Itoa(i),
			Methods: []string{"GET", "POST"},
			Prefix:  base,
			Cluster: "c-p-" + strconv.Itoa(i),
		})
	}

	// Variables
	for i := 0; i < nVars; i++ {
		if i%3 == 0 {
			out = append(out, RouteSpec{
				ID:      "var-" + strconv.Itoa(i),
				Methods: []string{"GET"},
				Path:    s.simplePattern("id"),
				Cluster: "c-v-" + strconv.Itoa(i),
			})
		} else {
			out = append(out, RouteSpec{
				ID:      "var2-" + strconv.Itoa(i),
				Methods: []string{"GET"},
				Path:    s.multiPattern("a", "b"),
				Cluster: "c-v2-" + strconv.Itoa(i),
			})
		}
	}

	// Exact Path
	for i := 0; i < nPath; i++ {
		out = append(out, RouteSpec{
			ID:      "pth-" + strconv.Itoa(i),
			Methods: []string{"GET"},
			Path:    "/api/v1/item/" + strconv.Itoa(i),
			Cluster: "c-x-" + strconv.Itoa(i),
		})
	}
	return out
}

func genRandomRequests(n int, seed int64) []*stdHttp.Request {
	rnd := rand.New(rand.NewSource(seed))
	reqs := make([]*stdHttp.Request, 0, n)
	methods := []string{"GET", "POST"}

	for i := 0; i < n; i++ {
		var path string
		switch rnd.Intn(5) { // NOSONAR
		case 0: // exact style
			path = "/api/v1/item/" + strconv.Itoa(rnd.Intn(50000)) // NOSONAR
		case 1: // prefix style
			path = "/api/v" + strconv.Itoa(1+rnd.Intn(3)) + "/svc" + strconv.Itoa(rnd.Intn(50)) + "/foo/bar" // NOSONAR
		case 2: // var
			path = "/users/" + strconv.Itoa(1000+rnd.Intn(9000)) // NOSONAR
		case 3: // var
			path = "/shops/" + strconv.Itoa(rnd.Intn(100)) + "/orders/" + strconv.Itoa(rnd.Intn(1000)) // NOSONAR
		default:
			path = "/unknown/" + strconv.Itoa(rnd.Intn(100000)) // NOSONAR
		}
		req, _ := stdHttp.NewRequest(methods[rnd.Intn(len(methods))], path, nil) // NOSONAR
		// header-only
		switch rnd.Intn(7) { // NOSONAR
		case 0:
			req.Header.Set("X-Env", "prod")
		case 1:
			req.Header.Set("X-Trace", "pixiu-"+strconv.FormatInt(rnd.Int63()&0xffffffff, 16)) // NOSONAR
		}
		reqs = append(reqs, req)
	}
	return reqs
}

func newTestRouter(id, method, path, cluster string) *model.Router {
	return &model.Router{
		ID: id,
		Match: model.RouterMatch{
			Methods: []string{method},
			Path:    path,
		},
		Route: model.RouteAction{
			Cluster: cluster,
		},
	}
}

func mustRouteCluster(t *testing.T, rm *RouterCoordinator, method, path string, want string) {
	t.Helper()
	req, err := stdHttp.NewRequest(method, path, nil)
	if err != nil {
		t.Fatalf("failed to build request: %v", err)
	}
	ctx := http.HttpContext{Request: req}
	act, err := rm.Route(&ctx)
	if err != nil {
		t.Fatalf("unexpected route error for %s %s: %v", method, path, err)
	}
	if act == nil {
		t.Fatalf("route action is nil for %s %s", method, path)
	}
	if act.Cluster != want {
		t.Fatalf("unexpected cluster for %s %s, want %q, got %q", method, path, want, act.Cluster)
	}
}

func mustRouteNotFound(t *testing.T, rm *RouterCoordinator, method, path string) {
	t.Helper()
	req, err := stdHttp.NewRequest(method, path, nil)
	if err != nil {
		t.Fatalf("failed to build request: %v", err)
	}
	ctx := http.HttpContext{Request: req}
	act, err := rm.Route(&ctx)
	if err == nil {
		t.Fatalf("expected error for %s %s, got nil (action=%#v)", method, path, act)
	}
}

func mustRouteByPathAndNameCluster(t *testing.T, rm *RouterCoordinator, method, path string, want string) {
	t.Helper()
	act, err := rm.RouteByPathAndName(path, method)
	if err != nil {
		t.Fatalf("unexpected RouteByPathAndName error for %s %s: %v", method, path, err)
	}
	if act == nil {
		t.Fatalf("route action is nil for %s %s", method, path)
	}
	if act.Cluster != want {
		t.Fatalf("unexpected cluster for %s %s, want %q, got %q", method, path, want, act.Cluster)
	}
}

func mustRouteByPathAndNameNotFound(t *testing.T, rm *RouterCoordinator, method, path string) {
	t.Helper()
	act, err := rm.RouteByPathAndName(path, method)
	if err == nil {
		t.Fatalf("expected RouteByPathAndName error for %s %s, got nil (action=%#v)", method, path, act)
	}
}

func TestCreateRouterCoordinatorInitialSnapshot(t *testing.T) {
	r1 := newTestRouter("r1", "GET", "/foo", "cluster-foo")
	r2 := newTestRouter("r2", "GET", "/bar", "cluster-bar")

	cfg := &model.RouteConfiguration{
		Routes:  []*model.Router{r1, r2},
		Dynamic: false,
	}

	rm := CreateRouterCoordinator(cfg)
	rm.debounce = 0

	mustRouteCluster(t, rm, "GET", "/foo", "cluster-foo")
	mustRouteCluster(t, rm, "GET", "/bar", "cluster-bar")

	mustRouteByPathAndNameCluster(t, rm, "GET", "/foo", "cluster-foo")
	mustRouteByPathAndNameCluster(t, rm, "GET", "/bar", "cluster-bar")
}

// TestRouterCoordinatorOnAddRouterUpdatesSnapshot 验证 Add 后新路由生效，旧路由保持
func TestRouterCoordinatorOnAddRouterUpdatesSnapshot(t *testing.T) {
	r1 := newTestRouter("r1", "GET", "/foo", "cluster-foo")

	cfg := &model.RouteConfiguration{
		Routes:  []*model.Router{r1},
		Dynamic: false,
	}

	rm := CreateRouterCoordinator(cfg)
	rm.debounce = 0

	// only /foo
	mustRouteCluster(t, rm, "GET", "/foo", "cluster-foo")
	mustRouteNotFound(t, rm, "GET", "/bar")
	mustRouteByPathAndNameCluster(t, rm, "GET", "/foo", "cluster-foo")
	mustRouteByPathAndNameNotFound(t, rm, "GET", "/bar")

	// Add /bar
	r2 := newTestRouter("r2", "GET", "/bar", "cluster-bar")
	rm.OnAddRouter(r2)

	// /foo available
	mustRouteCluster(t, rm, "GET", "/foo", "cluster-foo")
	mustRouteByPathAndNameCluster(t, rm, "GET", "/foo", "cluster-foo")

	// /bar available
	mustRouteCluster(t, rm, "GET", "/bar", "cluster-bar")
	mustRouteByPathAndNameCluster(t, rm, "GET", "/bar", "cluster-bar")
}

// TestRouterCoordinatorOnDeleteRouterUpdatesSnapshot 验证 Delete 后路由失效
func TestRouterCoordinatorOnDeleteRouterUpdatesSnapshot(t *testing.T) {
	r1 := newTestRouter("r1", "GET", "/foo", "cluster-foo")
	r2 := newTestRouter("r2", "GET", "/bar", "cluster-bar")

	cfg := &model.RouteConfiguration{
		Routes:  []*model.Router{r1, r2},
		Dynamic: false,
	}

	rm := CreateRouterCoordinator(cfg)
	rm.debounce = 0

	// /foo、/bar available
	mustRouteCluster(t, rm, "GET", "/foo", "cluster-foo")
	mustRouteCluster(t, rm, "GET", "/bar", "cluster-bar")
	mustRouteByPathAndNameCluster(t, rm, "GET", "/foo", "cluster-foo")
	mustRouteByPathAndNameCluster(t, rm, "GET", "/bar", "cluster-bar")

	// delete /bar
	rm.OnDeleteRouter(r2)

	// /foo available
	mustRouteCluster(t, rm, "GET", "/foo", "cluster-foo")
	mustRouteByPathAndNameCluster(t, rm, "GET", "/foo", "cluster-foo")

	// /bar fail
	mustRouteNotFound(t, rm, "GET", "/bar")
	mustRouteByPathAndNameNotFound(t, rm, "GET", "/bar")
}

func TestOldVsNew_OnAddRouterWithSameID(t *testing.T) {
	base := []*model.Router{
		{
			ID: "r1",
			Match: model.RouterMatch{
				Methods: []string{"POST"},
				Path:    "/api/v1/item/711",
			},
			Route: model.RouteAction{Cluster: "c-x-r1"},
		},
	}

	oldc := buildOldCoordinator(base)
	newc := buildNewCoordinator(base)
	newc.debounce = 0

	{
		req, _ := stdHttp.NewRequest("POST", "/api/v1/item/711", nil)
		ctxOld := http.HttpContext{Request: req}
		ctxNew := http.HttpContext{Request: req}

		oldRes, oldErr := oldc.Route(&ctxOld)
		newRes, newErr := newc.Route(&ctxNew)

		if oldErr != nil || newErr != nil {
			t.Fatalf("initial route error: oldErr=%v newErr=%v", oldErr, newErr)
		}
		if oldRes.Cluster != "c-x-r1" || newRes.Cluster != "c-x-r1" {
			t.Fatalf("initial cluster mismatch: old=%v new=%v", oldRes, newRes)
		}
	}

	delta := &model.Router{
		ID: "r1",
		Match: model.RouterMatch{
			Methods: []string{"POST"},
			Path:    "/api/v1/item/999999",
		},
		Route: model.RouteAction{Cluster: "c-x-r1"},
	}

	oldc.OnAddRouter(delta)
	newc.OnAddRouter(delta)

	{
		req, _ := stdHttp.NewRequest("POST", "/api/v1/item/711", nil)
		ctxOld := http.HttpContext{Request: req}
		ctxNew := http.HttpContext{Request: req}

		oldRes, oldErr := oldc.Route(&ctxOld)
		newRes, newErr := newc.Route(&ctxNew)

		assert.Equal(t, oldRes, newRes)
		assert.Equal(t, oldErr, newErr)
	}

	{
		req, _ := stdHttp.NewRequest("POST", "/api/v1/item/999999", nil)
		ctxOld := http.HttpContext{Request: req}
		ctxNew := http.HttpContext{Request: req}

		oldRes, oldErr := oldc.Route(&ctxOld)
		newRes, newErr := newc.Route(&ctxNew)

		assert.Equal(t, oldRes, newRes)
		assert.Equal(t, oldErr, newErr)
	}
}
