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
	"fmt"
	"strings"
)

const (
	AdapterAPIResource = "pixiu-api-resource"
	AdapterAPIMethod   = "pixiu-api-method"
	AdapterXDSListener = "pixiu-xds-listener"
	AdapterXDSCluster  = "pixiu-xds-cluster"
)

// RegisterBuiltinSchemas installs the first schema version for the four
// objects currently managed by Admin. These definitions are deliberately
// independent from pkg/config structs; adapters are the compatibility
// boundary to the current runtime.
func RegisterBuiltinSchemas(registry *Registry) error {
	if registry == nil {
		return fmt.Errorf("config schema registry is nil")
	}

	for _, objectSchema := range []ObjectSchema{
		resourceSchema(),
		methodSchema(),
		listenerSchema(),
		clusterSchema(),
	} {
		if err := registry.Register(objectSchema); err != nil {
			return err
		}
	}
	if err := registry.RegisterValidator(APIVersionV1Alpha1, KindMethod, validateMethodSemantics); err != nil {
		return err
	}
	if err := registry.RegisterConfigSetValidator(validateBuiltinReferences); err != nil {
		return err
	}
	return nil
}

func resourceSchema() ObjectSchema {
	return ObjectSchema{
		APIVersion:     APIVersionV1Alpha1,
		Kind:           KindResource,
		Description:    "HTTP resource path and inherited request metadata",
		RuntimeAdapter: AdapterAPIResource,
		Fields: map[string]*FieldSchema{
			"type": {
				Type:        FieldTypeString,
				Default:     "restful",
				Description: "resource protocol family",
				UI:          UIHints{Component: "select", Group: "basic", Order: 10},
			},
			"path": {
				Type:        FieldTypeString,
				Required:    true,
				Pattern:     `^/`,
				Description: "HTTP path prefix",
				UI:          UIHints{Component: "text", Group: "basic", Order: 20},
			},
			"timeout": {
				Type:        FieldTypeString,
				Default:     "1s",
				Description: "resource-level timeout retained for legacy compatibility",
				UI:          UIHints{Component: "duration", Group: "basic", Order: 30},
			},
			"description": stringField(false),
			"headers": {
				Type:                 FieldTypeMap,
				Default:              map[string]any{},
				AdditionalProperties: stringField(false),
				UI:                   UIHints{Component: "key-value", Group: "request", Order: 40},
			},
			"methodRefs": {
				Type:    FieldTypeArray,
				Default: []any{},
				Items:   stringField(false),
				UI:      UIHints{Component: "object-reference-list", Group: "methods", Order: 50},
			},
			"filters":    openObjectArray("filters", 60),
			"extensions": extensionsField(),
		},
	}
}

