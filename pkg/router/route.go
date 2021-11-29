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
	"net/url"
	"strings"
	"sync"
)

import (
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/api/config"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/router"

	"github.com/emirpasic/gods/trees/avltree"

	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
)

// Node defines the single method of the router configured API
type Node struct {
	fullPath string
	wildcard bool
	filters  []string
	methods  map[config.HTTPVerb]*config.Method
	headers  map[string]string
}

// Route defines the tree of router APIs
type Route struct {
	lock         sync.RWMutex
	tree         *avltree.Tree
	wildcardTree *avltree.Tree
}

// ClearAPI clear the api
func (rt *Route) ClearAPI() error {
	rt.lock.Lock()
	defer rt.lock.Unlock()
	rt.wildcardTree.Clear()
	rt.tree.Clear()
	return nil
}

func (r *Route) RemoveAPI(api router.API) {
	fullPath := api.URLPattern
	node, ok := r.findNode(fullPath)
	if !ok {
		return
	}
	if tempMethod, ok := node.methods[api.HTTPVerb]; ok {
		splitedURLs := strings.Split(tempMethod.IntegrationRequest.HTTPBackendConfig.URL, ",")
		afterRemoveedURL := make([]string, 0, len(splitedURLs))
		for _, v := range splitedURLs {
			if v != api.IntegrationRequest.HTTPBackendConfig.URL {
				afterRemoveedURL = append(afterRemoveedURL, v)
			}
		}
		if len(afterRemoveedURL) == 0 {
			delete(node.methods, api.HTTPVerb)
		}
		node.methods[api.HTTPVerb].IntegrationRequest.HTTPBackendConfig.URL = strings.Join(afterRemoveedURL, ",")
		return
	}
}

// PutAPI puts an api into the resource
func (rt *Route) PutAPI(api router.API) error {
	fullPath := api.URLPattern
	node, ok := rt.findNode(fullPath)
	rt.lock.Lock()
	defer rt.lock.Unlock()
	if !ok {
		wildcard := strings.Contains(fullPath, constant.PathParamIdentifier)
		rn := &Node{
			fullPath: fullPath,
			methods:  map[config.HTTPVerb]*config.Method{api.Method.HTTPVerb: &api.Method},
			wildcard: wildcard,
			headers:  api.Headers,
		}
		if wildcard {
			rt.wildcardTree.Put(fullPath, rn)
		}
		rt.tree.Put(fullPath, rn)
		return nil
	}
	return node.putMethod(api.Method, api.Headers)
}

func (node *Node) putMethod(method config.Method, headers map[string]string) error {
	// todo lock
	if tempMethod, ok := node.methods[method.HTTPVerb]; ok {
		splitedURLs := strings.Split(tempMethod.IntegrationRequest.HTTPBackendConfig.URL, ",")
		for _, v := range splitedURLs {
			if v == method.IntegrationRequest.HTTPBackendConfig.URL {
				return errors.Errorf("Method %s with address %s already exists in path %s",
					method.HTTPVerb, v, node.fullPath)
			}
		}
		splitedURLs = append(splitedURLs, method.IntegrationRequest.HTTPBackendConfig.URL)
		node.methods[method.HTTPVerb].IntegrationRequest.HTTPBackendConfig.URL = strings.Join(splitedURLs, ",")
		node.headers = headers
		return nil
	}
	node.methods[method.HTTPVerb] = &method
	node.headers = headers
	return nil
}

// UpdateAPI update the api method in the existing router node
func (rt *Route) UpdateAPI(api router.API) error {
	node, found := rt.findNode(api.URLPattern)
	if found {
		if _, ok := node.methods[api.Method.HTTPVerb]; ok {
			rt.lock.Lock()
			defer rt.lock.Unlock()
			node.methods[api.Method.HTTPVerb] = &api.Method
		}
	}
	return nil
}

// FindAPI returns the api that meets the
func (rt *Route) FindAPI(fullPath string, httpverb config.HTTPVerb) (*router.API, bool) {
	if n, found := rt.findNode(fullPath); found {
		rt.lock.RLock()
		defer rt.lock.RUnlock()
		if method, ok := n.methods[httpverb]; ok {
			return &router.API{
				URLPattern: n.fullPath,
				Method:     *method,
				Headers:    n.headers,
			}, ok
		}
	}
	return nil, false
}

// DeleteNode delete node by fullPath
func (rt *Route) DeleteNode(fullPath string) bool {
	rt.lock.RLock()
	defer rt.lock.RUnlock()
	rt.tree.Remove(fullPath)
	return true
}

// DeleteAPI delete api by fullPath and http verb
func (rt *Route) DeleteAPI(fullPath string, httpverb config.HTTPVerb) bool {
	if n, found := rt.findNode(fullPath); found {
		rt.lock.RLock()
		defer rt.lock.RUnlock()
		delete(n.methods, httpverb)
		return true
	}
	return false
}

func (rt *Route) findNode(fullPath string) (*Node, bool) {
	var n interface{}
	var found bool
	if n, found = rt.searchWildcard(fullPath); !found {
		rt.lock.RLock()
		defer rt.lock.RUnlock()
		if n, found = rt.tree.Get(fullPath); !found {
			return nil, false
		}
	}
	return n.(*Node), found
}

func (rt *Route) searchWildcard(fullPath string) (*Node, bool) {
	rt.lock.RLock()
	defer rt.lock.RUnlock()
	wildcardPaths := rt.wildcardTree.Keys()
	for _, p := range wildcardPaths {
		if wildcardMatch(p.(string), fullPath) != nil {
			n, ok := rt.wildcardTree.Get(p)
			return n.(*Node), ok
		}
	}
	return nil, false
}

// wildcardMatch validate if the checkPath meets the wildcardPath,
// for example /vought/12345 should match wildcard path /vought/:id;
// /vought/1234abcd/status should not match /vought/:id;
func wildcardMatch(wildcardPath string, checkPath string) url.Values {
	cPaths := strings.Split(strings.TrimLeft(checkPath, constant.PathSlash), constant.PathSlash)
	wPaths := strings.Split(strings.TrimLeft(wildcardPath, constant.PathSlash), constant.PathSlash)
	result := url.Values{}
	if len(cPaths) == 0 || len(wPaths) == 0 || len(cPaths) != len(wPaths) {
		return nil
	}
	for i := 0; i < len(cPaths); i++ {
		if !strings.EqualFold(cPaths[i], wPaths[i]) && !strings.HasPrefix(wPaths[i], constant.PathParamIdentifier) {
			return nil
		}
		if strings.HasPrefix(wPaths[i], constant.PathParamIdentifier) {
			result.Add(strings.TrimPrefix(wPaths[i], constant.PathParamIdentifier), cPaths[i])
		}
	}
	return result
}

// NewRoute returns an empty router tree
func NewRoute() *Route {
	return &Route{
		tree:         avltree.NewWithStringComparator(),
		wildcardTree: avltree.NewWithStringComparator(),
	}
}
