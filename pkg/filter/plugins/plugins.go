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

package plugins

import (
	"errors"
	"github.com/dubbogo/dubbo-go-proxy/pkg/filter"
	"plugin"
	"strings"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/constant"
	"github.com/dubbogo/dubbo-go-proxy/pkg/config"
	"github.com/dubbogo/dubbo-go-proxy/pkg/context"
	"github.com/dubbogo/dubbo-go-proxy/pkg/logger"
)

var (
	apiUrlWithPluginsMap = make(map[string]context.FilterChain)
	groupWithPluginsMap  = make(map[string]map[string]PluginsWithFunc)
	emptyPluginConfigErr = errors.New("Empty plugin config")
)

// single plugin details
type PluginsWithFunc struct {
	Name     string
	Priority int
	fn       context.FilterFunc
}

// Prase api_config.yaml(pluginsGroup) to map[string][]PluginsWithFunc
func InitPluginsGroup(groups []config.PluginsGroup, filePath string) {
	// load file.so
	pls, err := plugin.Open(filePath)
	if nil != err {
		panic(err)
	}

	for _, group := range groups {
		pwdMap := make(map[string]PluginsWithFunc, len(group.Plugins))

		// trans to context.FilterFunc
		for _, pl := range group.Plugins {
			pwf := PluginsWithFunc{pl.Name, pl.Priority, loadExternalPlugin(&pl, pls)}
			pwdMap[pl.Name] = pwf
		}

		groupWithPluginsMap[group.GroupName] = pwdMap
	}
}

// must behind InitPluginsGroup call
func InitApiUrlWithFilterChain(resources []config.Resource) {
	pairUrlWithFilterChain("", resources, context.FilterChain{})
}

func pairUrlWithFilterChain(parentPath string, resources []config.Resource, parentFilterFuncs []context.FilterFunc) {

	if len(resources) == 0 {
		return
	}

	groupPath := parentPath
	if parentPath == constant.PathSlash {
		groupPath = ""
	}

	for _, resource := range resources {
		fullPath := groupPath + resource.Path
		if !strings.HasPrefix(resource.Path, constant.PathSlash) {
			continue
		}

		currentFuncArr := getApiFilterFuncsWithPluginsGroup(&resource.Plugins)

		if len(currentFuncArr) > 0 {
			apiUrlWithPluginsMap[fullPath] = currentFuncArr
			parentFilterFuncs = currentFuncArr
		} else {
			if len(parentFilterFuncs) > 0 {
				apiUrlWithPluginsMap[fullPath] = parentFilterFuncs
			}
		}

		if len(resource.Resources) > 0 {
			pairUrlWithFilterChain(resource.Path, resource.Resources, parentFilterFuncs)
		}
	}

}

func GetApiFilterFuncsWithApiUrl(url string) context.FilterChain {
	// found from cache
	if funcs, found := apiUrlWithPluginsMap[url]; found {
		logger.Debugf("GetExternalPlugins is:%v,len:%d", funcs, len(funcs))
		return funcs
	}

	// or return empty FilterFunc
	return context.FilterChain{}
}

func loadExternalPlugin(p *config.Plugin, pl *plugin.Plugin) context.FilterFunc {
	if nil != p {
		logger.Infof("loadExternalPlugin name is :%s,version:%s,Priority:%d", p.Name, p.Version, p.Priority)
		sb, err := pl.Lookup(p.ExternalLookupName)
		if nil != err {
			panic(err)
		}

		sbf := sb.(func() filter.Filter)
		logger.Infof("loadExternalPlugin %s success", p.Name)
		return sbf().Do()
	}

	panic(emptyPluginConfigErr)
}

func getApiFilterFuncsWithPluginsGroup(plu *config.PluginsInUse) []context.FilterFunc {
	// not set plugins
	if nil == plu || nil != plu && len(plu.GroupNames) == 0 && len(plu.PluginNames) == 0 {
		return []context.FilterFunc{}
	}

	tmpMap := make(map[string]context.FilterFunc)

	// found with group name
	for _, groupName := range plu.GroupNames {
		pwfMap, found := groupWithPluginsMap[groupName]
		if found {
			for _, pwf := range pwfMap {
				tmpMap[pwf.Name] = pwf.fn
			}
		}
	}

	// found with with name from all group
	for _, group := range groupWithPluginsMap {
		for _, name := range plu.PluginNames {
			if pwf, found := group[name]; found {
				tmpMap[pwf.Name] = pwf.fn
			}
		}
	}

	ret := make([]context.FilterFunc, 0)
	for _, v := range tmpMap {
		if nil != v {
			ret = append(ret, v)
		}
	}

	return ret
}
