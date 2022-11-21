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

package deprecation

import (
	"fmt"
)

import (
	"istio.io/api/networking/v1alpha3"
	k8sext "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/config/analysis"
	"github.com/apache/dubbo-go-pixiu/pkg/config/analysis/msg"
	"github.com/apache/dubbo-go-pixiu/pkg/config/resource"
	"github.com/apache/dubbo-go-pixiu/pkg/config/schema/collection"
	"github.com/apache/dubbo-go-pixiu/pkg/config/schema/collections"
)

// FieldAnalyzer checks for deprecated Istio types and fields
type FieldAnalyzer struct{}

// Tracks Istio CRDs removed from manifests/charts/base/crds/crd-all.gen.yaml
var deprecatedCRDs = []k8sext.CustomResourceDefinitionSpec{
	{
		Group: "rbac.istio.io",
		Names: k8sext.CustomResourceDefinitionNames{Kind: "ClusterRbacConfig"},
	},
	{
		Group: "rbac.istio.io",
		Names: k8sext.CustomResourceDefinitionNames{Kind: "RbacConfig"},
	},
	{
		Group: "rbac.istio.io",
		Names: k8sext.CustomResourceDefinitionNames{Kind: "ServiceRole"},
	},
	{
		Group: "rbac.istio.io",
		Names: k8sext.CustomResourceDefinitionNames{Kind: "ServiceRoleBinding"},
	},
}

// Currently we don't have an Istio API that tells which Istio API fields are deprecated.
// Run `find . -name "*.proto" -exec grep -i "deprecated=true" \{\} \; -print`
// to see what is deprecated.  This analyzer is hand-crafted.

// Metadata implements analyzer.Analyzer
func (*FieldAnalyzer) Metadata() analysis.Metadata {
	deprecationInputs := collection.Names{
		collections.IstioNetworkingV1Alpha3Virtualservices.Name(),
		collections.IstioNetworkingV1Alpha3Sidecars.Name(),
		collections.K8SApiextensionsK8SIoV1Customresourcedefinitions.Name(),
	}

	return analysis.Metadata{
		Name:        "deprecation.DeprecationAnalyzer",
		Description: "Checks for deprecated Istio types and fields",
		Inputs: append(deprecationInputs,
			collections.Deprecated.CollectionNames()...),
	}
}

// Analyze implements analysis.Analyzer
func (fa *FieldAnalyzer) Analyze(ctx analysis.Context) {
	ctx.ForEach(collections.IstioNetworkingV1Alpha3Virtualservices.Name(), func(r *resource.Instance) bool {
		fa.analyzeVirtualService(r, ctx)
		return true
	})
	ctx.ForEach(collections.IstioNetworkingV1Alpha3Sidecars.Name(), func(r *resource.Instance) bool {
		fa.analyzeSidecar(r, ctx)
		return true
	})
	ctx.ForEach(collections.K8SApiextensionsK8SIoV1Customresourcedefinitions.Name(), func(r *resource.Instance) bool {
		fa.analyzeCRD(r, ctx)
		return true
	})
	for _, name := range collections.Deprecated.CollectionNames() {
		ctx.ForEach(name, func(r *resource.Instance) bool {
			ctx.Report(name,
				msg.NewDeprecated(r, crDeprecatedMessage(name.String())))
			return true
		})
	}
}

func (*FieldAnalyzer) analyzeCRD(r *resource.Instance, ctx analysis.Context) {
	for _, depCRD := range deprecatedCRDs {
		var group, kind string
		switch crd := r.Message.(type) {
		case *k8sext.CustomResourceDefinition:
			group = crd.Spec.Group
			kind = crd.Spec.Names.Kind
		case *k8sext.CustomResourceDefinitionSpec:
			group = crd.Group
			kind = crd.Names.Kind
		}
		if group == depCRD.Group && kind == depCRD.Names.Kind {
			ctx.Report(collections.K8SApiextensionsK8SIoV1Customresourcedefinitions.Name(),
				msg.NewDeprecated(r, crRemovedMessage(depCRD.Group, depCRD.Names.Kind)))
		}
	}
}

func (*FieldAnalyzer) analyzeSidecar(r *resource.Instance, ctx analysis.Context) {
	sc := r.Message.(*v1alpha3.Sidecar)

	if sc.OutboundTrafficPolicy != nil {
		if sc.OutboundTrafficPolicy.EgressProxy != nil {
			ctx.Report(collections.IstioNetworkingV1Alpha3Virtualservices.Name(),
				msg.NewDeprecated(r, ignoredMessage("OutboundTrafficPolicy.EgressProxy")))
		}
	}
}

func (*FieldAnalyzer) analyzeVirtualService(r *resource.Instance, ctx analysis.Context) {
	vs := r.Message.(*v1alpha3.VirtualService)

	for _, httpRoute := range vs.Http {
		if httpRoute.Fault != nil {
			if httpRoute.Fault.Delay != nil {
				// nolint: staticcheck
				if httpRoute.Fault.Delay.Percent > 0 {
					ctx.Report(collections.IstioNetworkingV1Alpha3Virtualservices.Name(),
						msg.NewDeprecated(r, replacedMessage("HTTPRoute.fault.delay.percent", "HTTPRoute.fault.delay.percentage")))
				}
			}
		}
	}
}

func replacedMessage(deprecated, replacement string) string {
	return fmt.Sprintf("%s is deprecated; use %s", deprecated, replacement)
}

func ignoredMessage(field string) string {
	return fmt.Sprintf("%s ignored", field)
}

func crDeprecatedMessage(typename string) string {
	return fmt.Sprintf("Custom resource type %q is deprecated", typename)
}

func crRemovedMessage(group, kind string) string {
	return fmt.Sprintf("Custom resource type %s %s is removed", group, kind)
}
