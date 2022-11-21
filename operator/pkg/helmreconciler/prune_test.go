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

package helmreconciler

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/yaml"
)

import (
	"github.com/apache/dubbo-go-pixiu/operator/pkg/apis/istio/v1alpha1"
	"github.com/apache/dubbo-go-pixiu/operator/pkg/name"
	"github.com/apache/dubbo-go-pixiu/operator/pkg/object"
	"github.com/apache/dubbo-go-pixiu/operator/pkg/util/clog"
	"github.com/apache/dubbo-go-pixiu/operator/pkg/util/progress"
	"github.com/apache/dubbo-go-pixiu/pkg/kube"
	"github.com/apache/dubbo-go-pixiu/pkg/test/env"
)

const (
	testRevision = "test"
)

func TestHelmReconciler_DeleteControlPlaneByManifest(t *testing.T) {
	t.Run("deleteControlPlaneByManifest", func(t *testing.T) {
		cl := fake.NewClientBuilder().Build()
		df := filepath.Join(env.IstioSrc, "manifests/profiles/default.yaml")
		iopStr, err := os.ReadFile(df)
		if err != nil {
			t.Fatal(err)
		}
		iop := &v1alpha1.IstioOperator{}
		if err := yaml.UnmarshalStrict(iopStr, iop); err != nil {
			t.Fatal(err)
		}
		iop.Spec.Revision = testRevision
		iop.Spec.InstallPackagePath = filepath.Join(env.IstioSrc, "manifests")

		h := &HelmReconciler{
			client:     cl,
			kubeClient: kube.NewFakeClientWithVersion("24"),
			opts: &Options{
				ProgressLog: progress.NewLog(),
				Log:         clog.NewDefaultLogger(),
			},
			iop:           iop,
			countLock:     &sync.Mutex{},
			prunedKindSet: map[schema.GroupKind]struct{}{},
		}
		manifestMap, err := h.RenderCharts()
		if err != nil {
			t.Fatalf("failed to render manifest: %v", err)
		}
		applyResourcesIntoCluster(t, h, manifestMap)
		if err := h.DeleteControlPlaneByManifests(manifestMap, testRevision, false); err != nil {
			t.Fatalf("HelmReconciler.DeleteControlPlaneByManifests() error = %v", err)
		}
		for _, gvk := range append(h.NamespacedResources(), ClusterCPResources...) {
			receiver := &unstructured.Unstructured{}
			receiver.SetGroupVersionKind(schema.GroupVersionKind{Group: gvk.Group, Version: gvk.Version, Kind: gvk.Kind})
			objKey := client.ObjectKey{Namespace: "istio-system", Name: "istiod-test"}
			if gvk.Kind == name.MutatingWebhookConfigurationStr {
				objKey.Name = "istio-sidecar-injector-test"
			}
			// List does not work well here as that requires adding all resource types to the fake client scheme
			if err := h.client.Get(context.TODO(), objKey, receiver); err != nil {
				// the error is expected because we expect resources do not exist any more in the cluster
				t.Logf(err.Error())
			}
			obj := receiver.Object
			if obj["spec"] != nil {
				t.Errorf("got resource: %s/%s from the cluster, expected to be deleted", receiver.GetKind(), receiver.GetName())
			}
		}
	})
}

func applyResourcesIntoCluster(t *testing.T, h *HelmReconciler, manifestMap name.ManifestMap) {
	for cn, ms := range manifestMap.Consolidated() {
		objects, err := object.ParseK8sObjectsFromYAMLManifest(ms)
		if err != nil {
			t.Fatalf("failed parse k8s objects from yaml: %v", err)
		}
		for _, obj := range objects {
			obju := obj.UnstructuredObject()
			if err := h.applyLabelsAndAnnotations(obju, cn); err != nil {
				t.Errorf("failed to apply label and annotations: %v", err)
			}
			if err := h.ApplyObject(obj.UnstructuredObject(), false); err != nil {
				t.Errorf("HelmReconciler.ApplyObject() error = %v", err)
			}
		}
	}
}
