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

package kube

import (
	"bytes"
	"fmt"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/kube"
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework/components/cluster"
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework/components/echo"
)

var _ echo.Cluster = &Cluster{}

// Cluster for a Kubernetes cluster. Provides access via a kube.Client.
type Cluster struct {
	// filename is the path to the kubeconfig file for this cluster.
	filename string

	// vmSupport indicates the cluster is being used for fake VMs
	vmSupport bool

	// ExtendedClient is embedded to interact with the kube cluster.
	kube.ExtendedClient

	// Topology is embedded to include common functionality.
	cluster.Topology
}

// CanDeploy for a kube cluster returns true if the config is a non-vm, or if the cluster supports
// fake pod-based VMs.
func (c *Cluster) CanDeploy(config echo.Config) (echo.Config, bool) {
	if config.DeployAsVM && !c.vmSupport {
		return echo.Config{}, false
	}
	return config, true
}

// OverrideTopology allows customizing the relationship between this and other clusters
// for a single suite. This practice is discouraged, and separate test jobs should be created
// on a per-topology bassis.
// TODO remove this when centralistiod test is isolated as it's own job
func (c *Cluster) OverrideTopology(fn func(cluster.Topology) cluster.Topology) {
	c.Topology = fn(c.Topology)
}

func (c *Cluster) String() string {
	buf := &bytes.Buffer{}

	_, _ = fmt.Fprint(buf, c.Topology.String())
	_, _ = fmt.Fprintf(buf, "Filename:           %s\n", c.filename)

	return buf.String()
}

// Filename of the kubeconfig file for this cluster.
// TODO(nmittler): Remove the need for this by changing operator to use provided kube clients directly.
func (c Cluster) Filename() string {
	return c.filename
}
