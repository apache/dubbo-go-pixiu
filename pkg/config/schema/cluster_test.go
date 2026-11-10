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
	"errors"
	"testing"
)

import (
	"github.com/stretchr/testify/require"

	"gopkg.in/yaml.v3"
)

import (
	legacyconfig "github.com/apache/dubbo-go-pixiu/pkg/config"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func validAdminCluster() AdminObject {
	return AdminObject{
		Kind:     KindAdminCluster,
		Metadata: ObjectMetadata{Name: "user-provider"},
		Spec: map[string]any{
			"endpoint": map[string]any{"address": "127.0.0.1", "port": 20880},
		},
	}
}

func TestAdminClusterSchemaAndCompile(t *testing.T) {
	registry, err := NewBuiltinRegistry()
	require.NoError(t, err)
	definition, ok := registry.Lookup(KindAdminCluster)
	require.True(t, ok)
	require.Equal(t, "Static", definition.Fields["type"].Default)
	require.Equal(t, []any{"Static", "StrictDNS", "LogicalDns", "EDS", "OriginalDst"}, definition.Fields["type"].Enum)
	runtimeTypes := make([]any, 0, len(model.DiscoveryTypeValue))
	for discoveryType := range model.DiscoveryTypeValue {
		runtimeTypes = append(runtimeTypes, discoveryType)
	}
	require.ElementsMatch(t, runtimeTypes, definition.Fields["type"].Enum)
	require.Equal(t, float64(1), *definition.Fields["endpoint"].Properties["port"].Minimum)
	require.Equal(t, float64(65535), *definition.Fields["endpoint"].Properties["port"].Maximum)
	*definition.Fields["endpoint"].Properties["port"].Maximum = 10
	unchanged, ok := registry.Lookup(KindAdminCluster)
	require.True(t, ok)
	require.Equal(t, float64(65535), *unchanged.Fields["endpoint"].Properties["port"].Maximum)

	input := validAdminCluster()
	compiled, err := CompileAdminCluster(registry, input)
	require.NoError(t, err)
	require.Equal(t, legacyconfig.Cluster{Name: "user-provider", Type: "Static", Address: "127.0.0.1", Port: 20880}, compiled.Cluster)
	require.Equal(t, "Static", compiled.Source.Spec["type"])
	require.NotContains(t, input.Spec, "type", "defaults must not mutate caller input")
	require.Equal(t, 42, compiled.RuntimeCluster(42).ID)
	require.Zero(t, compiled.Cluster.ID)

	preview, err := compiled.PreviewYAML()
	require.NoError(t, err)
	require.NotContains(t, string(preview), "id:")
	var decoded legacyconfig.Cluster
	require.NoError(t, yaml.Unmarshal(preview, &decoded))
	require.Equal(t, compiled.Cluster, decoded)
}

func TestAdminClusterValidationPaths(t *testing.T) {
	registry, err := NewBuiltinRegistry()
	require.NoError(t, err)
	cases := []struct {
		name string
		edit func(*AdminObject)
		path string
		code string
	}{
		{"empty name", func(o *AdminObject) { o.Metadata.Name = "  " }, "metadata.name", "required"},
		{"unsupported type", func(o *AdminObject) { o.Spec["type"] = "Other" }, "spec.type", "enum"},
		{"empty address", func(o *AdminObject) { o.Spec["endpoint"].(map[string]any)["address"] = " \t " }, "spec.endpoint.address", "required"},
		{"port 0", func(o *AdminObject) { o.Spec["endpoint"].(map[string]any)["port"] = 0 }, "spec.endpoint.port", "minimum"},
		{"port 65536", func(o *AdminObject) { o.Spec["endpoint"].(map[string]any)["port"] = 65536 }, "spec.endpoint.port", "maximum"},
		{"unknown spec field", func(o *AdminObject) { o.Spec["id"] = 7 }, "spec.id", "unknown"},
		{"unknown endpoint field", func(o *AdminObject) { o.Spec["endpoint"].(map[string]any)["weight"] = 10 }, "spec.endpoint.weight", "unknown"},
		{"missing endpoint", func(o *AdminObject) { delete(o.Spec, "endpoint") }, "spec.endpoint", "required"},
		{"missing address", func(o *AdminObject) { delete(o.Spec["endpoint"].(map[string]any), "address") }, "spec.endpoint.address", "required"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			object := validAdminCluster()
			tc.edit(&object)
			_, err := CompileAdminCluster(registry, object)
			var issues ValidationErrors
			require.True(t, errors.As(err, &issues), "expected field-level validation error: %v", err)
			require.NotEmpty(t, findIssueMessage(issues, tc.path, tc.code), "expected a nonempty message at %s (%s): %+v", tc.path, tc.code, issues)
		})
	}
	for _, port := range []int{1, 65535} {
		object := validAdminCluster()
		object.Spec["endpoint"].(map[string]any)["port"] = port
		_, err := CompileAdminCluster(registry, object)
		require.NoError(t, err)
	}
	object := validAdminCluster()
	object.Spec["endpoint"].(map[string]any)["port"] = float64(20880) // JSON requests decode into float64.
	compiled, err := CompileAdminCluster(registry, object)
	require.NoError(t, err)
	require.Equal(t, 20880, compiled.Cluster.Port)
	object.Spec["endpoint"].(map[string]any)["port"] = 20880.5
	_, err = CompileAdminCluster(registry, object)
	var issues ValidationErrors
	require.ErrorAs(t, err, &issues)
	require.NotEmpty(t, findIssueMessage(issues, "spec.endpoint.port", "type"))
	for discoveryType := range model.DiscoveryTypeValue {
		object := validAdminCluster()
		object.Spec["type"] = discoveryType
		compiled, err := CompileAdminCluster(registry, object)
		require.NoError(t, err)
		require.Equal(t, discoveryType, compiled.Cluster.Type)
	}
}

