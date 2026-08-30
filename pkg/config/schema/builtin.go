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

// RegisterBuiltinSchemas installs only the Admin-facing route binding. Legacy
// Resource and Method are compiler outputs rather than form objects.
func RegisterBuiltinSchemas(registry *Registry) error {
	if registry == nil {
		return fmt.Errorf("admin object schema registry is nil")
	}
	if err := registry.Register(adminRouteBindingSchema()); err != nil {
		return err
	}
	return registry.RegisterValidator(KindAdminRouteBinding, validateAdminRouteBinding)
}

func adminRouteBindingSchema() ObjectSchema {
	return ObjectSchema{
		Kind:        KindAdminRouteBinding,
		Description: "High-level HTTP entry to Dubbo invocation binding",
		Fields: map[string]*FieldSchema{
			"entry": {
				Type:     FieldTypeObject,
				Required: true,
				Properties: map[string]*FieldSchema{
					"protocol": {
						Type:    FieldTypeString,
						Default: "http",
						Enum:    []any{"http"},
						UI:      UIHints{Component: "select", Order: 10},
					},
					"path": {
						Type:        FieldTypeString,
						Required:    true,
						Pattern:     `^/`,
						Description: "HTTP path pattern exposed by Pixiu",
						UI:          UIHints{Component: "text", Order: 20, Placeholder: "/api/v1/users/:id"},
					},
					"method": {
						Type:     FieldTypeString,
						Required: true,
						Enum:     []any{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"},
						UI:       UIHints{Component: "select", Order: 30},
					},
				},
				UI: UIHints{Group: "entry", Order: 10},
			},
			"target": {
				Type:     FieldTypeObject,
				Required: true,
				Properties: map[string]*FieldSchema{
					"protocol": {
						Type:    FieldTypeString,
						Default: "dubbo",
						Enum:    []any{"dubbo"},
						UI:      UIHints{Component: "select", Order: 10},
					},
					"application": {
						Type:        FieldTypeString,
						Required:    true,
						Description: "Dubbo applicationName",
						UI:          UIHints{Component: "text", Order: 20},
					},
					"interface": {
						Type:     FieldTypeString,
						Required: true,
						UI:       UIHints{Component: "text", Order: 30},
					},
					"method": {
						Type:     FieldTypeString,
						Required: true,
						UI:       UIHints{Component: "text", Order: 40},
					},
					"version": {Type: FieldTypeString, UI: UIHints{Component: "text", Order: 50}},
					"group":   {Type: FieldTypeString, UI: UIHints{Component: "text", Order: 60}},
					"cluster": {
						Type:        FieldTypeString,
						Required:    true,
						Description: "legacy integrationRequest.clusterName",
						UI:          UIHints{Component: "text", Order: 70},
					},
				},
				UI: UIHints{Group: "target", Order: 20},
			},
			"params": {
				Type:        FieldTypeArray,
				Default:     []any{},
				Description: "Ordered HTTP source to Dubbo argument mappings",
				Items: &FieldSchema{
					Type: FieldTypeObject,
					Properties: map[string]*FieldSchema{
						"from": {
							Type:        FieldTypeString,
							Required:    true,
							Pattern:     `^(uri|queryStrings|headers|requestBody)\..+`,
							Description: "legacy mappingParams.name",
						},
						"to": {
							Type:        FieldTypeInteger,
							Required:    true,
							Minimum:     floatPointer(0),
							Description: "zero-based Dubbo argument index",
						},
						"type": {
							Type:        FieldTypeString,
							Required:    true,
							Description: "legacy mapType and parameterTypes entry",
						},
					},
				},
				UI: UIHints{Component: "parameter-binding-table", Group: "params", Order: 30},
			},
			"publish": {
				Type:    FieldTypeObject,
				Default: map[string]any{},
				Properties: map[string]*FieldSchema{
					"mode": {
						Type:    FieldTypeString,
						Default: "draft",
						Enum:    []any{"draft", "published"},
						UI:      UIHints{Component: "select", Order: 10},
					},
					"validate": {
						Type:    FieldTypeBoolean,
						Default: true,
						UI:      UIHints{Component: "switch", Order: 20},
					},
				},
				UI: UIHints{Group: "publish", Order: 40, Advanced: true},
			},
			"extensions": {
				Type:    FieldTypeObject,
				Default: map[string]any{},
				UI:      UIHints{Component: "extension-fields", Group: "advanced", Order: 50, Advanced: true},
			},
		},
	}
}

func validateAdminRouteBinding(object AdminObject) []ValidationIssue {
	params, ok := object.Spec["params"].([]any)
	if !ok {
		return nil
	}

	issues := make([]ValidationIssue, 0)
	seenIndexes := make(map[int]int, len(params))
	for index, rawParam := range params {
		param, ok := rawParam.(map[string]any)
		if !ok {
			continue
		}
		toNumber, ok := numericValue(param["to"])
		if !ok {
			continue
		}
		to := int(toNumber)
		if previous, exists := seenIndexes[to]; exists {
			issues = append(issues, issue(
				fmt.Sprintf("spec.params[%d].to", index),
				"duplicate",
				fmt.Sprintf("duplicates params[%d] argument index %d", previous, to),
			))
		} else {
			seenIndexes[to] = index
		}

		if paramType, ok := param["type"].(string); ok && !isSupportedParamType(paramType) {
			issues = append(issues, issue(
				fmt.Sprintf("spec.params[%d].type", index),
				"unsupported",
				fmt.Sprintf("unsupported scalar Dubbo map type %q", paramType),
			))
		}
	}

	for expected := 0; expected < len(params); expected++ {
		if _, exists := seenIndexes[expected]; !exists {
			issues = append(issues, issue(
				"spec.params",
				"non_contiguous",
				fmt.Sprintf("argument indexes must be contiguous from 0; index %d is missing", expected),
			))
		}
	}
	return issues
}

func isSupportedParamType(value string) bool {
	switch strings.TrimSpace(value) {
	case "string", "char", "short", "int", "long", "float", "double", "boolean", "byte", "date", "object",
		"java.lang.String", "java.lang.Character", "java.lang.Short", "java.lang.Integer", "java.lang.Long",
		"java.lang.Float", "java.lang.Double", "java.lang.Boolean", "java.lang.Byte", "java.lang.Object", "java.util.Date":
		return true
	default:
		return false
	}
}

func floatPointer(value float64) *float64 {
	return &value
}
