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

package server

import (
	"net"
	"strings"
)

import (
	"github.com/apache/dubbo-go-pixiu/pilot/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pilot/pkg/serviceregistry/provider"
	"github.com/apache/dubbo-go-pixiu/pkg/config/constants"
	dnsProto "github.com/apache/dubbo-go-pixiu/pkg/dns/proto"
)

// Config for building the name table.
type Config struct {
	Node *model.Proxy
	Push *model.PushContext

	// MulticlusterHeadlessEnabled if true, the DNS name table for a headless service will resolve to
	// same-network endpoints in any cluster.
	MulticlusterHeadlessEnabled bool
}

// BuildNameTable produces a table of hostnames and their associated IPs that can then
// be used by the agent to resolve DNS. This logic is always active. However, local DNS resolution
// will only be effective if DNS capture is enabled in the proxy
func BuildNameTable(cfg Config) *dnsProto.NameTable {
	if cfg.Node.Type != model.SidecarProxy {
		// DNS resolution is only for sidecars
		return nil
	}

	out := &dnsProto.NameTable{
		Table: make(map[string]*dnsProto.NameTable_NameInfo),
	}
	for _, svc := range cfg.Node.SidecarScope.Services() {
		svcAddress := svc.GetAddressForProxy(cfg.Node)
		var addressList []string
		hostName := svc.Hostname
		if svcAddress != constants.UnspecifiedIP {
			// Filter out things we cannot parse as IP. Generally this means CIDRs, as anything else
			// should be caught in validation.
			if addr := net.ParseIP(svcAddress); addr == nil {
				continue
			}
			addressList = append(addressList, svcAddress)
		} else {
			// The IP will be unspecified here if its headless service or if the auto
			// IP allocation logic for service entry was unable to allocate an IP.
			if svc.Resolution == model.Passthrough && len(svc.Ports) > 0 {
				for _, instance := range cfg.Push.ServiceInstancesByPort(svc, svc.Ports[0].Port, nil) {
					// TODO(stevenctl): headless across-networks https://github.com/istio/istio/issues/38327
					sameNetwork := cfg.Node.InNetwork(instance.Endpoint.Network)
					sameCluster := cfg.Node.InCluster(instance.Endpoint.Locality.ClusterID)
					// For all k8s headless services, populate the dns table with the endpoint IPs as k8s does.
					// And for each individual pod, populate the dns table with the endpoint IP with a manufactured host name.
					if instance.Endpoint.SubDomain != "" && sameNetwork {
						// Follow k8s pods dns naming convention of "<hostname>.<subdomain>.<pod namespace>.svc.<cluster domain>"
						// i.e. "mysql-0.mysql.default.svc.cluster.local".
						parts := strings.SplitN(hostName.String(), ".", 2)
						if len(parts) != 2 {
							continue
						}
						address := []string{instance.Endpoint.Address}
						shortName := instance.Endpoint.HostName + "." + instance.Endpoint.SubDomain
						host := shortName + "." + parts[1] // Add cluster domain.
						nameInfo := &dnsProto.NameTable_NameInfo{
							Ips:       address,
							Registry:  string(svc.Attributes.ServiceRegistry),
							Namespace: svc.Attributes.Namespace,
							Shortname: shortName,
						}

						if _, f := out.Table[host]; !f || sameCluster {
							// We may have the same pod in two clusters (ie mysql-0 deployed in both places).
							// We can only return a single IP for these queries. We should prefer the local cluster,
							// so if the entry already exists only overwrite it if the instance is in our own cluster.
							out.Table[host] = nameInfo
						}
					}
					skipForMulticluster := !cfg.MulticlusterHeadlessEnabled && !sameCluster
					if skipForMulticluster || !sameNetwork {
						// We take only cluster-local endpoints. While this seems contradictory to
						// our logic other parts of the code, where cross-cluster is the default.
						// However, this only impacts the DNS response. If we were to send all
						// endpoints, cross network routing would break, as we do passthrough LB and
						// don't go through the network gateway. While we could, hypothetically, send
						// "network-local" endpoints, this would still make enabling DNS give vastly
						// different load balancing than without, so its probably best to filter.
						// This ends up matching the behavior of Kubernetes DNS.
						continue
					}
					// TODO: should we skip the node's own IP like we do in listener?
					addressList = append(addressList, instance.Endpoint.Address)
				}
			}
			if len(addressList) == 0 {
				// could not reliably determine the addresses of endpoints of headless service
				// or this is not a k8s service
				continue
			}
		}

		nameInfo := &dnsProto.NameTable_NameInfo{
			Ips:      addressList,
			Registry: string(svc.Attributes.ServiceRegistry),
		}
		if svc.Attributes.ServiceRegistry == provider.Kubernetes &&
			!strings.HasSuffix(hostName.String(), "."+constants.DefaultClusterSetLocalDomain) {
			// The agent will take care of resolving a, a.ns, a.ns.svc, etc.
			// No need to provide a DNS entry for each variant.
			//
			// NOTE: This is not done for Kubernetes Multi-Cluster Services (MCS) hosts, in order
			// to avoid conflicting with the entries for the regular (cluster.local) service.
			nameInfo.Namespace = svc.Attributes.Namespace
			nameInfo.Shortname = svc.Attributes.Name
		}
		out.Table[hostName.String()] = nameInfo
	}
	return out
}
