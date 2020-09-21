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
	node, ok := rt.tree.Get(fullPath)

	if !ok {
		ms := make(map[config.HTTPVerb]config.Method)
		wildcard := isWildcard(fullPath)
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
	resources := strings.Split(fullPath[1:], "/")
	resources[0] = "/" + resources[0]
	wildcardPaths := rt.wildcardTree.Keys()

	for _, p := range wildcardPaths {
		if wildcardMatch(p.(string), fullPath) {
			n, ok := rt.wildcardTree.Get(p)
			return n.(*RouterNode), ok
		}
	}
	// parmPositions := []int{}
	// for i := range resources {
	// 	if isWildcard(resources[i]) {
	// 		parmPositions = append(parmPositions, i)
	// 	}
	// }
	// for _, position := range parmPositions {
	// 	searchPath := strings.Join(resources[:position], "/")
	// 	if node, ok := rt.tree.Get(searchPath); ok {
	// 		return node.(*RouterNode), ok
	// 	}
	// }
	return nil, false
}

func isWildcard(fullPath string) bool {
	for _, s := range fullPath {
		if s == ':' {
			return true
		}
	}
	return false
}

func wildcardMatch(wildcardPath string, targetPath string) bool {
	return false
}