func findIssueMessage(issues ValidationErrors, path, code string) string {
	for _, issue := range issues {
		if issue.Path == path && issue.Code == code {
			return issue.Message
		}
	}
	return ""
}

func TestAdminClusterLegacyYAMLCompatibility(t *testing.T) {
	registry, err := NewBuiltinRegistry()
	require.NoError(t, err)
	input := []byte("name: user-provider\ntype: StrictDNS\naddress: service.example.org\nport: 20880\nid: 17\n")
	object, id, err := DecodeLegacyClusterYAML(input)
	require.NoError(t, err)
	require.Equal(t, 17, id)
	require.NotContains(t, object.Spec, "id")
	compiled, err := CompileAdminCluster(registry, object)
	require.NoError(t, err)
	var original legacyconfig.Cluster
	require.NoError(t, yaml.Unmarshal(input, &original))
	require.Equal(t, original, compiled.RuntimeCluster(id))

	preview, err := compiled.PreviewYAML()
	require.NoError(t, err)
	previewObject, previewID, err := DecodeLegacyClusterYAML(preview)
	require.NoError(t, err)
	require.Equal(t, object, previewObject)
	require.Zero(t, previewID)

	_, _, err = DecodeLegacyClusterYAML(append(input, []byte("unexpected: value\n")...))
	require.ErrorContains(t, err, "unexpected")
	_, _, err = DecodeLegacyClusterYAML([]byte("name: one\n---\nname: two\n"))
	require.ErrorContains(t, err, "one document")
	_, _, err = DecodeLegacyClusterYAML([]byte("name: one\n---\nname: ["))
	require.ErrorContains(t, err, "decode legacy Cluster YAML")
}

func TestAdminClusterYAMLUnknownFieldIsReported(t *testing.T) {
	object, err := DecodeAdminObjectYAML([]byte("kind: AdminCluster\nmetadata:\n  name: backend\nspec:\n  type: Static\n  endpoint:\n    address: 127.0.0.1\n    port: 20880\n  id: 9\n"))
	require.NoError(t, err)
	registry, err := NewBuiltinRegistry()
	require.NoError(t, err)
	_, err = CompileAdminCluster(registry, object)
	var issues ValidationErrors
	require.ErrorAs(t, err, &issues)
	require.NotEmpty(t, findIssueMessage(issues, "spec.id", "unknown"))
}

func TestAdminClusterSchemaRejectsInvertedPortRange(t *testing.T) {
	definition := adminClusterSchema()
	*definition.Fields["endpoint"].Properties["port"].Minimum = 65536
	err := NewRegistry().Register(definition)
	require.ErrorContains(t, err, "minimum exceeds maximum")
}

func TestCompileAdminClusterRejectsInvalidCall(t *testing.T) {
	_, err := CompileAdminCluster(nil, validAdminCluster())
	require.ErrorContains(t, err, "registry is nil")
	registry, err := NewBuiltinRegistry()
	require.NoError(t, err)
	_, err = CompileAdminCluster(registry, validRouteBindingObject())
	require.ErrorContains(t, err, `cannot compile kind "AdminRouteBinding"`)
}
