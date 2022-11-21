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

package util

import (
	"os"
	"path/filepath"
	"testing"
)

import (
	meshconfig "istio.io/api/mesh/v1alpha1"
	v1alpha12 "istio.io/api/operator/v1alpha1"
	"sigs.k8s.io/yaml"
)

import (
	"github.com/apache/dubbo-go-pixiu/operator/pkg/apis/istio/v1alpha1"
	"github.com/apache/dubbo-go-pixiu/pkg/config/mesh"
	"github.com/apache/dubbo-go-pixiu/pkg/test/env"
	"github.com/apache/dubbo-go-pixiu/pkg/util/protomarshal"
)

func TestOverlayIOP(t *testing.T) {
	cases := []struct {
		path string
	}{
		{
			filepath.Join(env.IstioSrc, "manifests/profiles/default.yaml"),
		},
		{
			filepath.Join(env.IstioSrc, "manifests/profiles/demo.yaml"),
		},
		{
			filepath.Join("testdata", "overlay-iop.yaml"),
		},
	}

	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {
			b, err := os.ReadFile(tc.path)
			if err != nil {
				t.Fatal(err)
			}
			// overlaying tree over itself exercises all paths for merging
			if _, err := OverlayIOP(string(b), string(b)); err != nil {
				t.Fatal(err)
			}
		})
	}
}

// TestOverlayIOPExhaustiveness exhaustiveness check of `OverlayIOP`
// Once some one add a new `Provider` in api, we should update `wellknownProviders` and
// add to `meshConfigExtensionProvider`
func TestOverlayIOPExhaustiveness(t *testing.T) {
	wellknownProviders := map[string]struct{}{
		"prometheus":            {},
		"envoy_file_access_log": {},
		"stackdriver":           {},
		"envoy_otel_als":        {},
		"envoy_ext_authz_http":  {},
		"envoy_ext_authz_grpc":  {},
		"zipkin":                {},
		"lightstep":             {},
		"datadog":               {},
		"opencensus":            {},
		"skywalking":            {},
		"envoy_http_als":        {},
		"envoy_tcp_als":         {},
	}

	unexpectedProviders := make([]string, 0)

	msg := &meshconfig.MeshConfig_ExtensionProvider{}
	pb := msg.ProtoReflect()
	md := pb.Descriptor()

	of := md.Oneofs().Get(0)
	for i := 0; i < of.Fields().Len(); i++ {
		o := of.Fields().Get(i)
		n := string(o.Name())
		if _, ok := wellknownProviders[n]; ok {
			delete(wellknownProviders, n)
		} else {
			unexpectedProviders = append(unexpectedProviders, n)
		}
	}

	if len(wellknownProviders) != 0 || len(unexpectedProviders) != 0 {
		t.Errorf("unexpected provider not implemented in OverlayIOP, wellknownProviders: %v unexpectedProviders: %v", wellknownProviders, unexpectedProviders)
		t.Fail()
	}
}

func TestOverlayIOPDefaultMeshConfig(t *testing.T) {
	// Transform default mesh config into map[string]interface{} for inclusion in IstioOperator.
	m := mesh.DefaultMeshConfig()
	my, err := protomarshal.ToJSONMap(m)
	if err != nil {
		t.Fatal(err)
	}

	iop := &v1alpha1.IstioOperator{
		Spec: &v1alpha12.IstioOperatorSpec{
			MeshConfig: MustStruct(my),
		},
	}

	iy, err := yaml.Marshal(iop)
	if err != nil {
		t.Fatal(err)
	}

	// overlaying tree over itself exercises all paths for merging
	if _, err := OverlayIOP(string(iy), string(iy)); err != nil {
		t.Fatal(err)
	}
}

func TestOverlayIOPIngressGatewayLabel(t *testing.T) {
	l1, err := os.ReadFile("testdata/yaml/input/yaml_layer1.yaml")
	if err != nil {
		t.Fatal(err)
	}
	l2, err := os.ReadFile("testdata/yaml/input/yaml_layer2.yaml")
	if err != nil {
		t.Fatal(err)
	}

	if _, err := OverlayIOP(string(l1), string(l2)); err != nil {
		t.Fatal(err)
	}
}
