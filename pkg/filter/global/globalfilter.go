package global

import (
	"fmt"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

var globalFilterManagers = make(map[string]*filter.FilterManager)

func RegisterGlobalFilterManager(filterManagerName string, filterManager *filter.FilterManager) {
	if _, ok := globalFilterManagers[filterManagerName]; ok {
		logger.Warn("FilterManager already exists %s", filterManagerName)
		return
	}
	globalFilterManagers[filterManagerName] = filterManager
}

func GetGlobalFilterManager(name string) *filter.FilterManager {
	filterManager, ok := globalFilterManagers[name]
	if !ok {
		panic(fmt.Sprintf("globalFilterManagers %s does not exist", name))
	}
	return filterManager
}
