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

package authn

import (
	meshconfig "istio.io/api/mesh/v1alpha1"
)

import (
	"github.com/apache/dubbo-go-pixiu/pilot/pkg/features"
	"github.com/apache/dubbo-go-pixiu/pkg/util/sets"
)

func TrustDomainsForValidation(meshConfig *meshconfig.MeshConfig) []string {
	if features.SkipValidateTrustDomain {
		return nil
	}

	tds := append([]string{meshConfig.TrustDomain}, meshConfig.TrustDomainAliases...)
	return dedupTrustDomains(tds)
}

func dedupTrustDomains(tds []string) []string {
	known := sets.New()
	deduped := make([]string, 0, len(tds))

	for _, td := range tds {
		if td != "" && !known.Contains(td) {
			known.Insert(td)
			deduped = append(deduped, td)
		}
	}
	return deduped
}
