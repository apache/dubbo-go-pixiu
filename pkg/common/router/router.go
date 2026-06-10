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
	stdHttp "net/http"
	"slices"
	"sync"
	"sync/atomic"
	"time"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/router/trie"
	"github.com/apache/dubbo-go-pixiu/pkg/common/util/stringutil"
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
)

// RouterCoordinator the router coordinator for http connection manager
type RouterCoordinator struct {
	mainSnapshot atomic.Pointer[model.RouteSnapshot] // atomic snapshot
	mu           sync.Mutex

	nextSnapshot []*model.Router // temp store for dynamic update, DO NOT read directly

	timer    *time.Timer   // debounce timer
	debounce time.Duration // merge window, default 50ms

	needsRegistration atomic.Bool // whether needs to register as RouterListener
	dynamic           bool        // whether dynamic routing is enabled
}

// CreateRouterCoordinator create coordinator for http connection manager
func CreateRouterCoordinator(routeConfig *model.RouteConfiguration) *RouterCoordinator {
	rc := &RouterCoordinator{
		nextSnapshot: make([]*model.Router, 0, len(routeConfig.Routes)),
		debounce:     50 * time.Millisecond, // merge window
		dynamic:      routeConfig.Dynamic,
	}
	if routeConfig.Dynamic {
		rm := server.GetRouterManager()
		if rm != nil {
			rm.AddRouterListener(rc)
		} else {
			// RouterManager not initialized yet, will register later
			rc.needsRegistration.Store(true)
		}
	}
	// build initial config and store snapshot
	rc.mainSnapshot.Store(model.ToSnapshot(buildRouteConfiguration(routeConfig.Routes).Routes))
	// copy initial routes to store, keep origin order
	rc.nextSnapshot = append(rc.nextSnapshot, routeConfig.Routes...)
	return rc
}

func (rm *RouterCoordinator) Close() {
	if rm.dynamic {
		routerMgr := server.GetRouterManager()
		if routerMgr != nil {
			routerMgr.RemoveRouterListener(rm)
		}
	}
}

func (rm *RouterCoordinator) Route(hc *http.HttpContext) (*model.RouteAction, error) {
	if rm.dynamic && rm.needsRegistration.CompareAndSwap(true, false) {
		routerMgr := server.GetRouterManager()
		if routerMgr != nil {
			routerMgr.AddRouterListener(rm)
		} else {
			rm.needsRegistration.Store(true)
		}
	}
	return rm.route(hc.Request)
}

func (rm *RouterCoordinator) RouteByPathAndName(path, method string) (*model.RouteAction, error) {
	s := rm.mainSnapshot.Load()
	if s == nil {
		return nil, errors.New("router configuration is empty")
	}
	key := stringutil.GetTrieKey(method, path)
	// Try the method-specific trie first; fall back to the wildcard ("*") trie
	// so that routes declared with methods: ["*"] keep their "match any method"
	// semantics under the per-method trie layout.
	if act, err := matchInTrie(s.MethodTries[method], key); act != nil || err != nil {
		return act, err
	}
	if method != "*" {
		if act, err := matchInTrie(s.MethodTries["*"], key); act != nil || err != nil {
			return act, err
		}
	}
	return nil, errors.Errorf("route failed for %s, no rules matched", key)
}

func (rm *RouterCoordinator) route(req *stdHttp.Request) (*model.RouteAction, error) {
	s := rm.mainSnapshot.Load()
	if s == nil {
		return nil, errors.New("router configuration is empty")
	}

	// header-only first
	for _, hr := range s.HeaderOnly {
		if !model.MethodAllowed(hr.Methods, req.Method) {
			continue
		}
		if matchHeaders(hr.Headers, req) {
			if len(hr.Action.Cluster) == 0 {
				return nil, errors.New("action is nil. please check your configuration.")
			}
			return &hr.Action, nil
		}
	}

	key := stringutil.GetTrieKey(req.Method, req.URL.Path)
	// Method-specific trie first, then fall back to the wildcard ("*") trie
	// so routes declared with methods: ["*"] continue to match any method.
	if act, err := matchInTrie(s.MethodTries[req.Method], key); act != nil || err != nil {
		return act, err
	}
	if req.Method != "*" {
		if act, err := matchInTrie(s.MethodTries["*"], key); act != nil || err != nil {
			return act, err
		}
	}
	return nil, errors.Errorf("route failed for %s, no rules matched", key)
}

