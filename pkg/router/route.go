package router

import (
	"path"
	"strings"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/config"
	"github.com/emirpasic/gods/trees/avltree"
	"github.com/pkg/errors"
)

// RouterGroup a easy way to manage the actual router tree, provides the apis to group the routers
type RouterGroup struct {
	root       bool
	basePath   string
	routerTree *avltree.Tree
}

// Group deviates a new router group from current group. use the same routerTree.
func (rg *RouterGroup) Group(relativePath string) (*RouterGroup, error) {
	if len(relativePath) == 0 {
		return nil, errors.New("Cannot group router with empty path")
	}
	if relativePath[0] != '/' {
		return nil, errors.New("Path must start with '/'")
	}
	return &RouterGroup{
		basePath:   rg.absolutePath(relativePath),
		routerTree: rg.routerTree,
	}, nil
}

func (rg *RouterGroup) absolutePath(relativePath string) string {
	if len(relativePath) == 0 {
		return rg.basePath
	}
	return strings.TrimRight(path.Join(rg.basePath, relativePath), "/")
}

// Add adds the new router node to the group
func (rg *RouterGroup) Add(path string, method config.Method) error {
	rg.routerTree.Put(path, method)
	return nil
}

// NewRouter returns a nil tree root router group
func NewRouter() *RouterGroup {
	return &RouterGroup{
		root:       true,
		basePath:   "/",
		routerTree: avltree.NewWithStringComparator(),
	}
}
