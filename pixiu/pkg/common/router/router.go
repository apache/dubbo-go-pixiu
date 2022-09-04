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
	"strings"
	"sync"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/router/trie"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/util/stringutil"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/server"
)

type (
	// RouterCoordinator the router coordinator for http connection manager
	RouterCoordinator struct {
		activeConfig *model.RouteConfiguration
		rw           sync.RWMutex
	}
)

// CreateRouterCoordinator create coordinator for http connection manager
func CreateRouterCoordinator(routeConfig *model.RouteConfiguration) *RouterCoordinator {
	rc := &RouterCoordinator{activeConfig: routeConfig}
	if routeConfig.Dynamic {
		server.GetRouterManager().AddRouterListener(rc)
	}
	rc.initTrie()
	return rc
}

// Route find routeAction for request
func (rm *RouterCoordinator) Route(hc *http.HttpContext) (*model.RouteAction, error) {
	rm.rw.RLock()
	defer rm.rw.RUnlock()

	return rm.activeConfig.Route(hc.Request)
}

func (rm *RouterCoordinator) RouteByPathAndName(path, method string) (*model.RouteAction, error) {
	rm.rw.RLock()
	defer rm.rw.RUnlock()

	return rm.activeConfig.RouteByPathAndMethod(path, method)
}

func getTrieKey(method string, path string, isPrefix bool) string {
	if isPrefix {
		if !strings.HasSuffix(path, constant.PathSlash) {
			path = path + constant.PathSlash
		}
		path = path + "**"
	}
	return stringutil.GetTrieKey(method, path)
}

func (rm *RouterCoordinator) initTrie() {
	if rm.activeConfig.RouteTrie.IsEmpty() {
		rm.activeConfig.RouteTrie = trie.NewTrie()
	}
	for _, router := range rm.activeConfig.Routes {
		rm.OnAddRouter(router)
	}
}

// OnAddRouter add router
func (rm *RouterCoordinator) OnAddRouter(r *model.Router) {
	//TODO: lock move to trie node
	rm.rw.Lock()
	defer rm.rw.Unlock()
	if r.Match.Methods == nil {
		r.Match.Methods = []string{constant.Get, constant.Put, constant.Delete, constant.Post}
	}
	isPrefix := r.Match.Prefix != ""
	for _, method := range r.Match.Methods {
		var key string
		if isPrefix {
			key = getTrieKey(method, r.Match.Prefix, isPrefix)
		} else {
			key = getTrieKey(method, r.Match.Path, isPrefix)
		}
		_, _ = rm.activeConfig.RouteTrie.Put(key, r.Route)
	}
}

// OnDeleteRouter delete router
func (rm *RouterCoordinator) OnDeleteRouter(r *model.Router) {
	rm.rw.Lock()
	defer rm.rw.Unlock()

	if r.Match.Methods == nil {
		r.Match.Methods = []string{constant.Get, constant.Put, constant.Delete, constant.Post}
	}
	isPrefix := r.Match.Prefix != ""
	for _, method := range r.Match.Methods {
		var key string
		if isPrefix {
			key = getTrieKey(method, r.Match.Prefix, isPrefix)
		} else {
			key = getTrieKey(method, r.Match.Path, isPrefix)
		}
		_, _ = rm.activeConfig.RouteTrie.Remove(key)
	}
}