// matchInTrie looks up key in t. The tri-state return lets callers chain a
// fallback trie:
//
//	(action, nil) -> match found, return immediately
//	(nil, error)  -> match found but bizInfo has wrong type, do NOT fall back
//	(nil, nil)    -> miss, caller may try the next trie
func matchInTrie(t *trie.Trie, key string) (*model.RouteAction, error) {
	if t == nil {
		return nil, nil
	}
	node, _, ok := t.Match(key)
	if !ok || node == nil || node.GetBizInfo() == nil {
		return nil, nil
	}
	act, ok := node.GetBizInfo().(model.RouteAction)
	if !ok {
		return nil, errors.Errorf("route failed for %s, invalid route action type", key)
	}
	return &act, nil
}

// reset timer or publish directly
func (rm *RouterCoordinator) schedulePublishLocked() {
	if rm.debounce <= 0 {
		// fallback: immediate
		rm.publishLocked()
		return
	}
	if rm.timer == nil {
		rm.timer = time.NewTimer(rm.debounce)
		go rm.awaitAndPublish()
		return
	}
	// clear timer channel
	if !rm.timer.Stop() {
		select {
		case <-rm.timer.C:
		default:
		}
	}
	rm.timer.Reset(rm.debounce)
}

// wait for timer and publish
func (rm *RouterCoordinator) awaitAndPublish() {
	<-rm.timer.C
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.publishLocked()
	rm.timer = nil
}

// publish: clone from store -> build new config -> atomic switch
func (rm *RouterCoordinator) publishLocked() {
	// 1) clone routes
	next := make([]*model.Router, len(rm.nextSnapshot))
	copy(next, rm.nextSnapshot)
	// 2) build new config
	cfg := buildRouteConfiguration(next)
	// 3) atomic switch
	rm.mainSnapshot.Store(model.ToSnapshot(cfg.Routes))
}

func buildRouteConfiguration(routes []*model.Router) *model.RouteConfiguration {
	cfg := &model.RouteConfiguration{
		RouteTrie: trie.NewTrie(),
		Routes:    make([]*model.Router, 0, len(routes)),
		Dynamic:   false,
	}
	cfg.Routes = append(cfg.Routes, routes...)
	initRegex(cfg)
	fillTrieFromRoutes(cfg)
	return cfg
}

func initRegex(cfg *model.RouteConfiguration) {
	for _, router := range cfg.Routes {
		headers := router.Match.Headers
		for i := range headers {
			if headers[i].Regex && len(headers[i].Values) > 0 {
				if err := headers[i].SetValueRegex(headers[i].Values[0]); err != nil {
					logger.Warnf("invalid regexp in headers[%d]: %v", i, err)
				}
			}
		}
	}
}

// OnAddRouter add router, every call of OnAdd will ADD a new rule, instead of OVERWRITE the same rule
func (rm *RouterCoordinator) OnAddRouter(r *model.Router) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.nextSnapshot = append(rm.nextSnapshot, r)
	rm.schedulePublishLocked()
}

func fillTrieFromRoutes(cfg *model.RouteConfiguration) {
	for _, r := range cfg.Routes {
		methods := r.Match.Methods
		if len(methods) == 0 {
			methods = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"}
		}
		for _, m := range methods {
			key := stringutil.GetTrieKeyWithPrefix(m, r.Match.Path, r.Match.Prefix, r.Match.Prefix != "")
			_, _ = cfg.RouteTrie.Put(key, r.Route)
		}
	}
}

// OnDeleteRouter delete router
func (rm *RouterCoordinator) OnDeleteRouter(r *model.Router) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if len(rm.nextSnapshot) == 0 {
		return
	}
	out := rm.nextSnapshot[:0]
	for _, rr := range rm.nextSnapshot {
		if rr.ID == r.ID {
			continue
		}
		out = append(out, rr)
	}
	rm.nextSnapshot = out
	rm.schedulePublishLocked()
}

func matchHeaders(chs []model.CompiledHeader, r *stdHttp.Request) bool {
	for _, ch := range chs {
		if val := r.Header.Get(ch.Name); len(val) > 0 {
			if ch.Regex != nil {
				if ok := ch.Regex.MatchString(val); !ok {
					return false
				}
				continue
			}

			if !slices.Contains(ch.Values, val) {
				return false
			}
		} else {
			return false
		}
	}
	return true
}
