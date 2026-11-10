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
	"fmt"
	"io"
	"strings"
)

import (
	"gopkg.in/yaml.v3"
)

import (
	legacyconfig "github.com/apache/dubbo-go-pixiu/pkg/config"
)

// The enum matches the discovery types accepted by the current runtime.
var adminClusterTypes = []any{"Static", "StrictDNS", "LogicalDns", "EDS", "OriginalDst"}

func adminClusterSchema() ObjectSchema {
	return ObjectSchema{
		Kind:        KindAdminCluster,
		Description: "Cluster endpoint published to the existing Pixiu runtime",
		Fields: map[string]*FieldSchema{
			"type": {
				Type:     FieldTypeString,
				Required: true,
				Default:  "Static",
				Enum:     adminClusterTypes,
				UI:       UIHints{Component: "select", Order: 10},
			},
			"endpoint": {
				Type:     FieldTypeObject,
				Required: true,
				Properties: map[string]*FieldSchema{
					"address": {
						Type:     FieldTypeString,
						Required: true,
						UI:       UIHints{Component: "text", Order: 10, Placeholder: "127.0.0.1"},
					},
					"port": {
						Type:     FieldTypeInteger,
						Required: true,
						Minimum:  floatPointer(1),
						Maximum:  floatPointer(65535),
						UI:       UIHints{Component: "number", Order: 20},
					},
				},
				UI: UIHints{Group: "endpoint", Order: 20},
			},
		},
	}
}

func validateAdminCluster(object AdminObject) []ValidationIssue {
	endpoint, ok := object.Spec["endpoint"].(map[string]any)
	if !ok {
		return nil // The schema reports a missing or mistyped endpoint.
	}
	address, ok := endpoint["address"].(string)
	if ok && strings.TrimSpace(address) == "" {
		return []ValidationIssue{issue("spec.endpoint.address", "required", "address is required")}
	}
	return nil
}

// AdminClusterFromLegacy keeps the runtime ID out of editable form data.
// Invalid legacy values remain representable as drafts; compilation validates them.
func AdminClusterFromLegacy(cluster legacyconfig.Cluster) AdminObject {
	return AdminObject{
		Kind:     KindAdminCluster,
		Metadata: ObjectMetadata{Name: cluster.Name},
		Spec: map[string]any{
			"type": cluster.Type,
			"endpoint": map[string]any{
				"address": cluster.Address,
				"port":    cluster.Port,
			},
		},
	}
}

// DecodeLegacyClusterYAML returns the form object and its separate runtime ID.
// Unknown fields fail instead of being silently dropped during conversion.
func DecodeLegacyClusterYAML(data []byte) (AdminObject, int, error) {
	decoder := yaml.NewDecoder(strings.NewReader(string(data)))
	decoder.KnownFields(true)
	var cluster legacyconfig.Cluster
	if err := decoder.Decode(&cluster); err != nil {
		return AdminObject{}, 0, fmt.Errorf("decode legacy Cluster YAML: %w", err)
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		if err == nil {
			return AdminObject{}, 0, errors.New("legacy Cluster YAML must contain one document")
		}
		return AdminObject{}, 0, fmt.Errorf("decode legacy Cluster YAML: %w", err)
	}
	return AdminClusterFromLegacy(cluster), cluster.ID, nil
}

// CompiledCluster is a validated Admin object and its legacy runtime shape.
// ID is intentionally supplied by the storage layer, not by AdminObject.
type CompiledCluster struct {
	Source  AdminObject
	Cluster legacyconfig.Cluster
}

func CompileAdminCluster(registry *Registry, object AdminObject) (CompiledCluster, error) {
	if registry == nil {
		return CompiledCluster{}, errors.New("admin object schema registry is nil")
	}
	normalized, err := registry.Normalize(object)
	if err != nil {
		return CompiledCluster{}, err
	}
	if normalized.Kind != KindAdminCluster {
		return CompiledCluster{}, fmt.Errorf("cannot compile kind %q", normalized.Kind)
	}
	endpoint := normalized.Spec["endpoint"].(map[string]any)
	port, _ := numericValue(endpoint["port"])
	return CompiledCluster{
		Source: normalized,
		Cluster: legacyconfig.Cluster{
			Name:    normalized.Metadata.Name,
			Type:    normalized.Spec["type"].(string),
			Address: endpoint["address"].(string),
			Port:    int(port),
		},
	}, nil
}

// RuntimeCluster attaches the identity owned by the publisher.
func (c CompiledCluster) RuntimeCluster(id int) legacyconfig.Cluster {
	cluster := c.Cluster
	cluster.ID = id
	return cluster
}

// PreviewYAML shows the existing flat Cluster YAML fields. The runtime ID is
// omitted until publication assigns or reuses one.
func (c CompiledCluster) PreviewYAML() ([]byte, error) {
	return yaml.Marshal(struct {
		Name    string `yaml:"name"`
		Type    string `yaml:"type"`
		Address string `yaml:"address"`
		Port    int    `yaml:"port"`
	}{c.Cluster.Name, c.Cluster.Type, c.Cluster.Address, c.Cluster.Port})
}
