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
	"github.com/emirpasic/gods/trees/avltree"
	"github.com/pkg/errors"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/constant"
	"github.com/dubbogo/dubbo-go-proxy/pkg/config"
)

// Node defines the single method of the router configured API
type Node struct {
	fullPath string
	wildcard bool
	methods  map[config.HTTPVerb]*config.Method
}

// Route defines the tree of router APIs
type Route struct {
	lock         sync.RWMutex
	tree         *avltree.Tree
	wildcardTree *avltree.Tree
}

// Put put a key val into the tree
func (rt *Route) Put(fullPath string, method config.Method) error {
	lowerCasePath := strings.ToLower(fullPath)
	node, ok := rt.findNode(lowerCasePath)
	rt.lock.Lock()
	defer rt.lock.Unlock()
	if !ok {
		wildcard := strings.Contains(lowerCasePath, constant.PathParamIdentifier)
		rn := &Node{
			fullPath: lowerCasePath,
			methods:  map[config.HTTPVerb]*config.Method{method.HTTPVerb: &method},
			wildcard: wildcard,
		}
		if wildcard {
			rt.wildcardTree.Put(lowerCasePath, rn)
		}
		rt.tree.Put(lowerCasePath, rn)
		return nil
	}
	return node.putMethod(method)
}

// UpdateMethod update the api method in the existing router node
func (rt *Route) UpdateMethod(fullPath string, verb config.HTTPVerb, method config.Method) error {
	node, found := rt.findNode(fullPath)
	if found {
		if _, ok := node.methods[verb]; ok {
			rt.lock.Lock()
			defer rt.lock.Unlock()
			node.methods[verb] = &method
		}
	}
	return nil
}

func (rt *Route) findNode(fullPath string) (*Node, bool) {
	lowerPath := strings.ToLower(fullPath)
	var n interface{}
	var found bool
	if n, found = rt.searchWildcard(lowerPath); !found {
		rt.lock.RLock()
		defer rt.lock.RUnlock()
		if n, found = rt.tree.Get(lowerPath); !found {
			return nil, false
		}
	}
	return n.(*Node), found
}

// FindMethod returns the api that meets the
func (rt *Route) FindMethod(fullPath string, httpverb config.HTTPVerb) (*config.Method, bool) {
	if n, found := rt.findNode(fullPath); found {
		rt.lock.RLock()
		defer rt.lock.RUnlock()
		method, ok := n.methods[httpverb]
		return method, ok
	}
	return nil, false
}

func (rt *Route) searchWildcard(fullPath string) (*Node, bool) {
	rt.lock.RLock()
	defer rt.lock.RUnlock()
	wildcardPaths := rt.wildcardTree.Keys()
	for _, p := range wildcardPaths {
		if wildcardMatch(p.(string), fullPath) {
			n, ok := rt.wildcardTree.Get(p)
			return n.(*Node), ok
		}
	}
	return nil, false
}

func (node *Node) putMethod(method config.Method) error {
	if _, ok := node.methods[method.HTTPVerb]; ok {
		return errors.Errorf("Method %s already exists in path %s", method.HTTPVerb, node.fullPath)
	}
	node.methods[method.HTTPVerb] = &method
	return nil
}

// wildcardMatch validate if the checkPath meets the wildcardPath,
// for example /vought/12345 should match wildcard path /vought/:id;
// /vought/1234abcd/status should not match /vought/:id;
func wildcardMatch(wildcardPath string, checkPath string) bool {
	lowerWildcardPath := strings.ToLower(wildcardPath)
	lowerCheckPath := strings.ToLower(checkPath)
	wPathSplit := strings.Split(strings.TrimPrefix(lowerWildcardPath, constant.PathSlash), constant.PathSlash)
	cPathSplit := strings.Split(strings.TrimPrefix(lowerCheckPath, constant.PathSlash), constant.PathSlash)
	if len(wPathSplit) != len(cPathSplit) {
		return false
	}
	for i, s := range wPathSplit {
		if strings.Contains(s, constant.PathParamIdentifier) {
			cPathSplit[i] = s
		} else if wPathSplit[i] != cPathSplit[i] {
			return false
		}
	}
	return strings.Join(wPathSplit, constant.PathSlash) == strings.Join(cPathSplit, constant.PathSlash)
}

// NewRoute returns an empty router tree
func NewRoute() *Route {
	return &Route{
		tree:         avltree.NewWithStringComparator(),
		wildcardTree: avltree.NewWithStringComparator(),
	}
}