func methodSchema() ObjectSchema {
	parameterSchema := &FieldSchema{
		Type: FieldTypeObject,
		Properties: map[string]*FieldSchema{
			"in": {
				Type:     FieldTypeString,
				Required: true,
				Enum:     []any{"path", "query", "header", "body"},
				UI:       UIHints{Component: "select"},
			},
			"required": {Type: FieldTypeBoolean, Default: false},
			"schema": {
				Type:         FieldTypeObject,
				Required:     true,
				AllowUnknown: true,
				UI:           UIHints{Component: "json-schema"},
			},
		},
	}

	signatureItem := &FieldSchema{
		Type: FieldTypeObject,
		Properties: map[string]*FieldSchema{
			"index":    integerField(true, floatPointer(0), nil),
			"name":     stringField(false),
			"javaType": stringField(true),
		},
	}

	bindingItem := &FieldSchema{
		Type: FieldTypeObject,
		Properties: map[string]*FieldSchema{
			"from": stringField(true),
			"to":   stringField(true),
			"type": stringField(false),
		},
	}

	return ObjectSchema{
		APIVersion:     APIVersionV1Alpha1,
		Kind:           KindMethod,
		Description:    "HTTP operation, request schema and backend invocation",
		RuntimeAdapter: AdapterAPIMethod,
		Fields: map[string]*FieldSchema{
			"resourceRef": {
				Type:        FieldTypeString,
				Required:    true,
				Description: "owner Resource metadata.id",
				UI:          UIHints{Component: "object-reference", Group: "basic", Order: 10},
			},
			"enable":  {Type: FieldTypeBoolean, Default: true, UI: UIHints{Component: "switch", Group: "basic", Order: 20}},
			"timeout": {Type: FieldTypeString, Default: "1s", UI: UIHints{Component: "duration", Group: "basic", Order: 30}},
			"mock":    {Type: FieldTypeBoolean, Default: false, UI: UIHints{Component: "switch", Group: "basic", Order: 40}},
			"http": {
				Type:     FieldTypeObject,
				Required: true,
				Properties: map[string]*FieldSchema{
					"verb": {
						Type:     FieldTypeString,
						Required: true,
						Enum:     []any{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
						UI:       UIHints{Component: "select"},
					},
				},
				UI: UIHints{Group: "basic", Order: 50},
			},
			"request": {
				Type:    FieldTypeObject,
				Default: map[string]any{},
				Properties: map[string]*FieldSchema{
					"parameters": {
						Type:                 FieldTypeMap,
						Default:              map[string]any{},
						AdditionalProperties: parameterSchema,
						UI:                   UIHints{Component: "parameter-table"},
					},
					"body": {
						Type:         FieldTypeObject,
						AllowUnknown: true,
						UI:           UIHints{Component: "json-schema"},
					},
				},
				UI: UIHints{Group: "request", Order: 60},
			},
			"backend": {
				Type:          FieldTypeUnion,
				Required:      true,
				Discriminator: "kind",
				Variants: map[string]*FieldSchema{
					"dubbo": dubboBackendSchema(signatureItem),
					"http":  httpBackendSchema(),
				},
				UI: UIHints{Component: "discriminated-form", Group: "backend", Order: 70},
			},
			"bindings": {
				Type:    FieldTypeArray,
				Default: []any{},
				Items:   bindingItem,
				UI:      UIHints{Component: "binding-table", Group: "mapping", Order: 80},
			},
			"filters":    openObjectArray("filters", 90),
			"extensions": extensionsField(),
		},
	}
}

func listenerSchema() ObjectSchema {
	return ObjectSchema{
		APIVersion:     APIVersionV1Alpha1,
		Kind:           KindListener,
		Description:    "gateway ingress address, routes and HTTP filters",
		RuntimeAdapter: AdapterXDSListener,
		Fields: map[string]*FieldSchema{
			"address": {
				Type:     FieldTypeObject,
				Required: true,
				Properties: map[string]*FieldSchema{
					"host": {Type: FieldTypeString, Default: "0.0.0.0"},
					"port": integerField(true, floatPointer(1), floatPointer(65535)),
				},
				UI: UIHints{Group: "basic", Order: 10},
			},
			"protocol": {
				Type:    FieldTypeString,
				Default: "HTTP",
				Enum:    []any{"HTTP", "HTTPS", "HTTP2", "GRPC", "TRIPLE", "TCP", "UDP"},
				UI:      UIHints{Component: "select", Group: "basic", Order: 20},
			},
			"routes": {
				Type:    FieldTypeArray,
				Default: []any{},
				Items: &FieldSchema{
					Type: FieldTypeObject,
					Properties: map[string]*FieldSchema{
						"prefix":     {Type: FieldTypeString, Required: true, Pattern: `^/`},
						"clusterRef": stringField(true),
					},
				},
				UI: UIHints{Component: "route-table", Group: "routes", Order: 30},
			},
			"filters":    openObjectArray("filters", 40),
			"extensions": extensionsField(),
		},
	}
}

func clusterSchema() ObjectSchema {
	return ObjectSchema{
		APIVersion:     APIVersionV1Alpha1,
		Kind:           KindCluster,
		Description:    "upstream discovery and endpoint configuration",
		RuntimeAdapter: AdapterXDSCluster,
		Fields: map[string]*FieldSchema{
			"type": {
				Type:    FieldTypeString,
				Default: "Static",
				Enum:    []any{"Static", "StrictDNS", "LogicalDns", "EDS", "OriginalDst"},
				UI:      UIHints{Component: "select", Group: "basic", Order: 10},
			},
			"endpoints": {
				Type:    FieldTypeArray,
				Default: []any{},
				Items: &FieldSchema{
					Type: FieldTypeObject,
					Properties: map[string]*FieldSchema{
						"id":      stringField(false),
						"address": stringField(true),
						"port":    integerField(true, floatPointer(1), floatPointer(65535)),
					},
				},
				UI: UIHints{Component: "endpoint-table", Group: "endpoints", Order: 20},
			},
			"loadBalancer": {
				Type:    FieldTypeString,
				Default: "roundrobin",
				UI:      UIHints{Component: "select", Group: "policy", Order: 30},
			},
			"healthChecks": openObjectArray("healthChecks", 40),
			"extensions":   extensionsField(),
		},
	}
}

func dubboBackendSchema(signatureItem *FieldSchema) *FieldSchema {
	return &FieldSchema{
		Type: FieldTypeObject,
		Properties: map[string]*FieldSchema{
			"kind":            {Type: FieldTypeString, Required: true, Enum: []any{"dubbo"}},
			"clusterRef":      stringField(false),
			"applicationName": stringField(false),
			"interface":       stringField(true),
			"method":          stringField(true),
			"group":           stringField(false),
			"version":         stringField(false),
			"protocol":        {Type: FieldTypeString, Default: "dubbo"},
			"serialization":   stringField(false),
			"retries":         integerField(false, floatPointer(0), nil),
			"signature": {
				Type:    FieldTypeArray,
				Default: []any{},
				Items:   signatureItem,
				UI:      UIHints{Component: "signature-table"},
			},
			"extensions": extensionsField(),
		},
	}
}

func httpBackendSchema() *FieldSchema {
	return &FieldSchema{
		Type: FieldTypeObject,
		Properties: map[string]*FieldSchema{
			"kind":       {Type: FieldTypeString, Required: true, Enum: []any{"http"}},
			"clusterRef": stringField(false),
			"url":        stringField(false),
			"host":       stringField(false),
			"path":       stringField(false),
			"scheme":     {Type: FieldTypeString, Enum: []any{"http", "https"}},
			"extensions": extensionsField(),
		},
	}
}

func validateMethodSemantics(object ConfigObject) []ValidationIssue {
	backend, ok := object.Spec["backend"].(map[string]any)
	if !ok {
		return nil
	}

	issues := make([]ValidationIssue, 0)
	if backend["kind"] == "http" && allEmptyStrings(backend, "clusterRef", "url", "host") {
		issues = append(issues, issue(
			"spec.backend",
			"target_required",
			"HTTP backend requires clusterRef, url, or host",
		))
	}

	items, _ := backend["signature"].([]any)
	seenIndexes := make(map[int]int, len(items))
	for i, item := range items {
		parameter, ok := item.(map[string]any)
		if !ok {
			continue
		}
		indexValue, ok := numericValue(parameter["index"])
		if !ok {
			continue
		}
		index := int(indexValue)
		if previous, exists := seenIndexes[index]; exists {
			issues = append(issues, issue(
				fmt.Sprintf("spec.backend.signature[%d].index", i),
				"duplicate",
				fmt.Sprintf("duplicates signature[%d] index %d", previous, index),
			))
		} else {
			seenIndexes[index] = i
		}
	}
	return issues
}

func validateBuiltinReferences(configSet ConfigSet) []ValidationIssue {
	identities := make(map[string]struct{}, len(configSet.Objects))
	for _, object := range configSet.Objects {
		if object.Metadata.ID != "" {
			identities[object.Kind+"/"+object.Metadata.ID] = struct{}{}
		}
	}

	issues := make([]ValidationIssue, 0)
	for i, object := range configSet.Objects {
		prefix := fmt.Sprintf("objects[%d].", i)
		for j, reference := range object.Metadata.OwnerReferences {
			if reference.Kind == "" || reference.ID == "" {
				continue
			}
			if _, exists := identities[reference.Kind+"/"+reference.ID]; !exists {
				issues = append(issues, missingReferenceIssue(
					fmt.Sprintf("%smetadata.ownerReferences[%d].id", prefix, j),
					reference.Kind,
					reference.ID,
				))
			}
		}

		switch object.Kind {
		case KindResource:
			methodRefs, _ := object.Spec["methodRefs"].([]any)
			for j, rawReference := range methodRefs {
				reference, ok := rawReference.(string)
				if ok && reference != "" {
					issues = appendMissingReference(
						issues,
						identities,
						fmt.Sprintf("%sspec.methodRefs[%d]", prefix, j),
						KindMethod,
						reference,
					)
				}
			}
		case KindMethod:
			if reference, ok := object.Spec["resourceRef"].(string); ok && reference != "" {
				issues = appendMissingReference(issues, identities, prefix+"spec.resourceRef", KindResource, reference)
			}
			backend, _ := object.Spec["backend"].(map[string]any)
			if reference, ok := backend["clusterRef"].(string); ok && reference != "" {
				issues = appendMissingReference(issues, identities, prefix+"spec.backend.clusterRef", KindCluster, reference)
			}
		case KindListener:
			routes, _ := object.Spec["routes"].([]any)
			for j, rawRoute := range routes {
				route, _ := rawRoute.(map[string]any)
				if reference, ok := route["clusterRef"].(string); ok && reference != "" {
					issues = appendMissingReference(
						issues,
						identities,
						fmt.Sprintf("%sspec.routes[%d].clusterRef", prefix, j),
						KindCluster,
						reference,
					)
				}
			}
		}
	}
	return issues
}

func appendMissingReference(
	issues []ValidationIssue,
	identities map[string]struct{},
	path, kind, id string,
) []ValidationIssue {
	if _, exists := identities[kind+"/"+id]; exists {
		return issues
	}
	return append(issues, missingReferenceIssue(path, kind, id))
}

func missingReferenceIssue(path, kind, id string) ValidationIssue {
	return issue(path, "reference_not_found", fmt.Sprintf("referenced %s %q does not exist", kind, id))
}

func stringField(required bool) *FieldSchema {
	return &FieldSchema{Type: FieldTypeString, Required: required}
}

func integerField(required bool, minimum, maximum *float64) *FieldSchema {
	return &FieldSchema{Type: FieldTypeInteger, Required: required, Minimum: minimum, Maximum: maximum}
}

func extensionsField() *FieldSchema {
	return &FieldSchema{
		Type:         FieldTypeObject,
		Default:      map[string]any{},
		AllowUnknown: true,
		UI:           UIHints{Component: "extension-fields", Group: "advanced", Advanced: true},
	}
}

func openObjectArray(group string, order int) *FieldSchema {
	return &FieldSchema{
		Type:    FieldTypeArray,
		Default: []any{},
		Items:   &FieldSchema{Type: FieldTypeObject, AllowUnknown: true},
		UI:      UIHints{Component: "schema-list", Group: group, Order: order, Advanced: true},
	}
}

func floatPointer(value float64) *float64 {
	return &value
}

func allEmptyStrings(values map[string]any, keys ...string) bool {
	for _, key := range keys {
		if value, ok := values[key].(string); ok && strings.TrimSpace(value) != "" {
			return false
		}
	}
	return true
}
