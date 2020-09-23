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
	"fmt"
	"strings"
	"sync"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/config"
	"github.com/emirpasic/gods/trees/avltree"
	"github.com/pkg/errors"
)

// Node defines the single method of the router configured API
type Node struct {
	fullPath string
	wildcard bool
	methods  map[config.HTTPVerb]config.Method
	lock     sync.RWMutex
}

// Tree defines the tree of router APIs
type Tree struct {
	tree         *avltree.Tree
	wildcardTree *avltree.Tree
}

// Put put a key val into the tree
func (rt *Tree) Put(fullPath string, method config.Method) error {
	fullPath = strings.ToLower(fullPath)
	wildcard := containParam(fullPath)

	if wildcardNode, found := rt.searchWildcard(fullPath); found {
		return putMethod(wildcardNode, method)
	}

	node, ok := rt.tree.Get(fullPath)
	if !ok {
		ms := make(map[config.HTTPVerb]config.Method)
		rn := &Node{
			fullPath: fullPath,
			methods:  ms,
			wildcard: wildcard,
		}
		rn.methods[method.HTTPVerb] = method
		if wildcard {
			rt.wildcardTree.Put(fullPath, rn)
		}
		rt.tree.Put(fullPath, rn)
		return nil
	}
	if _, ok := node.(*Node).methods[method.HTTPVerb]; ok {
		return errors.New(fmt.Sprintf("Method %s already exists in path %s", method.HTTPVerb, fullPath))
	}
	node.(*Node).methods[method.HTTPVerb] = method
	return nil
}

func (rt *Tree) searchWildcard(fullPath string) (*Node, bool) {
	wildcardPaths := rt.wildcardTree.Keys()
	for _, p := range wildcardPaths {
		if wildcardMatch(p.(string), fullPath) {
			n, ok := rt.wildcardTree.Get(p)
			return n.(*Node), ok
		}
	}
	return nil, false
}

func putMethod(node *Node, method config.Method) error {
	if _, ok := node.methods[method.HTTPVerb]; ok {
		return errors.New(fmt.Sprintf("Method %s already exists in path %s", method.HTTPVerb, node.fullPath))
	}
	node.methods[method.HTTPVerb] = method
	return nil
}

func containParam(fullPath string) bool {
	for _, s := range fullPath {
		if s == ':' {
			return true
		}
	}
	return false
}

func wildcardMatch(wildcardPath string, checkPath string) bool {
	wildcardPath = strings.ToLower(wildcardPath)
	checkPath = strings.ToLower(checkPath)
	wPathSplit := strings.Split(wildcardPath[1:], "/")
	cPathSplit := strings.Split(checkPath[1:], "/")
	if len(wPathSplit) != len(cPathSplit) {
		return false
	}
	for i, s := range wPathSplit {
		if containParam(s) {
			cPathSplit[i] = s
		}
	}
	return strings.Join(wPathSplit, "/") == strings.Join(cPathSplit, "/")
}
