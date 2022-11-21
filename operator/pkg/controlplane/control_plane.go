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

package controlplane

import (
	"fmt"
	"sort"
)

import (
	"istio.io/api/operator/v1alpha1"
	"k8s.io/apimachinery/pkg/version"
)

import (
	iop "github.com/apache/dubbo-go-pixiu/operator/pkg/apis/istio/v1alpha1"
	"github.com/apache/dubbo-go-pixiu/operator/pkg/component"
	"github.com/apache/dubbo-go-pixiu/operator/pkg/name"
	"github.com/apache/dubbo-go-pixiu/operator/pkg/translate"
	"github.com/apache/dubbo-go-pixiu/operator/pkg/util"
	"github.com/apache/dubbo-go-pixiu/pkg/util/sets"
)

// IstioControlPlane is an installation of an Istio control plane.
type IstioControlPlane struct {
	// components is a slice of components that are part of the feature.
	components []component.IstioComponent
	started    bool
}

// NewIstioControlPlane creates a new IstioControlPlane and returns a pointer to it.
func NewIstioControlPlane(
	installSpec *v1alpha1.IstioOperatorSpec,
	translator *translate.Translator,
	filter []string,
	ver *version.Info,
) (*IstioControlPlane, error) {
	out := &IstioControlPlane{}
	opts := &component.Options{
		InstallSpec: installSpec,
		Translator:  translator,
		Filter:      sets.New(filter...),
		Version:     ver,
	}
	for _, c := range name.AllCoreComponentNames {
		o := *opts
		ns, err := name.Namespace(c, installSpec)
		if err != nil {
			return nil, err
		}
		o.Namespace = ns
		out.components = append(out.components, component.NewCoreComponent(c, &o))
	}

	if installSpec.Components != nil {
		for idx, c := range installSpec.Components.IngressGateways {
			o := *opts
			o.Namespace = defaultIfEmpty(c.Namespace, iop.Namespace(installSpec))
			out.components = append(out.components, component.NewIngressComponent(c.Name, idx, c, &o))
		}
		for idx, c := range installSpec.Components.EgressGateways {
			o := *opts
			o.Namespace = defaultIfEmpty(c.Namespace, iop.Namespace(installSpec))
			out.components = append(out.components, component.NewEgressComponent(c.Name, idx, c, &o))
		}
	}
	return out, nil
}

func orderedKeys(m map[string]*v1alpha1.ExternalComponentSpec) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func defaultIfEmpty(val, dflt string) string {
	if val == "" {
		return dflt
	}
	return val
}

// Run starts the Istio control plane.
func (i *IstioControlPlane) Run() error {
	for _, c := range i.components {
		if err := c.Run(); err != nil {
			return err
		}
	}
	i.started = true
	return nil
}

// RenderManifest returns a manifest rendered against
func (i *IstioControlPlane) RenderManifest() (manifests name.ManifestMap, errsOut util.Errors) {
	if !i.started {
		return nil, util.NewErrs(fmt.Errorf("istioControlPlane must be Run before calling RenderManifest"))
	}

	manifests = make(name.ManifestMap)
	for _, c := range i.components {
		ms, err := c.RenderManifest()
		errsOut = util.AppendErr(errsOut, err)
		manifests[c.ComponentName()] = append(manifests[c.ComponentName()], ms)
	}
	if len(errsOut) > 0 {
		return nil, errsOut
	}
	return
}

// componentsEqual reports whether the given components are equal to those in i.
func (i *IstioControlPlane) componentsEqual(components []component.IstioComponent) bool {
	if i.components == nil && components == nil {
		return true
	}
	if len(i.components) != len(components) {
		return false
	}
	for c := 0; c < len(i.components); c++ {
		if i.components[c].ComponentName() != components[c].ComponentName() {
			return false
		}
		if i.components[c].Namespace() != components[c].Namespace() {
			return false
		}
		if i.components[c].Enabled() != components[c].Enabled() {
			return false
		}
		if i.components[c].ResourceName() != components[c].ResourceName() {
			return false
		}
	}
	return true
}
