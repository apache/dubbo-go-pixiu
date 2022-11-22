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

package destinationrule

import (
	"fmt"
)

import (
	"istio.io/api/networking/v1alpha3"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/config/analysis"
	"github.com/apache/dubbo-go-pixiu/pkg/config/analysis/analyzers/util"
	"github.com/apache/dubbo-go-pixiu/pkg/config/analysis/msg"
	"github.com/apache/dubbo-go-pixiu/pkg/config/resource"
	"github.com/apache/dubbo-go-pixiu/pkg/config/schema/collection"
	"github.com/apache/dubbo-go-pixiu/pkg/config/schema/collections"
)

// CaCertificateAnalyzer checks if CaCertificate is set in case mode is SIMPLE/MUTUAL
type CaCertificateAnalyzer struct{}

var _ analysis.Analyzer = &CaCertificateAnalyzer{}

func (c *CaCertificateAnalyzer) Metadata() analysis.Metadata {
	return analysis.Metadata{
		Name:        "destinationrule.CaCertificateAnalyzer",
		Description: "Checks if caCertificates is set when TLS mode is SIMPLE/MUTUAL",
		Inputs: collection.Names{
			collections.IstioNetworkingV1Alpha3Destinationrules.Name(),
		},
	}
}

func (c *CaCertificateAnalyzer) Analyze(ctx analysis.Context) {
	ctx.ForEach(collections.IstioNetworkingV1Alpha3Destinationrules.Name(), func(r *resource.Instance) bool {
		c.analyzeDestinationRule(r, ctx)
		return true
	})
}

func (c *CaCertificateAnalyzer) analyzeDestinationRule(r *resource.Instance, ctx analysis.Context) {
	dr := r.Message.(*v1alpha3.DestinationRule)
	drNs := r.Metadata.FullName.Namespace
	drName := r.Metadata.FullName.String()
	mode := dr.GetTrafficPolicy().GetTls().GetMode()

	if mode == v1alpha3.ClientTLSSettings_SIMPLE || mode == v1alpha3.ClientTLSSettings_MUTUAL {
		if dr.GetTrafficPolicy().GetTls().GetCaCertificates() == "" {
			m := msg.NewNoServerCertificateVerificationDestinationLevel(r, drName,
				drNs.String(), mode.String(), dr.GetHost())

			if line, ok := util.ErrorLine(r, fmt.Sprintf(util.DestinationRuleTLSCert)); ok {
				m.Line = line
			}
			ctx.Report(collections.IstioNetworkingV1Alpha3Destinationrules.Name(), m)
		}
	}
	portSettings := dr.TrafficPolicy.GetPortLevelSettings()

	for i, p := range portSettings {
		mode = p.GetTls().GetMode()
		if mode == v1alpha3.ClientTLSSettings_SIMPLE || mode == v1alpha3.ClientTLSSettings_MUTUAL {
			if p.GetTls().GetCaCertificates() == "" {
				m := msg.NewNoServerCertificateVerificationPortLevel(r, drName,
					drNs.String(), mode.String(), dr.GetHost(), p.GetPort().String())

				if line, ok := util.ErrorLine(r, fmt.Sprintf(util.DestinationRuleTLSPortLevelCert, i)); ok {
					m.Line = line
				}
				ctx.Report(collections.IstioNetworkingV1Alpha3Destinationrules.Name(), m)
			}
		}
	}
}
