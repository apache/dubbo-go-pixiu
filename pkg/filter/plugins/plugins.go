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
	"plugin"
	"strings"
)

import (
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/api/config"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/context"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/filter"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

var (
	apiURLWithPluginsMap = make(map[string]FilterChain)
	groupWithPluginsMap  = make(map[string]map[string]WithFunc)
	errEmptyPluginConfig = errors.New("Empty plugin config")
)

// FilterChain include Pre & Post filters
type FilterChain struct {
	Pre    context.FilterChain
	Post context.FilterChain
}

// WithFunc is a single plugin details
type WithFunc struct {
	Name  string
	Priority int
	fn       context.FilterFunc
}

// InitPluginsGroup prase api_config.yaml(pluginsGroup) to map[string][]PluginsWithFunc
func InitPluginsGroup(groups []config.PluginsGroup, filePath string) {
	if "" == filePath || len(groups) == 0 {
		return
	}

	// load file.so
	pls, err := plugin.Open(filePath)
	if nil != err {
		panic(err)
	}

	for _, group := range groups {
		pwdMap := make(map[string]WithFunc, len(group.Plugins))

		// trans to context.FilterFunc
		for _, pl := range group.Plugins {
			pwf := WithFunc{pl.Name, pl.Priority, loadExternalPlugin(&pl, pls)}
			pwdMap[pl.Name] = pwf
		}

		groupWithPluginsMap[group.GroupName] = pwdMap
	}
}

// InitAPIURLWithFilterChain must behind InitPluginsGroup call
func InitAPIURLWithFilterChain(resources []config.Resource) {
	pairURLWithFilterChain("", resources, nil)
}

func pairURLWithFilterChain(parentPath string, resources []config.Resource, parentFilterChains *FilterChain) {
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

		currentFilterChains, err := getAPIFilterChains(&resource.Plugins)

		if err == nil {
			apiURLWithPluginsMap[fullPath] = *currentFilterChains
			parentFilterChains = currentFilterChains
		} else {
			if parentFilterChains != nil {
				apiURLWithPluginsMap[fullPath] = *parentFilterChains
			}
		}

		if len(resource.Resources) > 0 {
			pairURLWithFilterChain(resource.Path, resource.Resources, parentFilterChains)
		}
	}
}

// GetAPIFilterFuncsWithAPIURL is get filterchain with path
func GetAPIFilterFuncsWithAPIURL(url string) FilterChain {
	// found from cache
	if flc, found := apiURLWithPluginsMap[url]; found {
		logger.Debugf("GetExternalPlugins is:%v", flc)
		return flc
	}

	// or return empty FilterFunc
	return FilterChain{context.FilterChain{}, context.FilterChain{}}
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

	panic(errEmptyPluginConfig)
}

func getAPIFilterChains(pluginsConfig *config.PluginsConfig) (fcs *FilterChain, err error) {
	pre := getAPIFilterFuncsWithPluginsGroup(&pluginsConfig.PrePlugins)
	post := getAPIFilterFuncsWithPluginsGroup(&pluginsConfig.PostPlugins)

	if len(pre) == 0 && len(post) == 0 {
		return nil, errors.New("FilterChains is empty")
	}

	return &FilterChain{pre, post}, nil
}

func getAPIFilterFuncsWithPluginsGroup(plu *config.PluginsInUse) []context.FilterFunc {
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
