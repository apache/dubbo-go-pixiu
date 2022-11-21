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

package virtualservice

import (
	"fmt"
	"strings"
)

import (
	"istio.io/api/networking/v1alpha3"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/config/analysis"
	"github.com/apache/dubbo-go-pixiu/pkg/config/analysis/analyzers/util"
	"github.com/apache/dubbo-go-pixiu/pkg/config/analysis/msg"
	"github.com/apache/dubbo-go-pixiu/pkg/config/host"
	"github.com/apache/dubbo-go-pixiu/pkg/config/resource"
	"github.com/apache/dubbo-go-pixiu/pkg/config/schema/collection"
	"github.com/apache/dubbo-go-pixiu/pkg/config/schema/collections"
)

// GatewayAnalyzer checks the gateways associated with each virtual service
type GatewayAnalyzer struct{}

var _ analysis.Analyzer = &GatewayAnalyzer{}

// Metadata implements Analyzer
func (s *GatewayAnalyzer) Metadata() analysis.Metadata {
	return analysis.Metadata{
		Name:        "virtualservice.GatewayAnalyzer",
		Description: "Checks the gateways associated with each virtual service",
		Inputs: collection.Names{
			collections.IstioNetworkingV1Alpha3Gateways.Name(),
			collections.IstioNetworkingV1Alpha3Virtualservices.Name(),
		},
	}
}

// Analyze implements Analyzer
func (s *GatewayAnalyzer) Analyze(c analysis.Context) {
	c.ForEach(collections.IstioNetworkingV1Alpha3Virtualservices.Name(), func(r *resource.Instance) bool {
		s.analyzeVirtualService(r, c)
		return true
	})
}

func (s *GatewayAnalyzer) analyzeVirtualService(r *resource.Instance, c analysis.Context) {
	vs := r.Message.(*v1alpha3.VirtualService)
	vsNs := r.Metadata.FullName.Namespace
	vsName := r.Metadata.FullName

	for i, gwName := range vs.Gateways {
		// This is a special-case accepted value
		if gwName == util.MeshGateway {
			continue
		}

		gwFullName := resource.NewShortOrFullName(vsNs, gwName)

		if !c.Exists(collections.IstioNetworkingV1Alpha3Gateways.Name(), gwFullName) {
			m := msg.NewReferencedResourceNotFound(r, "gateway", gwName)

			if line, ok := util.ErrorLine(r, fmt.Sprintf(util.VSGateway, i)); ok {
				m.Line = line
			}

			c.Report(collections.IstioNetworkingV1Alpha3Virtualservices.Name(), m)
		}

		if !vsHostInGateway(c, gwFullName, vs.Hosts, vsNs.String()) {
			m := msg.NewVirtualServiceHostNotFoundInGateway(r, vs.Hosts, vsName.String(), gwFullName.String())

			if line, ok := util.ErrorLine(r, fmt.Sprintf(util.VSGateway, i)); ok {
				m.Line = line
			}

			c.Report(collections.IstioNetworkingV1Alpha3Virtualservices.Name(), m)
		}
	}
}

func vsHostInGateway(c analysis.Context, gateway resource.FullName, vsHosts []string, vsNamespace string) bool {
	var gatewayHosts []string
	var gatewayNs string

	c.ForEach(collections.IstioNetworkingV1Alpha3Gateways.Name(), func(r *resource.Instance) bool {
		if r.Metadata.FullName == gateway {
			s := r.Message.(*v1alpha3.Gateway)
			gatewayNs = r.Metadata.FullName.Namespace.String()
			for _, v := range s.Servers {
				sanitizeServerHostNamespace(v, gatewayNs)
				gatewayHosts = append(gatewayHosts, v.Hosts...)
			}
		}

		return true
	})

	gatewayHostNames := host.NamesForNamespace(gatewayHosts, vsNamespace)
	for _, gh := range gatewayHostNames {
		for _, vsh := range vsHosts {
			gatewayHost := gh
			vsHost := host.Name(vsh)

			if gatewayHost.Matches(vsHost) {
				return true
			}
		}
	}

	return false
}

// convert ./host to currentNamespace/Host
// */host to just host
// */* to just *
func sanitizeServerHostNamespace(server *v1alpha3.Server, namespace string) {
	for i, h := range server.Hosts {
		if strings.Contains(h, "/") {
			parts := strings.Split(h, "/")
			if parts[0] == "." {
				server.Hosts[i] = fmt.Sprintf("%s/%s", namespace, parts[1])
			} else if parts[0] == "*" {
				if parts[1] == "*" {
					server.Hosts = []string{"*"}
					return
				}
				server.Hosts[i] = parts[1]
			}
		}
	}
}
