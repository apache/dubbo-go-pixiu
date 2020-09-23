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

// RouterNode defines the single method of the router configured API
type RouterNode struct {
	fullPath string
	wildcard bool
	methods  map[config.HTTPVerb]config.Method
	lock     sync.RWMutex
}

// RouterTree defines the tree of router APIs
type RouterTree struct {
	tree         *avltree.Tree
	wildcardTree *avltree.Tree
}

// Put put a key val into the tree
func (rt *RouterTree) Put(fullPath string, method config.Method) error {
	fullPath = strings.ToLower(fullPath)
	wildcard := containParam(fullPath)

	if wildcardNode, found := rt.searchWildcard(fullPath); found {
		return putMethod(wildcardNode, method)
	}

	node, ok := rt.tree.Get(fullPath)
	if !ok {
		ms := make(map[config.HTTPVerb]config.Method)
		rn := &RouterNode{
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
	if _, ok := node.(*RouterNode).methods[method.HTTPVerb]; ok {
		return errors.New(fmt.Sprintf("Method %s already exists in path %s", method.HTTPVerb, fullPath))
	}
	node.(*RouterNode).methods[method.HTTPVerb] = method
	return nil
}

func (rt *RouterTree) searchWildcard(fullPath string) (*RouterNode, bool) {
	wildcardPaths := rt.wildcardTree.Keys()
	for _, p := range wildcardPaths {
		if wildcardMatch(p.(string), fullPath) {
			n, ok := rt.wildcardTree.Get(p)
			return n.(*RouterNode), ok
		}
	}
	return nil, false
}

func putMethod(node *RouterNode, method config.Method) error {
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
