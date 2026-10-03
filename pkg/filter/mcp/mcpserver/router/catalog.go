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
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

type StringSet map[string]struct{}

func NewStringSet(values []string) StringSet {
	if len(values) == 0 {
		return nil
	}
	set := make(StringSet, len(values))
	for _, value := range values {
		set[value] = struct{}{}
	}
	return set
}

func NewValidatedStringSet(owner string, values []string) (StringSet, error) {
	set := make(StringSet, len(values))
	for i, value := range values {
		name := strings.TrimSpace(value)
		if name == "" {
			return nil, fmt.Errorf("%s at index %d is empty", owner, i)
		}
		if _, exists := set[name]; exists {
			return nil, fmt.Errorf("%s duplicate %q at index %d", owner, name, i)
		}
		set[name] = struct{}{}
	}
	return set, nil
}

func (s StringSet) Contains(value string) bool {
	if len(s) == 0 {
		return false
	}
	_, ok := s[value]
	return ok
}

func (s StringSet) Intersects(other StringSet) bool {
	if len(s) == 0 || len(other) == 0 {
		return false
	}
	if len(s) > len(other) {
		s, other = other, s
	}
	for value := range s {
		if other.Contains(value) {
			return true
		}
	}
	return false
}

func (s StringSet) Names() []string {
	if len(s) == 0 {
		return nil
	}
	names := make([]string, 0, len(s))
	for name := range s {
		names = append(names, name)
	}
	return names
}

func ValidateUniqueNames(tools []model.ToolConfig) error {
	seen := make(StringSet, len(tools))
	for i, tool := range tools {
		if tool.Name == "" {
			return fmt.Errorf("tool name is required at index %d", i)
		}
		if seen.Contains(tool.Name) {
			return fmt.Errorf("duplicate tool name %q at index %d", tool.Name, i)
		}
		seen[tool.Name] = struct{}{}
	}
	return nil
}

func ValidateNonEmptyUniqueStrings(owner string, values []string) error {
	_, err := NewValidatedStringSet(owner, values)
	return err
}

type ToolCatalogView struct {
	ordered []model.ToolConfig
	byName  map[string]model.ToolConfig
}

func NewToolCatalogView(tools []model.ToolConfig) ToolCatalogView {
	view := ToolCatalogView{
		ordered: tools,
		byName:  make(map[string]model.ToolConfig, len(tools)),
	}
	for _, tool := range tools {
		if _, exists := view.byName[tool.Name]; !exists {
			view.byName[tool.Name] = tool
		}
	}
	return view
}

func (v ToolCatalogView) PickOrdered(names []string) []model.ToolConfig {
	if len(names) == 0 {
		return nil
	}
	out := make([]model.ToolConfig, 0, len(names))
	for _, name := range names {
		if tool, ok := v.byName[name]; ok {
			out = append(out, tool)
		}
	}
	return out
}

func (v ToolCatalogView) FilterSet(names StringSet) []model.ToolConfig {
	if len(v.ordered) == 0 || len(names) == 0 {
		return nil
	}
	out := make([]model.ToolConfig, 0, len(v.ordered))
	for _, tool := range v.ordered {
		if names.Contains(tool.Name) {
			out = append(out, tool)
		}
	}
	return out
}

func (v ToolCatalogView) Names() []string {
	names := make([]string, len(v.ordered))
	for i, tool := range v.ordered {
		names[i] = tool.Name
	}
	return names
}

func (v ToolCatalogView) DiscoverableNames() []string {
	names := make([]string, 0, len(v.ordered))
	for _, tool := range v.ordered {
		if IsToolDiscoverable(tool) {
			names = append(names, tool.Name)
		}
	}
	return names
}

func (v ToolCatalogView) ByName(name string) (model.ToolConfig, bool) {
	tool, ok := v.byName[name]
	return tool, ok
}

func IsToolDiscoverable(t model.ToolConfig) bool {
	return t.Meta == nil || t.Meta.DiscoveryVisibility == nil || *t.Meta.DiscoveryVisibility
}

func DiscoverableTools(tools []model.ToolConfig) []model.ToolConfig {
	if len(tools) == 0 {
		return nil
	}
	out := make([]model.ToolConfig, 0, len(tools))
	for _, tool := range tools {
		if IsToolDiscoverable(tool) {
			out = append(out, tool)
		}
	}
	return out
}

func DiscoverableToolNames(tools []model.ToolConfig) []string {
	return NewToolCatalogView(tools).DiscoverableNames()
}
