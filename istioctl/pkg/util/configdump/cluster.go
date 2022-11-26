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

package configdump

import (
	"sort"
)

import (
	adminapi "github.com/envoyproxy/go-control-plane/envoy/admin/v3"
	cluster "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
)

import (
	v3 "github.com/apache/dubbo-go-pixiu/pilot/pkg/xds/v3"
)

// GetDynamicClusterDump retrieves a cluster dump with just dynamic active clusters in it
func (w *Wrapper) GetDynamicClusterDump(stripVersions bool) (*adminapi.ClustersConfigDump, error) {
	clusterDump, err := w.GetClusterConfigDump()
	if err != nil {
		return nil, err
	}
	dac := clusterDump.GetDynamicActiveClusters()
	// Allow sorting to work even if we don't have the exact same type
	for i := range dac {
		dac[i].Cluster.TypeUrl = v3.ClusterType
	}
	sort.Slice(dac, func(i, j int) bool {
		cluster := &cluster.Cluster{}
		err = dac[i].Cluster.UnmarshalTo(cluster)
		if err != nil {
			return false
		}
		name := cluster.Name
		err = dac[j].Cluster.UnmarshalTo(cluster)
		if err != nil {
			return false
		}
		return name < cluster.Name
	})
	if stripVersions {
		for i := range dac {
			dac[i].VersionInfo = ""
			dac[i].LastUpdated = nil
		}
	}
	return &adminapi.ClustersConfigDump{DynamicActiveClusters: dac}, nil
}

// GetClusterConfigDump retrieves the cluster config dump from the ConfigDump
func (w *Wrapper) GetClusterConfigDump() (*adminapi.ClustersConfigDump, error) {
	clusterDumpAny, err := w.getSection(clusters)
	if err != nil {
		return nil, err
	}
	clusterDump := &adminapi.ClustersConfigDump{}
	err = clusterDumpAny.UnmarshalTo(clusterDump)
	if err != nil {
		return nil, err
	}
	return clusterDump, nil
}
