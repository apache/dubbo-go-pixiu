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
	"github.com/dubbo-go-pixiu/pixiu-api/pkg/api/config"
	"github.com/dubbo-go-pixiu/pixiu-api/pkg/router"
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/router/trie"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/util/stringutil"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
)

// Node defines the single method of the router configured API
type Node struct {
	fullPath string
	filters  []string
	method   *config.Method
	headers  map[string]string
}

// Route defines the tree of router APIs
type Route struct {
	lock sync.RWMutex
	tree trie.Trie
}

// ClearAPI clear the api
func (rt *Route) ClearAPI() error {
	rt.lock.Lock()
	defer rt.lock.Unlock()
	rt.tree.Clear()
	return nil
}

func (rt *Route) RemoveAPI(api router.API) {
	lowerCasePath := strings.ToLower(api.URLPattern)
	key := getTrieKey(api.Method.HTTPVerb, lowerCasePath, false)

	rt.lock.Lock()
	defer rt.lock.Unlock()
	// if api is exists
	if exists, err := rt.tree.Contains(key); err != nil {
		logger.Errorf("rt.tree.Get(key) err: %s", err.Error())
		return
	} else if exists {
		// append new cluster to the node
		if node, _, exists, err := rt.tree.Get(key); err != nil {
			logger.Errorf("rt.tree.Get(key) err: %s", err.Error())
			return
		} else if exists {
			bizInfoInterface := node.GetBizInfo()
			bizInfo, ok := bizInfoInterface.(*Node)
			if bizInfo == nil || !ok {
				logger.Error("bizInfoInterface.(*Node) failed")
				return
			}
			// avoid thread safe problem
			clusters := strings.Split(bizInfo.method.HTTPBackendConfig.URL, ",")
			if len(clusters) > 1 {
				var i int
				for i = 0; i < len(clusters); i++ {
					if clusters[i] == api.URL {
						break
					}
				}

				if i == len(clusters) {
					return
				}

				// operate pointer has no necessary to call update api
				bizInfo.method.HTTPBackendConfig.URL = strings.Join(append(clusters[:i], clusters[i+1:]...), ",")
			} else {
				// double check, avoid removing api that does not exists
				if !strings.Contains(bizInfo.method.HTTPBackendConfig.URL, api.URL) {
					return
				}
				// if backend has only one node, then just directly remove
				_, _ = rt.tree.Remove(key)
			}
		}
	}
}

func getTrieKey(method config.HTTPVerb, path string, isPrefix bool) string {
	if isPrefix {
		if !strings.HasSuffix(path, constant.PathSlash) {
			path = path + constant.PathSlash
		}
		path = path + "**"
	}
	return stringutil.GetTrieKey(string(method), path)
}

// PutAPI puts an api into the resource
func (rt *Route) PutAPI(api router.API) error {
	lowerCasePath := strings.ToLower(api.URLPattern)
	key := getTrieKey(api.Method.HTTPVerb, lowerCasePath, false)
	node, ok := rt.getNode(key)
	if !ok {
		rn := &Node{
			fullPath: lowerCasePath,
			method:   &api.Method,
			headers:  api.Headers,
		}
		rt.lock.Lock()
		defer rt.lock.Unlock()
		_, _ = rt.tree.Put(key, rn)
		return nil
	}
	return errors.Errorf("Method %s with address %s already exists in path %s",
		api.Method.HTTPVerb, lowerCasePath, node.fullPath)
}

// PutOrUpdateAPI puts or updates an api into the resource
func (rt *Route) PutOrUpdateAPI(api router.API) error {
	lowerCasePath := strings.ToLower(api.URLPattern)
	key := getTrieKey(api.Method.HTTPVerb, lowerCasePath, false)
	rn := &Node{
		fullPath: lowerCasePath,
		method:   &api.Method,
		headers:  api.Headers,
	}
	rt.lock.Lock()
	defer rt.lock.Unlock()

	// if api is exists
	if exists, err := rt.tree.Contains(key); err != nil {
		return err
	} else if exists {
		// append new cluster to the node
		if node, _, exists, err := rt.tree.Get(key); err != nil {
			return err
		} else if exists {
			bizInfoInterface := node.GetBizInfo()
			bizInfo, ok := bizInfoInterface.(*Node)
			if bizInfo == nil || !ok {
				return errors.New("bizInfoInterface.(*Node) failed")
			}
			if !strings.Contains(bizInfo.method.HTTPBackendConfig.URL, api.URL) {
				// operate pointer has no necessary to call update api
				bizInfo.method.HTTPBackendConfig.URL = bizInfo.method.HTTPBackendConfig.URL + "," + api.URL
			}
		}
	} else {
		if ok, err := rt.tree.PutOrUpdate(key, rn); !ok {
			return err
		}
	}

	return nil
}

// FindAPI return if api has path in trie,or nil
func (rt *Route) FindAPI(fullPath string, httpverb config.HTTPVerb) (*router.API, bool) {
	lowerCasePath := strings.ToLower(fullPath)
	key := getTrieKey(httpverb, lowerCasePath, false)
	if n, found := rt.getNode(key); found {
		rt.lock.RLock()
		defer rt.lock.RUnlock()
		return &router.API{
			URLPattern: n.fullPath,
			Method:     *n.method,
			Headers:    n.headers,
		}, found
	}
	return nil, false
}

// MatchAPI FindAPI returns the api that meets the rule
func (rt *Route) MatchAPI(fullPath string, httpverb config.HTTPVerb) (*router.API, bool) {
	lowerCasePath := strings.ToLower(fullPath)
	key := getTrieKey(httpverb, lowerCasePath, false)
	if n, found := rt.matchNode(key); found {
		rt.lock.RLock()
		defer rt.lock.RUnlock()
		return &router.API{
			URLPattern: n.fullPath,
			Method:     *n.method,
			Headers:    n.headers,
		}, found
	}
	return nil, false
}

// DeleteNode delete node by fullPath
func (rt *Route) DeleteNode(fullPath string) bool {
	rt.lock.RLock()
	defer rt.lock.RUnlock()
	methodList := [8]config.HTTPVerb{"ANY", "GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	for _, v := range methodList {
		key := getTrieKey(v, fullPath, false)
		_, _ = rt.tree.Remove(key)
	}
	return true
}

// DeleteAPI delete api by fullPath and http verb
func (rt *Route) DeleteAPI(fullPath string, httpverb config.HTTPVerb) bool {
	lowerCasePath := strings.ToLower(fullPath)
	key := getTrieKey(httpverb, lowerCasePath, false)
	if _, found := rt.getNode(key); found {
		rt.lock.RLock()
		defer rt.lock.RUnlock()
		_, _ = rt.tree.Remove(key)
		return true
	}
	return false
}

func (rt *Route) getNode(fullPath string) (*Node, bool) {
	var n interface{}
	var found bool
	rt.lock.RLock()
	defer rt.lock.RUnlock()
	trieNode, _, _, _ := rt.tree.Get(fullPath)
	found = trieNode != nil
	if !found {
		return nil, false
	}
	n = trieNode.GetBizInfo()
	if n == nil {
		return nil, false
	}
	return n.(*Node), found
}

func (rt *Route) matchNode(fullPath string) (*Node, bool) {
	var n interface{}
	var found bool
	rt.lock.RLock()
	defer rt.lock.RUnlock()
	trieNode, _, _ := rt.tree.Match(fullPath)
	found = trieNode != nil
	if !found {
		return nil, false
	}
	n = trieNode.GetBizInfo()
	if n == nil {
		return nil, false
	}
	return n.(*Node), found
}

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
		tree: trie.NewTrie(),
	}
}
