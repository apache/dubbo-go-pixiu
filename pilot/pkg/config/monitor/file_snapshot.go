// Copyright Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package monitor

import (
	"os"
	"path/filepath"
	"sort"
)

import (
	"github.com/apache/dubbo-go-pixiu/pilot/pkg/config/kube/crd"
	"github.com/apache/dubbo-go-pixiu/pkg/config"
	"github.com/apache/dubbo-go-pixiu/pkg/config/schema/collection"
	"github.com/apache/dubbo-go-pixiu/pkg/config/schema/collections"
)

var supportedExtensions = map[string]bool{
	".yaml": true,
	".yml":  true,
}

// FileSnapshot holds a reference to a file directory that contains crd
// config and filter criteria for which of those configs will be parsed.
type FileSnapshot struct {
	root             string
	domainSuffix     string
	configTypeFilter map[config.GroupVersionKind]bool
}

// NewFileSnapshot returns a snapshotter.
// If no types are provided in the descriptor, all Istio types will be allowed.
func NewFileSnapshot(root string, schemas collection.Schemas, domainSuffix string) *FileSnapshot {
	snapshot := &FileSnapshot{
		root:             root,
		domainSuffix:     domainSuffix,
		configTypeFilter: make(map[config.GroupVersionKind]bool),
	}

	ss := schemas.All()
	if len(ss) == 0 {
		ss = collections.Pilot.All()
	}

	for _, k := range ss {
		if _, ok := collections.Pilot.FindByGroupVersionKind(k.Resource().GroupVersionKind()); ok {
			snapshot.configTypeFilter[k.Resource().GroupVersionKind()] = true
		}
	}

	return snapshot
}

// ReadConfigFiles parses files in the root directory and returns a sorted slice of
// eligible model.Config. This can be used as a configFunc when creating a Monitor.
func (f *FileSnapshot) ReadConfigFiles() ([]*config.Config, error) {
	var result []*config.Config

	err := filepath.Walk(f.root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		} else if !supportedExtensions[filepath.Ext(path)] || (info.Mode()&os.ModeType) != 0 {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			log.Warnf("Failed to read %s: %v", path, err)
			return err
		}
		configs, err := parseInputs(data, f.domainSuffix)
		if err != nil {
			log.Warnf("Failed to parse %s: %v", path, err)
			return err
		}

		// Filter any unsupported types before appending to the result.
		for _, cfg := range configs {
			if !f.configTypeFilter[cfg.GroupVersionKind] {
				continue
			}
			result = append(result, cfg)
		}
		return nil
	})
	if err != nil {
		log.Warnf("failure during filepath.Walk: %v", err)
	}

	// Sort by the config IDs.
	sort.Sort(byKey(result))
	return result, err
}

// parseInputs is identical to crd.ParseInputs, except that it returns an array of config pointers.
func parseInputs(data []byte, domainSuffix string) ([]*config.Config, error) {
	configs, _, err := crd.ParseInputs(string(data))

	// Convert to an array of pointers.
	refs := make([]*config.Config, len(configs))
	for i := range configs {
		refs[i] = &configs[i]
		refs[i].Domain = domainSuffix
	}
	return refs, err
}

// byKey is an array of config objects that is capable or sorting by Namespace, GroupVersionKind, and Name.
type byKey []*config.Config

func (rs byKey) Len() int {
	return len(rs)
}

func (rs byKey) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

func (rs byKey) Less(i, j int) bool {
	return compareIds(rs[i], rs[j]) < 0
}
