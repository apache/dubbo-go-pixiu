//  Copyright Istio Authors
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package cluster

import (
	"bytes"
	"fmt"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/kube"
)

// Map can be given as a shared reference to multiple Topology/Cluster implementations.
// allowing clusters to find each other for lookups of Primary, ConfigCluster, etc.
type Map = map[string]Cluster

func NewTopology(config Config, allClusters Map) Topology {
	return Topology{
		ClusterName:        config.Name,
		ClusterKind:        config.Kind,
		Network:            config.Network,
		ClusterHTTPProxy:   config.HTTPProxy,
		PrimaryClusterName: config.PrimaryClusterName,
		ConfigClusterName:  config.ConfigClusterName,
		AllClusters:        allClusters,
		Index:              len(allClusters),
	}
}

// Topology gives information about the relationship between clusters.
// Cluster implementations can embed this struct to include common functionality.
type Topology struct {
	ClusterName        string
	ClusterKind        Kind
	Network            string
	ClusterHTTPProxy   string
	PrimaryClusterName string
	ConfigClusterName  string
	Index              int
	// AllClusters should contain all AllClusters in the context
	AllClusters Map
}

// NetworkName the cluster is on
func (c Topology) NetworkName() string {
	return c.Network
}

// Name provides the ClusterName this cluster used by Istio.
func (c Topology) Name() string {
	return c.ClusterName
}

// HTTPProxy to connect to the cluster
func (c Topology) HTTPProxy() string {
	return c.ClusterHTTPProxy
}

// knownClusterNames maintains a well-known set of cluster names. These will always be used with
// StableNames if encountered.
var knownClusterNames = map[string]struct{}{
	"primary":               {},
	"remote":                {},
	"cross-network-primary": {},
}

// StableName provides a name used for testcase names. Deterministic, so testgrid
// can be consistent when the underlying cluster names are dynamic.
func (c Topology) StableName() string {
	var prefix string
	switch c.Kind() {
	case Kubernetes:
		// If its a known cluster name, use that directly.
		// This will not be dynamic, and allows 1:1 correlation of cluster name and test name for simplicity.
		if _, f := knownClusterNames[c.Name()]; f {
			return c.Name()
		}
		if c.IsPrimary() {
			if c.IsConfig() {
				prefix = "primary"
			} else {
				prefix = "externalistiod"
			}
		} else if c.IsRemote() {
			if c.IsConfig() {
				prefix = "config"
			} else {
				prefix = "remote"
			}
		}
	default:
		prefix = string(c.Kind())
	}

	return fmt.Sprintf("%s-%d", prefix, c.Index)
}

func (c Topology) Kind() Kind {
	return c.ClusterKind
}

func (c Topology) IsPrimary() bool {
	return c.Primary().Name() == c.Name()
}

func (c Topology) IsConfig() bool {
	return c.Config().Name() == c.Name()
}

func (c Topology) IsRemote() bool {
	return !c.IsPrimary()
}

func (c Topology) IsExternalControlPlane() bool {
	return c.IsPrimary() && !c.IsConfig()
}

func (c Topology) Primary() Cluster {
	cluster, ok := c.AllClusters[c.PrimaryClusterName]
	if !ok || cluster == nil {
		panic(fmt.Errorf("cannot find %s, the primary cluster for %s", c.PrimaryClusterName, c.Name()))
	}
	return cluster
}

func (c Topology) PrimaryName() string {
	return c.PrimaryClusterName
}

func (c Topology) Config() Cluster {
	cluster, ok := c.AllClusters[c.ConfigClusterName]
	if !ok || cluster == nil {
		panic(fmt.Errorf("cannot find %s, the config cluster for %s", c.ConfigClusterName, c.Name()))
	}
	return cluster
}

func (c Topology) ConfigName() string {
	return c.ConfigClusterName
}

func (c Topology) WithPrimary(primaryClusterName string) Topology {
	// TODO remove this, should only be provided by external config
	c.PrimaryClusterName = primaryClusterName
	return c
}

func (c Topology) WithConfig(configClusterName string) Topology {
	// TODO remove this, should only be provided by external config
	c.ConfigClusterName = configClusterName
	return c
}

func (c Topology) MinKubeVersion(minor uint) bool {
	cluster := c.AllClusters[c.ClusterName]
	if cluster.Kind() != Kubernetes && cluster.Kind() != Fake {
		return c.Primary().MinKubeVersion(minor)
	}
	return kube.IsAtLeastVersion(cluster, minor)
}

func (c Topology) MaxKubeVersion(minor uint) bool {
	cluster := c.AllClusters[c.ClusterName]
	if cluster.Kind() != Kubernetes && cluster.Kind() != Fake {
		return c.Primary().MaxKubeVersion(minor)
	}
	return kube.IsLessThanVersion(cluster, minor+1)
}

func (c Topology) String() string {
	buf := &bytes.Buffer{}

	_, _ = fmt.Fprintf(buf, "Name:               %s\n", c.Name())
	_, _ = fmt.Fprintf(buf, "StableName:         %s\n", c.StableName())
	_, _ = fmt.Fprintf(buf, "Kind:               %s\n", c.Kind())
	_, _ = fmt.Fprintf(buf, "PrimaryCluster:     %s\n", c.Primary().Name())
	_, _ = fmt.Fprintf(buf, "ConfigCluster:      %s\n", c.Config().Name())
	_, _ = fmt.Fprintf(buf, "Network:            %s\n", c.NetworkName())
	_, _ = fmt.Fprintf(buf, "HTTPProxy:            %s\n", c.HTTPProxy())

	return buf.String()
}
