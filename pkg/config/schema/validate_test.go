/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to you under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package schema

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfigSetYAMLRoundTripAndNormalize(t *testing.T) {
	data, err := os.ReadFile("testdata/config_set.yaml")
	require.NoError(t, err)

	configSet, err := DecodeConfigSetYAML(data)
	require.NoError(t, err)
	require.Len(t, configSet.Objects, 4)

	registry, err := NewBuiltinRegistry()
	require.NoError(t, err)
	normalized, err := registry.NormalizeConfigSet(configSet)
	require.NoError(t, err)

	method := findObject(t, normalized.Objects, KindMethod)
	assert.Equal(t, true, method.Spec["enable"])
	assert.Equal(t, "1s", method.Spec["timeout"])
	backend := method.Spec["backend"].(map[string]any)
	assert.Equal(t, "dubbo", backend["protocol"])

	// Normalize works on a defensive copy.
	originalMethod := findObject(t, configSet.Objects, KindMethod)
	assert.NotContains(t, originalMethod.Spec, "enable")
	assert.NotContains(t, originalMethod.Spec, "timeout")

	encoded, err := EncodeConfigSetYAML(normalized)
	require.NoError(t, err)
	roundTripped, err := DecodeConfigSetYAML(encoded)
	require.NoError(t, err)
	assert.Equal(t, normalized, roundTripped)
}

func TestValidationReportsFieldPaths(t *testing.T) {
	registry, err := NewBuiltinRegistry()
	require.NoError(t, err)

	object := validMethodObject()
	object.Spec["httpVerb"] = "GET"
	object.Spec["http"] = map[string]any{"verb": "INVALID"}
	err = registry.Validate(object)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "spec.http.verb")
	assert.Contains(t, err.Error(), "spec.httpVerb")
}

func TestHTTPBackendRequiresTarget(t *testing.T) {
	registry, err := NewBuiltinRegistry()
	require.NoError(t, err)

	object := validMethodObject()
	object.Spec["backend"] = map[string]any{"kind": "http"}
	err = registry.Validate(object)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "HTTP backend requires clusterRef, url, or host")
}

func TestMethodSignatureIndexesMustBeUnique(t *testing.T) {
	registry, err := NewBuiltinRegistry()
	require.NoError(t, err)

	object := validMethodObject()
	object.Spec["backend"].(map[string]any)["signature"] = []any{
		map[string]any{"index": 0, "javaType": "java.lang.String"},
		map[string]any{"index": 0, "javaType": "java.lang.Integer"},
	}
	err = registry.Validate(object)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "duplicates signature[0] index 0")
}

func TestConfigSetRejectsDuplicateObjectIdentity(t *testing.T) {
	registry, err := NewBuiltinRegistry()
	require.NoError(t, err)

	cluster := validClusterObject()
	configSet := ConfigSet{
		APIVersion: APIVersionV1Alpha1,
		Kind:       KindConfigSet,
		Metadata:   ObjectMetadata{ID: "default", Name: "default"},
		Objects:    []ConfigObject{cluster, cluster.Clone()},
	}
	_, err = registry.NormalizeConfigSet(configSet)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "duplicates objects[0] identity Cluster/cluster-user")
}

func TestNormalizeDefaultsLifecycleWithoutMutatingInput(t *testing.T) {
	registry, err := NewBuiltinRegistry()
	require.NoError(t, err)

	object := validClusterObject()
	normalized, err := registry.Normalize(object)
	require.NoError(t, err)
	assert.Empty(t, object.Metadata.Lifecycle)
	assert.Equal(t, LifecycleDraft, normalized.Metadata.Lifecycle)
}

func TestNormalizeRejectsUnknownLifecycle(t *testing.T) {
	registry, err := NewBuiltinRegistry()
	require.NoError(t, err)

	object := validClusterObject()
	object.Metadata.Lifecycle = "deleted"
	_, err = registry.Normalize(object)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "metadata.lifecycle")
}

func TestConfigSetValidatesCrossObjectReferences(t *testing.T) {
	registry, err := NewBuiltinRegistry()
	require.NoError(t, err)

	method := validMethodObject()
	method.Spec["backend"].(map[string]any)["clusterRef"] = "cluster-missing"
	configSet := ConfigSet{
		APIVersion: APIVersionV1Alpha1,
		Kind:       KindConfigSet,
		Metadata:   ObjectMetadata{ID: "default", Name: "default"},
		Objects:    []ConfigObject{method},
	}
	_, err = registry.NormalizeConfigSet(configSet)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "objects[0].spec.resourceRef")
	assert.Contains(t, err.Error(), "referenced Resource \"resource-users\" does not exist")
	assert.Contains(t, err.Error(), "objects[0].spec.backend.clusterRef")
	assert.Contains(t, err.Error(), "referenced Cluster \"cluster-missing\" does not exist")
}

func validMethodObject() ConfigObject {
	return ConfigObject{
		APIVersion: APIVersionV1Alpha1,
		Kind:       KindMethod,
		Metadata:   ObjectMetadata{ID: "method-get-user", Name: "get-user"},
		Spec: map[string]any{
			"resourceRef": "resource-users",
			"http":        map[string]any{"verb": "GET"},
			"backend": map[string]any{
				"kind":      "dubbo",
				"interface": "com.example.UserService",
				"method":    "getUser",
			},
		},
	}
}

func validClusterObject() ConfigObject {
	return ConfigObject{
		APIVersion: APIVersionV1Alpha1,
		Kind:       KindCluster,
		Metadata:   ObjectMetadata{ID: "cluster-user", Name: "user-service"},
		Spec:       map[string]any{},
	}
}

func findObject(t *testing.T, objects []ConfigObject, kind string) ConfigObject {
	t.Helper()
	for _, object := range objects {
		if object.Kind == kind {
			return object
		}
	}
	t.Fatalf("object kind %s not found", kind)
	return ConfigObject{}
}
