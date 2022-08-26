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

package analyzers

import (
	"github.com/apache/dubbo-go-pixiu/pkg/config/analysis"
	"github.com/apache/dubbo-go-pixiu/pkg/config/analysis/analyzers/annotations"
	"github.com/apache/dubbo-go-pixiu/pkg/config/analysis/analyzers/authz"
	"github.com/apache/dubbo-go-pixiu/pkg/config/analysis/analyzers/deployment"
	"github.com/apache/dubbo-go-pixiu/pkg/config/analysis/analyzers/deprecation"
	"github.com/apache/dubbo-go-pixiu/pkg/config/analysis/analyzers/destinationrule"
	"github.com/apache/dubbo-go-pixiu/pkg/config/analysis/analyzers/envoyfilter"
	"github.com/apache/dubbo-go-pixiu/pkg/config/analysis/analyzers/gateway"
	"github.com/apache/dubbo-go-pixiu/pkg/config/analysis/analyzers/injection"
	"github.com/apache/dubbo-go-pixiu/pkg/config/analysis/analyzers/multicluster"
	"github.com/apache/dubbo-go-pixiu/pkg/config/analysis/analyzers/schema"
	"github.com/apache/dubbo-go-pixiu/pkg/config/analysis/analyzers/service"
	"github.com/apache/dubbo-go-pixiu/pkg/config/analysis/analyzers/serviceentry"
	"github.com/apache/dubbo-go-pixiu/pkg/config/analysis/analyzers/sidecar"
	"github.com/apache/dubbo-go-pixiu/pkg/config/analysis/analyzers/virtualservice"
	"github.com/apache/dubbo-go-pixiu/pkg/config/analysis/analyzers/webhook"
)

// All returns all analyzers
func All() []analysis.Analyzer {
	analyzers := []analysis.Analyzer{
		// Please keep this list sorted alphabetically by pkg.name for convenience
		&annotations.K8sAnalyzer{},
		&authz.AuthorizationPoliciesAnalyzer{},
		&deployment.ServiceAssociationAnalyzer{},
		&deployment.ApplicationUIDAnalyzer{},
		&deprecation.FieldAnalyzer{},
		&gateway.IngressGatewayPortAnalyzer{},
		&gateway.CertificateAnalyzer{},
		&gateway.SecretAnalyzer{},
		&gateway.ConflictingGatewayAnalyzer{},
		&injection.Analyzer{},
		&injection.ImageAnalyzer{},
		&injection.ImageAutoAnalyzer{},
		&multicluster.MeshNetworksAnalyzer{},
		&service.PortNameAnalyzer{},
		&sidecar.DefaultSelectorAnalyzer{},
		&sidecar.SelectorAnalyzer{},
		&virtualservice.ConflictingMeshGatewayHostsAnalyzer{},
		&virtualservice.DestinationHostAnalyzer{},
		&virtualservice.DestinationRuleAnalyzer{},
		&virtualservice.GatewayAnalyzer{},
		&virtualservice.JWTClaimRouteAnalyzer{},
		&virtualservice.RegexAnalyzer{},
		&destinationrule.CaCertificateAnalyzer{},
		&serviceentry.ProtocolAdressesAnalyzer{},
		&webhook.Analyzer{},
		&envoyfilter.EnvoyPatchAnalyzer{},
	}

	analyzers = append(analyzers, schema.AllValidationAnalyzers()...)

	return analyzers
}

// AllCombined returns all analyzers combined as one
func AllCombined() *analysis.CombinedAnalyzer {
	return analysis.Combine("all", All()...)
}
