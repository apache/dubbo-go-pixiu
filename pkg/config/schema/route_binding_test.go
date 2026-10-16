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
	"time"
)

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdminRouteBindingYAMLRoundTrip(t *testing.T) {
	data, err := os.ReadFile("testdata/admin_route_binding.yaml")
	require.NoError(t, err)
	object, err := DecodeAdminObjectYAML(data)
	require.NoError(t, err)
	assert.Equal(t, KindAdminRouteBinding, object.Kind)
	assert.Equal(t, "user-get", object.Metadata.Name)

	registry, err := NewBuiltinRegistry()
	require.NoError(t, err)
	normalized, err := registry.Normalize(object)
	require.NoError(t, err)

	encoded, err := EncodeAdminObjectYAML(normalized)
	require.NoError(t, err)
	roundTripped, err := DecodeAdminObjectYAML(encoded)
	require.NoError(t, err)
	assert.Equal(t, normalized, roundTripped)
}

func TestNormalizeAppliesFormDefaultsWithoutMutatingInput(t *testing.T) {
	registry, err := NewBuiltinRegistry()
	require.NoError(t, err)
	object := validRouteBindingObject()

	normalized, err := registry.Normalize(object)
	require.NoError(t, err)
	assert.NotContains(t, object.Spec, "publish")
	assert.NotContains(t, object.Spec, "params")

	entry := normalized.Spec["entry"].(map[string]any)
	target := normalized.Spec["target"].(map[string]any)
	assert.Equal(t, "http", entry["protocol"])
	assert.Equal(t, "dubbo", target["protocol"])
	assert.Equal(t, []any{}, normalized.Spec["params"])
	assert.Equal(t, true, normalized.Spec["enabled"])
	assert.NotContains(t, normalized.Spec, "publish")
}

func TestNormalizeDropsDeprecatedPublishMetadata(t *testing.T) {
	registry, err := NewBuiltinRegistry()
	require.NoError(t, err)
	object := validRouteBindingObject()
	object.Spec["publish"] = map[string]any{"mode": "draft", "validate": false}

	normalized, err := registry.Normalize(object)
	require.NoError(t, err)
	assert.NotContains(t, normalized.Spec, "publish")
	assert.Contains(t, object.Spec, "publish")
}

func TestCompileAdminRouteBindingToLegacyResourceAndMethod(t *testing.T) {
	object := loadExampleRouteBinding(t)
	registry, err := NewBuiltinRegistry()
	require.NoError(t, err)

	compiled, err := CompileAdminRouteBinding(registry, object)
	require.NoError(t, err)
	assert.Equal(t, "/api/v1/users/:id", compiled.Resource.Path)
	assert.Equal(t, "restful", compiled.Resource.Type)
	assert.Equal(t, time.Second, compiled.Resource.Timeout)
	assert.Empty(t, compiled.Resource.Methods)

	method := compiled.Method
	assert.Equal(t, "/api/v1/users/:id", method.ResourcePath)
	assert.Equal(t, "GET", method.HTTPVerb)
	assert.True(t, method.Enable)
	assert.Equal(t, time.Second, method.Timeout)
	assert.Equal(t, "http", method.InboundRequest.RequestType)
	assert.Equal(t, "dubbo", method.IntegrationRequest.RequestType)
	assert.Equal(t, "UserProvider", method.ApplicationName)
	assert.Equal(t, "com.example.UserService", method.Interface)
	assert.Equal(t, "GetUser", method.Method)
	assert.Equal(t, "1.0.0", method.Version)
	assert.Equal(t, "stable", method.Group)
	assert.Equal(t, "user-dubbo", method.ClusterName)
	assert.Equal(t, []string{"java.lang.String"}, method.ParameterTypes)
	require.Len(t, method.MappingParams, 1)
	assert.Equal(t, "uri.id", method.MappingParams[0].Name)
	assert.Equal(t, "0", method.MappingParams[0].MapTo)
	assert.Equal(t, "java.lang.String", method.MappingParams[0].MapType)

	legacy := compiled.LegacyAPIConfig()
	require.Len(t, legacy.Resources, 1)
	require.Len(t, legacy.Resources[0].Methods, 1)
}

func TestCompileAdminRouteBindingPreservesDisabledState(t *testing.T) {
	object := validRouteBindingObject()
	object.Spec["enabled"] = false
	registry, err := NewBuiltinRegistry()
	require.NoError(t, err)

	compiled, err := CompileAdminRouteBinding(registry, object)
	require.NoError(t, err)
	assert.False(t, compiled.Method.Enable)
}

func TestPreviewYAMLContainsOnlyLegacyAPIConfig(t *testing.T) {
	registry, err := NewBuiltinRegistry()
	require.NoError(t, err)
	compiled, err := CompileAdminRouteBinding(registry, loadExampleRouteBinding(t))
	require.NoError(t, err)

	preview, err := compiled.PreviewYAML()
	require.NoError(t, err)
	previewString := string(preview)
	assert.Contains(t, previewString, "resources:")
	assert.Contains(t, previewString, "integrationRequest:")
	assert.Contains(t, previewString, "mapTo: \"0\"")
	assert.Contains(t, previewString, "parameterTypes:")
	assert.NotContains(t, previewString, "AdminRouteBinding")
	assert.NotContains(t, previewString, "publish:")
}

func TestAdminRouteBindingValidationReportsActionablePaths(t *testing.T) {
	registry, err := NewBuiltinRegistry()
	require.NoError(t, err)
	object := validRouteBindingObject()
	object.Spec["params"] = []any{
		map[string]any{"from": "query.id", "to": 1, "type": "String"},
		map[string]any{"from": "uri.id", "to": 1, "type": "java.lang.String"},
	}

	err = registry.Validate(object)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "spec.params[0].from")
	assert.Contains(t, err.Error(), "spec.params[0].type")
	assert.Contains(t, err.Error(), "spec.params[1].to")
	assert.Contains(t, err.Error(), "argument indexes must be contiguous from 0")
}

func validRouteBindingObject() AdminObject {
	return AdminObject{
		Kind:     KindAdminRouteBinding,
		Metadata: ObjectMetadata{Name: "user-get"},
		Spec: map[string]any{
			"entry": map[string]any{
				"path":   "/api/v1/users/:id",
				"method": "GET",
			},
			"target": map[string]any{
				"application": "UserProvider",
				"interface":   "com.example.UserService",
				"method":      "GetUser",
				"cluster":     "user-dubbo",
			},
		},
	}
}

func loadExampleRouteBinding(t *testing.T) AdminObject {
	t.Helper()
	data, err := os.ReadFile("testdata/admin_route_binding.yaml")
	require.NoError(t, err)
	object, err := DecodeAdminObjectYAML(data)
	require.NoError(t, err)
	return object
}
